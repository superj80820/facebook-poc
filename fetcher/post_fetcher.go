package fetcher

import (
	"io"
	"net/http"
	"strings"

	"github.com/superj80820/facebook-poc/domain"
)

type postFetcher struct {
	url string
}

func NewFetcher(url string) domain.PageFetcherInf {
	return &postFetcher{
		url: url,
	}
}

func (f *postFetcher) GetHomePage() (string, error) {
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, f.url, nil)
	if err != nil {
		return "", err
	}

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// TODO: cursor arg to option
func (f *postFetcher) GetPosts(pageQuery *domain.PageQuery, cursor string) (string, error) {
	url := "https://www.facebook.com/api/graphql/"
	method := "POST"

	variables := []string{
		`"count":3`,
		`"id":"` + pageQuery.Identifier + `"`,
	}
	if cursor != "" {
		variables = append(variables, `"cursor":"`+cursor+`"`)
	}

	payload := strings.NewReader(`variables={` + strings.Join(variables, ",") + `}&doc_id=` + pageQuery.DocID)

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
