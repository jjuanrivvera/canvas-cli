package api

import (
	"net/http"
	"regexp"
	"strings"
)

var linkRegex = regexp.MustCompile(`<([^>]+)>;\s*rel="([^"]+)"`)

// ParsePaginationLinks parses the Link header for pagination
func ParsePaginationLinks(resp *http.Response) *PaginationLinks {
	linkHeader := resp.Header.Get("Link")
	if linkHeader == "" {
		return &PaginationLinks{}
	}

	links := &PaginationLinks{}

	// Parse each link from the header
	matches := linkRegex.FindAllStringSubmatch(linkHeader, -1)
	for _, match := range matches {
		if len(match) != 3 {
			continue
		}

		url := match[1]
		rel := match[2]

		switch rel {
		case "current":
			links.Current = url
		case "next":
			links.Next = url
		case "prev":
			links.Prev = url
		case "first":
			links.First = url
		case "last":
			links.Last = url
		}
	}

	return links
}

// HasNextPage checks if there is a next page
func (p *PaginationLinks) HasNextPage() bool {
	return p.Next != ""
}

// HasPrevPage checks if there is a previous page
func (p *PaginationLinks) HasPrevPage() bool {
	return p.Prev != ""
}

// GetPageNumber extracts the page number from a pagination URL
func GetPageNumber(url string) string {
	// Extract page parameter from URL
	parts := strings.Split(url, "page=")
	if len(parts) < 2 {
		return ""
	}

	// Get the page number (everything after page= until & or end)
	pageStr := strings.Split(parts[1], "&")[0]
	return pageStr
}

// GetPerPage extracts the per_page parameter from a URL
func GetPerPage(url string) string {
	parts := strings.Split(url, "per_page=")
	if len(parts) < 2 {
		return "10" // Default
	}

	perPageStr := strings.Split(parts[1], "&")[0]
	return perPageStr
}
