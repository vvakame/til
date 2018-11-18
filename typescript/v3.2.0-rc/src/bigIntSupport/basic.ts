let bi1: bigint = 1000n;

// 新しく導入された型
type f = bigint | BigInt | BigInt64Array | BigUint64Array;


let bi2 = BigInt(1);
let mod1 = BigInt.asIntN(64, 10000000000000000000n);
let mod2 = BigInt.asUintN(64, 10000000000000000000n);

// 1000n 1n -8446744073709551616n 10000000000000000000n と表示される
console.log(bi1, bi2, mod1, mod2);

// 四則演算
console.log(1n + 2n);
// numberと比較するやつ
// error TS2367: This condition will always return 'false' since the types '1n' and '1' have no overlap.
// console.log(1n == 1);
console.log(1n < 2);

// number とは非互換
let n: number = 0;
// bi1 = n;
// n = bi1;

function whatKindOfNumberIsIt(x: number | bigint) {
    if (typeof x === "bigint") {
        console.log("'x' is a bigint!");
    } else {
        console.log("'x' is a floating-point number");
    }
}

whatKindOfNumberIsIt(1);
whatKindOfNumberIsIt(1n);
