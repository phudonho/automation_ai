package automation

import (
	"fmt"
	"strconv"
	"time"
)

type mockExecutor struct{}

func newMockExecutor() *mockExecutor {
	return &mockExecutor{}
}

// Execute in ra console mô phỏng hành động mà không thực sự điều khiển chuột/phím
func (m *mockExecutor) Execute(action *ActionResponse) error {
	fmt.Printf("\n[MOCK EXECUTOR] Bắt đầu thực hiện hành động: %s\n", action.NextStep)
	switch action.NextStep {
	case "CLICK":
		fmt.Printf("=> [CHUỘT] Di chuyển đến tọa độ (%d, %d) và Click trái.\n", action.Coordinates.X, action.Coordinates.Y)

	case "TYPE":
		fmt.Printf("=> [CHUỘT] Di chuyển đến tọa độ (%d, %d) và Click trái để focus ô nhập liệu.\n", action.Coordinates.X, action.Coordinates.Y)
		fmt.Printf("=> [BÀN PHÍM] Sẽ nhập chuỗi văn bản: \"%s\"\n", action.Value)

	case "SCROLL":
		fmt.Printf("=> [CHUỘT] Cuộn chuột với tham số: \"%s\"\n", action.Value)

	case "WAIT":
		seconds, err := strconv.Atoi(action.Value)
		if err != nil {
			seconds = 2 // mặc định nếu parse lỗi
		}
		fmt.Printf("=> [HỆ THỐNG] Đang chờ %d giây...\n", seconds)
		time.Sleep(time.Duration(seconds) * time.Second)

	case "OPEN":
		fmt.Printf("=> [HỆ ĐIỀU HÀNH] Sẽ mở trình duyệt truy cập vào URL: \"%s\"\n", action.Value)

	case "FINISH":
		fmt.Printf("=> [HỆ THỐNG] NHIỆM VỤ ĐÃ HOÀN THÀNH. Dừng tiến trình.\n")

	default:
		return fmt.Errorf("không nhận diện được hành động: %s", action.NextStep)
	}
	fmt.Println("[MOCK EXECUTOR] Đã thực hiện xong.")
	return nil
}
