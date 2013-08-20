package main

import (
    "bufio"
    "bytes"
    "encoding/binary"
    "errors"
    "flag"
    "fmt"
    "io"
    "os"
    "sort"
)

var base_type_names = [14]string{
    "enum",
    "int8",
    "uint8",
    "int16",
    "uint16",
    "int32",
    "uint32",
    "string",
    "float32",
    "float64",
    "uint8z",
    "uint16z",
    "uint32z",
    "byte",
}

func addCRC(crc uint16, val byte) uint16 {
    lookup := [16]uint16{
        0x0000, 0xcc01, 0xd801, 0x1400, 0xf001, 0x3c00, 0x2800, 0xe401,
        0xa001, 0x6c00, 0x7800, 0xb401, 0x5000, 0x9c01, 0x8801, 0x4400,
    }

    var tmp uint16

    // compute checksum of lower four bits of value
    tmp = lookup[crc & 0xf]
    crc = uint16((crc >> 4) & 0xfff)
    crc = uint16(crc ^ tmp ^ lookup[val & 0xf])

    // compute checksum of upper four bits of value
    tmp = lookup[crc & 0xf]
    crc = uint16((crc >> 4) & 0xfff)
    crc = uint16(crc ^ tmp ^ lookup[(val >> 4) & 0xf])

    return crc
}

func checkCRC(rdr io.Reader, data []byte) error {
    buf := make([]byte, 2)

    n, err := rdr.Read(buf)
    if err != nil {
        return err
    } else if n != len(buf) {
        errfmt := "Tried to read %d byte CRC, only read %d bytes"
        return errors.New(fmt.Sprintf(errfmt, len(buf), n))
    }

    goodCRC := to_uint16(buf)
    if goodCRC == 0 {
        // CRC is not set, so we're done
        return nil
    }

    var crc uint16
    for i := 0; i < len(data); i++ {
        crc = addCRC(crc, data[i])
    }

    if goodCRC != crc {
        errfmt := "Bad header CRC: %04x != %04x"
        return errors.New(fmt.Sprintf(errfmt, crc, goodCRC))
    }

    return nil
}

type FitFile struct {
    filename string
    rdr io.Reader

    proto byte
    profile uint16
    datasize uint32

    defs []*FitDefinition
    data []FitMsg
}

type FitFieldDefinition struct {
    num byte
    size byte
    is_endian bool
    base_type byte
}

type FitDefinition struct {
    local_type byte
    little_endian bool
    global_num uint16
    fields []*FitFieldDefinition
    total_bytes uint16
}

// interface and methods used to sort field definition list
type Definitions []*FitFieldDefinition

func (s Definitions) Len() int {
    return len(s)
}

func (s Definitions) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}

type ByNum struct { Definitions }

func (s ByNum) Less(i, j int) bool {
    return s.Definitions[i].num < s.Definitions[j].num
}

// message interface

type FitMsg interface {
    name() string
    text() string
}

// file_id message

type MsgFileId struct {
    msgtype byte
    manufacturer uint16
    product uint16
    serial_number uint32
    time_created uint32
    number uint16
}

func (msg *MsgFileId) msgtype_name() string {
    switch msg.msgtype {
    case 1: return "device";
    case 2: return "settings";
    case 3: return "sport";
    case 4: return "activity";
    case 5: return "workout";
    case 6: return "course";
    case 7: return "schedules";
    case 9: return "weight";
    case 10: return "totals";
    case 11: return "goals";
    case 14: return "blood_pressure";
    case 15: return "monitoring";
    case 20: return "activity_summary";
    case 28: return "monitoring_daily";
    default: return fmt.Sprintf("invalid#%d", msg.msgtype);

    }
}

func (msg *MsgFileId) name() string {
    return "file_id"
}

func (msg *MsgFileId) text() string {
    return fmt.Sprintf("file_id #%d msgtype %s mfct %d prod %d ser# %d time %d",
        msg.number, msg.msgtype_name(), msg.manufacturer, msg.product,
        msg.serial_number, msg.time_created)
}

func NewMsgFileId(data []byte) (*MsgFileId, error) {
    const explen int = 15

    if len(data) != explen {
        errfmt := "FileId message should be %d bytes, not %d"
        return nil, errors.New(fmt.Sprintf(errfmt, explen, len(data)))
    }

    msg := new(MsgFileId)

    msg.msgtype = data[0]
    msg.manufacturer = to_uint16(data[1:3])
    msg.product = to_uint16(data[3:5])
    msg.serial_number = to_uint32(data[5:9])
    msg.time_created = to_uint32(data[9:13])
    msg.number = to_uint16(data[13:])

    return msg, nil
}

// event message

type MsgEvent struct {
    timestamp uint32
    event byte
    event_type byte
}

func (msg *MsgEvent) event_name() string {
    switch msg.event {
    case 0: return "timer"
    case 3: return "timer"
    case 4: return "timer"
    case 5: return "timer"
    case 6: return "timer"
    case 7: return "timer"
    case 8: return "timer"
    case 9: return "timer"
    case 10: return "timer"
    case 11: return "timer"
    case 12: return "timer"
    case 13: return "timer"
    case 14: return "timer"
    case 15: return "timer"
    case 16: return "timer"
    case 17: return "timer"
    case 18: return "timer"
    case 19: return "timer"
    case 20: return "timer"
    case 21: return "timer"
    case 22: return "timer"
    case 23: return "timer"
    case 24: return "timer"
    case 25: return "timer"
    case 26: return "timer"
    case 27: return "timer"
    case 28: return "timer"
    case 36: return "timer"
    default: return fmt.Sprintf("unknown#%d", msg.event)
    }
}

func (msg *MsgEvent) event_type_name() string {
    switch msg.event_type {
    case 0: return "start"
    case 1: return "stop"
    case 2: return "consecutive_deprecated"
    case 3: return "marker"
    case 4: return "stop_all"
    case 5: return "begin_deprecated"
    case 6: return "end_deprecated"
    case 7: return "end_all_deprecated"
    case 8: return "stop_disable"
    case 9: return "stop_disable_all"
    default: return fmt.Sprintf("unknown#%d", msg.event_type)
    }
}

func (msg *MsgEvent) name() string {
    return "event"
}

func (msg *MsgEvent) text() string {
    return fmt.Sprintf("event tstmp %d evt %s etyp %s", msg.timestamp,
        msg.event_name(), msg.event_type_name())
}

func NewMsgEvent(data []byte) (*MsgEvent, error) {
    const minlen int = 6

    if len(data) != minlen {
        errfmt := "Event message should be at least %d bytes, not %d"
        return nil, errors.New(fmt.Sprintf(errfmt, minlen, len(data)))
    }

    msg := new(MsgEvent)

    msg.timestamp = to_uint32(data[0:5])
    msg.event = data[5]
    msg.event_type = data[6]

    return msg, nil
}

// software message

type MsgSoftware struct {
    message_index uint16
    version uint16
    part_number string
}

func (msg *MsgSoftware) name() string {
    return "software"
}

func (msg *MsgSoftware) text() string {
    return fmt.Sprintf("software msgidx %d vers %d part# %d", msg.message_index,
        msg.version, msg.part_number)
}

func NewMsgSoftware(data []byte) (*MsgSoftware, error) {
    const minlen int = 5

    if len(data) < minlen {
        errfmt := "Software message should be at least %d bytes, not %d"
        return nil, errors.New(fmt.Sprintf(errfmt, minlen, len(data)))
    }

    msg := new(MsgSoftware)

    msg.message_index = to_uint16(data[0:2])
    msg.version = to_uint16(data[2:4])
    msg.part_number = to_string(data[4:])

    return msg, nil
}

// file_creator message

type MsgFileCreator struct {
    software_version uint16
    hardware_version byte
}

func (msg *MsgFileCreator) name() string {
    return "file_creator"
}

func (msg *MsgFileCreator) text() string {
    return fmt.Sprintf("file_creator soft %d hard %d", msg.software_version,
        msg.hardware_version)
}

func NewMsgFileCreator(data []byte) (*MsgFileCreator, error) {
    const explen int = 3

    if len(data) < explen {
        errfmt := "FileId message should be %d bytes, not %d"
        return nil, errors.New(fmt.Sprintf(errfmt, explen, len(data)))
    }

    msg := new(MsgFileCreator)

    msg.software_version = to_uint16(data[0:2])
    msg.hardware_version = data[2]

    return msg, nil
}

// unknown message

type MsgUnknown struct {
    global_num uint16
    data []byte
}

func (msg *MsgUnknown) name() string {
    return fmt.Sprintf("unknown#%d", msg.global_num)
}

func (msg *MsgUnknown) text() string {
    return fmt.Sprintf("unknown#%d", msg.global_num)
}

func NewMsgUnknown(global_num uint16, data []byte) (*MsgUnknown, error) {
    msg := new(MsgUnknown)

    msg.global_num = global_num
    msg.data = make([]byte, len(data))
    copy(msg.data, data)

    return msg, nil
}

func (ffile *FitFile) open(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return errors.New(fmt.Sprintf("Cannot open \"%s\"\n", filename))
    }
    defer file.Close()

    ffile.filename = filename
    ffile.rdr = bufio.NewReader(file)

    const minHeaderLen byte = 12

    buf := make([]byte, minHeaderLen)

    n, err := ffile.rdr.Read(buf)
    if err != nil {
        return err
    } else if n != int(minHeaderLen) {
        errfmt := "Tried to read %d byte header, only read %d bytes"
        return errors.New(fmt.Sprintf(errfmt, minHeaderLen, n))
    }

    size := buf[0]

    needCRC := size == minHeaderLen + 2
    if size != minHeaderLen && !needCRC {
        return errors.New(fmt.Sprintf("Unexpected header size %d", size))
    }

    // verify that the ASCII signature is correct
    signature := []byte(".FIT")
    for i := 0; i < 4; i++ {
        if buf[8 + i] != signature[i] {
            errfmt := "Bad signature char #%d: '%c' should be '%c'\n"
            return errors.New(fmt.Sprintf(errfmt, i, buf[8 + i],
                signature[i]))
        }
    }

    // verify that the CRC is correct (if present)
    if needCRC {
        err = checkCRC(ffile.rdr, buf)
        if err != nil {
            return err
        }
    }

    ffile.proto = buf[1]
    ffile.profile = to_uint16(buf[2:4])
    ffile.datasize = to_uint32(buf[4:9])

    ffile.defs = make([]*FitDefinition, 0)
    ffile.data = make([]FitMsg, 0)
    return nil
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
    case 0: return NewMsgFileId(buf)
    //case 21: return NewMsgEvent(buf)
    case 35: return NewMsgSoftware(buf)
    case 49: return NewMsgFileCreator(buf)
    default: return NewMsgUnknown(def.global_num, buf)
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
    def.global_num = to_uint16(buf[2:4])
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
    sort.Sort(ByNum{def.fields})

    if verbose {
        fmt.Printf("  def: ltyp %v little_endian %v glbl %d\n",
            def.local_type, def.little_endian, def.global_num)
        for i := 0; i < len(def.fields); i++ {
            var type_name string

            if def.fields[i].base_type >= 0 &&
                int(def.fields[i].base_type) < len(base_type_names) {
                type_name = base_type_names[def.fields[i].base_type]
            }

            if def.fields[i].num == 253 && def.fields[i].base_type == 6 {
                type_name = "timestamp"
            } else if def.fields[i].num == 254 {
                type_name = "message_index"
            }

            fmt.Printf("       :: num %d sz %d endian %v type %s\n",
                def.fields[i].num, def.fields[i].size, def.fields[i].is_endian,
                type_name)
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

func (ffile *FitFile) readMessage(verbose bool) (bool, error) {
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
            fmt.Println("  data:", data.text())
        }

        ffile.data = append(ffile.data, data)
    }

    return true, nil
}

func readFit(filename string, verbose bool) error {
    ffile := new(FitFile)

    err := ffile.open(filename)
    if err != nil {
        return err
    }

    if verbose {
        fmt.Printf("%s: proto %d profile %d data %d\n", ffile.filename,
            ffile.proto, ffile.profile, ffile.datasize)
    }

    for true {
        flag, err := ffile.readMessage(verbose)
        if err != nil {
            return err
        } else if !flag {
            break
        }
    }

    return nil
}

func to_uint16(data []byte) (ret uint16) {
    return to_uint16_endian(data, binary.LittleEndian)
}

func to_uint16_endian(data []byte, order binary.ByteOrder) (ret uint16) {
    buf := bytes.NewBuffer(data)
    binary.Read(buf, order, &ret)
    return
}

func to_string(data[] byte) string {
    if data[len(data) - 1] != 0 {
        return string(data)
    }

    return string(data[:len(data)-1])
}

func to_uint32(data []byte) (ret uint32) {
    return to_uint32_endian(data, binary.LittleEndian)
}

func to_uint32_endian(data []byte, order binary.ByteOrder) (ret uint32) {
    buf := bytes.NewBuffer(data)
    binary.Read(buf, order, &ret)
    return
}

func processArgs() (bool, []string) {
    usage := false

    verbosep := flag.Bool("verbose", false, "Verbose mode")

    flag.Parse()

    files := make([]string, 0)
    for _, f := range flag.Args() {
        if _, err := os.Stat(f); os.IsNotExist(err) {
            fmt.Println("File ", f, " does not exist")
            usage = true
        } else {
            files = append(files, f)
        }
    }

    if usage {
        fmt.Print("Usage: readfit.go")
        fmt.Print("[-verbose]")
        fmt.Print("file [file ...]")
        fmt.Println()

        os.Exit(1)
    }

    return *verbosep, files
}

func main() {
    verbose, files := processArgs()

    for _, f := range files {
        err := readFit(f, verbose)
        if err != nil {
            fmt.Printf("!! Cannot read %s: %s\n", f, err)
        }
    }
}
