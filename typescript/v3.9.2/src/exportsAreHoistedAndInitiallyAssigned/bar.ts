export * from "./foo";
// ./foo でも nameFromFoo を 定義して export している
// なので、ここで上書きできなければならない
export const nameFromFoo = 0;

// 出力されるJSコードの比較

// v3.8 までの出力 特に気にする点はない
// __export(require("./foo"));
// exports.nameFromFoo = 0;

// v3.9 からの出力 先に undefined で初期化されるようになった
// exports.nameFromFoo = void 0;
// __exportStar(require("./foo"), exports);
// exports.nameFromFoo = 0;

//  __exportStar の中で __createBinding を使うようになった
// その中で getter を定義してしまうので、後から上書きできない
// TypeError: Cannot set property nameFromFoo of #<Object> which has only a getter
// 上書きしたい場合、先にプロパティを作っておくと無視してくれるので別途後から値を設定する
