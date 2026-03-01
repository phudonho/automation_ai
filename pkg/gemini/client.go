package gemini

import (
	"encoding/json"
	"fmt"
	"helper/utils"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Client struct {
	HTTPClient *http.Client
	BaseURL    string
	Cookie     string
}

// CookieStruct dùng để map với object cookie định dạng JSON (ví dụ từ EditThisCookie)
type CookieStruct struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type CookieWrapper struct {
	Cookies []CookieStruct `json:"cookies"`
}

// ParseCookieJSON tự động chuyển chuỗi cookie JSON sang dạng chuỗi key=value standard
func ParseCookieJSON(cookieStr string) string {
	cookieStr = strings.TrimSpace(cookieStr)
	// Trả về nguyên bản nếu không giống format JSON array hoặc object
	if !strings.HasPrefix(cookieStr, "[") && !strings.HasPrefix(cookieStr, "{") {
		return cookieStr
	}

	var cookies []CookieStruct
	// Trường hợp format: {"url": "...", "cookies": [...]}
	if strings.HasPrefix(cookieStr, "{") {
		var wrapper CookieWrapper
		if err := json.Unmarshal([]byte(cookieStr), &wrapper); err == nil {
			cookies = wrapper.Cookies
		}
	} else {
		// Trường hợp format: [{"name": "...", "value": "..."}, ...]
		_ = json.Unmarshal([]byte(cookieStr), &cookies)
	}

	if len(cookies) == 0 {
		return cookieStr
	}

	var parts []string
	for _, c := range cookies {
		parts = append(parts, fmt.Sprintf("%s=%s", c.Name, c.Value))
	}
	return strings.Join(parts, "; ")
}

func NewClient(cookie string) *Client {
	return &Client{
		HTTPClient: &http.Client{},
		BaseURL:    "https://gemini.google.com/_/BardChatUi/data/assistant.lamda.BardFrontendService/StreamGenerate?bl=boq_assistant-bard-web-server_20260218.05_p0&f.sid=-24838734129403184&hl=vi&_reqid=5743406&rt=c",
		Cookie:     ParseCookieJSON(cookie),
	}
}

func (c *Client) Ask(prompt string) (string, error) {
	// 1. Prepare payload by encoding prompt
	// JSON encode to escape quotes/newlines, then strip surrounding quotes
	b, err := json.Marshal(prompt)
	if err != nil {
		return "", fmt.Errorf("lỗi encode prompt: %v", err)
	}
	s := string(b)
	s = s[1 : len(s)-1] // Bỏ `"` ở đầu và cuối

	// Escape for URL
	encodedPrompt := url.QueryEscape(s)

	// URL-encoded payload có sẵn `heo`
	basePayload := `f.req=%5Bnull%2C%22%5B%5B%5C%22heo%5C%22%2C0%2Cnull%2Cnull%2Cnull%2Cnull%2C0%5D%2C%5B%5C%22vi%5C%22%5D%2C%5B%5C%22c_0b94bb290d45c230%5C%22%2C%5C%22r_de7dae4c5531f0e6%5C%22%2C%5C%22rc_f73de017c4246f0c%5C%22%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2C%5C%22AwAAAAAAAAAQANM7mBjXKZSZPNM7yBk%5C%22%5D%2C%5C%22!jo2ljdXNAAaWZrTItmZCkubkhdxfR907ADQBEArZ1KvfV4--uWNkx-Ep46jvEtiqASvhdl0F-iWoMp3OshUdD3oJPoej_ga9uzGvijkVAgAAAGxSAAAABWgBB34AP4Oob2cZG69HKaRcPUZHjPUcaOEVGK3wipEtl-i8Xdb_SSTxVGma5LEiMTKABJc_En8iTX9Glh1C6DaYummRcZkDfRwSGIljC7pESXTVVIVnboUqR9bBC_ReeYPUg700X3Tn6SxYpEUxFovyAYW9cf72tafwPAUw6Qjwhk21eYujGHqWuV810175HLukj6k482aWAHfGNboeNfw39BraIeQ9AYS_BB5byG9ITkk4wKnSch4KeQ3u8X8hYBxQWj5aMrbrPIy2i1lMjotgLe3KuNXQDqoLE3yLdWtUJiwXni40fkOSQvddcOcYkjxLDB24vZb8nv_fJeEief4GR6SZp-3OxxkkY4FsezeclcyjySZgVgg2IaaadJKeGCpcdyUi21_Tiavb_ND9ok0ziUchxxtB_Hfu_weqOUhyJKT_o2bQQGFtncHovjhx8KqckSThzUnhColYaQwc_bkBztgdjEyR7ZvHSoFYlP29A0DITdLOIYnxAGePrpMN-a5BS4AAuIqYRpNSAP4ZSqJG5S54rXqXLp62j-r2NbT7i-UwONCfchqjYJp7WiOF2P2YBz2idpUPuqLrfCcL9xrXV8W6SS6Da4Dxl1IsbAM1d-TJg-4cHT6IV4g-aLe6fgbcdTRoNYkvxsCQAMdTxMTrnJnAtCSVH7n1G_06tXZ7B8bYXBtli0nD4-ftIxVn4sw2cWrHbJFkuGyudAI0i5M_HpDCe6gWEh_QIon-UGAG0leYLptad4VIAF3-AExleuSOTrd4prLrnmmwHMi5BUzh4oSZcC1VdMPb_Lpj9fmPShu41ki32xV7rGEwv3r0onws6LJNhasTVdTJ2c8gfBNaXjxa6iHn5qODPqxZPhIduY9coFvD6jyEF0mrJxZb8KUv_BjVyCBnFOsChuADSLhdQZ7mDSYGDHBCuxfG3BHYVe3aXnaefnB6fkLndJCXag8jUGAw2Nc_kxDVfSzseR6wJv1pY2YVDib-FqfLRxwxHMa1xM9Qv58l9gc0uQB05nbbgYM3nKHE1TSiBDGj6gYu_YxtRLbrKGIZsEkVWmzeQmKqGr8Jw6ENupcdg3b6TCDGZgKkjNAA7dV6AqOeyeW2B0nqtEcFW2XO7vtThljdmId03HVFOkt4R1EiTf5IgEdqhX5JBARY5aT9u6wAygbEQRmBPPWgdEt20LpsdqZxuUcn4qJqDUexezNTehS3EEWUnemBZamzRV3n9AfEwTx5-6lPnjBpBddYbD6RIdsCExL_BtH5y0WLWc4Gm4Uy3Eb7cvWM%5C%22%2C%5C%22fde32c2414b977290e21fe5587579977%5C%22%2Cnull%2C%5B0%5D%2C1%2Cnull%2Cnull%2C1%2C0%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2C%5B%5B2%5D%5D%2C0%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2C1%2Cnull%2Cnull%2C%5B4%5D%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2C%5B1%5D%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2C0%2Cnull%2Cnull%2Cnull%2Cnull%2Cnull%2C%5C%220AF73504-5D84-4EEA-96F4-6E6F037C978D%5C%22%2Cnull%2C%5B%5D%2Cnull%2Cnull%2Cnull%2Cnull%2C%5B1771736740%2C557000000%5D%2Cnull%2C1%5D%22%5D&at=AEHmXlHsCLgHZ05RpskY0Y1f5B5s%3A1771736602761&`

	// Thay thế "heo" bằng đoạn prompt đã encode
	payloadStr := strings.Replace(basePayload, "%5C%22heo%5C%22", "%5C%22"+encodedPrompt+"%5C%22", 1)

	// 2. Tạo Request
	req, err := http.NewRequest("POST", c.BaseURL, strings.NewReader(payloadStr))
	if err != nil {
		return "", fmt.Errorf("lỗi tạo request: %v", err)
	}

	// 3. Set Headers
	req.Header.Add("accept", "*/*")
	req.Header.Add("accept-language", "en,en-US;q=0.9,vi;q=0.8")
	req.Header.Add("content-type", "application/x-www-form-urlencoded;charset=UTF-8")
	req.Header.Add("origin", "https://gemini.google.com")
	req.Header.Add("referer", "https://gemini.google.com/")
	req.Header.Add("user-agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36")
	if c.Cookie != "" {
		req.Header.Add("cookie", c.Cookie)
	}

	// 4. Send Request
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("lỗi thực thi request: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("http status code không hợp lệ: %d", res.StatusCode)
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("lỗi đọc response: %v", err)
	}

	// 5. Parse dữ liệu trả về thông qua Utils
	return utils.ParseBardResponse(string(bodyBytes))
}
