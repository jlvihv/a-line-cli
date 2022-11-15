package service

import "github.com/hamster-shared/a-line-cli/pkg/model"

type IJobService interface {
	//SaveJob 保存Job
	SaveJob(name string, job *model.Job)

	// GetJob 获取Job
	GetJob(name string) *model.Job

	// SaveJobDetail 保存Job详情
	SaveJobDetail(name string, job *model.JobDetail)

	// GetJobDetail 获取Job详情
	GetJobDetail(name string) *model.JobDetail
	// GetJobLog 获取job日志
	GetJobLog(name string) *model.JobLog
	// GetJobStageLog 获取job的stage 日志
	GetJobStageLog(name string) map[string]*model.JobStageLog
}

type JobService struct {
	store   map[string]*model.Job
	details map[string]*model.JobDetail
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

func (svc *JobService) SaveJobDetail(name string, job *model.JobDetail) {
	// TODO .. 保存job文件
	svc.details[name] = job
}

func (svc *JobService) GetJobDetail(name string) *model.JobDetail {
	/// TODO ... 通过文件获取job
	return svc.details[name]
}

// GetJobLog 获取job日志
func (svc *JobService) GetJobLog(name string) *model.JobLog {

	//TODO ... 实现获取日志
	return nil
}

// GetJobStageLog 获取job的stage 日志
func (svc *JobService) GetJobStageLog(name string) map[string]*model.JobStageLog {

	//TODO... 实现获取阶段日志
	return nil
}
