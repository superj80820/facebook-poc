package crawler

import (
	"fmt"
	"strings"
	"time"

	"github.com/superj80820/facebook-poc/domain"
)

type Crawler struct {
	fetcherInf   domain.PageFetcherInf
	formatterInf domain.FormatterInf
	parserInf    domain.ParserInf
}

func NewCrawler(fetcherInf domain.PageFetcherInf, formatterInf domain.FormatterInf, parserInf domain.ParserInf) domain.CrawlerInf {
	return &Crawler{
		fetcherInf:   fetcherInf,
		formatterInf: formatterInf,
		parserInf:    parserInf,
	}
}

// TODO: endDate not working now
func (c *Crawler) FetchPagePosts(startDate, endDate time.Time) ([]*domain.PageInfo, error) {
	homePage, err := c.fetcherInf.GetHomePage()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch home page: %w", err)
	}

	pageQuery, err := c.parserInf.ParsePageQuery(homePage)
	if err != nil {
		return nil, fmt.Errorf("failed to parse page query: %w", err)
	}

	fmt.Println("doc id:", pageQuery.DocID, "identifier:", pageQuery.Identifier)

	endDate = time.Now()
	var (
		pageInfos []*domain.PageInfo
		cursor    string
	)
	for {
		pagePosts, err := c.fetcherInf.GetPosts(pageQuery, cursor)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch page posts: %w", err)
		}

		pageInfo, err := c.formatterInf.FormatPage(pagePosts, pageQuery.EntryPoint)
		if err != nil {
			return nil, fmt.Errorf("failed to format page: %w", err)
		}

		if startDate.After(pageInfo.CreatedTime) {
			break
		}

		printPageInfos(pageInfo)

		pageInfos = append(pageInfos, pageInfo)

		cursor = pageInfo.Cursor
	}

	return pageInfos, nil
}

func printPageInfos(pageInfo *domain.PageInfo) {
	message := pageInfo.Message
	if len(message) >= 30 {
		message = message[:30] + "..."
		message = strings.ReplaceAll(message, "\n", "")
	}

	fmt.Println("name:", pageInfo.Name, "message:", message, "created time:", pageInfo.CreatedTime)
}
