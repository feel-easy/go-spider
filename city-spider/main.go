package main

import (
	"encoding/csv"
	"log"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly"
)

type Movie struct {
	idx    string
	title  string
	year   string
	info   string
	rating string
	url    string
}

type City struct {
}

var startUrl = "http://www.ku51.net/area/"

// 起始Url

func main() {
	// 存储文件名
	fName := "area.csv"
	file, err := os.Create(fName)
	if err != nil {
		log.Fatalf("创建文件失败 %q: %s\n", fName, err)
		return
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()
	// 写CSV头部
	writer.Write([]string{"省", "市", "区", "街道"})

	// 创建Collector
	collector := colly.NewCollector(
		// 设置用户代理
		colly.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/84.0.4147.125 Safari/537.36"),
	)

	// 设置抓取频率限制
	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: time.Second, // 随机延迟
	})

	// 异常处理
	collector.OnError(func(response *colly.Response, err error) {
		log.Println(err.Error())
	})

	collector.OnRequest(func(request *colly.Request) {
		log.Println("start visit: ", request.URL.String())
	})

	// 解析列表
	collector.OnHTML(".homelist", func(element *colly.HTMLElement) {
		// 依次遍历所有的li节点
		element.DOM.Find("a").Each(func(i int, selection *goquery.Selection) {
			title := selection.Text()
			href, found := selection.Attr("href")
			writer.Write([]string{title})
			// 如果找到了详情页，则继续下一步的处理
			if found {
				parseDetail(collector, href, writer)
				log.Println(href)
			}
		})
	})

	// 起始入口
	collector.Visit(startUrl)
}

/**
* 处理市数据
 */
func cityDetail(collector *colly.Collector, url string, writer *csv.Writer) {
	collector = collector.Clone()
	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: time.Second,
	})
	collector.OnRequest(func(request *colly.Request) {
		log.Println("start visit: ", request.URL.String())
	})
	// 解析详情页数据
	collector.OnHTML("table#cont", func(element *colly.HTMLElement) {
		// 依次遍历所有的li节点
		element.DOM.Find("tr > td:first").Each(func(i int, selection *goquery.Selection) {

		})
	})
}

/**
 * 处理详情页
 */
func parseDetail(collector *colly.Collector, url string, writer *csv.Writer) {
	collector = collector.Clone()

	collector.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: time.Second,
	})

	collector.OnRequest(func(request *colly.Request) {
		log.Println("start visit: ", request.URL.String())
	})

	// 解析详情页数据
	collector.OnHTML("body", func(element *colly.HTMLElement) {
		selection := element.DOM.Find("div#content")
		idx := selection.Find("div.top250 > span.top250-no").Text()
		title := selection.Find("h1 > span").First().Text()
		year := selection.Find("h1 > span.year").Text()
		info := selection.Find("div#info").Text()
		info = strings.ReplaceAll(info, " ", "")
		info = strings.ReplaceAll(info, "\n", "; ")
		rating := selection.Find("strong.rating_num").Text()
		movie := Movie{
			idx:    idx,
			title:  title,
			year:   year,
			info:   info,
			rating: rating,
			url:    element.Request.URL.String(),
		}
		writer.Write([]string{
			idx,
			title,
			year,
			info,
			rating,
			element.Request.URL.String(),
		})
		log.Printf("%+v", movie)
	})

	collector.Visit(startUrl + url)
}
