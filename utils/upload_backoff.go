package utils

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// nếu sxcu của bạn trả lỗi dưới dạng string "server sxcu.net trả về code: 429"
// ta detect bằng chuỗi. Tốt hơn là bạn sửa UploadToSxcu trả về status code,
// nhưng wrapper này vẫn chạy được ngay.
func is429(err error) (bool, int) {
	if err == nil {
		return false, 0
	}
	msg := err.Error()
	if strings.Contains(msg, "code: 429") || strings.Contains(msg, "429") {
		return true, 429
	}
	return false, 0
}

// UploadToSxcuWithBackoff sẽ:
// - enforce minInterval giữa các upload (rate limit chủ động)
// - retry với exponential backoff khi gặp 429
// - respect Retry-After nếu bạn có thể lấy được header (nếu không, dùng backoff)
type UploadLimiter struct {
	MinInterval time.Duration
	last        time.Time
}

func NewUploadLimiter(minInterval time.Duration) *UploadLimiter {
	return &UploadLimiter{MinInterval: minInterval}
}

func (l *UploadLimiter) WaitIfNeeded() {
	if l.MinInterval <= 0 {
		return
	}
	if !l.last.IsZero() {
		elapsed := time.Since(l.last)
		if elapsed < l.MinInterval {
			time.Sleep(l.MinInterval - elapsed)
		}
	}
	l.last = time.Now()
}

// Nếu sau này bạn sửa UploadToSxcu để trả *http.Response hoặc status code,
// hãy dùng hàm này để parse Retry-After.
func retryAfterSeconds(resp *http.Response) time.Duration {
	if resp == nil {
		return 0
	}
	ra := resp.Header.Get("Retry-After")
	if ra == "" {
		return 0
	}
	// Retry-After có thể là seconds hoặc HTTP date. Ở đây xử seconds trước.
	if secs, err := strconv.Atoi(strings.TrimSpace(ra)); err == nil && secs > 0 {
		return time.Duration(secs) * time.Second
	}
	return 0
}

func jitter(d time.Duration) time.Duration {
	// jitter 0–300ms để tránh “đập đồng loạt”
	return d + time.Duration(rand.Intn(300))*time.Millisecond
}

// Wrapper gọi UploadToSxcu (của bạn) và retry nếu 429.
func UploadToSxcuWithBackoff(
	limiter *UploadLimiter,
	screenshotPath string,
	maxRetries int,
	baseDelay time.Duration,
) (string, error) {

	if limiter != nil {
		limiter.WaitIfNeeded()
	}

	var lastErr error
	delay := baseDelay
	if delay <= 0 {
		delay = 2 * time.Second
	}

	for attempt := 0; attempt <= maxRetries; attempt++ {
		url, err := UploadToSxcu(screenshotPath)
		if err == nil {
			return url, nil
		}

		lastErr = err

		isRateLimited, _ := is429(err)
		if !isRateLimited {
			// lỗi khác 429 thì không retry kiểu backoff (tuỳ bạn)
			return "", err
		}

		if attempt == maxRetries {
			break
		}

		// Exponential backoff: 2s, 4s, 8s, 16s...
		sleepFor := jitter(delay)
		time.Sleep(sleepFor)
		delay *= 2
		if delay > 60*time.Second {
			delay = 60 * time.Second
		}
	}

	return "", errors.New(fmt.Sprintf("upload bị rate limit liên tục sau %d retries: %v", maxRetries, lastErr))
}
