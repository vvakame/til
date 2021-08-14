# merged json

```shell
go test -bench . -benchmem
goos: darwin
goarch: arm64
pkg: github.com/vvakame/til/go/merge-json
BenchmarkMarshallers/a-8 	 3126638	       366.0 ns/op	     112 B/op	       3 allocs/op
BenchmarkMarshallers/b-8 	  615163	      1966 ns/op	    1240 B/op	      35 allocs/op
BenchmarkMarshallers/c-8 	  675745	      1720 ns/op	    1368 B/op	      21 allocs/op
PASS
ok  	github.com/vvakame/til/go/merge-json	4.155s
```
