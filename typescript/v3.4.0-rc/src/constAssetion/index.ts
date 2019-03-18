// as const でリテラルを具体的なreadonlyなオブジェクト型リテラル相当の表現に変換できる
// a は 10 型
let a = 10 as const;
// const の場合は昔から 10 型
const a1 = 10;

// 型アサーションと同じ記法なので前置の書き方もできる（この書き方使わなくなりましたねぇ…）
// b は readonly [10, 20] 型
let b = <const>[10, 20];

// オブジェクトリテラルにも適用できる 配列にもOK
// c1 は { readonly text: "hello" ; } 型
let c1 = { text: "hello" } as const;  // Type { readonly text: "hello" }
// c2 は [true, false] 型
let c2 = [true, false] as const;


// こういう複雑なオブジェクトもconstにできる
let d1 = { lunch: "saizeriya" };
let d = {
    name: "vvakame",
    love: {
        kind: "cat",
        name: "yukari",
    },
    location: "tokyo",
    note: d1,
} as const;

// NG! これは怒られる
// error TS2540: Cannot assign to 'note' because it is a read-only property.
// d.note = { lunch: "CoCo壱" };

// ここはreadonlyではない
d.note.lunch = "rigoletto";

export { }
