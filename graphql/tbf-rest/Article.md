# GraphQLやってみた話

GAE/SE Node.jsが出るという[話があった](https://www.youtube.com/watch?v=ogexnfng_hE)。
筆者はGAE/SEに魂がロックインされているので、GAE/FEでしか動かなかったNode.jsはサーバとしては歯牙にもかけなかったわけです。
しかし、SEで出るなら話は変わります。
これからは、Node.jsで何かをさせることも現実的に視野に入ってくるようになります。

さて、Node.jsでサーバを立てれるようになると、やってみたいのはGraphQLです。

## 注意

ここで掲載しているコードは参考のために公開するものであり、実際に動かしてAPIサーバにトラフィックを発生させることはしないでください。

## 今回の目的

* GraphQLの、特にサーバ側の構成について色々試す
* 既存のREST APIからGraphQLに移行するのは現実的か検討する
* データ反映速度について諦めない
    * 何らかの理由により大本のデータが更新されたらすぐ参照できること
    * キャッシュを高速（1秒以内くらい）に破棄できること

## 過程

とりあえず https://www.apollographql.com/engine を使いはじめてみる。
技術書典のREST APIを適当にラップしていきます。

### ターゲットのREST APIについて

REST APIの仕様について先にざっくり解説します。

データ構造的には イベント（Event）＞サークル情報（CircleExhibitInfo）＞頒布物情報（ProductInfo）という感じです。
サークル情報のリストを得るには、イベントのIDが必要ですし、頒布物情報を得るには、サークル情報のIDが必要です。

このサービスの特徴として、公開されているデータの更新頻度はSNSなどと比べると段違いに低く、キャッシュを効かせるメリットが大きいです。
また、GAEを使っているため下手にHTTPレベルで `Cache-Control` などを使ってしまうと、GoogleのEdge Cacheに載ってしまいます。
Edge CacheはAPIなどを使って任意のタイミングでコンテンツを消すことができません。
つまり、データ反映のタイミングを制御できなくなってしまうため、これを避ける必要があります。
fastly+VCLみたいなのGoogleがホストしてくんないかな〜〜〜。

### 投げたいクエリを考える

例えば、技術書典4のサークル一覧と、各サークルの頒布物情報を得たい場合。
これはWebページでいうとサークル一覧ページとサークル詳細ページで表示する内容に相当しますね。
こんなクエリを書けると幸せそうです。

```
{
  event(id: "tbf04") {
    id
    name
    place
    circles {
      id
      name
      products {
        name
        firstAtTechBookFest
        firstAppearanceEventName
        relatedURLs
        images {
          url
        }
      }
    }
  }
}
```

### 実際にやってみた

このレスポンスを生成するために、Schemaを書いてそれに対してresolverを書きます（[`src/schema.ts`](https://github.com/vvakame/til/blob/76399fef649a2e7cd7d528d0350c92c692974f34/graphql/tbf-rest/src/schema.ts)）。

先のクエリを実行すると、次のようなREST APIコールに分解され実行されます。

* tbf04なEventの取得 × 1
* tbf04に紐づく246サークル情報の取得 リスト取得1回あたり10件取得する設定にしているので × 25
* 各サークルに紐づく頒布物情報の取得 1サークル1リクエストなので × 246

なんと、合計277回のREST API Callが発生する…！
これを更にもう1段深くしたりすると余裕で1000を超えます。

キャッシュが全くない状態だと、このクエリが帰ってくるのに11.8秒かかりました。
これは僕が現在作業している環境からtokyo regionのGAEで動いているREST APIまでの往復を考えるとかなりのオーバーヘッドが含まれます。
実際に稼働させた場合、tokyo regionのGAE(Node.js)からtokyo regionのGAE(Go)へのリクエストであり、データセンタ内の通信になると考えられるため実際はもっと短い時間で済むでしょう。
とはいえ、キャッシュはどう考えても必要です。

### fetchにキャッシュレイヤーを入れてみた

キャッシュを考えるについて、最終的に"キャッシュを自在に消せる"必要がありますが、とりあえず今はこれについて考えません。
キャッシュ自体にはRedisを使います。
Googleが[Memorystore](https://cloud.google.com/memorystore/)というFull-managed Redisを出したからです。
GAE/SE Node.js+Memorystoreという構成でえんちゃう？という目論見ですね。

まずは、REST API Callそのものをキャッシュすることを考えます。
`window.fetch` は、Responseをネットワークアクセス無しに `BodyInit` や `ResponseInit` からひねり出すことができます。
なので、 [`src/fetch.ts`](https://github.com/vvakame/til/blob/a51e50afbbdaef6ee8c98311796da62cfb6cd512/graphql/tbf-rest/src/fetch.ts) でfetchをwrapし、キャッシュの仕組みを仕込んでいます。
これにより、さきのクエリが300msで応答が返ってくるようになりました。

しかし、まだ不満があります。
キャッシュが常に存在しているとは限りませんし、キャッシュがない状態でもレスポンスに1秒以上かけたくはありません。

### DataLoaderでBatchGet化する

[DataLoader](https://github.com/facebook/dataloader)について検討します。
DataLoaderは1つ1つのデータ取得のリクエストをまとめて、バッチ化してくれます。
[mercari/datastoreのBatch](https://godoc.org/go.mercari.io/datastore#Client)みたいな感じでしょうか。

これを実際に適用してみる（[`src/dataLoader.ts`](https://github.com/vvakame/til/blob/76399fef649a2e7cd7d528d0350c92c692974f34/graphql/tbf-rest/src/dataLoader.ts)）と、次のようにバッチ化されました。

```
1 'event' 1
2 'circles' 1
3 'productInfos' 100
4 'productInfos' 100
5 'productInfos' 46
```

277 → 5回への大削減です！
limitを指定して分割することを考えると素直に5回では着地できないんですが、大きな改善です。

で、実際に図ってみると10.8秒とかでした。
さほど早くない…。
fetchに仕込んだRedisキャッシュをオフにして色々試してみます。

1. 18.6秒 cursorを使ったシーケンシャルなサークル情報取得 limit=10×25 + 頒布物情報BatchGet
2. 3.6秒 cursorを使ったシーケンシャルなサークル情報取得 limit=100×3 + 頒布物情報取得×246並列
3. 3.9秒 cursorを使ったシーケンシャルなサークル情報取得 limit=100×3 + 頒布物情報BatchGet×3並列

各結果は1回実行した結果で、ベンチマーク的に何回もトライしたわけではありません。
家からのネットワークアクセスなのでどうせ安定しないでしょ…という雑さです。
結構以外な結果ですね…。

リクエスト本数なんかより、ネットワークアクセスの並列度に左右される、という結果です。
BatchGet化できるところを考慮するよりcursorを使ったシーケンシャルなリスト取得をなんとかしたほうが早くなる！

BatchGetなAPI作るにしても、内部的にはDatastoreのAncestorQueryですらない、普通のQueryを何回も何回もバンバン投げることになる（277回くらい）のでなかなか辛いですね。
これを根本的に改善し、Datastore的にもBatchGetで処理できるようにするためには、上位階層のEntityに下位のKeyを全部持たせる必要がありそうです。
これでかなり早くなる確信はあるんだけどEntity間の同期取るのがクソめんどくさいのでできればやりたくない…。
下位のEntityだけポンと更新しておけばよかったのが、上位のEntityをTx組んで更新しなきゃいけないとかちょっとつらすぐるでしょう？

### なんとか1秒切りたいわん

最速でも4秒弱かかることが分かっているので、なんとか1秒以内にレスポンスを返す方法を考えます。
とりあえず1秒以内に画面に何かを表示できればいいので、データ全件を1発で取得することを諦めます。

具体的に `searchCircle` queryを使って100件ずつ結果を取得する方針を試してみます。
試しに、100件、100件、46件という順で取得する場合、1.9秒、1.5秒、0.9秒でレスポンスが得られました。

とりあえず1回目のクエリでは30件取得するようにすると、1秒を切れる…かもしれない！
サーバ内でのDatastore＞Memcacheのキャッシュの具合やら何やらも影響するので細かいことはわからん！
GraphQLサーバがtokyo regionで稼働すればもっと良い結果になるでしょう。
と考えると、体感速度もさほど悪いものにはならないと期待できそうです。

### RedisキャッシュをDataLoaderに仕込む

メモリ上でのキャッシュをdisableにし、全てRedis上で管理するようにします。
オンメモリキャッシュはEntityが更新された時に任意のタイミングで削除するのが難しいですからね… Cloud Pub/Subとか使えばできなくもない気がするけど。
fetchにRedisでのキャッシュを差し込むのはキャッシュキーを工夫することができず、任意のタイミングで消し飛ばしにくいのでやめます。

DataLoaderのRedisサポートは[ガバガバ](https://github.com/facebook/dataloader/blob/master/examples/Redis.md)です。
DataLoaderのキャッシュはJSのMapを使ったインメモリキャッシュを前提にしているようで、ここのレイヤーをRedisに差し替えるのは避けることにします。
DataLoaderが使っているキャッシュは単なるMapで、`get`, `set`, `delete`, `clear` の4操作しかありません。
なので、RedisのMSETやMGETを使ったキャッシュのバッチ操作の恩恵を受けられないのです。

DataLoaderをwrapするCustomLoaderを作りました（[`src/dataLoader.ts`](https://github.com/vvakame/til/blob/76399fef649a2e7cd7d528d0350c92c692974f34/graphql/tbf-rest/src/dataLoader.ts#L48)）。
いい感じに抽象化できたので、Kind毎にLoaderを作るのがかなり簡単になりました。だいたいコピペでOK！

キャッシュがある状態なら300msくらいで応答が返ってくるようになったので満足。

### 脳内でキャッシュ破棄の仕組みを構築する

Go側でEntityを更新した時、TaskQueue経由でNode.js側に通知し、Node.js上でRedisの該当キャッシュを消します。
Queryをキャッシュする場合、得られたEntityのIDとQueryの名前をSADDとかで紐付けておいて後で消します。

よし実装できた！脳内で！
1秒程度のレイテンシは全然許容範囲だしこれでいいでしょ。
いざとなったらRedisの中身全部ぶっ飛ばせばええんや。

### デプロイする

TODO はやくGAE/SE Node.jsリリースされてくれ〜〜既存のプロジェクトにデプロイしたいんじゃ〜〜〜

## 感想と考察

* 副次的なものだけど、データ変換レイヤーをJSで書けてクライアント側からそれが見えないのはいい…
  * JS friendlyなデータ構造に変換するのをJSで書けるのでJS friendlyにしやすくてすごい(意訳)
* レスポンス全体の組み立てにresolverがあって、Promiseを返せばいいのは楽だしわかりやすい
  * GraphQLのSchemaと同じ木になるようにするだけなので直感的
* DataLoaderは保存するキャッシュのスキーマ毎に1種類作るのがよさそう
  * Kind単位で作るのはよくなさそう
* 最終的に、クライアント側がどれだけ嬉しくなるか？を評価しないと導入の是非について議論できない
  * Angularで組まれてるので組み込むのがダルい予感がしている
* DBのSchema設計の時点でBatchGet friendlyな設計にしておく必要性を感じる
* よく考えるとGraphQLの特性ってSpannerのインタリーブの仕様と相性がいいのでは…？
  * SpannerはGraphQL friendlyなDBだった可能性が微粒子レベルで存在している？
* どの言語でGraphQLサーバを構築するべきか？
  * 今のところNode.js…？
  * or golang? (メル社内はほぼGo統一のふいんき（何故か変換できない）だからね)
* 誰がGraphQLサーバを構築するべきか？？
  * Backend Engineer? or Frontend Engineer?
  * GAE/SE Node.jsを前提にするならFrontend Engineerでもよさそ
    * Backend EngineerがCI/CDのパイプライン組むとこまではやる
    * JavaScriptの都合のいい仕様に改変しやすいので握れると便利そう
  * キャッシュ戦略やメンテナンスはBackend Engineerがやりたいのでは？
    * データのTTLは一元化管理するべき
    * データのTTLとかそのあたりはDB周りを開発・運用している人が握りたいよね
    * JavaScript簡単だから誰でも書けるでしょ（偏見）
      * nodebrew入れてNode.js入れてnpm installしてVSCodeで開くだけだよ（雑）
  * どっちがとか考えずに触りたいやつが触れよもう
* 処理途中のデータ という概念がない
  * いやREST APIにも無いっちゃ無いんだけど
  * 例えばサークルデータを246件取るとして、それなりに時間がかかるけど、まず最初の100件表示、追加100件、追加46件、みたいな途中のデータが得られない
  * なので、GraphQLを投げて返ってくるのを待つ間は単純に待ち時間だ！
  * 結局クライアント側で何らかの最適化は必要（searchCircle使うみたいな）
* GraphQLを使いたい奴はどこにいるのか？
  * Web frontendは使いたいだろうな
    * React+Flux とは相性良さそう
    * Angularは… どうだろ？
  * Android...
    * Androidの人(ひつじ)に話したらREST APIよりはえんちゃう？みたいな感じだった
  * iOS...
    * 適切な友達がいない(小声)
  * Backend EngineerはGraphQLの嬉しさ伝わりにくい説
    * sinmetalに話してみたけどウケなかった

## まだ分かってない

* エラーハンドリング
  * 404の時とか
  * その他Promiseがエラー返した時とか
* なんかlimit=30とかにするとRedisにsetすると応答が返ってこなくなる…？
  * 気の所為っぽい
* Schema手書きするのやめたい
  * swagger.json を自動生成しているのでそこからなんとかならないか？
    * sw2dts でTypeScriptの型定義ファイルは自動生成している
* DataLoaderも自動生成できるんじゃないのこれ…
* クライアント側の設計や画面とクエリの粒度の対応関係について
* Apolloのキャッシュレイヤー邪魔では？
  * というかなぜApolloを使うのか…？Facebookのやつのほうでよくない？
* ユーザ個別のデータのハンドリングやスタッフ限定情報のハンドリング
  * 例えばサークル情報には代表者氏名などを含むが、これをどうするか？
  * GraphQLではPublicなデータだけやり取りする…というのはよくなさそう
    * サークルチェックの情報は個人に従属するがGraphQLで扱えたほうが便利である

