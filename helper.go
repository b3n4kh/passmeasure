package main

import (
    "fmt"
    "io"
    "bufio"
    "os"
)

func readline() {
    f, err := os.Open("/home/ben/bac/passwords/yahoo-passwords.txt")
    check(err)

    defer f.Close()
    r := bufio.NewReaderSize(f, 4*1024)
    line, isPrefix, err := r.ReadLine()
    for err == nil && !isPrefix {
        s := string(line)
        fmt.Println(s)
        line, isPrefix, err = r.ReadLine()
    }

    if isPrefix {
        fmt.Println("buffer size to small")
        return
    }

    if err != io.EOF {
        fmt.Println(err)
        return
    }
}
