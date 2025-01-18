package main

import (
	"fmt"
	"html/template"
	"net/http"
	"time"
)

var Templates map[string]*template.Template

type Post struct {
	title   string
	content template.HTML
	date    time.Time
	tags    []string
}

func main() {
	fmt.Println("Running server...")

	err := InitDatabase()
	if err != nil {
		fmt.Println(err)
		return
	}

	Templates = make(map[string]*template.Template)

	Templates["index.html"] = template.Must(template.ParseFiles("templates/index.html", "templates/layout.html"))
	Templates["post.html"] = template.Must(template.ParseFiles("templates/post.html", "templates/layout.html"))
	Templates["404.html"] = template.Must(template.ParseFiles("templates/404.html", "templates/layout.html"))
	Templates["list_posts.html"] = template.Must(template.New("list_posts").Funcs(template.FuncMap{
		"add": add,
		"sub": sub,
	}).ParseFiles("templates/list_posts.html", "templates/layout.html"))
	Templates["tags.html"] = template.Must(template.ParseFiles("templates/tags.html", "templates/layout.html"))
	Templates["posts_by_tag.html"] = template.Must(template.ParseFiles("templates/posts_by_tag.html", "templates/layout.html"))
	Templates["about.html"] = template.Must(template.ParseFiles("templates/about.html", "templates/layout.html"))

	http.Handle("/images/", http.StripPrefix("/images/", http.FileServer(http.Dir("images"))))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/post/", GetPostHandler)
	http.HandleFunc("/posts/", ListPostsHandler)
	http.HandleFunc("/tags/", GetAllTagsHandler)
	http.HandleFunc("/tag/", GetPostsByTagHandler)
	http.HandleFunc("/about/", AboutHandler)
	http.HandleFunc("/", IndexHandler)

	http.ListenAndServe(":3000", nil)

	DB.Close()
}

func add(a, b int) int { return a + b }
func sub(a, b int) int { return a - b }
