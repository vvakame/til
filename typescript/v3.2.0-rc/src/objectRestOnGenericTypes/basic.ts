interface XYZ { x: any; y: any; z: any; }

// T から Pickする 対象は TのプロパティからXYZのプロパティと重複する部分を除外したもの
type DropXYZ<T> = Pick<T, Exclude<keyof T, "x" | "y" | "z">>;

function dropXYZ<T extends XYZ>(obj: T): DropXYZ<T> {
    let { x, y, z, ...rest } = obj;
    return rest;
}

let obj = dropXYZ({ w: null, x: 1, y: "str", z: true });

// OK
obj.w;
// 存在しないのでコンパイルエラー！
// error TS2339: Property 'x' does not exist on type 'Pick<{ w: null; x: number; y: string; z: boolean; }, "w">'.
// obj.x;
// error TS2339: Property 'y' does not exist on type 'Pick<{ w: null; x: number; y: string; z: boolean; }, "w">'.
// obj.y;
// error TS2339: Property 'z' does not exist on type 'Pick<{ w: null; x: number; y: string; z: boolean; }, "w">'.
// obj.z;

export { }
