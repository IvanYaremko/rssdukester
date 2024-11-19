package reader

import (
	"net/http"
	"net/url"

	htmltomarkdown "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/go-shiori/go-readability"
)

func convertArticle(article string) (string, error) {
	markdown, err := htmltomarkdown.ConvertString(article)
	if err != nil {
		return "", err
	}

	return markdown, nil
}

func GetMarkdown(articleUrl string) (string, error) {
	req, err := http.NewRequest("GET", articleUrl, nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("User-Agent", "rssdukester")
	//	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")
	req.Header.Add("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Add("Accept-Language", "en-US,en;q=0.5")

	client := http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	structedUrl, err := url.Parse(articleUrl)
	if err != nil {
		return "", err
	}
	article, err := readability.FromReader(response.Body, structedUrl)
	if err != nil {
		return "", nil
	}

	markdown, err := convertArticle(article.Content)
	if err != nil {
		return "", err
	}
	return markdown, nil
}
