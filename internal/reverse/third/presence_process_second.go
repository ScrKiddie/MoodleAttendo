package third

import (
	"context"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"log/slog"
	"net/http"
	"strings"
	"sync"
	"time"
)

func PresenceProsesSecond(ctx context.Context, client http.Client, allUrl []string, session string, hostname string) (string, string, error) {
	once := new(sync.Once)
	var urlResult, urlBenar string
	wg := new(sync.WaitGroup)

	innerCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, url := range allUrl {
		wg.Add(1)
		go asyncGetAttendUrl(ctx, innerCtx, session, client, url, hostname, once, &urlResult, wg, &urlBenar, cancel)
	}

	wg.Wait()

	if urlBenar != "" && urlResult != "" {
		return urlBenar, urlResult, nil
	}
	str := strings.Join(allUrl, "\n")
	return "", "", errors.New("belum ada presensi di link berikut: \n" + str)
}

func asyncGetAttendUrl(ctx context.Context, innerCtx context.Context, session string, client http.Client, url string, hostname string, once *sync.Once, urlResult *string, wg *sync.WaitGroup, urlBenar *string, cancel context.CancelFunc) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case <-innerCtx.Done():
			return
		default:
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				once.Do(func() {
					slog.Error("masalah jaringan: " + err.Error())
				})
				time.Sleep(1 * time.Second)
				continue
			}

			req.Header.Set("Cookie", "MoodleSession="+session)
			resp, err := client.Do(req)
			if err != nil {
				once.Do(func() {
					slog.Error("masalah jaringan: " + err.Error())
				})
				time.Sleep(1 * time.Second)
				continue
			}
			defer resp.Body.Close()

			doc, err := goquery.NewDocumentFromReader(resp.Body)
			if err != nil {
				once.Do(func() {
					slog.Error("masalah jaringan: " + err.Error())
				})
				time.Sleep(1 * time.Second)
				continue
			}

			link, exists := doc.Find("a[href*='https://" + hostname + "/mod/attendance/attendance.php']").Attr("href")
			if exists {
				once.Do(func() {
					*urlResult = link
					*urlBenar = url
					cancel()
				})
				return
			}
			return
		}
	}
}
