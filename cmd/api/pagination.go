package main

import (
	"github.com/sergdort/Social/business/domain"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func ParsePaginatedFeedQuery(fq *domain.PaginatedFeedQuery, r *http.Request) {
	qs := r.URL.Query()

	if limit, err := strconv.Atoi(qs.Get("limit")); err == nil {
		if limit > 0 && limit < 50 {
			fq.Limit = limit
		}
	}
	if offset, err := strconv.Atoi(qs.Get("offset")); err == nil {
		if offset > 0 {
			fq.Offset = offset
		}
	}
	sortBy := qs.Get("sort_by")
	if sortBy == "ask" || sortBy == "desc" {
		fq.SortBy = sortBy
	}

	tags := qs.Get("tags")
	if tags != "" {
		fq.Tags = strings.Split(tags, ",")
	}

	search := qs.Get("search")
	if search != "" {
		fq.Search = search
	}

	since := qs.Get("since")
	if since != "" {
		fq.Since = parseTime(since)
	}

	until := qs.Get("until")
	if until != "" {
		fq.Until = parseTime(until)
	}
}

func parseTime(since string) string {
	t, err := time.Parse(time.DateTime, since)
	if err != nil {
		return ""
	}
	return t.Format(time.DateTime)
}
