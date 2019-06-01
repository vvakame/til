function foo1<T>(x: T): [T, string] {
    // TypeScript v3.5 ではエラー ベースが unknown なので
    // error TS2339: Property 'toString' does not exist on type 'T'.
    /// return [x, x.toString()];

    return [x, `${x}`];
}

foo1(null);

// extends {} で toString() が存在する制約を加える
// v3.4 以前と等価
function foo2<T extends {}>(x: T): [T, string] {
    return [x, x.toString()];
}

// null を渡そうとするとエラーとして検出できるようになった！
// error TS2345: Argument of type 'null' is not assignable to parameter of type '{}'.
// foo2(null);

export { }
