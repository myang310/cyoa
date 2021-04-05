package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

const (
	storyTemplate = "students/myang/templates/main.html"
)


func getFileBytes(fileName string) []byte {
	f, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(f)
	if err != nil {
		panic(err)
	}

	return buf.Bytes()
}

type StoryOption struct {
	Text string `json:"text"`
	Arc string `json:"arc"`
}

type Story struct {
	Title string `json:"title"`
	Story []string `json:"story"`
	Options []StoryOption `json:"options"`
}


func parseJSON(data []byte) (stories map[string]Story, err error) {
	err = json.Unmarshal(data, &stories)

	return stories, err
}

type RouteHandler struct {
	StoryArcs map[string]Story
	Template *template.Template
}

func (handler RouteHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var path = req.URL.Path
	if path == "/" {
		http.Redirect(res, req, "/intro", http.StatusPermanentRedirect)
	}
	path = strings.TrimLeft(path, "/")

	if storyArc, ok := handler.StoryArcs[path]; ok {
		err := handler.Template.Execute(res, storyArc)
		if err != nil {
			panic(err)
		}
	} else {
		log.Printf("No story exists for %s", path)
	}

}

func main() {
	var jsonFile = flag.String("jsonFile", "gopher.json", "The file containing the entire JSON for the story")
	parsedJSON, err := parseJSON(getFileBytes(*jsonFile))
	if err != nil {
		panic(err)
	}

	tmpl, err := template.ParseFiles(storyTemplate)
	if err != nil {
		panic(err)
	}

	http.ListenAndServe(":8080", RouteHandler{StoryArcs: parsedJSON, Template: tmpl})
}
