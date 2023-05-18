package internal

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/avast/retry-go"
)

type GeneralParser struct {
	Parser JobsParser
}

func NewGeneralParser(parser JobsParser) *GeneralParser {
	return &GeneralParser{
		Parser: parser,
	}
}

func (p *GeneralParser) Run(position Name) ([]*Job, error) {
	return nil, nil
	generalParseDom, err := p.parse(p.Parser.GeneralLink(position))
	if err != nil {
		return nil, fmt.Errorf("can't make general parse: %w", err)
	}

	pagesCount, err := p.Parser.PagesCount(generalParseDom)
	if err != nil {
		return nil, fmt.Errorf("can't get pages count: %w", err)
	}

	firstPageLinks, err := p.Parser.Links(generalParseDom)
	if err != nil {
		return nil, fmt.Errorf("can't get links in general parse: %w", err)
	}

	links := make(chan []string, 0)
	links <- firstPageLinks

	errChan := make(chan error)
	for i := uint16(0); i < pagesCount; i++ {
		go func(i uint16) {
			err := retry.Do(
				func() error {
					pageDom, err := p.parse(p.Parser.SearchPageLink(i + 1))
					if err != nil {
						return fmt.Errorf("can't parse page %d: %w", i+1, err)
					}

					nextPageLinks, err := p.Parser.Links(pageDom)
					if err != nil {
						return fmt.Errorf("can't get links in page %d: %w", i+1, err)
					}

					links <- nextPageLinks
					return nil
				},
				retry.Attempts(5),
				retry.Delay(1*time.Second),
			)
			if err != nil {
				errChan <- fmt.Errorf("can't parse: %w", err)
			}
		}(i)
	}

	jobsInfo := make([]*Job, 0, int(pagesCount*p.Parser.ItemsCount()))
	for {
		select {
		case linksBatch := <-links:
			for _, link := range linksBatch {
				err := retry.Do(
					func() error {
						detailPageDom, err := p.parse(link)
						if err != nil {
							return fmt.Errorf("can't get detail page dom %s: %w", link, err)
						}

						jobInfo, err := p.Parser.ParseDetail(detailPageDom)
						if err != nil {
							return fmt.Errorf("can't parse detail page dom %s: %w", link, err)
						}

						jobsInfo = append(jobsInfo, jobInfo)
						return nil
					},
					retry.Attempts(5),
					retry.Delay(1*time.Second),
				)
				if err != nil {
					return nil, fmt.Errorf("can't parse detail: %w", err)
				}
			}
		case <-errChan:
			return nil, fmt.Errorf("can't parse link: %w", err)
		}
	}

	return jobsInfo, nil
}

func (p *GeneralParser) parse(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("can't get url info: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http get is not ok, it has %d code", resp.StatusCode)
	}

	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("can't read response body: %w", err)
	}

	return string(content), nil
}
