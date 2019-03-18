// 高階関数の型推論がより賢くなった

// 2つの関数を引数を取り、T→U 変換して U → V 変換する場合 T→V な関数を返す
function compose<T, U, V>(f: (arg: T) => U, g: (arg: U) => V): (arg: T) => V {
    return (v1: T) => {
        const v2 = f(v1);
        const v3 = g(v2);
        return v3;
    }
}

function list<T>(x: T) { return [x]; }
function box<T>(value: T) { return { value }; }

// 今まではうまく推論できなかったので (arg: {}) => { value: {}[]; } になってた
// 3.4からちゃんとできるようになり f1 は <T>(arg: T) => { value: T[]; }
let f1 = compose(list, box);

// 今まではうまく推論できなかったので (arg: {}) => { value: {}; }[] になってた
// 3.4からちゃんとできるようになり f2 は <T>(arg: T) => { value: T; }[]
let f2 = compose(box, list);


let x1 = f1(100);
// T が {} ではなく正しく number にできるようになったのでエラーを検出できる！
// error TS2345: Argument of type '"hello"' is not assignable to parameter of type 'number'.
// x1.value.push("hello");

// 渡す関数にはGenericsの型パラメータが必要で、それがない場合は既存の挙動になる
// 推論できないパターン
const f3 = compose(x => [x], box);
const f4 = compose(function (x) { return [x]; }, box);
let x4 = f4(100);
// 検出に失敗する
x4.value.push("hello");

// 推論できるパターン
const f5 = compose(<T>(x: T) => [x], box);
const f6 = compose(function <T>(x: T) { return [x]; }, box);
let x6 = f6(100);
// ちゃんとエラーとして検出できる
// x6.value.push("hello");

// 複雑なパターンもいけるらしい
function compose2<A, B, C, D>(ab: (a: A) => B, cd: (c: C) => D): (a: [A, C]) => [B, D] {
    return ([a, c]) => {
        const b = ab(a);
        const d = cd(c);
        return [b, d];
    }
}
const f7 = compose2(list, box);
const f8 = compose2(box, list);
const f9 = compose2(list, list);

// rest parameterが絡むパターン
function compose3<A extends any[], B, C>(f: (...args: A) => B, g: (x: B) => C): (...args: A) => C {
    return (...args: A) => {
        const v1 = f(...args);
        const v2 = g(v1);
        return v2;
    }
}

// () => boolean
let f10 = compose3(() => true, b => !b);

// (x: any) => string
let f11 = compose3(x => "hello", s => s.length);

// <T, U>(x: T, y: U) => boolean … なんだけど
// T と U は比較しても常にfalseでは？と怒られる。偉い。
// error TS2367: This condition will always return 'false' since the types 'T' and 'U' have no overlap.
// let f12 = compose3(<T, U>(x: T, y: U) => ({ x, y }), o => o.x === o.y);

// (x: number) => string
let f13 = compose3((x: number) => x * x, x => `${x}`);


// 返り値の型にGenericsが含まれ、かつ文脈的に型が定まる場合、ちゃんと推論できるようになった
type Box<T> = { value: T };

function box2<T>(value: T): Box<T> {
    return { value }
}

// boxed1 の型から box2 への引数が正しいかどうかわかる
let boxed1: Box<'win' | 'draw'> = box2('draw');
// boxed2 の型から box2 への引数が正しくないことがわかる
// error TS2322: Type 'Box<"draw">' is not assignable to type 'Box<"win" | "lose">'.
// let boxed2: Box<'win' | 'lose'> = box2('draw');

// 返り値の型が明示的に宣言されていないとうまく動かない
// この定義だとvalueに何か変更の上returnされているかどうかがコードからはわからないため
function box3<T>(value: T) { return { value }; }
// error TS2322: Type '{ value: string; }' is not assignable to type 'Box<"win" | "draw">'.
// let boxed3: Box<'win' | 'draw'> = box3('draw');


export { }
