package ttpr

import (
	"io/ioutil"
	"log"
	"strings"
	"sync"

	"github.com/DavidBelicza/TextRank"
	"github.com/DavidBelicza/TextRank/convert"
	"github.com/DavidBelicza/TextRank/parse"
	"github.com/PuerkitoBio/goquery"
)

type Result struct {
	Url  string
	Summ []string
}

func Search(term string, linkCount int) []Result {
	urls := scrapeDDG(term)
	results := make([]Result, linkCount)

	var wg sync.WaitGroup
	for i, u := range urls[0:linkCount] {
		results[i].Url = u
		wg.Add(1)
		go GetRanked(&results[i], &wg)
	}
	wg.Wait()
	return results
}

func GetRanked(r *Result, wg *sync.WaitGroup) {
	defer wg.Done()

	paragraphs, err := pScrape(r.Url)
	if err != nil {
		log.Printf("%v for %v", err, r.Url)
		return
	}

	tr := textrank.NewTextRank()
	// Default Rule for parsing.
	//rule := textrank.NewDefaultRule()
	// Default Language for filtering stop words.
	language := textrank.NewDefaultLanguage()
	// Default algorithm for ranking text.
	algorithmDef := textrank.NewDefaultAlgorithm()

	// preparing *rank.Rank object for ranking
	text := parseText(paragraphs)
	for _, sentence := range text.GetSentences() {
		convert.TextToRank(sentence, language, tr.GetRankData())
	}
	// Run the ranking.
	tr.Ranking(algorithmDef)
	// Get the most important 4 sentences.
	sentences := textrank.FindSentencesByRelationWeight(tr, 4)
	// Put just the sentences in slice
	for _, s := range sentences {
		r.Summ = append(r.Summ, strings.TrimSpace(s.Value))
	}
}

func scrapeDDG(term string) []string {
	ddgPrefix := "https://duckduckgo.com/html?q="
	suffix := strings.Join(strings.Fields(term), "+")
	searchLink := ddgPrefix + suffix

	doc, err := goquery.NewDocument(searchLink)
	if err != nil {
		log.Fatal(err)
	}

	var urls []string
	results := doc.Find(".result")
	if len(results.Nodes) <= 0 {
		log.Println("goquery.Find give no results") // Watch for html/css structure change
	}
	for i := range results.Nodes {
		u := results.Eq(i).Find(".result__url").Text()
		urls = append(urls, "https://"+strings.TrimSpace(u))
	}
	return urls

}

// number of characters in p element to consider a content.
// to remove stuffs like ads and attribution in p.
const paraLimit = 175

// scrape contents of <p> tag
func pScrape(url string) ([]string, error) {
	var items []string

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return items, err
	}
	// find p tag element
	doc.Find("p").Each(func(i int, s *goquery.Selection) {
		paragraph := strings.TrimSpace(s.Text())
		lastDot := strings.LastIndex(paragraph, ".")
		// Remove insufficient length paragraph and cut string after last fullstop
		// todo: fix getting tripped by decimal
		if lastDot >= paraLimit {
			item := string(paragraph[0 : lastDot+1])
			items = append(items, item)

		}

	})
	return items, err
}

// put paragraph into parse.Text struct
func parseText(paragraphs []string) parse.Text {
	text := parse.Text{}
	for _, p := range paragraphs {
		// Split words from sentence
		for _, i := range strings.SplitAfter(p, ". ") {
			text.Append(i, strings.Fields(i))
		}
	}
	return text

}

// TODO
// - pScrape <nosript> case
