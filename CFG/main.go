package main

import (
	"fmt"
	"strings"
	"sort"
	"net/http"
	"os/exec"
	"os"
	"time"
)

type Grammar struct {
	first_rule string
	rules      map[string][]string
}
func (g *Grammar) AddRule(symb, rule string) {
	//fmt.Printf("Adding rule %s to grammar for symbol %s..\n", rule, symb)

	if len(g.rules) == 0 {
		g.first_rule = symb
	}

	g.rules[symb] = append(g.rules[symb], rule)
}

func (g Grammar) ToString() string {
	var output string

    // sort by rule symbols count
    var sizes []int
    symb_sizes := make(map[int][]string)
    for symb, _ := range g.rules {
        size := len(strings.Join(g.rules[symb], ""))
        sizes = append(sizes, size)
        symb_sizes[size] = append(symb_sizes[size], symb)
    }

    sort.Ints(sizes)

    var i int
    for a:=len(sizes)-1; a>=0; a-- {
        if a>=1 && sizes[a] == sizes[a-1] {
            i++
        } else {
            i = 0
        }

        curr_symb := symb_sizes[sizes[a]][i]
        output +=curr_symb+"->"+strings.Join(g.rules[curr_symb], "|")+"\n"
    }

	return output
}

func (g Grammar) ToCNF() string {
    var log string
    log += "Eliminate Unit rules (S->A)\n"

    for symb, _ := range g.rules {
        for i:=0; i < len(g.rules[symb]); i++ {
            s := &(g.rules[symb][i])
            if len(*s)==1 && "A" <= *s && *s <= "Z"  {
                g.rules[symb] = append(g.rules[symb], g.rules[*s]...)
                g.rules[symb] = append(g.rules[symb][:i], g.rules[symb][i+1:]...)
            }
        }
    }

    log += g.ToString() + "\n"

    alphabet := []string{}
		big_alphabet := false
		ini_big_alph := true
		var key string

    for i:='A'; i<'A'+26; i++ {
			if big_alphabet {
				if ini_big_alph {
					if _, exist := g.rules[string(i)]; !exist {
						alphabet = append(alphabet, string(i))
						if i == 'Z' {
							ini_big_alph = false
						}
						continue
					}
				}
				for y:='0'; y<'9'; y++ {
					key = string(i) + string(y)
					alphabet = append(alphabet, key)
				}
			} else {
				key = string(i)
	      if _, exist := g.rules[key]; !exist {
					alphabet = append(alphabet, key)
	      }
			}
    }


    log += "Eliminate the start symbol from right-hand sides (Adding S0 rule)\n"
    g.rules[g.first_rule+"0"] = g.rules[g.first_rule]

    log += g.ToString() + "\n"

    log += "Eliminate right-hand sides with more than 2 nonterminals (SFG->SX, X->FG) \n"

    is_done := false
    new_symbols := make(map[string]string)
    for !is_done {
        is_done = true

        for rule_symb, rules := range g.rules {
            for i, rule := range rules {
                if len(rule)==1 || len(rule)==2 {
                    continue
                }

                is_done = false

                var newS string
                replacing_str := rule[1:len(rule)]
                if _, exist := new_symbols[replacing_str]; exist {
                    newS = new_symbols[replacing_str]
                } else {
                    newS = alphabet[0]
                    new_symbols[replacing_str] = newS
                    alphabet = alphabet[1:len(alphabet)]
                    g.rules[newS] = []string{replacing_str}
                }

                g.rules[rule_symb][i] = rule[:1]+newS
            }
        }
    }


    log += g.ToString() + "\n"

    log += "Eliminate rules with nonsolitary terminals (Sa -> SX, X->a)\n"

    is_done = false
    for !is_done {
        is_done = true
        for rule_symb, rules := range g.rules {
            for i, rule := range rules {
                if len(rule)==1 {
                    continue
                }

                var terminal string
                if strings.ToLower(string(rule[0]))==string(rule[0]) { // find terminals, "a" to lower = "a"
                    terminal = string(rule[0])
                } else if strings.ToLower(string(rule[1]))==string(rule[1]) {
                    terminal = string(rule[1])
                } else {
                    continue
                }

                is_done = false

                if len(alphabet)==1 {
                    log += "the alphabet is not enought :(\n"
                    return log
                }

                var newS string
                if _, exist := new_symbols[terminal]; exist {
                    newS = new_symbols[terminal]
                } else {
                    newS = alphabet[0]
                    new_symbols[terminal] = newS
                    alphabet = alphabet[1:len(alphabet)-1]
                    g.rules[newS] = []string{terminal}
                }

                g.rules[rule_symb][i] = strings.Replace(g.rules[rule_symb][i], terminal, newS, 1)

            }
        }
    }


    log += g.ToString()

    return log
}

func (g Grammar) findRulesByRightSide(right_side string) []string {
    var result []string

    for symb, rules := range g.rules {
        for _, rule := range rules {
            if rule==right_side {
                result = append(result, symb)
            }
        }
    }

    return result
}

func (g Grammar) TestString(input string) bool {

    matrix := make([][]string, len(input))
    for i:=0; i<len(input); i++ {
        matrix[i] = make([]string, len(input)-i)
    }


    // first line
    for i:=0; i<len(input); i++ {
        if symbs := g.findRulesByRightSide(string(input[i])); len(symbs)>0 {
            matrix[0][i] = strings.Join(symbs, ",")
        } else {
            return false
        }
    }

    // other line
    for i:=1; i<len(input); i++ {
        for a:=0; a<len(input)-i; a++ {
            x1 := i-1
            y1 := a
            x2 := 0
            y2 := len(input)-((len(input)-i)-a)

            //fmt.Printf("%d,%d:\n", i, a)

            for loop_size:=0; loop_size<i; loop_size++ {
                //fmt.Printf("%d,%d + %d,%d\n", x1, y1, x2, y2)
                var result []string
                if len(matrix[x1][y1])>0 && len(matrix[x2][y2])>0 {
                    for _, symb1 := range strings.Split(matrix[x1][y1], ",") {
                        for _, symb2 := range strings.Split(matrix[x2][y2], ",") {
                            result = append(result, g.findRulesByRightSide(string(symb1)+string(symb2))...)
                        }
                    }
                }
                if len(result)>0 {
                    matrix[i][a] = strings.Join(result, ",")
                }
                x1--
                x2++
                y2--
            }
        }
    }

    //fmt.Println(matrix[len(input)-1])
    if len(matrix[len(input)-1][0]) > 0 {
        return true
    } else {
        return false
    }
}

func NewGrammarFromString(input string) Grammar {

	grammar := Grammar{}
	grammar.rules = make(map[string][]string)

	input = strings.Trim(input, "\n")
	lines := strings.Split(input, "\n")

	for _, line := range lines {
        line = strings.Replace(line, "->", "→", 1)
		input_arr := strings.FieldsFunc(line, func(c rune) bool {
			if c == '→' || c == '|' || c == ' ' || c == '\r' {
				return true
			}

			return false
		})

		symb := input_arr[0]
		rules := input_arr[1:len(input_arr)]

		for _, rule := range rules {
			grammar.AddRule(symb, rule)
		}
	}

	return grammar
}


func indexAction(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

   cfg_chan := make(chan string)
	var cfg_input, string_input string
	cfg_input = `
S -> ZZZZZZZZDD
Z -> a | b | c | d | e | f | g | h | i | j | k | l | m | n | o | p | q | r | s | t | u | v | w | x | y | z
D -> 1|0|2|3|4|5|6|7|8|9
`
	string_input = "password00"
	var post_reply string

	if len(r.Form) > 0 {
		cfg_input = r.Form["cfg"][0]
		string_input = r.Form["string"][0]
		go func() {
			cfg := NewGrammarFromString(cfg_input)
			cfg.ToCNF()
			reply := "<p class=\"alert alert-danger\">String \""+string_input+"\" is rejected.</p>"
			if cfg.TestString(string_input) {
			  reply = "<p class=\"alert alert-success\">String \""+string_input+"\" is accepted!</p>"
			}
			cnf_grammar := cfg.ToString()
			fmt.Printf(cnf_grammar)
			reply = reply + "\n<div class=\"col-xs-8\">" + strings.Replace(cnf_grammar, "\n", "<br />", -1) + "</div>"
			cfg_chan <- reply
		}()
		select {
		case reply := <-cfg_chan:
			post_reply = reply
		case <-time.After(time.Second * 4):
			post_reply = "<p class=\"alert alert-danger\">TIMEOUT</p>"
		}
	}

	template := `
        <html>
        <head>
				<title>CYK Checker</title>
				<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/3.3.7/css/bootstrap.min.css"
				integrity="sha384-BVYiiSIFeK1dGmJRAkycuHAHRg32OmUcww7on3RYdg4Va+PmSTsz/K68vbdEjh4u" crossorigin="anonymous">
        </head>
        <body>
        <div class="container-fluid">
        <h1>Enter PCFG:</h1>
        <form method="POST">
        <div class="form-group col-sm-10">
        <label for="cfg">Context-Free Grammar</label>
        <textarea name="cfg" id="cfg" class="form-control" style="height:200px;" required>` + cfg_input + `</textarea>
        <br />
        <label for="string">Password to Test</label>
        <input type="text" value="`+string_input+`" class="form-control" id="string" name="string" required/>
        </div>
        <div class="form-group col-sm-10">

        <input type="submit" class="btn btn-primary" value="Test" />
        </div>

        <div class="col-sm-10">
        ` + post_reply + `
        </div>
        </form>
        <div class="col-xs-8">
        </div>
        </div>
        </body>
        </html>`
	w.Write([]byte(template))
}

func main() {
	server_url := "localhost:9023"
	http.HandleFunc("/", indexAction)
	fmt.Println("Starting server at " + server_url)
	if len(os.Args) > 1 {
		exec.Command("google-chrome", "http://"+server_url).Run()
	}
	err := http.ListenAndServe(server_url, nil)
	if err != nil {
		panic("Error when starting server at " + server_url)
	}
}
