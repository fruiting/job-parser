package skillssorter

import "fruiting/job-parser/internal"

// todo rename
type SkillsSorter struct {
}

func NewSkillsSorter() *SkillsSorter {
	return &SkillsSorter{}
}

func (s *SkillsSorter) MostPopularSkills(jobs []*internal.Job, count uint16) []string {
	return nil
}
