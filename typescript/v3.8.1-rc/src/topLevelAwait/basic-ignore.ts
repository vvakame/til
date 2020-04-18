// target が es2017 以降で module が esnext か system じゃないとダメ
// es2017 以降だと await がキーワードになっているので

function timeout(milliseconds: number) {
    return new Promise(resolve => {
        setTimeout(resolve, milliseconds);
    })
}

console.log("Just 1 sec");
await timeout(1000);
console.log("dit it!");

// export とか import とかないとモジュールだと思ってもらえないやつ
//   --isolatedModules でもいいらしい(試してない)
// error TS1375: 'await' expressions are only allowed at the top level of a file when that file is a module, but this file has no imports or exports. Consider adding an empty 'export {}' to make this file a module.
export { }
