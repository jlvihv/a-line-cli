package model

import "fmt"

type Job struct {
	Version string           `yaml:"version"`
	Name    string           `yaml:"name"`
	Stages  map[string]Stage `yaml:"stages"`
}

type Status int

const (
	STATUS_NOTRUN  Status = 0
	STATUS_RUNNING Status = 1
	STATUS_FAIL    Status = 2
	STATUS_SUCCESS Status = 3
)

type JobWrapper struct {
	Id int
	Job
	Status Status
}

// StageSort job 排序
func (job *Job) StageSort() ([]StageWrapper, error) {
	stages := make(map[string]Stage)
	for key, stage := range job.Stages {
		stages[key] = stage
	}

	sortedMap := make(map[string]any)

	stageList := make([]StageWrapper, 0)
	for len(stages) > 0 {
		last := len(stages)
		for key, stage := range stages {
			allContains := true
			for _, needs := range stage.Needs {
				_, ok := sortedMap[needs]
				if !ok {
					allContains = false
				}
			}
			if allContains {
				sortedMap[key] = ""
				delete(stages, key)
				stageList = append(stageList, NewStageWrapper(key, stage))
			}
		}

		if len(stages) == last {
			return nil, fmt.Errorf("cannot resolve dependency, %v", stages)
		}

	}

	return stageList, nil

}
