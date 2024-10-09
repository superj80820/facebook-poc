package parser

import (
	"errors"
	"io"
	"net/http"
	"regexp"
	"strings"

	// TODO: for york
	"github.com/superj80820/facebook-poc/domain"

	"github.com/PuerkitoBio/goquery"
)

var (
	entryPointRegex = regexp.MustCompile(`"entryPoint":{"__dr":"(.*?)"}}`)

	identifierPattern1Regex = regexp.MustCompile(`"identifier":"{0,1}([0-9]{5,})"{0,1},`)
	identifierPattern2Regex = regexp.MustCompile(`fb://profile/(.*?)"`)
	identifierPattern3Regex = regexp.MustCompile(`content="fb://group/([0-9]{1,})" />`)
	identifierPattern4Regex = regexp.MustCompile(`"pageID":"{0,1}([0-9]{5,})"{0,1},`)

	docIDRegex = regexp.MustCompile(`e.exports="([0-9]{1,})"`)
)

type parser struct{}

func NewParser() domain.ParserInf {
	return &parser{}
}

func (p *parser) ParsePageQuery(homePage string) (*domain.PageQuery, error) {
	var query domain.PageQuery

	// get entry point
	matches := entryPointRegex.FindStringSubmatch(homePage)
	if len(matches) > 1 && strings.Contains(matches[1], `ProfilePlusCometLoggedOutRouteRoot.entrypoint`) {
		query.EntryPoint = domain.ProfilePlusCometLoggedOutRouteRoot
	} else if len(matches) > 1 && strings.Contains(matches[1], `CometGroupDiscussionRoot.entrypoint`) {
		query.EntryPoint = domain.CometGroupDiscussionRoot
	} else if len(matches) > 1 && strings.Contains(matches[1], `CometSinglePageHomeRoot.entrypoint`) {
		query.EntryPoint = domain.CometSinglePageHomeRoot
	} else {
		query.EntryPoint = domain.NoJSEntryPointType
	}

	// get identifier
	switch query.EntryPoint {
	case domain.ProfilePlusCometLoggedOutRouteRoot, domain.CometGroupDiscussionRoot:
		if matches := identifierPattern1Regex.FindStringSubmatch(homePage); len(matches) > 1 {
			query.Identifier = matches[1]
		} else if matches := identifierPattern2Regex.FindStringSubmatch(homePage); len(matches) > 1 {
			query.Identifier = matches[1]
		} else if matches := identifierPattern3Regex.FindStringSubmatch(homePage); len(matches) > 1 {
			query.Identifier = matches[1]
		}
	case domain.CometSinglePageHomeRoot, domain.NoJSEntryPointType:
		if matches := identifierPattern4Regex.FindStringSubmatch(homePage); len(matches) > 1 {
			query.Identifier = matches[1]
		}
	}

	// get doc id
	switch query.EntryPoint {
	case domain.NoJSEntryPointType:
		query.DocID = "NoDocid"
	default:
		doc, err := goquery.NewDocumentFromReader(strings.NewReader(homePage))
		if err != nil {
			return nil, err
		}

		// TODO: for york
		var hrefs []string
		doc.Find("link[rel='preload']").Each(func(i int, s *goquery.Selection) {
			href, exists := s.Attr("href")
			if !exists {
				return
			}

			hrefs = append(hrefs, href)
		})

		for _, href := range hrefs {
			resp, err := http.Get(href)
			if err != nil {
				return nil, err
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				return nil, errors.New("failed to fetch docid")
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return nil, err
			}
			lines := strings.Split(string(body), "\n")

			for _, line := range lines {
				if strings.Contains(line, "ProfileCometTimelineFeedRefetchQuery_") ||
					strings.Contains(line, "CometModernPageFeedPaginationQuery_") ||
					strings.Contains(line, "CometUFICommentsProviderQuery_") ||
					strings.Contains(line, "GroupsCometFeedRegularStoriesPaginationQuery") {
					if matches := docIDRegex.FindStringSubmatch(line); len(matches) > 1 {
						query.DocID = matches[1]
						break
					}
				}
			}
		}
	}

	return &query, nil
}
