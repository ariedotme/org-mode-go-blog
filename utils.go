package main

import (
	"fmt"
	"html/template"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/niklasfasching/go-org/org"
)

func GetPost(filePath string) (Post, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Post{}, err
	}

	orgDocument := org.New().Parse(file, "")
	parsedContent, err := orgDocument.Write(org.NewHTMLWriter())

	content, err := os.ReadFile(filePath)
	if err != nil {
		return Post{}, err
	}
	title, err := ExtractTitle(string(content))
	if err != nil {
		return Post{}, err
	}
	tags, err := ExtractTags(string(content))
	if err != nil {
		return Post{}, err
	}
	date, err := ExtractDate(string(content))
	if err != nil {
		return Post{}, err
	}

	parsedContent = RemoveHTMLCompiledTitle(parsedContent)

	return Post{title: title, content: template.HTML(parsedContent), date: date, tags: tags}, nil
}

func ExtractDate(content string) (time.Time, error) {
	re := regexp.MustCompile(`\#\+DATE:\s(\d{4}-\d{2}-\d{2})`)
	matches := re.FindStringSubmatch(content)

	if len(matches) > 1 {
		return time.Parse("2006-01-02", matches[1])
	}
	return time.Now(), fmt.Errorf("date not found in file")
}

func ExtractTags(content string) ([]string, error) {
	re := regexp.MustCompile(`\#\+TAGS:\s(.+)`)
	matches := re.FindStringSubmatch(content)

	if len(matches) > 1 {
		match := strings.Replace(matches[1], ",", "", -1)
		return strings.Split(match, " "), nil
	}

	return []string{}, fmt.Errorf("tags not found in file")
}

func ExtractTitle(content string) (string, error) {
	re := regexp.MustCompile(`\#\+TITLE:\s(.+)`)
	matches := re.FindStringSubmatch(content)

	if len(matches) > 1 {
		return matches[1], nil
	}

	return "", fmt.Errorf("title not found in file")
}

func RemoveHTMLCompiledTitle(content string) string {
	re := regexp.MustCompile(`(?m)^\s*<h1(.*)>.*<\/h1>\s*$`)
	return re.ReplaceAllString(content, "")
}
