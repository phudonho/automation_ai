package utils

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"os/exec"
	"runtime"

	"github.com/kbinani/screenshot"
)

// CaptureScreen chụp toàn bộ màn hình và lưu thành file sạch không có đường kẻ.
func CaptureScreen(outputPath string) error {
	var img image.Image
	var err error

	// Thử dùng gnome-screenshot trên Linux (để vượt qua giới hạn của Wayland/DBus)
	if runtime.GOOS == "linux" {
		tmpPath := "temp_gnome_screenshot.png"
		cmd := exec.Command("gnome-screenshot", "-f", tmpPath)
		if errCmd := cmd.Run(); errCmd == nil {
			defer os.Remove(tmpPath)
			f, errF := os.Open(tmpPath)
			if errF == nil {
				img, err = png.Decode(f)
				f.Close()
			}
		} else {
			fmt.Printf("Warning: gnome-screenshot failed: %v\n", errCmd)
		}
	}

	// Nếu cách trên thất bại hoặc không phải Linux, dùng thư viện kbinani/screenshot
	if img == nil {
		if n := screenshot.NumActiveDisplays(); n <= 0 {
			return fmt.Errorf("no active displays found")
		}
		bounds := screenshot.GetDisplayBounds(0)
		img, err = screenshot.CaptureRect(bounds)
		if err != nil {
			return fmt.Errorf("failed to capture screen fallback: %w", err)
		}
	}

	// Save to file (clean image)
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %w", outputPath, err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return fmt.Errorf("failed to encode png: %w", err)
	}

	return nil
}

// CaptureScreenAndDrawGrid captures the screen to a temporary clean file, draws
// a grid overlay, and writes the result to outputPath. The cleanPath will be
// removed after processing.
func CaptureScreenAndDrawGrid(cleanPath, outputPath string, gridSize int) error {
	if gridSize <= 0 {
		gridSize = 100
	}

	// Capture clean screenshot to cleanPath
	if err := CaptureScreen(cleanPath); err != nil {
		return fmt.Errorf("capture failed: %w", err)
	}
	defer os.Remove(cleanPath)

	// Open captured image
	f, err := os.Open(cleanPath)
	if err != nil {
		return fmt.Errorf("failed to open clean capture: %w", err)
	}
	defer f.Close()

	srcImg, err := png.Decode(f)
	if err != nil {
		return fmt.Errorf("failed to decode clean capture: %w", err)
	}

	bounds := srcImg.Bounds()
	rgba := image.NewRGBA(bounds)
	draw.Draw(rgba, bounds, srcImg, bounds.Min, draw.Src)

	// Grid line color: semi-transparent red
	lineCol := color.RGBA{255, 0, 0, 200}

	// Draw vertical lines
	for x := bounds.Min.X; x < bounds.Max.X; x += gridSize {
		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			rgba.Set(x, y, lineCol)
		}
	}

	// Draw horizontal lines
	for y := bounds.Min.Y; y < bounds.Max.Y; y += gridSize {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			rgba.Set(x, y, lineCol)
		}
	}

	// Save output
	out, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer out.Close()

	if err := png.Encode(out, rgba); err != nil {
		return fmt.Errorf("failed to encode output png: %w", err)
	}

	return nil
}
