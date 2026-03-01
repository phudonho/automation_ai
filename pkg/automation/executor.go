package automation

import (
	"fmt"
)

// ActionExecutor định nghĩa interface để thực thi các thao tác tương tác chuột/phím
type ActionExecutor interface {
	Execute(action *ActionResponse) error
}

// ExecutorType định nghĩa các loại công cụ thực thi
type ExecutorType string

const (
	MockExecutor    ExecutorType = "mock"
	XDoToolExecutor ExecutorType = "xdotool"
	RobotGoExecutor ExecutorType = "robotgo"
)

// NewExecutor sử dụng Factory pattern để khởi tạo công cụ thực thi
func NewExecutor(execType ExecutorType) (ActionExecutor, error) {
	switch execType {
	case MockExecutor:
		return newMockExecutor(), nil
	case XDoToolExecutor:
		return newXDoToolExecutor(), nil
	case RobotGoExecutor:
		return newRobotGoExecutor(), nil
	default:
		return nil, fmt.Errorf("loại executor không được hỗ trợ: %s", execType)
	}
}
