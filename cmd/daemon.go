package cmd

import (
	"github.com/hamster-shared/a-line-cli/pkg/dispatcher"
	"github.com/hamster-shared/a-line-cli/pkg/executor"
	model2 "github.com/hamster-shared/a-line-cli/pkg/model"
	"github.com/hamster-shared/a-line-cli/pkg/pipeline"
	"github.com/hamster-shared/a-line-cli/pkg/service"
	"io"
)

func Main(reader io.Reader) {

	channel := make(chan model2.QueueMessage)

	dispatch := dispatcher.NewDispatcher(channel)

	// 本地注册
	dispatch.Register(&model2.Node{
		Name:    "localhost",
		Address: "127.0.0.1",
	})

	jobService := service.NewJobService()

	// 启动executor

	executeClient := executor.NewExecutorClient(channel, jobService)
	defer close(channel)

	go executeClient.Main()

	job, _ := pipeline.GetJobFromReader(reader)
	jobService.SaveJob(job.Name, job)

	node := dispatch.DispatchNode(job)
	dispatch.SendJob(job, node)

	select {}
}
