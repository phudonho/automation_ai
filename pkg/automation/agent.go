package automation

import "fmt"

// Coordinates đại diện cho tọa độ trên màn hình
type Coordinates struct {
	X int `json:"x"`
	Y int `json:"y"`
}

// ActionResponse là cấu trúc Output JSON mà AI sẽ trả về
type ActionResponse struct {
	Analysis      string      `json:"analysis"`
	NextStep      string      `json:"next_step"`
	Coordinates   Coordinates `json:"coordinates"`
	Value         string      `json:"value"`
	MessageToUser string      `json:"message_to_user"`
}

// AutomationAgent là interface định nghĩa chung cho mọi loại AI trợ lý
type AutomationAgent interface {
	DecideNextAction(task string, imageState string, history []string) (*ActionResponse, error)
}

// AgentType định nghĩa các loại AI đang hỗ trợ
type AgentType string

const (
	Gemini AgentType = "gemini"
	// Thêm ChatGPT, Claude trong tương lai nếu cần
)

// Factory method khởi tạo Agent
func NewAgent(agentType AgentType, config interface{}) (AutomationAgent, error) {
	switch agentType {
	case Gemini:
		// Kiểm tra kiểu cấu hình
		cookie, ok := config.(string)
		if !ok {
			return nil, fmt.Errorf("gemini agent cần cấu hình kiểu string (cookie)")
		}
		// cookie = `{"url":"https://gemini.google.com","cookies":[{"domain":".gemini.google.com","expirationDate":1806296608.844867,"hostOnly":false,"httpOnly":false,"name":"_ga","path":"/","sameSite":"unspecified","secure":false,"session":false,"storeId":"0","value":"GA1.1.2053931787.1757043630"},{"domain":".google.com","expirationDate":1776679898.159614,"hostOnly":false,"httpOnly":false,"name":"SEARCH_SAMESITE","path":"/","sameSite":"strict","secure":false,"session":false,"storeId":"0","value":"CgQIn58B"},{"domain":".gemini.google.com","expirationDate":1772700509,"hostOnly":false,"httpOnly":false,"name":"_gcl_au","path":"/","sameSite":"unspecified","secure":false,"session":false,"storeId":"0","value":"1.1.1201602049.1764924509"},{"domain":".google.com","expirationDate":1781318085.3914,"hostOnly":false,"httpOnly":true,"name":"__Secure-BUCKET","path":"/","sameSite":"unspecified","secure":true,"session":false,"storeId":"0","value":"CJEE"},{"domain":".google.com","expirationDate":1804065999.055076,"hostOnly":false,"httpOnly":false,"name":"SID","path":"/","sameSite":"unspecified","secure":false,"session":false,"storeId":"0","value":"g.a0006AjibdaQVqq4Hq8WWhD2YY5xFfNnq1CWwmF_sC5gr51D_KinIb5a_nR8yfVUkuhqxkqM1QACgYKAcwSARcSFQHGX2Min0E6m0Q5CCXssJuNsBhwjRoVAUF8yKow9Dqjc9rubxvylwwVDtRR0076"},{"domain":".google.com","expirationDate":1804065999.055173,"hostOnly":false,"httpOnly":true,"name":"__Secure-1PSID","path":"/","sameSite":"unspecified","secure":true,"session":false,"storeId":"0","value":"g.a0006AjibdaQVqq4Hq8WWhD2YY5xFfNnq1CWwmF_sC5gr51D_KinWmAe85sr1EeGvm8xcdwOOQACgYKAegSARcSFQHGX2MihH0cUTz5slHVtOj6xK5ofhoVAUF8yKpZkzlXK6ikpamwrfX5WONv0076"},{"domain":".google.com","expirationDate":1804065999.055241,"hostOnly":false,"httpOnly":true,"name":"__Secure-3PSID","path":"/","sameSite":"no_restriction","secure":true,"session":false,"storeId":"0","value":"g.a0006AjibdaQVqq4Hq8WWhD2YY5xFfNnq1CWwmF_sC5gr51D_KinlhbPqiQV-4DT53c98YWJPQACgYKAeMSARcSFQHGX2Milf8o5HWicfaC-6DUFbA58BoVAUF8yKqE-QtDsNVpzENzW2ssbqGO0076"},{"domain":".google.com","expirationDate":1804065999.055437,"hostOnly":false,"httpOnly":true,"name":"HSID","path":"/","sameSite":"unspecified","secure":false,"session":false,"storeId":"0","value":"Avba3tjsNhxO5Rh7R"},{"domain":".google.com","expirationDate":1804065999.055485,"hostOnly":false,"httpOnly":true,"name":"SSID","path":"/","sameSite":"unspecified","secure":true,"session":false,"storeId":"0","value":"AWDnwcgnTKxK6a4LG"},{"domain":".google.com","expirationDate":1804065999.055532,"hostOnly":false,"httpOnly":false,"name":"APISID","path":"/","sameSite":"unspecified","secure":false,"session":false,"storeId":"0","value":"7lf7hhskmugrXdN6/AB2x1WXqQO9XHn9VJ"},{"domain":".google.com","expirationDate":1804065999.055583,"hostOnly":false,"httpOnly":false,"name":"SAPISID","path":"/","sameSite":"unspecified","secure":true,"session":false,"storeId":"0","value":"oLc0IqujjaTyCyvP/Ar3Fa5PiNr8k9DhO5"},{"domain":".google.com","expirationDate":1804065999.055634,"hostOnly":false,"httpOnly":false,"name":"__Secure-1PAPISID","path":"/","sameSite":"unspecified","secure":true,"session":false,"storeId":"0","value":"oLc0IqujjaTyCyvP/Ar3Fa5PiNr8k9DhO5"},{"domain":".google.com","expirationDate":1804065999.055684,"hostOnly":false,"httpOnly":false,"name":"__Secure-3PAPISID","path":"/","sameSite":"no_restriction","secure":true,"session":false,"storeId":"0","value":"oLc0IqujjaTyCyvP/Ar3Fa5PiNr8k9DhO5"},{"domain":".google.com","expirationDate":1772192423.670618,"hostOnly":false,"httpOnly":true,"name":"AEC","path":"/","sameSite":"lax","secure":true,"session":false,"storeId":"0","value":"AaJma5tm2zsNU6IZnOe3-i78-5bRpXc-AucoRmM2QqtfAADlF1jmm44DO80"},{"domain":".google.com","expirationDate":1787463223.816265,"hostOnly":false,"httpOnly":true,"name":"NID","path":"/","sameSite":"no_restriction","secure":true,"session":false,"storeId":"0","value":"529=shbFcq2GZZIQrOz40q1hguUmY-cNMFq-y7UQua9pS6LEw_P2sIEE2HlFvhZO5XU9z4AgH1xuDxATxbtZL97H0i5o4wGiywb9ZcSWoRvJNNstTUX1L_kGSgRbj-4lTtOyVxwFLG3pVzdTF1xr90LEhQLS7O9DzsK9bBSTkLJjmkkGS5RZ6r0nz2BfF2wPdaFNS7nRTvTGsuDmDFJMILlIhtaY9TWn7y1rzC-MA6Anih-MAeeJ8AhjukJrk7ChAIjD-dQQgl9iMJdDR5v3623-jeWp8vUZr_jcR17igPCa-uKOhBdtDfCSOAhdp3Qqv-uOgAzt5lm-u18zGXhPa1zq_880EUBgu0KLe7oXCPEhgUIVLSBBaITe_gCEYg"},{"domain":".gemini.google.com","expirationDate":1806298716.619349,"hostOnly":false,"httpOnly":false,"name":"_ga_WC57KJ50ZZ","path":"/","sameSite":"unspecified","secure":false,"session":false,"storeId":"0","value":"GS2.1.s1771738716$o96$g0$t1771738716$j60$l0$h0"},{"domain":".gemini.google.com","expirationDate":1806298716.631293,"hostOnly":false,"httpOnly":false,"name":"_ga_BF8Q35BMLM","path":"/","sameSite":"unspecified","secure":false,"session":false,"storeId":"0","value":"GS2.1.s1771738716$o73$g0$t1771738716$j60$l0$h0"},{"domain":".google.com","expirationDate":1803275018.283627,"hostOnly":false,"httpOnly":true,"name":"__Secure-1PSIDTS","path":"/","sameSite":"unspecified","secure":true,"session":false,"storeId":"0","value":"sidts-CjEBBj1CYpp0S_oeMN4LjLC6dWPcHV3MEzxyaW9nnSWmDL1B0_IynUEKnHj3xX0UvZfPEAA"},{"domain":".google.com","expirationDate":1803275018.284028,"hostOnly":false,"httpOnly":true,"name":"__Secure-3PSIDTS","path":"/","sameSite":"no_restriction","secure":true,"session":false,"storeId":"0","value":"sidts-CjEBBj1CYpp0S_oeMN4LjLC6dWPcHV3MEzxyaW9nnSWmDL1B0_IynUEKnHj3xX0UvZfPEAA"},{"domain":".google.com","expirationDate":1803275094.144819,"hostOnly":false,"httpOnly":false,"name":"SIDCC","path":"/","sameSite":"unspecified","secure":false,"session":false,"storeId":"0","value":"AKEyXzUyUGzhlgMTVhq4QAABraFH6CM80SIVbL_Rfv437B2YqMGpoJxkCGske_Zj8Fdf-oDzBA"},{"domain":".google.com","expirationDate":1803275094.145246,"hostOnly":false,"httpOnly":true,"name":"__Secure-1PSIDCC","path":"/","sameSite":"unspecified","secure":true,"session":false,"storeId":"0","value":"AKEyXzWulmWc4kGPMFk4DvUGk3KCdTmh9CCBWi8P5fE0R9w1t2quCyA9qkI_oxZ-FmkXFvK0FRU"},{"domain":".google.com","expirationDate":1803275094.145628,"hostOnly":false,"httpOnly":true,"name":"__Secure-3PSIDCC","path":"/","sameSite":"no_restriction","secure":true,"session":false,"storeId":"0","value":"AKEyXzW11BwWiAudVBOz4UzUXO8M7L35kqpCJ614LJtjcBrIsUEHPO1mUYlzY_3TuHjkwk_3zA"}]}`

		return newGeminiAgent(cookie), nil
	default:
		return nil, fmt.Errorf("loại agent không hỗ trợ: %s", agentType)
	}
}
