// forms.go
package main

import (
	"html/template"
	"net/http"
	"strconv"
)

type ContactDetails struct {
	Message string
	Lines   int
}

type FormData struct {
	Success  bool
	Messages []string
	Dada     []string
}

var Messages []string

func main() {
	tmpl := template.Must(template.ParseFiles("forms.html"))

	http.HandleFunc("/", formHandler(tmpl))

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

		formData := FormData{Success: true,
			Messages: Messages,
			Dada:     dada,
		}

		// tmpl.Execute(w, struct{ Success bool }{true})
		tmpl.Execute(w, formData)
	}
}

func addToMessages(Message string) {
	Messages = append(Messages, Message)
}
