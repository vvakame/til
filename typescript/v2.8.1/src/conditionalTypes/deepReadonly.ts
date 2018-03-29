type DeepReadonly<T> =
    T extends any[] ? DeepReadonlyArray<T[number]> :
    T extends object ? DeepReadonlyObject<T> :
    T;

interface DeepReadonlyArray<T> extends ReadonlyArray<DeepReadonly<T>> { }

type NonFunctionPropertyNames<T> = { [K in keyof T]: T[K] extends Function ? never : K }[keyof T];

type DeepReadonlyObject<T> = {
    readonly [P in NonFunctionPropertyNames<T>]: DeepReadonly<T[P]>;
};

function toReadonly<T>(obj: T): DeepReadonly<T> {
    return obj as any;
}

const obj = toReadonly({
    a: {
        a1: true,
    },
    b: [1, 2,],
});
// read-only だからエラーになる！
// obj.a.a1 = false;
// obj.b[3] = 3;
