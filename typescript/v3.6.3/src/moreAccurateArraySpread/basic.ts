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

export { }
