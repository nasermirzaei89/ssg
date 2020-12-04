package ssg

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

type PostData struct {
	Title   string
	URL     string
	Content template.HTML
}

type IndexData struct {
	Title string
	URL   string
	Posts []PostData
}

func Generate(dir, distPath, tplPath string) error {
	baseDir := path.Clean(dir)

	postsDir := path.Join(baseDir, "posts")
	postFiles, err := ioutil.ReadDir(postsDir)
	if err != nil {
		return fmt.Errorf("error on read dir: %w", err)
	}

	tplDir := path.Join(baseDir, ".template", tplPath)

	postTpl, err := template.ParseFiles(path.Join(tplDir, "post.html"))
	if err != nil {
		return fmt.Errorf("error on parse post template: %w", err)
	}

	var distDir string
	if path.IsAbs(distPath) {
		distDir = path.Clean(distPath)
	} else {
		distDir = path.Join(baseDir, distPath)
	}
	err = os.RemoveAll(distDir)
	if err != nil {
		return fmt.Errorf("error on remove directory: %w", err)
	}

	err = os.MkdirAll(distDir, 0755)
	if err != nil {
		return fmt.Errorf("error on make directory: %w", err)
	}

	posts := make([]PostData, 0)

	for i := range postFiles {
		if postFiles[i].IsDir() {
			continue
		}

		fileName := postFiles[i].Name()

		fileExt := path.Ext(fileName)

		if fileExt != ".md" {
			continue
		}

		filePath := path.Join(postsDir, fileName)

		outPath := path.Join(distDir, fileName[:len(fileName)-len(fileExt)], "index.html")

		err := os.MkdirAll(path.Dir(outPath), 0755)
		if err != nil {
			return fmt.Errorf("error on make dir: %w", err)
		}

		htmlFile, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY, 0744)
		if err != nil {
			return fmt.Errorf("error on open file: %w", err)
		}

		defer func() { _ = htmlFile.Close() }()

		var post PostData

		post.Title = strings.Replace(fileName[:len(fileName)-len(fileExt)], "-", " ", -1)

		post.URL = path.Join(strings.Split(path.Dir(outPath), "/")[1:]...)

		ifile, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("error on open file: %w", err)
		}

		buf := bytes.NewBuffer(nil)

		err = markdownToHtml(ifile, buf)
		if err != nil {
			return fmt.Errorf("error on convert markdown to html: %w", err)
		}

		post.Content = template.HTML(buf.String())

		err = postTpl.Execute(htmlFile, post)
		if err != nil {
			return fmt.Errorf("error on execute html template: %w", err)
		}

		posts = append(posts, post)
	}

	indexTpl, err := template.ParseFiles(path.Join(tplDir, "index.html"))
	if err != nil {
		return fmt.Errorf("error on parse index template: %w", err)
	}

	outPath := path.Join(distDir, "index.html")

	err = os.MkdirAll(path.Dir(outPath), 0755)
	if err != nil {
		return fmt.Errorf("error on make dir: %w", err)
	}

	htmlFile, err := os.OpenFile(outPath, os.O_CREATE|os.O_WRONLY, 0744)
	if err != nil {
		return fmt.Errorf("error on open file: %w", err)
	}

	defer func() { _ = htmlFile.Close() }()

	index := IndexData{
		Title: "",
		URL:   "/",
		Posts: posts,
	}

	err = indexTpl.Execute(htmlFile, index)
	if err != nil {
		return fmt.Errorf("error on execute html template: %w", err)
	}

	return nil
}
