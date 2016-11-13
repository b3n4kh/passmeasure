package main

import (
    "fmt"
    "io/ioutil"
    "io"
    "bufio"
    "os"
    "strings"
    "bytes"
)

type charDistribution struct {
	match string
	hits  int
  percent float32
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func countTuple(s string, size int) {
  var hits int
  var percent float32
  var buffer bytes.Buffer

  for letter := 'a'; letter <= 'z'; letter++ {
    for tuple := 'a'; tuple <= 'z'; tuple++ {
      buffer.WriteString(string(letter))
      buffer.WriteString(string(tuple))
      hits = strings.Count(s,buffer.String())
      percent = float32(hits) / float32(size) * 100
      fmt.Println(buffer.String(), hits,"\t", percent, "%")
      //cD = append(cD, charDistribution{buffer.String(), hits, percent})
      buffer.Reset()
    }
  }
}

func getalphabet(length int) []string {
  alphabet := []string{}

  for letter := 'a'; letter <= 'z'; letter++ {
    alphabet = append(alphabet, string(letter))
  }

  return alphabet
}

func countChars(s string, size int) {
  var hits int
  var percent float32
  var cD []charDistribution

  for letter := 'a'; letter <= 'z'; letter++ {
    hits = strings.Count(s,string(letter))
    percent = float32(hits) / float32(size) * 100

    //fmt.Println(string(letter), hits,"\t", percent, "%")
    cD = append(cD, charDistribution{string(letter), hits, percent})
  }

  fmt.Println(cD[20])

  fmt.Printf("%s", "\n\n\n")
}



func main() {
  dat, err := ioutil.ReadFile("/home/ben/bac/passwords/yahoo-passwords.txt")
  check(err)
  inputString := string(dat)
  size := strings.Count(inputString, "") - 1
  fmt.Println(size)
  countChars(inputString, size)
  fmt.Print(getalphabet(1))
}

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
