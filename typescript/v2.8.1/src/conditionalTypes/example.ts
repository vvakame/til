interface Id { id: number, /* other fields */ }
interface Name { name: string, /* other fields */ }

// 今までのやり方
declare function createLabelA(id: number): Id;
declare function createLabelA(name: string): Name;
declare function createLabelA(name: string | number): Id | Name;

// Conditional typesを使ったやり方
declare function createLabelB<T extends number | string>(idOrName: T):
    T extends number ? Id : Name;

// 今までのやり方だと…
let a1 = createLabelA("typescript");   // Name
let b1 = createLabelA(2.8);            // Id
let c1 = createLabelA("" as any);      // Id ← 最初にマッチしたものが採用されてしまう
let d1 = createLabelA("" as never);    // Id ← 最初にマッチしたものが採用されてしまう

// Conditional typesだと…
let a2 = createLabelB("typescript");   // Name
let b2 = createLabelB(2.8);            // Id
let c2 = createLabelB("" as any);      // Id | Name ← 偉い！
let d2 = createLabelB("" as never);    // never     ← 偉い！

export {}

{
    type Flatten<T> = T extends any[] ? T[0] : T;
    let z1: Flatten<string>;
    let z2: Flatten<string[]>;
}

{
    type Flatten<T> = T extends Array<infer U> ? U : T;
    let z1: Flatten<string>;
    let z2: Flatten<string[]>;
}

{
    type T10 = TypeName<string | (() => void)>;  // "string" | "function"
    type T12 = TypeName<string | string[] | undefined>;  // "string" | "object" | "undefined"
    type T11 = TypeName<string[] | number[]>;  // "object"
}
