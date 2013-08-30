package main

import (
    "flag"
    "fmt"
    "os"
    //"sort"
    "./src/ant-fit"
)

func readFit(filename string, verbose bool) error {
    ffile, err := ant_fit.NewFitFile(filename)
    if err != nil {
        return err
    }

    if verbose {
        fmt.Println(ffile.String())
    }

    for true {
        flag, err := ffile.ReadMessage(verbose)
        if err != nil {
            return err
        } else if !flag {
            break
        }
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

func main() {
    verbose, files := processArgs()

    for _, f := range files {
        err := readFit(f, verbose)
        if err != nil {
            fmt.Printf("!! Cannot read %s: %s\n", f, err)
        }
    }
}
