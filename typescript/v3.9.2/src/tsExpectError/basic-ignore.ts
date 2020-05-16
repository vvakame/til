function doStuff(a: string, b: string) {
    // @ts-expect-error
    // コンパイルエラーを抑制することができる (今回追加された)
    hello(1);

    // @ts-ignore
    // コンパイルエラーを抑制することができる (前からある)
    hello(1);

    // // @ts-expect-error
    // error TS2578: Unused '@ts-expect-error' directive.
    // 何もしなくてもコンパイルが正しく通る場合、コンパイルエラーになる！
    hello("cat");

    // @ts-ignore
    // 特に効果がなくても特になんもならない
    hello("cat");
}

function hello(word: string) {
    console.log(`Hello, ${word}`);
}
