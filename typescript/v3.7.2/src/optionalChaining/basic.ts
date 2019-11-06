let foo: any = {
    bar1: { buzz() { console.log("bar1"); } },
    bar2: void 0,
};

// bar1 と表示される
let x = foo?.bar1?.buzz();
// 何も表示されない
let y = foo?.bar2?.buzz();

// これはエラーになる
// error TS1109: Expression expected.
// ↓ 最後の ? が三項演算子だと思われてて面白い
// error TS1005: ':' expected.
// let z1 = foo?.bar?.buzz?();
// error TS1109: Expression expected.
// let z2 = foo?.bar?.buzz()?;

// ちなみにこれらはOK
// ?. でワンセット
let z3 = foo?.bar?.buzz?.();
let z4 = []?.[1];

export {}
