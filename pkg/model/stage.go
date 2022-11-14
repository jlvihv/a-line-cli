package model

type Stage struct {
	Steps []Step   `yaml:"steps"`
	Needs []string `yaml:"needs"`
}

type StageDetail struct {
	Name   string
	Stage  Stage
	Status Status
}

func NewStageDetail(name string, stage Stage) StageDetail {
	return StageDetail{
		Name:   name,
		Stage:  stage,
		Status: STATUS_NOTRUN,
	}
}
