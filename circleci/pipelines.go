package circleci

type VCS struct {
	Branch string `json:"branch"`
}

type Pipeline struct {
	ID  string `json:"id"`
	VCS VCS    `json:"vcs"`
}

type Pipelines []Pipeline

func (p Pipelines) FilteredPerBranch() map[string]Pipelines {
	filteredPipelines := make(map[string]Pipelines)
	for _, pipeline := range p {
		filteredPipelinesForBranch := filteredPipelines[pipeline.VCS.Branch]
		filteredPipelinesForBranch = append(filteredPipelinesForBranch, pipeline)
		filteredPipelines[pipeline.VCS.Branch] = filteredPipelinesForBranch
	}
	return filteredPipelines
}

func (p Pipelines) LatestPerBranch() map[string]Pipeline {
	latestPipelines := make(map[string]Pipeline)
	for _, pipeline := range p {
		if _, ok := latestPipelines[pipeline.VCS.Branch]; ok {
			continue
		}
		latestPipelines[pipeline.VCS.Branch] = pipeline
	}
	return latestPipelines
}
