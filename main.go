package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/dimchansky/utfbom"
	"github.com/jawher/mow.cli"
	md "gopkg.in/russross/blackfriday.v2"
)

func main() {
	app := cli.App("markdown", "")

	app.Command("note", "", func(note *cli.Cmd) {

		note.Command("add", "", func(add *cli.Cmd) {

			var (
				targetFile = add.StringOpt("f file", "", "file")
			)

			if *targetFile == "" {
				dirName, fileName := generateNames()
				_ = os.MkdirAll(dirName, os.ModePerm)
				*targetFile = fileName
			}

			add.Action = func() {

				if err := appendText(*targetFile, []byte("some text foobar")); err != nil {
					fmt.Fprintln(os.Stderr, "Error writing file text", err)
					os.Exit(-1)
				}
			}

		})

	})

	app.Command("html", "", func(html *cli.Cmd) {
		var (
			inputFile = html.StringOpt("f file", "", "file to convert")
		)

		html.Action = func() {
			var (
				input []byte
			)

			if *inputFile == "" {
				input = readStdin()
			} else {
				input = readFile(*inputFile)
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

func readStdin() []byte {
	var (
		input []byte
		err   error
	)
	if input, err = ioutil.ReadAll(os.Stdin); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading from Stdin:", err)
		os.Exit(-1)
	}
	return input
}

func readFile(fileName string) []byte {
	var (
		input []byte
		err   error
	)
	if input, err = ioutil.ReadFile(fileName); err != nil {
		fmt.Fprintln(os.Stderr, "Error reading from", fileName, ":", err)
		os.Exit(-1)
	}
	if input, err = ioutil.ReadAll(utfbom.SkipOnly(bytes.NewReader(input))); err != nil {
		fmt.Fprintln(os.Stderr, "Error removing bom", ":", err)
		os.Exit(-1)
	}
	return input
}

func generateNames() (dirName string, fileName string) {
	now := time.Now()
	_, week := now.ISOWeek()
	dirName = filepath.Join("notes", fmt.Sprintf("%s-w%d", now.Format("2006-01"), week))
	fileName = filepath.Join(dirName, now.Format("2006-01-02.md"))
	return
}

func fileExist(fileName string) bool {
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		return false
	}
	return true
}

func appendText(fileName string, input []byte) error {
	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if _, err := file.Write(input); err != nil {
		return err
	}
	return file.Close()
}
