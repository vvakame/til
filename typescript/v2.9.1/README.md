こんにちは[メルペイ社](https://www.merpay.com/jp/)な[@vvakame](https://twitter.com/vvakame)です。

[TypeScript 2.9.1](https://blogs.msdn.microsoft.com/typescript/2018/05/31/announcing-typescript-2-9/)がアナウンスされました。

[What's new in TypeScriptも更新](https://github.com/Microsoft/TypeScript/wiki/What's-new-in-TypeScript#typescript-29)されているようです。
[破壊的変更](https://github.com/Microsoft/TypeScript/wiki/Breaking-Changes#typescript-29)もあるよ！

[この辺](https://github.com/vvakame/til/tree/master/typescript/v2.9.1)に僕が試した時のコードを投げてあります。

ちなみに、次のバージョンは2.10じゃなくて[3.0らしい](https://github.com/Microsoft/TypeScript/wiki/Roadmap#30-july-2018)です。

## 変更点まとめ

* ファイル名のリネームのサポート [Add 'renameFile' command to services](https://github.com/Microsoft/TypeScript/pull/23573)
    * ファイル名をリファクタリングできるようになった
* 選択範囲を別ファイルに切り出す操作のサポート [Add 'move to new file' refactor](https://github.com/Microsoft/TypeScript/pull/23726)
    * 定義を別ファイルに切り出すリファクタリングができるようになった
* 使ってない定義があったら教えてくれるようになった [Show unused declarations as suggestions](https://github.com/Microsoft/TypeScript/pull/22361)
    * `--noUnusedLocals` や `--noUnusedParameters` でエラーになる箇所について、これらを指定しない場合警告として表示してくれるようになった
* プロパティをgetter/setterに変換できるようになった [Convert property to getter/setter](https://github.com/Microsoft/TypeScript/pull/22143)
    * そのまんま
* 型として `import(...)` をどこでも使えるようになった [Allow `import(...)`-ing types at any location](https://github.com/Microsoft/TypeScript/issues/14844)
    * `let p: import("./foo").Person` とか書ける
    * ついでにimportした関数とかから作った要素をexportする時の縛りがゆるくなった
* `--pretty` がデフォルトで有効になった [`--pretty` error output by default](https://github.com/Microsoft/TypeScript/pull/23408)
    * `tsc --pretty` とかするとエラー表示が見やすくなるオプションがあった これがデフォルト有効に
    * `--pretty false` で今までと同じ出力
* `.json` をimportした時にいい感じの型が付くようになった [New `--resolveJsonModule`](https://github.com/Microsoft/TypeScript/pull/22167)
    * `--resolveJsonModule` が生えた
    * `moduleResolution` が `node` のとき、jsonをimportするとjsonの中身に応じた型が自動的に付く
* タグ付きテンプレートリテラルにgenericsの型を指定できるようになった [Support for passing generics to tagged template calls](https://github.com/Microsoft/TypeScript/pull/23430)
    * ```styledComponent<MyProps>`font-size: 1.5em;` ```こんな感じ
* keyofとかMapped typesでstring以外のnumberとかsymbolも出てくるようになった [Support number and symbol named properties with keyof and mapped types](https://github.com/Microsoft/TypeScript/pull/23592)
    * `keyof T` で `string | number | symbol` が得られるようになった
    * `--keyofStringsOnly` も用意された
* JSXのタグの中でComponentにgenericsの型を指定できるようになった [Support for passing generics to JSX elements](https://github.com/Microsoft/TypeScript/pull/22415)
    * `<MyComponent<number> data={12} />` こゆ感じ
* `import.meta` がサポートされた [Support for `import.meta`](https://github.com/Microsoft/TypeScript/pull/23327)
    * https://github.com/tc39/proposal-import-meta この仕様
    * `import.meta.__dirname` とか書けるようになりたい！的な
    * globalに `ImportMeta` を生やす
    * `--target esnext` と `--module esnext` が必要
* `.d.ts.map` が出力できるようになった [Declaration source maps and code navigation via them](https://github.com/Microsoft/TypeScript/pull/22658)
    * `--declarationMap` が導入された
    * `--declaration` と併用する
    * `.d.ts.map` が生成され、 `.d.ts` に飛んだ時にオリジナルの `.ts` コードを見ることができる
    * `--inlineSource` とかは反映されないのでつらい
* "提案"レベルの情報の追加 [Show suggestion diagnostics for open files](https://github.com/Microsoft/TypeScript/pull/22204)
    * error, warning の他に suggestion(info) が追加された
* require 使ってたら import に変換してくれる [Convert require to import in .ts files](https://github.com/Microsoft/TypeScript/pull/23711)
    * 1個前の仕組みを使っている
    * .ts 限定らしい
    * 時代を感じますねぇ（互換性があるという判断が行われるようになった）
* Quick Fixとかで使われるクォートを " と ' のどっちかに指定できるオプション [Support setting quote style in quick fixes and refactorings](https://github.com/Microsoft/TypeScript/pull/22236)
    * とか言ってるけど実は(Language) Serviceに対して設定を追加できるAPIが生えた
        * 裏で UserPreferences というのが生えた
    * 今できるのはクォートの統一とモジュールをimportする時に相対パスにするか絶対パスにするかの設定ができる
        * `typescript.preferences.quoteStyle`
        * `typescript.preferences.importModuleSpecifier`
* Node.jsのビルトインモジュールを使おうとしたら `@types/node` 入れてくれる [Install `@types/node` for built-in node modules](https://github.com/Microsoft/TypeScript/pull/23807)
    * まぁあったほうが便利ですよね
* `strictNullChecks` で `object` を底にしてない型パラメータの値を `object` に割り振り不可にした [Unconstrained type variable not assignable to 'object'](https://github.com/Microsoft/TypeScript/pull/23806)
    * 一貫性を破壊できたのが修正された
* `never` は `for-in` とか `for-of` とかで繰り返し処理しようとした時にエラーになるようになった [Don't allow to iterate over 'never'](https://github.com/Microsoft/TypeScript/pull/22964)
    * そのまんま

### 破壊的変更！

上記リストのうち、破壊的変更を伴うのは次のものになります。

* keyofとかMapped typesでstring以外のnumberとかsymbolも出てくるようになった
* `--pretty` がデフォルトで有効になった
* 関数の引数に可変長引数がある場合、末尾のカンマを許さないように変更
* `strictNullChecks` で `object` に `object` を底にしてない型パラメータの値を割り振り不可にした
* `never` は `for-in` とか `for-of` とかで繰り返し処理しようとした時にエラーになるようになった

## ファイル名のリネームのサポート

そのまんまです。
VSCodeのファイルエクスプローラでファイル名を変更しようとすると、変更しようとしたファイルを参照しているコードのimportを書き換える？と聞かれます。
`No, never update imports`, `Yes, always update imports`, `No`, `Yes` の選択肢が現れます。
Yesを選ぶと自動的にimport句が更新されます。

仕様上、ファイルをrenameした後にダイアログを出しているようで、renameをキャンセルすることは（現時点では）できないようです。

`Yes, always update imports` を選ぶとUser Settingsに次の設定が追加されました。

```
{
    "typescript.updateImportsOnFileMove.enabled": "always"
}
```

大変便利なので常時Yesでいいんじゃないでしょうか。

## 選択範囲を別ファイルに切り出す操作のサポート

そのまんまです。
次のようなファイルがあるとします。

```ts
interface Foo {
    name: string;
}

let f: Foo = { name: "test" };

export { }
```

1〜3行目を選択して、Quick Fixで `Move to a new file` を選ぶと該当行を別ファイルに切り出してくれます。

```ts
import { Foo } from "./Foo";

let f: Foo = { name: "test" };
```

```ts
export interface Foo {
    name: string;
}
```

こんな感じですね。
ちゃんとexportを付与してくれたりimport句も作ってくれたりして偉いです。

元のファイルで `export {}` 的なESモジュールであると示しておかないと意図通りに動いてくれないので少しだけ注意が必要です。

## 使ってない定義があったら教えてくれるようになった

`--noUnusedLocals` や `--noUnusedParameters` というのがすでにあるんですが、これを使っていない時でも、エディタ上で使ってない変数とかをグレーアウトしてほんのりと使ってないことを教えてくれるようになった、という奴です。
JetBrains系のIDEではこういうほんのりと教えてくれる奴が充実していて役に立っているので嬉しいですね。
まぁ僕は上記オプションをONにしてコンパイルエラーにしますが…。

## プロパティをgetter/setterに変換できるようになった

そのまんまですね。
クラスのプロパティ名を選択してQuick Fixでgetter/setterに変換できるようになりました。

次のコードのaとbを変換してみます。

```ts
class Hoge {
    a?: string;
    b: string | undefined;
}
```

すると、こうなります。

```ts
class Hoge {
    private _a?: string;
    public get a(): string { // ← string | undefined じゃないのでコンパイルエラーになる
        return this._a;
    }
    public set a(value: string) {
        this._a = value;
    }
    private _b: string | undefined;
    public get b(): string | undefined {
        return this._b;
    }
    public set b(value: string | undefined) {
        this._b = value;
    }
}
```

ちょっとおバカですね…。
でも実用上の役にはしっかり立ちそうです。

## 型として `import(...)` をどこでも使えるようになった

ECMAScriptにdynamic importという仕様があります。
それと同じsyntaxで型注釈が書けるようになった…というものです。

```ts
export function hello(name: string) {
    return `Hello, ${name}!`;
}

export interface Data {
    id: string;
    content: any;
}
```

```ts
type fooType = typeof import("./foo");

function f(fn: typeof import("./foo").hello) {
    fn("import()");
}

const fn1: typeof import("./foo").hello = name => `Hi, ${name}`;
fn1("import()");

const fn2: fooType["hello"] = name => `Hi, ${name}`;
fn2("import()");


const data1: import("./foo").Data = {
    id: "foo",
    content: "bar",
};

type Data = import("./foo").Data;
interface Data1 extends Data { }
// この書き方はinterfaceのsyntax的にダメ
// interface Data2 extends import("./foo").Data { }
```

いろいろな書き方ができますね。

TypeScriptでは `import ... from "...";` で型をimportしても、値として利用しなければコンパイル後には消えてしまいます。
よって、この書き方ができなくても実用上問題はなかった気もしますが、わかりにくい仕様だったことに変わりはありません。
ここで、 `import(...)` と同様の書き方で型としても参照できるようになり、わかりやすさは改善されたと言えるでしょう。

使い所がいまいちわからない…。

また、この記法の導入により型定義を生成するときの制限が緩和されました。

TypeScript 2.8までは次のようなコードをコンパイルするとエラーになっていました。

```ts
import { createHash } from "crypto";

export const hash = createHash("sha256");
```

```
src/importAsTypes/relax.ts:3:14 - error TS4023: Exported variable 'hash' has or is using name 'Hash' from external module "crypto" but cannot be named.

3 export const hash = createHash("sha256");
               ~~~~
```

これのコンパイルを通すには `import { createHash, Hash } from "crypto";` と書く必要がありました。

これは、 `.d.ts` ファイルが次のように出力されていたためです。

```ts
/// <reference types="node" />
import { Hash } from "crypto";
export declare const hash: Hash;
```

TypeScript 2.9からは `Hash` をimportしていなくてもコンパイルが通ります。
そして、次のような `.d.ts` ファイルが生成されます。

```ts
/// <reference path="../../node_modules/@types/node/index.d.ts" />
/// <reference types="node" />
export declare const hash: import("crypto").Hash;
```

…1行目はいるんですかね…？

## `--pretty` がデフォルトで有効になった

前からあった `--pretty` というエラーをpretty printしてくれるオプションがありました。
今回からこれがデフォルトになり、2.8までと同様の出力にしたい場合、 `--pretty false` とします。

```ts
src/importAsTypes/index.ts:21:25 - error TS2499: An interface can only extend an identifier/qualified-name with optional type arguments.

21 interface Data2 extends import("./foo").Data { }
                           ~~~~~~~~~~~~~~~~~~~~
```

## `.json` をimportした時にいい感じの型が付くようになった

まんまです。

```ts
import data from "../../tsconfig.json";

// OK
console.log(data.compilerOptions.target);
// こっちは存在しないのでエラーになる
console.log(data.notExists);
```

なお、moduleResolutionをclassicにするとこの機能は使えません。

## タグ付きテンプレートリテラルにgenericsの型を指定できるようになった

タグ付きテンプレート用の関数にgenericsを使えるようになりました。
主にReactというかstyled-componentにいい感じに型をつけるために導入されたっぽいですね。

```ts
function tag<T>(strs: TemplateStringsArray, ...args: T[]) {
    console.log(strs, ...args);
}

// パラメタが全部 number なのでOK
tag<number>`hoge${1}${2}`;
// 2つ目のパラメタが number なのでエラー
// tag<string>`hoge${"A"}${2}`;
// 1つ目と2つ目のパラメタが一致してないのでエラー
// tag`${true}${1}`;
```

styled-componentだとこういう感じらしいです。

```ts
declare function styledComponent<Props>(strs: TemplateStringsArray): Component<Props>;

styledComponent<MyProps>`
    font-size: 1.5em;
`

interface Component<T> {
}

interface MyProps {
    name: string;
    age: number;
}
```

よかったですね。

## keyofとかMapped typesでstring以外のnumberとかsymbolも出てくるようになった

今までは string なプロパティしかkeyとして切り出せなかったのが、numberやsymbolもkeyとして取れるようになりました。
`keyof T` で `string | number | symbol` 的なものが取れる可能性がある、ということですね。

```ts
const sym = Symbol();

interface Foo {
    str: string;
    11: number;
    [sym]: symbol;
}

// typeof sym | "str" | 11 となる
type keyofFoo = keyof Foo;

// typeof sym
type A = Extract<keyofFoo, symbol>;
// str
type B = Extract<keyofFoo, string>;
// 11
type C = Extract<keyofFoo, number>;

// なお --keyofStringsOnly だと keyofFoo は "str" になる

// ちなみに…
type array = [1, 2, 3];
// "1", "2", "3", ...Arrayのメソッドいろいろ
type keyofArray = keyof array;
```

旧来の挙動にするために `--keyofStringsOnly` が用意されました。
部分的にやったいきたい場合、 `keyof Foo` を `Extract<keyof Foo, string>` に置き換えていくのがよい、とのことです。

## JSXのタグの中でComponentにgenericsの型を指定できるようになった

Reactのやり込みが足りないので活用ポイントがよくわかんないですね…。
こんな感じだそうです。

```tsx
import { Component } from "react";

interface Props<T> {
    data: T;
}

class MyComponent<T> extends Component<Props<T>> {
    render() {
        return <div>{this.props.data}</div>;
    }
}

<MyComponent<number> data={12} />
```

パーサさんも辛い気持ちになったりしているんだろうなぁ。

## `import.meta` がサポートされた

[import.meta](https://tc39.github.io/proposal-import-meta/)がサポートされました。
利用には `--target esnext` と `--module esnext` が必要です。

```ts
declare global {
    interface ImportMeta {
        foo: string;
    }
}

import.meta.foo;
// これは定義されてないのでエラー
// import.meta.notExist;
```

そういえばTypeScriptで[mjsが出力できないのはちょっと不便](https://github.com/Microsoft/TypeScript/issues/18442)ですね。

## `.d.ts.map` が出力できるようになった

.d.ts.map は .d.ts と .ts ファイルの対応関係を記すものです。
.js.map が .js と .ts の関係を記すのと同様ですね。

コードを書いているとき、定義にジャンプすると今まで .d.ts の定義に飛んでいたのが、ちゃんと対応関係が定義されている場合、 .ts ファイルの該当部分にジャンプできます。

有効にするには `--declarationMap` を使います。
`--declaration` との併用が必須です。
なお、 `--inlineSource` を使っても .d.ts.map に .ts の内容は現時点では保存されません。かなしい。

[.tsと.d.tsファイルが同じ場所にある場合、.tsの読み込みが優先](https://github.com/Microsoft/TypeScript/issues/4667)されます。
ですので、npm publishするときは .ts ファイルをignoreするか、 `--outDir` で別ディレクトリに出力しているかのどちらかだと思います。
もしあなたが後者なら、このオプションを使いやすく、ユーザの役に立つでしょう。
ちなみに僕は前者スタイルです(かなしい)

## "提案"レベルの情報の追加

Error, Warning の他に Suggestion が追加されました。

現時点では、次の項目である `require 使ってたら import に変換してくれる` が実装されています。

Language Service周りの開発に興味がある人は[PR](https://github.com/Microsoft/TypeScript/pull/22204)をよく読むとよさそうです。

## require 使ってたら import に変換してくれる

現時点では、 CommonJSモジュールをESモジュールに変換する？という提案が実装されています。

```ts
const fs = require("fs");
```

↑を↓に変換してくれます

```ts
import fs from "fs";
```

こういう変換を行ってもよい時代になった、とTypeScript teamが判断したこと自体が、趣深いですね。

## Quick Fixとかで使われるクォートを " と ' のどっちかに指定できるオプション

裏で(Language) Serviceに対して設定を追加できるAPIが生えました。
UserPreferences というのが生えてます。

とりあえず次の設定をVSCodeに対して行うとLanguage Serviceに設定が伝搬されます。
新規にQuick Fixやリファクタリングで生成された文字列リテラルがダブルクォートかシングルクォートかを指定できるというものです。

```json
{
    "typescript.preferences.quoteStyle": "double"
}
```

`typescript.preferences.importModuleSpecifier` というのも追加され、 `non-relative` と `relative` が選べますが `baseUrl` の設定などに左右されてめんどくさいので割愛します。

## Node.jsのビルトインモジュールを使おうとしたら `@types/node` 入れてくれる

そのまんまです。
今までも `.d.ts` が存在しないパッケージを `@types/` 無しに使っていると、Quick Fixで `@types/` を入れてくれていました。
この機能は[Node.jsのモジュール](https://github.com/Microsoft/TypeScript/blob/e53e56cf8212e45d0ebdd6affe462d161c7e0dc5/src/services/jsTyping.ts#L34)を `@types/node` 無しに使うと、 これを入れるか聞いてくれるというものです。

よくあるパターンなのであると嬉しいですね。

## `strictNullChecks` で `object` を底にしてない型パラメータの値を `object` に割り振り不可にした

[`object` 型](https://qiita.com/vvakame/items/eb6c054360868b88f9b1#object%E5%9E%8B%E3%81%AE%E3%82%B5%E3%83%9D%E3%83%BC%E3%83%88)に不正にstringとかを代入できるバグがあったので塞がれたというやつです。

詳細は次のコードの通りです。

```ts
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
```

## `never` は `for-in` とか `for-of` とかで繰り返し処理しようとした時にエラーになるようになった

次のようなコードが両方エラーになるようになりました

```ts
const neverVar: never = (() => { throw new Error() })();

// error TS2407: The right-hand side of a 'for...in' statement must be of type 'any', an object type or a type parameter, but here has type 'never'.
for (let v in neverVar) {
    console.log(v);
}

// error TS2488: Type 'never' must have a '[Symbol.iterator]()' method that returns an iterator.
for (let v of neverVar) {
    console.log(v);
}
```
