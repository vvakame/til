// この関数が値を返すことはない… (常に例外を投げるので)
function throwError(): never {
    throw new Error();
}

// TypeScript v3.6 ではコンパイルエラーになる
// error TS2366: Function lacks ending return statement and return type does not include 'undefined'.
// TypeScript v3.7 以降なら大丈夫
function multipler(v: any): string {
    if (typeof v === "string") {
        // 連結して2倍！
        return v + v;
    } else if (typeof v === "number") {
        // 2倍して2倍！(それはそう)
        return `${2 * v}`;
    }

    // v3.6 まではこう書くと あっ never ですね！返り値 string と矛盾しませんね！ ってなってた
    // return throwError();

    // v3.7 以降だとこれだけで あっ never ですね！ って伝わる
    throwError();
}
