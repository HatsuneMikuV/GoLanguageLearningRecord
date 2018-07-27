package github

import "time"

const (
	Ltm  = "Lessthanmonth"
 	Lty  = "Lessthanyear"
 	Mty  = "Morethanyear"
	IssuesURL = "https://api.github.com/search/issues"
)

type IssuesSearchResult struct {
	TotalCount int `json:"total_count"`
	Items          []*Issue
	Dict       map[string][]*Issue
}

type Issue struct {
	Number    int
	HTMLURL   string `json:"html_url"`
	Title     string
	State     string
	User      *User
	CreatedAt time.Time `json:"created_at"`
	Body      string    // in Markdown format
}

type User struct {
	Login   string
	HTMLURL string `json:"html_url"`
}