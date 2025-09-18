package fbgraph

import (
	"context"
	"fmt"
	"io"
	"net/http"
)

func NewRequest(method string, url string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, url, body)

	if DebugTrace {
		if err != nil {
			fmt.Printf("NewRequest(%s, %s, %t): %s\n", method, url, body != nil, err)
		} else {
			fmt.Printf("NewRequest(%s, %s, %t)\n", method, url, body != nil)
		}
	}

	return request, err
}

func NewRequestWithContext(ctx context.Context, method string, url string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequestWithContext(ctx, method, url, body)

	if DebugTrace {
		if err != nil {
			fmt.Printf("NewRequestWithContext(%s, %s, %t): %s\n", method, url, body != nil, err)
		} else {
			fmt.Printf("NewRequestWithContext(%s, %s, %t)\n", method, url, body != nil)
		}
	}
	return request, err
}
