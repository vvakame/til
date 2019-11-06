# TypeScript v3.7.2 変更点

更新 2019/11/06 正式版がリリースされて微妙に差分があったので更新
更新 2019/11/01 Call Hierarchy について更新

こんにちは[メルペイ社](https://www.merpay.com/)な[@vvakame](https://twitter.com/vvakame)です。

[TypeScript 3.7](https://devblogs.microsoft.com/typescript/announcing-typescript-3-7/)がアナウンスされました。

* [What's new in TypeScript in 3.7](https://www.typescriptlang.org/docs/handbook/release-notes/typescript-3-7.html)
* ~~[Breaking Changes in 3.7](https://github.com/microsoft/TypeScript/wiki/Breaking-Changes#typescript-37)~~
* [TypeScript 3.7 Iteration Plan](https://github.com/microsoft/TypeScript/issues/33352)
* [TypeScript Roadmap: July - December 2019](https://github.com/microsoft/TypeScript/issues/33118)

Roadmapは[こちら](https://github.com/Microsoft/TypeScript/wiki/Roadmap#37-november-2019)。

[この辺](https://github.com/vvakame/til/tree/master/typescript/v3.7.2)に僕が試した時のコードを投げてあります。

## 変更点まとめ

* [Optional Chaining](https://github.com/tc39/proposal-optional-chaining/)のサポート [Add support for Optional Chaining](https://github.com/microsoft/TypeScript/pull/33294)
    * stage 3になったので
    * `foo?.bar.baz()` とか書けるやつ
    * `?.` でワンセット
* [Nullish Coalescing](https://github.com/tc39/proposal-nullish-coalescing/)のサポート [nullish coalescing commit](https://github.com/microsoft/TypeScript/pull/32883)
    * stage 3になったので
    * `foo ?? bar()` とか書けるやつ
    * `||` と違って `null` と `undefined` しか相手にしない
* アサーションを行う関数のサポート [Assertion Functions](https://github.com/microsoft/TypeScript/pull/32695)
    * 関数の返り値に `asserts <expr>` 的なのが書けるようになる
    * ダメだったら `throw Error` とかさせる
* 関数の終端で `never` を返す関数を呼んだ時の型推論が賢くなった [Better Support for never-Returning Functions](https://github.com/microsoft/TypeScript/pull/32695)
    * `return process.exit(1)` ってやって `never` であることを伝えていたけど、それをしなくても分かってくれるようになった
* 再帰的な型エイリアスの制限のさらなる緩和 [(More) Recursive Type Aliases](https://github.com/microsoft/TypeScript/pull/33050)
    * サポート用のインタフェースが必要だった箇所で必要ないパターンが増えた
* `--declaration` と `--allowJs` が同時に使えるようになった [`--declaration` and `--allowJs`](https://github.com/microsoft/TypeScript/pull/32372)
    * マイグレーション途中のプロジェクトで便利
    * JSDocやコード実装から `.js` から `.d.ts` も頑張って出力される
* プロジェクト参照を使っている時ビルドフリーの編集が可能に [Build-Free Editing with Project References](https://github.com/microsoft/TypeScript/pull/32028)
    * プロジェクト参照を使ってコードを書いている時、 `.d.ts` ではなく `.ts` や `.tsx` をベースにエディターが動くようになる
    * コンパイル待ちのラグがなしに即時他のプロジェクトのコード変更が参照できて便利らしい
    * `disableSourceOfProjectReferenceRedirect` が tsconfig.json に追加
* 関数の存在チェックした後に呼んでなかったら怒ってくれる [Uncalled Function Checks](https://github.com/microsoft/TypeScript/pull/33178)
    * `if (obj.func)` とかした後に `obj.func()` してなかったら怒ってくれる
    * マジで確認だけして呼ばなくてよい場合は `if (!!obj.func)` とする
* `// @ts-nocheck` 導入 [`// @ts-nocheck` in TypeScript Files](https://github.com/microsoft/TypeScript/pull/33383)
    * JSからのマイグレーションの時に一時的に使うと便利〜〜みたいなやつ
* セミコロンをフォーマッタが 付ける/削除する を設定できるようになった [Semicolon Formatter Option](https://github.com/microsoft/TypeScript/pull/33402)
* ~~コールヒエラルキーのサポート [Call Hierarchy support](https://github.com/microsoft/TypeScript/issues/31863)~~
    * 主にVSCode側の機能っぽいけど…？
    * Find All References はすでにあるけどそれのことではないのか…？
    * Issueで聞いたらできてなくて大変なやつやで… という話になってた
* async内でawaitがいるような候補を選んだら自動的にawaitを挿入する [Auto-insert `await` for property accesses on Promise](https://github.com/microsoft/TypeScript/issues/31450)
    * 便利ですね
* より厳密なES仕様への追従と `useDefineForClassFields` オプションの追加 [Add useDefineForClassFields flag for Set -> Define property declaration](https://github.com/microsoft/TypeScript/pull/33509)
    * クラスのフィールドを定義した時の仕様がTypeScriptとESで異なる(歴史的経緯)
    * 一気に変えると色々ぶっ壊れそうなのでオプションでES準拠になるようにしたよ
    * 有効にするのを強く推奨するよ！
* 平らに均されたエラー表示 [Flatter Error Reporting](https://github.com/microsoft/TypeScript/pull/33473)
    * めっちゃネストした型エラーの表示とかがいい感じになるようになった

## 破壊的変更！

* DOMの変更 [DOM Changes](https://github.com/microsoft/TypeScript/pull/33710)
* `関数の存在チェックした後に呼んでなかったら怒ってくれる` のやつ [Function Truthy Checks](https://github.com/microsoft/TypeScript/pull/33178)
    * 改善のための破壊
* 他のモジュールのinterfaceと同名のinterfaceを作った時元のinterfaceの定義が拡張されるバグを修正 [Local and Imported Type Declarations Now Conflict](https://github.com/microsoft/TypeScript/pull/31231)
    * そんなんあったんか…
* API Changes
    * type aliasの再帰の改善のあおりで `TypeReference` から `typeArguments` が削除された 代わりに `TypeChecker#getTypeArguments` を使う

## Optional Chainingのサポート

[Optional Chaining](https://github.com/tc39/proposal-optional-chaining/)が stage 3 になりTypeScriptにも導入されました。

変数の値が null か undefined だった場合、評価を打ち切って undefined を返してくれるやつです。

```ts
let foo: any = {
    bar1: { buzz() { console.log("bar1"); } },
    bar2: void 0,
};

// bar1 と表示される
let x = foo?.bar1?.buzz();
// 何も表示されない
let y = foo?.bar2?.buzz();

// これはエラーになる
// error TS1109: Expression expected.
// ↓ 最後の ? が三項演算子だと思われてて面白い
// error TS1005: ':' expected.
// let z1 = foo?.bar?.buzz?();
// error TS1109: Expression expected.
// let z2 = foo?.bar?.buzz()?;

// ちなみにこれらはOK
// ?. でワンセット
let z3 = foo?.bar?.buzz?.();
let z4 = []?.[1];
```

あたかも、TypeScriptのoptionalと同じように `?` が導入されたように見えますが、実際に導入されたのは `?.` です。
`?` だけだと三項演算子と区別がつかないからですね。
そのため、 `buzz?.()` や `array?.[1]` のような一見珍奇な書き方をする必要があります。

もうちょっと例を見てみます。

```ts
let foo: any = { bar: { baz: true } };

// 今までのやり方
if (foo && foo.bar && foo.bar.baz) {
    console.log(foo.bar.baz);
}

// Optional Chainingを使うとこう書ける
if (foo?.bar?.baz) {
    console.log(foo.bar.baz);
}

// && と ? では厳密には挙動が異なる
// && は falsy な値 (null, undefined, "", 0, NaN, false) の場合処理を打ち切り、左辺の値を返す
// ? の場合、 null と undefined の時のみ処理を打ち切り、undefined を返す

// undefined と表示される (toStringは実行されないので)
console.log((null as any)?.toString());
// 実行時エラー Cannot read property 'toString' of null
console.log((null as any && true).toString());

function barPercentage(foo?: { bar: number }) {
    // こういうのもダメ foo?.bar の部分でエラーとなる
    // error TS2532: Object is possibly 'undefined'.
    // return foo?.bar / 100;
    // このように解釈されている
    // let tmp: number | undefined = (foo === null || foo === void 0) ? void 0 : foo.bar;
    // return tmp / 100;

    // こうすればOK
    return (foo?.bar ?? 0) / 100;
}
```

[prettierではまだこれをサポートしていない](https://github.com/prettier/prettier/issues/6595)ようです。
正式版までには使えるようになるといいですね。


## Nullish Coalescingにサポート

[Nullish Coalescing](https://github.com/tc39/proposal-nullish-coalescing/)が stage 3 になりTypeScriptにも導入されました。

変数の値が null や undefined のときに別の値を割り当てたいときに利用できます。

```ts
let foo: string | null = null as any;
let bar = "bar";

let a = foo ?? bar;
// bar と表示される foo が null なので
console.log(a);

foo = "" as any;
let b = foo ?? bar;
// 空文字列が表示される
// || と違って、null と undefined の時のみ右辺が評価される
// "" は当てはまらないので左辺の値が返る
console.log(b);

let c = foo || bar;
// bar と表示される
// "" は falsy な値なので 右辺が評価される
console.log(c);

// ?? と同じことをしてみる
let d = foo == null ? foo : bar;
// bar と表示される
// == null に当てはまるのは undefined と null のみ
console.log(d);
```

`||` の場合、falsyな値が対象ですが、 `??` の場合 null と undefined のときのみが対象になります。
つまり、 `0` や `""` は"存在する"ものとして扱われます。
基本的には、 `??` をメインに使い、faslyな値を潰したいときにのみ `||` を使うようにするのがよさそうです。


## アサーションを行う関数のサポート

`assert` などの、特定の条件下で例外を投げる関数に対するサポートが強化されました。
assert ではある変数が本当に制約を満たしているか？をチェックする用途で使う場合が多いでしょう。

この変更では、そこでチェックした内容をそれ以降のコントロールフローで利用できるようになります。

```ts
// asserts の後にどの仮引数が検査対象なのか書く
// この関数がエラーにならずに処理を返したら、someVariable は呼び出し元の型検査フローに対して正しい
function assert(someVariable: any, msg?: string): asserts someVariable {
    if (!someVariable) {
        // 例外を投げて処理の流れをぶった切る
        throw new Error(msg)
    }
}

function multiplyA(x: any, y: any) {
    // x, y が本当に number だったら assert は例外を投げない (という実装と型定義だった)
    assert(typeof x === "number");
    assert(typeof y === "number");

    // ここでは x と y はnumber型に絞られている
    return x * y;
}

function multiplyB(x: any, y: any) {
    // 今まではこうやって書いたりしていた
    // throw とかすると今までもControl Flow解析で x と y の型が定まっていた
    if (typeof x !== "number") {
        throw new Error();
    }
    if (typeof y !== "number") {
        throw new Error();
    }

    // ここでは x と y はnumber型に絞られている
    return x * y;
}


// この関数が true を返したら仮引数 val の型は string ですよというアレ(前からあるやつ)
// https://www.typescriptlang.org/docs/handbook/advanced-types.html#using-type-predicates で解説されている
function isString(val: any): val is string {
    return typeof val === "string";
}

// asserts の後に type predicates と同じ書き方をする
function assertIsString(val: any): asserts val is string {
    if (typeof val !== "string") {
        throw new Error("Not a string!");
    }
}

function usageC(str: string | null) {
    assertIsString(str);
    // assertIsString が 例外を投げなかったら str は string に絞られている
    str.toUpperCase();
}


function assertIsDefined<T>(val: T): asserts val is NonNullable<T> {
    if (val === undefined || val === null) {
        throw new Error(
            `Expected 'val' to be defined, but received ${val}`
        );
    }
}

function usageD(str: string | null) {
    assertIsDefined(str);
    // assertIsString が 例外を投げなかったら str から null の可能性が除外される
    str.toUpperCase();
}
```

便利ですね。
今まではこれができなかったがために、assert関数を使ってもあまり嬉しくなかったんですがこれが大幅に改善されました。


## 関数の終端で `never` を返す関数を呼んだ時の型推論が賢くなった

関数の終端で `never` を返す関数(Node.jsでいうとprocess.exitとか)を呼んだときの型推論が賢くなりました。

```ts
// この関数が値を返すことはない… (常に例外を投げるので)
function throwError(): never {
    throw new Error();
}

// TypeScript v3.6 ではコンパイルエラーになる
// error TS2366: Function lacks ending return statement and return type does not include 'undefined'.
// TypeScript v3.7 以降なら大丈夫
function multipler(v: any): string {
    if (typeof v === "string") {
        // 連結して2倍！
        return v + v;
    } else if (typeof v === "number") {
        // 2倍して2倍！(それはそう)
        return `${2 * v}`;
    }

    // v3.6 まではこう書くと あっ never ですね！返り値 string と矛盾しませんね！ ってなってた
    // return throwError();

    // v3.7 以降だとこれだけで あっ never ですね！ って伝わる
    throwError();
}
```

コンパイラのために余計な記述を行わなくてもよくなったので便利です。
筆者は今まで `throw new Error("unreachable")` とか書いてました…。


## 再帰的な型エイリアスの制限のさらなる緩和

今までは自分自身を参照するような構造を定義することができず、補助用のinterfaceなどを挟む必要がありました。
ここの評価が遅延されるようになったようで、循環構造を定義できるようになりました。

```ts
// TypeScript v3.6 までは直接自分自身を参照するような再帰構造は書けなかった
// error TS2456: Type alias 'Json' circularly references itself.
// TypeScript v3.7 以降は大丈夫
type Json =
    | string
    | number
    | boolean
    | null
    | { [property: string]: Json }
    | Json[];

let obj1: Json = 1;
let obj2: Json = "string";
let obj3: Json = {};
let obj4: Json = [];
let obj5: Json = {
    foo: [],
    bar: true,
};

{ // TypeScript v3.6 までは補助となるinterfaceとかの定義が必要だった
    type Json =
        | string
        | number
        | boolean
        | null
        | JsonObject
        | JsonArray;
    type JsonObject = {
        [property: string]: Json;
    };
    interface JsonArray extends Array<Json> { }
}
```

今までの書き方はなぜそうしなければいけないかが直感的ではなかったので、嬉しい変更です。

内部的には `TypeReference` から `typeArguments` が削除され、代わりに `TypeChecker#getTypeArguments` を使うようになっています。


## `--declaration` と `--allowJs` が同時に使えるようになった

らしいです。

```js
/**
 * Foo class.
 */
export class Foo {
    /**
     * @param {string} word
     * @returns {string}
     */
    bar(word) {
        return `Hello, ${word}`;
    }
}
```

こういうコードから

```ts
/**
 * Foo class.
 */
export class Foo {
    /**
     * @param {string} word
     * @returns {string}
     */
    bar(word: string): string;
}
```

こういう型定義ファイルが生成できます。
JSDocをしっかり書いていたプロジェクトであれば、かなりリッチな型定義ファイルが生成できそうです。


## プロジェクト参照を使っている時ビルドフリーの編集が可能に

らしいです。
[これ](https://qiita.com/vvakame/items/57a0559c45b88b2ae168#%E3%83%97%E3%83%AD%E3%82%B8%E3%82%A7%E3%82%AF%E3%83%88%E9%96%93%E3%81%AE%E5%8F%82%E7%85%A7%E3%81%AE%E3%82%B5%E3%83%9D%E3%83%BC%E3%83%88)の話題だと思うんですが試してみてもよくわかんなかったです…。
現在の VisualStudio Code - Insiders ではまだ使えない…？いやーでもtsserverが管理してるだろうしなぁ…？
謎です。

プロジェクト参照を使っている人は色々試してみてください。


## 関数の存在チェックした後に呼んでなかったら怒ってくれる

やりがちなミスなので嬉しいですね。

```ts
interface User {
    isAdministrator(): boolean;
    notify(): void;
    doNotDisturb?(): boolean;
}

function sudo() {
    console.log("exec sudu!");
}

// function doAdminThingA(user: User) {
//     // エラーになる！それ絶対存在するプロパティだから常にtrueなんだけど、ホントは呼び出したかったんじゃないの？
//     // error TS2774: This condition will always return true since the function is always defined. Did you mean to call it instead?
//     if (user.isAdministrator) {
//         sudo();
//     } else {
//         throw new Error("User is not an admin");
//     }
// }

function doAdminThingB(user: User) {
    // 当然、呼び出している場合はエラーにならない
    if (user.isAdministrator()) {
        sudo();
    } else {
        throw new Error("User is not an admin");
    }
}

function doAdminThingC(user: User) {
    // わざとだよ！という場合は !! として真偽値に変換することで意図を伝えることができる
    if (!!user.isAdministrator) {
        sudo();
    } else {
        throw new Error("User is not an admin");
    }
}

function doAdminThingD(user: User) {
    if (user.notify) {
        // その後、呼び出すならOK
        user.notify();
    }
    if (user.doNotDisturb) {
        // doNotDisturb は optional なのでOK
        sudo(); // 現実的にはOKじゃないかもね！
    }
}
```

かしこいですね。


## `// @ts-nocheck` 導入

JSからの移行で便利なやつです。
とりあえず拡張子を .js から .ts にしてしまって `@ts-nocheck` つければコンパイルは通る！

```ts
// @ts-nocheck

// やり放題だぜーーーーっ！！
class Foo {
}
// bar なんか存在しないぜーーーーっ！！
new Foo().bar();
```

JSDocをちゃんと書いてから移行するのが面倒な人はこっちのほうが手っ取り早そうです。


## セミコロンをフォーマッタが 付ける/削除する を設定できるようになった

これも試してみたんですがうまく動作しませんでした。
VSCode - Insiders 上に設定項目は存在していて、 `"typescript.format.semicolons": "insert"` みたいになるんですがフォーマッタ適用してもセミコロンが自動でついたりはしませんでした…。


## ~~コールヒエラルキーのサポート~~

これも謎です。
Find All References とは違うものなのか…？
[Issue](https://github.com/microsoft/TypeScript/issues/31863)にもほぼ情報がありません。

…と思ったんですがなんかまだできてなかったっぽいです。
詳細はIssueを見てみてください。

## async内でawaitがいるような候補を選んだら自動的にawaitを挿入する

これも謎です。
今回謎が多いな… 僕の検証方法が悪い可能性が微粒子レベルで存在している…？

```ts
async function asyncFunc(v: Promise<string>) {
    // v. とタイプすると then とか catch の他に toLowerCase などが候補に出るはずだが…？ :thinking_face:
    // v.
}
```

というわけで、これはちゃんと動いたら割と便利なはずのやつです。

## より厳密なES仕様への追従と `useDefineForClassFields` オプションの追加

v3.7.2 で追加されました。

```ts
class C {
    foo = 100;
    bar?: string;
}
```

これが今までは

```js
"use strict";
class C {
    constructor() {
        this.foo = 100;
    }
}
```

こうなっていたけど、ESの仕様に沿うのであれば

```js
"use strict";
class C {
    constructor() {
        Object.defineProperty(this, "foo", {
            enumerable: true,
            configurable: true,
            writable: true,
            value: 100
        });
        Object.defineProperty(this, "bar", {
            enumerable: true,
            configurable: true,
            writable: true,
            value: void 0
        });
    }
}
```

となるのが正しいらしいです。
ほんまか。

`useDefineForClassFields` が追加され、これを true にすると"本来あるべき"コードが生成されるようになります。
デフォルトは false で、今後のバージョンで(ES 3.8とか該当仕様がstage 4になるとか)デフォルトtrueになる予定のようです。

これによって色々と非互換な振る舞いの変更が発生するため、true にした後はテストを普段より念入りにやったほうがよいでしょう。
TypeScript team的には true を試してみることを推奨しています。

問題になる振る舞いの例を少し取り上げてみます。
基本的には問題のある場合はコンパイルエラーになって、明示的な対処を求められるようです。

```ts
class Base {
    set data(value: string) {
        console.log("data changed to " + value);
    }
}

class Derived1 extends Base {
    // ES仕様に準拠した振る舞い(useDefineForClassFieldsを有効にした場合)では、 set accessor が上書きされる
    // 結果、 data changed to ... の出力が行われなくなる
    // これを防ぐため、コンパイルエラーが発生するようになった
    // error TS2610: Class 'Base' defines instance member accessor 'data', but extended class 'Derived1' defines it as instance member property.
    // data = "foobar";

    constructor() {
        super();
        // エラーを解消するためには、コンストラクタ内で値を初期化すればよい
        this.data = "foobar";
    }
}

let obj = new Derived1();
obj.data = "fizzbuzz";
```

```ts
interface Animal { animalStuff: any }
interface Dog extends Animal { dogStuff: any }

class AnimalHouse {
    resident?: Animal;
    constructor(animal: Animal) {
        this.resident = animal;
    }
}

class DogHouse extends AnimalHouse {
    // useDefineForClassFields を使っている場合、値が undefined になる！
    // super 呼ぶ → AnimalHouse で resident がセットされる → DogHouse のfield initialize が走る → undefined に再設定される！
    // というのを防ぐため、コンパイルエラーになる
    // error TS2612: Property 'resident' will overwrite the base property in 'AnimalHouse'. If this is intentional, add an initializer. Otherwise, add a 'declare' modifier or remove the redundant declaration.
    // resident?: Dog;

    // 解消方法
    // declare を追加するとコード生成に関与しなくなる
    declare resident?: Dog;

    constructor(dog: Dog) {
        // super に Animal な値を渡すこともできてしまう
        // resident?: Dog が正しくなるかどうかはプログラマ次第…
        super(dog);

        // declare 使うのやめて、素直に自分で初期化したほうがいいかもね
        // this.resident = dog;
    }
}

let obj = new DogHouse({animalStuff: 1, dogStuff: 2});
// undefined と表示される useDefineForClassFields を有効にしている場合
console.log(obj.resident);
```

この一連の問題は[ここ](https://github.com/microsoft/TypeScript/issues/27644)でトラッキングされています。


## 平らに均されたエラー表示

割と便利なやつです。
型に関する整合性のエラーは入り組んだ、深くネストされた形になることも珍しくありません。
そのネストされたエラーがなるべく平易に表示されるようになったそうです。

```ts
type SomeVeryBigType = { a: { b: { c: { d: { e: { f(): string } } } } } }
type AnotherVeryBigType = { a: { b: { c: { d: { e: { f(): number } } } } } }

declare let x: SomeVeryBigType;
declare let y: AnotherVeryBigType;

// 短くてわかりやすいエラー！
// error TS2322: Type 'SomeVeryBigType' is not assignable to type 'AnotherVeryBigType'.
//   The types returned by 'a.b.c.d.e.f()' are incompatible between these types.
//   Type 'string' is not assignable to type 'number'.
// y = x;

// ちなみに今までのエラーはこんな感じ
// TSユーザは経験的に型関連のエラーは一番下から見ていけばいいと知っているし、この場合は下4行を読むと原因がわかる
// けど、長いものは長い。
// error TS2322: Type 'SomeVeryBigType' is not assignable to type 'AnotherVeryBigType'.
//   Types of property 'a' are incompatible.
//     Type '{ b: { c: { d: { e: { f(): string; }; }; }; }; }' is not assignable to type '{ b: { c: { d: { e: { f(): number; }; }; }; }; }'.
//       Types of property 'b' are incompatible.
//         Type '{ c: { d: { e: { f(): string; }; }; }; }' is not assignable to type '{ c: { d: { e: { f(): number; }; }; }; }'.
//           Types of property 'c' are incompatible.
//             Type '{ d: { e: { f(): string; }; }; }' is not assignable to type '{ d: { e: { f(): number; }; }; }'.
//               Types of property 'd' are incompatible.
//                 Type '{ e: { f(): string; }; }' is not assignable to type '{ e: { f(): number; }; }'.
//                   Types of property 'e' are incompatible.
//                     Type '{ f(): string; }' is not assignable to type '{ f(): number; }'.
//                       Types of property 'f' are incompatible.
//                         Type '() => string' is not assignable to type '() => number'.
//                           Type 'string' is not assignable to type 'number'.
```
