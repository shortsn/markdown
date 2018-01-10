package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dimchansky/utfbom"
	"github.com/jawher/mow.cli"
	md "gopkg.in/russross/blackfriday.v2"
)

func main() {
	app := cli.App("markdown", "")

	app.Command("html", "", func(cmd *cli.Cmd) {
		var (
			inputFile = cmd.StringOpt("f file", "", "file to convert")
		)

		cmd.Action = func() {

			var (
				input []byte
				err   error
			)

			switch *inputFile {
			case "":
				if input, err = ioutil.ReadAll(os.Stdin); err != nil {
					fmt.Fprintln(os.Stderr, "Error reading from Stdin:", err)
					os.Exit(-1)
				}
			default:
				if input, err = ioutil.ReadFile(*inputFile); err != nil {
					fmt.Fprintln(os.Stderr, "Error reading from", *inputFile, ":", err)
					os.Exit(-1)
				}
				input, _ = ioutil.ReadAll(utfbom.SkipOnly(bytes.NewReader(input)))
			}

			extensions := md.WithExtensions(md.CommonExtensions)
			htmlRenderer := md.WithRenderer(
				md.NewHTMLRenderer(
					md.HTMLRendererParameters{
						Flags: md.CommonHTMLFlags,
					}))
			output := md.Run(input, extensions, htmlRenderer)
			fmt.Println(string(output))
			os.Exit(0)
		}

	})

	app.Run(os.Args)
}
