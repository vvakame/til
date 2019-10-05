# TypeScript v3.6.3 変更点

こんにちは[メルペイ社](https://www.merpay.com/)な[@vvakame](https://twitter.com/vvakame)です。

[TypeScript 3.6](https://devblogs.microsoft.com/typescript/announcing-typescript-3-6/)がアナウンスされました。
…えぇ、3.7じゃなくて3.6の記事です。現時点では[3.7 Beta](https://devblogs.microsoft.com/typescript/announcing-typescript-3-7-beta/)が出ています。
技術書典の運営タスクに圧殺され3.6.0-rcは未実装なものが多いな… つってスルーしてたらこんな有様です。
実際の動作確認は現在のlatestであるv3.6.3で行っています。

* [What's new in TypeScript in 3.6](https://www.typescriptlang.org/docs/handbook/release-notes/typescript-3-6.html)
    * [wikiページ](https://github.com/microsoft/TypeScript/wiki/What%27s-new-in-TypeScript#typescript-36)はサイズでかすぎてレンダリングできなくなってた…
    * wikiページのmarkdownからhandbook側も生成しているはずなので問題ないはず
* [Breaking Changes in 3.6](https://github.com/microsoft/TypeScript/wiki/Breaking-Changes#typescript-36)
* [TypeScript 3.6 Iteration Plan](https://github.com/microsoft/TypeScript/issues/31639)
* [TypeScript Roadmap: July - December 2019](https://github.com/microsoft/TypeScript/issues/33118)

Roadmapは[こちら](https://github.com/Microsoft/TypeScript/wiki/Roadmap#36-august-2019)。

[この辺](https://github.com/vvakame/til/tree/master/typescript/v3.6.3)に僕が試した時のコードを投げてあります。

## 変更点まとめ

* より厳密なジェネレータへの型付け [Strongly typed iterators and generators](https://github.com/Microsoft/TypeScript/issues/2983)
    * next の引数の型がチェックされるようになった
    * doneの値によってvalueの型が違う場合に区別できるようになった
* Array spreadingの挙動の修正 [More accurate array spreads](https://github.com/microsoft/TypeScript/pull/31166)
    * `[...Array(5)]` のes5へのdownpileがより正確に行われるようになったらしい
    * tslibに `__spreadArrays` が追加された感じ
    * Issue立ってから実に3年越しの修正
* Promiseの使い方下手こいた時のUXの改善 [Improved UX around Promises](https://github.com/microsoft/TypeScript/issues/30646)
    * await いるよね？っていってQuick fixを提示してくれる
    * 今まで気にしてなかったけど確かに欲しいやつ…！
* セミコロンがない時のStatement追加時の挙動を改善 [Semicolon-aware auto-imports](https://github.com/microsoft/TypeScript/issues/19882)
    * セミコロンレス派が喜びそうなやつです
    * セミコロンちゃんと書け(or フォーマッタに自動追加させろ)派としてはいるのかこれ と思っています
* よりよい識別子でのUnicodeのサポート [Amend scanner to support astral characters in identifiers when parsing es6+](https://github.com/microsoft/TypeScript/pull/32096)
    * `𝓱𝓮𝓵𝓵𝓸` とかを使えるようになった
* `import.meta` がSystemJSでサポートされるようになった [Add support for import.meta in System modules](https://github.com/microsoft/TypeScript/pull/32797)
    * そのまんまです
* `get` と `set` が型定義中で使えるようになった [Allow accessors in ambient class declarations](https://github.com/microsoft/TypeScript/pull/32787)
    * より厳密にECMAScriptの仕様を表現できるように
    * TypeScript 3.7ではtsコードからの.d.ts生成の結果も変わるらしい
    * interfaceやobject type literalでの利用は未来送りになったようだ
* 型定義中で関数の定義とクラスの定義がマージできるようになった [Allow functions and ambient classes to merge](https://github.com/microsoft/TypeScript/pull/32584)
    * 今まで `var` と `new` とかで頑張ってたやつがより素直に書けるように
    * 主に `new Date()` と `Date()` の挙動を表現したりするのに使われていた
* Compiler APIで `--build` と `--incremental` が利用可能になった [Api for tsc --build and --incremental](https://github.com/microsoft/TypeScript/pull/31432)
    * 普通の人には関係ないやつです
    * `createSolutionBuilder` やらその他色々… PRのdiffもでかい
* モジュールインポートの時、空気を読んでシンタックスを選んでくれるようになった [Make auto-imports more likely to be valid for the file (including JS) & project settings](https://github.com/microsoft/TypeScript/pull/32684)
    * 今までは ES Module が自動的に選択されていたけど `import foo = require("foo")` とかが空気を読んで選択されるように
* 新しいTypeScript Playground
    * https://www.typescriptlang.org/play/ がリッチになったよ
    * https://github.com/agentcooper/typescript-play ベースだよ
    * https://github.com/microsoft/TypeScript-Website/ に置いてあるよ

<!--
Roadmap的にはv3.6となっているがmergeされたのは最近でv3.7.0のマイルストーンに入っている。
* `--declaration` と `--isolatedModules` の併用の改善 [`--declaration` and `--isolatedModules`](https://github.com/Microsoft/TypeScript/issues/29490)
* async内でawaitがいるような候補を選んだら自動的にawaitを挿入する [Auto-inserted `await` for completions](https://github.com/microsoft/TypeScript/issues/31450)

まだ実装が存在しない
* Call Hierarchyのサポート [Call Hierarchy support](https://github.com/microsoft/TypeScript/issues/31863)
    * Find Referencesはすでにあるけどガバガバ辿れたほうがいいよねー的な
-->

<!--
これバグだったんだ… シリーズ
* https://github.com/microsoft/TypeScript/issues/30471
-->

## 破壊的変更！

* `"constructor"` という名前のメソッドがコンストラクタ扱いされるようになった [Parse quoted constructors as constructors, not methods](https://github.com/microsoft/TypeScript/pull/31949)
* DOMの更新
    * `window` の定義が `Window` から `Window & typeof globalThis` に変更
    * `GlobalFetch` がなくなった 代わりに `WindowOrWorkerGlobalScope` を使う
    * `Navigator` にあった非標準のプロパティが消えた
    * `experimental-webgl` 　がなくなった 代わりに `webgl` か `webgl2` を使う
* JSDocコメントが複数ある場合にmergeされなくなった [Use only immediately preceding JSDoc](https://github.com/microsoft/TypeScript/pull/32181)
    * 最下部のコメントだけ有効
* キーワードにエスケープシーケンスを含められなくなった [Add error message for keywords with escapes in them](https://github.com/microsoft/TypeScript/pull/32718)
    * `\u0063ontinue` とかやった時に今までは `continue` に変換されてたけどエラーになるようになった

## より厳密なジェネレータへの型付け

ジェネレータは(一般的な利用頻度が低いこともあってか)今まであまりよい型付をすることができませんでした。
今回、これが改善され、 1. next の引数の型がチェックされるようになった 2. doneの値によってvalueの値がtype narrowingされるようになった という感じです。

```ts
// こういう感じの定義が
interface Generator extends Iterator<any> { }

interface Iterator<T> {
    next(value?: any): IteratorResult<T>;
    return?(value?: any): IteratorResult<T>;
    throw?(e?: any): IteratorResult<T>;
}

interface IteratorResult<T> {
    done: boolean;
    value: T;
}

// こうなったりした
interface Generator<T = unknown, TReturn = any, TNext = unknown> extends Iterator<T, TReturn, TNext> {
    // NOTE: 'next' is defined using a tuple to ensure we report the correct assignability errors in all places.
    next(...args: [] | [TNext]): IteratorResult<T, TReturn>;
    return(value: TReturn): IteratorResult<T, TReturn>;
    throw(e: any): IteratorResult<T, TReturn>;
    [Symbol.iterator](): Generator<T, TReturn, TNext>;
}

interface Iterator<T, TReturn = any, TNext = undefined> {
    // NOTE: 'next' is defined using a tuple to ensure we report the correct assignability errors in all places.
    next(...args: [] | [TNext]): IteratorResult<T, TReturn>;
    return?(value?: TReturn): IteratorResult<T, TReturn>;
    throw?(e?: any): IteratorResult<T, TReturn>;
}

interface IteratorYieldResult<TYield> {
    done?: false;
    value: TYield;
}

interface IteratorReturnResult<TReturn> {
    done: true;
    value: TReturn;
}

type IteratorResult<T, TReturn = any> = IteratorYieldResult<T> | IteratorReturnResult<TReturn>;
```

type narrowingが効きやすそうな定義になりましたね。
利用例を見てみます。

```ts
// 返り値の型は自動的にいい感じに推論される
function* counter() /* Generator<number, string, boolean> */ {
    console.log("Start!");
    let i = 0;
    while (true) {
        // ここの変数の型指定は必要 next の引数の型の推論に利用されるため
        //   なしの場合、 any な値が出てくる
        //   ジェネレータ関数自体の TNext は unknown のはずだが unknown が出てきちゃうとBCなので仕方なさそう
        const v: boolean = yield i;
        if (v) {
            break;
        }
        i++;
    }
    return "done!";
}

let iter = counter();
console.log("ready?");

// 最初の yeild までを実行
let curr = iter.next();

while (!curr.done) {
    console.log(curr.value);
    // whileループ内では curr.value は number とわかっている
    // curr.done は false だから
    curr = iter.next(curr.value === 3);

    // next の引数もチェックされる
    // error TS2345: Argument of type '[123]' is not assignable to parameter of type '[] | [boolean]'.
    // iter.next(123);

    // 残念ながらこれはvalid
    // [] or [boolean] を受け付けるため
    // iter.next();
}

// ループの外では curr.done === true なので curr.value は string とわかっている
console.log(curr.value.toUpperCase());

// 次のような出力になる
// ready?
// Start!
// 0
// 1
// 2
// 3
// DONE!
```

async generatorに関しても同様です。


## Array spreadingの挙動の修正

`[...Array(5)]` のes5へのdownpileがより正確に行われるようになったらしいです。
今までは `Array(5).slice()` に変換されていて、これは仕様に対して微妙に異なる挙動だったのが修正されました。

この挙動をサポートするため、tslibに `__spreadArrays` が追加されました。
[Issue](https://github.com/microsoft/TypeScript/issues/8856)が立ってから実に3年越しの修正でした。

ちなみに `--downlevelIteration` が使われている時は今までも `__spread` が使われ、仕様に沿った結果になっていました。

```ts
// tslib に __spreadArrays が追加されました 今回の変更をサポートするため
import { __spreadArrays } from "tslib";

// [empty × 3] と表示される in Chrome
console.log(Array(3));
// [ undefined, undefined, undefined ] と表示される in Chrome
console.log([...Array(3)]);


// false と表示される
// 長さは3だがプロパティが存在しないため
// 不正確だが雰囲気が伝わる記述をすると { length: 3 } みたいな感じ
console.log(1 in Array(3));

// false と表示される
// 上に同じくプロパティが存在しないため
console.log(1 in Array(3).slice());


// true と表示される
// [ undefined, undefined, undefined ] と解釈されるため
// 不正確だが雰囲気が伝わる記述をすると { 0: undefined, 1: undefined, 2: undefined, length: 3 } みたいな感じ
console.log(1 in [...Array(3)]);


// TypeScript 3.5 までは…
// [...Array(3)] は Array(3).slice() とdownpileされていた
// しかし、これはプロパティの有無という面で厳密に一致した挙動ではない
// これが今回改められた、という話

// false
console.log(1 in Array(3));
// true
console.log(1 in [...Array(3)]);
// true
console.log(1 in __spreadArrays(Array(3)));
```


## Promiseの使い方下手こいた時のUXの改善

Promiseをunwrap ( .then ) し忘れてた時に、やってなくない？と `await` を追加するQuick fixが追加されました。

```ts
interface User {
    name: string;
    age: number;
    location: string;
}

declare function getUserData(): Promise<User>;
declare function displayUser(user: User): void;

async function f1() {
    // 普通のエラーと改善方法の提案が出る
    // error TS2345: Argument of type 'Promise<User>' is not assignable to parameter of type 'User'.
    //   Type 'Promise<User>' is missing the following properties from type 'User': name, age, location
    //
    // `getUserData()` 部分に対して Did you forget to use 'await'?
    // displayUser(getUserData());
    
    // Quick fix を適用するとこうなる
    displayUser(await getUserData());
}

async function getCuteAnimals() {
    // error TS2339: Property 'json' does not exist on type 'Promise<Response>'.
    // 
    // `json` 部分に対して Did you forget to use 'await'?
    // fetch("https://reddit.com/r/aww.json").json();

    // Quick fix を適用するとこうなる
    (await fetch("https://reddit.com/r/aww.json")).json();
}
```


## セミコロンがない時のStatement追加時の挙動を改善

まんまです。
セミコロンつけろよ派なのでつけたほうがいいと思います。

```ts
// 3. 入力補完結果にセミコロン有無の好みが反映される
import { __spreadArrays } from "tslib"

// 1. ファイル中でセミコロンの有無を見て
console.log("foo")

// 2. 何らかのコードが自動的に補完されることをすると
__spreadArrays
```


## よりよい識別子でのUnicodeのサポート

芸人が喜びそう（暴言）

```ts
// 今回から es2015 target 以降で利用できるようになった
const 𝓱𝓮𝓵𝓵𝓸 = "world";
console.log(𝓱𝓮𝓵𝓵𝓸);
```

## `import.meta` がSystemJSでサポートされるようになった

らしいです。
SystemJSをもう使っているのでわからん…！
興味がある人は公式サイトの説明を見てください。

## `get` と `set` が型定義中で使えるようになった

らしいです。
今までは `error TS1086: An accessor cannot be declared in an ambient context.` とかいって怒られてました。

ECMA Script仕様のクラスフィールドへの対応の一環のようです。
フィールドとアクセサーを区別できないと、適切にエラーが出力ができなくなるのを回避するためだそうです。

最初はinterfaceやobject type literalでの利用もサポートしたかったようですが、一旦未来送りになりました。

```ts
declare class Foo {
    // 型定義の宣言で get, set は今まで使えなかった
    get x(): number;
    set x(val: number);
}

export { Foo };
```

これが許されるようになりました。

これを利用した時、メソッドのスタブをQuick fixに出力させると次のようになります。

```ts
import { Foo } from "./basic";

class FooImpl implements Foo {
    get x(): number {
        throw new Error("Method not implemented.");
    }    set x(val: number) {
        throw new Error("Method not implemented.");
    }

}
```

なんかインデントがずれてますね…。
ともあれ、生成コードコードもより意図が反映されたものになるわけです。

ちなみに、今はまだ `.ts` ファイルをコンパイルした結果の型定義ファイルはget, setなしの定義が出力されます。

```ts
// これが
class Foo {
    get x(): number { return 1; }
    set x(val: number) { }
}
```

```ts
// こうなる
declare class Foo {
    x: number;
}
```

TypeScript v3.7以降ではでは型定義の出力もset, getありのものになる予定です。

```ts
// npx typescript@next の出力を確認
declare class Foo {
    get x(): number;
    set x(val: number);
}
```

## 型定義中で関数の定義とクラスの定義がマージできるようになった

lib.d.ts とかを眺めたことがある人は、次のような定義を見たことがあると思います。

```ts
// 抜粋
declare var Date: DateConstructor;
interface DateConstructor {
    new(): Date;
    (): string;
}
```

こうなっているのには色々と理由があった気がしますが細かいことは忘れました…。
昔のJavaScriptは、ある識別子(上記の例だと `Date` )が関数だった場合、それを普通に呼び出したりコンストラクタとして使うことができました（今もやればできるが普通やらない）。
また、関数に好き勝手なプロパティを勝手に生やしたりすることができて、完全にカオスでした。
今はECMAScript 2015 afterの世界観なので、みんな平和に生きているのです…。
でまぁカオスを頑張って記述できる必要があるTypeScriptは様々な仕様の整合性を考えた結果、カオスを内包したりしていたわけです。

これを、素直（？）に次のように書けるようになりました。

```ts
declare function Date(): string;
declare class Date {
    constructor();
}
```

やっとか…！
今後、lib.d.tsなどの定義も次第にこのスタイルに書き換わっていくのではないでしょうか。


## Compiler APIで `--build` と `--incremental` が利用可能になった

らしいです。
話が若干マニアックな方向に行くのと、筆者が今のところあまり興味を持っていないパートなので割愛します。
公式の説明や該当のPRをチェックしてください。


## モジュールインポートの時、空気を読んでシンタックスを選んでくれるようになった

らしいです。
例えば、次のようなCommonJS形式のモジュールとES Module形式のモジュールを用意します。

```ts
function hello1() {
    console.log("Hello, world!");
}

export = hello1;
```

```ts
export function hello2() {
    console.log("Hello, world!");
}
```

それぞれ `hello1`, `hello2` を入力しようとして、importのパートを自動で補完させると次のようになります。

```ts
// hello1 は CommonJS 形式で書いたモジュール
import hello1 = require("./hello1");
// hello2 は ES Module 形式で書いたモジュール
import { hello2 } from "./hello2";

hello1();
hello2();
```

うーん、便利…かな？

tsconfig.jsonの設定値や、import対象がJSかTSかによって挙動が異なります。
`esModuleInterop` か `allowSyntheticDefaultImport` が有効な場合、ES Module形式が利用される… とPRの概要に書いてあるんですがなんかそうでもない気がする…。
基本的にはimport元の定義方法に依存していそう。

`@types/moment` とかがCommonJS形式時代の型定義のままなので、自分で試してみてください…。


## 新しいTypeScript Playground

https://www.typescriptlang.org/play/#code/IYGwpgTgLgFARIZQZB2DICIZBiDIawZA8XoBQYmAMGOASgG4AoIA


## `"constructor"` という名前のメソッドがコンストラクタ扱いされるようになった

らしいです。
ECMAScriptの仕様は複雑ですね…。

```ts
class A {
    // 普通のコンストラクタ
    constructor() {
        console.log("A");
    }
}

class B {
    // コンストラクタと認識されるようになった
    "constructor"() {
        console.log("B");
    }

    // 2つ定義することになるのでエラーになる
    // error TS2392: Multiple constructor implementations are not allowed.
    // constructor() {}
}

class C {
    // computed propertyの場合コンストラクタにはならない
    ["constructor"]() {
        console.log("C");
    }

    // 重複定義にはならないのでエラーにならない
    constructor() {}
}

// A と表示
new A();
// B と表示
new B();
// なにも表示されない
new C(); // .constructor() が生えてる
```


## DOMの更新

* `window` の定義が `Window` から `Window & typeof globalThis` に変更
* `GlobalFetch` がなくなった 代わりに `WindowOrWorkerGlobalScope` を使う
* `Navigator` にあった非標準のプロパティが消えた
* `experimental-webgl` 　がなくなった 代わりに `webgl` か `webgl2` を使う

って感じらしいです。


## キーワードにエスケープシーケンスを含められなくなった

らしいです。

```ts
while (true) {
    // error TS1260: Keywords cannot contain escape characters.
    \u0063ontinue;
}
```


## おまけ：`import console = require("console")` 勝手に入れられるのは実はバグだった

[Fix export=global auto-import exclusion](https://github.com/microsoft/TypeScript/pull/32898)

割と不便だったんだけどそれバグだったんだ…。
