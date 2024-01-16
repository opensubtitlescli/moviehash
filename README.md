# Moviehash

[![audit](https://github.com/opensubtitlescli/moviehash/actions/workflows/audit.yml/badge.svg)](https://github.com/opensubtitlescli/moviehash/actions/workflows/audit.yml)

A package of one hash function to match subtitle files against movie files. You can read more about algorithm on the [OpenSubtitles API](https://opensubtitles.stoplight.io/docs/opensubtitles-api/e3750fd63a100-getting-started#calculating-moviehash-of-video-file).

```go
package main

import (
	"fmt"
	"os"
	"github.com/opensubtitlescli/moviehash"
)

func main() {
	f, _ := os.Open("./breakdance.avi")
	h, _ := moviehash.Sum(f)
	fmt.Println(h)
}
```

Result:

```txt
8e245d9679d31e12
```

## License

The repository code is distributed under the [MIT License](./LICENSE).
