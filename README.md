# Try

<img title="The image of a baseball gopher catcher" align="right" width="260px" src="https://user-images.githubusercontent.com/55307887/129087852-317182cf-ef93-4fe6-a856-2d3d37952b8f.png">

Go package that enhances Go's error handling capabilities in microservices with complex partitioning into multiple layers.

__experiment; do not use it in production__

_Illustrated on right side image is a baseball gopher catcher_

#### Example

```go
package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/kiselev-nikolay/try"
)

var body []byte

func doWork(tc try.TryContext) {
	rawData := make(map[string]string)
	err := json.Unmarshal([]byte("{\"try\":catch}"), &rawData)
	tc.Catch(err)
	data := url.Values{}
	for k, v := range rawData {
		data.Set(k, v)
	}
	resp, err := http.PostForm("http://example.com/form", data)
	tc.Catch(err)
	body, err = io.ReadAll(resp.Body)
	tc.Catch(err)
	log.Println(body)
}

func main() {
	ctx := context.Background()
	err := try.Try(ctx, doWork)
	if err != nil {
		log.Fatal("sad, doWork failed: ", err)
	}
	log.Println(body)
}

```

## Why is that a thing

Because developers in programmes with complex partitioning into multiple layers must return any errors that occur up through the layers. In an ideal world and serious libraries, this looks different and is justified. But in simple projects it looks like this:

```go
func doWork() ([]byte, error) {
	rawData := make(map[string]string)
	err := json.Unmarshal([]byte("{\"try\":catch}"), &rawData)
	if err != nil {
		return nil, err
	}
	data := url.Values{}
	for k, v := range rawData {
		data.Set(k, v)
	}
	resp, err := http.PostForm("http://example.com/form", data)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func doWork2(try.TryContext, []byte) []byte {
	// there also code with error handling as doWork
	panic("...")
}

func doWork3(try.TryContext, []byte) []byte {
	// there also code with error handling as doWork
	panic("...")
}

func UseDoWorks() error {
	b, err := doWork()
	if err != nil {
		return err
	}
	b, err = doWork2(b)
	if err != nil {
		return err
	}
	b, err = doWork3(b)
	if err != nil {
		return err
	}
}
```

This __try__ package helps with that! So code will be:

```go
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

func UseDoWorks() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	return try.Try(ctx, func(tc try.TryContext) {
		b := doWork(tc)
		b = doWork2(tc, b)
		b = doWork3(tc, b)
	})
}
```
