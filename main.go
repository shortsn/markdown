package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dimchansky/utfbom"
	"github.com/jawher/mow.cli"
	md "gopkg.in/russross/blackfriday.v2"
)

func main() {
	executableFullname, _ := os.Executable()
	executableBasename := filepath.Base(executableFullname)
	app := cli.App(executableBasename, "")
	app.Version("v version", "0.0.1")

	app.Command("note", "", func(note *cli.Cmd) {

		note.Command("add", "", func(add *cli.Cmd) {

			var (
				targetFile = add.StringOpt("o output", "", "output file")
			)

			if *targetFile == "" {
				dirName, fileName := generateNames()
				_ = os.MkdirAll(dirName, os.ModePerm)
				*targetFile = fileName
			}

			add.Action = func() {

				if err := appendText(*targetFile, readStdin()); err != nil {
					fmt.Fprintln(os.Stderr, "Error writing file text", err)
					os.Exit(-1)
				}

				fmt.Fprintln(os.Stdout, *targetFile)
				os.Exit(0)
			}

		})

	})

	app.Command("convert", "", func(convert *cli.Cmd) {

		convert.Spec = "[-s | FILES...]"

		var (
			fromStdin  = convert.BoolOpt("s stdin", false, "")
			inputFiles = convert.StringsArg("FILES", []string{}, "files to convert")
		)

		convert.Action = func() {

			if *fromStdin {
				input := convertToHTML(readStdin())
				fmt.Fprintln(os.Stdout, string(input))
				os.Exit(0)
			}

			for _, inputFile := range *inputFiles {
				inputFile, _ = filepath.Abs(inputFile)
				outputFile := strings.Replace(inputFile, filepath.Ext(inputFile), ".html", -1)
				content := readFile(inputFile)

				if err := ioutil.WriteFile(outputFile, content, 0644); err != nil {
					fmt.Fprintln(os.Stderr, "Error writing file", err)
					os.Exit(-1)
				}

				fmt.Fprintln(os.Stdout, outputFile)
			}

			os.Exit(0)
		}

	})

	app.Run(os.Args)
}

func convertToHTML(input []byte) []byte {
	extensions := md.WithExtensions(md.CommonExtensions)
	htmlRenderer := md.WithRenderer(
		md.NewHTMLRenderer(
			md.HTMLRendererParameters{
				Flags: md.CommonHTMLFlags,
			}))
	return md.Run(input, extensions, htmlRenderer)
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
