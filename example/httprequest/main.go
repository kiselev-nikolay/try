package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/kiselev-nikolay/try"
)

func doWork(tc try.TryContext) []byte {
	rawData := make(map[string]string)
	err := json.Unmarshal([]byte("{\"try\":catch}"), &rawData)
	tc.Catch(err)
	data := url.Values{}
	for k, v := range rawData {
		data.Set(k, v)
	}
	resp, err := http.PostForm("http://example.com/form", data)
	tc.Catch(err)
	body, err := io.ReadAll(resp.Body)
	tc.Catch(err)
	return body
}

func doWork2(try.TryContext, []byte) []byte {
	// there also code with error handling as doWork
	panic("...")
}

func doWork3(try.TryContext, []byte) []byte {
	// there also code with error handling as doWork
	panic("...")
}

func useDoWorks() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return try.Try(ctx, func(tc try.TryContext) {
		b := doWork(tc)
		b = doWork2(tc, b)
		b = doWork3(tc, b)
	})
}
