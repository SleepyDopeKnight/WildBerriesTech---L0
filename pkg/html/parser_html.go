package html

import (
	"html/template"
	"log"
	"net/http"
)

func ParseTemplate(w http.ResponseWriter, templateFile string, data any) {
	page, err := template.ParseFiles(templateFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	err = page.Execute(w, data)
	if err != nil {
		log.Println(err)
	}
}
