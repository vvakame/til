# TypeScript v3.4.0-rc 変更点

こんにちは[メルペイ社](https://www.merpay.com/jp/)な[@vvakame](https://twitter.com/vvakame)です。

[TypeScript 3.4 RC](https://devblogs.microsoft.com/typescript/announcing-typescript-3-4-rc/)がアナウンスされました。

[What's new in TypeScriptは消滅](https://github.com/Microsoft/TypeScript/wiki/What's-new-in-TypeScript#typescript-34)したようです。
Roadmapは[こちら](https://github.com/Microsoft/TypeScript/wiki/Roadmap#34-march-2019)。
v3.4.0での[破壊的変更](https://github.com/Microsoft/TypeScript/wiki/Breaking-Changes#typescript-34)はこちら。

今回から進捗管理の方法に変化があるようですね。
外部から変更を後追いしやすくなった印象です。

* [TypeScript 3.4 Iteration Plan](https://github.com/Microsoft/TypeScript/issues/30281)
* [TypeScript Roadmap: January - June 2019](https://github.com/Microsoft/TypeScript/issues/29288)

とかがあります。
かなりボリュームがあるので逐次追うのはかなりMPが必要そうですね…。

[この辺](https://github.com/vvakame/til/tree/master/typescript/v3.4.0-rc)に僕が試した時のコードを投げてあります。

## 変更点まとめ

* `--incremental` の追加 [`--incremental` builds with `.tsbuildinfo` files](https://github.com/Microsoft/TypeScript/pull/29813)
    * `--watch` と違ってプロセスを跨いで(cold build)インクリメンタルビルド可能に
    * `.tsbuildinfo` に出力がある
    * tsconfig.jsonに `incremental` と `tsBuildInfoFile` 追加
* `ReadonlyArray` と `readonly` tupleの改善 [Improved support for read-only arrays and tuples](https://github.com/Microsoft/TypeScript/pull/29435)
    * 新しい記法の追加
        * 配列に対しての記法
        * tupleに対しての記法
    * `readonly` を使ったmapped typeの挙動の変更
* `const` assertionの追加 [Const contexts for literal expressions](https://github.com/Microsoft/TypeScript/pull/29510)
    * `as const` とか `<const>` とか書ける
    * リテラルで書いた値の型がliteral typesになるみたいなやつ
* `globalThis` に型がついた [`globalThis`](https://github.com/Microsoft/TypeScript/pull/29332)
    * Top-levelの `this` にも型がついた
* Genericsの型の推論が強化された [Higher order function type inference](https://github.com/Microsoft/TypeScript/pull/30215) [*1](https://github.com/Microsoft/TypeScript/pull/30114) [*2](https://github.com/Microsoft/TypeScript/pull/29478)
    * *1 と *2 も十把一絡げ扱いされてるっぽい
* リファクタリングの改善
    * 関数へ複数の引数があるパターンをOptional型に変換するリファクタリングの追加 [Convert to named parameters](https://github.com/Microsoft/TypeScript/pull/30089)
* TypeScript本体のビルドが jake から gulp に変更になった [Remove jake (hopefully for real this time)](https://github.com/Microsoft/TypeScript/pull/29085)
    * ちょっとショック

## 破壊的変更！

* Top-levelの `this` にも型がついた
* Genericsの型の推論が強化された

に既存のコードと非互換なパーツが含まれています。

## `--incremental` の追加

`--incremental` と `--tsBuildInfoFile` がtsconfig.jsonに追加になっています。
これはビルド時に `*.tsbuildinfo` 的なものを出力して、ビルド間で情報を引き継ぎビルドを早くしよう的なものです。

tsconfig.jsonに追加して利用します。

```json
{
  "compilerOptions": {
    "target": "es5",
    "module": "commonjs",
    "incremental": true,
    "tsBuildInfoFile": "./.tsbuildinfo",
    "strict": true,
    "esModuleInterop": true
  }
}
```

手元で適当に試した限りでは、index.ts単体のプロジェクトでも 有効: 1.5秒 無効: 3.2秒 程度の差が出ています。
プロジェクトが育つにつれ恩恵もでかくなると思いますのでとりあえず試してみるとよいのではないでしょうか。
ただし、筆者が手元で試した時に修正したはずのコンパイルエラーが何故か治らない場合などがあったため、挙動がおかしい、と感じたら手で `.tsbuildinfo` を消す作業をする必要があります。（なお `@next` で試すと発生しないため次のバージョンでは解消されてそう）

`tsBuildInfoFile` の指定なしのデフォルトでは tsconfig.json に対して tsconfig.tsbuildinfo が生成されます。
`tsc -p tsconfig.test.json` とかすると `tsconfig.test.tsbuildinfo` が生成されます。
一般的には `tsBuildInfoFile` 指定無しで運用で良さそうです。
本記事でも特段の断りがない場合、指定無しで使った場合の解説とします。

`tsc --outDir dist` とかすると `dist/tsconfig.tsbuildinfo` に生成されます。

`tsc --module system --out dist/foo.js` とかすると `dist/foo.tsbuildinfo` に生成されます。

[composite project](https://qiita.com/vvakame/items/57a0559c45b88b2ae168#%E3%83%97%E3%83%AD%E3%82%B8%E3%82%A7%E3%82%AF%E3%83%88%E9%96%93%E3%81%AE%E5%8F%82%E7%85%A7%E3%81%AE%E3%82%B5%E3%83%9D%E3%83%BC%E3%83%88)を使った場合、 `--incremental` は自動的に `true` になります。

tsbuildinfoファイルの中身は次のようなJSONファイルです。

```json
{
  "program": {
    "fileInfos": {
      "/users/vvakame/.nodebrew/node/v11.10.1/lib/node_modules/typescript/lib/lib.d.ts": {},
      "/users/vvakame/.nodebrew/node/v11.10.1/lib/node_modules/typescript/lib/lib.es5.d.ts": {},
      "/users/vvakame/.nodebrew/node/v11.10.1/lib/node_modules/typescript/lib/lib.dom.d.ts": {},
      "/users/vvakame/.nodebrew/node/v11.10.1/lib/node_modules/typescript/lib/lib.webworker.importscripts.d.ts": {},
      "/users/vvakame/.nodebrew/node/v11.10.1/lib/node_modules/typescript/lib/lib.scripthost.d.ts": {},
      "/users/vvakame/work/til/typescript/v3.4.0-rc/src/incrementalbuild/foobar.ts": {
        "signature": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
      },
      "/users/vvakame/work/til/typescript/v3.4.0-rc/src/incrementalbuild/index.ts": {
        "signature": "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"
      }
    },
    "options": {
      "target": 1,
      "module": 1,
      "incremental": true,
      "strict": true,
      "esModuleInterop": true,
      "configFilePath": "/Users/vvakame/work/til/typescript/v3.4.0-rc/tsconfig.json"
    },
    "referencedMap": {},
    "exportedModulesMap": {},
    "semanticDiagnosticsPerFile": [
      "/users/vvakame/.nodebrew/node/v11.10.1/lib/node_modules/typescript/lib/lib.d.ts",
      "/users/vvakame/.nodebrew/node/v11.10.1/lib/node_modules/typescript/lib/lib.es5.d.ts",
      "/users/vvakame/.nodebrew/node/v11.10.1/lib/node_modules/typescript/lib/lib.dom.d.ts",
      "/users/vvakame/.nodebrew/node/v11.10.1/lib/node_modules/typescript/lib/lib.webworker.importscripts.d.ts",
      "/users/vvakame/.nodebrew/node/v11.10.1/lib/node_modules/typescript/lib/lib.scripthost.d.ts",
      "/users/vvakame/work/til/typescript/v3.4.0-rc/src/incrementalbuild/foobar.ts",
      "/users/vvakame/work/til/typescript/v3.4.0-rc/src/incrementalbuild/index.ts"
    ]
  },
  "version": "3.4.0-rc"
}
```

完全に他の人とはシェアしないタイプのやつですね。
素直に `.gitignore` に `*.tsbuildinfo` 追加でよさそうです。
`.tsbuildfile` はいつでも好きな時に消してしまって大丈夫です。

`--incremental` と `--tsBuildInfoFile` と書いてあるのでオプションとしても指定できそうに見えますが `tsc --incremental` などとすると

```bash
$ tsc --incremental
error TS6064: Option 'incremental' can only be specified in 'tsconfig.json' file.
```

と怒られるので素直に `tsconfig.json` に設定しましょう。

`.tsbuildfile` の書式については特に何も保証はしない、とのことなので今後のTypeScriptのバージョンアップに伴い出力される内容も変化していくものと思われます。

## `ReadonlyArray` と `readonly` tupleの改善

新しく `readonly` という修飾子が追加されたのと、それが配列やタプルの各要素に適用可能になりました。

`number[]` が `Array<number>` と等価であるように、 `readonly numer[]` が `ReadonlyArray<number>` と等価な表現として使えます。
また、 `readonly` はタプルにも適用できます。

配列での例とタプルでの例を見てみましょう。

```ts
// この2つは等価 ReadonlyArray は前からあったやつ
let readonlyArray1: readonly string[] =  ["a", "b", "c"];
let readonlyArray2: ReadonlyArray<string> =  ["a", "b", "c"];

// これは怒られる
// error TS1354: 'readonly' type modifier is only permitted on array and tuple literal types.
// let readonlyArray3: readonly Array<string> =  ["a", "b", "c"];

// 過去のおさらい

let array1 =  ["a", "b", "c"];

readonlyArray1.slice(); // OK
array1.slice(); // OK

// NG! 破壊的変更があるメソッドは ReadonlyArray には存在しない
// error TS2339: Property 'push' does not exist on type 'readonly string[]'.
// readonlyArray1.push('d');
// error TS2339: Property 'sort' does not exist on type 'readonly string[]'.
// readonlyArray1.sort();

// 普通のarrayならもちろんOK
array1.push('d');
array1.sort();

// NG! まぁ普通に代入もだめ
// error TS2542: Index signature in type 'readonly string[]' only permits reading.
// readonlyArray1[0] = "d";
```

```ts
let readonlyTuple1: readonly [string, string] = ["a", "b"];
// ↓の書き方は前からできたけど readonlyTuple2[0] = "b" とかから守ってくれなかった（後述）
let readonlyTuple2: Readonly<[string, string]> = ["a", "b"];

// 両方NG!
// error TS2540: Cannot assign to '0' because it is a read-only property.
// readonlyTuple1[0] = "a";
// readonlyTuple2[0] = "b";
```

はい。
この変更に伴い、mapped typesでの挙動にも変更がありました。
型の計算を行う時に、Arrayに対しての操作は要素への操作に展開されていました。
ただし、readonlyについてはこの限りではなかったのですが、今回その制限が撤廃され、適切なルールが定義された形になります。

```ts
// mapped typesの挙動が変わった

// 今までもこれからも (string | undefined)[]
// Arrayに対してMapped Typesを適用すると各要素に適用された
type A0 = Partial<string[]>;

// Readonlyについては上記のルールは当てはまらなかった！が、今回から適用されるようになった
// これから readonly string[]
// これまで string[]
type A1 = Readonly<string[]>;

// -readonly による属性剥がしも同様
type Writable<T> = {
    -readonly [K in keyof T]: Writable<T[K]>;
}

// これから string[]
// これまで ReadonlyArray<any>
type A2 = Writable<ReadonlyArray<string>>;
```

Apolloで黒魔術ごっこしてるとReadonlyArrayを剥がしたくなることがあるので、これは嬉しいですね。

なお、arrayかtupleのリテラル型以外に対して `readonly` をつけると怒られるので気をつけましょう。

```ts
// やると怒られる系
// readonly Array<numer> とかが許されないのはなんとなく不便なので将来的には治るのではなかろうか（てきとう）
// error TS1354: 'readonly' type modifier is only permitted on array and tuple literal types.
// let err1: readonly Set<number>; // error!
// let err2: readonly Array<number>; // error!
```

## `const` assertionの追加

`as const` または `<const>` のような型アサーションの記法で、リテラルの類の型をリテラルに沿った型リテラルとして扱い、readonlyにする記法です。
複雑にネストしたリテラルにも適用することができます。

例を見てみましょう。

```ts
// as const でリテラルを具体的なreadonlyなオブジェクト型リテラル相当の表現に変換できる
// a は 10 型
let a = 10 as const;
// const の場合は昔から 10 型
const a1 = 10;

// 型アサーションと同じ記法なので前置の書き方もできる（この書き方使わなくなりましたねぇ…）
// b は readonly [10, 20] 型
let b = <const>[10, 20];

// オブジェクトリテラルにも適用できる 配列にもOK
// c1 は { readonly text: "hello" ; } 型
let c1 = { text: "hello" } as const;  // Type { readonly text: "hello" }
// c2 は [true, false] 型
let c2 = [true, false] as const;


// こういう複雑なオブジェクトもconstにできる
let d1 = { lunch: "saizeriya" };
let d = {
    name: "vvakame",
    love: {
        kind: "cat",
        name: "yukari",
    },
    location: "tokyo",
    note: d1,
} as const;

// NG! これは怒られる
// error TS2540: Cannot assign to 'note' because it is a read-only property.
// d.note = { lunch: "CoCo壱" };

// ここはreadonlyではない
d.note.lunch = "rigoletto";
```

なんでこんなもんがいるんだ…？という気もしますが、次のような用途に使えるようです。

```ts
// 実用例 as const 無しだと array1 は { kind: string; language?: string[]; endpoints?: string[]; } 的な型になってしまう
let array1 = [
    { kind: "AppEngine", services: ["default", "worker"] },
    { kind: "Cloud Functions", endpoints: ["Hello", "Bye"] },
] as const;
for (let value of array1) {
    // 各要素の持つ値がはっきりしているのでtype narrowingで安全にアクセスできる
    if (value.kind === "AppEngine") {
        value.services.forEach(v => console.log(v));
    } else {
        value.endpoints.forEach(v => console.log(v));
    }
}

// 既存の何かの型にあわせるみたいなのもできる
type CloudService = { kind: "AppEngine"; services: readonly string[]; } | { kind: "Cloud Functions", endpoints: readonly string[] };
let services: ReadonlyArray<CloudService> = array1;

// 素直にこう書けばよくない？という説もなきにしもあらず
let array2: CloudService[] = [
    { kind: "AppEngine", services: ["default", "worker"] },
    { kind: "Cloud Functions", endpoints: ["Hello", "Bye"] },
];
```

ダメパターンも紹介しておきます。
基本的にはリテラルに対してのみ利用可能で、なんらかの計算を挟む(適切に推論しないと導出できないような場合)とエラーになります。

```ts
// 型アサーション自体は値の世界の住人なので型の世界で利用することはできない
// …というかその必要がないよね
// type A1 = { name: string; } as const;
// 次のように書けばよい
type A2 = Readonly<{ name: string; }>;

// こういうのもダメ
// error TS1355: A 'const' assertion can only be applied to a string, number, boolean, array, or object literal.
// let b1 = (Math.random() < 0.5 ? 0 : 1) as const;

// OK! b2 は 0 | 1 型になる
let b2 = Math.random() < 0.5 ?
    0 as const :
    1 as const;
```

## `globalThis` に型がついた

`globalThis` が[stage 3になった](https://github.com/tc39/proposal-global)ということで導入されたようです。

`globalThis` は今や静的に型付けされています…！

```ts
declare global {
    namespace foo {
        var bar: string;
    }

    // こういうことするとコンパイラがコケるのでやらないこと
    // https://github.com/Microsoft/TypeScript/issues/30459
    // namespace globalThis  {
    //     var test: string;
    // }
}

export {}
```

```ts
import "./extention";

console.log(globalThis.foo.bar);
console.log(this.foo.bar);

// globalThis.test = "a-b";
// console.log(globalThis.test.split("-"));
```

ざっくりこんな感じ。
`globalThis` の型がほしい時は素直に `typeof globalThis` を使えばいいそうな。

破壊的変更として、`this` にも型が付きました。
`typeof globalThis` って感じです。

コンテキストごとにglobalThisの型はどういつプロジェクト内でも異なる可能性があると思うんですがどうするのかな…。
なんかそういう仕組みって既にあったっけ…？(忘却の彼方)

## Genericsの型の推論が強化された

言語化して説明するのが重たい…
この破壊的変更が含まれます！

ざっくりサンプルコードを書いたので見てわかってほしい…！

```ts
// 2つの関数を引数を取り、T→U 変換して U → V 変換する場合 T→V な関数を返す
function compose<T, U, V>(f: (arg: T) => U, g: (arg: U) => V): (arg: T) => V {
    return (v1: T) => {
        const v2 = f(v1);
        const v3 = g(v2);
        return v3;
    }
}

function list<T>(x: T) { return [x]; }
function box<T>(value: T) { return { value }; }

// 今まではうまく推論できなかったので (arg: {}) => { value: {}[]; } になってた
// 3.4からちゃんとできるようになり f1 は <T>(arg: T) => { value: T[]; }
let f1 = compose(list, box);

// 今まではうまく推論できなかったので (arg: {}) => { value: {}; }[] になってた
// 3.4からちゃんとできるようになり f2 は <T>(arg: T) => { value: T; }[]
let f2 = compose(box, list);


let x1 = f1(100);
// T が {} ではなく正しく number にできるようになったのでエラーを検出できる！
// error TS2345: Argument of type '"hello"' is not assignable to parameter of type 'number'.
// x1.value.push("hello");

// 渡す関数にはGenericsの型パラメータが必要で、それがない場合は既存の挙動になる
// 推論できないパターン
const f3 = compose(x => [x], box);
const f4 = compose(function (x) { return [x]; }, box);
let x4 = f4(100);
// 検出に失敗する
x4.value.push("hello");

// 推論できるパターン
const f5 = compose(<T>(x: T) => [x], box);
const f6 = compose(function <T>(x: T) { return [x]; }, box);
let x6 = f6(100);
// ちゃんとエラーとして検出できる
// x6.value.push("hello");

// 複雑なパターンもいけるらしい
function compose2<A, B, C, D>(ab: (a: A) => B, cd: (c: C) => D): (a: [A, C]) => [B, D] {
    return ([a, c]) => {
        const b = ab(a);
        const d = cd(c);
        return [b, d];
    }
}
const f7 = compose2(list, box);
const f8 = compose2(box, list);
const f9 = compose2(list, list);


// rest parameterが絡むパターン
function compose3<A extends any[], B, C>(f: (...args: A) => B, g: (x: B) => C): (...args: A) => C {
    return (...args: A) => {
        const v1 = f(...args);
        const v2 = g(v1);
        return v2;
    }
}

// () => boolean
let f10 = compose3(() => true, b => !b);

// (x: any) => string
let f11 = compose3(x => "hello", s => s.length);

// <T, U>(x: T, y: U) => boolean … なんだけど
// T と U は比較しても常にfalseでは？と怒られる。偉い。
// error TS2367: This condition will always return 'false' since the types 'T' and 'U' have no overlap.
// let f12 = compose3(<T, U>(x: T, y: U) => ({ x, y }), o => o.x === o.y);

// (x: number) => string
let f13 = compose3((x: number) => x * x, x => `${x}`);


// 返り値の型にGenericsが含まれ、かつ文脈的に型が定まる場合、ちゃんと推論できるようになった
type Box<T> = { value: T };

function box2<T>(value: T): Box<T> {
    return { value }
}

// boxed1 の型から box2 への引数が正しいかどうかわかる
let boxed1: Box<'win' | 'draw'> = box2('draw');
// boxed2 の型から box2 への引数が正しくないことがわかる
// error TS2322: Type 'Box<"draw">' is not assignable to type 'Box<"win" | "lose">'.
// let boxed2: Box<'win' | 'lose'> = box2('draw');

// 返り値の型が明示的に宣言されていないとうまく動かない
// この定義だとvalueに何か変更の上returnされているかどうかがコードからはわからないため
function box3<T>(value: T) { return { value }; }
// error TS2322: Type '{ value: string; }' is not assignable to type 'Box<"win" | "draw">'.
// let boxed3: Box<'win' | 'draw'> = box3('draw');
```

## 関数へ複数の引数があるパターンをOptional型に変換するリファクタリングの追加

そのまんま。

```ts
function foo(a: number, b?: boolean, c = "foo") {
    return { a, b, c };
}
```

みたいなコードに対して `Convert to named parameter` すると

```ts
function foo({ a, b, c = "foo" }: { a: number; b?: boolean; c?: string; }) {
    return { a, b, c };
}
```

となります。
便利ですね。

さらに手で

```ts
type Option = { a: number; b?: boolean; c?: string; };
function foo({ a, b, c = "foo" }: Option) {
    return { a, b, c };
}
```

とかする感じでしょうか。

## TypeScript本体のビルドが jake から gulp に変更になった

gulpあまり好きじゃないのでちょっと悲しい…（それだけ）。

## おまけ

この前したこのツイートが妙にRTされたのでここにも貼っておきます。

<blockquote class="twitter-tweet" data-lang="ja"><p lang="ja" dir="ltr">僕も結構前からTypeScriptでinterfaceはあまり使わなくなりましたね(extendsしたいとかの理由付けができる時だけ使う) だいたいtypeで済ませている。理由としてinterfaceのほうがエラーメッセージやらがわかりやすいみたいなアドバンテージがほぼ消えた(typeが強くなった)ため。</p>&mdash; わかめ@毎日猫がいる (@vvakame) <a href="https://twitter.com/vvakame/status/1105738354862100481?ref_src=twsrc%5Etfw">2019年3月13日</a></blockquote>
<script async src="https://platform.twitter.com/widgets.js" charset="utf-8"></script>
