package pricesorter

import "fruiting/job-parser/internal"

// todo rename
type PriceSorter struct {
}

func NewPriceSorter() *PriceSorter {
	return &PriceSorter{}
}

func (s *PriceSorter) PricesFromJobs(jobs []*internal.Job) (internal.Salary, internal.Salary, internal.Salary) {
	return 0, 0, 0
}
