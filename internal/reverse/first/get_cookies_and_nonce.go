package first

import (
	"context"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"log/slog"
	"moodle_attendo/internal/model"
	"net/http"
	"time"
)

func GetCookiesAndNonce(ctx context.Context, client http.Client, hostname string) (*model.AuthModel, error) {
	alert := true

	for {
		select {
		case <-ctx.Done():
			return nil, errors.New("program dihentikan karena server tidak menanggapi request")
		default:
			req, err := http.NewRequest("GET", "https://"+hostname+"/login/index.php", nil)
			if err != nil {
				return nil, err
			}

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

			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				if alert {
					slog.Error("masalah jaringan: " + err.Error())
					alert = false
				}
				time.Sleep(1 * time.Second)
				continue
			}

			var nonce string
			doc.Find(".login-form input[name='logintoken']").Each(func(i int, s *goquery.Selection) {
				nonce, _ = s.Attr("value")
			})
			if nonce == "" {
				return nil, errors.New("hostname tidak disupport oleh program ini")
			}

			cookies := resp.Cookies()

			var moodleSession string
			for _, cookie := range cookies {
				if cookie.Name == "MoodleSession" {
					moodleSession = cookie.Value
				}
			}

			return &model.AuthModel{MoodleSession: moodleSession, Nonce: nonce}, nil
		}
	}
}
