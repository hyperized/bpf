# bpf/files

Utility library for bpf files.

Useful for BSD: Darwin, Dragonfly, FreeBSD, NetBSD & OpenBSD

## Files

Example implementation:

```go
package main

import (
	"github.com/hyperized/bpf/file"
	"log"
)

func main() {
	// Use logging, off by default
	bpf, err := file.GetBpfDevice(true)
	defer func() {
		if err := bpf.Close(); err != nil {
			log.Printf(err.Error())
		}
	}()

	// Showing off methods available
	if bpf != nil {
		log.Println(bpf.File())
		log.Println(bpf.Path())
		log.Println(bpf.FileDescriptor())
	} else {
		log.Println(err)
	}

	// Open new bpf device, this will be a new one, since the previous one is already open
	bpf2, err := file.GetBpfDevice(true)
	defer func() {
		if err := bpf2.Close(); err != nil {
			log.Printf(err.Error())
		}
	}()

	log.Println(bpf2.Path())
}

```

Results in:

```
2019/10/09 14:35:40 github.com/hyperized/bpf/files: /dev/bpf0 cannot be accessed, skipping
2019/10/09 14:35:40 github.com/hyperized/bpf/files: /dev/bpf1 cannot be accessed, skipping
2019/10/09 14:35:40 github.com/hyperized/bpf/files: using /dev/bpf10
2019/10/09 14:35:40 &{0xc0000221e0}
2019/10/09 14:35:40 /dev/bpf10
2019/10/09 14:35:40 3
2019/10/09 14:35:40 github.com/hyperized/bpf/files: /dev/bpf0 cannot be accessed, skipping
2019/10/09 14:35:40 github.com/hyperized/bpf/files: /dev/bpf1 cannot be accessed, skipping
2019/10/09 14:35:40 github.com/hyperized/bpf/files: /dev/bpf10 cannot be accessed, skipping
2019/10/09 14:35:40 github.com/hyperized/bpf/files: using /dev/bpf100
2019/10/09 14:35:40 /dev/bpf100
2019/10/09 14:35:40 github.com/hyperized/bpf/files: closing /dev/bpf100
2019/10/09 14:35:40 github.com/hyperized/bpf/files: closing /dev/bpf10
```

## Author

Gerben Geijteman <gerben@hyperized.net>