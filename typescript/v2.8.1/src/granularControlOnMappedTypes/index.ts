{
    // ビルトインの奴の定義を再掲
    type Readonly<T> = {
        readonly [P in keyof T]: T[P];
    };
    type Partial<T> = {
        [P in keyof T]?: T[P];
    };
}
{
    // 今回新しく追加されたビルトイン型
    type Required<T> = {
        [P in keyof T]-?: T[P];
    };
}
// ビルトインではないけどとりあえず
type Mutable<T> = {
    -readonly [P in keyof T]: T[P];
}

// 実験台
class Foo {
    str: string = "";
}

// ここの部分は前からあった
// { readonly str: string; } と評価される
type A = Readonly<Foo>;
// { readonly str?: string; } と評価される
type B = Partial<A>;

// ここから先は今回から
// { readonly str: string; } と評価される
type C = Required<B>;
// { str: string; } と評価される
type D = Mutable<C>;


export { }
