package main

import (
    "encoding/binary"
    "bufio"
    "bytes"
    "errors"
    "flag"
    "fmt"
    "io"
    "os"
)

type FitData struct {
    proto byte
    profile uint16
    datasize uint32
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
    } else if n != 2 {
        errfmt := "Tried to read 2 byte CRC, only read %d bytes"
        return errors.New(fmt.Sprintf(errfmt, n))
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

func readHeader(rdr io.Reader) (*FitData, error) {
    const minHeaderLen byte = 12

    buf := make([]byte, minHeaderLen)

    n, err := rdr.Read(buf)
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
        err = checkCRC(rdr, buf)
        if err != nil {
            return nil, err
        }
    }

    fdata := new(FitData)

    fdata.proto = buf[1]
    fdata.profile = to_uint16(buf[2:4])
    fdata.datasize = to_uint32(buf[4:9])

    return fdata, nil
}

func readFit(filename string) {
    file, err := os.Open(filename)
    if err != nil {
        fmt.Printf("Cannot open \"%s\"\n", filename)
        return
    }
    defer file.Close()

    rdr := bufio.NewReader(file)

    fdata, err := readHeader(rdr)
    if err != nil {
        fmt.Printf("Cannot read %s: %s\n", filename, err)
        return
    }

    fmt.Printf("%s: proto %d profile %d data %d\n", filename, fdata.proto,
        fdata.profile, fdata.datasize)
}

func to_uint16(data []byte) (ret uint16) {
    buf := bytes.NewBuffer(data)
    binary.Read(buf, binary.LittleEndian, &ret)
    return
}

func to_uint32(data []byte) (ret uint32) {
    buf := bytes.NewBuffer(data)
    binary.Read(buf, binary.LittleEndian, &ret)
    return
}

func main() {
    verbose, files := processArgs()

    if verbose {
        fmt.Println("Verbose mode")
    }

    for _, f := range files {
        readFit(f)
    }
}
