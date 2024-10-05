package third

import (
	"context"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"log/slog"
	"net/http"
	"time"
)

func PresenceProcessFirst(ctx context.Context, client http.Client, idCourse string, session string, hostname string) (string, error) {
	alert := true

	for {
		select {
		case <-ctx.Done():
			return "", errors.New("program dihentikan karena server tidak menanggapi request")
		default:
			url := fmt.Sprintf("https://"+hostname+"/course/view.php?id=%s", idCourse)
			req, err := http.NewRequest("GET", url, nil)
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
			var lastURL string

			doc.Find("a[href*='https://" + hostname + "/mod/attendance/view.php']").Each(func(i int, s *goquery.Selection) {
				href, _ := s.Attr("href")
				lastURL = href
			})

			if lastURL == "" {
				return "", errors.New("link presensi pada " + url + " tidak ditemukan")
			}
			return lastURL, nil
		}
	}
}
