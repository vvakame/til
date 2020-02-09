let obj1: { [x: string]: number } | { a: number };

// 余計かつ型的にマッチしないものがあると怒られるようになった
// error TS2322: Type 'string' is not assignable to type 'number'.
obj1 = { a: 5, c: 'abc' };
// もちろんこれはOK
obj1 = { a: 5, c: 3 };

// object literalの直接代入じゃなければ許されるのは相変わらず
let obj1a = { a: 5, c: 'abc' };
obj1 = obj1a;


let obj2: { a: number };
// 余計なものがあると怒られるヤツとはエラーメッセージが異なる…！ (昔からエラー)
// error TS2322: Type '{ a: number; c: string; }' is not assignable to type '{ a: number; }'.
//  Object literal may only specify known properties, and 'c' does not exist in type '{ a: number; }'.
obj2 = { a: 5, c: 'abc' };


let obj3: { [x: string]: number } | { [x: number]: number };
// これも怒られる。今までなんでOKだったんや…
// error TS2322: Type 'string' is not assignable to type 'number'.
obj3 = { a: 'abc' };

export { }
