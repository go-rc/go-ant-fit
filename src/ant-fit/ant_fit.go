package ant_fit

import (
    "errors"
    "fmt"
)

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

// message interface

type FitMsg interface {
    Name() string
    Text() string
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
    case 1: return "device"
    case 2: return "settings"
    case 3: return "sport"
    case 4: return "activity"
    case 5: return "workout"
    case 6: return "course"
    case 7: return "schedules"
    case 9: return "weight"
    case 10: return "totals"
    case 11: return "goals"
    case 14: return "blood_pressure"
    case 15: return "monitoring"
    case 20: return "activity_summary"
    case 28: return "monitoring_daily"
    default: return fmt.Sprintf("invalid#%d", msg.msgtype)
    }
}

func (msg *MsgFileId) Name() string {
    return "file_id"
}

func (msg *MsgFileId) Text() string {
    return fmt.Sprintf("file_id #%d msgtype %s mfct %d prod %d ser# %d time %d",
        msg.number, msg.msgtype_name(), msg.manufacturer, msg.product,
        msg.serial_number, msg.time_created)
}

func NewMsgFileId(def *FitDefinition, data []byte) (*MsgFileId, error) {
    const explen int = 15

    if len(data) != explen {
        errfmt := "FileId message should be %d bytes, not %d"
        return nil, errors.New(fmt.Sprintf(errfmt, explen, len(data)))
    }

    msg := new(MsgFileId)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 0: msg.msgtype, pos = get_uint8_pos(data, pos)
        case 1: msg.manufacturer, pos = get_uint16_pos(data, pos)
        case 2: msg.product, pos = get_uint16_pos(data, pos)
        case 3: msg.serial_number, pos = get_uint32_pos(data, pos)
        case 4: msg.time_created, pos = get_uint32_pos(data, pos)
        case 5: msg.number, pos = get_uint16_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad file_id field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// event message

type MsgEvent struct {
    event byte
    event_type byte
    data16 uint16
    data uint32
    event_group uint8
    timestamp uint32
}

func (msg *MsgEvent) event_name() string {
    switch msg.event {
    case 0: return "timer"
    case 3: return "workout"
    case 4: return "workout_step"
    case 5: return "power_down"
    case 6: return "power_up"
    case 7: return "off_course"
    case 8: return "session"
    case 9: return "lap"
    case 10: return "course_point"
    case 11: return "battery"
    case 12: return "virtual_partner_pace"
    case 13: return "hr_high_alert"
    case 14: return "hr_low_alert"
    case 15: return "speed_high_alert"
    case 16: return "speed_low_alert"
    case 17: return "cad_high_alert"
    case 18: return "cad_low_alert"
    case 19: return "power_high_alert"
    case 20: return "power_low_alert"
    case 21: return "recovery_hr"
    case 22: return "battery_low"
    case 23: return "time_duration_alert"
    case 24: return "distance_duration_alert"
    case 25: return "calorie_duration_alert"
    case 26: return "activity"
    case 27: return "fitness_equipment"
    case 28: return "length"
    case 36: return "calibration"
    case 255: return "invalid"
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

func (msg *MsgEvent) Name() string {
    return "event"
}

func (msg *MsgEvent) Text() string {
    return fmt.Sprintf("event tstmp %d evt %s etyp %s", msg.timestamp,
        msg.event_name(), msg.event_type_name())
}

func NewMsgEvent(def *FitDefinition, data []byte) (*MsgEvent, error) {
    const minlen int = 6

    if len(data) < minlen {
        errfmt := "Event message should be at least %d bytes, not %d"
        return nil, errors.New(fmt.Sprintf(errfmt, minlen, len(data)))
    }

    msg := new(MsgEvent)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 0: msg.event, pos = get_uint8_pos(data, pos)
        case 1: msg.event_type, pos = get_uint8_pos(data, pos)
        case 2: msg.data16, pos = get_uint16_pos(data, pos)
        case 3: msg.data, pos = get_uint32_pos(data, pos)
        case 4: msg.event_group, pos = get_uint8_pos(data, pos)
        case 253: msg.timestamp, pos = get_uint32_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad event field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// software message

type MsgSoftware struct {
    message_index uint16
    version uint16
    part_number string
}

func (msg *MsgSoftware) Name() string {
    return "software"
}

func (msg *MsgSoftware) Text() string {
    return fmt.Sprintf("software msgidx %d vers %d part# %d", msg.message_index,
        msg.version, msg.part_number)
}

func NewMsgSoftware(def *FitDefinition, data []byte) (*MsgSoftware, error) {
    const minlen int = 5

    if len(data) < minlen {
        errfmt := "Software message should be at least %d bytes, not %d"
        return nil, errors.New(fmt.Sprintf(errfmt, minlen, len(data)))
    }

    msg := new(MsgSoftware)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 0: msg.message_index, pos = get_uint16_pos(data, pos)
        case 1: msg.version, pos = get_uint16_pos(data, pos)
        case 2: msg.part_number, pos = get_string_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad software field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// file_creator message

type MsgFileCreator struct {
    software_version uint16
    hardware_version byte
}

func (msg *MsgFileCreator) Name() string {
    return "file_creator"
}

func (msg *MsgFileCreator) Text() string {
    return fmt.Sprintf("file_creator soft %d hard %d", msg.software_version,
        msg.hardware_version)
}

func NewMsgFileCreator(def *FitDefinition,
    data []byte) (*MsgFileCreator, error) {
    const explen int = 3

    if len(data) < explen {
        errfmt := "FileId message should be %d bytes, not %d"
        return nil, errors.New(fmt.Sprintf(errfmt, explen, len(data)))
    }

    msg := new(MsgFileCreator)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 0: msg.software_version, pos = get_uint16_pos(data, pos)
        case 1: msg.hardware_version, pos = get_uint8_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad event field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// device info message

type MsgDeviceInfo struct {
    timestamp uint32
    device_index byte
    device_type byte
    manufacturer uint16
    serial_number uint32
    product uint16
    software_version uint16
    hardware_version uint16
    cum_operating_time uint32
    battery_voltage uint32
    battery_status uint32
}

func (msg *MsgDeviceInfo) device_type_name() string {
    switch msg.device_type {
    case 1: return "antfs"
    case 11: return "bike_power"
    case 12: return "environment_sensor_legacy"
    case 15: return "multi_sport_speed_distance"
    case 16: return "control"
    case 17: return "fitness_equipment"
    case 18: return "blood_pressure"
    case 19: return "geocache_node"
    case 20: return "light_electric_vehicle"
    case 25: return "env_sensor"
    case 119: return "weight_scale"
    case 120: return "heart_rate"
    case 121: return "bike_speed_cadence"
    case 122: return "bike_cadence"
    case 123: return "bike_speed"
    case 124: return "stride_speed_distance"
    default: return fmt.Sprintf("unknown#%d", msg.device_type)
    }
}

func (msg *MsgDeviceInfo) Name() string {
    return "device_info"
}

func (msg *MsgDeviceInfo) Text() string {
    return fmt.Sprintf("device_info tstmp %d idx %d dtyp %s mfr %d ser# %d" +
        " prod %d soft %d hard %d optime %d volt %d stat %d", msg.timestamp,
        msg.device_index, msg.device_type_name(), msg.manufacturer,
        msg.serial_number, msg.product, msg.software_version,
        msg.hardware_version, msg.cum_operating_time, msg.battery_voltage,
        msg.battery_status)
}

func NewMsgDeviceInfo(def *FitDefinition, data []byte) (*MsgDeviceInfo, error) {
    const minlen int = 6

    if len(data) < minlen {
        errfmt := "Device info message should be at least %d bytes, not %d"
        return nil, errors.New(fmt.Sprintf(errfmt, minlen, len(data)))
    }

    msg := new(MsgDeviceInfo)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 0: msg.device_index, pos = get_uint8_pos(data, pos)
        case 1: msg.device_type, pos = get_uint8_pos(data, pos)
        case 2: msg.manufacturer, pos = get_uint16_pos(data, pos)
        case 3: msg.serial_number, pos = get_uint32_pos(data, pos)
        case 4: msg.product, pos = get_uint16_pos(data, pos)
        case 5: msg.software_version, pos = get_uint16_pos(data, pos)
        case 6: msg.hardware_version, pos = get_uint16_pos(data, pos)
        case 7: msg.cum_operating_time, pos = get_uint32_pos(data, pos)
        case 8: fmt.Printf("Ignoring device_info field #%d\n",
            def.fields[i].num)
        case 9: fmt.Printf("Ignoring device_info field #%d\n",
            def.fields[i].num)
        case 10: msg.battery_voltage, pos = get_uint32_pos(data, pos)
        case 11: msg.battery_status, pos = get_uint32_pos(data, pos)
        case 15: fmt.Printf("Ignoring device_info field #%d\n",
            def.fields[i].num)
        case 16: fmt.Printf("Ignoring device_info field #%d\n",
            def.fields[i].num)
        case 253: msg.timestamp, pos = get_uint32_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad device_info field #%d",
                def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// unknown message

type MsgUnknown struct {
    global_num uint16
    data []byte
}

func (msg *MsgUnknown) Name() string {
    return fmt.Sprintf("unknown#%d", msg.global_num)
}

func (msg *MsgUnknown) Text() string {
    return fmt.Sprintf("unknown#%d", msg.global_num)
}

func NewMsgUnknown(def *FitDefinition, data []byte,
    global_num uint16) (*MsgUnknown, error) {
    msg := new(MsgUnknown)

    msg.global_num = global_num
    msg.data = make([]byte, len(data))
    copy(msg.data, data)

    return msg, nil
}
