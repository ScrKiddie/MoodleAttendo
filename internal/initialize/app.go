package initialize

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"moodle_attendo/internal/model"
	"moodle_attendo/internal/reverse/first"
	"moodle_attendo/internal/reverse/second"
	"moodle_attendo/internal/reverse/third"
	"moodle_attendo/internal/util"
	"net/http"
	"time"
)

func App(ctx context.Context, client http.Client, courseId string, account model.AccountModel) {

	cookiesAndNonce, err := first.GetCookiesAndNonce(ctx, client, account.Hostname)
	if err != nil {
		log.Fatal(err)
	}

	realCookies, err := second.GetRealCookies(ctx, client, cookiesAndNonce, account)
	if err != nil {
		log.Fatal(err)
	}

	if courseId == "testing" {
		currentTime := time.Now().In(time.FixedZone("WIB", 7*60*60))
		message := fmt.Sprintf("[%s] %s", currentTime.Format("2006-01-02 15:04:05"), "berhasil terhubung dengan telegram")
		screenshot, err := util.TakeScreenshot(ctx, realCookies, "https://"+account.Hostname+"/user/profile.php", account.Hostname)
		if err != nil {
			slog.Warn("gagal melakukan screenshot: " + err.Error())
			html, err := util.ExportHtml(ctx, client, "https://"+account.Hostname+"/user/profile.php", realCookies)
			if err != nil {
				slog.Warn("gagal export html: " + err.Error())
			} else {
				if err = util.SendDocument(ctx, client, account.BotToken, account.ChatId, html, fmt.Sprintf("%s.html", time.Now().Format("2006-01-02_15-04-05")), ""); err != nil {
					slog.Warn(err.Error())
				}
			}
		} else {
			if err = util.SendDocument(ctx, client, account.BotToken, account.ChatId, screenshot, fmt.Sprintf("%s.png", time.Now().Format("2006-01-02_15-04-05")), ""); err != nil {
				slog.Warn(err.Error())
			}
		}
		if err = util.SendMessage(ctx, client, account.BotToken, account.ChatId, message); err != nil {
			slog.Warn(err.Error())
		}
		slog.Info(message)
		return
	}

	arrLink, err := third.PresenceProcessFirst(ctx, client, courseId, realCookies, account.Hostname)
	if err != nil {
		currentTime := time.Now().In(time.FixedZone("WIB", 7*60*60))
		message := fmt.Sprintf("[%s] %s", currentTime.Format("2006-01-02 15:04:05"), err.Error())
		if err := util.SendMessage(ctx, client, account.BotToken, account.ChatId, message); err != nil {
			slog.Warn(err.Error())
		}
		log.Fatal(message)
	}

	var formLink string
	var link string
	if len(arrLink) > 0 {
		link, formLink, err = third.PresenceProsesSecond(ctx, client, arrLink, realCookies, account.Hostname)
		if err != nil {
			currentTime := time.Now().In(time.FixedZone("WIB", 7*60*60))
			message := fmt.Sprintf("[%s] %s", currentTime.Format("2006-01-02 15:04:05"), err.Error())
			if err := util.SendMessage(ctx, client, account.BotToken, account.ChatId, message); err != nil {
				slog.Warn(err.Error())
			}
			log.Fatal(message)
		}
	}

	payloadPart := new(map[string]string)
	if formLink != "" {
		payloadPart, err = third.PresenceProcessThird(ctx, client, formLink, realCookies)
		if err != nil {
			currentTime := time.Now().In(time.FixedZone("WIB", 7*60*60))
			message := fmt.Sprintf("[%s] %s", currentTime.Format("2006-01-02 15:04:05"), err.Error())
			if err := util.SendMessage(ctx, client, account.BotToken, account.ChatId, message); err != nil {
				slog.Warn(err.Error())
			}
			log.Fatal(message)
		}
		if payloadPart == nil {
			currentTime := time.Now().In(time.FixedZone("WIB", 7*60*60))
			message := fmt.Sprintf("[%s] berhasil melakukan presensi pada %s", currentTime.Format("2006-01-02 15:04:05"), link)
			if err := util.CloseSidebar(ctx, client, realCookies, account.Hostname); err != nil {
				slog.Warn(err.Error())
			}
			screenshot, err := util.TakeScreenshot(ctx, realCookies, link, account.Hostname)
			if err != nil {
				slog.Warn("gagal melakukan screenshot: " + err.Error())
				html, err := util.ExportHtml(ctx, client, link, realCookies)
				if err != nil {
					slog.Warn("gagal export html: " + err.Error())
					if err := util.SendMessage(ctx, client, account.BotToken, account.ChatId, message); err != nil {
						slog.Warn(err.Error())
					}
				} else {
					if err := util.SendDocument(ctx, client, account.BotToken, account.ChatId, html, fmt.Sprintf("%s.html", time.Now().Format("2006-01-02_15-04-05")), message); err != nil {
						slog.Warn(err.Error())
					}
				}
			} else {
				if err := util.SendDocument(ctx, client, account.BotToken, account.ChatId, screenshot, fmt.Sprintf("%s.png", time.Now().Format("2006-01-02_15-04-05")), message); err != nil {
					slog.Warn(err.Error())
				}
			}
			slog.Info(message)
			return
		}
	}

	if err := third.PresenceProcessFourth(ctx, client, *payloadPart, realCookies, account.Hostname); err != nil {
		currentTime := time.Now().In(time.FixedZone("WIB", 7*60*60))
		message := fmt.Sprintf("[%s] %s", currentTime.Format("2006-01-02 15:04:05"), err.Error())
		if err := util.SendMessage(ctx, client, account.BotToken, account.ChatId, message); err != nil {
			slog.Warn(err.Error())
		}
		log.Fatal(message)
	} else {
		currentTime := time.Now().In(time.FixedZone("WIB", 7*60*60))
		message := fmt.Sprintf("[%s] berhasil melakukan presensi pada %s", currentTime.Format("2006-01-02 15:04:05"), link)
		if err := util.CloseSidebar(ctx, client, realCookies, account.Hostname); err != nil {
			slog.Warn(err.Error())
		}
		screenshot, err := util.TakeScreenshot(ctx, realCookies, link, account.Hostname)
		if err != nil {
			slog.Warn("gagal melakukan screenshot: " + err.Error())
			html, err := util.ExportHtml(ctx, client, link, realCookies)
			if err != nil {
				slog.Warn("gagal export html: " + err.Error())
				if err := util.SendMessage(ctx, client, account.BotToken, account.ChatId, message); err != nil {
					slog.Warn(err.Error())
				}
			} else {
				if err := util.SendDocument(ctx, client, account.BotToken, account.ChatId, html, fmt.Sprintf("%s.html", time.Now().Format("2006-01-02_15-04-05")), message); err != nil {
					slog.Warn(err.Error())
				}
			}
		} else {
			if err := util.SendDocument(ctx, client, account.BotToken, account.ChatId, screenshot, fmt.Sprintf("%s.png", time.Now().Format("2006-01-02_15-04-05")), message); err != nil {
				slog.Warn(err.Error())
			}
		}
		slog.Info(message)
		return
	}
}
