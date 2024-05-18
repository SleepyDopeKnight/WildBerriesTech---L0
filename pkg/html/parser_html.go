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

	if err = page.Execute(w, data); err != nil {
		log.Println(err)
	}
}
