package internal

import "io"

type HtmlParser interface {
	MaxPage(r io.Reader) (uint16, error)
	Links(r io.Reader) ([]string, error)
}
