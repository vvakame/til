type A = { a: string };
type B = { b: string };

// 今は "a" | "b" と評価される
// 前も "a" | "b" と評価される
type T1 = keyof (A & B);

// 今は keyof T | "b" と評価される
// 前は keyof (T & { b: string; }) と評価されていた
type T2<T> = keyof (T & B);

// 今は "a" | keyof U と評価される
// 前は keyof ({ a: string; } & U) と評価されていた
type T3<U> = keyof (A & U);

// 今は keyof T | keyof U と評価される
// 前は keyof (T & U) と評価されていた
type T4<T, U> = keyof (T & U);

// ここから下のやつは全部
// 今は "b" | "a" と評価される
// 前も "a" | "b" と評価される
type T5 = T2<A>;
type T6 = T3<B>;
type T7 = T4<A, B>;

export {}
