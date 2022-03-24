package convert

import (
	"context"
	"log"
	"net/url"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func ConvertHTMLToPDF(u *url.URL) []byte {
	var buf []byte
	// log the CDP messages so that you can find the one to use.
	ctx, cancel := chromedp.NewContext(context.Background(), chromedp.WithDebugf(log.Printf))
	defer cancel()

	if err := chromedp.Run(ctx,
		chromedp.Navigate(u.String()),
		chromedp.Sleep(500*time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
			b, _, err := page.PrintToPDF().WithPrintBackground(true).Do(ctx)
			if err != nil {
				return err
			}

			buf = b

			return nil
		}),
	); err != nil {
		log.Fatal(err)
	}

	return buf
}
