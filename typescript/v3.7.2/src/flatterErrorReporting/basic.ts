type SomeVeryBigType = { a: { b: { c: { d: { e: { f(): string } } } } } }
type AnotherVeryBigType = { a: { b: { c: { d: { e: { f(): number } } } } } }

declare let x: SomeVeryBigType;
declare let y: AnotherVeryBigType;

// 短くてわかりやすいエラー！
// error TS2322: Type 'SomeVeryBigType' is not assignable to type 'AnotherVeryBigType'.
//   The types returned by 'a.b.c.d.e.f()' are incompatible between these types.
//   Type 'string' is not assignable to type 'number'.
// y = x;

// ちなみに今までのエラーはこんな感じ
// TSユーザは経験的に型関連のエラーは一番下から見ていけばいいと知っているし、この場合は下4行を読むと原因がわかる
// けど、長いものは長い。
// error TS2322: Type 'SomeVeryBigType' is not assignable to type 'AnotherVeryBigType'.
//   Types of property 'a' are incompatible.
//     Type '{ b: { c: { d: { e: { f(): string; }; }; }; }; }' is not assignable to type '{ b: { c: { d: { e: { f(): number; }; }; }; }; }'.
//       Types of property 'b' are incompatible.
//         Type '{ c: { d: { e: { f(): string; }; }; }; }' is not assignable to type '{ c: { d: { e: { f(): number; }; }; }; }'.
//           Types of property 'c' are incompatible.
//             Type '{ d: { e: { f(): string; }; }; }' is not assignable to type '{ d: { e: { f(): number; }; }; }'.
//               Types of property 'd' are incompatible.
//                 Type '{ e: { f(): string; }; }' is not assignable to type '{ e: { f(): number; }; }'.
//                   Types of property 'e' are incompatible.
//                     Type '{ f(): string; }' is not assignable to type '{ f(): number; }'.
//                       Types of property 'f' are incompatible.
//                         Type '() => string' is not assignable to type '() => number'.
//                           Type 'string' is not assignable to type 'number'.

export {}
