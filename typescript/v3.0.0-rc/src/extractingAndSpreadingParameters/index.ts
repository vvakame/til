function call<TS extends any[], R>(fn: (...args: TS) => R, ...args: TS): R {
    return fn(...args);
}

function hello(word = "TypeScript") {
    return `Hello, ${word}`;
}

// TS は [(string | undefined)?] と推論されている
let str1 = call(hello, "JavaScript");
let str2 = call(hello, void 0);
let str3 = call(hello);

// これはちゃんとエラーになる！えらい！
// index.ts:18:17 - error TS2345: Argument of type '(word?: string) => string' is not assignable to parameter of type '(args_0: boolean) => string'.
//   Types of parameters 'word' and 'args_0' are incompatible.
//     Type 'boolean' is not assignable to type 'string | undefined'.
// let str4 = call(hello, true);

// 引数から TS を推論させて word の型指定を省略することもできる かしこい
call(word => `Hello, ${word.toUpperCase()}`, "TypeScript");
