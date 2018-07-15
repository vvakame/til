こんにちは[メルペイ社](https://www.merpay.com/jp/)な[@vvakame](https://twitter.com/vvakame)です。

[TypeScript 3.0 RC](https://blogs.msdn.microsoft.com/typescript/2018/07/12/announcing-typescript-3-0-rc/)がアナウンスされました。

[What's new in TypeScriptも更新](https://github.com/Microsoft/TypeScript/wiki/What's-new-in-TypeScript#typescript-30)されているようです。
[破壊的変更](https://github.com/Microsoft/TypeScript/wiki/Breaking-Changes#typescript-30)もあるよ！

[この辺](https://github.com/vvakame/til/tree/master/typescript/v3.0.0-rc)に僕が試した時のコードを投げてあります。

## 変更点まとめ

* プロジェクト間の参照のサポート [Support for project references/composite projects](https://github.com/Microsoft/TypeScript/issues/3469) [2](https://github.com/Microsoft/TypeScript/issues/25600) [3](https://github.com/Microsoft/TypeScript/pull/25283) [4](https://github.com/Microsoft/TypeScript/pull/24850) [5](https://github.com/Microsoft/TypeScript/pull/23944) [6](https://github.com/Microsoft/TypeScript/pull/23354)
    * monorepo的構造の中で過ごしやすくなったっぽいですね
    * `compilerOptions.composite` と `references` がtsconfig.jsonに追加
        * 依存先プロジェクトの型が変わるような変更があったら適当にビルドしてくれる
    * `--build, -b` の追加
        * 実質 `--build` はサブコマンドみたいな感じのようだ
        * `--verbose`, `--dry`, `--clean`, `--force` オプションあり
        * なお、その他のオプション(`--ourDir` とか)と併用不可 [*](https://github.com/Microsoft/TypeScript/issues/25613)
    * share, server, client に個別にtsconfig.jsonある構成で `--outDir` したときに意図しないクソ出力がされるのを抑制的な
* 可変長引数にGenericsの型パラメータが割当可能に [Extracting and spreading parameter lists with tuples](https://github.com/Microsoft/TypeScript/pull/24897)
    * Function#call とかに型パラメータずらずら並べたoverloadしなくてよくなるやつ
    * 逆方向の型推論(callの第二引数以降から第一引数に渡した関数の仮引数の型が推論される)
    * ≒タプルの型周りの強化
* タプル周りの型の扱いの強化 `Richer tuple types` [1](https://github.com/Microsoft/TypeScript/pull/25408) [2](https://github.com/Microsoft/TypeScript/pull/22850) ?
    * `[number?]` のように書けるようになった
    * `[number, ...string[]]` のように書けるようになった
    * `[]` のように書けるようになった
* `unknown` typeの導入 [The `unknown` type](https://github.com/Microsoft/TypeScript/pull/24439)
    * type-safeな `any` の代替
    * どんな値でも `unknown` 型の変数に投げ込める
    * 逆は無理
    * 型アサーションやなローイング無しに何かに使うことができない
    * `unknown` は予約語になりました
* `defaultProps`がtsxでサポートされた [Support for `defaultProps` in JSX](https://github.com/Microsoft/TypeScript/pull/24422)
    * `JSX.LibraryManagedAttributes` が導入された
        * ここにうまいこと型を定義してやるとJSX syntaxでattributeの型チェックがいい感じに行われる
    * 上記を使ってReactの `defaultProps` をサポートできる
        * Propsでnon-nullな要素があっても、defaultPropsに値が定義されていればattributeに書かなくてもコンパイル時にハネられなくなる
        * ついでに `propTypes` も考慮にいれて型を構築できるらしい
        * 無指定の場合 `defaultProps` は `Partial<Props>` 相当なので注意
    * これで明示的なPropsの型を宣言せずにJSコード相当のものだけでいい感じに型が定まるようになり神
* `<reference lib="..." />` の追加 [`/// <reference lib="..." />` reference directives](https://github.com/Microsoft/TypeScript/pull/23893)
    * built-inのlibを参照させられる
        * core-js とか es6-shim に記述されてたら便利だよね的な
        * target `es5` でlibに `es2015.promise` 指定したいとかよくある
    * `no-default-lib` あっても無視される
* エラーが発生する箇所とエラー発生源の両方でDiagnotice出したい的な？ [Related error spans](https://github.com/Microsoft/TypeScript/issues/25257)
    * なんかまだ作業途中で一部が入ってきているっぽい
* 同上かな [Improved message quality and suggestions](https://github.com/Microsoft/TypeScript/issues/25310)
* `*` なimportとnamed importを相互にリファクタリング可能に [Convert named imports to namespace imports and back](https://github.com/Microsoft/TypeScript/pull/24469)
    * `import * as m from "mo";` と `import { a } from "mo";` を相互に変換可能に
* アロー関数のbodyの `{}` をつけたり剥がしたりのリファクタリングが可能に [Add or remove braces from arrow function](https://github.com/Microsoft/TypeScript/pull/23423)
    * `a => a` と `a => { return a; }` を相互に変換可能に
* いらんラベルを剥ぐQuickFixの追加 [Remove unused labels](https://github.com/Microsoft/TypeScript/pull/24037)
    * まんまです
* 到達不能なコードを除去するQuickFixの追加 [Remove unreachable code](https://github.com/Microsoft/TypeScript/pull/24028)
    * まんまです
* JSXタグの折りたたみサポート [Outlining spans for JSX expressions](https://github.com/Microsoft/TypeScript/issues/23273)
    * `<hoge><fuga><piyo /></fuga></hoge>` 的なやつをエディタの見た目上 `<hoge>...</hoge>` って折り畳めるようにするやつのはず
    * [3.0.1](https://github.com/Microsoft/TypeScript/milestone/75) っぽい
* エディタ側からJSXの閉じタグ何使えばいいか教えてくれ！って要求できるやつ [Auto-closing JSX tags](https://github.com/Microsoft/TypeScript/issues/20811)
    * まんまです

* [API breaking changesもちょいちょいある](https://github.com/Microsoft/TypeScript/wiki/API-Breaking-Changes#typescript-30)


### 破壊的変更！

上記リストのうち、破壊的変更を伴うのは次のものになります。

* `unknown` typeの導入

## プロジェクト間の参照のサポート

monorepo向けの機能として、プロジェクト間の参照がサポートされるようになりました。
現時点ではちょっと使いづらいです。

これに伴い、全プロジェクトを一括ビルドするために `--build` というオプション(実質サブコマンド)が追加されています。

先に罠のまとめを書いておきます。

* `"composite": true` したら `"declaration": true` も絶対指定しろ絶対にだ
* tsconfig.json生成したり変更したら該当プロジェクトでのビルドを絶対1回やれ
* 関連するプロジェクトは全部referencesに書け 参照先の参照先もだ
* 関数などのリファクタリングはかならず参照元プロジェクトでやれ
* `tsc --build --clean` と `tsc --build` は絶対ワンセット
* `--outDir` 使っている場合は参照先に注意
* エディタがおかしいと思ったらおとなしくtsserverをリスタートさせる

> continuing to improve the experience around editor support using

と書かれていたのでそのうち罠も少なくなると思います。
RCが取れる時に治ってたら神。

細かいことは[ここ](https://github.com/vvakame/til/tree/master/typescript/v3.0.0-rc/project-refs/)に置いてあるコードを色々と自分で試してみてください。

### tsconfig.jsonに書くべき内容

tsconfig.jsonに設定するときのキーとなる要素は次に抜粋した通りです。

```tsconfig.json
{
  "references": [
    {
      "path": "../core/"
    }
  ],
  "compilerOptions": {
    "declaration": true,                      /* Generates corresponding '.d.ts' file. */
    "composite": true                         /* Enable project compilation */
  }
}
```

`"composite": true` を指定し、かつ `"declaration": true` を指定する必要があります。
declaration無しでもエディタ上では一瞬それっぽく動くんですがコンパイル時に怒られたりするので絶対両方指定します。
現時点ではエディタの言うことを信じてデバッグしようとすると辛い目にあうのでtscの出力を信用しましょう。

### `--build` について

tscに `--build` というオプションが追加されています。
`-b` でもOKです。
これは `--verbose`, `--dry`, `--clean`, `--force` というオプションを持ちます。
しかし、他の、例えば `--outDir` みたいなものと併用することができません。
つまり、 `--build` は実質サブコマンドです。

`--verbose` はビルドの詳細を出力する。
`--dry` はビルドしてみるけどファイルの出力はなし。
`--clean` はビルドした結果出力されるファイルを削除する。
`--force` はビルド時に更新の必要がないコードもビルドしなおす。

という感じです。

### `--outDir` などの処理対象がちゃんとプロジェクト毎になる

特定のプロジェクトで `tsc --outDir dist` などの単体ビルドすると、そのプロジェクトの生成物のみ出力されます。
昔はプロジェクトという概念がなかったので client と server がある時に両方のコードが変なところに変な構造で出力されていたんですが、これが解消されます。
なかなか良いですね。

### 細かい説明

この説明では、 client → shared → core というプロジェクト構成である前提で説明をします。

referencesには間接的に参照するプロジェクトもすべて指定します。
具体的に、clientにはsharedもcoreも両方指定する必要があります。

この機能は現在では `.d.ts` の存在やタイムスタンプに様々な処理が依存しているようで、これに起因する不便が色々なところにあります。
例えば、 `tsc --build --clean` した後に `tsc` すると `.d.ts` が生成されてないのでエラーになったりします。
素直に `tsc --build` をワンセットにして使いましょう。

参照先プロジェクトが `--outDir` で `.d.ts` の出力先変えている場合、モジュールパスの指定を間違えると実行時エラーになります。
例えば、coreが `src` にソースがあり `dist` に出力している場合、他のところからは `../core/src` とかするとビルド成功はするけど実行時エラーになります。
実行時エラーになるのはクッソつらいので頑張って気をつけましょう…。

.d.tsに依存する都合上、clientでcoreの関数をリネームするようなリファクタリングを行うとcoreの.d.tsファイルのみが変更されて全体としては壊れます。
なので、coreの関数をリネームしたい場合はcoreでリネームし、`tsc --build` してエラー箇所を修正して回ることになります。
全体が一気に変わってくれるとめっちゃ嬉しいんだけどなー。

参照先プロジェクトのtsconfig.jsonにエラーがあっても通知されない場合があります。
例えば、coreのtsconfig.jsonにエラーがあっても、 `tsc --build` する時にcoreのキャッシュが使える場合、tsconfig.jsonのチェックもスキップされます。
どこかのtsconfig.jsonをいじった場合、 `tsc --build --force` してエラーが発生しないか確認するのがよいでしょう。

Visual Studio Code - Insiders でも、まだエラー表示などが正しくない場合がありそうです。
エディタ上での表示を信用せず、tscコマンドの出力を信用しましょう。
これは、tsserverが `tsc --build --dry` 相当の操作をせず、 `tsc --build` の生成物に頼って処理を行っているのだと思います。

`error TS2307: Cannot find module '../../core/src'.`　的なエラーが意図せず出た場合、tsconfig.jsonの設定を全体的に見直すか、`tsc --build --clean && tsc --build --force` で治る or 原因が判明する場合が多いのではないかと思います。

コンパイルにかかる時間がそこそこ長い気がしますね。
内部的な処理を考えると仕方ない気もしますが…。

困ったら `tsc --build --clean && tsc --build --force --verbose` しておきましょう。
ということです。

## 可変長引数にGenericsの型パラメータが割当可能に

可変長引数に対して型パラメータを割り当て、これを流用することができるようになりました。
ざっくり次のようなコードが書けます。

```ts
function call<TS extends any[], R>(fn: (...args: TS) => R, ...args: TS): R {
    return fn(...args);
}

function hello(word = "TypeScript") {
    return `Hello, ${word}`;
}

// TS は [(string | undefined)?] と推論されている
let str1 = call(hello, "JavaScript");
let str2 = call(hello, void 0);
let str3 = call(hello);

// これはちゃんとエラーになる！えらい！
// index.ts:18:17 - error TS2345: Argument of type '(word?: string) => string' is not assignable to parameter of type '(args_0: boolean) => string'.
//   Types of parameters 'word' and 'args_0' are incompatible.
//     Type 'boolean' is not assignable to type 'string | undefined'.
// let str4 = call(hello, true);

// 引数から TS を推論させて word の型指定を省略することもできる かしこい
call(word => `Hello, ${word.toUpperCase()}`, "TypeScript");
```

強いですね。
TypeScriptバンドルのlib.d.ts系でこれを使っているものはまだなさそうに見えます。

## タプル周りの型の扱いの強化

次のコードの通りです。
`[number?]`, `[...number[]]`, `[]` 的な表記ができるようになりました。

```ts
// optionalな要素を簡単に書けるようになった 前は (number | undefined) でしたね
type Coordinate = [number, number, number?];

function coordinates(...args: Coordinate) {
    const [x, y, z] = args;
    return { x, y, z };
}
// v はちゃんと { x: number; y: number; z: number | undefined; } と推論される えらい！
let v = coordinates(1, 2);

// arrayのspreadingみたいな記法がサポートされた
// string[] に評価される
type SpreadedStrings = [...string[]];

// 空のtupleも作れるようになった
type Empty = [];
```

前述の"可変長引数にGenericsの型パラメータが割当可能に"と併せて、可変長引数にtupleな型を指定できるようになっています。

## `unknown` typeの導入

"なんだかわからん"ことを表す型として `unknown` が登場しました。
なんだかわからん時に使いましょう。

ざっくりした性質は次のコードの通り。

```ts
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
```

なお、`unknown` は予約語となりました。
と言っても型の宣言空間での予約語なので変数名としては普通に使えます。

`unknown` が絡む型の演算についてどうなるかは[PRのdescription](https://github.com/Microsoft/TypeScript/pull/24439)に詳しく書かれているのでそちらを見てください。

## `defaultProps`がtsxでサポートされた

というとちょっと語弊があるんですが、 `JSX.LibraryManagedAttributes` が導入されました。
[handbookのJSXのページ](https://www.typescriptlang.org/docs/handbook/jsx.html)にはまだ書かれてないのでアレなんですが…。
JSX関連のfeatureをサポートする時に直接的な実装をせずに拡張用のフックポイントを実装するのは相変わらずですね。

JSXライブラリがJSXタグに要求するattributeを制御できます。
[実装のPR](https://github.com/Microsoft/TypeScript/pull/24422)のテストケースを見るとわかりますが、自力実装は型魔術師じゃないと難しいと思います。
筆者は厳しそうです。
というか型周りにどういう記法があって何ができるかがそろそろ脳メモリに乗らない分量…。

実装のPRのテストケースから引っ張ってきた型定義です。
これがReact向けの完全な定義ではないので細かいところの解説は避けます。

```ts
// defaultPropsを考慮したPropsを計算する
type Defaultize<TProps, TDefaults> =
    // defaultPropsに含まれるプロパティはoptionalにする
    & { [K in Extract<keyof TProps, keyof TDefaults>]?: TProps[K] }
    // defaultPropsに含まれないものはそのまま
    & { [K in Exclude<keyof TProps, keyof TDefaults>]: TProps[K] }
    // defaultPropsにしか含まれないものはoptionalにする
    & Partial<TDefaults>;

// propTypesに指定されたPropTypesの値から求められる型を計算する
type InferredPropTypes<P> = { [K in keyof P]: P[K] extends PropTypeChecker<infer T, infer U> ? PropTypeChecker<T, U>[typeof checkedType] : {} };

// 型計算用のシンボルを定義
declare const checkedType: unique symbol;
// propTypesに定義する値 何の型を要求するかとそれがrequiredかどうか
interface PropTypeChecker<U, TRequired = false> {
    (props: any, propName: string, componentName: string, location: any, propFullName: string): boolean;
    isRequired: PropTypeChecker<U, true>;
    [checkedType]: TRequired extends true ? U : U | null | undefined;
}

// 値空間に存在するPropTypesの定義
declare namespace PropTypes {
    export const number: PropTypeChecker<number>;
    export const string: PropTypeChecker<string>;
    export const node: PropTypeChecker<ReactNode>;
}

type ReactNode = string | number | ReactComponent<{}, {}>;

declare class ReactComponent<P={}, S={}> {
    constructor(props: P);
    props: P & Readonly<{ children: ReactNode[] }>;
    setState(s: Partial<S>): S;
    render(): ReactNode;
}

declare namespace JSX {
    interface Element extends ReactComponent { }
    interface IntrinsicElements { }

    type LibraryManagedAttributes<TComponent, TProps> =
        // defaultProps, propTypes の両方が存在するか
        TComponent extends { defaultProps: infer D; propTypes: infer P; }
            ? Defaultize<TProps & InferredPropTypes<P>, D>
            // defaultProps が存在するか
            : TComponent extends { defaultProps: infer D }
                ? Defaultize<TProps, D>
                // propTypes が存在するか
                : TComponent extends { propTypes: infer P }
                    ? TProps & InferredPropTypes<P>
                    : TProps;
}
```

この定義がスルスル読める人はTypeScriptマスター度がかなり高いと思います。
こんなん書くのは頭のよろしい人にまかせておけばよろしおす…。

先の定義の利用例はこんな感じです。
([Stage 2の仕様が通れば](https://github.com/tc39/proposal-static-class-features))esnext validな、型周りの情報が全く存在しない形で書いたのにtype safeだった！
という感じです。

```ts
class Component extends ReactComponent {
    static propTypes = {
        foo: PropTypes.number,
        bar: PropTypes.node,
        baz: PropTypes.string.isRequired,
    };
    static defaultProps = {
        foo: 42,
    }
}

// ざっくりこういうイメージ 型なので = 42 とかは本来は書けない
// type Props = {
//     foo: number = 42;
//     bar: ReactNode;
//     baz: string;
// };

// OK
const a = <Component foo={12} bar="yes" baz="yeah" />;
const b = <Component bar="yes" baz="yeah" />;
const c = <Component foo={12} bar={null} baz="cool" />;

// エラーになる奴ら
// barはundefinedを受け付けるけどプロパティ自体を省略できるわけではない(≒型定義の不備に近い)
// const d = <Component foo={12} baz="yeah" />;
// bazがないのでエラー
// const e = <Component foo={12} bar={void 0} />;
// batはPropsの定義に存在しないのでエラー
// const f = <Component bar="yes" baz="yo" bat="ohno" />;
// bazはnon-nullなのでエラー
// const g = <Component foo={12} bar="yeah" baz={null} />;
```

曲芸か何かかな…？

[DefinitelyTypedのReactに実装するPR](https://github.com/DefinitelyTyped/DefinitelyTyped/pull/27267)はこちら。

## `<reference lib="..." />` の追加

主に型定義ファイル向けの機能です。
`<reference lib="..." />` 的な感じで参照したいlibを書いておくとtargetなどに関係なくその定義が使えるようになります。
以前は core-js とか es6-shim を導入した時にtsconfig.jsonに使える型の範囲を自分で追加しないといけなかったけど、該当ライブラリの型定義にこれが書いてあればユーザは何も考えなくてもそれが使える、というわけです。

```ts
/// <reference lib="es2015.promise" />

// target=es5 だけどPromiseが使えます
Promise.resolve("test");
```

## エラーが発生する箇所とエラー発生源の両方でDiagnotice出したい的な？ & 同上

エラーになっている箇所とエラーの発生源がコード上の別の箇所であることはままあります。
そういう時に、関連箇所の表示ができるよう、tsserver側もエディタ側も拡張する必要がある… 的な文脈です。

例えば次のようなコードだと、`afterDeclared` が定義される前に利用されているのでエラーになります。
この機能では、2行目で使われているやつやで！と教えてくれて、該当の箇所にジャンプすることもできます。

```ts
console.log(afterDeclared);
let afterDeclared = true;
```

tscコマンド的にも次のような表示にしてくれます。

```
src/errors/indexIgnore.ts:1:13 - error TS2448: Block-scoped variable 'afterDeclared' used before its declaration.

1 console.log(afterDeclared);
              ~~~~~~~~~~~~~

  src/errors/indexIgnore.ts:2:5
    2 let afterDeclared = true;
          ~~~~~~~~~~~~~
    'afterDeclared' was declared here.
```

## `*` なimportとnamed importを相互にリファクタリング可能に

まんまです。
1行全体を選択すると使えるQuick Fixですね。

これが

```ts
import * as m from "../../project-refs/core/dist/";

console.log(m.hello("TypeScript"));
```

こうなって

```ts
import { hello } from "../../project-refs/core/dist/";

console.log(hello("TypeScript"));
```

こう

```ts
import * as dist from "../../project-refs/core/dist/";

console.log(dist.hello("TypeScript"));
```

名前が干渉する場合もいい感じに避けてくれます。

```ts
import { hello as hello_1 } from "../../project-refs/core/dist/";

const hello = "hello";

console.log(hello_1("TypeScript"));
```

後は細かく自分でrenameしましょう。

## アロー関数のbodyの `{}` をつけたり剥がしたりのリファクタリングが可能に

アロー関数全体か一部を選択して実行するQuick Fixです。

```ts
const f = (word:string) => `Hello, ${word}`;
```

こう

```ts
const f = (word:string) => {
    return `Hello, ${word}`;
};
```

こう

```ts
const f = (word:string) => `Hello, ${word}`;
```

よさそう。

次のような単一式じゃないものは変換できないようです。
そらそやな。

```ts
const f2 = (word: string) => {
    word = `${word}!`;
    return word;
};
```

## いらんラベルを剥ぐQuickFixの追加

はい。

```ts
loop: for (let i of [1, 2, 3]) {
    console.log(i);
}
```

```ts
for (let i of [1, 2, 3]) {
    console.log(i);
}
```

## 到達不能なコードを除去するQuickFixの追加

まんまです。
実装PRのテストコードより。

```ts
function f() {
    return f();
    return 1;
    function f() {}
    return 2;
    type T = number;
    interface I {}
    const enum E {}
    enum E {}
    namespace N { export type T = number; }
    namespace N { export const x = 0; }
    var x;
    var y = 0;
}
```

これが

```ts
function f() {
    return f();
    function f() {}
    type T = number;
    interface I {}
    const enum E {}
    namespace N { export type T = number; }
    var x;
}
```

こうなる

## JSXタグの折りたたみサポート

まんまのはずだけど3.0.1マイルストーンなので多分まだ未実装。

## エディタ側からJSXの閉じタグ何使えばいいか教えてくれ！って要求できるやつ

これもまんまです。
実際に試してみると `<Hoge>` まで入力すると `</Hoge>` が自動入力されました。
わかりにくい。
