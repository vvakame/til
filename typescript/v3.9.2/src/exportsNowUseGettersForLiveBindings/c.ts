import { data, setData } from "./b";

console.log("c-1", data);
setData(3);
console.log("c-2", data);

// v3.8 だと --module commonjs の場合 次のように出力される
// c-1 1
// c-2 1
// v3.9 だと --module commonjs の場合 次のように出力される
// c-1 2
// c-2 3

export { data, setData };
