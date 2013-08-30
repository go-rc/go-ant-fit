package ant_fit

import (
    "bytes"
    "encoding/binary"
)

// data buffer extraction functions

func get_string_pos(data []byte, pos int) (string, int) {
    n := pos
    for n < len(data) {
        if data[n] == 0 {
            n++
            break
        }
        n++
    }

    return string(data[pos:n]), n
}

func get_uint8_pos(data []byte, pos int) (uint8, int) {
    return data[pos], pos + 1
}

func get_uint16_pos(data []byte, pos int) (uint16, int) {
    var ret uint16
    buf := bytes.NewBuffer(data[pos:pos + 2])
    binary.Read(buf, binary.LittleEndian, &ret)
    return ret, pos + 2
}

func get_uint32_pos(data []byte, pos int) (uint32, int) {
    var ret uint32
    buf := bytes.NewBuffer(data[pos:pos + 4])
    binary.Read(buf, binary.LittleEndian, &ret)
    return ret, pos + 4
}
