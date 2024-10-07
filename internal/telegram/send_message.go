package telegram

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/url"
	"time"
)

func SendMessage(ctx context.Context, client http.Client, token, chat, message string) error {
	alert := true

	data := url.Values{}
	data.Set("chat_id", chat)
	data.Set("text", message)

	for {
		select {
		case <-ctx.Done():
			return errors.New("request ke telegram dihentikan karena server tidak menanggapi")
		default:
			resp, err := client.PostForm(fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token), data)
			if err != nil {
				if alert {
					slog.Error("masalah jaringan: " + err.Error())
					alert = false
				}
				time.Sleep(1 * time.Second)
				continue
			}

			defer resp.Body.Close()
			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("response %s dari telegram pastikan chat id dan bot token sudah benar", resp.Status)
			}
			return nil
		}
	}
}

func SendDocument(ctx context.Context, client http.Client, token, chat string, document []byte, filename string, caption string) error {
	alert := true

	for {
		select {
		case <-ctx.Done():
			return errors.New("request ke telegram dihentikan karena server tidak menanggapi")
		default:
			body := new(bytes.Buffer)
			writer := multipart.NewWriter(body)

			part, err := writer.CreateFormFile("document", filename)
			if err != nil {
				return fmt.Errorf("gagal membuat form file: %v", err)
			}
			if _, err := part.Write(document); err != nil {
				return fmt.Errorf("gagal menulis dokumen: %v", err)
			}
			if err := writer.WriteField("chat_id", chat); err != nil {
				return fmt.Errorf("gagal menulis chat_id: %v", err)
			}
			if err := writer.WriteField("caption", caption); err != nil {
				return fmt.Errorf("gagal menulis caption: %v", err)
			}
			if err := writer.Close(); err != nil {
				return fmt.Errorf("gagal menutup writer: %v", err)
			}

			req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("https://api.telegram.org/bot%s/sendDocument", token), body)
			if err != nil {
				if alert {
					slog.Error("masalah membuat request: " + err.Error())
					alert = false
				}
				time.Sleep(1 * time.Second)
				continue
			}

			req.Header.Set("Content-Type", writer.FormDataContentType())

			resp, err := client.Do(req)
			if err != nil {
				if alert {
					slog.Error("masalah jaringan: " + err.Error())
					alert = false
				}
				time.Sleep(1 * time.Second)
				continue
			}

			bodyBytes, _ := io.ReadAll(resp.Body)
			resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return fmt.Errorf("response %s dari telegram: %s", resp.Status, string(bodyBytes))
			}

			return nil
		}
	}
}
