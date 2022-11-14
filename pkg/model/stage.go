package model

type Stage struct {
	Steps []Step   `yaml:"steps"`
	Needs []string `yaml:"needs"`
}

type StageWrapper struct {
	Name  string
	Stage Stage
}

func NewStageWrapper(name string, stage Stage) StageWrapper {
	return StageWrapper{
		Name:  name,
		Stage: stage,
	}
}
