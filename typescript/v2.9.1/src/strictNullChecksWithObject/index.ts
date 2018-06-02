function f1<T>(x: T) {
    // 今回からエラー！ T を object に代入できちゃうのはヤバい！エラーになるのは正しい！
    const y: object = x;
    console.log(y);
}
f1("string");
f1({});

// T の底を object にする
function f2<T extends object>(x: T) {
    const y: object = x;
    console.log(y);
}
// string は object ではないのでエラーになる
f2("string");
f2({});

function f3<T>(x: T) {
    // object をやめて {} にする
    const y: {} = x;
    console.log(y);
}
// 両方OK
f3("string");
f3({});
