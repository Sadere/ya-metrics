package agent

import "errors"

// Ошибки агента
var (
	ErrAgentSendFailed = errors.New("agent couldn't transfer data to server")
)
