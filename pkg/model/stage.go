package model

import (
	"fmt"
	"time"
)

type Stage struct {
	Steps []Step   `yaml:"steps"`
	Needs []string `yaml:"needs"`
}

type StageDetail struct {
	Name      string
	Stage     Stage
	Status    Status
	StartTime time.Time
	Duration  time.Duration
}

func NewStageDetail(name string, stage Stage) StageDetail {
	return StageDetail{
		Name:   name,
		Stage:  stage,
		Status: STATUS_NOTRUN,
	}
}

func (s *StageDetail) ToString() string {
	return fmt.Sprintf("StageName: %s, status: %d, StartTime: %s, Duration: %d ", s.Name, s.Status, s.StartTime, s.Duration)
}
