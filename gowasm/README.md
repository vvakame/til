# Go+wasm

https://github.com/golang/go/issues/18892
https://blog.gopheracademy.com/advent-2017/go-wasm/

```
$ cd $HOME/Dropbox/work/gowasm-work
$ git clone https://go.googlesource.com/go
$ cd go
$ git remote add neelance https://github.com/neelance/go
$ git fetch --all
$ git checkout wasm-wip
$ ch src
$ ./make.bash
```

```
export GOROOT=$HOME/Dropbox/work/gowasm-work/go
export GOPATH=$HOME/Dropbox/work/gowasm-work/gopath
export PATH=$HOME/Dropbox/work/gowasm-work/go/bin:$PATH
```
