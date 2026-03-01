package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// SxcuResponse là cấu trúc trả về từ sxcu.net
type SxcuResponse struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Thumb  string `json:"thumb"`
	DelURL string `json:"del_url"`
}

// UploadToSxcu tải file lên sxcu.net và trả về direct URL
func UploadToSxcu(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("không thể mở file: %v", err)
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// API của sxcu.net yêu cầu param 'endpoint=sxcu' (như user đã curl)
	_ = writer.WriteField("endpoint", "sxcu")

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", fmt.Errorf("lỗi tạo form file: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return "", fmt.Errorf("lỗi copy dữ liệu file: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return "", fmt.Errorf("lỗi đóng writer: %v", err)
	}

	req, err := http.NewRequest("POST", "https://sxcu.net/api/files/create", body)
	if err != nil {
		return "", fmt.Errorf("lỗi tạo http request: %v", err)
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("User-Agent", "curl/7.81.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("lỗi upload file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("server sxcu.net trả về code: %d", resp.StatusCode)
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("lỗi đọc response: %v", err)
	}

	var sxcuResp SxcuResponse
	err = json.Unmarshal(respBody, &sxcuResp)
	if err != nil {
		return "", fmt.Errorf("lỗi parse JSON từ sxcu.net: %v", err)
	}

	return sxcuResp.URL, nil
}
