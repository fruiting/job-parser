package parser

import (
	"fmt"
	"io"
	"net/http"

	"fruiting/job-parser/internal"
)

type GeneralParser struct {
	parser internal.JobsParser
}

func NewGeneralParser() *GeneralParser {
	return &GeneralParser{}
}

func (p *GeneralParser) Run(position internal.Name) (*internal.JobsInfo, error) {
	generalParseDom, err := p.parse(p.parser.GeneralLink(position))
	if err != nil {
		return nil, fmt.Errorf("can't make general parse: %w", err)
	}

	pagesCount, err := p.parser.PagesCount(generalParseDom)
	if err != nil {
		return nil, fmt.Errorf("can't get pages count: %w", err)
	}

	firstPageLinks, err := p.parser.Links(generalParseDom)
	if err != nil {
		return nil, fmt.Errorf("can't get links in general parse: %w", err)
	}

	links := make(chan []string, 0)
	links <- firstPageLinks

	errChan := make(chan error)
	for i := uint16(0); i < pagesCount; i++ {
		go func(i uint16) {
			//todo вьебать ретраи
			pageDom, err := p.parse(p.parser.SearchPageLink(i + 1))
			if err != nil {
				errChan <- fmt.Errorf("can't parse page %d: %w", i+1, err)
			}

			nextPageLinks, err := p.parser.Links(pageDom)
			if err != nil {
				errChan <- fmt.Errorf("can't get links in page %d: %w", i+1, err)
			}

			links <- nextPageLinks
		}(i)
	}

	jobsInfo := make([]*internal.Job, 0, int(pagesCount*p.parser.ItemsCount()))
	for {
		select {
		case linksBatch := <-links:
			for _, link := range linksBatch {
				detailPageDom, err := p.parse(link)
				if err != nil {
					return nil, fmt.Errorf("can't get detail page dom %s: %w", link, err)
				}

				jobInfo, err := p.parser.ParseDetail(detailPageDom)
				if err != nil {
					return nil, fmt.Errorf("can't parse detail page dom %s: %w", link, err)
				}

				jobsInfo = append(jobsInfo, jobInfo)
			}
		case <-errChan:
			return nil, fmt.Errorf("can't parse link: %w", err)
		}
	}

	return &internal.JobsInfo{
		PositionToParse: position,
		MinSalary:       0,   // todo sort
		MaxSalary:       0,   // todo sort
		AverageSalary:   0,   // todo sort
		PopularSkills:   nil, // todo sort
		Jobs:            jobsInfo,
	}, nil
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
