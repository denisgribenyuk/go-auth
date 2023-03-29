package notificator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func SendNotification(email string, title string, msg string) error {

	message := make(map[string]string)
	message["email"] = email
	message["title"] = title
	message["message"] = msg
	message["token"] = os.Getenv("SERVICE_SECRET_KEY")

	messageBody, _ := json.Marshal(message)

	notificationUrl := os.Getenv("BACK_NOTIFICATION_URL") + "/send_notification"

	notificationRequest, err := http.NewRequest(http.MethodPost, notificationUrl, bytes.NewBuffer(messageBody))
	if err != nil {
		return fmt.Errorf("error form notification request: %v", err)
	}

	notificationRequest.Header.Set("Content-Type", "application/json")

	client := http.Client{Timeout: time.Duration(10) * time.Second}

	resp, err := client.Do(notificationRequest)
	if err != nil {
		return fmt.Errorf("error send notification request: %v", err)
	}
	defer resp.Body.Close()

	br, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading send notification response: %v. Body: %s", err, string(br))
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("error in send notification response: %s.\nStatus code: %v", br, resp.StatusCode)
	}

	return nil
}
