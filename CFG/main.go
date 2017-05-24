package main

import (
	"fmt"
	"net/http"
	"os/exec"
	"os"
	"github.com/hhh0pE/cfg-to-cnf/CFG"
)

func indexAction(w http.ResponseWriter, r *http.Request) {

	r.ParseForm()

	var cfg_input, string_input string
	var post_reply string
	if len(r.Form) > 0 {
		cfg_input = r.Form["cfg"][0]
        string_input = r.Form["string"][0]
		cfg := CFG.NewGrammarFromString(cfg_input)
		cfg.ToCNF()
        if cfg.TestString(string_input) {
            post_reply = "<p class=\"alert alert-success\">String \""+string_input+"\" is accepted!"
        } else {
            post_reply = "<p class=\"alert alert-danger\">String \""+string_input+"\" is rejected.</p>"
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
        <textarea name="cfg" id="cfg" class="form-control" style="height:200px;">` + cfg_input + `</textarea>
        <br />
        <label for="string">Password to Test</label>
        <input type="text" value="`+string_input+`" class="form-control" id="string" name="string" />
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
