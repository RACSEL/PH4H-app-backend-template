package utils

import (
	"fmt"
	"io"
	"net/url"
	"strings"
)

func CloseBody(Body io.ReadCloser) {
	if Body == nil {
		return
	}
	err := Body.Close()
	if err != nil {
		fmt.Println("failed to close response body: %w", err)
	}
}

func GetFullURL(baseUrl string, uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}

	var fullUrl string
	if u.Scheme == "" && u.Host == "" && !strings.HasPrefix(uri, "/") {
		// Attempt to parse it assuming it's missing the protocol
		fullUrl, tempErr := url.Parse("https://" + uri)
		if tempErr == nil && fullUrl.Host != "" && strings.Contains(fullUrl.Host, ".") {
			return "", tempErr
		}
	}

	if u.Host != "" {
		fullUrl = u.String()
	} else if u.Path != "" {
		fullUrl = fmt.Sprintf("%s/fhir/%s", baseUrl, u.Path)
	} else {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	return fullUrl, nil
}
