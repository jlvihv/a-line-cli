package cmd

import (
	"fmt"
	"github.com/hamster-shared/a-line-cli/pkg/dispatcher"
	"github.com/hamster-shared/a-line-cli/pkg/executor"
	"github.com/hamster-shared/a-line-cli/pkg/http"
	"github.com/hamster-shared/a-line-cli/pkg/model"
	"github.com/hamster-shared/a-line-cli/pkg/pipeline"
	"github.com/hamster-shared/a-line-cli/pkg/service"
	"io"
	"time"
)

func Main(reader io.Reader) {

	channel := make(chan model.QueueMessage)

	dispatch := dispatcher.NewDispatcher(channel)

	// 本地注册
	dispatch.Register(&model.Node{
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

	ticker := time.NewTicker(1000 * time.Millisecond)
	done := make(chan bool)

	detail := jobService.GetJobDetail(job.Name)

	go func() {
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
				detail = jobService.GetJobDetail(job.Name)
				if detail != nil {
					fmt.Println(detail.ToString())

					if detail.Status == model.STATUS_SUCCESS || detail.Status == model.STATUS_FAIL {
						ticker.Stop()
					}
				}
			}
		}
	}()

	http.NewHttpService(jobService).StartHttpServer()
}
