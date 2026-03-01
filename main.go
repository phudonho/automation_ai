package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"helper/pkg/automation"
	"helper/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/go-vgo/robotgo"
	"github.com/joho/godotenv"
)

func main() {
	// Khởi tạo env
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Warning: Không tìm thấy file .env, sẽ đọc thẳng từ biến hệ thống.")
	}

	// Khởi tạo client gemini với Cookie
	cookie := os.Getenv("GEMINI_COOKIE")
	if cookie == "" {
		fmt.Println("[LỖI] Vui lòng cấu hình biến GEMINI_COOKIE trong file .env trước khi định tuyến AI.")
		return
	}

	agent, err := automation.NewAgent(automation.Gemini, cookie)
	if err != nil {
		fmt.Printf("Lỗi khởi tạo agent: %v\n", err)
		return
	}

	executor, err := automation.NewExecutor(automation.RobotGoExecutor)
	if err != nil {
		fmt.Printf("Lỗi khởi tạo executor: %v\n", err)
		return
	}

	// Yêu cầu người dùng setup ENV Token cho Telegram Bot
	token := os.Getenv("TELEGRAM_BOT_TOKEN")
	if token == "" || token == "your_telegram_bot_token" {
		fmt.Println("[LỖI] Vui lòng cấu hình biến môi trường TELEGRAM_BOT_TOKEN trong file .env trước khi chạy chương trình.")
		return
	}

	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		fmt.Printf("[LỖI] Kết nối Telegram thất bại: %v\n", err)
		return
	}

	bot.Debug = false
	fmt.Printf("[THÔNG BÁO] Hệ thống Helper Bot đã khởi động. Xin chào: %s\n", bot.Self.UserName)
	fmt.Println("Vui lòng mở Telegram Chat với Bot này và gửi lệnh (vd: 'mở facebook và like bài mới').")

	// Cấu hình ID của Group/Channel nếu người dùng muốn giới hạn Bot
	allowedChatIDStr := os.Getenv("TELEGRAM_CHAT_ID")
	var allowedChatID int64
	if allowedChatIDStr != "" {
		allowedChatID, _ = strconv.ParseInt(allowedChatIDStr, 10, 64)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	// Ép buộc API GetUpdate của Telegram gửi cả 4 loại tin nhắn & Sự kiện Click nút
	u.AllowedUpdates = []string{"message", "edited_message", "channel_post", "edited_channel_post", "callback_query"}
	updates := bot.GetUpdatesChan(u)

	// Single thread automation flag (để không nhận 2 task cùng lúc)
	isDoingTask := false

	for update := range updates {
		var msg *tgbotapi.Message

		// Bắt thông điệp từ chat cá nhân, Group hoặc Channel (Kể cả tin nhắn bị edit)
		if update.Message != nil {
			msg = update.Message
		} else if update.ChannelPost != nil {
			msg = update.ChannelPost
		} else if update.EditedMessage != nil {
			msg = update.EditedMessage
		} else if update.EditedChannelPost != nil {
			msg = update.EditedChannelPost
		} else if update.CallbackQuery != nil {
			// Xử lý khi user bấm vào nút Menu Inline (Bàn phím ảo)
			chatID := update.CallbackQuery.Message.Chat.ID
			actionData := update.CallbackQuery.Data

			// Trả lời callback để Telegram tắt trạng thái "Loading..."
			bot.Request(tgbotapi.NewCallback(update.CallbackQuery.ID, ""))

			if actionData == "btn_task" {
				bot.Send(tgbotapi.NewMessage(chatID, "✍️ Vui lòng gõ mệnh lệnh mà bạn muốn tôi thực hiện (VD: mở trình duyệt tìm kiếm bài hát...)."))
			} else if actionData == "btn_status" {
				statusStr := "Đang rảnh rỗi (Chờ lệnh) 🟢"
				if isDoingTask {
					statusStr = "Đang bận chạy nhiệm vụ 🔴"
				}
				bot.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("Trạng thái hiện tại: *%s*", statusStr)))
			}
			continue
		} else {
			// Bỏ qua các sự kiện còn lại
			continue
		}

		if msg == nil {
			continue
		}

		chatID := msg.Chat.ID
		taskStr := msg.Text

		if taskStr == "" {
			continue
		}

		// Xử lý menu khởi tạo /start
		if taskStr == "/start" {
			welcomeMsg := tgbotapi.NewMessage(chatID, "Xin chào! Bạn muốn yêu cầu tôi làm gì?")
			// Hiển thị 2 nút Inline trong Channel
			keyboard := tgbotapi.NewInlineKeyboardMarkup(
				tgbotapi.NewInlineKeyboardRow(
					tgbotapi.NewInlineKeyboardButtonData("📥 Nhận nhiệm vụ", "btn_task"),
					tgbotapi.NewInlineKeyboardButtonData("ℹ️ Trạng thái", "btn_status"),
				),
			)
			welcomeMsg.ReplyMarkup = keyboard
			bot.Send(welcomeMsg)
			continue
		}

		// Khóa API: Chỉ làm việc với đúng Group / Channel được phép (nếu khai báo biến môi trường)
		if allowedChatID != 0 && chatID != allowedChatID {
			continue
		}

		senderName := "Channel"
		if msg.From != nil {
			senderName = msg.From.UserName
		}
		fmt.Printf("[%s] Người dùng phân công nhiệm vụ: %s\n", senderName, taskStr)

		if isDoingTask {
			msg := tgbotapi.NewMessage(chatID, "🚫 Tôi đang bận thực hiện một lệnh cũ. Vui lòng chờ lát nữa nhé!")
			bot.Send(msg)
			continue
		}

		// Nhận Task
		isDoingTask = true
		msgStart := "✅ Đã nhận lệnh bắt đầu thao tác auto: *" + taskStr + "*. Cùng xem tôi làm việc nhé!"
		msgStartTg := tgbotapi.NewMessage(chatID, msgStart)
		msgStartTg.ParseMode = "Markdown"
		bot.Send(msgStartTg)

		go func() {
			defer func() { isDoingTask = false }() // Thả cờ sau khi task kết thúc
			runAutoLoop(bot, chatID, taskStr, agent, executor)
		}()
	}
}

func runAutoLoop(bot *tgbotapi.BotAPI, chatID int64, task string, agent automation.AutomationAgent, executor automation.ActionExecutor) {
	history := []string{}
	fmt.Println("BẮT ĐẦU CHUỖI NHIỆM VỤ AUTO...")

	limiter := utils.NewUploadLimiter(12 * time.Second)
	for step := 1; ; step++ {
		fmt.Printf("\n=== BƯỚC %d ===\n", step)

		// 1. Chụp màn hình
		fmt.Println("Đang chụp màn hình...")
		screenshotPath := "current_screen.png"

		err := utils.CaptureScreen(screenshotPath)
		if err != nil {
			msg := fmt.Sprintf("❌ Lỗi chụp màn hình: %v", err)
			bot.Send(tgbotapi.NewMessage(chatID, msg))
			break
		}

		// 2. Upload file
		fmt.Println("Đang upload ảnh lên sxcu.net...")
		imageURL, errUpload := utils.UploadToSxcuWithBackoff(
			limiter,
			screenshotPath,
			5,             // maxRetries
			2*time.Second, // baseDelay
		)
		if errUpload != nil {
			msg := fmt.Sprintf("❌ Lỗi upload API: %v", errUpload)
			bot.Send(tgbotapi.NewMessage(chatID, msg))
			os.Remove(screenshotPath)
			break
		}

		// 3. Phân tích qua Gemini chuyên biệt Spatial
		fmt.Printf("Đang phân tích toạ độ nội nhúng từ ảnh gốc: %s\n", imageURL)
		action, err := agent.DecideNextAction(task, imageURL, history)
		if err != nil {
			msg := fmt.Sprintf("❌ Lỗi khi request AI: %v", err)
			bot.Send(tgbotapi.NewMessage(chatID, msg))
			os.Remove(screenshotPath)
			break
		}

		// 4. CHUYỂN ĐỔI TOẠ ĐỘ: Scale 1000x1000 sang kích thước màn hình thật
		screenWidth, screenHeight := robotgo.GetScreenSize()
		originalX, originalY := action.Coordinates.X, action.Coordinates.Y

		if action.NextStep != "FINISH" && originalX > 0 && originalY > 0 {
			action.Coordinates.X = (originalX * screenWidth) / 1000
			action.Coordinates.Y = (originalY * screenHeight) / 1000
		}

		fmt.Println("\n---- KẾT QUẢ TỪ AI ----")
		fmt.Printf("Phân tích: %s\n", action.Analysis)
		fmt.Printf("Lệnh Next Step: %s\n", action.NextStep)
		fmt.Printf("Tọa độ khung (1000s): (%d, %d) => Ánh xạ màn hình thật: (%d, %d)\n", originalX, originalY, action.Coordinates.X, action.Coordinates.Y)

		// 5. Báo cáo bằng tin nhắn (và ảnh) vào Telegram
		reportTxt := fmt.Sprintf("🤖 *Bước %d:*\n_%s_\n\n➡ Thao tác máy tính: **%s** (%d, %d)", step, action.MessageToUser, action.NextStep, action.Coordinates.X, action.Coordinates.Y)

		photoMsg := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(imageURL))
		photoMsg.Caption = reportTxt
		photoMsg.ParseMode = "Markdown"
		_, errSend := bot.Send(photoMsg)
		if errSend != nil { // Fallback về text nếu link ảnh upload lỗi định dạng với Telegram
			txtMsg := tgbotapi.NewMessage(chatID, reportTxt+"\n[Link ảnh màn hình]("+imageURL+")")
			txtMsg.ParseMode = "Markdown"
			bot.Send(txtMsg)
		}

		// Nhận lệnh Kết thúc thì thoát
		if action.NextStep == "FINISH" {
			bot.Send(tgbotapi.NewMessage(chatID, "🎉 HOÀN THÀNH NHIỆM VỤ RỒI NHÉ!"))
			os.Remove(screenshotPath)
			break
		}

		// 6. Thực thi thao tác Chuột/Phím qua OS
		err = executor.Execute(action)
		if err != nil {
			msg := fmt.Sprintf("❌ Lỗi hệ điều hành khi chạy thao tác: %v", err)
			bot.Send(tgbotapi.NewMessage(chatID, msg))
		}

		os.Remove(screenshotPath)

		// 7. Cập nhật lịch sử
		historyStr := fmt.Sprintf("Nhịp %d: Lệnh %s tại toạ độ %d,%d, nhập: %s", step, action.NextStep, action.Coordinates.X, action.Coordinates.Y, action.Value)
		history = append(history, historyStr)

		time.Sleep(3 * time.Second)
	}
}
