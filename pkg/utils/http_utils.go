package utils

import (
	"fmt"
	"io"
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
