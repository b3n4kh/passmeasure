package main

import (
    "fmt"
    "io/ioutil"
    "io"
    "bufio"
    "os"
    "strings"
    "bytes"
    "flag"
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

func getDictonary(length int, alphabet []string) []string {
  dict := []string{}


  return dict
}

func dictonaryRecurs(length int, alphabet []string, prefix string, position int) []string {
  if length == 0 {
    fmt.Println(prefix)
  }

  for _,char := range alphabet {
    dict = append(dict, char)
  }
}

//  for _,char := range alphabet {
//      if len(buffer.Bytes()) < length {
//        buffer.WriteString(string(letter))
//      }
//      alphabet = append(alphabet, buffer.String())
//      buffer.Reset()
//    }
//    length--
//  }




func getalphabet() []string {
  alphabet := []string{}

  for letter := 'a'; letter <= 'z'; letter++ {
    alphabet = append(alphabet, string(letter))
  }

  return alphabet
}

func countChars(input string, alphabet []string, size int) []charDistribution {
  var hits int
  var percent float32
  var cD []charDistribution

  for _,element := range alphabet {
    hits = strings.Count(input,string(element))
    percent = float32(hits) / float32(size) * 100

    //fmt.Println(string(letter), hits,"\t", percent, "%")
    cD = append(cD, charDistribution{string(element), hits, percent})
  }

  return cD
}

func toString(output []charDistribution) {
  for _,element := range output {
    fmt.Println(element.match, " ", element.hits, "\t", element.percent, "%")
  }
  fmt.Printf("%s", "\n\n\n")
}


func main() {
  var alphabetsize int
  dat, err := ioutil.ReadFile("/home/ben/bac/passwords/yahoo-passwords.txt")
  check(err)
  flag.IntVar(&alphabetsize, "s", 1, "size of the alphabet")
  flag.Parse()
  inputString := string(dat)
  size := strings.Count(inputString, "") - 1
  alphabet := getalphabet()
  dict := getDictonary(alphabetsize, alphabet)
  //result := countChars(inputString, dict, size)
  fmt.Println(size)
  //toString(result)
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
