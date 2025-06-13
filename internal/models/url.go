package models

import "time"

type Url struct {
	Id           int64      `json:"id"`
	Original_url string     `json:"url"`
	Alias        string     `json:"alias"`
	Created_at   time.Time  `json:"created_at"`
	Expires_at   *time.Time `json:"expires_at"`
	Is_active    bool       `json:"is_active"`
	Visit_count  int        `json:"visit_count"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
}

type Visit struct {
	Id         int64     `json:"id"`
	Url_id     int64     `json:"url_id"`
	Created_at time.Time `json:"created_at"`
}

type UrlInfo struct {
	Id           int64      `json:"id"`
	Original_url string     `json:"url"`
	Alias        string     `json:"alias"`
	Created_at   time.Time  `json:"created_at"`
	Expires_at   *time.Time `json:"expires_at"`
	Is_active    bool       `json:"is_active"`
	Visit_count  int        `json:"visit_count"`
	Title        string     `json:"title"`
	Description  string     `json:"description"`
	Visits       []Visit    `json:"visits"`
}

type ShortenRequest struct {
	URL         string
	Alias       string
	MaxVisits   *int
	ExpiresAt   *time.Time
	Title       *string
	Description *string
	Step        string
	SkipClicked bool
}
