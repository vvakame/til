const sym = Symbol();

interface Foo {
    str: string;
    11: number;
    [sym]: symbol;
}

// typeof sym | "str" | 11 となる
type keyofFoo = keyof Foo;

// typeof sym
type A = Extract<keyofFoo, symbol>;
// str
type B = Extract<keyofFoo, string>;
// 11
type C = Extract<keyofFoo, number>;

// なお --keyofStringsOnly だと keyofFoo は "str" になる

// ちなみに…
type array = [1, 2, 3];
// "1", "2", "3", ...Arrayのメソッドいろいろ
type keyofArray = keyof array;
