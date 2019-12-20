package main

import (
	"log"
	"strings"

	"github.com/DavidBelicza/TextRank"
	"github.com/DavidBelicza/TextRank/convert"
	"github.com/DavidBelicza/TextRank/parse"
	"github.com/PuerkitoBio/goquery"
)

func main() {
	urls := ScrapeUrls("huawei p30 review")

	for _, u := range urls {
		paragraphs, err := pScrape(u)
		if err != nil {
			log.Printf("%v for %v", err, u)
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

		var ranked []string
		// Put just the sentences in slice
		for _, s := range sentences {
			ranked = append(ranked, strings.TrimSpace(s.Value))
		}
		log.Printf("%v for %v", ranked[3], u)
	}

}

// TODO 1
func ScrapeUrls(_ string) []string {
	return []string{"https://www.androidcentral.com/huawei-p30-pro-review-3-months-later",
		"https://www.digitaltrends.com/cell-phone-reviews/huawei-p30-pro-review/"}
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

// paragraph into parse.Text struct
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

// TODO 2
// Write to html
func (urls,rankedText []string){}