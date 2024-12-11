package handlers

import (
	"embed"
	"log"
	"text/template"
)

//go:embed templates/*
var templateFS embed.FS

var Template *template.Template

func init() {
	var err error
	Template, err = template.ParseFS(templateFS, "templates/template.html")
	if err != nil {
		log.Fatal(err)
	}
}
