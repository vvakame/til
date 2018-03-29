{
    type ToArray<T> = T extends any ? T[] : never;

    // これは次の形に展開される
    //   (string extends any ? string[] : never) |
    //   (number extends any ? number[] : never)
    // つまり、簡約すると string[] | number[]
    // union typesは個別にconditional typesの適用を受けるのだ
    type Tmp = ToArray<string | number>;

    // 括弧でくくっても string[] | number[] にされる(げせぬ)
    type TmpB = ToArray<(string | number)>;
}
{
    // T を [T] で展開しない！
    // また紛らわしい記法が…
    type ToArray<T> = [T] extends any ? T[] : never;

    // (string | number)[] と評価される
    type Tmp = ToArray<string | number>;
}
{
    // なお、neverが混ざると最終的に消えます
    type Foo<T> = T extends any ? T : never;
    // never は消えて string | number になる
    type Tmp = Foo<string | number | never>;
}
{
    // any も [any] にしても結果は変わらない
    // なぜなのか…
    type ToArray<T> = [T] extends [any] ? T[] : never;

    // (string | number)[] と評価される
    type Tmp = ToArray<string | number>;
}
{
    // もうおじちゃんの脳は型注釈を正しくパースできなくなってきましたよ…
    const array: [number] = [1];
}
