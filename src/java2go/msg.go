package java2go

import (
    "bufio"
    "errors"
    "fmt"
    "os"
    "path"
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

var name_entry_pat = regexp.MustCompile(`^\s+(.*)\(\(short\)(\d+)\),\s*$`)

type NameEntry struct {
    name string
    num int
}

type NameFunc struct {
    class string
    name string
    list []NameEntry
}

func NewNameFunc(cls string, dir string, file string,
    name string) (*NameFunc, error) {
    pathstr := path.Join(dir, file + ".java")
    fd, err := os.Open(pathstr)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("Cannot open \"%s\"\n", pathstr))
    }
    defer fd.Close()

    namefunc := new(NameFunc)
    namefunc.class = cls
    namefunc.name = convertClass(name)

    scan := bufio.NewScanner(fd)
    for scan.Scan() {
        line := scan.Text()

        m := name_entry_pat.FindStringSubmatch(line)
        if m != nil {
            val, err := strconv.ParseInt(m[2], 0, 32)
            if err != nil {
                return nil, err
            }

            namefunc.list = append(namefunc.list, NameEntry{m[1], int(val)})
        }
    }

    return namefunc, nil
}

func (nf *NameFunc) IsName(name string) bool {
    return name == nf.name
}

func (nf *NameFunc) PrintFunc() {
    lowcls := nf.name

    fmt.Printf("func (msg *Msg%s) %s_name() string {\n", nf.class, lowcls)
    fmt.Printf("    switch msg.%s {\n", lowcls)
    for _, entry := range nf.list {
        fmt.Printf("    case %d: return \"%s\"\n", entry.num,
            strings.ToLower(entry.name))
    }
    fmt.Printf("    default: return fmt.Sprintf(\"unknown#%%d\", msg.%s)\n",
        lowcls)
    fmt.Println("    }")
    fmt.Println("}")
    fmt.Println()
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
    []string{"ware_vers", ""},
}

func NewField(flds []string) (*Field, error) {
    fld := new(Field)

    fld.name = strings.Trim(flds[0], `'"`)
    if fld.name == "type" {
        fld.name = "msgtype"
    }

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

func (fld *Field) FormatString(has_func bool) string {
    if has_func || fld.ftype == 7 {
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

func (fld *Field) Number() int {
    return fld.num
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
    namefuncs []*NameFunc
}

var msg_class_pat = regexp.MustCompile(`^public\s+class\s+(.*)Mesg\s+` +
    `extends\s+Mesg.*$`)
var msg_field_pat = regexp.MustCompile(`^\s*.*Mesg\.addField\(new\s+` +
    `Field\((.*)\)\);\s*$`)

func NewMessage(dir string, filename string) (*Message, error) {
    var fullpath string
    if dir == "" {
        if _, err := os.Stat(filename); os.IsNotExist(err) {
            errmsg := fmt.Sprintf("Cannot find \"%s\"", filename)
            return nil, errors.New(errmsg)
        }

        fullpath = filename
    } else {
        fullpath = path.Join(dir, filename)
        if _, err := os.Stat(fullpath); os.IsNotExist(err) {
            if !strings.HasSuffix(filename, "Mesg.java") {
                fullpath = fullpath + "Mesg.java"
            } else if !strings.HasSuffix(filename, ".java") {
                fullpath = fullpath + ".java"
            }

            if _, err := os.Stat(fullpath); os.IsNotExist(err) {
                errmsg := fmt.Sprintf("Cannot find \"%s\"", filename)
                return nil, errors.New(errmsg)
            }
        }
    }

    fd, err := os.Open(fullpath)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("Cannot open \"%s\"\n", filename))
    }
    defer fd.Close()

    msg := new(Message)

    scan := bufio.NewScanner(fd)
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
        return nil, errors.New("Cannot find class name in " + filename)
    }

    if msg.cls == "Event" {
        nf, err := NewNameFunc(msg.cls, path.Dir(fullpath), "Event", "event")
        if err != nil {
            return nil, err
        }
        msg.namefuncs = append(msg.namefuncs, nf)

        nf, err = NewNameFunc(msg.cls, path.Dir(fullpath), "EventType",
            "event_type")
        if err != nil {
            return nil, err
        }
        msg.namefuncs = append(msg.namefuncs, nf)
    } else if msg.cls == "FileId" {
        nf, err := NewNameFunc(msg.cls, path.Dir(fullpath), "File", "msgtype")
        if err != nil {
            return nil, err
        }
        msg.namefuncs = append(msg.namefuncs, nf)
    }

    return msg, nil
}

func (msg *Message) hasFunc(fld *Field) bool {
    for _, nf := range msg.namefuncs {
        if nf.IsName(fld.Name()) {
            return true
        }
    }

    return false
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

    for _, nf := range msg.namefuncs {
        nf.PrintFunc()
    }

    fmt.Printf("func (msg *Msg%s) Name() string {\n", msg.cls)
    fmt.Printf("    return \"%s\"\n", lowcls)
    fmt.Println("}")
    fmt.Println()

    fmt.Printf("func (msg *Msg%s) Text() string {\n", msg.cls)
    fmtstr := "    return fmt.Sprintf(\"" + lowcls
    for _, f := range msg.flds {
        fldfmt := " " + f.ShortName() + " " + f.FormatString(msg.hasFunc(f))
        newfmt := fmtstr + fldfmt
        if len(newfmt) <= 77 {
            fmtstr = newfmt
        } else {
            fmt.Println(fmtstr + "\" +")
            fmtstr = "        \"" + fldfmt
        }
    }
    fmtstr += "\""
    for _, f := range msg.flds {
        attr := "msg." + f.Name()
        if msg.hasFunc(f) {
            attr += "_name()"
        }

        newfld := ", " + attr
        newfmt := fmtstr + newfld
        if len(newfmt) <= 79 {
            fmtstr = newfmt
        } else {
            fmt.Println(fmtstr + ",")
            fmtstr = "        " + attr
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
    for _, f := range msg.flds {
        fmt.Printf("        case %d: msg.%s, pos = get_%s_pos(data, pos)\n",
            f.Number(), f.Name(), f.GoType())
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
