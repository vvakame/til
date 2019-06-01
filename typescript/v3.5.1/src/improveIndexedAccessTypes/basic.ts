// read の時はうまく動く
type A = {
    s: string;
    n: number;
};

function read<K extends keyof A>(arg: A, key: K): A[K] {
    return arg[key];
}

const a: A = { s: "", n: 0 };
// x は string となる
const x = read(a, "s");


// write の時は危ない操作ができてしまいがち
function write<K extends keyof A>(arg: A, key: K, value: A[K]): void {
    // TypeScript v3.4 ではエラーにならない
    // TypeScript v3.5 では
    //   error TS2322: Type '"hello, world"' is not assignable to type 'A[K]'.
    //     Type '"hello, world"' is not assignable to type 'string & number'.
    //       Type '"hello, world"' is not assignable to type 'number'.
    // arg[key] = "hello, world";
}
// n: number だが string の値に置き換えてしまえる
write(a, "n", 1);

// 今まで検出できてなかったシリーズ！

function f1(obj: { a: number, b: string }, key: keyof typeof obj) {
    // v の型は string | number 取得できる値としてはただしい
    let v = obj[key];

    // obj の a は number だし b は string...
    // 次の2つの代入はどちらかが型的に不正な変更になる

    // error TS2322: Type '1' is not assignable to type 'number & string'.
    // obj[key] = 1;
    // error TS2322: Type '"x"' is not assignable to type 'number & string'.
    // obj[key] = 'x';

    // type narrowing で key の値を確定させればエラーにならない
    if (key === "a") {
        obj[key] = 1;
    }
    if (key === "b") {
        obj[key] = 'x';
    }

    // こうもできてほしいがまぁできない
    if (typeof key === "number") {
        // obj[key] = 1;
    }
}

function f2(obj: { a: number, b: 0 | 1 }, key: keyof typeof obj) {
    // v の型は number 取得できる値としてはまぁそう
    let v = obj[key];

    // a は number だし b は 0 | 1
    // b に 2 を入れるのは不正だができてしまう！

    obj[key] = 1;
    // error TS2322: Type '2' is not assignable to type '0 | 1'.
    // obj[key] = 2;

    if (key === "a") {
        obj[key] = 2;
    }
}

function f3<T extends { [key: string]: any }>(obj: T) {
    // any な値が取れる
    let v1 = obj['foo'];
    let v2 = obj['bar'];

    // T は { [key: string]: any } を底にしているが実際は
    // { [key: string]: boolean; } などになりうる(わかるまで2-3分悩んだ)

    // error TS2536: Type '"foo"' cannot be used to index type 'T'.
    // obj['foo'] = 123;
    // error TS2536: Type '"bar"' cannot be used to index type 'T'.
    // obj['bar'] = 'x';
}
// f3の実装のままだとまずいような使い方
f3<{ [key: string]: boolean; }>({ foo: true });

function f4<K extends string>(a: { [P in K]: number }, b: { [key: string]: number }) {
    // K は型であり、 a のプロパティと b のプロパティは当然一致しない
    // error TS2322: Type '{ [key: string]: number; }' is not assignable to type '{ [P in K]: number; }'.
    // a = b;
    // しかし、今のところ逆はいいらしい いいのか？(よくない気がする) (自分の足を撃つ権利の範疇？)
    b = a;
    b["b"] = 4;
}
let arg1 = {
    a: 1,
    b: true,
};
let arg2 = {
    b: 2,
    c: 3,
}
// f4の実装のままだとまずいような使い方
f4<"a">(arg1, arg2);

export { }
