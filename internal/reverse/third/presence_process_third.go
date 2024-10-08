package third

import (
	"context"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"log/slog"
	"net/http"
	"time"
)

func PresenceProcessThird(ctx context.Context, client http.Client, link string, session string) (*map[string]string, error) {
	alert := true

	for {
		select {
		case <-ctx.Done():
			return nil, errors.New("program dihentikan karena server tidak menanggapi request")
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

			firstElement := doc.Find("label.form-check-inline").First()
			input := firstElement.Find("input[type='radio']")
			status, exists := input.Attr("value")

			if status == "" {
				slog.Warn("pilihan pada presensi tidak ditemukan sehingga presensi dianggap sukses")
				return nil, nil
			}

			sessid, exists := doc.Find("input[name='sessid']").Attr("value")
			if !exists {
				return nil, errors.New("sessid tidak ditemukan")
			}

			sesskey, exists := doc.Find("input[name='sesskey']").Attr("value")
			if !exists {
				return nil, errors.New("sesskey tidak ditemukan")
			}

			return &map[string]string{
				"sessid":  sessid,
				"sesskey": sesskey,
				"status":  status,
			}, nil
		}
	}
}
