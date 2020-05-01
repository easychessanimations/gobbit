package uci

import (
    "bufio"    
    //"log"
    "os"
    "strings"
)

func IterateTextFile(path string, processLineFunc func(string)) {
    file, err := os.Open(path)
    if err != nil {
        //log.Fatal(err)
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        rawLine := scanner.Text()

        line := strings.TrimSpace(rawLine)

        processLineFunc(line)
    }

    if err := scanner.Err(); err != nil {
        //log.Fatal(err)
    }
}
