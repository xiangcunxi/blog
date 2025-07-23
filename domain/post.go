package domain

import "time"

type Post struct {
	Id      int64
	Title   string
	Content string
	Author  Author
	Ctime   time.Time
	Utime   time.Time
}

type Author struct {
	ID   int64
	Name string
}

type ArticleVO struct {
	Id       int64  `json:"id"`
	Title    string `json:"title"`
	Abstract string `json:"abstract"`
	Content  string `json:"content"`
	Author   string `json:"author"`
	Ctime    string `json:"ctime"`
	Utime    string `json:"utime"`
}
