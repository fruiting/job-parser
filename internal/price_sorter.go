package internal

//go:generate mockgen -source=price_sorter.go -destination=./price_sorter_mock.go -package=internal

// PriceSorter sorts prices
type PriceSorter interface {
	// PricesFromJobs returns min, max and median salaries
	PricesFromJobs(jobs []*Job) (Salary, Salary, Salary)
}
