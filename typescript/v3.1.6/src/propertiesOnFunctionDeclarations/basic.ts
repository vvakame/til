// 今までのやり方。namespace はTypeScriptの独自要素。
function fooA() {
    console.log("fooA");
}
namespace fooA {
  export var barA = () => {
    console.log("fooA.barA");
  };
}

fooA();
fooA.barA();

// JSだったら普通素直にこう書くよね…？という書き方ができるようになった。
function fooB() {
    console.log("fooB");
}
fooB.barB = () => {
    console.log("fooB.barB");
};

fooB();
fooB.barB();

// const+関数でもできる
const fooC = () => {
    console.log("fooC");
}
fooC.barC = () => {
    console.log("fooC.barC");
};
fooC();
fooC.barC();

// let(やvar)+関数ではNG
let fooD = () => {
    console.log("fooD");
}
// この書き方はエラーになる fooDの値は差し替え可能だからね。仕方ないね。
// error TS2339: Property 'barD' does not exist on type '() => void'.
// fooD.barD = () => {
//     console.log("fooD.barD");
// };
// fooD();
// fooD.barD();
