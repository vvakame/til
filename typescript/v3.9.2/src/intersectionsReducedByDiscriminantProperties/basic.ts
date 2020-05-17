declare function smushObjects<T, U>(x: T, y: U): T & U;

interface Circle {
    kind: "circle";
    radius: number;
}

interface Square {
    kind: "square";
    sideLength: number;
}

declare let x: Circle;
declare let y: Square;

// v3.8 までは z の型は Circle & Square だった
// v3.9 からは z の型は never になった
// "circle" & "square" が両立することはありえないので never になる
let z = smushObjects(x, y);

// z が never になったので、zのプロパティへのアクセスはコンパイルエラーになるようになった
// 今までは z.kind の型は never だったので、このコード自体はコンパイルエラーにならなかった
// error TS2339: Property 'kind' does not exist on type 'never'.
//   The intersection 'Circle & Square' was reduced to 'never' because property 'kind' has conflicting types in some constituents.
// console.log(z.kind);

export {}
