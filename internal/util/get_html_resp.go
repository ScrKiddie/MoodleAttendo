package util

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"time"
)

func ExportHtml(ctx context.Context, client http.Client, link string, session string) ([]byte, error) {
	alert := true

	for {
		select {
		case <-ctx.Done():
			return []byte{}, errors.New("program dihentikan karena server tidak menanggapi request")
		default:
			req, err := http.NewRequest("GET", link, nil)
			if err != nil {
				if alert {
					slog.Error("masalah jaringan: " + err.Error())
					alert = false
				}
				time.Sleep(1 * time.Second)
				continue
			}

			req.Header.Set("Cookie", "MoodleSession="+session)

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

			html, err := io.ReadAll(resp.Body)
			if err != nil {
				if alert {
					slog.Error("masalah jaringan: " + err.Error())
					alert = false
				}
				time.Sleep(1 * time.Second)
				continue
			}
			return html, nil
		}
	}
}
