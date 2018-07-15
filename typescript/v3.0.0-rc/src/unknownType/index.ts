let fooAny: any = 10;
let fooUnknown: unknown = 10;
let fooObject: {} = 10;

// any は存在しないプロパティをにアクセスしてもエラーにならない
console.log(fooAny.notExists);

// unknownの場合ちゃんとエラーになる！
// error TS2571: Object is of type 'unknown'.
// console.log(fooUnknown.notExists);

// まぁ {} でも似たようなことができる
// error TS2339: Property 'notExists' does not exist on type '{}'.
// console.log(fooObject.notExists);

function double(n: number) {
    return n * 2;
}

// 型アサーションや型の絞り込みを行うとその型として扱える
if (typeof fooUnknown === "number") {
    console.log(double(fooUnknown));
}
if (typeof fooObject === "number") {
    console.log(double(fooObject));
}

// unknown にはなんでも入る
fooUnknown = true;
fooUnknown = void 0;

// {} はただの空オブジェクトなので undefined とかは無理
fooObject = true;
// fooObject = void 0;

// unknown は unknown と any にのみ代入可能
let unknown: unknown = fooUnknown;
let any: any = fooUnknown;
// error TS2322: Type 'unknown' is not assignable to type 'string'.
// let str: string = fooUnknown;

// 予約後なので上書き不可
// type unknown = string;
