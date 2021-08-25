package dict

import (
	"fmt"
	"github.com/gocolly/colly"
	"log"
	"strings"
)

type SearchResults struct {
	ready   bool
	Query   string
	Results []Result
}

type Result struct {
	Word     string
	Meanings []Meaning
	Synonyms []string
	Contexts []Context
}

type Meaning struct {
	Author  string
	Meaning string
}

type Context struct {
	Text   string
	Source string
}

var dictURL string = "https://rustxt.ru/dict/"

func ParseWordData(word string, ch chan *Result) {

	c := colly.NewCollector()

	c.OnRequest(func(r *colly.Request) {
		fmt.Printf("Visiting: %s", dictURL+strings.ToLower(word))
	})

	var result *Result

	c.OnHTML(".col-md-8.col-12", func(e *colly.HTMLElement) {
		var meanings []Meaning
		e.ForEach(".answer.border.p-2", func(i int, e *colly.HTMLElement) {
			var meaning string
			var author string

			e.ForEach(".mb-1", func(i int, e *colly.HTMLElement) {
				meaning += e.Text
			})
			author = e.ChildText(".d-block.text-right.mb-0.mt-1")

			meanings = append(meanings, Meaning{
				Meaning: meaning,
				Author:  author,
			})
		})

		var synonyms []string
		e.ForEach(".wrap", func(i int, e *colly.HTMLElement) {
			e.ForEach("li", func(i int, e *colly.HTMLElement) {
				synonyms = append(synonyms, e.Text)
			})
		})

		var contexts []Context
		e.ForEach(".row.mb-4.sentence", func(i int, e *colly.HTMLElement) {
			contexts = append(contexts, Context{
				Text:   e.ChildText(".mb-2.sentence-text"),
				Source: e.ChildText(".text-muted .text-muted"),
			})
		})

		result = &Result{
			Word:     word,
			Meanings: meanings,
			Synonyms: synonyms,
			Contexts: contexts,
		}

		ch <- result
	})

	if err := c.Visit(dictURL + strings.ToLower(word)); err != nil {
		log.Fatalln("Incorrect request")
	}
}
