// TypeScript v3.4 でできるようになったやつ
function box<T>(value: T) {
    return {
        kind: "box" as const,
        value,
    };
}
function bag<T>(value: T) {
    return {
        kind: "bag" as const,
        value,
    };
}

function composeFunc<T, U, V>(
    f1: (x: T) => U,
    f2: (y: U) => V,
): (x: T) => V {
    return x => f2(f1(x));
}

let f1 = composeFunc(box, bag);
let a1 = f1(42);
// TypeScript v3.4 でこれがコンパイルエラーにできるようになった
// a1.value.value = "string";


// TypeScript v3.5 でできるようになったやつ
class Box<T> {
    kind: "box" = "box";
    constructor(public value: T) { }
}

class Bag<U> {
    kind: "bag" = "bag";
    constructor(public value: U) { }
}

function composeCtor<T, U, V>(
    C1: new (x: T) => U,
    C2: new (y: U) => V,
): (x: T) => V {
    return x => new C2(new C1(x));
}

// f2 は <T>(x: T) => Bag<Box<T>> になる
// v3.4 までは (x: {}) => Bag<{}> であった
let f2 = composeCtor(Box, Bag);
let a2 = f2(42);

// TypeScript v3.4.5 Property 'value' does not exist on type '{}'.
a2.value.value;

// TypeScript v3.5.1 Type '"string"' is not assignable to type 'number'.
// a2.value.value = "string";

export { a1, a2 };
