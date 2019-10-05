function* bar() {
    // ここで変数の型からnextに投げ込める値の型が決められるのすごい
    let x: { hello(): void } = yield;
    x.hello();
}

let iter = bar();
iter.next();

// 次のものを投げ込もうとするとコンパイルエラーになる
// iter.next(123);

// 型があってればちゃんとコンパイル通る えらい
iter.next({ hello() { } });

export {}
