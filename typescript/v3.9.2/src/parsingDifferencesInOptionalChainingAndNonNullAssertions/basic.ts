declare var foo: undefined | {
    bar: undefined | {
        baz: string;
    };
};

// v3.8 までは v1 の型は string になっていた (実行時エラーになる可能性がある)
// v3.9 からは v1 の型は string | undefined になるようになった
let v1 = foo?.bar!.baz;

// --target es2019 で出力した場合
// v3.8 までの結果は次のとおり (undefined).baz という実行パスを通る可能性がある
// let v1 = (foo === null || foo === void 0 ? void 0 : foo.bar).baz;
// v3.9 からの結果は次のとおり optional chaining がぶった切られる問題がなくなった
// let v1 = foo === null || foo === void 0 ? void 0 : foo.bar.baz;


// v3.9 で今までと同じ結果にするにはこう書く必要がある 同じ結果にしたいことはないと思うけど…
let v2 = (foo?.bar)!.baz;
