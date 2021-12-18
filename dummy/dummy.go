package main

import (
	"html/template"
	"log"
	"os"
)

type Todo struct {
	Title string
	Done  bool
}

type PageData struct {
	PageTitle string
	Todos     []Todo
}

func main() {

	const tmpl = `<h1>{{.PageTitle}}</h1>
    <ul>
        {{range .Todos}}
            {{if .Done}}
                <li>{{.Title}}</li>
            {{else}}
                <li>{{.Title}}</li>
            {{end}}
        {{end}}
    </ul>`

	// Make and parse the HTML template
	t, err := template.New("webpage").Parse(tmpl)
	if err != nil {
		log.Fatal(err)
	}

	// Initialze a struct storing page data and todo data
	data := PageData{
		PageTitle: "My TODO list",
		Todos: []Todo{
			{Title: "Task 1", Done: false},
			{Title: "Task 2", Done: true},
			{Title: "Task 3", Done: true},
		},
	}

	// Render the data and output using standard output
	t.Execute(os.Stdout, data)
}
