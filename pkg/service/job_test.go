package service

import (
	"github.com/hamster-shared/a-line-cli/pkg/model"
	"github.com/stretchr/testify/assert"
	ass "gotest.tools/v3/assert"
	"log"
	"testing"
)

func Test_SaveJob(t *testing.T) {
	step1 := model.Step{
		Name: "sun",
		Uses: "",
		With: map[string]string{
			"pipelie": "string",
			"data":    "data",
		},
		RunsOn: "open",
		Run:    "stage",
	}
	var steps []model.Step
	var strs []string
	strs = append(strs, "strings")
	steps = append(steps, step1)
	job := model.Job{
		Version: "1",
		Name:    "mysql",
		Stages: map[string]model.Stage{
			"node": {
				Steps: steps,
				Needs: strs,
			},
		},
	}
	jobService := NewJobService()
	jobService.SaveJob("guo", &job)
}

func Test_SaveJobDetail(t *testing.T) {
	step1 := model.Step{
		Name: "sun",
		Uses: "",
		With: map[string]string{
			"pipelie": "string",
			"data":    "data",
		},
		RunsOn: "open",
		Run:    "stage",
	}
	var steps []model.Step
	var strs []string
	strs = append(strs, "strings")
	steps = append(steps, step1)
	stageDetail := model.StageDetail{
		Name: "string",
		Stage: model.Stage{
			Steps: steps,
			Needs: strs,
		},
		Status: model.STATUS_FAIL,
	}
	var stageDetails []model.StageDetail
	stageDetails = append(stageDetails, stageDetail)
	jobDetail := model.JobDetail{
		Id: 2,
		Job: model.Job{
			Version: "2",
			Name:    "mysql",
			Stages: map[string]model.Stage{
				"node": {
					Steps: steps,
					Needs: strs,
				},
			},
		},
		Status: model.STATUS_NOTRUN,
		Stages: stageDetails,
	}
	jobService := NewJobService()
	jobService.SaveJobDetail("guo", &jobDetail)
}

func Test_GetJob(t *testing.T) {
	jobService := NewJobService()
	data := jobService.GetJob("guo")
	log.Println(data)
	assert.NotNil(t, data)
}

func Test_UpdateJob(t *testing.T) {
	jobService := NewJobService()
	step1 := model.Step{
		Name: "jian",
		Uses: "",
		With: map[string]string{
			"pipelie": "string",
			"data":    "data",
		},
		RunsOn: "open",
		Run:    "stage",
	}
	var steps []model.Step
	var strs []string
	strs = append(strs, "strings")
	steps = append(steps, step1)
	job := model.Job{
		Version: "1",
		Name:    "mysql",
		Stages: map[string]model.Stage{
			"node": {
				Steps: steps,
				Needs: strs,
			},
		},
	}
	err := jobService.UpdateJob("guo", "jian", &job)
	ass.NilError(t, err)
}

func Test_GetJobDetail(t *testing.T) {
	jobService := NewJobService()
	data := jobService.GetJobDetail("sun", 3)
	assert.NotNil(t, data)
}

func Test_DeleteJob(t *testing.T) {
	jobService := NewJobService()
	err := jobService.DeleteJob("sun")
	ass.NilError(t, err)
}

func Test_DeleteJobDetail(t *testing.T) {
	jobService := NewJobService()
	err := jobService.DeleteJobDetail("cdqadqa92d3if4r9n8j0", 1)
	ass.NilError(t, err)
}

func Test_JobList(t *testing.T) {
	jobService := NewJobService()
	data := jobService.JobList("cdqadqa92d3if4r9n8j0", 1, 10)
	assert.NotNil(t, data)
}

func Test_JobDetailList(t *testing.T) {
	jobService := NewJobService()
	data := jobService.JobDetailList("sun", 2, 10)
	log.Println(data)
	assert.NotNil(t, data)
}

func Test_ExecuteJob(t *testing.T) {
	jobService := NewJobService()
	jobService.ExecuteJob("sun")
}
