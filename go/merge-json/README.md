# merged json

愚直ver

```shell
$ go test -bench . -benchmem
goos: darwin
goarch: arm64
pkg: github.com/vvakame/til/go/merge-json
BenchmarkMarshallers/a-8 	 3126638	       366.0 ns/op	     112 B/op	       3 allocs/op
BenchmarkMarshallers/b-8 	  615163	      1966 ns/op	    1240 B/op	      35 allocs/op
BenchmarkMarshallers/c-8 	  675745	      1720 ns/op	    1368 B/op	      21 allocs/op
PASS
ok  	github.com/vvakame/til/go/merge-json	4.155s
```

aでsync.Pool使う版

```shell
$ go test -bench . -benchmem                                                                                                                                                                             130 ↵
goos: darwin
goarch: arm64
pkg: github.com/vvakame/til/go/merge-json
BenchmarkMarshallers/a-8 	 3338532	       347.2 ns/op	      48 B/op	       2 allocs/op
BenchmarkMarshallers/b-8 	  601790	      1948 ns/op	    1240 B/op	      35 allocs/op
BenchmarkMarshallers/c-8 	  708204	      1699 ns/op	    1368 B/op	      21 allocs/op
PASS
ok  	github.com/vvakame/til/go/merge-json	4.197s
```

cでreflect系の処理をキャッシュする版

```shell
$ go test -bench . -benchmem
goos: darwin
goarch: arm64
pkg: github.com/vvakame/til/go/merge-json
BenchmarkMarshallers/a-8 	 3400077	       342.7 ns/op	      48 B/op	       2 allocs/op
BenchmarkMarshallers/b-8 	  610495	      1962 ns/op	    1240 B/op	      35 allocs/op
BenchmarkMarshallers/c-8 	 2389891	       488.6 ns/op	     144 B/op	       3 allocs/op
PASS
ok  	github.com/vvakame/til/go/merge-json	4.510s
```
