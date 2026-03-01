package automation

import (
	"fmt"
	"os/exec"
	"strconv"
	"time"
)

type xdotoolExecutor struct{}

func newXDoToolExecutor() *xdotoolExecutor {
	return &xdotoolExecutor{}
}

// Execute gọi lệnh xdotool thật sự trên hệ điều hành Linux sử dụng server X11
// Chỉ hoạt động trên môi trường XOrg (không chạy trực tiếp trên Wayland)
func (x *xdotoolExecutor) Execute(action *ActionResponse) error {
	switch action.NextStep {
	case "CLICK":
		// Lệnh: xdotool mousemove X Y
		sx := strconv.Itoa(action.Coordinates.X)
		sy := strconv.Itoa(action.Coordinates.Y)

		fmt.Printf("[XDoTool] Đang thực thi: xdotool mousemove %s %s\n", sx, sy)
		cmdMove := exec.Command("xdotool", "mousemove", sx, sy)
		if err := cmdMove.Run(); err != nil {
			return fmt.Errorf("lỗi di chuyển chuột: %v", err)
		}

		time.Sleep(300 * time.Millisecond)

		fmt.Println("[XDoTool] Đang thực thi: xdotool click 1")
		cmdClick := exec.Command("xdotool", "click", "1")
		if err := cmdClick.Run(); err != nil {
			return fmt.Errorf("lỗi click chuột: %v", err)
		}
	case "TYPE":
		sx := strconv.Itoa(action.Coordinates.X)
		sy := strconv.Itoa(action.Coordinates.Y)

		// 1. Di chuyển chuột tới tọa độ ô nhập liệu
		fmt.Printf("[XDoTool] Đang thực thi focus: xdotool mousemove %s %s\n", sx, sy)
		cmdMove := exec.Command("xdotool", "mousemove", sx, sy)
		if err := cmdMove.Run(); err != nil {
			return fmt.Errorf("lỗi focus trường nhập liệu (mousemove): %v", err)
		}

		time.Sleep(200 * time.Millisecond)

		fmt.Println("[XDoTool] Đang thực thi focus: xdotool click 1")
		if err := exec.Command("xdotool", "click", "1").Run(); err != nil {
			// Không báo lỗi quá gắt nếu click focus bị lỗi trên Wayland
			fmt.Printf("Cảnh báo focus click: %v\n", err)
		}

		time.Sleep(300 * time.Millisecond)

		// 2. Gõ phím
		fmt.Printf("[XDoTool] Đang thực thi: xdotool type --delay 50 \"%s\"\n", action.Value)
		cmdType := exec.Command("xdotool", "type", "--delay", "50", action.Value)
		if err := cmdType.Run(); err != nil {
			return fmt.Errorf("lỗi khi gõ bàn phím: %v", err)
		}

		// 3. Tự động bấm Enter sau khi type
		time.Sleep(100 * time.Millisecond)
		fmt.Println("[XDoTool] Đang bấm Enter: xdotool key Return")
		cmdEnter := exec.Command("xdotool", "key", "Return")
		_ = cmdEnter.Run()
	case "SCROLL":
		// xdotool click 4 (scroll up), 5 (scroll down)
		btn := "5" // mặc định cuộn xuống
		if action.Value == "up" {
			btn = "4"
		}
		cmd := exec.Command("xdotool", "click", btn)
		_ = cmd.Run()
	case "WAIT":
		seconds, err := strconv.Atoi(action.Value)
		if err != nil {
			seconds = 3 // Mặc định 3s
		}
		time.Sleep(time.Duration(seconds) * time.Second)
	case "OPEN":
		cmd := exec.Command("xdg-open", action.Value)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("lỗi khi mở trình duyệt: %v", err)
		}
		// Đợi một chút để trình duyệt kịp mở lên
		time.Sleep(3 * time.Second)
	case "FINISH":
		fmt.Println("Đã nhận được lệnh kết thúc từ AI.")
	default:
		return fmt.Errorf("hành động không được hỗ trợ bởi xdotool: %s", action.NextStep)
	}

	return nil
}
