package third

import (
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"time"
)

func PresenceProcessFourth(ctx context.Context, client http.Client, payloadPart map[string]string, session string, hostname string) error {
	alert := true

	payload := url.Values{}
	payload.Set("sessid", payloadPart["sessid"])
	payload.Add("sesskey", payloadPart["sesskey"])
	payload.Set("_qf__mod_attendance_form_studentattendance", "1")
	payload.Set("mform_isexpanded_id_session", "1")
	payload.Set("status", payloadPart["status"])
	payload.Set("submitbutton", "Save changes")

	for {
		select {
		case <-ctx.Done():
			return errors.New("program dihentikan karena server tidak menanggapi request")
		default:
			req, err := http.NewRequest("POST", "https://"+hostname+"/mod/attendance/attendance.php", bytes.NewBufferString(payload.Encode()))
			if err != nil {
				if alert {
					slog.Error("masalah jaringan: " + err.Error())
					alert = false
				}
				time.Sleep(1 * time.Second)
				continue
			}

			req.Header.Add("Host", hostname)
			req.Header.Add("Cookie", "MoodleSession="+session)
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

			return nil
		}
	}
}
