# README

Super simple golang disk buffer.  Doesn't feature wrap-around yet.  Pretty sure it's thread safe.  

```golang
	import (
	"github.com/mvult/diskBuffer"
	"os"
	)

	f, err := os.Create("buffer.b")
	if err != nil {
		panic(err)
	}
	db := New(f)

	db.Write([]byte("Hey"))

	chunk := make([]byte, 4)
	n, _ := db.Read(chunk)

	db.CloseAndRemove()

```

## TODO
- Merge tsbWriter and diskBuffer libraries into one, so that overflowing memory buffers automatically switch to disk option. 
- Make it wrap around