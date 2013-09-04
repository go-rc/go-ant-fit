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
    default: return fmt.Sprintf("unknown#%d", msg.msgtype)
    }
}

func (msg *MsgFileId) Name() string {
    return "file_id"
}

func (msg *MsgFileId) Text() string {
    return fmt.Sprintf("file_id msgtyp %s mfct %d prod %d ser# %d timecre %d" +
        " # %d", msg.msgtype_name(), msg.manufacturer, msg.product,
        msg.serial_number, msg.time_created, msg.number)
}

func NewMsgFileId(def *FitDefinition, data []byte) (*MsgFileId, error) {
    msg := new(MsgFileId)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 0: msg.msgtype, pos = get_byte_pos(data, pos)
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

// capabilities message

type MsgCapabilities struct {
    languages uint8
    sports uint8
    workouts_supported uint32
}

func (msg *MsgCapabilities) Name() string {
    return "capabilities"
}

func (msg *MsgCapabilities) Text() string {
    return fmt.Sprintf("capabilities languages %d sports %d" +
        " workoutssupported %d", msg.languages, msg.sports,
        msg.workouts_supported)
}

func NewMsgCapabilities(def *FitDefinition, data []byte) (*MsgCapabilities, error) {
    msg := new(MsgCapabilities)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 0: msg.languages, pos = get_uint8_pos(data, pos)
        case 1: msg.sports, pos = get_uint8_pos(data, pos)
        case 21: msg.workouts_supported, pos = get_uint32_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad capabilities field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// device_settings message

type MsgDeviceSettings struct {
    utc_offset uint32
}

func (msg *MsgDeviceSettings) Name() string {
    return "device_settings"
}

func (msg *MsgDeviceSettings) Text() string {
    return fmt.Sprintf("device_settings utcoffset %d", msg.utc_offset)
}

func NewMsgDeviceSettings(def *FitDefinition, data []byte) (*MsgDeviceSettings, error) {
    msg := new(MsgDeviceSettings)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 1: msg.utc_offset, pos = get_uint32_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad device_settings field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// user_profile message

type MsgUserProfile struct {
    message_index uint16
    friendly_name string
    gender byte
    age uint8
    height uint8
    weight uint16
    language byte
    elev_setting byte
    weight_setting byte
    resting_heart_rate uint8
    default_max_running_heart_rate uint8
    default_max_biking_heart_rate uint8
    default_max_heart_rate uint8
    hr_setting byte
    speed_setting byte
    dist_setting byte
    power_setting byte
    activity_class byte
    position_setting byte
    temperature_setting byte
    local_id uint16
    global_id byte
}

func (msg *MsgUserProfile) Name() string {
    return "user_profile"
}

func (msg *MsgUserProfile) Text() string {
    return fmt.Sprintf("user_profile msgidx %d friendlyname %s gender %d" +
        " age %d height %d weight %d language %d elevsetting %d" +
        " weightsetting %d restingheartrate %d defaultmaxrunningheartrate %d" +
        " defaultmaxbikingheartrate %d defaultmaxheartrate %d hrsetting %d" +
        " speedsetting %d distsetting %d powersetting %d activityclass %d" +
        " possetting %d tempsetting %d localid %d globalid %d",
        msg.message_index, msg.friendly_name, msg.gender, msg.age, msg.height,
        msg.weight, msg.language, msg.elev_setting, msg.weight_setting,
        msg.resting_heart_rate, msg.default_max_running_heart_rate,
        msg.default_max_biking_heart_rate, msg.default_max_heart_rate,
        msg.hr_setting, msg.speed_setting, msg.dist_setting, msg.power_setting,
        msg.activity_class, msg.position_setting, msg.temperature_setting,
        msg.local_id, msg.global_id)
}

func NewMsgUserProfile(def *FitDefinition, data []byte) (*MsgUserProfile, error) {
    msg := new(MsgUserProfile)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 254: msg.message_index, pos = get_uint16_pos(data, pos)
        case 0: msg.friendly_name, pos = get_string_pos(data, pos)
        case 1: msg.gender, pos = get_byte_pos(data, pos)
        case 2: msg.age, pos = get_uint8_pos(data, pos)
        case 3: msg.height, pos = get_uint8_pos(data, pos)
        case 4: msg.weight, pos = get_uint16_pos(data, pos)
        case 5: msg.language, pos = get_byte_pos(data, pos)
        case 6: msg.elev_setting, pos = get_byte_pos(data, pos)
        case 7: msg.weight_setting, pos = get_byte_pos(data, pos)
        case 8: msg.resting_heart_rate, pos = get_uint8_pos(data, pos)
        case 9: msg.default_max_running_heart_rate, pos = get_uint8_pos(data, pos)
        case 10: msg.default_max_biking_heart_rate, pos = get_uint8_pos(data, pos)
        case 11: msg.default_max_heart_rate, pos = get_uint8_pos(data, pos)
        case 12: msg.hr_setting, pos = get_byte_pos(data, pos)
        case 13: msg.speed_setting, pos = get_byte_pos(data, pos)
        case 14: msg.dist_setting, pos = get_byte_pos(data, pos)
        case 16: msg.power_setting, pos = get_byte_pos(data, pos)
        case 17: msg.activity_class, pos = get_byte_pos(data, pos)
        case 18: msg.position_setting, pos = get_byte_pos(data, pos)
        case 21: msg.temperature_setting, pos = get_byte_pos(data, pos)
        case 22: msg.local_id, pos = get_uint16_pos(data, pos)
        case 23: msg.global_id, pos = get_byte_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad user_profile field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// hrm_profile message

type MsgHrmProfile struct {
    message_index uint16
    enabled byte
    hrm_ant_id uint16
    log_hrv byte
    hrm_ant_id_trans_type uint8
}

func (msg *MsgHrmProfile) Name() string {
    return "hrm_profile"
}

func (msg *MsgHrmProfile) Text() string {
    return fmt.Sprintf("hrm_profile msgidx %d enabled %d hrmantid %d" +
        " loghrv %d hrmantidtranstyp %d", msg.message_index, msg.enabled,
        msg.hrm_ant_id, msg.log_hrv, msg.hrm_ant_id_trans_type)
}

func NewMsgHrmProfile(def *FitDefinition, data []byte) (*MsgHrmProfile, error) {
    msg := new(MsgHrmProfile)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 254: msg.message_index, pos = get_uint16_pos(data, pos)
        case 0: msg.enabled, pos = get_byte_pos(data, pos)
        case 1: msg.hrm_ant_id, pos = get_uint16_pos(data, pos)
        case 2: msg.log_hrv, pos = get_byte_pos(data, pos)
        case 3: msg.hrm_ant_id_trans_type, pos = get_uint8_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad hrm_profile field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// sdm_profile message

type MsgSdmProfile struct {
    message_index uint16
    enabled byte
    sdm_ant_id uint16
    sdm_cal_factor uint16
    odometer uint32
    speed_source byte
    sdm_ant_id_trans_type uint8
    odometer_rollover uint8
}

func (msg *MsgSdmProfile) Name() string {
    return "sdm_profile"
}

func (msg *MsgSdmProfile) Text() string {
    return fmt.Sprintf("sdm_profile msgidx %d enabled %d sdmantid %d" +
        " sdmcalfactor %d odometer %d speedsource %d sdmantidtranstyp %d" +
        " odometerrollover %d", msg.message_index, msg.enabled, msg.sdm_ant_id,
        msg.sdm_cal_factor, msg.odometer, msg.speed_source,
        msg.sdm_ant_id_trans_type, msg.odometer_rollover)
}

func NewMsgSdmProfile(def *FitDefinition, data []byte) (*MsgSdmProfile, error) {
    msg := new(MsgSdmProfile)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 254: msg.message_index, pos = get_uint16_pos(data, pos)
        case 0: msg.enabled, pos = get_byte_pos(data, pos)
        case 1: msg.sdm_ant_id, pos = get_uint16_pos(data, pos)
        case 2: msg.sdm_cal_factor, pos = get_uint16_pos(data, pos)
        case 3: msg.odometer, pos = get_uint32_pos(data, pos)
        case 4: msg.speed_source, pos = get_byte_pos(data, pos)
        case 5: msg.sdm_ant_id_trans_type, pos = get_uint8_pos(data, pos)
        case 7: msg.odometer_rollover, pos = get_uint8_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad sdm_profile field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// bike_profile message

type MsgBikeProfile struct {
    message_index uint16
    name string
    sport byte
    sub_sport byte
    odometer uint32
    bike_spd_ant_id uint16
    bike_cad_ant_id uint16
    bike_spdcad_ant_id uint16
    bike_power_ant_id uint16
    custom_wheelsize uint16
    auto_wheelsize uint16
    bike_weight uint16
    power_cal_factor uint16
    auto_wheel_cal byte
    auto_power_zero byte
    id uint8
    spd_enabled byte
    cad_enabled byte
    spdcad_enabled byte
    power_enabled byte
    crank_length uint8
    enabled byte
    bike_spd_ant_id_trans_type uint8
    bike_cad_ant_id_trans_type uint8
    bike_spdcad_ant_id_trans_type uint8
    bike_power_ant_id_trans_type uint8
    odometer_rollover uint8
}

func (msg *MsgBikeProfile) Name() string {
    return "bike_profile"
}

func (msg *MsgBikeProfile) Text() string {
    return fmt.Sprintf("bike_profile msgidx %d name %s sport %d subsport %d" +
        " odometer %d bikespdantid %d bikecadantid %d bikespdcadantid %d" +
        " bikepowerantid %d customwheelsize %d autowheelsize %d bikeweight %d" +
        " powercalfactor %d autowheelcal %d autopowerzero %d id %d" +
        " spdenabled %d cadenabled %d spdcadenabled %d powerenabled %d" +
        " cranklen %d enabled %d bikespdantidtranstyp %d" +
        " bikecadantidtranstyp %d bikespdcadantidtranstyp %d" +
        " bikepowerantidtranstyp %d odometerrollover %d", msg.message_index,
        msg.name, msg.sport, msg.sub_sport, msg.odometer, msg.bike_spd_ant_id,
        msg.bike_cad_ant_id, msg.bike_spdcad_ant_id, msg.bike_power_ant_id,
        msg.custom_wheelsize, msg.auto_wheelsize, msg.bike_weight,
        msg.power_cal_factor, msg.auto_wheel_cal, msg.auto_power_zero, msg.id,
        msg.spd_enabled, msg.cad_enabled, msg.spdcad_enabled, msg.power_enabled,
        msg.crank_length, msg.enabled, msg.bike_spd_ant_id_trans_type,
        msg.bike_cad_ant_id_trans_type, msg.bike_spdcad_ant_id_trans_type,
        msg.bike_power_ant_id_trans_type, msg.odometer_rollover)
}

func NewMsgBikeProfile(def *FitDefinition, data []byte) (*MsgBikeProfile, error) {
    msg := new(MsgBikeProfile)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 254: msg.message_index, pos = get_uint16_pos(data, pos)
        case 0: msg.name, pos = get_string_pos(data, pos)
        case 1: msg.sport, pos = get_byte_pos(data, pos)
        case 2: msg.sub_sport, pos = get_byte_pos(data, pos)
        case 3: msg.odometer, pos = get_uint32_pos(data, pos)
        case 4: msg.bike_spd_ant_id, pos = get_uint16_pos(data, pos)
        case 5: msg.bike_cad_ant_id, pos = get_uint16_pos(data, pos)
        case 6: msg.bike_spdcad_ant_id, pos = get_uint16_pos(data, pos)
        case 7: msg.bike_power_ant_id, pos = get_uint16_pos(data, pos)
        case 8: msg.custom_wheelsize, pos = get_uint16_pos(data, pos)
        case 9: msg.auto_wheelsize, pos = get_uint16_pos(data, pos)
        case 10: msg.bike_weight, pos = get_uint16_pos(data, pos)
        case 11: msg.power_cal_factor, pos = get_uint16_pos(data, pos)
        case 12: msg.auto_wheel_cal, pos = get_byte_pos(data, pos)
        case 13: msg.auto_power_zero, pos = get_byte_pos(data, pos)
        case 14: msg.id, pos = get_uint8_pos(data, pos)
        case 15: msg.spd_enabled, pos = get_byte_pos(data, pos)
        case 16: msg.cad_enabled, pos = get_byte_pos(data, pos)
        case 17: msg.spdcad_enabled, pos = get_byte_pos(data, pos)
        case 18: msg.power_enabled, pos = get_byte_pos(data, pos)
        case 19: msg.crank_length, pos = get_uint8_pos(data, pos)
        case 20: msg.enabled, pos = get_byte_pos(data, pos)
        case 21: msg.bike_spd_ant_id_trans_type, pos = get_uint8_pos(data, pos)
        case 22: msg.bike_cad_ant_id_trans_type, pos = get_uint8_pos(data, pos)
        case 23: msg.bike_spdcad_ant_id_trans_type, pos = get_uint8_pos(data, pos)
        case 24: msg.bike_power_ant_id_trans_type, pos = get_uint8_pos(data, pos)
        case 37: msg.odometer_rollover, pos = get_uint8_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad bike_profile field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// zones_target message

type MsgZonesTarget struct {
    max_heart_rate uint8
    threshold_heart_rate uint8
    functional_threshold_power uint16
    hr_calc_type byte
    pwr_calc_type byte
}

func (msg *MsgZonesTarget) Name() string {
    return "zones_target"
}

func (msg *MsgZonesTarget) Text() string {
    return fmt.Sprintf("zones_target maxheartrate %d thresholdheartrate %d" +
        " functionalthresholdpower %d hrcalctyp %d pwrcalctyp %d",
        msg.max_heart_rate, msg.threshold_heart_rate,
        msg.functional_threshold_power, msg.hr_calc_type, msg.pwr_calc_type)
}

func NewMsgZonesTarget(def *FitDefinition, data []byte) (*MsgZonesTarget, error) {
    msg := new(MsgZonesTarget)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 1: msg.max_heart_rate, pos = get_uint8_pos(data, pos)
        case 2: msg.threshold_heart_rate, pos = get_uint8_pos(data, pos)
        case 3: msg.functional_threshold_power, pos = get_uint16_pos(data, pos)
        case 5: msg.hr_calc_type, pos = get_byte_pos(data, pos)
        case 7: msg.pwr_calc_type, pos = get_byte_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad zones_target field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// hr_zone message

type MsgHrZone struct {
    message_index uint16
    high_bpm uint8
    name string
}

func (msg *MsgHrZone) Name() string {
    return "hr_zone"
}

func (msg *MsgHrZone) Text() string {
    return fmt.Sprintf("hr_zone msgidx %d highbpm %d name %s",
        msg.message_index, msg.high_bpm, msg.name)
}

func NewMsgHrZone(def *FitDefinition, data []byte) (*MsgHrZone, error) {
    msg := new(MsgHrZone)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 254: msg.message_index, pos = get_uint16_pos(data, pos)
        case 1: msg.high_bpm, pos = get_uint8_pos(data, pos)
        case 2: msg.name, pos = get_string_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad hr_zone field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// power_zone message

type MsgPowerZone struct {
    message_index uint16
    high_value uint16
    name string
}

func (msg *MsgPowerZone) Name() string {
    return "power_zone"
}

func (msg *MsgPowerZone) Text() string {
    return fmt.Sprintf("power_zone msgidx %d highvalue %d name %s",
        msg.message_index, msg.high_value, msg.name)
}

func NewMsgPowerZone(def *FitDefinition, data []byte) (*MsgPowerZone, error) {
    msg := new(MsgPowerZone)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 254: msg.message_index, pos = get_uint16_pos(data, pos)
        case 1: msg.high_value, pos = get_uint16_pos(data, pos)
        case 2: msg.name, pos = get_string_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad power_zone field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// met_zone message

type MsgMetZone struct {
    message_index uint16
    high_bpm uint8
    calories uint16
    fat_calories uint8
}

func (msg *MsgMetZone) Name() string {
    return "met_zone"
}

func (msg *MsgMetZone) Text() string {
    return fmt.Sprintf("met_zone msgidx %d highbpm %d cals %d fatcals %d",
        msg.message_index, msg.high_bpm, msg.calories, msg.fat_calories)
}

func NewMsgMetZone(def *FitDefinition, data []byte) (*MsgMetZone, error) {
    msg := new(MsgMetZone)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 254: msg.message_index, pos = get_uint16_pos(data, pos)
        case 1: msg.high_bpm, pos = get_uint8_pos(data, pos)
        case 2: msg.calories, pos = get_uint16_pos(data, pos)
        case 3: msg.fat_calories, pos = get_uint8_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad met_zone field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// sport message

type MsgSport struct {
    sport byte
    sub_sport byte
    name string
}

func (msg *MsgSport) Name() string {
    return "sport"
}

func (msg *MsgSport) Text() string {
    return fmt.Sprintf("sport sport %d subsport %d name %s", msg.sport,
        msg.sub_sport, msg.name)
}

func NewMsgSport(def *FitDefinition, data []byte) (*MsgSport, error) {
    msg := new(MsgSport)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 0: msg.sport, pos = get_byte_pos(data, pos)
        case 1: msg.sub_sport, pos = get_byte_pos(data, pos)
        case 3: msg.name, pos = get_string_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad sport field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// goal message

type MsgGoal struct {
    message_index uint16
    sport byte
    sub_sport byte
    start_date uint32
    end_date uint32
    msgtype byte
    value uint32
    repeat byte
    target_value uint32
    recurrence byte
    recurrence_value uint16
    enabled byte
}

func (msg *MsgGoal) Name() string {
    return "goal"
}

func (msg *MsgGoal) Text() string {
    return fmt.Sprintf("goal msgidx %d sport %d subsport %d startdate %d" +
        " enddate %d msgtyp %d value %d repeat %d targetvalue %d" +
        " recurrence %d recurrencevalue %d enabled %d", msg.message_index,
        msg.sport, msg.sub_sport, msg.start_date, msg.end_date, msg.msgtype,
        msg.value, msg.repeat, msg.target_value, msg.recurrence,
        msg.recurrence_value, msg.enabled)
}

func NewMsgGoal(def *FitDefinition, data []byte) (*MsgGoal, error) {
    msg := new(MsgGoal)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 254: msg.message_index, pos = get_uint16_pos(data, pos)
        case 0: msg.sport, pos = get_byte_pos(data, pos)
        case 1: msg.sub_sport, pos = get_byte_pos(data, pos)
        case 2: msg.start_date, pos = get_uint32_pos(data, pos)
        case 3: msg.end_date, pos = get_uint32_pos(data, pos)
        case 4: msg.msgtype, pos = get_byte_pos(data, pos)
        case 5: msg.value, pos = get_uint32_pos(data, pos)
        case 6: msg.repeat, pos = get_byte_pos(data, pos)
        case 7: msg.target_value, pos = get_uint32_pos(data, pos)
        case 8: msg.recurrence, pos = get_byte_pos(data, pos)
        case 9: msg.recurrence_value, pos = get_uint16_pos(data, pos)
        case 10: msg.enabled, pos = get_byte_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad goal field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// session message

type MsgSession struct {
    message_index uint16
    timestamp uint32
    event byte
    event_type byte
    start_time uint32
    start_position_lat int32
    start_position_long int32
    sport byte
    sub_sport byte
    total_elapsed_time uint32
    total_timer_time uint32
    total_distance uint32
    total_cycles uint32
    total_calories uint16
    total_fat_calories uint16
    avg_speed uint16
    max_speed uint16
    avg_heart_rate uint8
    max_heart_rate uint8
    avg_cadence uint8
    max_cadence uint8
    avg_power uint16
    max_power uint16
    total_ascent uint16
    total_descent uint16
    total_training_effect uint8
    first_lap_index uint16
    num_laps uint16
    event_group uint8
    trigger byte
    nec_lat int32
    nec_long int32
    swc_lat int32
    swc_long int32
    normalized_power uint16
    training_stress_score uint16
    intensity_factor uint16
    left_right_balance uint16
    avg_stroke_count uint32
    avg_stroke_distance uint16
    swim_stroke byte
    pool_length uint16
    pool_length_unit byte
    num_active_lengths uint16
    total_work uint32
    avg_altitude uint16
    max_altitude uint16
    gps_accuracy uint8
    avg_grade int16
    avg_pos_grade int16
    avg_neg_grade int16
    max_pos_grade int16
    max_neg_grade int16
    avg_temperature int8
    max_temperature int8
    total_moving_time uint32
    avg_pos_vertical_speed int16
    avg_neg_vertical_speed int16
    max_pos_vertical_speed int16
    max_neg_vertical_speed int16
    min_heart_rate uint8
    time_in_hr_zone uint32
    time_in_speed_zone uint32
    time_in_cadence_zone uint32
    time_in_power_zone uint32
    avg_lap_time uint32
    best_lap_index uint16
    min_altitude uint16
}

func (msg *MsgSession) Name() string {
    return "session"
}

func (msg *MsgSession) Text() string {
    return fmt.Sprintf("session msgidx %d tstmp %d evt %d evttyp %d" +
        " starttime %d startposlat %d startposlong %d sport %d subsport %d" +
        " totalelapsedtime %d totaltimertime %d totaldist %d totalcycles %d" +
        " totalcals %d totalfatcals %d avgspeed %d maxspeed %d" +
        " avgheartrate %d maxheartrate %d avgcadence %d maxcadence %d" +
        " avgpower %d maxpower %d totalascent %d totaldescent %d" +
        " totaltrainingeffect %d firstlapidx %d numlaps %d evtgrp %d" +
        " trigger %d neclat %d neclong %d swclat %d swclong %d" +
        " normalizedpower %d trainingstressscore %d intensityfactor %d" +
        " leftrightbalance %d avgstrokecount %d avgstrokedist %d" +
        " swimstroke %d poollen %d poollenunit %d numactivelens %d" +
        " totalwork %d avgalt %d maxalt %d gpsaccuracy %d avggrade %d" +
        " avgposgrade %d avgneggrade %d maxposgrade %d maxneggrade %d" +
        " avgtemp %d maxtemp %d totalmovingtime %d avgposvertspeed %d" +
        " avgnegvertspeed %d maxposvertspeed %d maxnegvertspeed %d" +
        " minheartrate %d timeinhrzone %d timeinspeedzone %d" +
        " timeincadencezone %d timeinpowerzone %d avglaptime %d bestlapidx %d" +
        " minalt %d", msg.message_index, msg.timestamp, msg.event,
        msg.event_type, msg.start_time, msg.start_position_lat,
        msg.start_position_long, msg.sport, msg.sub_sport,
        msg.total_elapsed_time, msg.total_timer_time, msg.total_distance,
        msg.total_cycles, msg.total_calories, msg.total_fat_calories,
        msg.avg_speed, msg.max_speed, msg.avg_heart_rate, msg.max_heart_rate,
        msg.avg_cadence, msg.max_cadence, msg.avg_power, msg.max_power,
        msg.total_ascent, msg.total_descent, msg.total_training_effect,
        msg.first_lap_index, msg.num_laps, msg.event_group, msg.trigger,
        msg.nec_lat, msg.nec_long, msg.swc_lat, msg.swc_long,
        msg.normalized_power, msg.training_stress_score, msg.intensity_factor,
        msg.left_right_balance, msg.avg_stroke_count, msg.avg_stroke_distance,
        msg.swim_stroke, msg.pool_length, msg.pool_length_unit,
        msg.num_active_lengths, msg.total_work, msg.avg_altitude,
        msg.max_altitude, msg.gps_accuracy, msg.avg_grade, msg.avg_pos_grade,
        msg.avg_neg_grade, msg.max_pos_grade, msg.max_neg_grade,
        msg.avg_temperature, msg.max_temperature, msg.total_moving_time,
        msg.avg_pos_vertical_speed, msg.avg_neg_vertical_speed,
        msg.max_pos_vertical_speed, msg.max_neg_vertical_speed,
        msg.min_heart_rate, msg.time_in_hr_zone, msg.time_in_speed_zone,
        msg.time_in_cadence_zone, msg.time_in_power_zone, msg.avg_lap_time,
        msg.best_lap_index, msg.min_altitude)
}

func NewMsgSession(def *FitDefinition, data []byte) (*MsgSession, error) {
    msg := new(MsgSession)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 254: msg.message_index, pos = get_uint16_pos(data, pos)
        case 253: msg.timestamp, pos = get_uint32_pos(data, pos)
        case 0: msg.event, pos = get_byte_pos(data, pos)
        case 1: msg.event_type, pos = get_byte_pos(data, pos)
        case 2: msg.start_time, pos = get_uint32_pos(data, pos)
        case 3: msg.start_position_lat, pos = get_int32_pos(data, pos)
        case 4: msg.start_position_long, pos = get_int32_pos(data, pos)
        case 5: msg.sport, pos = get_byte_pos(data, pos)
        case 6: msg.sub_sport, pos = get_byte_pos(data, pos)
        case 7: msg.total_elapsed_time, pos = get_uint32_pos(data, pos)
        case 8: msg.total_timer_time, pos = get_uint32_pos(data, pos)
        case 9: msg.total_distance, pos = get_uint32_pos(data, pos)
        case 10: msg.total_cycles, pos = get_uint32_pos(data, pos)
        case 11: msg.total_calories, pos = get_uint16_pos(data, pos)
        case 13: msg.total_fat_calories, pos = get_uint16_pos(data, pos)
        case 14: msg.avg_speed, pos = get_uint16_pos(data, pos)
        case 15: msg.max_speed, pos = get_uint16_pos(data, pos)
        case 16: msg.avg_heart_rate, pos = get_uint8_pos(data, pos)
        case 17: msg.max_heart_rate, pos = get_uint8_pos(data, pos)
        case 18: msg.avg_cadence, pos = get_uint8_pos(data, pos)
        case 19: msg.max_cadence, pos = get_uint8_pos(data, pos)
        case 20: msg.avg_power, pos = get_uint16_pos(data, pos)
        case 21: msg.max_power, pos = get_uint16_pos(data, pos)
        case 22: msg.total_ascent, pos = get_uint16_pos(data, pos)
        case 23: msg.total_descent, pos = get_uint16_pos(data, pos)
        case 24: msg.total_training_effect, pos = get_uint8_pos(data, pos)
        case 25: msg.first_lap_index, pos = get_uint16_pos(data, pos)
        case 26: msg.num_laps, pos = get_uint16_pos(data, pos)
        case 27: msg.event_group, pos = get_uint8_pos(data, pos)
        case 28: msg.trigger, pos = get_byte_pos(data, pos)
        case 29: msg.nec_lat, pos = get_int32_pos(data, pos)
        case 30: msg.nec_long, pos = get_int32_pos(data, pos)
        case 31: msg.swc_lat, pos = get_int32_pos(data, pos)
        case 32: msg.swc_long, pos = get_int32_pos(data, pos)
        case 34: msg.normalized_power, pos = get_uint16_pos(data, pos)
        case 35: msg.training_stress_score, pos = get_uint16_pos(data, pos)
        case 36: msg.intensity_factor, pos = get_uint16_pos(data, pos)
        case 37: msg.left_right_balance, pos = get_uint16_pos(data, pos)
        case 41: msg.avg_stroke_count, pos = get_uint32_pos(data, pos)
        case 42: msg.avg_stroke_distance, pos = get_uint16_pos(data, pos)
        case 43: msg.swim_stroke, pos = get_byte_pos(data, pos)
        case 44: msg.pool_length, pos = get_uint16_pos(data, pos)
        case 46: msg.pool_length_unit, pos = get_byte_pos(data, pos)
        case 47: msg.num_active_lengths, pos = get_uint16_pos(data, pos)
        case 48: msg.total_work, pos = get_uint32_pos(data, pos)
        case 49: msg.avg_altitude, pos = get_uint16_pos(data, pos)
        case 50: msg.max_altitude, pos = get_uint16_pos(data, pos)
        case 51: msg.gps_accuracy, pos = get_uint8_pos(data, pos)
        case 52: msg.avg_grade, pos = get_int16_pos(data, pos)
        case 53: msg.avg_pos_grade, pos = get_int16_pos(data, pos)
        case 54: msg.avg_neg_grade, pos = get_int16_pos(data, pos)
        case 55: msg.max_pos_grade, pos = get_int16_pos(data, pos)
        case 56: msg.max_neg_grade, pos = get_int16_pos(data, pos)
        case 57: msg.avg_temperature, pos = get_int8_pos(data, pos)
        case 58: msg.max_temperature, pos = get_int8_pos(data, pos)
        case 59: msg.total_moving_time, pos = get_uint32_pos(data, pos)
        case 60: msg.avg_pos_vertical_speed, pos = get_int16_pos(data, pos)
        case 61: msg.avg_neg_vertical_speed, pos = get_int16_pos(data, pos)
        case 62: msg.max_pos_vertical_speed, pos = get_int16_pos(data, pos)
        case 63: msg.max_neg_vertical_speed, pos = get_int16_pos(data, pos)
        case 64: msg.min_heart_rate, pos = get_uint8_pos(data, pos)
        case 65: msg.time_in_hr_zone, pos = get_uint32_pos(data, pos)
        case 66: msg.time_in_speed_zone, pos = get_uint32_pos(data, pos)
        case 67: msg.time_in_cadence_zone, pos = get_uint32_pos(data, pos)
        case 68: msg.time_in_power_zone, pos = get_uint32_pos(data, pos)
        case 69: msg.avg_lap_time, pos = get_uint32_pos(data, pos)
        case 70: msg.best_lap_index, pos = get_uint16_pos(data, pos)
        case 71: msg.min_altitude, pos = get_uint16_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad session field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// lap message

type MsgLap struct {
    message_index uint16
    timestamp uint32
    event byte
    event_type byte
    start_time uint32
    start_position_lat int32
    start_position_long int32
    end_position_lat int32
    end_position_long int32
    total_elapsed_time uint32
    total_timer_time uint32
    total_distance uint32
    total_cycles uint32
    total_calories uint16
    total_fat_calories uint16
    avg_speed uint16
    max_speed uint16
    avg_heart_rate uint8
    max_heart_rate uint8
    avg_cadence uint8
    max_cadence uint8
    avg_power uint16
    max_power uint16
    total_ascent uint16
    total_descent uint16
    intensity byte
    lap_trigger byte
    sport byte
    event_group uint8
    num_lengths uint16
    normalized_power uint16
    left_right_balance uint16
    first_length_index uint16
    avg_stroke_distance uint16
    swim_stroke byte
    sub_sport byte
    num_active_lengths uint16
    total_work uint32
    avg_altitude uint16
    max_altitude uint16
    gps_accuracy uint8
    avg_grade int16
    avg_pos_grade int16
    avg_neg_grade int16
    max_pos_grade int16
    max_neg_grade int16
    avg_temperature int8
    max_temperature int8
    total_moving_time uint32
    avg_pos_vertical_speed int16
    avg_neg_vertical_speed int16
    max_pos_vertical_speed int16
    max_neg_vertical_speed int16
    time_in_hr_zone uint32
    time_in_speed_zone uint32
    time_in_cadence_zone uint32
    time_in_power_zone uint32
    repetition_num uint16
    min_altitude uint16
    min_heart_rate uint8
    wkt_step_index uint16
}

func (msg *MsgLap) Name() string {
    return "lap"
}

func (msg *MsgLap) Text() string {
    return fmt.Sprintf("lap msgidx %d tstmp %d evt %d evttyp %d starttime %d" +
        " startposlat %d startposlong %d endposlat %d endposlong %d" +
        " totalelapsedtime %d totaltimertime %d totaldist %d totalcycles %d" +
        " totalcals %d totalfatcals %d avgspeed %d maxspeed %d" +
        " avgheartrate %d maxheartrate %d avgcadence %d maxcadence %d" +
        " avgpower %d maxpower %d totalascent %d totaldescent %d intensity %d" +
        " laptrigger %d sport %d evtgrp %d numlens %d normalizedpower %d" +
        " leftrightbalance %d firstlenidx %d avgstrokedist %d swimstroke %d" +
        " subsport %d numactivelens %d totalwork %d avgalt %d maxalt %d" +
        " gpsaccuracy %d avggrade %d avgposgrade %d avgneggrade %d" +
        " maxposgrade %d maxneggrade %d avgtemp %d maxtemp %d" +
        " totalmovingtime %d avgposvertspeed %d avgnegvertspeed %d" +
        " maxposvertspeed %d maxnegvertspeed %d timeinhrzone %d" +
        " timeinspeedzone %d timeincadencezone %d timeinpowerzone %d" +
        " repetitionnum %d minalt %d minheartrate %d wktstepidx %d",
        msg.message_index, msg.timestamp, msg.event, msg.event_type,
        msg.start_time, msg.start_position_lat, msg.start_position_long,
        msg.end_position_lat, msg.end_position_long, msg.total_elapsed_time,
        msg.total_timer_time, msg.total_distance, msg.total_cycles,
        msg.total_calories, msg.total_fat_calories, msg.avg_speed,
        msg.max_speed, msg.avg_heart_rate, msg.max_heart_rate, msg.avg_cadence,
        msg.max_cadence, msg.avg_power, msg.max_power, msg.total_ascent,
        msg.total_descent, msg.intensity, msg.lap_trigger, msg.sport,
        msg.event_group, msg.num_lengths, msg.normalized_power,
        msg.left_right_balance, msg.first_length_index, msg.avg_stroke_distance,
        msg.swim_stroke, msg.sub_sport, msg.num_active_lengths, msg.total_work,
        msg.avg_altitude, msg.max_altitude, msg.gps_accuracy, msg.avg_grade,
        msg.avg_pos_grade, msg.avg_neg_grade, msg.max_pos_grade,
        msg.max_neg_grade, msg.avg_temperature, msg.max_temperature,
        msg.total_moving_time, msg.avg_pos_vertical_speed,
        msg.avg_neg_vertical_speed, msg.max_pos_vertical_speed,
        msg.max_neg_vertical_speed, msg.time_in_hr_zone, msg.time_in_speed_zone,
        msg.time_in_cadence_zone, msg.time_in_power_zone, msg.repetition_num,
        msg.min_altitude, msg.min_heart_rate, msg.wkt_step_index)
}

func NewMsgLap(def *FitDefinition, data []byte) (*MsgLap, error) {
    msg := new(MsgLap)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 254: msg.message_index, pos = get_uint16_pos(data, pos)
        case 253: msg.timestamp, pos = get_uint32_pos(data, pos)
        case 0: msg.event, pos = get_byte_pos(data, pos)
        case 1: msg.event_type, pos = get_byte_pos(data, pos)
        case 2: msg.start_time, pos = get_uint32_pos(data, pos)
        case 3: msg.start_position_lat, pos = get_int32_pos(data, pos)
        case 4: msg.start_position_long, pos = get_int32_pos(data, pos)
        case 5: msg.end_position_lat, pos = get_int32_pos(data, pos)
        case 6: msg.end_position_long, pos = get_int32_pos(data, pos)
        case 7: msg.total_elapsed_time, pos = get_uint32_pos(data, pos)
        case 8: msg.total_timer_time, pos = get_uint32_pos(data, pos)
        case 9: msg.total_distance, pos = get_uint32_pos(data, pos)
        case 10: msg.total_cycles, pos = get_uint32_pos(data, pos)
        case 11: msg.total_calories, pos = get_uint16_pos(data, pos)
        case 12: msg.total_fat_calories, pos = get_uint16_pos(data, pos)
        case 13: msg.avg_speed, pos = get_uint16_pos(data, pos)
        case 14: msg.max_speed, pos = get_uint16_pos(data, pos)
        case 15: msg.avg_heart_rate, pos = get_uint8_pos(data, pos)
        case 16: msg.max_heart_rate, pos = get_uint8_pos(data, pos)
        case 17: msg.avg_cadence, pos = get_uint8_pos(data, pos)
        case 18: msg.max_cadence, pos = get_uint8_pos(data, pos)
        case 19: msg.avg_power, pos = get_uint16_pos(data, pos)
        case 20: msg.max_power, pos = get_uint16_pos(data, pos)
        case 21: msg.total_ascent, pos = get_uint16_pos(data, pos)
        case 22: msg.total_descent, pos = get_uint16_pos(data, pos)
        case 23: msg.intensity, pos = get_byte_pos(data, pos)
        case 24: msg.lap_trigger, pos = get_byte_pos(data, pos)
        case 25: msg.sport, pos = get_byte_pos(data, pos)
        case 26: msg.event_group, pos = get_uint8_pos(data, pos)
        case 32: msg.num_lengths, pos = get_uint16_pos(data, pos)
        case 33: msg.normalized_power, pos = get_uint16_pos(data, pos)
        case 34: msg.left_right_balance, pos = get_uint16_pos(data, pos)
        case 35: msg.first_length_index, pos = get_uint16_pos(data, pos)
        case 37: msg.avg_stroke_distance, pos = get_uint16_pos(data, pos)
        case 38: msg.swim_stroke, pos = get_byte_pos(data, pos)
        case 39: msg.sub_sport, pos = get_byte_pos(data, pos)
        case 40: msg.num_active_lengths, pos = get_uint16_pos(data, pos)
        case 41: msg.total_work, pos = get_uint32_pos(data, pos)
        case 42: msg.avg_altitude, pos = get_uint16_pos(data, pos)
        case 43: msg.max_altitude, pos = get_uint16_pos(data, pos)
        case 44: msg.gps_accuracy, pos = get_uint8_pos(data, pos)
        case 45: msg.avg_grade, pos = get_int16_pos(data, pos)
        case 46: msg.avg_pos_grade, pos = get_int16_pos(data, pos)
        case 47: msg.avg_neg_grade, pos = get_int16_pos(data, pos)
        case 48: msg.max_pos_grade, pos = get_int16_pos(data, pos)
        case 49: msg.max_neg_grade, pos = get_int16_pos(data, pos)
        case 50: msg.avg_temperature, pos = get_int8_pos(data, pos)
        case 51: msg.max_temperature, pos = get_int8_pos(data, pos)
        case 52: msg.total_moving_time, pos = get_uint32_pos(data, pos)
        case 53: msg.avg_pos_vertical_speed, pos = get_int16_pos(data, pos)
        case 54: msg.avg_neg_vertical_speed, pos = get_int16_pos(data, pos)
        case 55: msg.max_pos_vertical_speed, pos = get_int16_pos(data, pos)
        case 56: msg.max_neg_vertical_speed, pos = get_int16_pos(data, pos)
        case 57: msg.time_in_hr_zone, pos = get_uint32_pos(data, pos)
        case 58: msg.time_in_speed_zone, pos = get_uint32_pos(data, pos)
        case 59: msg.time_in_cadence_zone, pos = get_uint32_pos(data, pos)
        case 60: msg.time_in_power_zone, pos = get_uint32_pos(data, pos)
        case 61: msg.repetition_num, pos = get_uint16_pos(data, pos)
        case 62: msg.min_altitude, pos = get_uint16_pos(data, pos)
        case 63: msg.min_heart_rate, pos = get_uint8_pos(data, pos)
        case 71: msg.wkt_step_index, pos = get_uint16_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad lap field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// record message

type MsgRecord struct {
    timestamp uint32
    position_lat int32
    position_long int32
    altitude uint16
    heart_rate uint8
    cadence uint8
    distance uint32
    speed uint16
    power uint16
    compressed_speed_distance byte
    grade int16
    resistance uint8
    time_from_course int32
    cycle_length uint8
    temperature int8
    speed_1s uint8
    cycles uint8
    total_cycles uint32
    compressed_accumulated_power uint16
    accumulated_power uint32
    left_right_balance uint8
    gps_accuracy uint8
    vertical_speed int16
    calories uint16
    left_torque_effectiveness uint8
    right_torque_effectiveness uint8
    left_pedal_smoothness uint8
    right_pedal_smoothness uint8
    combined_pedal_smoothness uint8
    cadence256 uint16
}

func (msg *MsgRecord) Name() string {
    return "record"
}

func (msg *MsgRecord) Text() string {
    return fmt.Sprintf("record tstmp %d poslat %d poslong %d alt %d" +
        " heartrate %d cadence %d dist %d speed %d power %d" +
        " compressedspeeddist %d grade %d resistance %d timefromcourse %d" +
        " cyclelen %d temp %d speed1s %d cycles %d totalcycles %d" +
        " compressedaccumpower %d accumpower %d leftrightbalance %d" +
        " gpsaccuracy %d vertspeed %d cals %d lefttorqueeffectiveness %d" +
        " righttorqueeffectiveness %d leftpedalsmooth %d rightpedalsmooth %d" +
        " combinedpedalsmooth %d cadence256 %d", msg.timestamp,
        msg.position_lat, msg.position_long, msg.altitude, msg.heart_rate,
        msg.cadence, msg.distance, msg.speed, msg.power,
        msg.compressed_speed_distance, msg.grade, msg.resistance,
        msg.time_from_course, msg.cycle_length, msg.temperature, msg.speed_1s,
        msg.cycles, msg.total_cycles, msg.compressed_accumulated_power,
        msg.accumulated_power, msg.left_right_balance, msg.gps_accuracy,
        msg.vertical_speed, msg.calories, msg.left_torque_effectiveness,
        msg.right_torque_effectiveness, msg.left_pedal_smoothness,
        msg.right_pedal_smoothness, msg.combined_pedal_smoothness,
        msg.cadence256)
}

func NewMsgRecord(def *FitDefinition, data []byte) (*MsgRecord, error) {
    msg := new(MsgRecord)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 253: msg.timestamp, pos = get_uint32_pos(data, pos)
        case 0: msg.position_lat, pos = get_int32_pos(data, pos)
        case 1: msg.position_long, pos = get_int32_pos(data, pos)
        case 2: msg.altitude, pos = get_uint16_pos(data, pos)
        case 3: msg.heart_rate, pos = get_uint8_pos(data, pos)
        case 4: msg.cadence, pos = get_uint8_pos(data, pos)
        case 5: msg.distance, pos = get_uint32_pos(data, pos)
        case 6: msg.speed, pos = get_uint16_pos(data, pos)
        case 7: msg.power, pos = get_uint16_pos(data, pos)
        case 8: msg.compressed_speed_distance, pos = get_byte_pos(data, pos)
        case 9: msg.grade, pos = get_int16_pos(data, pos)
        case 10: msg.resistance, pos = get_uint8_pos(data, pos)
        case 11: msg.time_from_course, pos = get_int32_pos(data, pos)
        case 12: msg.cycle_length, pos = get_uint8_pos(data, pos)
        case 13: msg.temperature, pos = get_int8_pos(data, pos)
        case 17: msg.speed_1s, pos = get_uint8_pos(data, pos)
        case 18: msg.cycles, pos = get_uint8_pos(data, pos)
        case 19: msg.total_cycles, pos = get_uint32_pos(data, pos)
        case 28: msg.compressed_accumulated_power, pos = get_uint16_pos(data, pos)
        case 29: msg.accumulated_power, pos = get_uint32_pos(data, pos)
        case 30: msg.left_right_balance, pos = get_uint8_pos(data, pos)
        case 31: msg.gps_accuracy, pos = get_uint8_pos(data, pos)
        case 32: msg.vertical_speed, pos = get_int16_pos(data, pos)
        case 33: msg.calories, pos = get_uint16_pos(data, pos)
        case 43: msg.left_torque_effectiveness, pos = get_uint8_pos(data, pos)
        case 44: msg.right_torque_effectiveness, pos = get_uint8_pos(data, pos)
        case 45: msg.left_pedal_smoothness, pos = get_uint8_pos(data, pos)
        case 46: msg.right_pedal_smoothness, pos = get_uint8_pos(data, pos)
        case 47: msg.combined_pedal_smoothness, pos = get_uint8_pos(data, pos)
        case 52: msg.cadence256, pos = get_uint16_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad record field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// event message

type MsgEvent struct {
    timestamp uint32
    event byte
    event_type byte
    data16 uint16
    data uint32
    event_group uint8
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
    default: return fmt.Sprintf("unknown#%d", msg.event)
    }
}

func (msg *MsgEvent) event_type_name() string {
    switch msg.event_type {
    case 0: return "start"
    case 1: return "stop"
    case 2: return "consecutive_depreciated"
    case 3: return "marker"
    case 4: return "stop_all"
    case 5: return "begin_depreciated"
    case 6: return "end_depreciated"
    case 7: return "end_all_depreciated"
    case 8: return "stop_disable"
    case 9: return "stop_disable_all"
    default: return fmt.Sprintf("unknown#%d", msg.event_type)
    }
}

func (msg *MsgEvent) Name() string {
    return "event"
}

func (msg *MsgEvent) Text() string {
    return fmt.Sprintf("event tstmp %d evt %s evttyp %s data16 %d data %d" +
        " evtgrp %d", msg.timestamp, msg.event_name(), msg.event_type_name(),
        msg.data16, msg.data, msg.event_group)
}

func NewMsgEvent(def *FitDefinition, data []byte) (*MsgEvent, error) {
    msg := new(MsgEvent)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 253: msg.timestamp, pos = get_uint32_pos(data, pos)
        case 0: msg.event, pos = get_byte_pos(data, pos)
        case 1: msg.event_type, pos = get_byte_pos(data, pos)
        case 2: msg.data16, pos = get_uint16_pos(data, pos)
        case 3: msg.data, pos = get_uint32_pos(data, pos)
        case 4: msg.event_group, pos = get_uint8_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad event field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// device_info message

type MsgDeviceInfo struct {
    timestamp uint32
    device_index uint8
    device_type uint8
    manufacturer uint16
    serial_number uint32
    product uint16
    software_version uint16
    hardware_version uint8
    cum_operating_time uint32
    battery_voltage uint16
    battery_status uint8
}

func (msg *MsgDeviceInfo) Name() string {
    return "device_info"
}

func (msg *MsgDeviceInfo) Text() string {
    return fmt.Sprintf("device_info tstmp %d devidx %d devtyp %d mfct %d" +
        " ser# %d prod %d soft %d hard %d cumoptime %d battvolt %d" +
        " battstat %d", msg.timestamp, msg.device_index, msg.device_type,
        msg.manufacturer, msg.serial_number, msg.product, msg.software_version,
        msg.hardware_version, msg.cum_operating_time, msg.battery_voltage,
        msg.battery_status)
}

func NewMsgDeviceInfo(def *FitDefinition, data []byte) (*MsgDeviceInfo, error) {
    msg := new(MsgDeviceInfo)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 253: msg.timestamp, pos = get_uint32_pos(data, pos)
        case 0: msg.device_index, pos = get_uint8_pos(data, pos)
        case 1: msg.device_type, pos = get_uint8_pos(data, pos)
        case 2: msg.manufacturer, pos = get_uint16_pos(data, pos)
        case 3: msg.serial_number, pos = get_uint32_pos(data, pos)
        case 4: msg.product, pos = get_uint16_pos(data, pos)
        case 5: msg.software_version, pos = get_uint16_pos(data, pos)
        case 6: msg.hardware_version, pos = get_uint8_pos(data, pos)
        case 7: msg.cum_operating_time, pos = get_uint32_pos(data, pos)
        case 10: msg.battery_voltage, pos = get_uint16_pos(data, pos)
        case 11: msg.battery_status, pos = get_uint8_pos(data, pos)
        default:
            fmt.Printf("Ignoring device_info field #%d\n", def.fields[i].num)
        }
    }

    return msg, nil
}

// workout message

type MsgWorkout struct {
    sport byte
    capabilities uint32
    num_valid_steps uint16
    wkt_name string
}

func (msg *MsgWorkout) Name() string {
    return "workout"
}

func (msg *MsgWorkout) Text() string {
    return fmt.Sprintf("workout sport %d capabilities %d numvalidsteps %d" +
        " wktname %s", msg.sport, msg.capabilities, msg.num_valid_steps,
        msg.wkt_name)
}

func NewMsgWorkout(def *FitDefinition, data []byte) (*MsgWorkout, error) {
    msg := new(MsgWorkout)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 4: msg.sport, pos = get_byte_pos(data, pos)
        case 5: msg.capabilities, pos = get_uint32_pos(data, pos)
        case 6: msg.num_valid_steps, pos = get_uint16_pos(data, pos)
        case 8: msg.wkt_name, pos = get_string_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad workout field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// workout_step message

type MsgWorkoutStep struct {
    message_index uint16
    wkt_step_name string
    duration_type byte
    duration_value uint32
    target_type byte
    target_value uint32
    custom_target_value_low uint32
    custom_target_value_high uint32
    intensity byte
}

func (msg *MsgWorkoutStep) Name() string {
    return "workout_step"
}

func (msg *MsgWorkoutStep) Text() string {
    return fmt.Sprintf("workout_step msgidx %d wktstepname %s durationtyp %d" +
        " durationvalue %d targettyp %d targetvalue %d" +
        " customtargetvaluelow %d customtargetvaluehigh %d intensity %d",
        msg.message_index, msg.wkt_step_name, msg.duration_type,
        msg.duration_value, msg.target_type, msg.target_value,
        msg.custom_target_value_low, msg.custom_target_value_high,
        msg.intensity)
}

func NewMsgWorkoutStep(def *FitDefinition, data []byte) (*MsgWorkoutStep, error) {
    msg := new(MsgWorkoutStep)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 254: msg.message_index, pos = get_uint16_pos(data, pos)
        case 0: msg.wkt_step_name, pos = get_string_pos(data, pos)
        case 1: msg.duration_type, pos = get_byte_pos(data, pos)
        case 2: msg.duration_value, pos = get_uint32_pos(data, pos)
        case 3: msg.target_type, pos = get_byte_pos(data, pos)
        case 4: msg.target_value, pos = get_uint32_pos(data, pos)
        case 5: msg.custom_target_value_low, pos = get_uint32_pos(data, pos)
        case 6: msg.custom_target_value_high, pos = get_uint32_pos(data, pos)
        case 7: msg.intensity, pos = get_byte_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad workout_step field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// schedule message

type MsgSchedule struct {
    manufacturer uint16
    product uint16
    serial_number uint32
    time_created uint32
    completed byte
    msgtype byte
    scheduled_time uint32
}

func (msg *MsgSchedule) Name() string {
    return "schedule"
}

func (msg *MsgSchedule) Text() string {
    return fmt.Sprintf("schedule mfct %d prod %d ser# %d timecre %d" +
        " completed %d msgtyp %d scheduledtime %d", msg.manufacturer,
        msg.product, msg.serial_number, msg.time_created, msg.completed,
        msg.msgtype, msg.scheduled_time)
}

func NewMsgSchedule(def *FitDefinition, data []byte) (*MsgSchedule, error) {
    msg := new(MsgSchedule)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 0: msg.manufacturer, pos = get_uint16_pos(data, pos)
        case 1: msg.product, pos = get_uint16_pos(data, pos)
        case 2: msg.serial_number, pos = get_uint32_pos(data, pos)
        case 3: msg.time_created, pos = get_uint32_pos(data, pos)
        case 4: msg.completed, pos = get_byte_pos(data, pos)
        case 5: msg.msgtype, pos = get_byte_pos(data, pos)
        case 6: msg.scheduled_time, pos = get_uint32_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad schedule field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// weight_scale message

type MsgWeightScale struct {
    timestamp uint32
    weight uint16
    percent_fat uint16
    percent_hydration uint16
    visceral_fat_mass uint16
    bone_mass uint16
    muscle_mass uint16
    basal_met uint16
    physique_rating uint8
    active_met uint16
    metabolic_age uint8
    visceral_fat_rating uint8
    user_profile_index uint16
}

func (msg *MsgWeightScale) Name() string {
    return "weight_scale"
}

func (msg *MsgWeightScale) Text() string {
    return fmt.Sprintf("weight_scale tstmp %d weight %d percentfat %d" +
        " percenthydration %d visceralfatmass %d bonemass %d musclemass %d" +
        " basalmet %d physiquerating %d activemet %d metabolicage %d" +
        " visceralfatrating %d userprofileidx %d", msg.timestamp, msg.weight,
        msg.percent_fat, msg.percent_hydration, msg.visceral_fat_mass,
        msg.bone_mass, msg.muscle_mass, msg.basal_met, msg.physique_rating,
        msg.active_met, msg.metabolic_age, msg.visceral_fat_rating,
        msg.user_profile_index)
}

func NewMsgWeightScale(def *FitDefinition, data []byte) (*MsgWeightScale, error) {
    msg := new(MsgWeightScale)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 253: msg.timestamp, pos = get_uint32_pos(data, pos)
        case 0: msg.weight, pos = get_uint16_pos(data, pos)
        case 1: msg.percent_fat, pos = get_uint16_pos(data, pos)
        case 2: msg.percent_hydration, pos = get_uint16_pos(data, pos)
        case 3: msg.visceral_fat_mass, pos = get_uint16_pos(data, pos)
        case 4: msg.bone_mass, pos = get_uint16_pos(data, pos)
        case 5: msg.muscle_mass, pos = get_uint16_pos(data, pos)
        case 7: msg.basal_met, pos = get_uint16_pos(data, pos)
        case 8: msg.physique_rating, pos = get_uint8_pos(data, pos)
        case 9: msg.active_met, pos = get_uint16_pos(data, pos)
        case 10: msg.metabolic_age, pos = get_uint8_pos(data, pos)
        case 11: msg.visceral_fat_rating, pos = get_uint8_pos(data, pos)
        case 12: msg.user_profile_index, pos = get_uint16_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad weight_scale field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// course message

type MsgCourse struct {
    sport byte
    name string
    capabilities uint32
}

func (msg *MsgCourse) Name() string {
    return "course"
}

func (msg *MsgCourse) Text() string {
    return fmt.Sprintf("course sport %d name %s capabilities %d", msg.sport,
        msg.name, msg.capabilities)
}

func NewMsgCourse(def *FitDefinition, data []byte) (*MsgCourse, error) {
    msg := new(MsgCourse)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 4: msg.sport, pos = get_byte_pos(data, pos)
        case 5: msg.name, pos = get_string_pos(data, pos)
        case 6: msg.capabilities, pos = get_uint32_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad course field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// course_point message

type MsgCoursePoint struct {
    message_index uint16
    timestamp uint32
    position_lat int32
    position_long int32
    distance uint32
    msgtype byte
    name string
}

func (msg *MsgCoursePoint) Name() string {
    return "course_point"
}

func (msg *MsgCoursePoint) Text() string {
    return fmt.Sprintf("course_point msgidx %d tstmp %d poslat %d poslong %d" +
        " dist %d msgtyp %d name %s", msg.message_index, msg.timestamp,
        msg.position_lat, msg.position_long, msg.distance, msg.msgtype,
        msg.name)
}

func NewMsgCoursePoint(def *FitDefinition, data []byte) (*MsgCoursePoint, error) {
    msg := new(MsgCoursePoint)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 254: msg.message_index, pos = get_uint16_pos(data, pos)
        case 1: msg.timestamp, pos = get_uint32_pos(data, pos)
        case 2: msg.position_lat, pos = get_int32_pos(data, pos)
        case 3: msg.position_long, pos = get_int32_pos(data, pos)
        case 4: msg.distance, pos = get_uint32_pos(data, pos)
        case 5: msg.msgtype, pos = get_byte_pos(data, pos)
        case 6: msg.name, pos = get_string_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad course_point field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// totals message

type MsgTotals struct {
    message_index uint16
    timestamp uint32
    timer_time uint32
    distance uint32
    calories uint32
    sport byte
    elapsed_time uint32
    sessions uint16
    active_time uint32
}

func (msg *MsgTotals) Name() string {
    return "totals"
}

func (msg *MsgTotals) Text() string {
    return fmt.Sprintf("totals msgidx %d tstmp %d timertime %d dist %d" +
        " cals %d sport %d elapsedtime %d sessions %d activetime %d",
        msg.message_index, msg.timestamp, msg.timer_time, msg.distance,
        msg.calories, msg.sport, msg.elapsed_time, msg.sessions,
        msg.active_time)
}

func NewMsgTotals(def *FitDefinition, data []byte) (*MsgTotals, error) {
    msg := new(MsgTotals)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 254: msg.message_index, pos = get_uint16_pos(data, pos)
        case 253: msg.timestamp, pos = get_uint32_pos(data, pos)
        case 0: msg.timer_time, pos = get_uint32_pos(data, pos)
        case 1: msg.distance, pos = get_uint32_pos(data, pos)
        case 2: msg.calories, pos = get_uint32_pos(data, pos)
        case 3: msg.sport, pos = get_byte_pos(data, pos)
        case 4: msg.elapsed_time, pos = get_uint32_pos(data, pos)
        case 5: msg.sessions, pos = get_uint16_pos(data, pos)
        case 6: msg.active_time, pos = get_uint32_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad totals field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// activity message

type MsgActivity struct {
    timestamp uint32
    total_timer_time uint32
    num_sessions uint16
    msgtype byte
    event byte
    event_type byte
    local_timestamp uint32
    event_group uint8
}

func (msg *MsgActivity) Name() string {
    return "activity"
}

func (msg *MsgActivity) Text() string {
    return fmt.Sprintf("activity tstmp %d totaltimertime %d numsessions %d" +
        " msgtyp %d evt %d evttyp %d localtstmp %d evtgrp %d", msg.timestamp,
        msg.total_timer_time, msg.num_sessions, msg.msgtype, msg.event,
        msg.event_type, msg.local_timestamp, msg.event_group)
}

func NewMsgActivity(def *FitDefinition, data []byte) (*MsgActivity, error) {
    msg := new(MsgActivity)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 253: msg.timestamp, pos = get_uint32_pos(data, pos)
        case 0: msg.total_timer_time, pos = get_uint32_pos(data, pos)
        case 1: msg.num_sessions, pos = get_uint16_pos(data, pos)
        case 2: msg.msgtype, pos = get_byte_pos(data, pos)
        case 3: msg.event, pos = get_byte_pos(data, pos)
        case 4: msg.event_type, pos = get_byte_pos(data, pos)
        case 5: msg.local_timestamp, pos = get_uint32_pos(data, pos)
        case 6: msg.event_group, pos = get_uint8_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad activity field #%d", def.fields[i].num)
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
    return fmt.Sprintf("software msgidx %d vers %d part# %s", msg.message_index,
        msg.version, msg.part_number)
}

func NewMsgSoftware(def *FitDefinition, data []byte) (*MsgSoftware, error) {
    msg := new(MsgSoftware)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 254: msg.message_index, pos = get_uint16_pos(data, pos)
        case 3: msg.version, pos = get_uint16_pos(data, pos)
        case 5: msg.part_number, pos = get_string_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad software field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// file_capabilities message

type MsgFileCapabilities struct {
    message_index uint16
    msgtype byte
    flags uint8
    directory string
    max_count uint16
    max_size uint32
}

func (msg *MsgFileCapabilities) Name() string {
    return "file_capabilities"
}

func (msg *MsgFileCapabilities) Text() string {
    return fmt.Sprintf("file_capabilities msgidx %d msgtyp %d flags %d" +
        " directory %s maxcount %d maxsize %d", msg.message_index, msg.msgtype,
        msg.flags, msg.directory, msg.max_count, msg.max_size)
}

func NewMsgFileCapabilities(def *FitDefinition, data []byte) (*MsgFileCapabilities, error) {
    msg := new(MsgFileCapabilities)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 254: msg.message_index, pos = get_uint16_pos(data, pos)
        case 0: msg.msgtype, pos = get_byte_pos(data, pos)
        case 1: msg.flags, pos = get_uint8_pos(data, pos)
        case 2: msg.directory, pos = get_string_pos(data, pos)
        case 3: msg.max_count, pos = get_uint16_pos(data, pos)
        case 4: msg.max_size, pos = get_uint32_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad file_capabilities field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// mesg_capabilities message

type MsgMesgCapabilities struct {
    message_index uint16
    file byte
    mesg_num uint16
    count_type byte
    count uint16
}

func (msg *MsgMesgCapabilities) Name() string {
    return "mesg_capabilities"
}

func (msg *MsgMesgCapabilities) Text() string {
    return fmt.Sprintf("mesg_capabilities msgidx %d file %d mesgnum %d" +
        " counttyp %d count %d", msg.message_index, msg.file, msg.mesg_num,
        msg.count_type, msg.count)
}

func NewMsgMesgCapabilities(def *FitDefinition, data []byte) (*MsgMesgCapabilities, error) {
    msg := new(MsgMesgCapabilities)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 254: msg.message_index, pos = get_uint16_pos(data, pos)
        case 0: msg.file, pos = get_byte_pos(data, pos)
        case 1: msg.mesg_num, pos = get_uint16_pos(data, pos)
        case 2: msg.count_type, pos = get_byte_pos(data, pos)
        case 3: msg.count, pos = get_uint16_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad mesg_capabilities field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// field_capabilities message

type MsgFieldCapabilities struct {
    message_index uint16
    file byte
    mesg_num uint16
    field_num uint8
    count uint16
}

func (msg *MsgFieldCapabilities) Name() string {
    return "field_capabilities"
}

func (msg *MsgFieldCapabilities) Text() string {
    return fmt.Sprintf("field_capabilities msgidx %d file %d mesgnum %d" +
        " fieldnum %d count %d", msg.message_index, msg.file, msg.mesg_num,
        msg.field_num, msg.count)
}

func NewMsgFieldCapabilities(def *FitDefinition, data []byte) (*MsgFieldCapabilities, error) {
    msg := new(MsgFieldCapabilities)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 254: msg.message_index, pos = get_uint16_pos(data, pos)
        case 0: msg.file, pos = get_byte_pos(data, pos)
        case 1: msg.mesg_num, pos = get_uint16_pos(data, pos)
        case 2: msg.field_num, pos = get_uint8_pos(data, pos)
        case 3: msg.count, pos = get_uint16_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad field_capabilities field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// file_creator message

type MsgFileCreator struct {
    software_version uint16
    hardware_version uint8
}

func (msg *MsgFileCreator) Name() string {
    return "file_creator"
}

func (msg *MsgFileCreator) Text() string {
    return fmt.Sprintf("file_creator soft %d hard %d", msg.software_version,
        msg.hardware_version)
}

func NewMsgFileCreator(def *FitDefinition, data []byte) (*MsgFileCreator, error) {
    msg := new(MsgFileCreator)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 0: msg.software_version, pos = get_uint16_pos(data, pos)
        case 1: msg.hardware_version, pos = get_uint8_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad file_creator field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// blood_pressure message

type MsgBloodPressure struct {
    timestamp uint32
    systolic_pressure uint16
    diastolic_pressure uint16
    mean_arterial_pressure uint16
    map_3_sample_mean uint16
    map_morning_values uint16
    map_evening_values uint16
    heart_rate uint8
    heart_rate_type byte
    status byte
    user_profile_index uint16
}

func (msg *MsgBloodPressure) Name() string {
    return "blood_pressure"
}

func (msg *MsgBloodPressure) Text() string {
    return fmt.Sprintf("blood_pressure tstmp %d systolicpressure %d" +
        " diastolicpressure %d meanarterialpressure %d map3samplemean %d" +
        " mapmorningvalues %d mapeveningvalues %d heartrate %d" +
        " heartratetyp %d stat %d userprofileidx %d", msg.timestamp,
        msg.systolic_pressure, msg.diastolic_pressure,
        msg.mean_arterial_pressure, msg.map_3_sample_mean,
        msg.map_morning_values, msg.map_evening_values, msg.heart_rate,
        msg.heart_rate_type, msg.status, msg.user_profile_index)
}

func NewMsgBloodPressure(def *FitDefinition, data []byte) (*MsgBloodPressure, error) {
    msg := new(MsgBloodPressure)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 253: msg.timestamp, pos = get_uint32_pos(data, pos)
        case 0: msg.systolic_pressure, pos = get_uint16_pos(data, pos)
        case 1: msg.diastolic_pressure, pos = get_uint16_pos(data, pos)
        case 2: msg.mean_arterial_pressure, pos = get_uint16_pos(data, pos)
        case 3: msg.map_3_sample_mean, pos = get_uint16_pos(data, pos)
        case 4: msg.map_morning_values, pos = get_uint16_pos(data, pos)
        case 5: msg.map_evening_values, pos = get_uint16_pos(data, pos)
        case 6: msg.heart_rate, pos = get_uint8_pos(data, pos)
        case 7: msg.heart_rate_type, pos = get_byte_pos(data, pos)
        case 8: msg.status, pos = get_byte_pos(data, pos)
        case 9: msg.user_profile_index, pos = get_uint16_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad blood_pressure field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// speed_zone message

type MsgSpeedZone struct {
    message_index uint16
    high_value uint16
    name string
}

func (msg *MsgSpeedZone) Name() string {
    return "speed_zone"
}

func (msg *MsgSpeedZone) Text() string {
    return fmt.Sprintf("speed_zone msgidx %d highvalue %d name %s",
        msg.message_index, msg.high_value, msg.name)
}

func NewMsgSpeedZone(def *FitDefinition, data []byte) (*MsgSpeedZone, error) {
    msg := new(MsgSpeedZone)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 254: msg.message_index, pos = get_uint16_pos(data, pos)
        case 0: msg.high_value, pos = get_uint16_pos(data, pos)
        case 1: msg.name, pos = get_string_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad speed_zone field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// monitoring message

type MsgMonitoring struct {
    timestamp uint32
    device_index uint8
    calories uint16
    distance uint32
    cycles uint32
    active_time uint32
    activity_type byte
    activity_subtype byte
    compressed_distance uint16
    compressed_cycles uint16
    compressed_active_time uint16
    local_timestamp uint32
}

func (msg *MsgMonitoring) Name() string {
    return "monitoring"
}

func (msg *MsgMonitoring) Text() string {
    return fmt.Sprintf("monitoring tstmp %d devidx %d cals %d dist %d" +
        " cycles %d activetime %d activitytyp %d activitysubtyp %d" +
        " compresseddist %d compressedcycles %d compressedactivetime %d" +
        " localtstmp %d", msg.timestamp, msg.device_index, msg.calories,
        msg.distance, msg.cycles, msg.active_time, msg.activity_type,
        msg.activity_subtype, msg.compressed_distance, msg.compressed_cycles,
        msg.compressed_active_time, msg.local_timestamp)
}

func NewMsgMonitoring(def *FitDefinition, data []byte) (*MsgMonitoring, error) {
    msg := new(MsgMonitoring)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 253: msg.timestamp, pos = get_uint32_pos(data, pos)
        case 0: msg.device_index, pos = get_uint8_pos(data, pos)
        case 1: msg.calories, pos = get_uint16_pos(data, pos)
        case 2: msg.distance, pos = get_uint32_pos(data, pos)
        case 3: msg.cycles, pos = get_uint32_pos(data, pos)
        case 4: msg.active_time, pos = get_uint32_pos(data, pos)
        case 5: msg.activity_type, pos = get_byte_pos(data, pos)
        case 6: msg.activity_subtype, pos = get_byte_pos(data, pos)
        case 8: msg.compressed_distance, pos = get_uint16_pos(data, pos)
        case 9: msg.compressed_cycles, pos = get_uint16_pos(data, pos)
        case 10: msg.compressed_active_time, pos = get_uint16_pos(data, pos)
        case 11: msg.local_timestamp, pos = get_uint32_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad monitoring field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// hrv message

type MsgHrv struct {
    time uint16
}

func (msg *MsgHrv) Name() string {
    return "hrv"
}

func (msg *MsgHrv) Text() string {
    return fmt.Sprintf("hrv time %d", msg.time)
}

func NewMsgHrv(def *FitDefinition, data []byte) (*MsgHrv, error) {
    msg := new(MsgHrv)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 0: msg.time, pos = get_uint16_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad hrv field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// length message

type MsgLength struct {
    message_index uint16
    timestamp uint32
    event byte
    event_type byte
    start_time uint32
    total_elapsed_time uint32
    total_timer_time uint32
    total_strokes uint16
    avg_speed uint16
    swim_stroke byte
    avg_swimming_cadence uint8
    event_group uint8
    total_calories uint16
    length_type byte
}

func (msg *MsgLength) Name() string {
    return "length"
}

func (msg *MsgLength) Text() string {
    return fmt.Sprintf("length msgidx %d tstmp %d evt %d evttyp %d" +
        " starttime %d totalelapsedtime %d totaltimertime %d totalstrokes %d" +
        " avgspeed %d swimstroke %d avgswimmingcadence %d evtgrp %d" +
        " totalcals %d lentyp %d", msg.message_index, msg.timestamp, msg.event,
        msg.event_type, msg.start_time, msg.total_elapsed_time,
        msg.total_timer_time, msg.total_strokes, msg.avg_speed, msg.swim_stroke,
        msg.avg_swimming_cadence, msg.event_group, msg.total_calories,
        msg.length_type)
}

func NewMsgLength(def *FitDefinition, data []byte) (*MsgLength, error) {
    msg := new(MsgLength)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 254: msg.message_index, pos = get_uint16_pos(data, pos)
        case 253: msg.timestamp, pos = get_uint32_pos(data, pos)
        case 0: msg.event, pos = get_byte_pos(data, pos)
        case 1: msg.event_type, pos = get_byte_pos(data, pos)
        case 2: msg.start_time, pos = get_uint32_pos(data, pos)
        case 3: msg.total_elapsed_time, pos = get_uint32_pos(data, pos)
        case 4: msg.total_timer_time, pos = get_uint32_pos(data, pos)
        case 5: msg.total_strokes, pos = get_uint16_pos(data, pos)
        case 6: msg.avg_speed, pos = get_uint16_pos(data, pos)
        case 7: msg.swim_stroke, pos = get_byte_pos(data, pos)
        case 9: msg.avg_swimming_cadence, pos = get_uint8_pos(data, pos)
        case 10: msg.event_group, pos = get_uint8_pos(data, pos)
        case 11: msg.total_calories, pos = get_uint16_pos(data, pos)
        case 12: msg.length_type, pos = get_byte_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad length field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// monitoring_info message

type MsgMonitoringInfo struct {
    timestamp uint32
    local_timestamp uint32
}

func (msg *MsgMonitoringInfo) Name() string {
    return "monitoring_info"
}

func (msg *MsgMonitoringInfo) Text() string {
    return fmt.Sprintf("monitoring_info tstmp %d localtstmp %d", msg.timestamp,
        msg.local_timestamp)
}

func NewMsgMonitoringInfo(def *FitDefinition, data []byte) (*MsgMonitoringInfo, error) {
    msg := new(MsgMonitoringInfo)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 253: msg.timestamp, pos = get_uint32_pos(data, pos)
        case 0: msg.local_timestamp, pos = get_uint32_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad monitoring_info field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// pad message

type MsgPad struct {
}

func (msg *MsgPad) Name() string {
    return "pad"
}

func (msg *MsgPad) Text() string {
    return fmt.Sprintf("pad")
}

func NewMsgPad(def *FitDefinition, data []byte) (*MsgPad, error) {
    msg := new(MsgPad)

    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        default:
            errmsg := fmt.Sprintf("Bad pad field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// slave_device message

type MsgSlaveDevice struct {
    manufacturer uint16
    product uint16
}

func (msg *MsgSlaveDevice) Name() string {
    return "slave_device"
}

func (msg *MsgSlaveDevice) Text() string {
    return fmt.Sprintf("slave_device mfct %d prod %d", msg.manufacturer,
        msg.product)
}

func NewMsgSlaveDevice(def *FitDefinition, data []byte) (*MsgSlaveDevice, error) {
    msg := new(MsgSlaveDevice)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 0: msg.manufacturer, pos = get_uint16_pos(data, pos)
        case 1: msg.product, pos = get_uint16_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad slave_device field #%d", def.fields[i].num)
            return nil, errors.New(errmsg)
        }
    }

    return msg, nil
}

// cadence_zone message

type MsgCadenceZone struct {
    message_index uint16
    high_value uint8
    name string
}

func (msg *MsgCadenceZone) Name() string {
    return "cadence_zone"
}

func (msg *MsgCadenceZone) Text() string {
    return fmt.Sprintf("cadence_zone msgidx %d highvalue %d name %s",
        msg.message_index, msg.high_value, msg.name)
}

func NewMsgCadenceZone(def *FitDefinition, data []byte) (*MsgCadenceZone, error) {
    msg := new(MsgCadenceZone)

    pos := 0
    for i := 0; i < len(def.fields); i++ {
        switch def.fields[i].num {
        case 254: msg.message_index, pos = get_uint16_pos(data, pos)
        case 0: msg.high_value, pos = get_uint8_pos(data, pos)
        case 1: msg.name, pos = get_string_pos(data, pos)
        default:
            errmsg := fmt.Sprintf("Bad cadence_zone field #%d", def.fields[i].num)
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
