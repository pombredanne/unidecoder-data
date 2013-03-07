package main

import (
    "bufio"
    "io"
    "fmt"
    "log"
    "os"
    "strings"
)

func main() {
    programFile := "data.go"
    program, err := os.Create(programFile)
    if err != nil {
        log.Fatalln("Cannot open %s", programFile)
    }
    defer program.Close()

    codePoint := make(map[byte][]string)

    program.WriteString("package unidecoder\n\n")

    for i := 0; i < 256; i++ {
        var errorReadFile error

        filename := fmt.Sprintf("x%02x.yml", i)
        var content []string

        file, err := os.Open(filename)
        if err != nil {
            log.Printf("skipping %s", filename)
            errorReadFile = err
            continue
        }
        defer file.Close()

        reader := bufio.NewReader(file)

        for {
            lineBytes, _, err := reader.ReadLine()
            if err == io.EOF {
                errorReadFile = nil
                break
            } else if err != nil {
                log.Printf("while reading %s, %v", filename, err)
                errorReadFile = err
                break
            }

            line := string(lineBytes)

            if line == "---" {
                continue
            }

            if strings.HasPrefix(line, "- ") {
                line = line[2:]
            }

            if line[0] == '\'' || line[0] == '"' {
                line = line[1:len(line) - 1]
            }

            content = append(content, line)
        }

        if errorReadFile == nil {
            codePoint[byte(i)] = content
        }
    }

    program.WriteString("const CodePoint = " + fmt.Sprintf("%#v", codePoint))
}
