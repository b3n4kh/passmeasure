package main

import (
    "fmt"
    "io/ioutil"
    "strings"
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

var dictonary []string

func getDictonary(length int, alphabet []string) []string {
  dictonaryRecurs(length, alphabet, "", len(alphabet))
  return dictonary
}

func dictonaryRecurs(length int, alphabet []string, prefix string, position int) string {

  if length == 0 {
    dictonary = append(dictonary, prefix)
    return ""
  }

  for i := 0; i < position; i++ {
    newPrefix := prefix + alphabet[i]
    dictonaryRecurs(length - 1, alphabet, newPrefix, position)
  }
  return prefix
}

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
  dat, err := ioutil.ReadFile("./10-million-combos.txt")
  check(err)
  flag.IntVar(&alphabetsize, "s", 1, "size of the alphabet")
  flag.Parse()
  inputString := string(dat)
  size := strings.Count(inputString, "") - 1
  alphabet := getalphabet()
  dict := getDictonary(alphabetsize, alphabet)
  result := countChars(inputString, dict, size)
  fmt.Println(size)
  toString(result)
}
