type T0 = { done: false, value: number };
type T1 = { done: true, value: number };
type T = T0 | T1;

let target: T;

// T0 を満たすのでコンパイルが通る
target = { done: true, value: 1 };
// T1 を満たすのでコンパイルが通る
target = { done: false, value: 1 };

// done が true とも false とも決まっていない bool の場合は…
type S = { done: boolean, value: number };
let source: S = { done: true, value: 1 };

// TypeScript v3.4 までは
//   error TS2322: Type 'S' is not assignable to type 'T'.
//     Type 'S' is not assignable to type 'T1'.
//       Types of property 'done' are incompatible.
//         Type 'boolean' is not assignable to type 'true'.
// TypeScript v3.5 ではコンパイルできる！
//   S と T0 を比べた結果 { done: false, value: number } の可能性は回収された！
//   残る T1 と比べた結果 { done: true, value: number } の可能性も回収された！
//   S が取りうる全ての可能性が回収されたので S と T は互換性があることがわかった
target = source;

export {}
