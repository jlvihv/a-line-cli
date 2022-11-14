package cmd

import (
	"state-example/pkg/dispatcher"
	"state-example/pkg/executor"
	model2 "state-example/pkg/model"
)

func Main() {

	channel := make(chan model2.QueueMessage)

	dispatch := dispatcher.NewDispatcher(channel)

	// 本地注册
	dispatch.Register(&model2.Node{
		Name:    "localhost",
		Address: "127.0.0.1",
	})

	// 启动executor

	executeClient := executor.NewExecutorClient(channel)
	defer close(channel)

	go executeClient.Main()
}
