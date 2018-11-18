function fooA<T>(x: T | (() => string)) {
  if (typeof x === "function") {
    // 3.1からこれがエラーになる
    //   今まで x は () => string と推論されていた
    //   これからは x は (() => string) | (T & Function) と推論されるようになるため
    //   () => string と Function には共通のcall signatureがない
    // error TS2349: Cannot invoke an expression whose type lacks a call signature.
    //   Type '(() => string) | (T & Function)' has no compatible call signatures.
    // x();
  }
}
// () => string こうなってない関数を弾きたいとかそういう
console.log(fooA((word: string) => `Hello, ${word}!`));
console.log(fooA(() => `Hello, world!`));

// 回避するために、Tに実際求める追加の制約を与えたりしよう
function fooB<T extends { word: string }>(x: T | (() => string)) {
  if (typeof x === "function") {
    // T は制約として word プロパティを持つ必要があるので T & function は除外されるのでコンパイル通る
    return x();
  }
  return `Hello, ${x.word}`;
}

// T & Function 側にマッチしなくなるので平和
console.log(fooB(() => `Hello, TypeScript!`));
console.log(fooB({ word: "world" }));

// このパターンはinvalidなんだけど現状すり抜けられてしまう
const func = (name: string) => `Hello, ${name}`;
func.word = "invalid";
console.log(fooB(func));

export {};
