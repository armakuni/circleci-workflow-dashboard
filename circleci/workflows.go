package circleci

type Workflow struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}

type Workflows []Workflow

func (w *Workflow) BuildError() bool {
	return w.Name == "Build Error"
}

func (w *Workflows) BuildError() bool {
	for _, workflow := range *w {
		if workflow.BuildError() {
			return true
		}
	}
	return false
}
