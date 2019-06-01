// TypeScript v3.5 からダメになった
//   error TS2322: Type '() => void' is not assignable to type '{ [k: string]: unknown; }'.
// let obj1: { [k: string]: unknown } = () => {};

// これは今まで通り
let obj2: { [k: string]: any } = () => {};


declare function someFunc(): void;
declare function fn<T>(arg: { [k: string]: T }): void;

// TypeScript v3.4 までは
//   error TS2345: Argument of type '() => void' is not assignable to parameter of type '{ [k: string]: {}; }'.
//     Index signature is missing in type '() => void'.
//   と評価されたのでコンパイルエラーとして検出できた。
//   { [k: string]: {}; } が { [k: string]: unknown; } に変わると、v3.4 のルールのままだと取りこぼしてしまう
// TypeScript v3.5 でのエラーは次のように変わった
// error TS2345: Argument of type '() => void' is not assignable to parameter of type '{ [k: string]: unknown; }'.
//   Index signature is missing in type '() => void'.
// fn(someFunc);

export {}
