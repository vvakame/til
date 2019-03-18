// この2つは等価 ReadonlyArray は前からあったやつ
let readonlyArray1: readonly string[] =  ["a", "b", "c"];
let readonlyArray2: ReadonlyArray<string> =  ["a", "b", "c"];

// これは怒られる
// error TS1354: 'readonly' type modifier is only permitted on array and tuple literal types.
// let readonlyArray3: readonly Array<string> =  ["a", "b", "c"];

// 過去のおさらい

let array1 =  ["a", "b", "c"];

readonlyArray1.slice(); // OK
array1.slice(); // OK

// NG! 破壊的変更があるメソッドは ReadonlyArray には存在しない
// error TS2339: Property 'push' does not exist on type 'readonly string[]'.
// readonlyArray1.push('d');
// error TS2339: Property 'sort' does not exist on type 'readonly string[]'.
// readonlyArray1.sort();

// 普通のarrayならもちろんOK
array1.push('d');
array1.sort();

// NG! まぁ普通に代入もだめ
// error TS2542: Index signature in type 'readonly string[]' only permits reading.
// readonlyArray1[0] = "d";

export { }
