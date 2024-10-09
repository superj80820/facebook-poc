package domain

import "time"

type PageFetcherInf interface {
	GetHomePage() (string, error)
	GetPosts(pageQuery *PageQuery, cursor string) (string, error)
}

type FormatterInf interface {
	FormatPage(posts string, entryPointType EntryPoint) (*PageInfo, error)
}

type ParserInf interface {
	ParsePageQuery(homePage string) (*PageQuery, error)
}

type EntryPoint string

const (
	ProfilePlusCometLoggedOutRouteRoot EntryPoint = "ProfilePlusCometLoggedOutRouteRoot.entrypoint"
	CometGroupDiscussionRoot           EntryPoint = "CometGroupDiscussionRoot.entrypoint"
	CometSinglePageHomeRoot            EntryPoint = "CometSinglePageHomeRoot.entrypoint"
	NoJSEntryPointType                 EntryPoint = "NoJS"
)

type PageQuery struct {
	EntryPoint EntryPoint
	Identifier string
	DocID      string
}

type PageInfo struct {
	Name        string
	Message     string
	PostURL     string
	Cursor      string
	CreatedTime time.Time
}

type CrawlerInf interface {
	FetchPagePosts(startData, endDate time.Time) ([]*PageInfo, error)
}
