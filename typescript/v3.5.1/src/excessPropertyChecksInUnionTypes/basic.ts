type Style = {
    alignment: string;
    color?: string;
};

const s: Style = {
    alignment: "center",
    // よくあるミスの検出例 color を colour にtypo
    // error TS2322: Type '{ alignment: string; colour: string; }' is not assignable to type 'Style'.
    //   Object literal may only specify known properties, but 'colour' does not exist in type 'Style'. Did you mean to write 'color'?
    // colour: "grey",
};


type Point = {
    x: number;
    y: number;
};

type Label = {
    name: string;
};

// x, y, name は Point | Label として全部存在しうるのでOK
const obj1: Point | Label = {
    x: 0,
    y: 0,
    name: "foobar",
};

// このパターンは純粋に余計なものがあるのでエラーになる
// error TS2322: Type '{ x: number; y: number; name: boolean; }' is not assignable to type 'Point'.
//   Object literal may only specify known properties, and 'name' does not exist in type 'Point'.
// const obj2: Point = {
//     x: 0,
//     y: 0,
//     name: true,
// };

// v3.5 からこれがエラーになるようになった！ Point としては valid だけど余計かつ型が一致しない
// error TS2326: Types of property 'name' are incompatible.
//   Type 'boolean' is not assignable to type 'string | undefined'.
// const obj3: Point | Label = {
//     x: 0,
//     y: 0,
//     name: true,
// };

type PointOrLabel = (Point & Partial<Label>) | (Label & Partial<Point>);
const obj4: PointOrLabel = {
    x: 0,
    y: 0,
    // TypeScript v3.4 でもエラーになる
    //   error TS2322: Type 'true' is not assignable to type 'string | undefined'.
    // name: true,
};


export { }
