# TypeScript v3.3.0-rc 変更点

こんにちは[メルペイ社](https://www.merpay.com/jp/)な[@vvakame](https://twitter.com/vvakame)です。

[TypeScript 3.3 RC](https://blogs.msdn.microsoft.com/typescript/2019/01/23/announcing-typescript-3-3-rc/)がアナウンスされました。

[What's new in TypeScriptも更新](https://github.com/Microsoft/TypeScript/wiki/What's-new-in-TypeScript#typescript-33)されています。
v3.3.0では破壊的変更<!-- https://github.com/Microsoft/TypeScript/wiki/Breaking-Changes#typescript-33 -->は存在しない予定です。エライ！

[この辺](https://github.com/vvakame/til/tree/master/typescript/v3.3.0-rc)に僕が試した時のコードを投げてあります。

## 変更点まとめ

* 関数などの呼び出し時にunion typesが絡む場合の挙動を改善 [Relaxed rules on methods of union types](https://github.com/Microsoft/TypeScript/pull/29011)
    * 今までコンパイル通らなかったけど通るパターンが出た
* `--build` の `--watch` でインクリメンタルビルドがサポートされた [File-incremental builds in `--build --watch` mode for composite projects](https://github.com/Microsoft/TypeScript/pull/29161)
    * `--build` は v3.0.0 で出たやつ [この辺](https://qiita.com/vvakame/items/57a0559c45b88b2ae168#%E3%83%97%E3%83%AD%E3%82%B8%E3%82%A7%E3%82%AF%E3%83%88%E9%96%93%E3%81%AE%E5%8F%82%E7%85%A7%E3%81%AE%E3%82%B5%E3%83%9D%E3%83%BC%E3%83%88)
    * 今までは変更検知するとフルビルドしなおしてたけどインクリメンタルビルドするようになった
    * 50%から75%くらいビルド時間削減できた

## 破壊的変更！

なし！ えらい！

## 関数などの呼び出し時にunion typesが絡む場合の挙動を改善

call signatureを持つ型同士でunion typesを構成した時、引数の型についてintersection typesを使った評価がされるようになりました。
何を言っているかわからないと思いますが要するに下のコードみたいなことが行われるようになりました。

```ts
type Fruit = "apple" | "orange";
type Color = "red" | "orange";

type FruitEater = (fruit: Fruit) => number;     // eats and ranks the fruit
type ColorConsumer = (color: Color) => string;  // consumes and describes the colors

declare let f: FruitEater | ColorConsumer;

// 今まではコレはダメだった。3.3以降はOK！
// error TS2349: Cannot invoke an expression whose type lacks a call signature.
//   Type 'FruitEater | ColorConsumer' has no compatible call signatures.
f("orange");
// 次の2つはNG 共通の引数ではないので
// error TS2345: Argument of type '"apple"' is not assignable to parameter of type '"orange"'.
// f("apple");
// error TS2345: Argument of type '"red"' is not assignable to parameter of type '"orange"'.
// f("red");

// TypeScript 3.3以降、次のように引数の型が評価される（はず）
type f_0 = (arg: Parameters<FruitEater>[0] & Parameters<ColorConsumer>[0]) => ReturnType<FruitEater> | ReturnType<ColorConsumer>;
// とりあえず展開できるところを普通に展開
type f_1 = (arg: Fruit & Color) => number | string;
// さらに展開
type f_2 = (arg: ("apple" & "red") | ("apple" & "orange") | ("orange" & "red") | ("orange" & "orange")) => number | string;
// 成立しえない部分は never になる
type f_3 = (arg: (never) | (never) | (never) | ("orange")) => number | string;
// そして "orange" だけが残った
type f_4 = (arg: "orange") => number | string;
```

実際これが役に立つシチュエーションはちょっと想像しにくいですね…。
なにかのパーサーとか…？

ちなみに、construct signatureにも適用できるようです。

```ts
type Fruit = "apple" | "orange";
type Color = "red" | "orange";

interface FruitConstructor {
    new(arg: Fruit): { fruit: string };
}
interface ColorConstructor {
    new(arg: Color): { color: Color };
}

declare let Ctor: FruitConstructor | ColorConstructor;
let obj = new Ctor("orange");
```

## `--build` の `--watch` でインクリメンタルビルドがサポートされた

そのまんまです。

`--build` は[以前紹介した](https://qiita.com/vvakame/items/57a0559c45b88b2ae168#%E3%83%97%E3%83%AD%E3%82%B8%E3%82%A7%E3%82%AF%E3%83%88%E9%96%93%E3%81%AE%E5%8F%82%E7%85%A7%E3%81%AE%E3%82%B5%E3%83%9D%E3%83%BC%E3%83%88)ように、プロジェクト間の参照があるような場合に一括でビルドするためのコマンドです。

昔から `-b` と `-w` の併用はできたんですが、なにか変更があるたびにフルビルドしていました。
ビルダAPIに新しくインクリメンタルビルド用のAPIを作り、内部的にそれを使うようにしたそうです。

以上。終わり！

今回はめっちゃ薄味でしたね…。
