package markdown

import htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"

func ConvertArticle(article string) (string, error) {
	markdown, err := htmltomarkdown.ConvertString(article)
	if err != nil {
		return "", err
	}

	return markdown, nil
}
