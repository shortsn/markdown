package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/dimchansky/utfbom"
)

// Reads all .txt files in the current folder
// and encodes them as strings literals in textfiles.go
func main() {
	sourceDir := "./static"
	fs, _ := ioutil.ReadDir(sourceDir)
	out, _ := os.Create("static.go")
	out.Write([]byte("package main \n\nconst (\n"))
	for _, f := range fs {
		fileName := f.Name()
		out.Write([]byte(strings.TrimSuffix(fileName, filepath.Ext(fileName)) + " = `"))
		f, _ := os.Open(filepath.Join(sourceDir, fileName))

		fmt.Fprintln(os.Stdout, fileName)
		io.Copy(out, utfbom.SkipOnly(f))
		out.Write([]byte("`\n"))
	}
	out.Write([]byte(")\n"))
}
