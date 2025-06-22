package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
)

// renders frontend for root page
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "test.page.gohtml")
	})

	fmt.Println("Starting front end service on port 8081")
	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		log.Panic(err)
	}
}

//go:embed templates
var templateFs embed.FS

// when passing in something.page.html, combines with base, footer, header and then returns
func render(w http.ResponseWriter, t string) {

	// partials := []string{
	// 	"./cmd/web/templates/base.layout.gohtml",
	// 	"./cmd/web/templates/header.partial.gohtml",
	// 	"./cmd/web/templates/footer.partial.gohtml",
	// }

	partials := []string{
		"templates/base.layout.gohtml",
		"templates/header.partial.gohtml",
		"templates/footer.partial.gohtml",
	}

	var templateSlice []string
	// templateSlice = append(templateSlice, fmt.Sprintf("./cmd/web/templates/%s", t))
	templateSlice = append(templateSlice, fmt.Sprintf("templates/%s", t))

	for _, x := range partials {
		templateSlice = append(templateSlice, x)
	}

	// tmpl, err := template.ParseFiles(templateSlice...)
	tmpl, err := template.ParseFS(templateFs, templateSlice...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// NEW
	var data struct {
		BrokerURL string
	}

	data.BrokerURL = os.Getenv("BROKER_URL")

	// second parameter stores data that we want available to template
	// allows us to access data in templates
	if err := tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
