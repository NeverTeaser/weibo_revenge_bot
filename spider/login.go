package spider

import (
	"context"
	"io/ioutil"
	"log"
	"time"

	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
)

func main() {
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", false),
		chromedp.Flag("mute-audio", false),
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()

	var screenShot []byte
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	var imgUrl string
	var cookieStr string

	err := chromedp.Run(ctx,
		chromedp.Navigate(`https://weibo.com`),

		chromedp.WaitVisible(`.UG_box`),
		chromedp.Sleep(time.Second/2),
		chromedp.Click(`//div[@class='info_header']//a[@node-type='qrcode_tab']`, chromedp.NodeVisible),
		chromedp.WaitVisible(`//img[@node-type='qrcode_src']`),
		chromedp.OuterHTML(`//img[@node-type='qrcode_src']/@src`, &imgUrl, chromedp.BySearch),
		// 等待扫码成功
		chromedp.WaitVisible(`//div[@node-type='scan_success']`, chromedp.NodeVisible),
		// 等待登陆成功
		chromedp.Sleep(time.Second/2),
		chromedp.WaitVisible(`//div[@class='nameBox']`, chromedp.NodeVisible),
		chromedp.Sleep(time.Second/2),
		chromedp.ActionFunc(func(ctx context.Context) error {
			cookies, err := network.GetAllCookies().Do(ctx)
			if err != nil {
				return err
			}
			for _, v := range cookies {
				cookieStr += v.Name + "=" + v.Value + ";"
			}
			return nil
		}),
		chromedp.CaptureScreenshot(&screenShot),
	)
	if err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile("fullScreenshot.png", screenShot, 0644); err != nil {
		log.Fatal(err)
	}

}
