package main

import (
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/ffmiyo/glince"
)

func main() {
	linkCount := 4
	term := "panasonic g7 review"

	reviews := glince.Search(term, linkCount)
	writeHtml(term, reviews)
}

func writeHtml(term string, pages []glince.Result) {
	t := template.Must(template.ParseFiles("./layout.gohtml"))
	f, err := os.Create(strings.Join(strings.Fields(term), "-") + ".html")
	if err != nil {
		log.Println("create file: ", err)
		return
	}
	defer f.Close()
	err = t.Execute(f, pages)
	if err != nil {
		log.Print("execute: ", err)
		return
	}
	log.Println("File written")
}
