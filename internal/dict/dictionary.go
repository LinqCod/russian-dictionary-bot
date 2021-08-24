package dict

import (
	"fmt"
	"github.com/gocolly/colly"
	"io"
	"net/http"
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

func DictionaryAPI(word string) string {

	var answer string

	if response, err := http.Get(dictURL + strings.ToLower(word)); err != nil {
		answer = "Dictionary is not respond"
	} else {
		defer response.Body.Close()

		//TODO:parse web page and read data
		//results, err :=
		ParseWordData(response.Body, word)
		if err != nil {
			answer = "Word not found"
		} else {
			//TODO: insert word to database
		}
	}

	return answer
}

//(SearchResults, error)
func ParseWordData(body io.ReadCloser, word string) {
	c := colly.NewCollector()

	c.OnHTML(".row", func(e *colly.HTMLElement) {

		var meanings []Meaning
		e.ForEach(".p-2", func(i int, e *colly.HTMLElement) {
			var meaning string
			var author string

			e.ForEach("mb-1", func(i int, e *colly.HTMLElement) {
				meaning += e.Text
			})
			author = e.ChildText(".mt-1")

			meanings = append(meanings, Meaning{
				Meaning: meaning,
				Author:  author,
			})
		})

		var synonyms []string
		e.ForEach("li", func(i int, e *colly.HTMLElement) {
			synonyms = append(synonyms, e.Text)
		})

		var contexts []Context
		e.ForEach(".row .mb-4 .sentence", func(i int, e *colly.HTMLElement) {
			contexts = append(contexts, Context{
				Text:   e.ChildText(".sentence-text"),
				Source: e.ChildText(".text-muted"),
			})
		})

		result := Result{
			Word:     word,
			Meanings: meanings,
			Synonyms: synonyms,
			Contexts: contexts,
		}

		fmt.Println(result)
	})
}
