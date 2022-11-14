package executor

import (
	"state-example/pkg/logger"
	"state-example/pkg/model"
	"state-example/pkg/pipeline"
)

func NewExecutorClient(channel chan model.QueueMessage) *ExecutorClient {
	return &ExecutorClient{
		executor: &Executor{
			cancelMap: make(map[string]func()),
		},
		channel: channel,
	}
}

type ExecutorClient struct {
	executor IExecutor
	channel  chan model.QueueMessage
}

func (c *ExecutorClient) Main() {
	//1. TODO... 注册节点

	//2. TODO...  gorouting 发送定时心跳

	for { //

		//3. 监听队列
		queueMessage, ok := <-c.channel
		if !ok {
			return
		}

		//4.TODO...，获取job信息

		// TODO ... 计算jobId
		jobName := queueMessage.JobName
		jobId := queueMessage.JobId

		pipelineReader, err := c.executor.FetchJob(jobName)

		if err != nil {
			logger.Error(err)
			continue
		}

		//5. 解析pipeline
		job, err := pipeline.GetJobFromReader(pipelineReader)

		//6. 异步执行pipeline
		go c.executor.Execute(jobId, job)

	}
}
