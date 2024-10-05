package third

import (
	"context"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"log/slog"
	"net/http"
	"time"
)

func PresenceProsesSecond(ctx context.Context, client http.Client, link string, session string, hostname string) (string, error) {
	alert := true

	for {
		select {
		case <-ctx.Done():
			return "", errors.New("program dihentikan karena server tidak menanggapi request")
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

			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				if alert {
					slog.Error("masalah jaringan: " + err.Error())
					alert = false
				}
				time.Sleep(1 * time.Second)
				continue
			}

			url, exists := doc.Find("a[href*='https://" + hostname + "/mod/attendance/attendance.php']").Attr("href")
			if !exists {
				return "", errors.New("belum ada presensi di course " + link)
			}
			return url, nil
		}
	}
}
