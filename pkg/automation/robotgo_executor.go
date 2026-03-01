package automation

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/go-vgo/robotgo"
)

type robotGoExecutor struct{}

func newRobotGoExecutor() *robotGoExecutor {
	return &robotGoExecutor{}
}

func (r *robotGoExecutor) Execute(action *ActionResponse) error {
	switch action.NextStep {
	case "CLICK":
		fmt.Printf("[RobotGo] Đang di chuyển chuột tới (%d, %d)\n", action.Coordinates.X, action.Coordinates.Y)
		robotgo.Move(action.Coordinates.X, action.Coordinates.Y)
		time.Sleep(100 * time.Millisecond)
		fmt.Printf("[RobotGo] Đang click trái...\n")
		robotgo.Click("left")

	case "TYPE":
		// 1. Di chuyển chuột tới tọa độ ô nhập liệu và click trái để focus
		fmt.Printf("[RobotGo] Đang focus chuột tới (%d, %d)\n", action.Coordinates.X, action.Coordinates.Y)
		robotgo.Move(action.Coordinates.X, action.Coordinates.Y)
		time.Sleep(100 * time.Millisecond)
		robotgo.Click("left")

		time.Sleep(200 * time.Millisecond)

		// 2. Gõ phím
		fmt.Printf("[RobotGo] Đang nhập liệu: \"%s\"\n", action.Value)
		robotgo.TypeStr(action.Value)

		time.Sleep(100 * time.Millisecond)
		// 3. Tự động bấm Enter sau khi type
		fmt.Printf("[RobotGo] Đang gửi phím Enter...\n")
		robotgo.KeyTap("enter")

	case "SCROLL":
		amount := 5 // Tốc độ cuộn mặc định
		fmt.Printf("[RobotGo] Đang cuộn %d dòng về hướng: %s\n", amount, action.Value)
		if action.Value == "up" {
			robotgo.ScrollDir(amount, "up")
		} else {
			robotgo.ScrollDir(amount, "down")
		}

	case "WAIT":
		seconds, err := strconv.Atoi(action.Value)
		if err != nil {
			seconds = 3 // Mặc định 3s
		}
		fmt.Printf("[RobotGo] Đang chờ %d giây...\n", seconds)
		time.Sleep(time.Duration(seconds) * time.Second)

	case "OPEN":
		fmt.Printf("[Hệ Thống] Đang mở trình duyệt với URL: \"%s\"\n", action.Value)
		err := openBrowser(action.Value)
		if err != nil {
			return fmt.Errorf("lỗi mở trình duyệt: %v", err)
		}
		time.Sleep(3 * time.Second)

	case "FINISH":
		fmt.Println("Đã nhận được lệnh kết thúc từ AI.")
	default:
		return fmt.Errorf("hành động không được hỗ trợ bởi robotgo: %s", action.NextStep)
	}

	return nil
}

// Hàm hỗ trợ mở URL đa nền tảng
func openBrowser(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start", url}
	case "darwin":
		cmd = "open"
		args = []string{url}
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
		args = []string{url}
	}

	return exec.Command(cmd, args...).Run()
}
