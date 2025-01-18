package main

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDatabase() error {
	db, err := sql.Open("sqlite3", "./ariextech.db")

	if err != nil {
		return err
	}
	DB = db

	if err = createPostsTable(); err != nil {
		return err
	}
	if err = createTagsTable(); err != nil {
		return err
	}
	if err = createPostTagsTable(); err != nil {
		return err
	}

	err = populateDatabase()
	if err != nil {
		return err
	}
	return nil
}

func populateDatabase() error {
	err := filepath.WalkDir("./posts", func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if filepath.Ext(path) == ".org" {
			post, err := GetPost(path)
			if err != nil {
				return err
			}

			var exists bool
			err = DB.QueryRow("SELECT EXISTS (SELECT 1 FROM posts WHERE path = ?)", path).Scan(&exists)
			if err != nil {
				return err
			}

			if !exists {
				// Insert post
				result, err := DB.Exec("INSERT INTO posts (title, path, date, tags) VALUES (?, ?, ?, ?)",
					post.title, path, post.date.Format("2006-01-02 15:04:05"), strings.Join(post.tags, ","))
				if err != nil {
					return err
				}

				postID, _ := result.LastInsertId()

				for _, tag := range post.tags {
					var tagID int64
					err := DB.QueryRow("SELECT id FROM tags WHERE name = ?", tag).Scan(&tagID)
					if err == sql.ErrNoRows {
						tagResult, err := DB.Exec("INSERT INTO tags (name) VALUES (?)", tag)
						if err != nil {
							return err
						}
						tagID, _ = tagResult.LastInsertId()
					}

					_, err = DB.Exec("INSERT INTO post_tags (post_id, tag_id) VALUES (?, ?)", postID, tagID)
					if err != nil {
						return err
					}
				}
			}
		}
		return nil
	})
	return err
}

func createPostsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS posts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			path TEXT NOT NULL,
			date DATETIME NOT NULL,
			tags TEXT NOT NULL
		);
	`
	_, err := DB.Exec(query)
	return err
}

func createTagsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS tags (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE
		);
	`
	_, err := DB.Exec(query)
	return err
}

func createPostTagsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS post_tags (
			post_id INTEGER NOT NULL,
			tag_id INTEGER NOT NULL,
			FOREIGN KEY (post_id) REFERENCES posts (id),
			FOREIGN KEY (tag_id) REFERENCES tags (id),
			PRIMARY KEY (post_id, tag_id)
		);
	`
	_, err := DB.Exec(query)
	return err
}

func GetPaginatedPosts(limit, offset int) ([]map[string]any, error) {
	rows, err := DB.Query("SELECT title, path, date, tags FROM posts ORDER BY date DESC LIMIT ? OFFSET ?", limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []map[string]any
	for rows.Next() {
		var title, path, date, tags string
		err := rows.Scan(&title, &path, &date, &tags)
		if err != nil {
			return nil, err
		}

		parsedDate, err := time.Parse("2006-01-02T15:04:05Z", date)
		if err != nil {
			return []map[string]any{}, err
		}
		date = parsedDate.Format("2006-01-02")
		path = strings.Replace(path, "posts", "post", 1)
		path = strings.Replace(path, ".org", "", 1)
		posts = append(posts, map[string]any{
			"Title": title,
			"Path":  path,
			"Date":  date,
			"Tags":  strings.Split(tags, ","),
		})
	}

	return posts, nil
}

func GetTotalPostCount() (int, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM posts").Scan(&count)
	return count, err
}

func GetPostsByTag(tagName string, limit, offset int) ([]map[string]any, error) {
	query := `
		SELECT posts.title, posts.path, posts.date, posts.tags
		FROM posts
		INNER JOIN post_tags ON posts.id = post_tags.post_id
		INNER JOIN tags ON post_tags.tag_id = tags.id
		WHERE tags.name = ?
		ORDER BY posts.date DESC
		LIMIT ? OFFSET ?
	`

	rows, err := DB.Query(query, tagName, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []map[string]any
	for rows.Next() {
		var title, path, date, tags string
		err := rows.Scan(&title, &path, &date, &tags)
		if err != nil {
			return nil, err
		}

		posts = append(posts, map[string]any{
			"Title": title,
			"Path":  path,
			"Date":  date,
			"Tags":  strings.Split(tags, ","),
		})
	}

	return posts, nil
}
