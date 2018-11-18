こんにちは[メルペイ社](https://www.merpay.com/jp/)な[@vvakame](https://twitter.com/vvakame)です。

[TypeScript 3.2 RC](https://blogs.msdn.microsoft.com/typescript/2018/11/15/announcing-typescript-3-2-rc/)がアナウンスされました。

[What's new in TypeScriptも更新](https://github.com/Microsoft/TypeScript/wiki/What's-new-in-TypeScript#typescript-32)されて…いません！！ナンデ？？？
[破壊的変更](https://github.com/Microsoft/TypeScript/wiki/Breaking-Changes#typescript-32)はアップデートされてました。

[この辺](https://github.com/vvakame/til/tree/master/typescript/v3.2.0-rc)に僕が試した時のコードを投げてあります。

## 変更点まとめ

* Functionのbind, call, apply周りの型推論が更に賢くなった [strictBindCallApply](https://github.com/Microsoft/TypeScript/pull/27028)
  * bind, call, apply についての型推論がさらにいい感じになった
  * `--strictBindCallApply` が導入された
  * `--strict` に含まれる
  * `CallableFunction` と `NewableFunction` という組み込み型も導入
* Genericsを使った型をObject spreadでいい感じに扱えるようになった [Object spread on generic types](https://github.com/Microsoft/TypeScript/pull/28234)
  * `{...T, ...U}` 的概念の操作ができるようになった
  * `object spread type` と言ってるらしい
* Genericsが絡んだrest parameterで余り物の型がいい感じに出るようになった [Object rest on generic types](https://github.com/Microsoft/TypeScript/pull/28312)
  * `let {x, y, ...rest} = obj` 的なやつのrest部分
* BigIntサポートの追加 [BigInt Support](https://github.com/Microsoft/TypeScript/issues/15096)
  * [これヤツ](https://tc39.github.io/proposal-bigint/)
* unionの式による判別の範囲のnon-unit typesにまで拡張 [Allow non-unit types in union discriminants](https://github.com/Microsoft/TypeScript/pull/27695)
  * undefinedかどうかとかを見てunionのどの型か判別できるようになった
  * unit types === literal, symbol, nullable
* tsconfig.jsonの `extends` が普通のパッケージのようなルールで解決されるようになった [Configuration inheritance through node packages](https://github.com/Microsoft/TypeScript/pull/27348)
  * `./node_modules/foobar/tsconfig.json` とか書いてたのを `foobar/tsconfig.json` とか書けるようになった
  * 今まで書けなかったんだっけ…？
* 最終的に使われるtsconfig.jsonをダンプする [Support printing the implied configuration object to the console with `--showConfig`](https://github.com/Microsoft/TypeScript/pull/27353)
  * extendsとかも取りまとめた結果を出してくれる
* [Supporting Object.defineProperty property assignments in JS](https://github.com/Microsoft/TypeScript/pull/27208)
  * salsa系の話題なのでスルー(TSには関係ない)
* フォーマッタの改善 [Improved formatting and indentation for lists and chained calls](https://github.com/Microsoft/TypeScript/pull/28340)
* 型定義が存在しないパッケージに対して自分で `@types` するための雛形を生成 [Scaffold local `@types` packages with dts-gen](https://github.com/Microsoft/TypeScript/issues/25746)
  * 便利そう
* 非道なtype assertion入れようとした時に unknown を挟むCodeFix [Add intermediate unknown type assertions](https://github.com/Microsoft/TypeScript/issues/28067)
  * `0 as string` とかすると `0 as unknown as string` ってしないとダメだよ〜って教えてくれる
* new付け忘れてそうな時に教えてくれるCodeFix [Add missing `new` keyword](https://github.com/Microsoft/TypeScript/issues/26580)
  * ままありそう
* JSXと普通の関数で呼び出しを処理する部分を共通化した [JSX resolution changes](https://github.com/Microsoft/TypeScript/pull/27627)
  * 互換性はなんも壊れてないはずだけどやらかしを見つけたら教えてください だそうです
* lib.d.ts changes
  * lib.d.ts 周りの改善が続けられている
  * More specific types
  * More platform-specific deprecations
  * wheelDelta and friends have been removed.


### 破壊的変更！

上記リストのうち、破壊的変更を伴うのは次のものになります。

* lib.d.ts changes

## Functionのbind, call, apply周りの型推論が更に賢くなった

`--strictBindCallApply` が導入されました。
`--strict` に含まれているので `--strict` 使っておきましょう。

`--strictBindCallApply` を有効にしておくと、Functionのbind, call, applyなどの型の推論がより賢く行われるようになります。
次のコード例は、今までは全てコンパイルエラーとして検出できなかったものです。

```ts
function foo(a: number, b: string): string {
    return a + b;
}

// 引数が足りないのでエラー であることが検出できる！
let a = foo.apply(undefined, [10]);
// 2番目の引数があわないのでエラー  であることが検出できる！
let b = foo.apply(undefined, [10, 20]);
// 引数の数が多すぎるのでエラー であることが検出できる！
let c = foo.apply(undefined, [10, "hello", 30]);
// OKなパターン
let d = foo.apply(undefined, [10, "hello"]);

// 引数あわない
let e = foo.call(undefined, 1, 1);
// OK！
let f = foo.call(undefined, 1, "string");

// 1つ目の引数を固定する
let g = foo.bind(undefined, 1);
// NGであることが検出できる！
g(true);
// OK
g("foobar");
```

この仕組を実現するために、 `CallableFunction` と `NewableFunction` という型が導入されました。
普通の関数がCallableFunctionで、classのconstructorとかがNewableFunctionになります。
`--strictBindCallApply` を有効にすると、これらのオブジェクトの型が Function ではなくそれぞれのものに振り分けられるようになります。

call, applyなどは使わないほうが良いですが、bindの型がしっかりつくようになったのは嬉しいですね。

## Genericsを使った型をObject spreadでいい感じに扱えるようになった

こんなん。

```ts
function merge<T, U>(x: T, y: U) {
    // 今まではこの書き方はエラーになってた
    // error TS2698: Spread types may only be created from object types.
    // 3.2からイケる！
    return { ...x, ...y };
}

let obj = merge({ name: "vv" }, { cat: "yukari" });
```

これを `object spread type` とか言ってるらしいです。
[この辺](https://github.com/Microsoft/TypeScript/issues/10727)でユースケースについていろいろ検討されたそうな。

便利ですね。

## Genericsが絡んだrest parameterで余り物の型がいい感じに出るようになった

上のものとワンセット的に、rest parameterにもGenerics対応が入りました。
こんな感じ。

```ts
interface XYZ { x: any; y: any; z: any; }

// T から Pickする 対象は TのプロパティからXYZのプロパティと重複する部分を除外したもの
type DropXYZ<T> = Pick<T, Exclude<keyof T, keyof XYZ>>;

function dropXYZ<T extends XYZ>(obj: T): DropXYZ<T> {
    let { x, y, z, ...rest } = obj;
    return rest;
}

let obj = dropXYZ( { w: null, x: 1, y: "str", z: true });

// OK
obj.w;
// 存在しないのでコンパイルエラー！
// error TS2339: Property 'x' does not exist on type 'Pick<{ w: null; x: number; y: string; z: boolean; }, "w">'.
// obj.x;
// error TS2339: Property 'y' does not exist on type 'Pick<{ w: null; x: number; y: string; z: boolean; }, "w">'.
// obj.y;
// error TS2339: Property 'z' does not exist on type 'Pick<{ w: null; x: number; y: string; z: boolean; }, "w">'.
// obj.z;
```

## BigIntサポートの追加

[BigInt](https://tc39.github.io/proposal-bigint/)にサポートが入りました。
使うには `"target": "esnext"` が必要です。
また、 tsconfig.jsonのlib用に `"esnext.bigint"` も追加されています。

早速試していきましょう。

```ts
let bi1: bigint = 1000n;

// 新しく導入された型
type f = bigint | BigInt | BigInt64Array | BigUint64Array;


let bi2 = BigInt(1);
let mod1 = BigInt.asIntN(64, 10000000000000000000n);
let mod2 = BigInt.asUintN(64, 10000000000000000000n);

// 1000n 1n -8446744073709551616n 10000000000000000000n と表示される
console.log(bi1, bi2, mod1, mod2);

// 四則演算
console.log(1n + 2n);
// numberと比較するやつ
// error TS2367: This condition will always return 'false' since the types '1n' and '1' have no overlap.
// console.log(1n == 1);
console.log(1n < 2);
```

bigintとnumberは非互換です。
また、`typeof` によるtype narrowingもすでにサポートされています。

## unionの式による判別の範囲のnon-unit typesにまで拡張

短い日本語で表すのが難しいですね…。
もう細かい仕様を僕は覚えられてないんですが、union typesな値が実際何であるか式によって判別することができます。
今まではunit typesなプロパティしか判別に使えなかったのですが、それが拡張されました。
union typesとは、コンパイラのコードをシュッと見た感じでは literal, symbol, nullable を指すようです。

今回は次のようなコードで正しくどの型か判別できるようになりました。

```ts
type Result<T> = { error?: undefined, value: T } | { error: Error };
// 昔もこんな感じだったらまぁイケた
// type Result<T> = { error: false, value: T } | { error: true };

function test(x: Result<number>) {
    if (!x.error) {
        // 今まではこのやり方でtype narrowingできなかった
        // error TS2339: Property 'value' does not exist on type 'Result<number>'.
        //   Property 'value' does not exist on type '{ error: Error; }'.
        // error が true | false だったりとunit typeである必要があった…
        // だがしかし今はこれでイケる！
        x.value;  // number
    }
    else {
        x.error.message;  // string
    }
}

test({ value: 10 });
test({ error: new Error("boom") });
```

よいですね。
そろそろ、TypeScriptコンパイラのサポート範囲を考慮せずに、自由にJSっぽいコードを書いてもちゃんと型が効くようになってきたのではないでしょうか。

## tsconfig.jsonの `extends` が普通のパッケージのようなルールで解決されるようになった

今まで `"extends": "./node_modules/foo-config"` とか書いてたやつが `"extends": "foo-config"` で読み込めるようになった、という感じのようです。
会社で共通設定ファイルを作っている人は多少設定ファイルがきれいになるのではないでしょうか。

## 最終的に使われるtsconfig.jsonをダンプする

`--showConfig` が追加されました。

```tsconfig.json
{
    "extends": "./other.tsconfig.json",
    "compilerOptions": {
        "target": "esnext",
        "module": "commonjs",
        "lib": [
            "esnext"
        ],
        "strict": true,
        "esModuleInterop": true
    }
}
```

```other.tsconfig.json
{
    "compilerOptions": {
        "noUnusedLocals": true,
        "noUnusedParameters": true,
        "noImplicitReturns": true,
        "noFallthroughCasesInSwitch": true
    }
}
```

という状況だと

```
$ tsc --showConfig
{
    "compilerOptions": {
        "noUnusedLocals": true,
        "noUnusedParameters": true,
        "noImplicitReturns": true,
        "noFallthroughCasesInSwitch": true,
        "target": "esnext",
        "module": "commonjs",
        "lib": [
            "esnext"
        ],
        "strict": true,
        "esModuleInterop": true
    },
}
```

となります。
`extends` 適用後の設定が見られるのはデバッグ的が楽になりますね。

この時、 `"extends"` の指定先が間違っている場合、 `--showConfig` では単に無視されてエラーになりませんので注意が必要です。
これは不便だと思うので明日あたりにIssueとして報告しておこうと思います。

## Supporting Object.defineProperty property assignments in JS のやつ

JSをtscとかで処理した時の振る舞いで興味ないのでカットだ！

## フォーマッタの改善

なんか改善されたらしいです。
リストやらメソッドチェインのインデントがよくなったっぽいです。
この辺は不満があったので嬉しいところですね。

## 型定義が存在しないパッケージに対して自分で `@types` するための雛形を生成

`require` を撲滅するために、 `@types` すらなかった場合[dts-gen](https://github.com/Microsoft/dts-gen/)的なものを使って対象パッケージの型定義のスケルトンを生成する的なやつのようです。
まだ、internalなAPIでエディタから叩くことはできない…ように見えるのですがもしできたら筆者までお知らせください(他力本願)。

## 非道なtype assertion入れようとした時に unknown を挟むCodeFix

`0 as string` とかしようとすると `error TS2352: Conversion of type 'number' to type 'string' may be a mistake because neither type sufficiently overlaps with the other. If this was intentional, convert the expression to 'unknown' first.` と怒られます。
長文指摘おじさんじゃん。

ここで `Add 'unknown' conversion for non-overlapping types` のCodeFixを選ぶと `0 as unknown as string;` に変換してくれます。
世界の平和を乱す行為ですが実際必要になることがあるのも事実…。

この修正は `good first issue` が付けられていて、「これ初心者向けってマジ？全然わからんのやけどポインタちょうだい」みたいな質問をする人が現れ、MSの人とキャッキャうふふしつつ実装された機能です。
いい話。

## new付け忘れてそうな時に教えてくれるCodeFix

こういうのを修正してくれる… やつのはずですが、3.2.1に入るっぽくてまだ使えませんでした。
まま便利そう。

```ts
class C {}

// error TS2348: Value of type 'typeof C' is not callable. Did you mean to include 'new'?
C();
```

## JSXと普通の関数で呼び出しを処理する部分を共通化した

なんかリグレッションあったら教えてくれ〜〜って感じっぽいです。
既存のテストは変更無しで通ったということらしいので、まぁ大丈夫でしょきっと…(知らんけど)

## lib.d.ts changes

ちょいちょい変わっていってるみたいですね。
DOM周りで `null` を受け付けないべきであるものがそうなったり、interface名がspecに沿うものになったりしていってるっぽいです。
`WheelEvents` の `wheelDelta` とかが削除されたりもしているそうです。
