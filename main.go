package main

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"

	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/debug"
)

var log = logrus.New()

func main() {

	log.Out = os.Stdout
	log.SetFormatter(&logrus.TextFormatter{ForceColors: true})
	log.SetOutput(colorable.NewColorableStdout())
	file, err := os.OpenFile("logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("Failed to log to file, using default stderr")
	}

	re, err := regexp.Compile(`https://arxiv\.org/[a-zA-Z]{3}/\d{4}.\d{5}`)
	if err != nil {
		log.Panic(err)
	}

	re_check_for_arxiv, err := regexp.Compile(`https://arxiv\.org/`)
	if err != nil {
		log.Panic(err)
	}

	c := colly.NewCollector(colly.Debugger(&debug.LogDebugger{}), colly.Async())

	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: 2, RandomDelay: 10 * time.Second})

	c.OnHTML("div#content a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		link = e.Request.AbsoluteURL(link)
		rez := re.MatchString(link)
		if rez {
			fmt.Println(e.Request.AbsoluteURL(link))
			date := e.ChildText("div.dateline")
			title := e.ChildText("h1.title")
			descriptor := e.ChildText("span.descriptor")
			primary_subject := e.ChildText("span.primary-subject")
			act_link := link
			DbInsert(date, title, descriptor, primary_subject, act_link)
		} else {
			check_base := re_check_for_arxiv.MatchString(link)
			if check_base {
				e.Request.Visit(e.Request.AbsoluteURL(link))
			}
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		log.Println("Request URL:", r.Request.URL, "failed with response:", r.StatusCode, "Error:", err)
	})

	c.Visit("https://arxiv.org/")
	c.Wait()
}
