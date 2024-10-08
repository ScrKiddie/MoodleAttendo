package util

import (
	"context"
	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"log/slog"
	"time"
)

func TakeScreenshot(ctx context.Context, session string, link string, hostname string) ([]byte, error) {
	var buf []byte

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.Flag("window-size", "1980,1080"),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(ctx, opts...)
	defer cancel()

	ctx, cancel = chromedp.NewContext(allocCtx, chromedp.WithLogf(slog.Info))
	defer cancel()

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

		return nil, err
	}
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
