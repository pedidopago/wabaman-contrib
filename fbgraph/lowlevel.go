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

// CompareGraphAPIVersions compares two Graph API versions.
// Returns -1 if v1 < v2, 0 if v1 == v2, and 1 if v1 > v2.
// Versions should be in the format "vX.Y" where X and Y are integers.
func CompareGraphAPIVersions(v1, v2 string) (int, error) {
	v1major, v1minor := splitGraphAPIVersion(v1)
	v2major, v2minor := splitGraphAPIVersion(v2)

	if v1major == -1 || v2major == -1 {
		return 0, fmt.Errorf("invalid version: %s", v1)
	}

	if v1major < v2major {
		return -1, nil
	}
	if v1major > v2major {
		return 1, nil
	}

	if v1minor < v2minor {
		return -1, nil
	}

	if v1minor > v2minor {
		return 1, nil
	}

	return 0, nil
}

func splitGraphAPIVersion(version string) (major, minor int) {
	if len(version) < 2 || version[0] != 'v' {
		return -1, -1
	}

	n, err := fmt.Sscanf(version[1:], "%d.%d", &major, &minor)
	if err != nil || n < 1 {
		return -1, -1
	}

	return major, minor
}
