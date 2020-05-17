let str = "a";
let num : number | undefined = 1 as any;

let strP = Promise.resolve(str);
let numP = Promise.resolve(num);

// rStr1 の型
// TypeScript 3.8 では string | undefined
// TypeScript 3.9 から string
// undefinedのコンタミがなくなった！
let [rStr1, rNum1] = await Promise.all([strP, numP]);


let bool: boolean | null = true as any;
let boolP = Promise.resolve(bool);
// rStr2 の型
// TypeScript 3.8 では string | null | undefined
// TypeScript 3.9 から string
// rNum2 の型
// TypeScript 3.8 では number | null | undefined
// TypeScript 3.9 から number | undefined
// undefined, null のコンタミがなくなって直感的な型に
let [rStr2, rNum2, rBool2] = await Promise.all([strP, numP, boolP]);

export {};
