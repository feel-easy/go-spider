package main

import (
	"fmt"

	"github.com/crawlab-team/crawlab-go-sdk/entity"
	"github.com/gocolly/colly/v2"
)

func main() {
	// 生成 colly 采集器
	c := colly.NewCollector(
		colly.AllowedDomains("www.baidu.com"),
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_14_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.135 Safari/537.36"),
	)

	// 抓取结果数据钩子函数
	c.OnHTML(".result.c-container", func(e *colly.HTMLElement) {
		// 抓取结果实例
		item := entity.Item{
			"title": e.ChildText("h3.t > a"),
			"url":   e.ChildAttr("h3.t > a", "href"),
		}

		// 打印抓取结果
		fmt.Println(item)

		// 取消注释调用 Crawlab Go SDK 存入数据库
		//_ = crawlab.SaveItem(item)
	})

	// 分页钩子函数
	c.OnHTML("a.n", func(e *colly.HTMLElement) {
		_ = c.Visit("https://www.baidu.com" + e.Attr("href"))
	})

	// 访问初始 URL
	startUrl := "https://www.baidu.com/s?wd=crawlab"
	_ = c.Visit(startUrl)

	// 等待爬虫结束
	c.Wait()
}
