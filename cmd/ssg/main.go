package main

import (
	"github.com/nasermirzaei89/ssg/internal/ssg"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "Static Site Generator"
	app.Description = "Generates static site from simple markdown files"
	app.Copyright = "@Copyleft All wrongs reserved"

	app.Commands = append(app.Commands, &cli.Command{
		Name:        "generate",
		Description: "generates static site from markdowns",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "path",
				Value: ".",
				Usage: "root path of site repo",
			},
			&cli.StringFlag{
				Name:  "dist",
				Value: "dist",
				Usage: "path to generate in",
			},
			&cli.StringFlag{
				Name:  "theme",
				Value: "default",
				Usage: "theme name to generate with",
			},
		},
		Action: func(ctx *cli.Context) error {
			err := ssg.Generate(ctx.String("path"), ctx.String("dist"), ctx.String("theme"))
			if err != nil {
				return errors.Wrap(err, "error on generate")
			}

			return nil
		},
	})

	app.Commands = append(app.Commands, &cli.Command{
		Name:        "serve",
		Description: "serves generated static site",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "path",
				Value: ".",
				Usage: "root path of site repo",
			},
			&cli.StringFlag{
				Name:  "port",
				Value: "8080",
				Usage: "port of serve",
			},
		},
		Action: func(ctx *cli.Context) error {
			err := ssg.Serve(ctx.String("path"), ctx.String("port"))
			if err != nil {
				return errors.Wrap(err, "error on serve")
			}

			return nil
		},
	})

	err := app.Run(os.Args)
	if err != nil {
		log.Fatalln(errors.Wrap(err, "error on run cli"))
	}
}
