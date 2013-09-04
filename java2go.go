package main

import (
    "bufio"
    "errors"
    "flag"
    "fmt"
    "os"
    "path"
    "regexp"
    "strconv"
    "unicode"
    "./src/java2go"
)

func processArgs() (string, []string, bool) {
    usage := false

    dirp := flag.String("d", "", "ANT+ Fit Java source directory")
    readp := flag.Bool("r", false, "If set, print readData() function")

    flag.Parse()

    files := make([]string, 0)
    for _, f := range flag.Args() {
        files = append(files, f)
    }

    if usage {
        fmt.Print("Usage: java2go.go")
        fmt.Print("[-d srcdir]")
        fmt.Print("[file file ...]")
        fmt.Println()

        os.Exit(1)
    }

    return *dirp, files, *readp
}

var msg_pat = regexp.MustCompile(`^\s+public\s+static\s+final\s+int\s+` +
    `(.*)\s+=\s+(\d+);\s*$`)

type MesgNum struct {
    num int
    name string
}

func readMessages(dir string) ([]*MesgNum, error) {
    fullpath := path.Join(dir, "MesgNum.java")
    if _, err := os.Stat(fullpath); os.IsNotExist(err) {
        return nil, errors.New(fmt.Sprintf("Cannot find \"%s\"", fullpath))
    }

    fd, err := os.Open(fullpath)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("Cannot open \"%s\"", fullpath))
    }
    defer fd.Close()

    var list []*MesgNum

    scan := bufio.NewScanner(fd)
    for scan.Scan() {
        line := scan.Text()

        m := msg_pat.FindStringSubmatch(line)
        if m != nil {
            cls := toClassName(m[1])

            val, err := strconv.ParseInt(m[2], 0, 32)
            if err != nil {
                return nil, err
            }

            list = append(list, &MesgNum{int(val), cls})
        }
    }

    return list, nil
}

func printMsgUnknown() {
    fmt.Println()
    fmt.Println("// unknown message")
    fmt.Println()
    fmt.Println("type MsgUnknown struct {")
    fmt.Println("    global_num uint16")
    fmt.Println("    data []byte")
    fmt.Println("}")
    fmt.Println()
    fmt.Println("func (msg *MsgUnknown) Name() string {")
    fmt.Println("    return fmt.Sprintf(\"unknown#%d\", msg.global_num)")
    fmt.Println("}")
    fmt.Println()
    fmt.Println("func (msg *MsgUnknown) Text() string {")
    fmt.Println("    return fmt.Sprintf(\"unknown#%d\", msg.global_num)")
    fmt.Println("}")
    fmt.Println()
    fmt.Println("func NewMsgUnknown(def *FitDefinition, data []byte,")
    fmt.Println("    global_num uint16) (*MsgUnknown, error) {")
    fmt.Println("    msg := new(MsgUnknown)")
    fmt.Println()
    fmt.Println("    msg.global_num = global_num")
    fmt.Println("    msg.data = make([]byte, len(data))")
    fmt.Println("    copy(msg.data, data)")
    fmt.Println()
    fmt.Println("    return msg, nil")
    fmt.Println("}")
}

func printInitial() {
    fmt.Println("package ant_fit")
    fmt.Println()

    fmt.Println("import (")
    fmt.Println("    \"errors\"")
    fmt.Println("    \"fmt\"")
    fmt.Println(")")
    fmt.Println()
    fmt.Println("type FitFieldDefinition struct {")
    fmt.Println("    num byte")
    fmt.Println("    size byte")
    fmt.Println("    is_endian bool")
    fmt.Println("    base_type byte")
    fmt.Println("}")
    fmt.Println()
    fmt.Println("type FitDefinition struct {")
    fmt.Println("    local_type byte")
    fmt.Println("    little_endian bool")
    fmt.Println("    global_num uint16")
    fmt.Println("    fields []*FitFieldDefinition")
    fmt.Println("    total_bytes uint16")
    fmt.Println("}")
    fmt.Println()
    fmt.Println("// message interface")
    fmt.Println()
    fmt.Println("type FitMsg interface {")
    fmt.Println("    Name() string")
    fmt.Println("    Text() string")
    fmt.Println("}")
    fmt.Println()
}

func printMessages(dir string, list []*MesgNum) {
    printInitial()

    for _, m := range list {
        msg, err := java2go.NewMessage(dir, m.name)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Cannot read %s: %s\n", m.name, err)
        } else {
            msg.PrintFuncs()
        }
    }

    printMsgUnknown()
}

func printReadData(list []*MesgNum) {
    fmt.Println()
    fmt.Println("func (ffile *FitFile) readData(def *FitDefinition,")
    fmt.Println("    time_offset uint32, verbose bool) (FitMsg, error) {")
    fmt.Println()
    fmt.Println("    buf := make([]byte, def.total_bytes)")
    fmt.Println()
    fmt.Println("    n, err := ffile.rdr.Read(buf)")
    fmt.Println("    if err != nil {")
    fmt.Println("        return nil, err")
    fmt.Println("    } else if n != len(buf) {")
    fmt.Println("        return nil, errors.New(fmt.Sprintf(\"Read %d" +
        " bytes, not %d\", n,")
    fmt.Println("            len(buf)))")
    fmt.Println("    }")
    fmt.Println()
    fmt.Println("    switch def.global_num {")

    for _, m := range list {
        fmt.Printf("    case %d: return NewMsg%s(def, buf)\n", m.num, m.name)
    }

    fmt.Println("    default: return NewMsgUnknown(def, buf, def.global_num)")
    fmt.Println("    }")
    fmt.Println("}")
}

func toClassName(mesgnum string) string {
    var class []rune

    capitalize := true
    for _, c := range mesgnum {
        if capitalize {
            class = append(class, rune(c))
            capitalize = false
        }  else if (c == '_') {
            capitalize = true
        } else {
            class = append(class, unicode.ToLower(rune(c)))
        }
    }

    return string(class)
}

func main() {
    dir, files, addReadDataFunc := processArgs()

    if len(files) == 0 {
        if dir == "" {
            fmt.Fprintln(os.Stderr, "Please specify a directory")
        } else {
            list, err := readMessages(dir)
            if err != nil {
                fmt.Fprintf(os.Stderr, "Cannot read MesgNum: %s\n", err)
            } else {
                printMessages(dir, list)

                if addReadDataFunc {
                    printReadData(list)
                }
            }
        }
    } else {
        for _, f := range files {
            msg, err := java2go.NewMessage(dir, f)
            if err != nil {
                fmt.Fprintf(os.Stderr, "Cannot read %s: %s\n", f, err)
            }

            msg.PrintFuncs()
        }
    }
}
