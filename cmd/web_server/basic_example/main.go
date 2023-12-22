package main

// ===========================================================================
//
// source: https://go.dev/doc/articles/wiki/
//
// ===========================================================================

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

var validPath = regexp.MustCompile("^/(edit|save|view)/([a-zA-Z0-9]+)$")

var templates = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html"))

func main() {
	http.HandleFunc("/view/", makeHandler(viewHandler))
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	fileName := fmt.Sprintf("data/%s.txt", p.Title)
	return os.WriteFile(fileName, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	fileName := fmt.Sprintf("data/%s.txt", title)
	body, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, fmt.Sprintf("/edit/%s", title), http.StatusFound)
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	if err := p.save(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/view/%s", title), http.StatusFound)
}

func renderTemplate(w http.ResponseWriter, tmpl string, page *Page) {
	if err := templates.ExecuteTemplate(w, fmt.Sprintf("%s.html", tmpl), page); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}
