package second

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"moodle_attendo/internal/model"
	"net/http"
	"net/url"
	"time"
)

func GetRealCookies(ctx context.Context, client http.Client, a *model.AuthModel, b model.AccountModel) (string, error) {
	alert := true

	payload := url.Values{}
	payload.Set("logintoken", a.Nonce)
	payload.Set("username", b.Username)
	payload.Set("password", b.Password)

	for {
		select {
		case <-ctx.Done():
			return "", errors.New("program dihentikan karena server tidak menanggapi request")
		default:
			payloadBytes := []byte(payload.Encode())
			req, err := http.NewRequest("POST", "https://"+b.Hostname+"/login/index.php", bytes.NewBuffer(payloadBytes))
			if err != nil {
				if alert {
					slog.Error("masalah jaringan: " + err.Error())
					alert = false
				}
				time.Sleep(1 * time.Second)
				continue
			}

			req.Header.Add("Host", b.Hostname)
			req.Header.Add("Cookie", "MoodleSession="+a.MoodleSession)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.5112.102 Safari/537.36")

			resp, err := client.Do(req)
			if err != nil {
				if alert {
					slog.Error("masalah jaringan: " + err.Error())
					alert = false
				}
				time.Sleep(1 * time.Second)
				continue
			}
			defer resp.Body.Close()

			cookies := resp.Cookies()
			for _, cookie := range cookies {
				if cookie.Name == "MoodleSession" {
					return cookie.Value, nil
				}
			}

			return "", errors.New("username atau password salah")
		}
	}
}
