package ssg

import (
	"fmt"
	"net/http"
	"os"
	"path"
)

func Serve(dir, port string) error {
	baseDir := path.Clean(dir)

	distDir := path.Join(baseDir, "dist")
	_, err := os.Stat(distDir)
	if os.IsNotExist(err) {
		distDir = baseDir
	}

	fs := http.FileServer(http.Dir(distDir))

	http.Handle("/", fs)

	fmt.Printf("Serving at http://localhost:%s\n", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		return fmt.Errorf("error on listen and serve http: %w", err)
	}

	return nil
}
