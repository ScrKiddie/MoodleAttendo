package chromium

import (
	"context"
	"errors"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"log/slog"
	"net"
	"time"
)

func TakeScreenshot(ctx context.Context, session string, link string, hostname string) ([]byte, error) {
	var buf []byte

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("window-size", "1920,1080"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(allocCtx, chromedp.WithLogf(slog.Info))
	defer cancel()

	alert := true

	for {
		err := chromedp.Run(ctx,
			setCookie(
				"MoodleSession",
				session,
				hostname,
				"/",
				false,
				false,
			),
			chromedp.Navigate(link),
			chromedp.WaitVisible("body", chromedp.ByQuery),
			chromedp.Screenshot("body", &buf, chromedp.NodeVisible, chromedp.ByQuery),
		)

		if err == nil {
			return buf, nil
		}

		if isTemporaryError(err) {
			if alert {
				slog.Error("masalah jaringan: " + err.Error())
				alert = false
			}
			time.Sleep(1 * time.Second)
			if ctx.Err() != nil {
				return nil, errors.New("proses screenshot dihentikan karena server tidak menanggapi request")
			}
			continue
		}

		return nil, err
	}
}

func isTemporaryError(err error) bool {
	if _, ok := err.(net.Error); ok {
		return true
	}
	return false
}

func setCookie(name, value, domain, path string, httpOnly, secure bool) chromedp.Action {
	return chromedp.ActionFunc(func(ctx context.Context) error {
		expr := cdp.TimeSinceEpoch(time.Now().Add(180 * 24 * time.Hour))
		err := network.SetCookie(name, value).
			WithExpires(&expr).
			WithDomain(domain).
			WithPath(path).
			WithHTTPOnly(httpOnly).
			WithSecure(secure).
			Do(ctx)
		if err != nil {
			return err
		}
		return nil
	})
}
