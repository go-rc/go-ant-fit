package ant_fit

import (
    "errors"
    "fmt"
    "io"
    //"sort"
)

// general utility functions

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

func get_type_name(fld *FitFieldDefinition) string {
    if fld.base_type >= 0 &&
        int(fld.base_type) < len(base_type_names) {
        return base_type_names[fld.base_type]
    }

    if fld.num == 253 && fld.base_type == 6 {
        return "timestamp"
    } else if fld.num == 254 {
        return "message_index"
    }

    return fmt.Sprintf("unknown#%d", fld.num)
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

    goodCRC, _ := get_uint16_pos(buf, 0)
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
