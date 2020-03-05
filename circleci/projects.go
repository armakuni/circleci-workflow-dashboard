package circleci

import "fmt"

type Project struct {
	VCSType  string                 `json:"vcs_type"`
	Username string                 `json:"username"`
	Reponame string                 `json:"reponame"`
	Branches map[string]interface{} `json:"branches"`
}

type Projects []Project

func (p *Project) Slug() string {
	return fmt.Sprintf("%s/%s/%s", p.VCSType, p.Username, p.Reponame)
}

func (p *Project) Name() string {
	return fmt.Sprintf("%s/%s", p.Username, p.Reponame)
}

func (p Projects) Filter(filter *Filter) Projects {
	if len(*filter) == 0 {
		return p
	}
	var keepProjects Projects
	for _, project := range p {
		for projectFitler, _ := range *filter {
			if project.Name() == projectFitler {
				keepProjects = append(keepProjects, project)
			}
		}
	}
	return keepProjects
}
