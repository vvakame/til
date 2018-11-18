function merge<T, U>(x: T, y: U) {
    // 今まではこの書き方はエラーになってた
    // error TS2698: Spread types may only be created from object types.
    // 3.2からイケる！
    return { ...x, ...y };
}

let obj = merge({ name: "vv" }, { cat: "yukari" });
