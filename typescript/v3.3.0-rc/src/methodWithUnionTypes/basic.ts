type Fruit = "apple" | "orange";
type Color = "red" | "orange";

type FruitEater = (fruit: Fruit) => number;     // eats and ranks the fruit
type ColorConsumer = (color: Color) => string;  // consumes and describes the colors

declare let f: FruitEater | ColorConsumer;

// 今まではコレはダメだった。3.3以降はOK！
// error TS2349: Cannot invoke an expression whose type lacks a call signature.
//   Type 'FruitEater | ColorConsumer' has no compatible call signatures.
f("orange");
// 次の2つはNG 共通の引数ではないので
// error TS2345: Argument of type '"apple"' is not assignable to parameter of type '"orange"'.
// f("apple");
// error TS2345: Argument of type '"red"' is not assignable to parameter of type '"orange"'.
// f("red");

// TypeScript 3.3以降、次のように引数の型が評価される（はず）
type f_0 = (arg: Parameters<FruitEater>[0] & Parameters<ColorConsumer>[0]) => ReturnType<FruitEater> | ReturnType<ColorConsumer>;
// とりあえず展開できるところを普通に展開
type f_1 = (arg: Fruit & Color) => number | string;
// さらに展開
type f_2 = (arg: ("apple" & "red") | ("apple" & "orange") | ("orange" & "red") | ("orange" & "orange")) => number | string;
// 成立しえない部分は never になる
type f_3 = (arg: (never) | (never) | (never) | ("orange")) => number | string;
// そして "orange" だけが残った
type f_4 = (arg: "orange") => number | string;

export {}
