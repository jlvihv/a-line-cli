package service

import (
	"github.com/hamster-shared/a-line-cli/pkg/consts"
	"github.com/hamster-shared/a-line-cli/pkg/model"
	"github.com/hamster-shared/a-line-cli/pkg/utils"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

type IJobService interface {
	//SaveJob 保存Job
	SaveJob(name string, job *model.Job)

	// GetJob 获取Job
	GetJob(name string) *model.Job

	//UpdateJob update job
	UpdateJob(name string, job *model.Job) error

	//DeleteJob delete job
	DeleteJob(name string) error

	// SaveJobDetail 保存Job详情
	SaveJobDetail(name string, job *model.JobDetail)

	// GetJobDetail 获取Job详情
	GetJobDetail(name string, id int) *model.JobDetail

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
		store:   make(map[string]*model.Job),
		details: make(map[string]*model.JobDetail),
	}
}

func (svc *JobService) SaveJob(name string, job *model.Job) {
	// serializes yaml struct
	data, err := yaml.Marshal(job)
	if err != nil {
		log.Println("serializes yaml failed", err)
	}
	//determine whether the folder exists, and create it if it does not exist
	utils.CreateJobDir()
	//file directory path
	dir := filepath.Join(utils.DefaultConfigDir(), consts.JOB_DIR_NAME+"/"+name)
	src := filepath.Join(utils.DefaultConfigDir(), consts.JOB_DIR_NAME+"/"+name+"/"+name+".yml")
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Println("create jobs dir failed", err.Error())
		}
	}
	//write data to yaml file
	err = os.WriteFile(src, data, 0777)
	if err != nil {
		log.Println("write data to yaml file failed", err)
	}
}

func (svc *JobService) GetJob(name string) *model.Job {
	var jobData model.Job
	//job file path
	src := filepath.Join(utils.DefaultConfigDir(), consts.JOB_DIR_NAME+"/"+name+"/"+name+".yml")
	//judge whether the job file exists
	_, err := os.Stat(src)
	//not exist
	if os.IsNotExist(err) {
		log.Println("get job failed,job file not exist", err.Error())
		return &jobData
	}
	//exist
	fileContent, err := os.ReadFile(src)
	if err != nil {
		log.Println("get job read file failed", err.Error())
		return &jobData
	}
	//deserialization job yml file
	err = yaml.Unmarshal(fileContent, &jobData)
	if err != nil {
		log.Println("get job,deserialization job file failed", err.Error())
		return &jobData
	}
	return &jobData
}

func (svc *JobService) UpdateJob(name string, job *model.Job) error {
	// job file path
	src := filepath.Join(utils.DefaultConfigDir(), consts.JOB_DIR_NAME+"/"+name+"/"+name+".yml")
	//judge whether the job detail file exists
	_, err := os.Stat(src)
	//not exist
	if os.IsNotExist(err) {
		log.Println("update job failed,job file not exist", err.Error())
		return err
	}
	// serializes yaml struct
	data, err := yaml.Marshal(job)
	if err != nil {
		log.Println("serializes yaml failed", err)
		return err
	}
	//write data to yaml file
	err = os.WriteFile(src, data, 0777)
	if err != nil {
		log.Println("write data to yaml file failed", err)
		return err
	}
	return nil
}

func (svc *JobService) DeleteJob(name string) error {
	// job file path
	src := filepath.Join(utils.DefaultConfigDir(), consts.JOB_DIR_NAME+"/"+name+"/"+name+".yml")
	//judge whether the job detail file exists
	_, err := os.Stat(src)
	//not exist
	if os.IsNotExist(err) {
		log.Println("delete job failed,job file not exist", err.Error())
		return err
	}
	err = os.Remove(src)
	if err != nil {
		log.Println("delete job failed", err.Error())
		return err
	}
	return nil
}

func (svc *JobService) SaveJobDetail(name string, job *model.JobDetail) {
	// serializes yaml struct
	data, err := yaml.Marshal(job)
	if err != nil {
		log.Println("serializes yaml failed", err)
	}
	//determine whether the folder exists, and create it if it does not exist
	utils.CreateJobDetailDir()
	//file directory path
	dir := filepath.Join(utils.DefaultConfigDir(), consts.JOB_DIR_NAME+"/"+name+"/"+consts.JOB_DETAIL_DIR_NAME)
	_, err = os.Stat(dir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			log.Println("create job details failed", err.Error())
		}
	}
	src := filepath.Join(utils.DefaultConfigDir(), consts.JOB_DIR_NAME+"/"+name+"/"+consts.JOB_DETAIL_DIR_NAME+"/"+strconv.Itoa(job.Id)+".yml")
	//write data to yaml file
	err = os.WriteFile(src, data, 0777)
	if err != nil {
		log.Println("write data to yaml file failed", err)
	}
}

func (svc *JobService) GetJobDetail(name string, id int) *model.JobDetail {
	var jobDetailData model.JobDetail
	//job file path
	src := filepath.Join(utils.DefaultConfigDir(), consts.JOB_DIR_NAME+"/"+name+"/"+consts.JOB_DETAIL_DIR_NAME+"/"+strconv.Itoa(id)+".yml")
	//judge whether the job detail file exists
	_, err := os.Stat(src)
	//not exist
	if os.IsNotExist(err) {
		log.Println("get job detail failed,job detail file not exist", err.Error())
		return &jobDetailData
	}
	//exist
	fileContent, err := os.ReadFile(src)
	if err != nil {
		log.Println("get job read detail file failed", err.Error())
		return &jobDetailData
	}
	//deserialization job detail yml file
	err = yaml.Unmarshal(fileContent, &jobDetailData)
	if err != nil {
		log.Println("get job,deserialization job detail file failed", err.Error())
		return &jobDetailData
	}
	return &jobDetailData
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
