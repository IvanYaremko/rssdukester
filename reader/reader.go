package reader

import (
	"time"

	"github.com/go-shiori/go-readability"
)

func GetArticle(url, title string) (string, error) {
	article, err := readability.FromURL(url, 10*time.Second)
	if err != nil {
		return "", err
	}

	return article.Content, nil
}
