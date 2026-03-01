package utils

import (
	"os"
	"testing"
)

func TestCaptureScreenAndDrawGrid(t *testing.T) {
	// Khai báo đường dẫn
	outputPath := "test_screenshot_grid.png"

	// Dọn dẹp file cũ nếu có
	os.Remove(outputPath)
	// Bỏ defer xoá file để bạn có thể xem kết quả sau khi test chạy xong
	// defer os.Remove(outputPath)

	// Chạy hàm chụp màn hình và vẽ lưới
	err := CaptureScreenAndDrawGrid("clean_"+outputPath, outputPath, 100)
	if err != nil {
		t.Skipf("Skipping screen capture test (no display or permission): %v", err)
	}

	// Xác nhận file đã được tạo thành công
	info, err := os.Stat(outputPath)
	if os.IsNotExist(err) {
		t.Fatalf("Expected output file %s to be created, but it was not", outputPath)
	}
	if err != nil {
		t.Fatalf("Failed to stat output file: %v", err)
	}

	// Đảm bảo file không trống
	if info.Size() == 0 {
		t.Fatalf("Output file %s is empty", outputPath)
	}

	t.Logf("Successfully captured screen and drew grid to %s (size: %d bytes)", outputPath, info.Size())
	t.Logf("Bạn có thể mở tệp %s để xem thử!", outputPath)
}
