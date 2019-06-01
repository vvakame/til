{
    type S = { done: boolean, value: number };
    type T =
        | { done: true, value: number }     // T0
        | { done: false, value: number };   // T1

    let s: S = void 0 as any;
    let t: T;

    // S は T0 部分に done: true の状態であれば割当可能
    // S は T1 部分に done: false の状態であれば割当可能
    // 各パターン合体させると { done: true | false; value: number; } となる
    // よってコンパイルは可能である！
    t = s;

    // なお ↑ TypeScript v3.4 では
    // error TS2322: Type 'S' is not assignable to type 'T'.
    //   Type 'S' is not assignable to type '{ done: false; value: number; }'.
    //     Types of property 'done' are incompatible.
    //       Type 'boolean' is not assignable to type 'false'.
}
{
    // done に null のパターンを追加
    type S = { done: boolean | null, value: number };
    type T =
        | { done: true, value: number }     // T0
        | { done: false, value: number };   // T1
    // | { done: null, value: number }  ← この可能性が残る

    let s: S = { done: true, value: 1 };
    let t: T;

    // この S の定義だと done: null のパターンがケアされないので T と互換ではない
    // error TS2322: Type 'S' is not assignable to type 'T'.
    //   Type 'S' is not assignable to type '{ done: false; value: number; }'.
    //     Types of property 'done' are incompatible.
    //       Type 'boolean | null' is not assignable to type 'false'.
    //         Type 'null' is not assignable to type 'false'.
    // t = s;
}

{
    type S = { a: 0 | 2, b: 4 };
    type T = { a: 0, b: 1 | 4 }     // T0
        | { a: 1, b: 2 }            // T1
        | { a: 2, b: 3 | 4 };       // T2

    let s: S = void 0 as any;
    let t: T

    // S は T0 部分に { a: 0, b: 4 } の状態であれば割当可能
    // S は T2 部分に { a: 2, b: 4 } の状態であれば割当可能
    // よってコンパイルは可能である！
    t = s;

    // なお ↑ TypeScript v3.4 では
    // error TS2322: Type 'S' is not assignable to type 'T'.
    //   Type 'S' is not assignable to type '{ a: 2; b: 4 | 3; }'.
    //     Types of property 'a' are incompatible.
    //       Type '0 | 2' is not assignable to type '2'.
    //         Type '0' is not assignable to type '2'.
}
{
    type S = { a: 0 | 2, b: 2 | 4 };
    type T = { a: 0, b: 1 | 4 }     // T0
        | { a: 1, b: 2 }            // T1
        | { a: 2, b: 3 | 4 }        // T2
        | { a: 0, b: 2 };           // T3

    let s: S = void 0 as any;
    let t: T

    // この場合、S が { a: 2, b: 2 } のパターンがカバーされてないのでエラーになる
    // t = s;
}

{
    // 判別可能な型の組み合わせは25パターンまで…
    //   Nで3パターン、Sで3プロパティがあり、3×3×3の27パターン…
    //   今回の仕組みは作動しない
    type N = 0 | 1 | 2;
    type S = { a: N, b: N, c: N };
    type T = { a: 0, b: N, c: N }
        | { a: 1, b: N, c: N }
        | { a: 2, b: N, c: N }
        | { a: N, b: 0, c: N }
        | { a: N, b: 1, c: N }
        | { a: N, b: 2, c: N }
        | { a: N, b: N, c: 0 }
        | { a: N, b: N, c: 1 }
        | { a: N, b: N, c: 2 };

    let s: S = void 0 as any;
    let t: T

    // この場合 S は T に割付可能だが組み合わせが複雑なので怒られる
    // error TS2322: Type 'S' is not assignable to type 'T'.
    //   Type 'S' is not assignable to type '{ a: 0 | 2 | 1; b: 0 | 2 | 1; c: 2; }'.
    //     Types of property 'c' are incompatible.
    //       Type '0 | 2 | 1' is not assignable to type '2'.
    //         Type '0' is not assignable to type '2'.
    // t = s;
}

export {}
