# merged json

ざっくりした内訳。

* a → 文字列操作で無理やりくっつける
* b → mergo で map を経由する
* c → [zoncoen](https://twitter.com/zoncoen)さんが考えてくれたreflectでfieldの合成 

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

---- >8 ここまで2つのobjectを対象にしてた ここで10つのobjectを対象にしてみる >8 ----

mergeする数が多いとcが一番はやくなるっぽい？

```shell
$ go test -bench . -benchmem
goos: darwin
goarch: arm64
pkg: github.com/vvakame/til/go/merge-json
Benchmark_2objects/a-8  	 3401150	       341.1 ns/op	      48 B/op	       2 allocs/op
Benchmark_2objects/b-8  	  593203	      1955 ns/op	    1240 B/op	      35 allocs/op
Benchmark_2objects/c-8  	 2440664	       491.7 ns/op	     144 B/op	       3 allocs/op
Benchmark_10objects/a-8 	  719061	      1608 ns/op	     192 B/op	      10 allocs/op
Benchmark_10objects/b-8 	  166652	      7122 ns/op	    3701 B/op	      97 allocs/op
Benchmark_10objects/c-8 	  833755	      1444 ns/op	     512 B/op	       3 allocs/op
PASS
ok  	github.com/vvakame/til/go/merge-json	8.146s
```

```shell
$ go test -bench . -benchmem
goos: darwin
goarch: arm64
pkg: github.com/vvakame/til/go/merge-json
Benchmark_2objects/a-8              	 3224755	       349.1 ns/op	      48 B/op	       2 allocs/op
Benchmark_2objects/b-8              	  611488	      1954 ns/op	    1240 B/op	      35 allocs/op
Benchmark_2objects/c-8              	 2452960	       489.7 ns/op	     144 B/op	       3 allocs/op
Benchmark_10objects/a-8             	  708416	      1618 ns/op	     192 B/op	      10 allocs/op
Benchmark_10objects/b-8             	  165700	      7129 ns/op	    3700 B/op	      97 allocs/op
Benchmark_10objects/c-8             	  822145	      1442 ns/op	     512 B/op	       3 allocs/op
Benchmark_10objectsWithLongText/a-8 	     121	   9723101 ns/op	11350365 B/op	      12 allocs/op
Benchmark_10objectsWithLongText/b-8 	     124	   9473763 ns/op	12055171 B/op	      98 allocs/op
Benchmark_10objectsWithLongText/c-8 	     126	   9461746 ns/op	11622189 B/op	       3 allocs/op
PASS
ok  	github.com/vvakame/til/go/merge-json	14.668s
```

json.NewEncoder使ったら早くない？という指摘を反映
構造上 a は文字列操作で先頭や末尾を切り落とす都合上適用が簡単にはできなさそうなので諦め。
c がやっぱり一番よさそう。

```shell
go test -bench . -benchmem                                                                                                                                                                             130 ↵
goos: darwin
goarch: arm64
pkg: github.com/vvakame/til/go/merge-json
Benchmark_2objects/a-8              	 3384633	       344.1 ns/op	      48 B/op	       2 allocs/op
Benchmark_2objects/b-8              	  597603	      1961 ns/op	    1192 B/op	      34 allocs/op
Benchmark_2objects/c-8              	 2425002	       490.8 ns/op	      96 B/op	       2 allocs/op
Benchmark_10objects/a-8             	  721414	      1611 ns/op	     192 B/op	      10 allocs/op
Benchmark_10objects/b-8             	  167528	      7092 ns/op	    3573 B/op	      96 allocs/op
Benchmark_10objects/c-8             	  816115	      1435 ns/op	     384 B/op	       2 allocs/op
Benchmark_10objectsWithLongText/a-8 	     122	   9734517 ns/op	11438749 B/op	      12 allocs/op
Benchmark_10objectsWithLongText/b-8 	     128	   9283150 ns/op	  292878 B/op	      96 allocs/op
Benchmark_10objectsWithLongText/c-8 	     128	   9292385 ns/op	  289698 B/op	       2 allocs/op
PASS
ok  	github.com/vvakame/til/go/merge-json	14.561s
```
