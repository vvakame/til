let foo: any = { bar: { baz: true } };

// 今までのやり方
if (foo && foo.bar && foo.bar.baz) {
    console.log(foo.bar.baz);
}

// Optional Chainingを使うとこう書ける
if (foo?.bar?.baz) {
    console.log(foo.bar.baz);
}

// && と ? では厳密には挙動が異なる
// && は falsy な値 (null, undefined, "", 0, NaN, false) の場合処理を打ち切り、左辺の値を返す
// ? の場合、 null と undefined の時のみ処理を打ち切り、undefined を返す

// undefined と表示される (toStringは実行されないので)
console.log((null as any)?.toString());
// 実行時エラー Cannot read property 'toString' of null
console.log((null as any && true).toString());

function barPercentage(foo?: { bar: number }) {
    // こういうのもダメ foo?.bar の部分でエラーとなる
    // error TS2532: Object is possibly 'undefined'.
    // return foo?.bar / 100;
    // このように解釈されている
    // let tmp: number | undefined = (foo === null || foo === void 0) ? void 0 : foo.bar;
    // return tmp / 100;

    // こうすればOK
    return (foo?.bar ?? 0) / 100;
}

export { }
