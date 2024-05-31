// forms.go
package main

import (
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
)

type ContactDetails struct {
	Message string
	Lines   int
}

type FormData struct {
	Success   bool
	Messages  []string
	Dada      []string
	PoemTitle string
}

var Messages []string

func main() {
	tmpl := template.Must(template.ParseFiles("forms.html"))

	http.HandleFunc("/", formHandler(tmpl))

	// Serve static files
	http.ListenAndServe(":8090", nil)
}

func formHandler(tmpl *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			postData := FormData{Success: false, Messages: Messages}

			// tmpl.Execute(w, nil)
			tmpl.Execute(w, postData)
			return
		}

		linesStr := r.FormValue("lines")
		intVal, err := strconv.Atoi(linesStr)

		if err != nil {
			intVal = 5
		}

		details := ContactDetails{
			Message: r.FormValue("message"),
			Lines:   intVal,
		}

		// do something with details
		_ = details

		addToMessages(details.Message)
		dada := createDada(Messages, intVal)
		poemTitle := getPoemTitle(dada)

		formData := FormData{Success: true,
			Messages:  Messages,
			Dada:      dada,
			PoemTitle: poemTitle,
		}

		// tmpl.Execute(w, struct{ Success bool }{true})
		tmpl.Execute(w, formData)
	}
}

func getPoemTitle(dada []string) string {
	rowToPick := rand.Intn(len(dada))

	row := dada[rowToPick]
	wordsInRow := strings.Fields(row)
	wordToPick := rand.Intn(len(wordsInRow))
	title := wordsInRow[wordToPick]
	return title
}

func addToMessages(Message string) {
	Messages = append(Messages, Message)
}
