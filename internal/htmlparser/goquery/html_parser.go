package goquery

import (
	"fmt"
	"io"
	"strconv"

	"github.com/PuerkitoBio/goquery"
)

type HtmlParser struct {
}

func NewHtmlParser() *HtmlParser {
	return &HtmlParser{}
}

func (p *HtmlParser) MaxPage(r io.Reader) (uint16, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return 0, fmt.Errorf("can't init goquery document: %w", err)
	}

	var maxPageStr string
	doc.Find(".bloko-button[data-qa=\"pager-page\"]").Each(func(i int, s *goquery.Selection) {
		maxPageStr = s.Nodes[0].LastChild.LastChild.Data
	})

	maxPage, err := strconv.ParseUint(maxPageStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("can't parse uint: %w", err)
	}

	return uint16(maxPage), nil
}

func (p *HtmlParser) Links(r io.Reader) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("can't init goquery document: %w", err)
	}

	links := make([]string, 0)
	doc.Find(".serp-item__title").Each(func(i int, s *goquery.Selection) {
		val, ok := s.Attr("href")
		if !ok {
			return
		}

		links = append(links, val)
	})

	return links, nil
}
