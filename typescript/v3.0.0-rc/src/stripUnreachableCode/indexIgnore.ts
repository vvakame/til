function f() {
    return f();
    function f() {}
    type T = number;
    interface I {}
    const enum E {}
    namespace N { export type T = number; }
    var x;
}