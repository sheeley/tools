package grabnytimesrecipe

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Input struct {
	Verbose         bool
	CreateBearEntry bool
	URLs            []string
}

type Output struct {
	Results map[string]*Recipe
}

type Recipe struct {
	Title string
	Body  string
}

func GrabNytimesRecipe(in *Input) (*Output, error) {
	o := &Output{
		Results: make(map[string]*Recipe, len(in.URLs)),
	}
	for _, url := range in.URLs {
		o.Results[url] = processUrl(url)
	}
	return o, nil
}

func processUrl(url string) *Recipe {
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

	template := `## Ingredients
%s

## Instructions
%s

#recipes/not made yet#
`
	body := fmt.Sprintf(template, makeList("+", ingredients), makeList("*", instructions))
	return &Recipe{
		Title: title,
		Body:  body,
	}
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
