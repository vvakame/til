declare type TypeName<T> = T extends string ? "string" : T extends number ? "number" : T extends boolean ? "boolean" : T extends undefined ? "undefined" : T extends Function ? "function" : "object";
declare const a1: TypeName<string>;
declare const a2: TypeName<"a">;
declare const b: TypeName<true>;
declare const c: TypeName<undefined>;
declare const d: TypeName<() => void>;
declare const e: TypeName<Date>;
