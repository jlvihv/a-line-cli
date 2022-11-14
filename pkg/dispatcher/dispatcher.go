package dispatcher

import (
	"state-example/pkg/executor"
	model2 "state-example/pkg/model"
)

type IDispatcher interface {
	// DispatchNode 选择节点
	DispatchNode(job *model2.Job) *model2.Node
	// Register 节点注册
	Register(node *model2.Node)
	// UnRegister 节点注销
	UnRegister(node *model2.Node)

	// HealthcheckNode 节点心跳
	HealthcheckNode(node *model2.Node)

	// SendJob 发送任务
	SendJob(job *model2.Job, node *model2.Node)

	// CancelJob 取消任务
	CancelJob(job *model2.Job, node *model2.Node)

	// GetExecutor 根据节点获取执行器
	// TODO ... 这个方法设计的不好，分布式机构后应当用api代替
	GetExecutor(node *model2.Node) executor.IExecutor
}

type Dispatcher struct {
	Channel chan model2.QueueMessage
}

func NewDispatcher(channel chan model2.QueueMessage) *Dispatcher {
	return &Dispatcher{
		channel,
	}
}

// DispatchNode 选择节点
func (d *Dispatcher) DispatchNode(job *model2.Job) *model2.Node {

	//TODO ... 单机情况直接返回 本地
	return nil
}

// Register 节点注册
func (d *Dispatcher) Register(node *model2.Node) {
	return
}

// UnRegister 节点注销
func (d *Dispatcher) UnRegister(node *model2.Node) {
	return
}

// HealthcheckNode 节点心跳
func (d *Dispatcher) HealthcheckNode(*model2.Node) {
	// TODO  ... 检查注册的心跳信息，超过3分钟没有更新的节点，踢掉
	return
}

// SendJob 发送任务
func (d *Dispatcher) SendJob(job *model2.Job, node *model2.Node) {

	// TODO ... 单机情况下 不考虑节点，直接发送本地
	// TODO ... 集群情况下 通过注册的ip 地址进行api接口调用

	d.Channel <- model2.NewStartQueueMsg(job.Name, 1)

	return
}

// CancelJob 取消任务
func (d *Dispatcher) CancelJob(job *model2.Job, node *model2.Node) {

	d.Channel <- model2.NewStopQueueMsg(job.Name, 1)
	return
}
