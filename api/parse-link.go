package api

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

var regLinkHeader = regexp.MustCompile(`(?m)\<(?P<uri>[^\>]+)\>; rel="(?P<rel>[^"]+)"`)

func parseLinkHeader(linkHeader string) (*map[string]*url.URL, error) {
	matches := regLinkHeader.FindAllStringSubmatch(linkHeader, -1)

	links := make(map[string]*url.URL, len(matches))

	for _, group := range matches {
		var (
			linkUrl *url.URL
			rel     *string
		)

		for groupId, groupName := range regLinkHeader.SubexpNames() {
			if groupName == "uri" {
				uri := group[groupId]

				u, err := url.Parse(uri)
				if err != nil {
					return nil, fmt.Errorf("parse link `%s`: %w", uri, err)
				}

				linkUrl = u
			}

			if groupName == "rel" {
				rel = &group[groupId]
			}
		}

		if rel == nil || linkUrl == nil {
			return nil, fmt.Errorf("group link `%s` not parsed URI or REL", group[0])
		}

		links[*rel] = linkUrl
	}

	return &links, nil
}

func getQueryPage(u *url.URL) (int, error) {
	ps := u.Query().Get("page")

	if ps == "" {
		return 0, nil
	}

	return strconv.Atoi(ps)
}

func getAllPagesFromHeader(headers *http.Header) (int, error) {
	headerLink := headers.Get("Link")

	links, err := parseLinkHeader(headerLink)
	if err != nil {
		return 0, fmt.Errorf("parse link header: %w", err)
	}

	lastUrl, ok := (*links)["last"]
	if !ok {
		return 0, nil
	}

	page, err := getQueryPage(lastUrl)
	if err != nil {
		return 0, fmt.Errorf("get query page from query: %w", err)
	}

	return page, nil
}
