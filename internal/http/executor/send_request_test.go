package executor

import (
	"net/http"
	"io"
	"net/url"
	"strings"
	"testing"

	//"io"
	"fmt"

	"github.com/alp1n3-eth/cast/pkg/models"
)

func BenchmarkSendFastHTTPRequest(b *testing.B) {
	parsedURL, err := url.Parse("https://www.google.com")
	if err != nil {
    	fmt.Println("Error parsing URL:", err)
    return
	}

	parsedURL2, err := url.Parse("https://echo.free.beeceptor.com")
	if err != nil {
    	fmt.Println("Error parsing URL:", err)
    return
	}

	var body io.Reader
	body = strings.NewReader("test1=test2")


	headers := http.Header{}
	headers.Add("header1", "value2")

	tests := []models.ExecutionResult{
	{
		Request: models.Request{
			Method: "GET",
			URL: parsedURL,
		},
	},
		{
			Request: models.Request{
			Method: "POST",
			URL: parsedURL2,
			Body: body,
			Headers: headers,
		},
	},

	}

	for _, tt := range tests {
        b.Run(tt.Request.Method, func(b *testing.B) {
            b.ResetTimer()
            for i := 0; i < b.N; i++ {
                SendFastHTTPRequest(tt)
            }
        })
    }
}
