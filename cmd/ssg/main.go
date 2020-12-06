package main

import (
	"flag"
	"fmt"
	"github.com/nasermirzaei89/ssg/internal/ssg"
	"log"
)

func main() {
	pathFlag := flag.String("path", ".", "root path of site repo")
	distFlag := flag.String("dist", "dist", "path to generate in")
	themeFlag := flag.String("theme", "default", "theme name to generate with")
	portFlag := flag.String("port", "8080", "port of serve")
	flag.Parse()

	switch {
	case flag.Arg(0) == "generate":
		err := ssg.Generate(*pathFlag, *distFlag, *themeFlag)
		if err != nil {
			log.Fatalln(fmt.Errorf("error on generate: %w", err))
		}
	case flag.Arg(0) == "serve":
		err := ssg.Serve(*pathFlag, *portFlag)
		if err != nil {
			log.Fatalln(fmt.Errorf("error on serve: %w", err))
		}
	case flag.Arg(0) == "version":
		err := ssg.PrintVersion()
		if err != nil {
			log.Fatalln(fmt.Errorf("error on print version: %w", err))
		}
	default:
		err := ssg.PrintHelp()
		if err != nil {
			log.Fatalln(fmt.Errorf("error on print help: %w", err))
		}
	}
}
