package reader

import (
	"time"

	"github.com/IvanYaremko/rssdukester/markdown"
	"github.com/go-shiori/go-readability"
)

func GetMarkdown(url, title string) (string, error) {
	article, err := readability.FromURL(url, 10*time.Second)
	if err != nil {
		return "", err
	}

	markdown, err := markdown.ConvertArticle(article.Content)
	if err != nil {
		return "", err
	}
	return markdown, nil
}
