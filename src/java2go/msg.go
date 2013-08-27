package java2go

import (
    "bufio"
    "errors"
    "fmt"
    "os"
    "regexp"
    "strconv"
    "strings"
    "unicode"
)

func convertClass(java string) string {
    var newname []rune

    for _, c := range java {
        if unicode.IsUpper(c) {
            if len(newname) > 0 {
                newname = append(newname, rune('_'))
            }
            newname = append(newname, rune(unicode.ToLower(c)))
        } else {
            newname = append(newname, rune(c))
        }
    }

    return string(newname)
}

var base_type_names = [][]string{
    []string{"enum", "byte"},
    []string{"int8", "int8"},
    []string{"uint8", "uint8"},
    []string{"int16", "int16"},
    []string{"uint16", "uint16"},
    []string{"int32", "int32"},
    []string{"uint32", "uint32"},
    []string{"string", "string"},
    []string{"float32", "float32"},
    []string{"float64", "float64"},
    []string{"uint8z", "uint8"},
    []string{"uint16z", "uint16"},
    []string{"uint32z", "uint32"},
    []string{"byte", "byte"},
}

func fitType(num int, base_type int) string {
    low_type := base_type & 0x7f

    if num == 253 && low_type == 6 {
        return "timestamp"
    }

    if num == 254 {
        return "message_index"
    }

    if low_type >= 0 && low_type < len(base_type_names) {
        return base_type_names[low_type][0]
    }

    return fmt.Sprintf("unknown#%d", base_type)
}

func goType(num int, base_type int) string {
    low_type := base_type & 0x7f

    if low_type >= 0 && low_type < len(base_type_names) {
        return base_type_names[low_type][1]
    }

    return fmt.Sprintf("unknown#%d", base_type)
}

type Field struct {
    name string
    num int
    ftype int
    scale float32
    offset float32
    units string
    accumulated bool
}

var short_name_pairs = [][]string{
    []string{"created", "cre"},
    []string{"event", "evt"},
    []string{"group", "grp"},
    []string{"index", "idx"},
    []string{"manufacturer", "mfct"},
    []string{"message", "msg"},
    []string{"number", "#"},
    []string{"product", "prod"},
    []string{"serial", "ser"},
    []string{"timestamp", "tstmp"},
    []string{"type", "typ"},
    []string{"version", "vers"},
    []string{"warevers", ""},
}

func NewField(flds []string) (*Field, error) {
    fld := new(Field)

    fld.name = strings.Trim(flds[0], `'"`)

    for i := 1; i < 3; i++ {
        val, err := strconv.ParseInt(flds[i], 0, 32)
        if err != nil {
            return nil, err
        }

        if i == 1 {
            fld.num = int(val)
        } else {
            fld.ftype = int(val)
        }
    }

    for i := 3; i < 5; i++ {
        val, err := strconv.ParseFloat(flds[i], 32)
        if err != nil {
            return nil, err
        }

        if i == 3 {
            fld.scale = float32(val)
        } else {
            fld.offset = float32(val)
        }
    }

    fld.units = strings.Trim(flds[5], `'"`)
    fld.accumulated = strings.Trim(flds[6], `'"`) == "true"

    return fld, nil
}

func (fld *Field) FormatString() string {
    if fld.ftype == 7 {
        return "%s"
    } else if fld.ftype == 8 || fld.ftype == 9 {
        return "%f"
    }

    return "%d"
}

func (fld *Field) GoType() string {
    return goType(fld.num, fld.ftype)
}

func (fld *Field) Name() string {
    return fld.name
}

func (fld *Field) ShortName() string {
    name := fld.name
    for i := 0; i < len(short_name_pairs); i++ {
        idx := strings.Index(name, short_name_pairs[i][0])
        if idx < 0 {
            continue
        }

        iend := idx + len(short_name_pairs[i][0])
        name = name[0:idx] + short_name_pairs[i][1] + name[iend:]
    }

    return strings.Replace(name, "_", "", -1)
}

func (fld *Field) String() string {
    return fmt.Sprintf("#%d %s %s scal %.1f off %.1f units %s acc %s",
        fld.num, fld.name, fitType(fld.num, fld.ftype), fld.scale, fld.offset,
        fld.units, fld.accumulated)
}

type Message struct {
    cls string
    flds []*Field
}

var msg_class_pat = regexp.MustCompile(`^public\s+class\s+(.*)Mesg\s+` +
    `extends\s+Mesg.*$`)
var msg_field_pat = regexp.MustCompile(`^\s*.*Mesg\.addField\(new\s+` +
    `Field\((.*)\)\);\s*$`)

func NewMessage(path string) (*Message, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("Cannot open \"%s\"\n", path))
    }
    defer file.Close()

    msg := new(Message)

    scan := bufio.NewScanner(file)
    for scan.Scan() {
        line := scan.Text()

        if msg.cls == "" {
            m := msg_class_pat.FindStringSubmatch(line)
            if m != nil {
                msg.cls = m[1]
            }

            continue
        }

        m := msg_field_pat.FindStringSubmatch(line)
        if m == nil {
            continue
        }

        flds := strings.Split(m[1], ", ")
        if len(flds) != 7 {
            fmt.Println("Bad Field line:", line)
            continue
        }

        fld, err := NewField(flds)
        if err != nil {
            fmt.Println("Unusable Field line:", line)
            continue
        }

        msg.flds = append(msg.flds, fld)
    }

    if err := scan.Err(); err != nil {
        return nil, err
    }

    if msg.cls == "" {
        return nil, errors.New("Cannot find class name in " + path)
    }

    return msg, nil
}

func (msg *Message) PrintFuncs() {
    lowcls := convertClass(msg.cls)

    fmt.Printf("// %s message\n", lowcls)
    fmt.Println()
    fmt.Printf("type Msg%s struct {\n", msg.cls)
    for _, f := range msg.flds {
        fmt.Printf("    %s %s\n", f.Name(), f.GoType())
    }
    fmt.Println("}")
    fmt.Println()
    fmt.Printf("func (msg *Msg%s) name() string {\n", msg.cls)
    fmt.Printf("    return \"%s\"\n", lowcls)
    fmt.Println("}")
    fmt.Println()
    fmt.Printf("func (msg *Msg%s) text() string {\n", msg.cls)
    fmtstr := "    return fmt.Sprintf(\"" + lowcls
    for _, f := range msg.flds {
        fldfmt := " " + f.ShortName() + " " + f.FormatString()
        newfmt := fmtstr + fldfmt
        if len(newfmt) <= 77 {
            fmtstr = newfmt
        } else {
            fmt.Println(fmtstr + "\" +")
            fmtstr = "       \"" + fldfmt
        }
    }
    fmtstr += "\""
    for _, f := range msg.flds {
        newfld := ", msg." + f.Name()
        newfmt := fmtstr + newfld
        if len(newfmt) <= 79 {
            fmtstr = newfmt
        } else {
            fmt.Println(fmtstr + ",")
            fmtstr = "        msg." + f.Name()
        }
    }
    fmt.Println(fmtstr + ")")
    fmt.Println("}")
    fmt.Println()
    fmt.Printf("func NewMsg%s(def *FitDefinition, data []byte)" +
        " (*Msg%s, error) {\n", msg.cls, msg.cls)
    fmt.Printf("    msg := new(Msg%s)\n", msg.cls)
    fmt.Println()
    fmt.Println("    pos := 0")
    fmt.Println("    for i := 0; i < len(def.fields); i++ {")
    fmt.Println("        switch def.fields[i].num {")
    n := 0
    for _, f := range msg.flds {
        fmt.Printf("        case %d: msg.%s, pos = get_%s_pos(data, pos)\n",
            n, f.Name(), f.GoType())
        n += 1
    }
    fmt.Println("        default:")
    fmt.Printf("            errmsg := fmt.Sprintf(\"Bad %s field #%%d\"," +
        " def.fields[i].num)\n", lowcls)
    fmt.Println("            return nil, errors.New(errmsg)")
    fmt.Println("        }")
    fmt.Println("    }")
    fmt.Println()
    fmt.Println("    return msg, nil")
    fmt.Println("}")
    fmt.Println()
}
