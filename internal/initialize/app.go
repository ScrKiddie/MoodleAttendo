package initialize

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"moodle_attendo/internal/chromium"
	"moodle_attendo/internal/model"
	"moodle_attendo/internal/reverse/first"
	"moodle_attendo/internal/reverse/second"
	"moodle_attendo/internal/reverse/third"
	"moodle_attendo/internal/telegram"
	"net/http"
	"time"
)

func App(ctx context.Context, client http.Client, courseId string, account model.AccountModel) {

	cookiesAndNonce, err := first.GetCookiesAndNonce(ctx, client, account.Hostname)
	if err != nil {
		log.Fatal(err.Error())
	}

	realCookies, err := second.GetRealCookies(ctx, client, cookiesAndNonce, account)
	if err != nil {
		log.Fatal(err.Error())
	}

	if courseId == "testing" {
		currentTime := time.Now().In(time.FixedZone("WIB", 7*60*60))
		message := fmt.Sprintf("[%s] %s", currentTime.Format("2006-01-02 15:04:05"), "berhasil terhubung dengan telegram")
		screenshot, err := chromium.TakeScreenshot(ctx, realCookies, "https://jollycontrarian.com/images/6/6c/Rickroll.jpg", account.Hostname)
		if err != nil {
			slog.Info(message)
			slog.Warn("gagal melakukan screenshot dengan chromium: " + err.Error())
		}
		err = telegram.SendDocument(ctx, client, account.BotToken, account.ChatId, screenshot, fmt.Sprintf("%s.png", time.Now().Format("2006-01-02_15-04-05")), "")
		if err != nil {
			slog.Info(message)
			slog.Warn(err.Error())
		}
		err2 := telegram.SendMessage(ctx, client, account.BotToken, account.ChatId, message)
		if err2 != nil {
			slog.Info(message)
			slog.Warn(err2.Error())
		}
		return
	}

	link, err := third.PresenceProcessFirst(ctx, client, courseId, realCookies, account.Hostname)
	if err != nil {
		currentTime := time.Now().In(time.FixedZone("WIB", 7*60*60))
		message := fmt.Sprintf("[%s] %s", currentTime.Format("2006-01-02 15:04:05"), err.Error())
		err2 := telegram.SendMessage(ctx, client, account.BotToken, account.ChatId, message)
		if err2 != nil {
			slog.Warn(err2.Error())
		}
		log.Fatal(err.Error())
	}

	var formLink string
	if link != "" {
		formLink, err = third.PresenceProsesSecond(ctx, client, link, realCookies, account.Hostname)
		if err != nil {
			currentTime := time.Now().In(time.FixedZone("WIB", 7*60*60))
			message := fmt.Sprintf("[%s] %s", currentTime.Format("2006-01-02 15:04:05"), err.Error())
			err2 := telegram.SendMessage(ctx, client, account.BotToken, account.ChatId, message)
			if err2 != nil {
				slog.Warn(err2.Error())
			}
			log.Fatal(err.Error())
		}
	}

	payloadPart := new(map[string]string)
	if formLink != "" {
		payloadPart, err = third.PresenceProcessThird(ctx, client, formLink, realCookies)
		if err != nil {
			log.Fatal(err.Error())
		}
		if payloadPart == nil {
			currentTime := time.Now().In(time.FixedZone("WIB", 7*60*60))
			message := fmt.Sprintf("[%s] berhasil melakukan presensi pada %s", currentTime.Format("2006-01-02 15:04:05"), link)
			screenshot, err := chromium.TakeScreenshot(ctx, realCookies, link, account.Hostname)
			if err != nil {
				slog.Warn("gagal melakukan screenshot, pastikan chromium sudah terinstall dengan benar: " + err.Error())
				return
			}
			err = telegram.SendDocument(ctx, client, account.BotToken, account.ChatId, screenshot, fmt.Sprintf("%s.png", time.Now().Format("2006-01-02_15-04-05")), message)
			if err != nil {
				slog.Warn(err.Error())
				return
			}
		}
	}

	if payloadPart != nil {
		if err := third.PresenceProcessFourth(ctx, client, *payloadPart, realCookies, account.Hostname); err != nil {
			currentTime := time.Now().In(time.FixedZone("WIB", 7*60*60))
			message := fmt.Sprintf("[%s] %s", currentTime.Format("2006-01-02 15:04:05"), err.Error())
			err2 := telegram.SendMessage(ctx, client, account.BotToken, account.ChatId, message)
			if err2 != nil {
				slog.Warn(err2.Error())
			}
			log.Fatal(err.Error())
		} else {
			currentTime := time.Now().In(time.FixedZone("WIB", 7*60*60))
			message := fmt.Sprintf("[%s] berhasil melakukan presensi pada %s", currentTime.Format("2006-01-02 15:04:05"), link)

			screenshot, err := chromium.TakeScreenshot(ctx, realCookies, link, account.Hostname)
			if err != nil {
				slog.Info(message)
				slog.Warn("gagal melakukan screenshot, pastikan chromium sudah terinstall dengan benar: " + err.Error())
				return
			}
			err = telegram.SendDocument(ctx, client, account.BotToken, account.ChatId, screenshot, fmt.Sprintf("%s.png", time.Now().Format("2006-01-02_15-04-05")), message)
			if err != nil {
				slog.Info(message)
				slog.Warn(err.Error())
				return
			}
		}
	}
}
