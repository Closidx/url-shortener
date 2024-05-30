package random_test

import (
	"github/closidx/url-shortener/internal/lib/random"
	"regexp"
	"testing"
)

func TestNewRandomString(t *testing.T) {
    testCases := []struct {
        name     string
        size     int
        expected string
    }{
        {
            name:     "generate string of length 10",
            size:     10,
            expected: "[A-Za-z0-9]{10}",
        },
        {
            name:     "generate string of length 20",
            size:     20,
            expected: "[A-Za-z0-9]{20}",
        },
        {
            name:     "generate string of length 0",
            size:     0,
            expected: "",
        },
        {
            name:     "generate string of length -1",
            size:     -1,
            expected: "",
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            result := random.NewRandomString(tc.size)
            if tc.size <= 0 {
                if len(result) > 0 {
                    t.Errorf("expected length %d, got %d", tc.size, len(result))
                }
                return
            }
            if len(result) != tc.size {
                t.Errorf("expected length %d, got %d", tc.size, len(result))
            }
            if !regexp.MustCompile(tc.expected).MatchString(result) {
                t.Errorf("expected result to match pattern %s, got %s", tc.expected, result)
            }
        })
    }
}
