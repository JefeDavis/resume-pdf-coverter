package convert

import (
	"context"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

func ConvertHTMLToPDF(u *url.URL) []byte {
	var buf []byte

	// Run Chrome without the setuid sandbox and without /dev/shm so the
	// container can drop the SYS_ADMIN capability and does not require a
	// large shared-memory mount.
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoSandbox,
		chromedp.Flag("disable-dev-shm-usage", true),
	)

	// browserless v2 ships Chromium outside of PATH; let the image point us
	// at it explicitly.
	if path := os.Getenv("CHROME_PATH"); path != "" {
		opts = append(opts, chromedp.ExecPath(path))
	}

	allocCtx, allocCancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer allocCancel()

	// log the CDP messages so that you can find the one to use.
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithDebugf(log.Printf))
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
