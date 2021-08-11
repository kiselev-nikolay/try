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
	"log"

	"github.com/kiselev-nikolay/try"
)

var v interface{}

func processJSON(tc try.TryContext) {
	err := json.Unmarshal([]byte(""), &v)
	tc.Catch(err)
}

func main() {
	ctx := context.Background()
	err := try.Try(ctx, processJSON)
	if err != nil {
		log.Fatal(err)
	}
	log.Fatal(v)
}

```

