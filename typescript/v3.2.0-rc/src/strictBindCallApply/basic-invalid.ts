function foo(a: number, b: string): string {
    return a + b;
}

// 引数が足りないのでエラー であることが検出できる！
let a = foo.apply(undefined, [10]);
// 2番目の引数があわないのでエラー  であることが検出できる！
let b = foo.apply(undefined, [10, 20]);
// 引数の数が多すぎるのでエラー であることが検出できる！
let c = foo.apply(undefined, [10, "hello", 30]);
// OKなパターン
let d = foo.apply(undefined, [10, "hello"]);

// 引数あわない
let e = foo.call(undefined, 1, 1);
// OK！
let f = foo.call(undefined, 1, "string");

// 1つ目の引数を固定する
let g = foo.bind(undefined, 1);
// NGであることが検出できる！
g(true);
// OK
g("foobar");

type Func = CallableFunction | NewableFunction;

export { a, b, c, d, e, f, g, Func }
