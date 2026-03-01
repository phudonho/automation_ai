package utils

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
)

// ParseBardResponse extracts the final generated text from a Bard-like streaming response body.
func ParseBardResponse(body string) (string, error) {
	// Dùng Regex tìm "rc_" theo sau là ký tự a-z0-9 sinh bất kỳ (đã bị escape)
	// Trả về tất cả các mảng index tìm được.
	re := regexp.MustCompile(`\\"rc_[a-zA-Z0-9]+\\"`)
	matches := re.FindAllStringIndex(body, -1)
	if len(matches) == 0 {
		return "", errors.New("không tìm thấy chuỗi rc_ trong dữ liệu")
	}

	// Lấy rc_ ở vị trí xuất hiện cuối cùng trong chuỗi.
	lastMatch := matches[len(matches)-1]
	lastRCIdx := lastMatch[0]

	// Từ vị trí rc_ cuối cùng, tìm ngược lên nơi chuỗi tham số JSON bắt đầu
	startIdx := strings.LastIndex(body[:lastRCIdx], "\"[null,")
	if startIdx == -1 {
		return "", errors.New("không tìm thấy block JSON chứa rc_")
	}

	// Xác định index dấu ngoặc kép kết thúc chuỗi chứa JSON do chuỗi này được escape trước đó rồi
	endIdx := -1
	for i := startIdx + 1; i < len(body); i++ {
		if body[i] == '\\' {
			i++ // Bỏ qua escape của ký tự sau backslash
			continue
		}
		if body[i] == '"' {
			endIdx = i
			break
		}
	}

	if endIdx == -1 {
		return "", errors.New("không tìm thấy kết thúc chuỗi JSON chứa rc_")
	}

	// Trích xuất đoạn string json payload
	jsonStr := body[startIdx : endIdx+1]

	// Fix nhanh một số kí tự escape lỗi từ google (như \=)
	jsonStr = strings.ReplaceAll(jsonStr, "\\=", "=")

	var innerJSON string
	err := json.Unmarshal([]byte(jsonStr), &innerJSON)
	if err != nil {
		return "", err
	}

	var inner interface{}
	err = json.Unmarshal([]byte(innerJSON), &inner)
	if err != nil {
		return "", err
	}

	// Đệ quy lấy nội dung output
	if txt, found := findRCText(inner); found {
		return txt, nil
	}

	return "", errors.New("không tìm thấy đoạn text trong dữ liệu unmarshal")
}

// findRCText Đệ quy lấy chuỗi văn bản kế tiếp chuỗi có prefix "rc_"
func findRCText(v interface{}) (string, bool) {
	arr, ok := v.([]interface{})
	if !ok {
		return "", false
	}
	if len(arr) >= 2 {
		if s, ok := arr[0].(string); ok && len(s) >= 3 && s[:3] == "rc_" {
			if textArr, ok := arr[1].([]interface{}); ok && len(textArr) > 0 {
				if text, ok := textArr[0].(string); ok {
					return text, true
				}
			}
		}
	}
	for _, item := range arr {
		if text, found := findRCText(item); found {
			return text, true
		}
	}
	return "", false
}
