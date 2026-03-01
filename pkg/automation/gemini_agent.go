package automation

import (
	"encoding/json"
	"fmt"
	"helper/pkg/gemini"
	"strings"
)

type geminiAgent struct {
	client *gemini.Client
}

func newGeminiAgent(cookie string) *geminiAgent {
	return &geminiAgent{
		client: gemini.NewClient(cookie),
	}
}

// DecideNextAction sử dụng Gemini để phân tích ảnh và quyết định hành động bằng Bounding Box
func (a *geminiAgent) DecideNextAction(task string, imageState string, history []string) (*ActionResponse, error) {
	// Thay đổi cấu trúc prompt sử dụng tọa độ nhúng gốc của Gemini. Trả về toạ độ [0-1000] scale.
	promptTemplate := "Lịch sử thực thi: %s. Task thiết lập: %s. Ảnh màn hình thiết bị hiện tại: `%s`. Hãy quan sát đối tượng trên màn hình cần thao tác. CHÚ Ý QUAN TRỌNG: Trước khi quyết định dùng FINISH, hãy phân tích kỹ màn hình xem mục tiêu CỦA TASK ĐÃ THỰC SỰ ĐƯỢC HOÀN TẤT CHƯA. Hãy tìm toạ độ TÂM ĐIỂM của đối tượng sẽ thao tác (Ví dụ nút bấm, thanh tìm kiếm), quy đổi và chuẩn hoá chính xác toạ độ này theo thang đo không gian 1000x1000 pixels (X=0 ở mép trái, X=1000 ở mép phải. Y=0 đỉnh, Y=1000 góc dưới cùng). Trả về JSON key: 'analysis', 'next_step' (CLICK, TYPE, SCROLL, WAIT, OPEN, FINISH), 'coordinates' ({x, y} ở thang 1000x1000), 'value' (giá trị điền hoặc url), 'message_to_user'. Chỉ trả về JSON."

	historyStr := "trống"
	if len(history) > 0 {
		historyStr = strings.Join(history, "; ")
	}

	prompt := fmt.Sprintf(promptTemplate, historyStr, task, imageState)

	// Gửi request cho Gemini
	resultStr, err := a.client.Ask(prompt)
	if err != nil {
		return nil, fmt.Errorf("lỗi khi gọi Gemini: %w", err)
	}

	// Xử lý chuỗi trả về để trích xuất JSON (đôi khi Gemini tự thêm text mô tả)
	firstBrace := strings.Index(resultStr, "{")
	lastBrace := strings.LastIndex(resultStr, "}")
	if firstBrace != -1 && lastBrace != -1 && lastBrace > firstBrace {
		resultStr = resultStr[firstBrace : lastBrace+1]
	}

	var action ActionResponse
	err = json.Unmarshal([]byte(resultStr), &action)
	if err != nil {
		return nil, fmt.Errorf("không thể parse JSON từ AI: %w\nData: %s", err, resultStr)
	}

	return &action, nil
}
