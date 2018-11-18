こんにちは[メルペイ社](https://www.merpay.com/jp/)な[@vvakame](https://twitter.com/vvakame)です。

[TypeScript 3.1 RC](https://blogs.msdn.microsoft.com/typescript/2018/09/13/announcing-typescript-3-1-rc/)と[3.1](https://blogs.msdn.microsoft.com/typescript/2018/09/27/announcing-typescript-3-1/)が(とっくに)アナウンスされました。
アナウンスされたのは9月13日とかで、技術書典の準備がハチャメチャ忙しいのにやってられるか！ってスルーしてました。
そしたら[3.2 RC](https://blogs.msdn.microsoft.com/typescript/2018/11/15/announcing-typescript-3-2-rc/)のアナウンスが出て、まぁ今までずっとやってるしv3超えたけど惰性継続するか…ということで3.2のことを書く前に3.1を書きます。

[What's new in TypeScriptも更新](https://github.com/Microsoft/TypeScript/wiki/What's-new-in-TypeScript#typescript-31)されているようです。
[破壊的変更](https://github.com/Microsoft/TypeScript/wiki/Breaking-Changes#typescript-31)もあるよ！

[この辺](https://github.com/vvakame/til/tree/master/typescript/v3.1.6)に僕が試した時のコードを投げてあります。

## 変更点まとめ

* タプルと配列にmapped types適用した時の挙動を改善 [Mappable tuple and array types](https://github.com/Microsoft/TypeScript/pull/26063)
    * `.length` とか `forEach` とか全てのプロパティが影響受けてたのを直した
* 名前付き関数かconstに関数を割り当てた場合、JSっぽい書き方でプロパティを生やせるようにした [Properties on function declarations](https://github.com/Microsoft/TypeScript/pull/26368)
    * 主にReactの `defaultProps` のためだって
* `typesVersions` を package.json に導入 [Version selection with typesVersions](https://github.com/Microsoft/TypeScript/issues/22605)
    * 1つのpackageで複数バージョンのTypeScriptに対応できるようになった
* `lib.d.ts` からベンダプリフィクス付きの定義をのけた [Vendor-specific declarations removed](https://github.com/Microsoft/TypeScript/pull/25944)
    * `webkitRequestFullScreen` とかいろいろ
    * Edge WebIDLから標準のWebIDLベースに変えたらしい
* `typeof` によるGenericsのナローイングの挙動の変更 [Differences in narrowing functions](https://github.com/Microsoft/TypeScript/pull/25243)
    * Genericsの `T` とかを `typeof` に突っ込んだ時の型のナローイングのされ方が変わった
    * それに伴って今まで通ってた関数についてのナローイングが通らなくなったりする
* エラーのUXの改善 [Error UX improvements](https://github.com/Microsoft/TypeScript/issues/26077)
    * なんかまぁいろいろ変わったっぽい
    * 3.0の時のIssueが引き継ぎで使われていて調査がめんどいので諦める！
* import/exportのパスを変えるリファクタリングでファイル名のリネームもついでにやる [Rename files from import/export paths](https://github.com/Microsoft/TypeScript/pull/26373)
    * その名のとおり
* Promiseのthen,catcを使ったコードをasync/awaitとtry/catchを使ったコードに変換する [Convert from Promise#then/catch to async/await](https://github.com/Microsoft/TypeScript/pull/26373)
    * その名のとおり

### 破壊的変更！

上記リストのうち、破壊的変更を伴うのは次のものになります。

* `lib.d.ts` からベンダプリフィクス付きの定義をのけた
* `typeof` によるGenericsのナローイングの挙動の変更

## タプルと配列にmapped types適用した時の挙動を改善

```ts
{ [P in keyof T]: X }
```
的な変換をした時に、Arrayに直接適用すると、配列の要素以外のプロパティ、例えば `.length` とか `.forEach` とかも一切合切変換されてしまっていました。
仕様的には正しいのですが、一般的なシチュエーションとしてこの挙動が求められることはほぼ無いでしょう。
ということで、Array(とTuple)にmapped typesを適用した場合、要素にのみその操作を適用した型を返す、という挙動になったようです。

```ts
// 全プロパティの値をstring型にした型を作る
type Stringify<T> = { [K in keyof T]: string };

function stringifyProps<T>(v: T): Stringify<T> {
  const result = {} as Stringify<T>;
  for (const prop in v) {
    result[prop] = String(v[prop]);
  }
  return result;
}

{ // よく見る一般的なArrayとか相手ではないmapped typesの操作
  const obj: { no: string } = stringifyProps({ no: 151 });
  console.log(obj);
}

function stringifyAll<T extends unknown[]>(...args: T): Stringify<T> {
  return args.map(v => String(v)) as any;
}

{
  const array = stringifyAll([1, true]);
  // TypeScript 3.1以前だと forEach も length も string と解釈される
  // 一般的には配列の要素部分だけ変換されれば十分だよなぁ…？
  array.forEach(v => console.log(v));
  const len: number = array.length;
  console.log(array, len);
}
{
  const tuple: [1, true] = [1, true];
  const array = stringifyAll(tuple);
  // TypeScript 3.1以前だと forEach も length も string と解釈される
  array.forEach(v => console.log(v));
  const len: number = array.length;
  console.log(tuple, array, len);
}
```

## 名前付き関数かconstに関数を割り当てた場合、JSっぽい書き方でプロパティを生やせるようにした

見出しのまんまですね。
主にReactのdefaultPropsの書き方をよりJSっぽくできるようにしたいということで追加された機能のようです。

```ts
// 今までのやり方。namespace はTypeScriptの独自要素。
function fooA() {
    console.log("fooA");
}
namespace fooA {
  export var barA = () => {
    console.log("fooA.barA");
  };
}

fooA();
fooA.barA();

// JSだったら普通素直にこう書くよね…？という書き方ができるようになった。
function fooB() {
    console.log("fooB");
}
fooB.barB = () => {
    console.log("fooB.barB");
};

fooB();
fooB.barB();

// const+関数でもできる
const fooC = () => {
    console.log("fooC");
}
fooC.barC = () => {
    console.log("fooC.barC");
};
fooC();
fooC.barC();

// let(やvar)+関数ではNG
let fooD = () => {
    console.log("fooD");
}
// この書き方はエラーになる fooDの値は差し替え可能だからね。仕方ないね。
// error TS2339: Property 'barD' does not exist on type '() => void'.
// fooD.barD = () => {
//     console.log("fooD.barD");
// };
// fooD();
// fooD.barD();
```

namespaceの用途がまた1つ減って喜ばしいですね。

## `typesVersions` を package.json に導入

npmパッケージが複数の型定義を提供し、利用するTypeScriptのバージョンによってどれを使うかを切り替えられるようになったそうです。
こんな感じ。

```json
{
  "name": "package-name",
  "version": "1.0",
  "types": "index",
  "typesVersions": {
    ">=3.1.0-0": { "*": ["ts3.1/*"] }
  }
}
```

バージョンの指定は[semverのrange](https://github.com/npm/node-semver#ranges)に従う。
マッチングは上から順番に試される。
詳しくは[What's new](https://github.com/Microsoft/TypeScript/wiki/What%27s-new-in-TypeScript#version-selection-with-typesversions)の記述を読むのがよさそう。

## `lib.d.ts` からベンダプリフィクス付きの定義をのけた

`lib.dom.d.ts` の生成元をEdge WebIDLから標準のWebIDLに変えたことによるものらしいですね…。
具体的に何が消えたかは[Breaking Changes](https://github.com/Microsoft/TypeScript/wiki/Breaking-Changes#some-vendor-specific-types-are-removed-from-libdts)の該当項目を読むのがよさそう。

## `typeof` によるGenericsのナローイングの挙動の変更

`typeof` による型のナローイングを行った時の挙動が変更されました。
対象の変数がunionだった場合、全ての部分にintersectionが配られる的な挙動になったようです。

```ts
function fooA<T>(x: T | (() => string)) {
  if (typeof x === "function") {
    // 3.1からこれがエラーになる
    //   今まで x は () => string と推論されていた
    //   これからは x は (() => string) | (T & Function) と推論されるようになるため
    //   () => string と Function には共通のcall signatureがない
    // error TS2349: Cannot invoke an expression whose type lacks a call signature.
    //   Type '(() => string) | (T & Function)' has no compatible call signatures.
    x();
  }
}
// () => string ←これじゃない関数を弾きたいとかそういう
fooA((name: string) => `Hello, ${name}!`);
```

```ts
// 回避するために、Tに実際求める追加の制約を与えたりしよう
function fooB<T extends { word: string }>(x: T | (() => string)) {
  if (typeof x === "function") {
    // T は制約として word プロパティを持つ必要があるので T & function は除外されるのでコンパイル通る
    return x();
  }
  return `Hello, ${x.word}`;
}

// T & Function 側にマッチしなくなるので平和
fooB(() => `Hello, TypeScript!`);
fooB({ word: "world" });

// このパターンはinvalidなんだけど現状すり抜けられてしまう
const func = (word: string) => `Hello, ${word}`;
func.word = "invalid";
fooB(func);
```

この辺りの挙動は更に[3.3で修正が入る可能性](https://github.com/Microsoft/TypeScript/issues/26970)がありそうです。

## エラーのUXの改善

地道に改善されているようです。
偉い！

## import/exportのパスを変えるリファクタリングでファイル名のリネームもついでにやる

そのまんまです。

```ts
import { hello } from "./b";

hello();
```

で `"./b"` のところでF2とかでrenameの操作をすると対象のファイルのパスが書き換わります。
簡単に試した限りでは別のモジュール内のパスも書き換わりました。便利。

LanguageServerの常なのですが、これは機能的に "bファイルがcにリネームされました" という情報を送るのみなので、実際のリネーム操作はエディタ側が実装に追従しなければなりません。

## Promiseのthen,catcを使ったコードをasync/awaitとtry/catchを使ったコードに変換する

リファクタリングの機能ですね。

```ts
function timeout(timeout?: number) {
  return new Promise(resolve => {
    setTimeout(() => {
      resolve();
    }, timeout);
  });
}

function exec() {
  return timeout(100)
    .then(() => {
      if (Math.random() > 0.5) {
        throw new Error("random failed");
      }
      console.log("hi!");
    })
    .catch(err => {
      console.error(err);
    });
}
```

こういうコードのexecに対して `Convert to async function` を適用すると

```ts
async function exec() {
  try {
    await timeout(100);
    if (Math.random() > 0.5) {
      throw new Error("random failed");
    }
    console.log("hi!");
  }
  catch (err) {
    console.error(err);
  }
}
```

こうなります。
良いですね。
