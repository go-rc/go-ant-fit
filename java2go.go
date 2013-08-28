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

func processArgs() (string, []string) {
    usage := false

    dirp := flag.String("d", "", "ANT+ Fit Java source directory")

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

    return *dirp, files
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

func printMessages(dir string, list []*MesgNum) {
    for _, m := range list {
        msg, err := java2go.NewMessage(dir, m.name)
        if err != nil {
            fmt.Fprintf(os.Stderr, "Cannot read %s: %s\n", m.name, err)
        } else {
            msg.PrintFuncs()
        }
    }

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
    dir, files := processArgs()

    if len(files) == 0 {
        if dir == "" {
            fmt.Fprintln(os.Stderr, "Please specify a directory")
        } else {
            list, err := readMessages(dir)
            if err != nil {
                fmt.Fprintf(os.Stderr, "Cannot read MesgNum: %s\n", err)
            } else {
                printMessages(dir, list)
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
