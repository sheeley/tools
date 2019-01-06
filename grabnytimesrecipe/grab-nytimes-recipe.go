package grabnytimesrecipe

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Input struct {
	Verbose bool
	URLs    []string
}

type Output struct {
}

func GrabNytimesRecipe(in *Input) (*Output, error) {
	for _, url := range in.URLs {
		processUrl(url)
	}
	return &Output{}, nil
}

func processUrl(url string) {
	qIdx := strings.Index(url, "?")
	if qIdx > -1 {
		url = url[0:qIdx]
	}

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		panic(err)
	}

	title := strings.TrimSpace(doc.Find(".recipe-title").Text())

	var ingredients []string
	doc.Find(".recipe-ingredients li").Each(func(_ int, s *goquery.Selection) {
		t := cleanString(s.Text())
		if strings.HasPrefix(t, "Nutritional Information") || strings.HasPrefix(t, "Nutritional analysis") {
			return
		}
		ingredients = append(ingredients, t)
	})

	var instructions []string
	doc.Find(".recipe-steps-wrap li").Each(func(_ int, s *goquery.Selection) {
		t := cleanString(s.Text())
		instructions = append(instructions, t)
	})

	fmt.Printf(`
# %s

## Ingredients
%s

## Instructions
%s

#recipes/not made yet#
`, title, makeList("+", ingredients), makeList("*", instructions))
}

func makeList(sep string, entries []string) string {
	return sep + " " + strings.Join(entries, "\n"+sep+" ")
}

func cleanString(s string) string {
	s = strings.TrimSpace(s)
	s = strings.Replace(s, "\n", "", -1)
	return r.ReplaceAllString(s, " ")
}

var r = regexp.MustCompile(`\s+`)
