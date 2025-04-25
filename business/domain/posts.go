package domain

type Post struct {
	ID        int64     `json:"id"`
	Content   string    `json:"content"`
	Title     string    `json:"title"`
	UserID    int64     `json:"user_id"`
	CreatedAt string    `json:"created_at"`
	UpdatedAt string    `json:"updated_at"`
	Tags      []string  `json:"tags"`
	Comments  []Comment `json:"comments"`
	Version   int64     `json:"version"`
	User      User      `json:"user"`
}

type PostWithMetadata struct {
	Post
	CommentsCount int64 `json:"comments_count"`
}
