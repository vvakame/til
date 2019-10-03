# TypeScript v3.6.0 Beta 変更点

こんにちは[メルペイ社](https://www.merpay.com/)な[@vvakame](https://twitter.com/vvakame)です。

[TypeScript 3.6 Beta](https://devblogs.microsoft.com/typescript/announcing-typescript-3-6-beta/)がアナウンスされました。

* What's new in TypeScript in 3.6 → まだ <!-- https://github.com/Microsoft/TypeScript/wiki/What's-new-in-TypeScript#typescript-36 -->
* Breaking Changes in 3.6 → まだ <!-- https://github.com/Microsoft/TypeScript/wiki/Breaking-Changes#typescript-35 -->
* [TypeScript 3.6 Iteration Plan](https://github.com/microsoft/TypeScript/issues/31639)
* TypeScript Roadmap: July - December 2019 → まだ <!-- https://github.com/Microsoft/TypeScript/issues/29288 -->

Roadmapは[こちら](https://github.com/Microsoft/TypeScript/wiki/Roadmap#36-august-2019)。

[この辺](https://github.com/vvakame/til/tree/master/typescript/v3.6.0-rc)に僕が試した時のコードを投げてあります。

## 変更点まとめ

* ジェネレータの型の改善 [Strongly typed iterators and generators](https://github.com/Microsoft/TypeScript/issues/2983)
    * next と done で異なる型の値を返してそれを区別できるようになった
* Array spreadingの挙動の修正 [More accurate array spreads](https://github.com/microsoft/TypeScript/pull/31166)
    * `[...Array(5)]` のes5へのdownpileがより正確に行われるようになったらしい
    * tslibに `__spreadArrays` が追加された感じ
    * Issue立ってから実に3年越しの修正
* Promiseの使い方下手こいた時のUXの改善 [Improved UX around Promises](https://github.com/microsoft/TypeScript/issues/30646)
    * 今まで気にしてなかったけど確かに欲しいやつ…！
* セミコロンがない時のStatement追加時の挙動を改善 [Semicolon-aware auto-imports](https://github.com/microsoft/TypeScript/issues/19882)
    * セミコロンレス派が喜びそうなやつです
    * セミコロンちゃんと書け(or フォーマッタに自動追加させろ)派としてはいるのかこれ と思っています
* `--declaration` と `--isolatedModules` の併用の改善 [`--declaration` and `--isolatedModules`](https://github.com/Microsoft/TypeScript/issues/29490)
    * v3.6.1 で入りそう まだmergeとかされてない
* async内でawaitがいるような候補を選んだら自動的にawaitを挿入する [Auto-inserted `await` for completions](https://github.com/microsoft/TypeScript/issues/31450)
    * v3.6.1 で入りそう まだmergeとかされてない
* Call Hierarchyのサポート [Call Hierarchy support](https://github.com/microsoft/TypeScript/issues/31863)
    * v3.6.1 で入りそう まだmergeとかされてない
    * Find Referencesはすでにあるけどガバガバ辿れたほうがいいよねー的な

## 破壊的変更！

* 

## ジェネレータの型の改善
