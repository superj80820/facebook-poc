package formatter

import (
	// TODO: for york
	"fmt"
	"time"

	"github.com/superj80820/facebook-poc/domain"

	"github.com/bitly/go-simplejson"
)

type formatter struct{}

func NewFormatter() domain.FormatterInf {
	return &formatter{}
}

func (f *formatter) FormatPage(postsStr string, entryPointType domain.EntryPoint) (*domain.PageInfo, error) {
	postsJSON, err := simplejson.NewJson([]byte(postsStr))
	if err != nil {
		return nil, fmt.Errorf("failed to parse posts: %w", err)
	}

	edge := postsJSON.Get("data").Get("node").Get("timeline_list_feed_units").Get("edges").GetIndex(0)
	cometSections := edge.Get("node").Get("comet_sections")

	name := cometSections.Get("context_layout").Get("story").Get("comet_sections").Get("actor_photo").Get("story").Get("actors").GetIndex(0).Get("name").MustString()
	createTime := cometSections.Get("context_layout").Get("story").Get("comet_sections").Get("metadata").GetIndex(0).Get("story").Get("creation_time").MustInt64()
	message := cometSections.Get("content").Get("story").Get("comet_sections").Get("message").Get("story").Get("message").Get("text").MustString()
	postURL := cometSections.Get("content").Get("story").Get("wwwURL").MustString()
	cursor := edge.Get("cursor").MustString()

	return &domain.PageInfo{
		Name:        name,
		CreatedTime: time.Unix(createTime, 0),
		Message:     message,
		PostURL:     postURL,
		Cursor:      cursor,
	}, nil
}
