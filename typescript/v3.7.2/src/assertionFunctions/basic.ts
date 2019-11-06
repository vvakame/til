// asserts の後にどの仮引数が検査対象なのか書く
// この関数がエラーにならずに処理を返したら、someVariable は呼び出し元の型検査フローに対して正しい
function assert(someVariable: any, msg?: string): asserts someVariable {
    if (!someVariable) {
        // 例外を投げて処理の流れをぶった切る
        throw new Error(msg)
    }
}

function multiplyA(x: any, y: any) {
    // x, y が本当に number だったら assert は例外を投げない (という実装と型定義だった)
    assert(typeof x === "number");
    assert(typeof y === "number");

    // ここでは x と y はnumber型に絞られている
    return x * y;
}

function multiplyB(x: any, y: any) {
    // 今まではこうやって書いたりしていた
    // throw とかすると今までもControl Flow解析で x と y の型が定まっていた
    if (typeof x !== "number") {
        throw new Error();
    }
    if (typeof y !== "number") {
        throw new Error();
    }

    // ここでは x と y はnumber型に絞られている
    return x * y;
}


// この関数が true を返したら仮引数 val の型は string ですよというアレ(前からあるやつ)
// https://www.typescriptlang.org/docs/handbook/advanced-types.html#using-type-predicates で解説されている
function isString(val: any): val is string {
    return typeof val === "string";
}

// asserts の後に type predicates と同じ書き方をする
function assertIsString(val: any): asserts val is string {
    if (typeof val !== "string") {
        throw new Error("Not a string!");
    }
}

function usageC(str: string | null) {
    assertIsString(str);
    // assertIsString が 例外を投げなかったら str は string に絞られている
    str.toUpperCase();
}


function assertIsDefined<T>(val: T): asserts val is NonNullable<T> {
    if (val === undefined || val === null) {
        throw new Error(
            `Expected 'val' to be defined, but received ${val}`
        );
    }
}

function usageD(str: string | null) {
    assertIsDefined(str);
    // assertIsString が 例外を投げなかったら str から null の可能性が除外される
    str.toUpperCase();
}

export {}
