{
    // infer 使わないやつ
    type Flatten<T> = T extends any[] ? T[number] : T;

    // number
    type a = Flatten<number>;
    // number
    type b = Flatten<number[]>;
    // number[]
    type c = Flatten<number[][]>;
}
{
    // T[number] パット見て一瞬でわからない… 次も同じ意味
    type Flatten<T> = T extends any[] ? T[0] : T;
}
{
    // infer を使うとわりとわかりやすい！
    type Flatten<T> = T extends Array<infer U> ? U : T;
}
{
    // 残念ながら再帰的に自分を参照することはできない
    // type Flatten<T> = T extends Array<infer U> ? Flatten<U> : T;
    // type a = Flatten<number[][][][][]>;
}
