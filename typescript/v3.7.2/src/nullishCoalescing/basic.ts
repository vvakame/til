let foo: string | null = null as any;
let bar = "bar";

let a = foo ?? bar;
// bar と表示される foo が null なので
console.log(a);

foo = "" as any;
let b = foo ?? bar;
// 空文字列が表示される
// || と違って、null と undefined の時のみ右辺が評価される
// "" は当てはまらないので左辺の値が返る
console.log(b);

let c = foo || bar;
// bar と表示される
// "" は falsy な値なので 右辺が評価される
console.log(c);

// ?? と同じことをしてみる
let d = foo == null ? foo : bar;
// bar と表示される
// == null に当てはまるのは undefined と null のみ
console.log(d);

export {}
