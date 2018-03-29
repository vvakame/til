# TypeScript 2.8.1 変更点

こんにちは[ソウゾウ社](https://www.souzoh.com/jp/)な[@vvakame](https://twitter.com/vvakame)です。

[TypeScript 2.8.1](https://blogs.msdn.microsoft.com/typescript/2018/03/27/announcing-typescript-2-8/)がアナウンスされました。

[What's new in TypeScriptも更新](https://github.com/Microsoft/TypeScript/wiki/What's-new-in-TypeScript#typescript-28)されているようです。
[破壊的変更](https://github.com/Microsoft/TypeScript/wiki/Breaking-Changes#typescript-28)もあるよ！

今回から[この辺](https://github.com/vvakame/til/tree/master/typescript/v2.8.1)に僕が試した時のコードを投げておくことにしました。

## 変更点まとめ

* Conditional types（条件付きの型） [Conditional types](https://github.com/Microsoft/TypeScript/pull/21316)
    * `T extends U ? X : Y` みたいなの書ける
* Conditional typesでの型推論 [Type inference in conditional types](https://github.com/Microsoft/TypeScript/pull/21496)
    * `infer` の導入
    * ある型に式を当てはめた結果、得られた型推論の型を新たな型パラメタとして利用可能
    * `type ReturnType<T> = T extends (...args: any[]) => infer R ? R : T;` こういう
        * `R` が新しく導出された型 この場合、関数の返り値の型が取れる
* ビルトインのConditional typesを使った型の追加 [Predefined conditional types](https://github.com/Microsoft/TypeScript/pull/21847)
    * `Exclude<T, U>`, `Extract<T, U>`, `NonNullable<T>`, `ReturnType<T>`, `InstanceType<T>` の追加
* `.d.ts` だけ出力する `--emitDeclarationOnly` オプションの追加 [--emitDeclarationOnly flag to enable declarations only output](https://github.com/Microsoft/TypeScript/pull/20735)
    * そのまんま
* `@jsx` プラグマコメントのサポート [Add support for transpiling per-file jsx pragmas](https://github.com/Microsoft/TypeScript/pull/21218)
    * `/** @jsx dom */` とかやると `React.createElement` の代わりに `dom` が使われるようになる
    * 今までプロジェクト全体で変えることはできたけどファイル個別にはできなかった
        * ReactとPreactとかを混ぜて使えるようになった
* JSX名前空間の探索にファクトリ関数の名前空間が使われるようになった [Lookup JSX namespace within factory function](https://github.com/Microsoft/TypeScript/pull/22207)
    * `React.createElement` をファクトリ関数とするなら、まず `React.JSX` にある定義を探して使ってくれるようになった
    * 実はJSXの型チェックは[カスタマイズ可能](http://www.typescriptlang.org/docs/handbook/jsx.html#type-checking)なんだけど、複数のJSXを使うライブラリを混在させた時にコンパイルが通せるように頑張った形
* [Mapped types](https://qiita.com/vvakame/items/fc7e37d0296c63f41f4f#%E3%81%82%E3%82%8B%E5%9E%8B%E3%81%AE%E3%83%95%E3%82%A3%E3%83%BC%E3%83%AB%E3%83%89%E3%81%AE%E4%BF%AE%E9%A3%BE%E5%AD%90%E3%81%AE%E5%A4%89%E6%8F%9Bmap%E5%87%A6%E7%90%86%E3%81%8C%E5%8F%AF%E8%83%BD%E3%81%AB)で `-readonly` とか `-?` で修飾子の引き剥がしが可能に
    * そのまんま
    * `+readonly` とか `+?` も記法としては書ける(前から+無しでできたやつ)
* importの整理ができるようになった [Introduce an organizeImports command](https://github.com/Microsoft/TypeScript/pull/21909)
    * エディタからやるやつです
    * 今まではtslintでやってた人も多いかも？
* 初期化忘れのフィールドを雑に修正できるようになった [add support of codefix for Strict Class Initialization](https://github.com/Microsoft/TypeScript/pull/21528)
    * `| undefined` を追加したり `!` を追加したりゼロ値を初期値としてセットしたり
    * 地味に便利そう
* `keyof` を intersection type に適用した時個別に簡約されるようになった [Distribute 'keyof' intersection types](https://github.com/Microsoft/TypeScript/pull/22300)
    * `type T1 = keyof (A & B)` こんなんがちゃんと展開される
* [More special declaration types in JS](https://github.com/Microsoft/TypeScript/pull/21974)
    * なんか名前空間の定義の解釈できるバリエーションが増えたらしい
    * JSにはあまり興味はぬい…！
* `--noUnusedParameters` で使ってない型パラメータも怒られるようになった [Unused type parameters should be checked by --noUnusedParameters, not --noUnusedLocals](https://github.com/Microsoft/TypeScript/pull/21167)
    * 前は `--noUnusedLocals` を使ってた時に報告されてたけど違うでしょ的な
* HTMLObjectElementがalt属性持ってたのを直した [HTMLObjectElement has alt attribute only in Internet Explorer (and edge?)](https://github.com/Microsoft/TypeScript/issues/21386)
    * あわせて、僕がかなり前に[なんかやった](https://github.com/Microsoft/TSJS-lib-generator/pull/176)のもひっそりとshipされた

## Conditional types（条件付きの型）

三項演算子のように、ある条件を満たす時は型X、そうでなければ型Yを表す。
ということが `T extends U ? X : Y` 的に書けるようになりました。

もちろん、構文自体は三項演算子のようですが、話題としては型なので自由な式が書けるわけではありません。
今のところ、`T extends U` 形式しか使えないようです。
`T` を `U` に代入可能であれば真側、不可能であれば偽側の型と評価されます。

```ts
// T の型に応じて文字列のリテラル型に変換
type TypeName<T> =
    T extends string ? "string" :
    T extends number ? "number" :
    T extends boolean ? "boolean" :
    T extends undefined ? "undefined" :
    T extends Function ? "function" :
    "object";

// string と互換性のある型をTに指定 → string
const a1: TypeName<string> = "string";
const a2: TypeName<"a"> = "string";

// 同様にそれぞれ互換性のある型に落ち着く
const b: TypeName<true> = "boolean";
const c: TypeName<undefined> = "undefined";
const d: TypeName<() => void> = "function";
const e: TypeName<Date> = "object";
```

実際の利用例として、オーバーロードの置き換えなどにも利用できます。
公式ブログのコードを使って試してみると、今までのやり方よりもより良い結果が得られます。

```ts
interface Id { id: number, /* other fields */ }
interface Name { name: string, /* other fields */ }

// 今までのやり方
declare function createLabelA(id: number): Id;
declare function createLabelA(name: string): Name;
declare function createLabelA(name: string | number): Id | Name;

// Conditional typesを使ったやり方
declare function createLabelB<T extends number | string>(idOrName: T):
    T extends number ? Id : Name;

// 今までのやり方だと…
let a1 = createLabelA("typescript");   // Name
let b1 = createLabelA(2.8);            // Id
let c1 = createLabelA("" as any);      // Id ← 最初にマッチしたものが採用されてしまう
let d1 = createLabelA("" as never);    // Id ← 最初にマッチしたものが採用されてしまう

// Conditional typesだと…
let a2 = createLabelB("typescript");   // Name
let b2 = createLabelB(2.8);            // Id
let c2 = createLabelB("" as any);      // Id | Name ← 偉い！
let d2 = createLabelB("" as never);    // never     ← 偉い！
```

lib.d.tsでも[keyofを使ったoverloadの最適化](https://github.com/Microsoft/TypeScript/blob/66bf5b4e9d2796dd03f5a5b65a9f5b51f3b282a8/lib/lib.d.ts#L7471-L7472)もありますし、宣言的に型を記述できる幅が増えていきますね。

その他、[アンダース・ヘルスバーグ御大がPRに書いた概要](https://github.com/Microsoft/TypeScript/pull/21316)を見ると様々な応用例が考えられています。
`FunctionProperties` や `DeepReadonly` の実装や発想は参考になります…！

### Conditional typesでの型推論

Conditional typesは先に紹介した以外にも協力な機能を持っています。
例えば、`infer` の導入です。
`infer` は既存の型パラメタ以外の、型推論から導出される型を新しい型パラメタとして利用できます。
説明がめっちゃ難しいですが、次の例を見てください。

```ts
// infer R で新しい型パラメタを推論結果から導入する！つよい！
type ReturnType<T> = T extends (...args: any[]) => infer R ? R : T;

// 1. () => number は (...args: any[]) => any に代入可能である
// 2. infer R の位置が number である
// よって TmpType == number ！
type TmpType = ReturnType<() => number>;

// number[] と評価される
type ConcatReturn = ReturnType<Array<number>["concat"]>;
// 関数ではないのでDateと評価される（個人的にはneverになってほしい気がする）
type DateType = ReturnType<Date>;

// TypeScript 2.7以前で使ってたハック
const tmp = (false as true) && [1].concat([2]);
type ConcatReturnOld = typeof tmp;

// 残念ながらRの型の範囲は指定できない… っぽい？たぶん
// type ReturnType2<T, R extends string> = T extends (...args: any[]) => infer (R extends string) ? R : T;

// 型が書ける場所だったらだいたいイケるっぽい…
{
    type SecondArg<T> = T extends (a: string, b: infer R) => void ? R : never;
    // boolean！ ばばーん
    type TmpType = SecondArg<(a: string, b: boolean) => void>;
}
{
    type FooParamType<T> = T extends { Foo: infer R; } ? R : never;
    // RegExp！ ばばーん
    type TmpType1 = FooParamType<{ Foo: RegExp }>;
    // never！ ぼぼーん
    type TmpType2 = FooParamType<{ Bar: RegExp }>;
}
{
    type InstanceType<T> = T extends { new(...args: any[]): infer R; } ? R : never;
    // Date！ びびーん
    type DateInstance = InstanceType<typeof Date>;
}
```

これにより、Flatten的な型も直感的に書けるようになっています。


```ts
{
    // infer 使わないやつ
    type Flatten<T> = T extends any[] ? T[number] : T;

    // number
    type a = Flatten<number>;
    // number
    type b = Flatten<number[]>;
    // number[]
    type c = Flatten<number[][]>;
}
{
    // T[number] パット見て一瞬でわからない… 次も同じ意味
    type Flatten<T> = T extends any[] ? T[0] : T;
}
{
    // infer を使うとわりとわかりやすい！
    type Flatten<T> = T extends Array<infer U> ? U : T;
}
{
    // 残念ながら再帰的に自分を参照することはできない
    // type Flatten<T> = T extends Array<infer U> ? Flatten<U> : T;
    // type a = Flatten<number[][][][][]>;
}
```

[型パラメタの値としてinferを使えるようにするPR](https://github.com/Microsoft/TypeScript/pull/22368)もあるようです。

### Conditional typesはUnion typesを分配して処理する

はい。
Conditional typesでは、基本的にUnion typesを指定した場合、含まれる個別の型それぞれに計算を適用します。
コード例を見たほうがわかりやすいと思うのでどうぞ…。

```ts
{
    type ToArray<T> = T extends any ? T[] : never;

    // これは次の形に展開される
    //   (string extends any ? string[] : never) |
    //   (number extends any ? number[] : never)
    // つまり、簡約すると string[] | number[]
    // union typesは個別にconditional typesの適用を受けるのだ
    type Tmp = ToArray<string | number>;

    // 括弧でくくっても string[] | number[] にされる(げせぬ)
    type TmpB = ToArray<(string | number)>;
}
{
    // T を [T] で展開しない！
    // また紛らわしい記法が…
    type ToArray<T> = [T] extends any ? T[] : never;

    // (string | number)[] と評価される
    type Tmp = ToArray<string | number>;
}
{
    // なお、neverが混ざると最終的に消えます
    type Foo<T> = T extends any ? T : never;
    // never は消えて string | number になる
    type Tmp = Foo<string | number | never>;
}
{
    // any も [any] にしても結果は変わらない
    // なぜなのか…
    type ToArray<T> = [T] extends [any] ? T[] : never;

    // (string | number)[] と評価される
    type Tmp = ToArray<string | number>;
}
{
    // もうおじちゃんの脳は型注釈を正しくパースできなくなってきましたよ…
    const array: [number] = [1];
}
```

### ビルトインのConditional typesを使った型の追加

Mapped typesの時と同じように、ビルトインのConditional typesを使った型が追加されています。
`Exclude<T, U>`, `Extract<T, U>`, `NonNullable<T>`, `ReturnType<T>`, `InstanceType<T>` です。

まずは、lib.d.tsから抜粋した定義をペッとしておきます。
前述のUnion typesに対して分配処理するやつを上手く使ってる定義もあるので頑張ってついていきましょう。

```ts
/**
 *  T（union types）から、Uで指定した型を除外したものを返す
 */
type Exclude<T, U> = T extends U ? never : T;

/**
 * T（union types）から、Uで指定した型に含まれるもののみを返す
 */
type Extract<T, U> = T extends U ? T : never;

/**
 * T（union types）から、null, undefinedを除外して返す
 * Exclude<T, null | undefined> と一緒（のはず）
 */
type NonNullable<T> = T extends null | undefined ? never : T;

/**
 * 関数の返り値の型を切り出す
 */
type ReturnType<T extends (...args: any[]) => any> = T extends (...args: any[]) => infer R ? R : any;

/**
 * コンストラクタを持つ型からインスタンスの型を切り出す
 */
type InstanceType<T extends new (...args: any[]) => any> = T extends new (...args: any[]) => infer R ? R : any;
```

利用例を見ていきます。

```ts
{
    type A = string | number | boolean | Date | null;

    // string | number | boolean | null と評価される
    type TmpA = Exclude<A, Date | RegExp>;
    // Date と評価される
    type TmpB = Extract<A, Date | RegExp>;
    // string | number | boolean | Date と評価される
    type TmpC = NonNullable<A>;
}
{
    // number と評価される
    type TmpA = ReturnType<() => number>;
    // string と評価される
    type TmpB = ReturnType<{ m(): string; }["m"]>;

    // Date と評価される
    type TmpC = InstanceType<typeof Date>;
    // RegExp と評価される
    type TmpD = InstanceType<{ new(): RegExp; }>;
}

// 気合で応用していくぞ…！
{
    class Clazz {
        str?: string;
        func?: (num: number) => number;
        method() { }
    }

    // Clazzの各要素の集合を取得して…
    type A = Clazz[keyof Clazz];
    // Functionであるもの（つまり関数とメソッド）だけを切り出す
    // ((num: number) => number) | (() => void) と評価される
    type B = Extract<A, Function>;
}
{
    // Pickの逆 指定したプロパティを含まない型を返す
    type Flip<T, K extends keyof T> = {
        [P in Exclude<keyof T, K>]: T[P];
    }
    // 指定したプロパティだけNonNullableを適用して返す
    // -? という記法は後述
    type PickWithNonNullable<T, K extends keyof T> = {
        [P in K]-?: NonNullable<T[P]>;
    };
    // 上記2つを組合せて指定したプロパティだけNonNullableにして返す
    type PropertyNonNullable<T, K extends keyof T> = PickWithNonNullable<T, K> & Flip<T, K>;

    class Clazz {
        str?: string;
        func?: (num: number) => number;
        method() { }
    }

    // { func: (num: number) => number; } & { str: string | undefined; method: () => void; }
    type A = PropertyNonNullable<Clazz, "func">;
    // { str: string; } & { func: (num: number) => number; method: () => void; }
    type B = PropertyNonNullable<A, "str">;

    const obj = new Clazz();

    obj.func = num => num;
    obj.func(1);
    const a: A = obj as A; // control flow的にはfuncはundefinedではないとわかっているがキャスト必要

    obj.str = "";
    obj.str.charAt(0);
    const b: B = a as B; // 同上
}
{
    // PRのdescriptionに書いてあったこっちのほうがまだしも有用っぽい
    type FunctionPropertyNames<T> = { [K in keyof T]: T[K] extends Function ? K : never }[keyof T];
    type FunctionProperties<T> = Pick<T, FunctionPropertyNames<T>>;
    type NonFunctionPropertyNames<T> = { [K in keyof T]: T[K] extends Function ? never : K }[keyof T];
    type NonFunctionProperties<T> = Pick<T, NonFunctionPropertyNames<T>>;
}
```

どこで使うんだこれはァ…！みたいな感じがしますね。
こんなん何も見ずにスラスラ書けるようになんないよ〜〜〜。
後述しますが、あるオブジェクトの全プロパティからundefineとnullを引剥したい時は `Required` を使うとよいです。

## `.d.ts` だけ出力する `--emitDeclarationOnly` オプションの追加

読んで字のごとく。
`tsc --emitDeclarationOnly` とかやると `.d.ts` だけ出力されます。
`tsc --outDir typings --emitDeclarationOnly` とかやると任意のディレクトリに型定義ファイルだけ出力されて便利そうです。

## `@jsx` プラグマコメントのサポート

`@jsx` を使って、そのファイル固有の[`--jsxFactory`](https://www.typescriptlang.org/docs/handbook/compiler-options.html)相当のものを指定できるようになりました。
これは、複数のJSXライブラリを混在させて使いたい時に便利です。

```ts
/** @jsx React.element作るマン */
import React from "react";

export function render() {
    return <span></span>;
}
```

とかすると

```ts
/** @jsx React.element作るマン */
import React from "react";
export function render() {
    return React.element作るマン("span", null);
}
```

となる。
もちろん他の `.tsx` の生成結果はデフォルトの `React.createElement` のままです。
今の所、 `React.element作るマン` のような存在しない要素を指定してもエラーにならないようなので注意が必要です。

## JSX名前空間の探索にファクトリ関数の名前空間が使われるようになった

JSXを使った時の型チェックのルールは[カスタマイズ可能](http://www.typescriptlang.org/docs/handbook/jsx.html#type-checking)です。
`@jsx` の導入により、複数のJSXを使ったライブラリが混在する可能性があります。
これに対処するため、JSXの型チェックルールも混在できるようにしよう、という仕様です。

今までは、`JSX`名前空間はglobalに置かなければいませんでした。
e.g. [Reactの場合](https://github.com/DefinitelyTyped/DefinitelyTyped/blob/1ad206c6160b2fc5d208e7ffe1856aad03b60802/types/react/index.d.ts#L3746)

```ts
declare global {
    namespace JSX {
        // 色々…
    }
}
```

これを、各ライブラリないし、各ファクトリ関数毎に別々のJSXの定義が使えます。
あまり直感的ではないのですが、importの書き方によって処理が分かれます。

まずは定義を見てみましょう。
`namespace JSX` が2つ出てきますが、それぞれ モジュールのルート と ファクトリ関数と同じ場所 に配置されています。

```ts
export function myjsx(): any;

// /** @jsx myjsx */ と import { myjsx } from "./lib"; の場合 こっちが使われる
export namespace myjsx {
    namespace JSX {
        interface IntrinsicElements {
            [e: string]: {};
        }
        interface Element {
            __myjsxBrandA: void; // ここが違う
            children: Element[];
            props: {};
        }
        interface ElementAttributesProperty { props: any; }
        interface ElementChildrenAttribute { children: any; }
    }
}

// /** @jsx Lib.myjsx */ と import Lib from "./lib"; の場合 こっちが使われる
export namespace JSX {
    interface IntrinsicElements {
        [e: string]: {};
    }
    interface Element {
        __myjsxBrandB: void; // ここが違う
        children: Element[];
        props: {};
    }
    interface ElementAttributesProperty { props: any; }
    interface ElementChildrenAttribute { children: any; }
}
```

どちらのJSX定義が利用されるかは、importの仕方によって違います。

```ts
/** @jsx myjsx */
import { myjsx } from "./lib";

export function render() {
    return <span></span>;
}

// "__myjsxBrandA" | "children" | "props" と評価される
type El = keyof ReturnType<typeof render>;
```

```ts
/** @jsx Lib.myjsx */
import Lib from "./lib";

export function render() {
    return <span></span>;
}

// "__myjsxBrandB" | "children" | "props" と評価される
type El = keyof ReturnType<typeof render>;
```

ちょっとややこしいですね。

`export import JSX = myjsx.JSX;` とかやってごまかせないか試してみましたが、意図通りに動かなかったです。

## Mapped typesで `-readonly` とか `-?` で修飾子の引き剥がしが可能に

[Mapped typesという複雑な奴](https://qiita.com/vvakame/items/fc7e37d0296c63f41f4f#%E3%81%82%E3%82%8B%E5%9E%8B%E3%81%AE%E3%83%95%E3%82%A3%E3%83%BC%E3%83%AB%E3%83%89%E3%81%AE%E4%BF%AE%E9%A3%BE%E5%AD%90%E3%81%AE%E5%A4%89%E6%8F%9Bmap%E5%87%A6%E7%90%86%E3%81%8C%E5%8F%AF%E8%83%BD%E3%81%AB)があったと思うんですが、アレについて `-readonly` と `-?` が追加され、readonlyやoptionalを剥がすことができるようになりました。

あ、 `Required` が新しくビルトイン型として追加されました。

```ts
{
    // ビルトインの奴の定義を再掲
    type Readonly<T> = {
        readonly [P in keyof T]: T[P];
    };
    type Partial<T> = {
        [P in keyof T]?: T[P];
    };
}
{
    // 今回新しく追加されたビルトイン型
    type Required<T> = {
        [P in keyof T]-?: T[P];
    };
}
// ビルトインではないけどとりあえず
type Mutable<T> = {
    -readonly [P in keyof T]: T[P];
}

// 実験台
class Foo {
    str: string;
}

// ここの部分は前からあった
// { readonly str: string; } と評価される
type A = Readonly<Foo>;
// { readonly str?: string; } と評価される
type B = Partial<A>;

// ここから先は今回から
// { readonly str: string; } と評価される
type C = Required<B>;
// { str: string; } と評価される
type D = Mutable<C>;
```

ちなみに、 `+readonly` と `+?` という書き方もできるようになりました。
（lib.d.ts内部では使われていないけど）

## importの整理ができるようになった

Language Serviceにimportの整理が追加されました。
VSCode(Insiders)で、 `TypeScript: Organize Imports` が追加されている。

```ts
import { createElement } from "react";
import { createFactory } from "react";
```

が

```ts
import { createElement, createFactory } from "react";
```

になったりする。

[公式ブログのgif](https://blogs.msdn.microsoft.com/typescript/2018/03/27/announcing-typescript-2-8/#organize-imports)見たほうが早いかも。

<!-- TODO tsfmtでサポートする？ -->

## 初期化忘れのフィールドを雑に修正できるようになった

Language ServiceのQuick fixに追加された感じです。

```ts
class Foo {
    str: string;
    date: Date;
}

// Add 'undefined' type to property 'str' とかを実行するとこうなる
class FooA {
    str: string | undefined;
    date: Date | undefined;
}

// Add definite assignment to property 'str: string' とかを実行するとこうなる
class FooB {
    str!: string;
    date!: Date;
}

// Add initializer to property 'str' とかを実行するとこうなる
// Date の初期値は自明ではないので使えない…
class FooC {
    str: string = "";
    date: Date;
}
```

## `keyof` を intersection type に適用した時個別に簡約されるようになった

ざっくりこんな感じ。
全体的に読みやすくなりました。
早いタイミングで簡約されるようになったのでコンパイラにも優しいのではなかろうか。

```ts
type A = { a: string };
type B = { b: string };

// 今は "a" | "b" と評価される
// 前も "a" | "b" と評価される
type T1 = keyof (A & B);

// 今は keyof T | "b" と評価される
// 前は keyof (T & { b: string; }) と評価されていた
type T2<T> = keyof (T & B);

// 今は "a" | keyof U と評価される
// 前は keyof ({ a: string; } & U) と評価されていた
type T3<U> = keyof (A & U);

// 今は keyof T | keyof U と評価される
// 前は keyof (T & U) と評価されていた
type T4<T, U> = keyof (T & U);

// ここから下のやつは全部
// 今は "b" | "a" と評価される
// 前も "a" | "b" と評価される
type T5 = T2<A>;
type T6 = T3<B>;
type T7 = T4<A, B>;
```

## More special declaration types in JS

JSにあまり興味ないので自分で見てください…
[この辺](https://github.com/Microsoft/TypeScript/wiki/What%27s-new-in-TypeScript#better-handling-for-namespace-patterns-in-js-files)

## `--noUnusedParameters` で使ってない型パラメータも怒られるようになった

前は `--noUnusedLocals` の時に怒られてたらしいです。
正しいといえば正しいけど結局どっちもONなので大差ないぜ！

↓こういうのが怒られる

```ts
class Foo<T> {
}
```

## HTMLObjectElementがalt属性持ってたのを直した

そのまんまです。
https://github.com/Microsoft/TSJS-lib-generator の最新がmergeされたっぽいのでコレ以外にも変わったりしているものがあります。
例えば `document.createElement("button", {is: "pretty-button"})` 的なコードが通るようになった。[この辺](https://github.com/Microsoft/TSJS-lib-generator/pull/176)
