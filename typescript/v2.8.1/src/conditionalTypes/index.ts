// T の型に応じて文字列のリテラル型に変換
type TypeName<T> =
    T extends string ? "string" :
    T extends number ? "number" :
    T extends boolean ? "boolean" :
    T extends undefined ? "undefined" :
    T extends Function ? "function" :
    "object";

// string と互換性のある型をTに指定 → string
const a1: TypeName<string> = "string";
const a2: TypeName<"a"> = "string";

// 同様にそれぞれ互換性のある型に落ち着く
const b: TypeName<true> = "boolean";
const c: TypeName<undefined> = "undefined";
const d: TypeName<() => void> = "function";
const e: TypeName<Date> = "object";
