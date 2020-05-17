import { data, setData } from "./a";

console.log("b-1", data);
setData(2);
console.log("b-2", data);

// v3.8, v3.9 ともに --module commonjs の場合 次のように出力される
// b-1 1
// b-2 2

export { data, setData };
