package ssg

import (
	"bytes"
	"github.com/otiai10/copy"
	"github.com/pkg/errors"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type PageData struct {
	Title   string
	URL     string
	Content template.HTML
	Type    string
	Layout  string
}

func Generate(dir, distPath, themePath string) error {
	baseDir := path.Clean(dir)

	themeDir := path.Join(baseDir, ".themes", themePath)

	tpls := make(map[string]*template.Template)

	err := filepath.Walk(themeDir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		fileExt := path.Ext(p)
		fileName := strings.Replace(p, themeDir, "", 1)
		fileName = strings.TrimLeft(fileName, "/")
		fileName = fileName[:len(fileName)-len(fileExt)]

		if strings.HasPrefix(fileName, ".") {
			return nil
		}

		if ext := strings.ToLower(path.Ext(p)); ext != ".html" {
			return nil
		}

		tpl, err := template.ParseFiles(p)
		if err != nil {
			return errors.Wrap(err, "error on parse template")
		}

		tpls[fileName] = tpl

		return nil
	})
	if err != nil {
		return errors.Wrap(err, "error on walk filepath")
	}

	var distDir string
	if path.IsAbs(distPath) {
		distDir = path.Clean(distPath)
	} else {
		distDir = path.Join(baseDir, distPath)
	}

	err = os.RemoveAll(distDir)
	if err != nil {
		return errors.Wrap(err, "error on remove directory")
	}

	err = os.MkdirAll(distDir, 0755)
	if err != nil {
		return errors.Wrap(err, "error on make directory")
	}

	pages := make([]string, 0)

	err = filepath.Walk(baseDir, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasPrefix(p, ".") {
			return nil
		}

		if ext := strings.ToLower(path.Ext(p)); ext != ".md" && ext != ".markdown" {
			return nil
		}

		pages = append(pages, p)
		return nil
	})
	if err != nil {
		return errors.Wrap(err, "error on walk filepath")
	}

	for i := range pages {
		b, err := ioutil.ReadFile(pages[i])
		if err != nil {
			return errors.Wrap(err, "error on read file")
		}

		var pageData PageData

		fileExt := path.Ext(pages[i])
		fileName := path.Base(pages[i])
		fileName = fileName[:len(fileName)-len(fileExt)]

		pageData.Title = strings.Replace(fileName, "-", " ", -1)
		pageData.URL = pages[i][:len(pages[i])-len(fileExt)]
		pageData.Type = "page"
		pageData.Layout = "page"

		buf := bytes.NewBuffer(nil)

		md := goldmark.New(
			goldmark.WithExtensions(
				meta.New(),
			),
		)

		ctx := parser.NewContext()

		err = md.Convert(b, buf, parser.WithContext(ctx))
		if err != nil {
			return errors.Wrap(err, "error on convert markdown to html")
		}

		cfg, err := meta.TryGet(ctx)
		if err != nil {
			return errors.Wrap(err, "error on get yaml data from markdown: %w")
		}

		if iv, ok := cfg["title"]; ok {
			if v, ok2 := iv.(string); ok2 {
				pageData.Title = v
			} else {
				return errors.Errorf("invalid type for title field in meta: %T", iv)
			}
		}

		if iv, ok := cfg["permalink"]; ok {
			if v, ok2 := iv.(string); ok2 {
				pageData.URL = v
			} else {
				return errors.Errorf("invalid type for permalink field in meta: %T", iv)
			}
		}

		if iv, ok := cfg["type"]; ok {
			if v, ok2 := iv.(string); ok2 {
				pageData.Type = v
			} else {
				return errors.Errorf("invalid type for type field in meta: %T", iv)
			}
		}

		if iv, ok := cfg["layout"]; ok {
			if v, ok2 := iv.(string); ok2 {
				pageData.Layout = v
			} else {
				return errors.Errorf("invalid type for layout field in meta: %T", iv)
			}
		}

		pageData.Content = template.HTML(buf.String())

		pagePath := path.Join(distDir, pageData.URL, "index.html")

		err = os.MkdirAll(path.Dir(pagePath), 0755)
		if err != nil {
			return errors.Wrap(err, "error on make directory")
		}

		htmlFile, err := os.OpenFile(pagePath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return errors.Wrap(err, "error on open file")
		}

		pageTpl, ok := tpls[pageData.Layout]
		if !ok {
			return errors.Errorf("layout '%s' doesn't exists", pageData.Layout)
		}

		err = pageTpl.Execute(htmlFile, pageData)
		if err != nil {
			return errors.Wrap(err, "error on execute page template")
		}
	}

	// copy static

	staticDir := path.Join(baseDir, "static")

	err = copy.Copy(staticDir, distDir)
	if err != nil {
		return errors.Wrap(err, "error on copy statics")
	}

	return nil
}
