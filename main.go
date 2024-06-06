// forms.go
package main

import (
	"embed"
	"html/template"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"unicode"
)

type FormData struct {
	Success   bool
	Messages  []string
	Dada      []string
	PoemTitle string
}

type ResetData struct {
	Messages []string
}

type PoemData struct {
	Dada      []string
	PoemTitle string
}

var Messages []string

//go:embed templates
var templatesFS embed.FS

func main() {
	formsTmpl := template.Must(template.ParseFS(templatesFS, "templates/forms.html"))
	resetTmpl := template.Must(template.ParseFS(templatesFS, "templates/reset.html"))
	poemTmpl := template.Must(template.ParseFS(templatesFS, "templates/poem.html"))

	http.HandleFunc("/", formHandler(formsTmpl))
	http.HandleFunc("/reset", resetHandler(resetTmpl))
	http.HandleFunc("/poem", poetHandler(poemTmpl))

	// Serve static files
	http.ListenAndServe(":8090", nil)
}

func poetHandler(tmpl *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		intVal := rand.Intn(10) + 5
		dada := createDada(Messages, intVal)
		poemTitle := getPoemTitle(dada)

		poemData := PoemData{
			Dada:      dada,
			PoemTitle: poemTitle,
		}

		tmpl.Execute(w, poemData)
	}

}

func resetHandler(tmpl *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		resetData := ResetData{
			Messages: Messages,
		}

		// tmpl.Execute(w, struct{ Success bool }{true})
		tmpl.Execute(w, resetData)
	}
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

func sanatizeString(input string) string {
	prev := string(input)
	prefiltered := strings.Map(func(r rune) rune {
		if unicode.IsNumber(r) || unicode.IsLetter(r) || unicode.IsSpace(r) {
			return unicode.ToLower(r)
		}
		return -1
	}, prev)

	return prefiltered
}

func shuffleString(s string) string {
	runes := []rune(s)
	rand.Shuffle(len(runes), func(i, j int) {
		runes[i], runes[j] = runes[j], runes[i]
	})
	return string(runes)
}

func createDada(input []string, poetLines int) []string {
	var asOneString string
	for _, a := range input {
		asOneString += " " + a
	}

	words := shuffleWords(getWords(sanatizeString(asOneString)))
	var lines = make([]string, poetLines)

	globalWordCount := 0
	wordCountForLines := getRandomWordCountForLines(poetLines)
	for currentLine, numWordsForThisLine := range wordCountForLines {
		for lineWords := 1; lineWords <= numWordsForThisLine; lineWords++ {
			lines[currentLine] = lines[currentLine] + words[globalWordCount%len(words)] + " "
			globalWordCount += 1
		}
	}
	return lines
}

func getRandomWordCountForLines(poetLines int) []int {
	var lineWordCounts = make([]int, poetLines)
	for i := 0; i < poetLines; i++ {
		lineWordCounts[i] = getRandomWordCountForLine()
	}

	return lineWordCounts
}

func getRandomWordCountForLine() int {
	minWordsPerLine := 1
	maxWordsPerLine := 10
	return rand.Intn(maxWordsPerLine-minWordsPerLine) + minWordsPerLine
}

func getWords(input string) []string {
	return strings.Fields(input)
}

func shuffleWords(input []string) []string {
	retVal := input

	rand.Shuffle(len(retVal), func(i, j int) {
		retVal[i], retVal[j] = retVal[j], retVal[i]
	})

	return retVal
}
