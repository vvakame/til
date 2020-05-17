# TypeScript v3.9.2 変更点

こんにちは[メルペイ社](https://www.merpay.com/)な[@vvakame](https://twitter.com/vvakame)です。

[TypeScript 3.9](https://devblogs.microsoft.com/typescript/announcing-typescript-3-9/)がアナウンスされました。

* [What's new in TypeScript in 3.9](https://www.typescriptlang.org/docs/handbook/release-notes/typescript-3-9.html)
* [Breaking Changes in 3.9](https://github.com/microsoft/TypeScript/wiki/Breaking-Changes#typescript-39)
* [TypeScript 3.9 Iteration Plan](https://github.com/microsoft/TypeScript/issues/37198)
* [TypeScript Roadmap: January - June 2020](https://github.com/microsoft/TypeScript/issues/36948)

Roadmapは[こちら](https://github.com/Microsoft/TypeScript/wiki/Roadmap#39-may-2020)。

[この辺](https://github.com/vvakame/til/tree/master/typescript/v3.9.2)に僕が試した時のコードを投げてあります。

## 変更点まとめ

* `Promise.all` 関係の型推論を改善 [Improvements in Inference and Promise.all](https://github.com/microsoft/TypeScript/pull/34501)
    * 3.7で `Promise.all` 周りの改善を入れたけどリグレッションが発生してた
    * `undefined` などが絡んだときの型推論がおかしくなってたのを修正
* `awaited` という演算子の導入をしませんでした [What About the `awaited` Type?](https://github.com/microsoft/TypeScript/pull/35998)
    * Promiseを解決したときに得られる値についての挙動を `awaited` としてモデル化しようとした
    * … というのを3.9で出したかったけど、もうちょっと磨いてから出します
* 速度面の改善 **Speed Improvements** [*1](https://github.com/microsoft/TypeScript/pull/36576) [*2](https://github.com/microsoft/TypeScript/pull/36590) [*3](https://github.com/microsoft/TypeScript/pull/36607) [*4](https://github.com/microsoft/TypeScript/pull/36622) [*5](https://github.com/microsoft/TypeScript/pull/36754) [*6](https://github.com/microsoft/TypeScript/pull/36696)
    * コンパイルにかかる速度を改善した 色々と！
    * material-ui や styled-components といった複雑な型を持つライブラリで編集やコンパイルにかかる時間がかなり長かった
    * union types, intersections, conditional types, mappted types など色々
* `@ts-expect-error` の追加 [`// @ts-expect-error` Comments](https://github.com/microsoft/TypeScript/pull/36014)
    * コンパイルエラーを無視させる
    * 主にテスト時に不正な使い方をしたいときなどに使える
    * `@ts-ignore` は次の行にエラーがなくても怒らないが、 `@ts-expect-error` はunusedだと怒られる
* 三項演算子で関数を呼ばずに値として評価している箇所をチェックするようにした [Uncalled Function Checks in Conditional Expressions](https://github.com/microsoft/TypeScript/pull/36402)
    * quick fixもあるよ！
* VSCodeが複数のバージョンのTypeScriptの利用をサポート [selecting different versions of TypeScript](https://code.visualstudio.com/docs/typescript/typescript-compiling#_using-newer-typescript-versions)
    * [JavaScript and TypeScript Nightly extention](https://marketplace.visualstudio.com/items?itemName=ms-vscode.vscode-typescript-next)とか
* Visual Studio関連の [何](https://marketplace.visualstudio.com/items?itemName=TypeScriptTeam.TypeScript-39) [か](https://www.nuget.org/packages/Microsoft.TypeScript.MSBuild)
* Sublime Text 3が[異なるバージョンのTypeScriptの選択をサポート](https://github.com/microsoft/TypeScript-Sublime-Plugin#note-using-different-versions-of-typescript)しただとか
* JSでCommonJSスタイルのimportの自動補完をサポート [CommonJS Auto-Imports in JavaScript](https://github.com/microsoft/TypeScript/pull/37027)
* コードアクションでなにかを操作したときにより元の改行が尊重されるようになった [Code Actions Preserve Newlines](https://github.com/microsoft/TypeScript/pull/36688)
* 関数でreturnを書き忘れていそうなときのQuick Fixを追加 [Quick Fixes for Missing Return Expressions](https://github.com/microsoft/TypeScript/pull/26434)
    * `return` 足してくれたり `{}` を削除してくれたり
* ソリューションスタイルのtsconfig.jsonのサポート [Support for “Solution Style” `tsconfig.json` Files](https://github.com/microsoft/TypeScript/pull/37239)
    * ソリューションスタイル = そのtsconfig.jsonの対象になるファイルがゼロ (≒ references を持つ)
    * 今までは自分が含まれるプロジェクトのtsconfig.json以外で参照されている型も参照可能だったけど(正しく)見えなくなった
* マルチプロジェクトでのシンボルナビゲーションのサポート [Editor Support for Multi-Project Symbol Navigation](https://github.com/microsoft/TypeScript/pull/38027)
* Quick Fixで private でプロパティを生成できるようになった [Declare a `private` property](https://github.com/microsoft/TypeScript/pull/36632)

## 破壊的変更！

今回はかなり色々ありますね。とはいえ実プロジェクトで問題になるような変更は少ないでしょう。

* `?.` と `!` と組み合わせたときの挙動を修正 [Parsing Differences in Optional Chaining and Non-Null Assertions](https://github.com/microsoft/TypeScript/pull/36539)
    * `foo?.bar!.baz` が `(foo?.bar).baz` と評価されていたが、これだとchainが崩れるので修正した
* JSXのテキストとして `}` や `>` ベタが許されなくなった [`}` and `>` are Now Invalid JSX Text Characters](https://github.com/microsoft/TypeScript/pull/36636) ([Quick Fix](https://github.com/microsoft/TypeScript/pull/37436))
    * `}` とか `>` はJSX的にはテキストの位置にでてきてはいけない
    * 今までこれをエラーにしていなかったが、エラーにするようにした
    * `>` であれば `{">"}` か `&gt;` を、 `}` であれば `{"}"}` か `&rbrace;` を使うこと
* IntersectionとOptionalが絡むときのチェックを改善 [Stricter Checks on Intersections and Optional Properties](https://github.com/microsoft/TypeScript/pull/37195)
* Intersectionsでチェックを強化 [Intersections Reduced By Discriminant Properties](https://github.com/microsoft/TypeScript/pull/36696)
    * 無理なものは無理
* getter/setterが enumerable: true だったのをECMAScriptの仕様にそってfalseでdownpileするように変更 [Getters/Setters are No Longer Enumerable](https://github.com/microsoft/TypeScript/pull/32264)
    * IssueにSkateJSで困るんだわ… とか2016年に僕がコメントしてて時代を感じる…
* `T extends any` とかしてたやつでTがany相当の振る舞いをしないようにした [Type Parameters That Extend `any` No Longer Act as `any`](https://github.com/microsoft/TypeScript/pull/37124)
* `export *` してたらとりあえずJSに出力されるようになった [`export *` is Always Retained](https://github.com/microsoft/TypeScript/pull/37124)
* `export *` とかを使ったときの live binding で getter を使うよう修正 [Exports Now Use Getters for Live Bindings](https://github.com/microsoft/TypeScript/pull/35967)
    * tslib に `__createBinding` が追加されたり
* exportするときにプロパティを先に初期化するようにした [Exports are Hoisted and Initially Assigned](https://github.com/microsoft/TypeScript/pull/37093)
    * 1個上の変更に対する互換性維持のため

`tsc --init` で生成した tsconfig.json に https://aka.ms/tsconfig.json の案内が出るようになりましたね。
多分今回から…？

次のVersionは[4.0らしい](https://github.com/microsoft/TypeScript/issues/38510)です。
このまとめを続けるかどうか悩ましい…
ところで[僕のGitHub Sponsor](https://github.com/sponsors/vvakame)はこちらからどうぞ。

## `Promise.all` 関係の型推論を改善

Promise.all 周りの改善です。
標準の型定義を改善したけどリグレッションで正確じゃない場合があったのを直したらしいです。

```ts
let str = "a";
let num : number | undefined = 1 as any;

let strP = Promise.resolve(str);
let numP = Promise.resolve(num);

// rStr1 の型
// TypeScript 3.8 では string | undefined
// TypeScript 3.9 から string
// undefinedのコンタミがなくなった！
let [rStr1, rNum1] = await Promise.all([strP, numP]);


let bool: boolean | null = true as any;
let boolP = Promise.resolve(bool);
// rStr2 の型
// TypeScript 3.8 では string | null | undefined
// TypeScript 3.9 から string
// rNum2 の型
// TypeScript 3.8 では number | null | undefined
// TypeScript 3.9 から number | undefined
// undefined, null のコンタミがなくなって直感的な型に
let [rStr2, rNum2, rBool2] = await Promise.all([strP, numP, boolP]);
```

なるほどねー…。

## `awaited` という演算子の導入をしませんでした

しなかったらしいです。
色々問題が発生したのでもうちょっと揉むわー的な。

https://github.com/microsoft/TypeScript/pull/35998
https://github.com/microsoft/TypeScript/issues/37897
https://github.com/microsoft/TypeScript/pull/37610

## 速度面の改善

なんか色々やったらしいです。

https://github.com/microsoft/TypeScript/pull/36576
https://github.com/microsoft/TypeScript/pull/36590
https://github.com/microsoft/TypeScript/pull/36607
https://github.com/microsoft/TypeScript/pull/36622
https://github.com/microsoft/TypeScript/pull/36754
https://github.com/microsoft/TypeScript/pull/36696

## `@ts-expect-error` の追加

`@ts-expect-error` が追加されました。
すでにある `@ts-ignore` に類似のものです。
`@ts-ignore` がコンパイルエラーを無視させるのに対して、 `@ts-expect-error` もコンパイルエラーを無視させますが、エラーがない場合エラーになります。

```ts
function doStuff(a: string, b: string) {
    // @ts-expect-error
    // コンパイルエラーを抑制することができる (今回追加された)
    hello(1);

    // @ts-ignore
    // コンパイルエラーを抑制することができる (前からある)
    hello(1);

    // // @ts-expect-error
    // error TS2578: Unused '@ts-expect-error' directive.
    // 何もしなくてもコンパイルが正しく通る場合、コンパイルエラーになる！
    hello("cat");

    // @ts-ignore
    // 特に効果がなくても特になんもならない
    hello("cat");
}

function hello(word: string) {
    console.log(`Hello, ${word}`);
}
```

`@ts-ignore` から `@ts-expect-error` に乗り換えるべきか？という話題について。

* `@ts-expect-error` がおすすめな状況
  * テストコードを記述していて、わざとエラーになるような操作を記述したい場合
  * 近いうちに修正される予定で、急場をしのぐ回避策がほしい場合
  * 働きアリがもりもり働いていて、これらの回避策をすぐ修正してくれそうな場合
* `@ts-ignore` を使ったほうがよい状況
  * 大きなプロジェクトで、コードのオーナーが不明瞭な箇所がエラーになった場合
  * TypeScriptのバージョンを更新しようとして、どちらかのバージョンではエラーになるがもう片方ではエラーにならない場合
  * どっちがマシか判断する時間もないような場合

だそうです。

## 三項演算子で関数を呼ばずに値として評価している箇所をチェックするようにした

タイトルのまんまです。

```ts
function hasImportantPermissions(): boolean {
    return true;
}

function deleteAllTheImportantFiles() {
}

// 関数の呼び出しを忘れて条件式に使った場合、エラーになる (TypeScript v3.7 で入った)
// error TS2774: This condition will always return true since the function is always defined. Did you mean to call it instead?
if (hasImportantPermissions) {
    deleteAllTheImportantFiles();
}

// 三項演算子でもこのチェックが働くようになった (今回から)
// error TS2774: This condition will always return true since the function is always defined. Did you mean to call it instead?
hasImportantPermissions ? deleteAllTheImportantFiles() : void 0;

// なお、 Add missing call parentheses のQuick Fixも追加されたので、エラーになっている箇所を見つけたらｼｭｯと直せる
```

## VSCodeが複数のバージョンのTypeScriptの利用をサポート

前からできなかったっけ…？

## Visual Studio関連の何か

興味がないので調査してない

## Sublime Text 3が異なるバージョンのTypeScriptの選択をサポート

興味がないので調査してない

## JSでCommonJSスタイルのimportの自動補完をサポート

らしいです。

```js
const { } = require("fs");

// ここで readFileSync とかが補完候補にあがり、選ぶとimport文にいい感じに追加される
```

## コードアクションでなにかを操作したときにより元の改行が尊重されるようになった

地味に不便してたやつです。

```ts
const maxValue = 100;

// start から end の範囲を Extract to function in Global scope のrefactorを行うとする
/*start*/
for (let i = 0; i <= maxValue; i++) {
    // First get the squared value.
    let square = i ** 2;

    // Now print the squared value.
    console.log(square);
}
/*end*/

// 今まで 改行のみ の行は尊重されず、潰されてしまっていた
function newFunction1() {
    for (let i = 0; i <= maxValue; i++) {
        // First get the squared value.
        let square = i ** 2;
        // Now print the squared value.
        console.log(square);
    }
}

// 3.9 からちゃんと元の改行が保持されるようになった
function newFunction2() {
    for (let i = 0; i <= maxValue; i++) {
        // First get the squared value.
        let square = i ** 2;

        // Now print the squared value.
        console.log(square);
    }
}
```

## 関数でreturnを書き忘れていそうなときのQuick Fixを追加

まんまです。

```ts
let numberArray = [1, 9, 2, 8, 3, 7, 4, 6, 5];

// sort(compareFn?: (a: T, b: T) => number): this; という定義
// 次の書き方だと sort に渡した関数が値を返していない(のでコンパイルエラーになる)
// numberArray.sort((a, b) => { a - b });

// Quick Fix が追加されているので試してみよう！

// Add a return statement を選ぶと次のように変形される
numberArray.sort((a, b) => { return a - b; });

// Remove block body braces を選ぶと次のように変形される
numberArray.sort((a, b) => a - b);
```

## ソリューションスタイルのtsconfig.jsonのサポート

ソリューションスタイル = そのtsconfig.jsonの対象になるファイルがゼロ (≒ references を持つ) 状態だそうです。
Visual Studioのslnとprj的な感じなんですかね。

今までは自分が含まれるプロジェクトのtsconfig.json以外で参照されている型も参照可能だったけど(正しく)見えなくなったらしいです。

https://github.com/microsoft/TypeScript/issues/27372 で報告されているやつ。

報告社は再現用のリポジトリも用意している。
https://github.com/OliverJAsh/jest-ts-project-references

次のような tsconfig.json を使っているとする。

```json5
{
  // By default, all files will be included. We want to use our project references instead, so we
  // clear these.
  "files": [],
  "include": [],
  "references": [
    {
      "path": "./tsconfig-src.json"
    },
    {
      "path": "./tsconfig-tests.json"
    }
  ]
}
```

srcのほうでは `"types": []` 、testsのほうでは `"types": ["jest"]` だったとしても、src側でjestの型定義が参照してコンパイルとかも通っていた。
のを、しっかり管理するようにしてそれぞれのファイルが自分が影響される設定に従ってチェックされるようになった… という感じのようだ。

## マルチプロジェクトでのシンボルナビゲーションのサポート

なんかできるようになったらしいです。

## Quick Fixで private でプロパティを生成できるようになった

そのまんま。

```ts
class A {
    constructor() {
        // Declare private property '_x' が追加
        this._x = 10;
    }
}
```

## `?.` と `!` と組み合わせたときの挙動を修正

JSにdownpileされるときに生成されるコードが修正されたという話ですね。

```ts
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
```

## JSXのテキストとして `}` や `>` ベタが許されなくなった

許されてたのはinvalidだったらしいです。
そうだったんか…！

```tsx
// JSXのテキストで、JSXの仕様的に使ってはいけない > とか } の利用がコンパイルエラーになるようになった
// error TS1382: Unexpected token. Did you mean `{'>'}` or `&gt;`?
// error TS1381: Unexpected token. Did you mean `{'}'}` or `&rbrace;`?
let node = (
    <div>
        テキストに > とか } があったらダメ！
        これからは {">"} とか {"}"} を使うか、
        &gt; や &rbrace; を使おう！
    </div>
);

// Quick Fix も2パターン追加されている
// Wrap invalid character in an expression container
// Convert invalid character to its html entity code
// …が、 v3.9.2 + VSCode 1.46.0-insider の組み合わせだと選んでも直してくれないぽい
```

## IntersectionとOptionalが絡むときのチェックを改善

まだこの手のが発掘されるのってすげーな… と思ってしまう…。

```ts
interface A {
    a: number; // ここの a は number
}

interface B {
    b: string;
}

interface C {
    a?: boolean; // ここの a は boolean
    b: string;
}

declare let x: A & B;
declare let y: C;

// y = x;
// ↑ 上記コードは A と C は一致しないが B と C は一致するので許されていた
// が、当然実行時エラーになる可能性がある
// v3.9 から、 A & B は C に一致しないので怒られるようになった
// error TS2322: Type 'A & B' is not assignable to type 'C'.
//   Types of property 'a' are incompatible.
//   Type 'number' is not assignable to type 'boolean | undefined'.
```

## getter/setterが enumerable: true だったのをECMAScriptの仕様にそってfalseでdownpileするように変更

classでgetter, setterがある場合、それがは `enumerable: false` で生成されるのがECMAScriptの仕様的には正しいが `true` で生成されたのを修正。
Issueを見ると、2016年に僕がSkateJSでbabelと互換性無くて困る〜〜みたいなことを言ってて、時代を感じますね。
その時はあまり簡単に直せるようなコードベースだった気がするのでPR送らなかったんだと思いますが、わりとｼｭｯと治っていてびびる。

```ts
class A {
    get a() {
        return "a";
    }
    set a(v: string) {
        console.log(v);
    }
}

// v3.8 での出力
// var A = /** @class */ (function () {
//     function A() {
//     }
//     Object.defineProperty(A.prototype, "a", {
//         get: function () {
//             return "a";
//         },
//         set: function (v) {
//             console.log(v);
//         },
//         enumerable: true,
//         configurable: true
//     });
//     return A;
// }());

// v3.9 での出力
// var A = /** @class */ (function () {
//     function A() {
//     }
//     Object.defineProperty(A.prototype, "a", {
//         get: function () {
//             return "a";
//         },
//         set: function (v) {
//             console.log(v);
//         },
//         enumerable: false,
//         configurable: true
//     });
//     return A;
// }());
```

## `T extends any` とかしてたやつでTがany相当の振る舞いをしないようにした

そのまんまです。

```ts
function foo<T extends any>(arg: T) {
    // v3.8 までは T は any として扱われてたのでエラーなし
    // v3.9 からは T は extends unknown と同等に振る舞いっぽ
    // error TS2339: Property 'notExists' does not exist on type 'T'.
    arg.notExists;
}
```

これはまぁなんというか元のままでもよかったのでは…という気もする。
lintが `extends any` とかするなっ！って言ってくれればよいようなまぁ改善のような…。

## `export *` してたらとりあえずJSに出力されるようになった

何もexportしてないような .ts があるとして

```empty.ts
export {}
```

これを適当に `export *` してた場合、生成するJSからはomitされてたんだけどこれが何も考えずにベタに出るようになった。

```ts
export * from "./empty";

// v3.8 では参照先が何もexportしていない場合、JSの出力ではexportは削られていた
// v3.9 では何も考えずとりあえず出力するようになった
```

マジカルが減る方向性なのでよいと思います。
v3.8で `import type` も導入されているので、適宜必要なところで使っていきたいですね。

## `export *` とかを使ったときの live binding で getter を使うよう修正

そもそもECMAScriptがそういう仕様なの知らなかった…。

変数をexportし、その変数を変えられるような構造を持つとする。

```a.ts
export let data = 1;

export function setData(v: number) {
    data = v;
}
```

これを参照し使った場合、定義していたモジュールを直接参照している場合はそれっぽく動いている。

```b.ts
import { data, setData } from "./a";

console.log("b-1", data);
setData(2);
console.log("b-2", data);

// v3.8, v3.9 ともに --module commonjs の場合 次のように出力される
// b-1 1
// b-2 2

export { data, setData };
```

が、間接的に参照された場合、変数がgetter経由で提供されていないので、関係性が断絶してしまう。
のが、今回治った、というお話。

```c.ts
import { data, setData } from "./b";

console.log("c-1", data);
setData(3);
console.log("c-2", data);

// v3.8 だと --module commonjs の場合 次のように出力される
// c-1 1
// c-2 1
// v3.9 だと --module commonjs の場合 次のように出力される
// c-1 2
// c-2 3

export { data, setData };
```

この挙動をサポートするため、 tslib に `__createBinding` みたいなヘルパ関数が追加されました。
tslib も v2.0.0 になり、 v1 系は TypeScript 3.8 まで用、 v2 で TypeScript 3.9 以降用と使い分ける必要があるようです。

実際に `__createBinding` が使われる出力を得たい場合、次のようなコードを試す。

```d.ts
import * as b from "./b";
export * from "./b";

// __createBinding みたいなヘルパ関数が増えた
// tslib も v2.0.0 になった様子 v1 系は TypeScript 3.8 まで用、 v2 で TypeScript 3.9 以降用

b.data;
```

## exportするときにプロパティを先に初期化するようにした

1つ上の変更による余波。

```foo.ts
export * from "./foo";
// ./foo でも nameFromFoo を 定義して export している
// なので、ここで上書きできなければならない
export const nameFromFoo = 0;

// 出力されるJSコードの比較

// v3.8 までの出力 特に気にする点はない
// __export(require("./foo"));
// exports.nameFromFoo = 0;

// v3.9 からの出力 先に undefined で初期化されるようになった
// exports.nameFromFoo = void 0;
// __exportStar(require("./foo"), exports);
// exports.nameFromFoo = 0;
```

```bar.ts
export * from "./foo";
// ./foo でも nameFromFoo を 定義して export している
// なので、ここで上書きできなければならない
export const nameFromFoo = 0;

// 出力されるJSコードの比較

// v3.8 までの出力 特に気にする点はない
// __export(require("./foo"));
// exports.nameFromFoo = 0;

// v3.9 からの出力 先に undefined で初期化されるようになった
// exports.nameFromFoo = void 0;
// __exportStar(require("./foo"), exports);
// exports.nameFromFoo = 0;

//  __exportStar の中で __createBinding を使うようになった
// その中で getter を定義してしまうので、後から上書きできない
// TypeError: Cannot set property nameFromFoo of #<Object> which has only a getter
// 上書きしたい場合、先にプロパティを作っておくと無視してくれるので別途後から値を設定する
```

`__exportStar` の中で `__createBinding` を使うようになった。
それで、その中で getter を定義してしまうので、後から上書きできない。
先にプロパティの存在を主張しておかないと、 `TypeError: Cannot set property nameFromFoo of #<Object> which has only a getter` とか言われてしまうため。
上書きしたい場合、先にプロパティを作っておくと無視してくれるので別途後から値を設定するようなコードが生成されるようになった。


```use.ts
import { nameFromFoo } from "./bar";

// 0 と表示される
console.log(nameFromFoo);
```
