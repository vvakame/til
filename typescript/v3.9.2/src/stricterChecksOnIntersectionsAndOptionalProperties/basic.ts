interface A {
    a: number; // ここの a は number
}

interface B {
    b: string;
}

interface C {
    a?: boolean; // ここの a は boolean
    b: string;
}

declare let x: A & B;
declare let y: C;

// y = x;
// ↑ 上記コードは A と C は一致しないが B と C は一致するので許されていた
// が、当然実行時エラーになる可能性がある
// v3.9 から、 A & B は C に一致しないので怒られるようになった
// error TS2322: Type 'A & B' is not assignable to type 'C'.
//   Types of property 'a' are incompatible.
//   Type 'number' is not assignable to type 'boolean | undefined'.
// これ系の重箱の隅系、毎回色々潰されていると思うんだけどまだあるんだなぁ とびっくりしますね
