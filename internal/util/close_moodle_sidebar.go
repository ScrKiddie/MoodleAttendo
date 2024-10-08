package util

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

func CloseSidebar(ctx context.Context, client http.Client, session string, hostname string) error {
	alert := true

	data := []map[string]interface{}{
		{
			"index":      0,
			"methodname": "core_user_set_user_preferences",
			"args": map[string]interface{}{
				"preferences": []map[string]interface{}{
					{
						"name":   "drawer-open-index",
						"value":  false,
						"userid": 0,
					},
				},
			},
		},
	}

	for {
		select {
		case <-ctx.Done():
			return errors.New("close sidebar dihentikan karena server tidak menanggapi request")
		default:
			jsonData, err := json.Marshal(data)
			if err != nil {
				slog.Warn("gagal encode JSON: " + err.Error())
			}
			req, err := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("https://%s/lib/ajax/service.php?info=core_user_set_user_preferences", hostname), bytes.NewBuffer(jsonData))
			if err != nil {
				if alert {
					slog.Error("masalah jaringan: " + err.Error())
					alert = false
				}
				time.Sleep(1 * time.Second)
				continue
			}

			req.Header.Set("Content-Type", "application/json")
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

			if resp.StatusCode != http.StatusOK {
				return errors.New("gagal menutup sidebar: " + err.Error())
			}

			return nil
		}
	}
}
