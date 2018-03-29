// infer R で新しい型パラメタを推論結果から導入する！つよい！
type ReturnType<T> = T extends (...args: any[]) => infer R ? R : T;

// 1. () => number は (...args: any[]) => any に代入可能である
// 2. infer R の位置が number である
// よって TmpType == number ！
type TmpType = ReturnType<() => number>;

// number[] と評価される
type ConcatReturn = ReturnType<Array<number>["concat"]>;
// 関数ではないのでDateと評価される（個人的にはneverになってほしい気がする）
type DateType = ReturnType<Date>;

// TypeScript 2.7以前で使ってたハック
const tmp = (false as true) && [1].concat([2]);
type ConcatReturnOld = typeof tmp;

// 残念ながらRの型の範囲は指定できない… っぽい？たぶん
// type ReturnType2<T, R extends string> = T extends (...args: any[]) => infer (R extends string) ? R : T;

// 型が書ける場所だったらだいたいイケるっぽい…
{
    type SecondArg<T> = T extends (a: string, b: infer R) => void ? R : never;
    // boolean！ ばばーん
    type TmpType = SecondArg<(a: string, b: boolean) => void>;
}
{
    type FooParamType<T> = T extends { Foo: infer R; } ? R : never;
    // RegExp！ ばばーん
    type TmpType1 = FooParamType<{ Foo: RegExp }>;
    // never！ ぼぼーん
    type TmpType2 = FooParamType<{ Bar: RegExp }>;
}
{
    type InstanceType<T> = T extends { new(...args: any[]): infer R; } ? R : never;
    // Date！ びびーん
    type DateInstance = InstanceType<typeof Date>;
}

export { }
