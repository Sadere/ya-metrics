package agent

import "errors"

var (
	ErrAgentSendFailed = errors.New("agent couldn't transfer data to server")
)
