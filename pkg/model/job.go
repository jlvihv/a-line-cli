package model

import (
	"fmt"
	"time"
)

type Status int

const (
	STATUS_NOTRUN  Status = 0
	STATUS_RUNNING Status = 1
	STATUS_FAIL    Status = 2
	STATUS_SUCCESS Status = 3
)

type Job struct {
	Version string           `yaml:"version"`
	Name    string           `yaml:"name"`
	Stages  map[string]Stage `yaml:"stages"`
}

type JobDetail struct {
	Id int
	Job
	Status    Status
	Stages    []StageDetail
	StartTime time.Time
	Duration  time.Duration
}

func (jd *JobDetail) ToString() string {
	return fmt.Sprintf("job: %s, Status: %d, StartTime: %s , Duration: %d, stages: [%v]", jd.Name, jd.Status, jd.StartTime, jd.Duration, jd.Stages)
}

// StageSort job 排序
func (job *Job) StageSort() ([]StageDetail, error) {
	stages := make(map[string]Stage)
	for key, stage := range job.Stages {
		stages[key] = stage
	}

	sortedMap := make(map[string]any)

	stageList := make([]StageDetail, 0)
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
				stageList = append(stageList, NewStageDetail(key, stage))
			}
		}

		if len(stages) == last {
			return nil, fmt.Errorf("cannot resolve dependency, %v", stages)
		}

	}

	return stageList, nil
}

type JobLog struct {
	// 开始时间
	StartTime time.Time
	// 持续时间
	Duration time.Duration

	//日志内容
	Content string

	//最后一行 行号
	LastLine int
}

type JobStageLog struct {
	// 开始时间
	StartTime time.Time
	// 持续时间
	Duration time.Duration

	//日志内容
	Content string

	//最后一行 行号
	LastLine int
}
