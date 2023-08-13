package internal

import (
	"time"
)

type Job struct {
	PositionName Name
	Link         Link
	Salary       Salary
	Skills       skills
}

//easyjson:json JobsInfo
type JobsInfo struct {
	PositionsToParse []Name    `json:"position_to_parse"`
	MinSalary        Salary    `json:"min_salary"`
	MaxSalary        Salary    `json:"max_salary"`
	MedianSalary     Salary    `json:"median_salary"`
	PopularSkills    skills    `json:"popular_skills"`
	Parser           Parser    `json:"parser"`
	Jobs             []*Job    `json:"jobs"`
	Time             time.Time `json:"time"`
}
