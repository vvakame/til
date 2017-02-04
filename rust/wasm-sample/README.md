# rust + wasm sample

original instruction is [here](https://hackernoon.com/compiling-rust-to-webassembly-guide-411066a69fde).

```
$ cargo run
...
Hello, world!
```

try with Google Chrome Canary.
```
$ rustc --target=wasm32-unknown-emscripten src/main.rs -o index.html
$ python -m SimpleHTTPServer 8989
```
