package ant_fit

import (
    "bufio"
    "errors"
    "fmt"
    "io"
    "os"
    //"sort"
)

type FitFile struct {
    filename string
    rdr io.Reader

    proto byte
    profile uint16
    datasize uint32

    defs []*FitDefinition
    data []FitMsg
}

func NewFitFile(filename string) (*FitFile, error) {
    file, err := os.Open(filename)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("Cannot open \"%s\"\n", filename))
    }
    defer file.Close()

    ffile := new(FitFile)

    ffile.filename = filename
    ffile.rdr = bufio.NewReader(file)

    const minHeaderLen byte = 12

    buf := make([]byte, minHeaderLen)

    n, err := ffile.rdr.Read(buf)
    if err != nil {
        return nil, err
    } else if n != int(minHeaderLen) {
        errfmt := "Tried to read %d byte header, only read %d bytes"
        return nil, errors.New(fmt.Sprintf(errfmt, minHeaderLen, n))
    }

    size := buf[0]

    needCRC := size == minHeaderLen + 2
    if size != minHeaderLen && !needCRC {
        return nil, errors.New(fmt.Sprintf("Unexpected header size %d", size))
    }

    // verify that the ASCII signature is correct
    signature := []byte(".FIT")
    for i := 0; i < 4; i++ {
        if buf[8 + i] != signature[i] {
            errfmt := "Bad signature char #%d: '%c' should be '%c'\n"
            return nil, errors.New(fmt.Sprintf(errfmt, i, buf[8 + i],
                signature[i]))
        }
    }

    // verify that the CRC is correct (if present)
    if needCRC {
        err = checkCRC(ffile.rdr, buf)
        if err != nil {
            return nil, err
        }
    }

    ffile.proto = buf[1]
    ffile.profile, _ = get_uint16_pos(buf, 2)
    ffile.datasize, _ = get_uint32_pos(buf, 4)

    ffile.defs = make([]*FitDefinition, 0)
    ffile.data = make([]FitMsg, 0)
    return ffile, nil
}

func (ffile *FitFile) findDefinition(local_type byte) (*FitDefinition, error) {
    var def *FitDefinition
    for i := 0; i < len(ffile.defs); i++ {
        if ffile.defs[i].local_type == local_type {
            def = ffile.defs[i]
        }
    }

    var err error
    if def == nil {
        err = errors.New(fmt.Sprintf("Unknown local_type %d", local_type))
    }

    return def, err
}

func (ffile *FitFile) readData(def *FitDefinition,
    time_offset uint32, verbose bool) (FitMsg, error) {

    buf := make([]byte, def.total_bytes)

    n, err := ffile.rdr.Read(buf)
    if err != nil {
        return nil, err
    } else if n != len(buf) {
        return nil, errors.New(fmt.Sprintf("Read %d bytes, not %d", n,
            len(buf)))
    }

    switch def.global_num {
    case 0: return NewMsgFileId(def, buf)
    case 21: return NewMsgEvent(def, buf)
    case 23: return NewMsgDeviceInfo(def, buf)
    case 35: return NewMsgSoftware(def, buf)
    case 49: return NewMsgFileCreator(def, buf)
    default: return NewMsgUnknown(def, buf, def.global_num)
    }
}

func (ffile *FitFile) readDefinition(local_type byte,
    verbose bool) (*FitDefinition, error) {
    buf := make([]byte, 5)

    n, err := ffile.rdr.Read(buf)
    if err != nil {
        return nil, err
    } else if n != len(buf) {
        return nil, errors.New(fmt.Sprintf("Read %d bytes, not %d", n,
            len(buf)))
    }

    def := new(FitDefinition)

    def.local_type = local_type
    def.little_endian = buf[1] == 0
    def.global_num, _ = get_uint16_pos(buf, 2)
    def.total_bytes = 0

    num := int(buf[4])

    def.fields = make([]*FitFieldDefinition, num)
    for i := 0; i < num; i++ {
        def.fields[i], err = ffile.readFieldDef(buf)
        if err != nil {
            return nil, err
        }
        def.total_bytes += uint16(def.fields[i].size)
    }
    //sort.Sort(ByNum{def.fields})

    if verbose {
        fmt.Printf("  def: ltyp %v little_endian %v glbl %d\n",
            def.local_type, def.little_endian, def.global_num)
        for i := 0; i < len(def.fields); i++ {
            fmt.Printf("       :: num %d sz %d endian %v type %s\n",
                def.fields[i].num, def.fields[i].size, def.fields[i].is_endian,
                get_type_name(def.fields[i]))
        }
    }

    return def, nil
}

func (ffile *FitFile) readFieldDef(buf []byte) (*FitFieldDefinition, error) {
    n, err := ffile.rdr.Read(buf[:3])
    if err != nil {
        return nil, err
    } else if n != 3 {
        return nil, errors.New(fmt.Sprintf("Read %d bytes, not %d", n, 3))
    }

    fld := new(FitFieldDefinition)

    fld.num = buf[0]
    fld.size = buf[1]
    fld.is_endian = buf[2] & 0x80 == 0x80
    fld.base_type = buf[2] & 0xf

    return fld, nil
}

func (ffile *FitFile) ReadMessage(verbose bool) (bool, error) {
    buf := make([]byte, 1)

    n, err := ffile.rdr.Read(buf)
    if err != nil {
        return false, err
    } else if n != len(buf) {
        return false, errors.New(fmt.Sprintf("Read %d bytes, not %d", n,
            len(buf)))
    }

    var is_def bool
    var local_type byte
    var time_offset uint32

    compressed := buf[0] & 0x80 == 0x80
    if !compressed {
        is_def = buf[0] & 0x40 == 0x40
        local_type = buf[0] & 0x0f
        time_offset = 0
    } else {
        is_def = false
        local_type = (buf[0] >> 4) & 0x3
        time_offset = uint32(buf[0] & 0xf)
    }

    if is_def {
        def, derr := ffile.readDefinition(local_type, verbose)
        if derr != nil {
            return false, derr
        }

        ffile.defs = append(ffile.defs, def)
    } else {
        def, err2 := ffile.findDefinition(local_type)
        if err2 != nil {
            return false, err2
        }

        data, err3 := ffile.readData(def, time_offset, verbose)
        if err3 != nil {
            return false, err3
        }

        if verbose {
            fmt.Println("  data:", data.Text())
        }

        ffile.data = append(ffile.data, data)
    }

    return true, nil
}

func (ffile *FitFile) String() string {
    return fmt.Sprintf("%s: proto %d profile %d data %d", ffile.filename,
        ffile.proto, ffile.profile, ffile.datasize)
}
