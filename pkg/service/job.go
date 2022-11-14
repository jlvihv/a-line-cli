package service

import "github.com/hamster-shared/a-line-cli/pkg/model"

type IJobService interface {
	SaveJob(name string, job *model.Job)

	GetJob(name string) *model.Job

	SaveJobDetail(name string, job *model.JobDetail)

	GetJobDetail(name string) *model.JobDetail
}

type JobService struct {
	store map[string]*model.Job
}

func NewJobService() *JobService {
	return &JobService{
		store: make(map[string]*model.Job),
	}
}

func (svc *JobService) SaveJob(name string, job *model.Job) {
	// TODO .. 保存job文件
	svc.store[name] = job
}

func (svc *JobService) GetJob(name string) *model.Job {
	/// TODO ... 通过文件获取job
	return svc.store[name]
}
