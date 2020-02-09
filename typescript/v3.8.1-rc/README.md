# TypeScript v3.8.1 RC 変更点

こんにちは[メルペイ社](https://www.merpay.com/)な[@vvakame](https://twitter.com/vvakame)です。

[TypeScript 3.8 RC](https://devblogs.microsoft.com/typescript/announcing-typescript-3-8-rc/)がアナウンスされました。

* ~~[What's new in TypeScript in 3.8](https://www.typescriptlang.org/docs/handbook/release-notes/typescript-3-8.html)~~
* [Breaking Changes in 3.8](https://github.com/microsoft/TypeScript/wiki/Breaking-Changes#typescript-38)
* [TypeScript 3.8 Iteration Plan](https://github.com/microsoft/TypeScript/issues/34898)
* [TypeScript Roadmap: July - December 2019](https://github.com/microsoft/TypeScript/issues/33118)

Roadmapは[こちら](https://github.com/Microsoft/TypeScript/wiki/Roadmap#38-february-2020)。

[この辺](https://github.com/vvakame/til/tree/master/typescript/v3.8.1-rc)に僕が試した時のコードを投げてあります。

## 変更点まとめ

* 型のみのimport, exportのサポート [Type-Only Imports and Exports](https://github.com/microsoft/TypeScript/pull/35200)
    * `import type { foo } from "./foo";` や `export type { foo };` をサポートした
    * flowtypeには前からあったやつですね
* ECMAScript Private Fieldsのサポート [ECMAScript Private Fields](https://github.com/Microsoft/TypeScript/pull/30829)
    * `class { #name: string; }` 的なやつのサポート
    * `#` やだなぁ…
* `export * as foo ...` 的な構文のサポート [export * as ns Syntax](https://github.com/microsoft/TypeScript/pull/34903)
    * `export * as utils from "./utils";` 的な構文のサポート
* トップレベルのawaitのサポート [Top-Level await](https://github.com/microsoft/TypeScript/issues/25988)
    * そのまんま
* `es2020` が `target` と `module` に入った [Add ES2020 matchAll APIs](https://github.com/microsoft/TypeScript/pull/30936)
* JSDocでの修飾子のサポート [JSDoc Property Modifiers *1](https://github.com/microsoft/TypeScript/pull/35731) [*2](https://github.com/microsoft/TypeScript/issues/17233)
    * JS環境下で `@public`, `@private`, `@protected`, `@readonly` が使えるように
* より良いディレクトリ監視と `watchOptions`の追加 [Better Directory Watching and watchOptions](https://github.com/microsoft/TypeScript/pull/35615)
    * `node_modules` がクソ重たいのでインストール直後などに過敏に反応しすぎないようにした という感じ？
    * tsconfig.json, jsconfig.json に `watchOptions` の項目を追加
        * `watchFile`, `watchDirectory`, `fallbackPolling`, `synchronousWatchDirectory` が下にある
* ちょっと怠惰なインクリメンタルチェック [“Fast and Loose” Incremental Checking](https://github.com/microsoft/TypeScript/issues/33329)
    * `assumeChangesOnlyAffectDirectDependencies` の追加
    * 変更されたファイルの依存ツリーをあまり厳密に追わないかわりに早くする

EditorとかProducitivity系の変更は新機能紹介で取り上げられなくなってきた気がしますね…。
Iteration Planからたどるとこれ以外にも結構色々改善されています。

## 破壊的変更！

* index signatureを含むunion型への値の代入のチェックをもうちょっと厳格にした [Stricter Assignability Checks to Unions with Index Signatures](https://github.com/microsoft/TypeScript/pull/34927)
* JSDocで `object` と `noImplicitAny` が併用された時 `any` にならないようにした [`object` in JSDoc is No Longer `any` Under `noImplicitAny`](https://github.com/microsoft/TypeScript/issues/34980)
    * 歴史的経緯が色々あった
    * `Object` と `object` はTypeScript的には区別される
    * JavaScriptでは両方 `any` になっていた
    * `object` は `noImplicitAny` と併用した場合、TypeScriptと同様の振る舞いになるようになる
* オプショナルな引数について型推論が行われない場合、 `any` になるようになった (変更されたPRわからず…)
    * これによって `noImplicitAny` に引っかかるようになった

## 型のみのimport, exportのサポート

`import type { foo } from "./foo";` や `export type { foo };` がサポートされました。

```exports.ts
export interface Options {
    name: string;
}
```

```imports.ts
import type { Options } from "./exports";

export type { Options };
```

こういう感じ。
コンパイルするとimports.tsは次のようになります。

```imports.js
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
```

全部消えましたね。

基本的にTypeScriptは各ファイルの内容を把握していますし、どの値が型としてのみ使われていて、どの値が本当に値としても使われているかも判別できます。
しかし、先の例のようなimport, exportのみのような場合だと一部の例で差し障りがあるようです。
例えば `Babel` 経由で動かしている場合、 `transpileModule` API経由で使っている場合などです。

これらのために、"importしたいのは型だけだよ"と構文上明示できる必要があります。
[flowtypeには5年くらい前からある](https://flow.org/blog/2015/02/18/Import-Types/)やつです。

型だけなので当然値として使おうとすると怒られる。

```ts
import type { Foo } from "./exports";

// 型のみなので当然継承できずに怒られる
// error TS1361: 'Foo' cannot be used as a value because it was imported using 'import type'.
export class Bar extends Foo {
    sayHello() {
        return "Hello, bar!";
    }
}
```

構文上の曖昧さをなくすため、現時点ではdefault importとnamed importを混在して使うことは禁止されています。

```ts
// typeのみの場合、default importとnamed importを両方同時にすることはできない
//   ↓ は現時点ではエラーになる type import は A と B 両方？ A のみ？ パット見意味が一意にならない
// import type A, { B } from "./exports";

// 分割すればOK
import type A from "./exports";
import type { B } from "./exports";
```

さらに、今で同様のimportの使い方を制御するためのオプションも追加されました。
`compilerOptions` に `importsNotUsedAsValues` が指定できるようになっています。
使えるのは `remove` （デフォルト、今までと同じ）、 `preserve`、 `error` です。

```ts
// importsNotUsedAsValues が remove の場合、消える
// importsNotUsedAsValues が preserve の場合、 require("./exports"); などが残る
// importsNotUsedAsValues が error の場合、怒られる
//   error TS1371: This import is never used as a value and must use 'import type' because the 'importsNotUsedAsValues' is set to 'error'.
import { A } from "./exports";
// ↓ならOK
// import type { A } from "./exports";

function foo(): typeof A {
    return "a";
}
```

`preserve` は `import "./exports";` 相当を書いたのと同様の結果が残ります。
importによるside-effectが欲しい場合ですね。

`error` の場合、エラーになります。
強火で行きたい人はこれにしましょう。

今後、この `import type` を積極的に使うべきかどうかは、プロジェクトの事情によるでしょう。
そのうちこれを検出して自動的に直してくれるlinterなりが出る気はするので、困ったらしっかり使う、くらいの感じでよいのではないでしょうか。
今のところimportの自動補完でも `type` 無しで行われるっぽいので手で直すのめんどくさいし…。

## ECMAScript Private Fieldsのサポート

出ましたね `#foo` とかいうキモいsyntaxのあいつです…。
利用する時は `WeakMap` をpolyfillとして使うので `target` が `es2015` 以上が必要です。

```ts
// error TS18028: Private identifiers are only available when targeting ECMAScript 2015 and higher.

class Person {
    #name: string
    private fav?: string;
    // 名前空間が分かれてるので同名のプロパティが存在してもよい
    name: string;

    // 他の可視性制御の修飾子と併用したらダメ
    // error TS18010: An accessibility modifier cannot be used with a private identifier.
    // public #other1?: string;
    // private #other2?: string;
    // protected #other3?: string;
    // これもダメ abstract class で定義しても次のエラーが出る
    // error TS18019: 'abstract' modifier cannot be used with a private identifier
    // abstract #other4?: string;

    // これはOKらしい
    readonly #other5?: string;

    constructor(name: string, fav?: string) {
        this.#name = name;
        this.name = name;
        this.fav = fav;
    }

    greet() {
        console.log(`Hello, my name is ${this.#name}!${this.fav ? ` I like ${this.fav}.` : ""}`);
    }
}

class Person2 extends Person {
    // 継承しても親から引き継がない 独立している
    #name: string;

    constructor(name: string) {
        super(name + "👍");
        this.#name = name;
    }

    greet2() {
        console.log(`Hi, I'm ${this.#name}.`);
    }
}

let papyrus = new Person("Papyrus", "🐾");
// Hello, my name is Jeremy Bearimy! I like 🐾. と表示される
papyrus.greet();

let sans = new Person2("Sans");
// Hello, my name is Sans👍! と表示される
sans.greet();
// Hi, I'm Sans. と表示される
sans.greet2();

// クラスの外側からはアクセスできない
// error TS18013: Property '#name' is not accessible outside class 'Person' because it has a private identifier.
// papyrus.#name

// 型上にはちゃんと存在しているので、該当のprivate fieldがないとエラーになる
// 実質、Person型かそれを継承したクラスのインスタンス縛りになる
// error TS2741: Property '#name' is missing in type '{ fav: string; name: string; greet(): void; }' but required in type 'Person'.
// let mimic1: Person = {
//     fav: "",
//     name: "",
//     greet() { },
// };

// これはOK
let mimic2: Person = sans;

// esnextで .js 出せばわかるけど #name 的なのは宣言必須
```

private fieldは完全に他の名前空間から独立し、そのクラス自身にのみ紐付きます。
クラス外部からアクセスすることはできません。
これには子クラスも含まれます！
なので、private fieldは完全に隠蔽され追加・変更・削除を自由にやっても他を壊す/壊されることがないので安心です。

transpileされたコードを見るとどうなっているかわかりやすいので、コンパイル後のJSコードを抜粋します。

```js
var __classPrivateFieldSet = (this && this.__classPrivateFieldSet) || function (receiver, privateMap, value) {
    if (!privateMap.has(receiver)) {
        throw new TypeError("attempted to set private field on non-instance");
    }
    privateMap.set(receiver, value);
    return value;
};
var __classPrivateFieldGet = (this && this.__classPrivateFieldGet) || function (receiver, privateMap) {
    if (!privateMap.has(receiver)) {
        throw new TypeError("attempted to get private field on non-instance");
    }
    return privateMap.get(receiver);
};
var _name;
class Person {
    constructor(name) {
        _name.set(this, void 0);
        __classPrivateFieldSet(this, _name, name);
    }
    greet() {
        console.log(`Hello, my name is ${__classPrivateFieldGet(this, _name)}!`);
    }
}
_name = new WeakMap();
```

クラスの種類毎にプロパティ保持用のWeakMapが作成されています。
このWeakMapは他のモジュールや子クラスなどの外部に公開されないことがわかります。

## `export * as foo ...` 的な構文のサポート

できるようになったらしいです。

```libs.ts
export function hello() {
    console.log("Hello!");
}

export function bye() {
    console.log("Bye...");
}
```

```index.ts
// 今までのやり方
import * as a from "./libs";
export { a };

// まぁ普通にアクセスできる
// a.hello();

// 上記のimport, exportを一度にやる！
export * as b from "./libs";

// b はこのモジュール中には公開されないっぽい
// b.hello();
```

```basic.ts
import { a, b } from "./";

a.hello();
a.bye();

b.hello();
b.bye();
```

この構文でexportした識別子は同モジュール中で参照できないようです。
トランスパイルされたJSを見るとprivateっぽい振る舞いになっていることが汲み取れます。

個人的には `export * from "./foo";` すら稀に使うなーくらいだったので使う機会あまりなさそう。

## トップレベルのawaitのサポート

割と嬉しいやつ〜〜〜
`target` が `es2017` 以降、 `module` が `esnext` か `system` じゃないとダメだそうです。

```ts
// target が es2017 以降で module が esnext か system じゃないとダメ
// es2017 以降だと await がキーワードになっているので

function timeout(milliseconds: number) {
    return new Promise(resolve => {
        setTimeout(resolve, milliseconds);
    })
}

console.log("Just 1 sec");
await timeout(1000);
console.log("dit it!");

// export とか import とかないとモジュールだと思ってもらえないやつ
//   --isolatedModules でもいいらしい(試してない)
// error TS1375: 'await' expressions are only allowed at the top level of a file when that file is a module, but this file has no imports or exports. Consider adding an empty 'export {}' to make this file a module.
export { }
```

bundlerがいい感じにしてくれないと生変換でそのまま使うのは時代を先取りすぎそう…。

## `es2020` が `target` と `module` に入った

らしいです。
optional chaining, nullish coalescing, `export * as ns`, dynamic import, 
bigint literalあたりが含まれているそうな。

## JSDocでの修飾子のサポート

らしいです。
前からなかったっけこれ…？
JavaScriptとか全然やらないのでわからん。

```js
// @ts-check

class Foo {
    constructor() {
        /** @private */
        this.stuff = 100;
    }

    printStuff() {
        console.log(this.stuff);
    }
}

// error TS2341: Property 'stuff' is private and only accessible within class 'Foo'.
new Foo().stuff;
```

`@public`, `@private`, `@protected`, `@readonly` が使えるようになったらしいですね。

## より良いディレクトリ監視と `watchOptions`の追加

tsconfig.json(とjsconfig.json)に `compilerOptions` と同じレベルに `watchOptions` が追加されました。
なんか今までも `TSC_WATCHFILE`, `TSC_WATCHDIRECTORY` というのが[あった](https://www.typescriptlang.org/docs/handbook/configuring-watch.html)らしいですね…。
まったく知らんかった…。[v2.8くらいからあった](https://github.com/microsoft/TypeScript/pull/21243)らしいです。

`node_modules` の更新時とかに激烈にCPU使ってしまったりするのを制御する目的っぽい。

```json
{
  // 既存の設定色々…

  "watchOptions": {
    // fixedPollingInterval, priorityPollingInterval, dynamicPriorityPolling, useFsEvents, useFsEventsOnParentDirectory
    // useFsEvents がデフォルトで今まで通り
    "watchFile": "useFsEvents",
    // fixedPollingInterval, dynamicPriorityPolling, useFsEvents
    // useFsEvents がデフォルトで今まで通り
    "watchDirectory": "useFsEvents",
    // useFsEvents を選択してた時に利用できなかったらどの戦略で行くか かな
    // fixedPollingInterval, priorityPollingInterval, dynamicPriorityPolling
    "fallbackPolling": "dynamicPriority",
    "synchronousWatchDirectory": false
  }
}
```

こんな感じ。
設定項目は4つある。
英語が正しく読めてるかちょっと自信ない…！

* `watchFile` 単体のファイルをどう監視するかの設定
    * `fixedPollingInterval` 一定間隔で毎秒何回かチェック
    * `priorityPollingInterval` 一定間隔でチェックするがあんま重要そうじゃなさそうなやつはちょっと頻度落とす
    * `dynamicPriorityPolling` 動的なキューに入れてチェック 微妙そうなヤツは微妙なりの頻度
    * `useFsEvents` デフォルト。OSネイティブのファイルシステムからのイベントで駆動
    * `useFsEventsOnParentDirectory` 対象ファイルを含むディレクトリの変更を利用する。負荷は下がるが精度も下がるかも
* `watchDirectory` 再帰的なディレクトリ監視ができるシステムでも振る舞いの設定
    * `fixedPollingInterval` 一定間隔で毎秒何回かチェック
    * `dynamicPriorityPolling` 動的なキューに入れてチェック 微妙そうなヤツは微妙なりの頻度
    * `useFsEvents` デフォルト。OSネイティブのファイルシステムからのイベントで駆動
* `fallbackPolling` ファイルシステムからのイベントで動かす設定かつ、それがサポートされない環境だった時の設定
    * `fixedPollingInterval` 前に同じ
    * `priorityPollingInterval` 前に同じ
    * `dynamicPriorityPolling` 前に同じ
* `synchronousWatchDirectory` ディレクトリの遅延監視を無効にする。ごく少数のセットアップでは役に立つかも。

って感じらしいです。
現状特に困ってないなら使わなくてよいんじゃないでしょうか…。

## ちょっと怠惰なインクリメンタルチェック

`assumeChangesOnlyAffectDirectDependencies` というオプションが導入されたそうです。
`--watch` モードに影響があるそうな…。

ファイルA.ts←ファイルB.ts←ファイルC.ts←.... 

というファイル間の依存関係があった場合、今まではファイルAを変更するとファイルB, Cという具合に参照先すべてをチェックしなおしていました。
`assumeChangesOnlyAffectDirectDependencies` を `true` に設定すると、ファイルAを変更した後、1つ先のファイルBまでしかチェックされないようになるらしいです。
多くの場合ではそれで十分な挙動だと思います…！

普段 `--watch` 相当の機能や `--incremental` とか全然使ってないのでわからん…！
もしかしたら陰ながら恩恵の預かれたりするのかな…。

## index signatureを含むunion型への値の代入のチェックをもうちょっと厳格にした

説明が難しいのでコード見て感じてください。

```ts
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
```

なんかこれは今までもエラーだったのでは…？という気がするけどそうでもないっぽい。

```
$ npx -p typescript@v3.7 tsc ./src/indexSignatureAndUnion/basic.ts
$ npx -p typescript@rc tsc ./src/indexSignatureAndUnion/basic.ts
```

v3.7.5 と v3.8.1 を比べるとたしかに前は素通しになっているのであった。

## JSDocで `object` と `noImplicitAny` が併用された時 `any` にならないようにした

らしいです。
歴史的経緯とか文化より、JSDocでは `Object` は `any` 相当に評価されていた。
これはまぁなんでもいいよの意図と捉えるのが自然なので…。

で、 `object` (先頭小文字) も `any` としてたんだけど、v3.8からTypeScriptの型の評価にあわせる感じにした、という経緯のようです。

## オプショナルな引数について型推論が行われない場合、 `any` になるようになった

らしいです。

```ts
function foo(f: () => void) {
}

// error TS7006: Parameter 'param' implicitly has an 'any' type.
foo((param?) => {
});
```

型アノテーションなしで `param?` だけ書くという発想がなかった…。
なるほどね？
