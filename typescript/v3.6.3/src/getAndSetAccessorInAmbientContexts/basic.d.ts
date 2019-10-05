declare class Foo {
    // 型定義の宣言で get, set は今まで使えなかった
    get x(): number;
    set x(val: number);
}

export { Foo };
