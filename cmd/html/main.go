package main

import (
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/ffmiyo/ttpr"
)

func main() {
	linkCount := 4
	term := "the gereg review"
	reviews := ttpr.Search(term, linkCount)
	writeHtml(term, reviews)
}

func writeHtml(term string, pages []ttpr.Result) {
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
