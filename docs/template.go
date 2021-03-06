package main

import (
	"hash/adler32"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"
)

var shapes = []string{
	"circle",
	"hexagon",
	"pentagon",
	"square",
	"triangle",
}

var colours = []string{
	"b", // blue
	"c", // cyan
	"g", // green
	"p", // pink
	"r", // red
	"t", // turquoise
	"v", // violet
	"y", // yellow
}

var rotations = []string{
	"rotate1",
	"rotate2",
	"rotate3",
	"rotate4",
	"rotate5",
	"rotate6",
	"rotate7",
}

var pageTitles = map[string]string{
	"advanced.html":         "Advanced Please",
	"acknowledgements.html": "Acknowledgements",
	"basics.html":           "Please basics",
	"cache.html":            "Please caching system",
	"commands.html":         "Please commands",
	"config.html":           "Please config file reference",
	"faq.html":              "Please FAQ",
	"index.html":            "Please",
	"intermediate.html":     "Intermediate Please",
	"language.html":         "The Please BUILD language",
	"lexicon.html":          "Please Lexicon",
	"pleasings.html":        "Extra rules (aka. Pleasings)",
	"quickstart.html":       "Please quickstart",
	"error.html":            "plz op...",
}

func mustRead(filename string) string {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	return string(b)
}

func main() {
	filename := os.Args[2]
	basename := path.Base(filename)
	basenameIndex := int(adler32.Checksum([]byte(basename)))
	modulo := func(s []string, i int) string { return s[(basenameIndex+i)%len(s)] }
	random := func(x, min, max int) int { return (x*basenameIndex+min)%(max-min) + min }
	funcs := template.FuncMap{
		"menuItem": func(s string) string {
			if basename[:len(basename)-5] == s {
				return ` class="selected"`
			}
			return ""
		},
		"shape":        func(i int) string { return modulo(shapes, i) },
		"colour":       func(i int) string { return modulo(colours, i) },
		"rotate":       func(i int) string { return modulo(rotations, i) },
		"random":       func(x, min, max int) int { return (x*basenameIndex+min)%(max-min) + min },
		"randomoffset": func(x, min, max, step int) int { return x*step + random(x, min, max) },
	}
	data := struct {
		Title, Header, Contents string
		SideImages              []int
		Player, IsIndex         bool
	}{
		Title:    pageTitles[basename],
		Header:   mustRead(os.Args[1]),
		Contents: mustRead(filename),
		Player:   basename == "faq.html",
		IsIndex:  basename == "index.html",
	}
	for i := 0; i <= strings.Count(data.Contents, "\n")/150; i++ {
		// Awkwardly this seems to have to be a slice to range over in the template.
		data.SideImages = append(data.SideImages, i+1)
	}

	tmpl := template.Must(template.New("tmpl").Funcs(funcs).Parse(mustRead(os.Args[1])))
	err := tmpl.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}
}
