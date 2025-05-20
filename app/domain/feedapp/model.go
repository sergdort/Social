package feedapp

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/sergdort/Social/business/domain"
)

type PostFeedItem struct {
	ID            int64    `json:"id" example:"117"`
	Content       string   `json:"content" example:"I will not become a queen of ashes."`
	Title         string   `json:"title" example:"The King of Ashes"`
	UserID        int      `json:"user_id" example:"38"`
	CreatedAt     string   `json:"created_at" example:"2025-03-19 10:08:25 +0000 UTC"`
	UpdatedAt     string   `json:"updated_at" example:"2025-03-19 10:08:25 +0000 UTC"`
	Tags          []string `json:"tags" example:"Dothraki,Lannister,BattleOfBastards,KingsLanding"`
	CommentsCount int64    `json:"comments_count" example:"4"`
	User          FeedUser `json:"user"`
}

// Needed for swagger docs, should not be used
type FeedData struct {
	Data []PostFeedItem `json:"data"`
}

type FeedUser struct {
	ID       int    `json:"id" example:"38"`
	Username string `json:"username" example:"GendryBaratheon"`
}

func parsePaginatedFeedQuery(fq *domain.PaginatedFeedQuery, r *http.Request) {
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

func toPostFeedItem(p domain.PostWithMetadata) PostFeedItem {
	return PostFeedItem{
		ID:            p.ID,
		Content:       p.Content,
		Title:         p.Title,
		UserID:        int(p.UserID),
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
		Tags:          p.Tags,
		CommentsCount: p.CommentsCount,
		User:          toFeedUser(p.User),
	}
}

func toFeedUser(u domain.User) FeedUser {
	return FeedUser{
		ID:       int(u.ID),
		Username: u.Username,
	}
}
