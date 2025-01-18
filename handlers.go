package main

import (
	"net/http"
	"strconv"
	"strings"
)

func GetPostHandler(w http.ResponseWriter, r *http.Request) {
	postName := r.URL.Path[len("/post/"):]
	if len(postName) < 1 {
		fallbackHandler(w, r, "post name missing")
		return
	}

	post, err := GetPost("posts/" + postName + ".org")
	if err != nil {
		fallbackHandler(w, r, "post not found")
		return
	}

	Templates["post.html"].ExecuteTemplate(w, "layout.html", map[string]any{
		"Title":   post.title,
		"Content": post.content,
		"Tags":    post.tags,
		"Date":    post.date.Format("2006-01-02"),
	})
}

func ListPostsHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page := 1
	if p := query.Get("page"); p != "" {
		page, _ = strconv.Atoi(p)
	}
	limit := 10
	offset := (page - 1) * limit

	posts, err := GetPaginatedPosts(limit, offset)
	if err != nil {
		fallbackHandler(w, r, "failed to fetch posts")
		return
	}

	totalCount, err := GetTotalPostCount()
	if err != nil {
		fallbackHandler(w, r, "failed to count posts")
		return
	}

	hasNextPage := offset+limit < totalCount

	Templates["list_posts.html"].ExecuteTemplate(w, "layout.html", map[string]any{
		"Posts":       posts,
		"Page":        page,
		"hasNextPage": hasNextPage,
	})
}

func GetAllTagsHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := DB.Query("SELECT name FROM tags ORDER BY name ASC")
	if err != nil {
		fallbackHandler(w, r, "Failed to fetch tags")
		return
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			fallbackHandler(w, r, "Failed to scan tag")
			return
		}
		tags = append(tags, tag)
	}

	Templates["tags.html"].ExecuteTemplate(w, "layout.html", map[string]any{
		"Tags": tags,
	})
}

func GetPostsByTagHandler(w http.ResponseWriter, r *http.Request) {
	tag := strings.TrimPrefix(r.URL.Path, "/tag/")
	if len(tag) == 0 {
		fallbackHandler(w, r, "Tag not specified")
		return
	}

	query := `
		SELECT posts.title, posts.path, posts.date
		FROM posts
		INNER JOIN post_tags ON posts.id = post_tags.post_id
		INNER JOIN tags ON post_tags.tag_id = tags.id
		WHERE tags.name = ?
		ORDER BY posts.date DESC
	`
	rows, err := DB.Query(query, tag)
	if err != nil {
		fallbackHandler(w, r, "Failed to fetch posts for tag")
		return
	}
	defer rows.Close()

	var posts []map[string]string
	for rows.Next() {
		var title, path, date string
		if err := rows.Scan(&title, &path, &date); err != nil {
			fallbackHandler(w, r, "Failed to scan post")
			return
		}

		path = strings.Replace(path, "posts/", "", 1)
		path = strings.Replace(path, ".org", "", 1)
		posts = append(posts, map[string]string{
			"Title": title,
			"Path":  path,
			"Date":  date,
		})
	}

	Templates["posts_by_tag.html"].ExecuteTemplate(w, "layout.html", map[string]any{
		"Tag":      tag,
		"Posts":    posts,
		"HasPosts": len(posts) > 0,
	})
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		fallbackHandler(w, r, "")
		return
	}
	Templates["index.html"].ExecuteTemplate(w, "layout.html", nil)
}

func AboutHandler(w http.ResponseWriter, r *http.Request) {
	Templates["about.html"].ExecuteTemplate(w, "layout.html", nil)
}

func fallbackHandler(w http.ResponseWriter, _ *http.Request, optionalMessage string) {
	if optionalMessage == "" {
		optionalMessage = "this page does not exist"
	}
	Templates["404.html"].ExecuteTemplate(w, "layout.html", map[string]any{
		"Message": optionalMessage,
	})
}
