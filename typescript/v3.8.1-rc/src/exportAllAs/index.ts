// 今までのやり方
import * as a from "./libs";
export { a };

// まぁ普通にアクセスできる
// a.hello();

// 上記のimport, exportを一度にやる！
export * as b from "./libs";

// b はこのモジュール中には公開されないっぽい
// b.hello();
