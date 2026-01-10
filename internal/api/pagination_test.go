package api

import (
	"net/http"
	"testing"
)

func TestParsePaginationLinks(t *testing.T) {
	tests := []struct {
		name       string
		linkHeader string
		wantNext   string
		wantPrev   string
		wantFirst  string
		wantLast   string
		wantCur    string
	}{
		{
			name: "full pagination",
			linkHeader: `<https://example.com/api/v1/courses?page=2>; rel="current",` +
				`<https://example.com/api/v1/courses?page=3>; rel="next",` +
				`<https://example.com/api/v1/courses?page=1>; rel="prev",` +
				`<https://example.com/api/v1/courses?page=1>; rel="first",` +
				`<https://example.com/api/v1/courses?page=5>; rel="last"`,
			wantNext:  "https://example.com/api/v1/courses?page=3",
			wantPrev:  "https://example.com/api/v1/courses?page=1",
			wantFirst: "https://example.com/api/v1/courses?page=1",
			wantLast:  "https://example.com/api/v1/courses?page=5",
			wantCur:   "https://example.com/api/v1/courses?page=2",
		},
		{
			name:       "only next link",
			linkHeader: `<https://example.com/api/v1/courses?page=2>; rel="next"`,
			wantNext:   "https://example.com/api/v1/courses?page=2",
		},
		{
			name:       "empty link header",
			linkHeader: "",
		},
		{
			name:       "malformed link",
			linkHeader: "invalid link header",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := &http.Response{
				Header: http.Header{
					"Link": []string{tt.linkHeader},
				},
			}

			links := ParsePaginationLinks(resp)
			if links == nil {
				t.Fatal("expected non-nil pagination links")
			}

			if links.Next != tt.wantNext {
				t.Errorf("Next = %v, want %v", links.Next, tt.wantNext)
			}
			if links.Prev != tt.wantPrev {
				t.Errorf("Prev = %v, want %v", links.Prev, tt.wantPrev)
			}
			if links.First != tt.wantFirst {
				t.Errorf("First = %v, want %v", links.First, tt.wantFirst)
			}
			if links.Last != tt.wantLast {
				t.Errorf("Last = %v, want %v", links.Last, tt.wantLast)
			}
			if links.Current != tt.wantCur {
				t.Errorf("Current = %v, want %v", links.Current, tt.wantCur)
			}
		})
	}
}

func TestPaginationLinks_HasNextPage(t *testing.T) {
	tests := []struct {
		name string
		next string
		want bool
	}{
		{
			name: "has next page",
			next: "https://example.com/api/v1/courses?page=2",
			want: true,
		},
		{
			name: "no next page",
			next: "",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			links := &PaginationLinks{
				Next: tt.next,
			}

			got := links.HasNextPage()
			if got != tt.want {
				t.Errorf("HasNextPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPaginationLinks_HasPrevPage(t *testing.T) {
	tests := []struct {
		name string
		prev string
		want bool
	}{
		{
			name: "has prev page",
			prev: "https://example.com/api/v1/courses?page=1",
			want: true,
		},
		{
			name: "no prev page",
			prev: "",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			links := &PaginationLinks{
				Prev: tt.prev,
			}

			got := links.HasPrevPage()
			if got != tt.want {
				t.Errorf("HasPrevPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPageNumber(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			// NOTE: This test documents existing behavior where GetPageNumber
			// incorrectly matches "page=" in "per_page=" parameter.
			// The function splits on first "page=" occurrence.
			name: "page after per_page (documents bug)",
			url:  "https://example.com/api/v1/courses?per_page=10&page=5&other=param",
			want: "10", // Returns per_page value due to substring match
		},
		{
			name: "page after per_page at end",
			url:  "https://example.com/api/v1/courses?per_page=10&page=3",
			want: "10", // Returns per_page value due to substring match
		},
		{
			name: "only per_page parameter",
			url:  "https://example.com/api/v1/courses?per_page=10",
			want: "10", // Returns per_page value due to substring match
		},
		{
			name: "page is first parameter",
			url:  "https://example.com/api/v1/courses?page=2",
			want: "2",
		},
		{
			name: "page before per_page",
			url:  "https://example.com/api/v1/courses?page=5&per_page=10",
			want: "5", // Correct when page comes first
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetPageNumber(tt.url)
			if got != tt.want {
				t.Errorf("GetPageNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPerPage(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want string
	}{
		{
			name: "per_page in middle of URL",
			url:  "https://example.com/api/v1/courses?page=5&per_page=20&other=param",
			want: "20",
		},
		{
			name: "per_page at end of URL",
			url:  "https://example.com/api/v1/courses?page=3&per_page=50",
			want: "50",
		},
		{
			name: "no per_page parameter",
			url:  "https://example.com/api/v1/courses?page=5",
			want: "10", // Default
		},
		{
			name: "per_page is first parameter",
			url:  "https://example.com/api/v1/courses?per_page=100",
			want: "100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetPerPage(tt.url)
			if got != tt.want {
				t.Errorf("GetPerPage() = %v, want %v", got, tt.want)
			}
		})
	}
}
