package action

// RemoteAction 执行远程命令
type RemoteAction struct {
	name string
}

func (a *RemoteAction) Pre() error {

	return nil
}

func (a *RemoteAction) Hook() error {

	return nil
}

func (a *RemoteAction) Post() error {

	return nil
}
