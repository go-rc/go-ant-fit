package main

import (
    "errors"
    "flag"
    "fmt"
    "os"
    "path"
    "strings"
    "./src/java2go"
)

func checkFile(dir string, filename string) (string, error) {
    var fullpath string
    if dir == "" {
        if _, err := os.Stat(filename); os.IsNotExist(err) {
            return "", errors.New("Cannot find " + filename)
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
                return "", errors.New("Cannot find " + filename)
            }
        }
    }

    return fullpath, nil
}

func processArgs() []string {
    usage := false

    dirp := flag.String("d", "", "ANT+ Fit Java source directory")

    flag.Parse()

    files := make([]string, 0)
    for _, f := range flag.Args() {
        file, err := checkFile(*dirp, f)
        if err != nil {
            fmt.Printf("File %s does not exist\n", f)
            usage = true
            continue
        }

        files = append(files, file)
    }

    if usage {
        fmt.Print("Usage: java2go.go")
        fmt.Print("[-d srcdir]")
        fmt.Print("file [file ...]")
        fmt.Println()

        os.Exit(1)
    }

    return files
}

func main() {
    files := processArgs()

    for _, f := range files {
        msg, err := java2go.NewMessage(f)
        if err != nil {
            fmt.Printf("!! Cannot read %s: %s\n", f, err)
        }

        msg.PrintFuncs()
    }
}
