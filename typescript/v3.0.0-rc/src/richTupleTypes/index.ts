// optionalな要素を簡単に書けるようになった 前は (number | undefined) でしたね
type Coordinate = [number, number, number?];

function coordinates(...args: Coordinate) {
    const [x, y, z] = args;
    return { x, y, z };
}
// v はちゃんと { x: number; y: number; z: number | undefined; } と推論される えらい！
let v = coordinates(1, 2);

// arrayのspreadingみたいな記法がサポートされた
// string[] に評価される
type SpreadedStrings = [...string[]];

// 空のtupleも作れるようになった
type Empty = [];
