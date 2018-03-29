{
    type A = string | number | boolean | Date | null;

    // string | number | boolean | null と評価される
    type TmpA = Exclude<A, Date | RegExp>;
    // Date と評価される
    type TmpB = Extract<A, Date | RegExp>;
    // string | number | boolean | Date と評価される
    type TmpC = NonNullable<A>;
}
{
    // number と評価される
    type TmpA = ReturnType<() => number>;
    // string と評価される
    type TmpB = ReturnType<{ m(): string; }["m"]>;

    // Date と評価される
    type TmpC = InstanceType<typeof Date>;
    // RegExp と評価される
    type TmpD = InstanceType<{ new(): RegExp; }>;
}

// 気合で応用していくぞ…！
{
    class Clazz {
        str?: string;
        func?: (num: number) => number;
        method() { }
    }

    // Clazzの各要素の集合を取得して…
    type A = Clazz[keyof Clazz];
    // Functionであるもの（つまり関数とメソッド）だけを切り出す
    // ((num: number) => number) | (() => void) と評価される
    type B = Extract<A, Function>;
}
{
    // Pickの逆 指定したプロパティを含まない型を返す
    type Flip<T, K extends keyof T> = {
        [P in Exclude<keyof T, K>]: T[P];
    }
    // 指定したプロパティだけNonNullableを適用して返す
    // -? という記法は後述
    type PickWithNonNullable<T, K extends keyof T> = {
        [P in K]-?: NonNullable<T[P]>;
    };
    // 上記2つを組合せて指定したプロパティだけNonNullableにして返す
    type PropertyNonNullable<T, K extends keyof T> = PickWithNonNullable<T, K> & Flip<T, K>;

    class Clazz {
        str?: string;
        func?: (num: number) => number;
        method() { }
    }

    // { func: (num: number) => number; } & { str: string | undefined; method: () => void; }
    type A = PropertyNonNullable<Clazz, "func">;
    // { str: string; } & { func: (num: number) => number; method: () => void; }
    type B = PropertyNonNullable<A, "str">;

    const obj = new Clazz();

    obj.func = num => num;
    obj.func(1);
    const a: A = obj as A; // control flow的にはfuncはundefinedではないとわかっているがキャスト必要

    obj.str = "";
    obj.str.charAt(0);
    const b: B = a as B; // 同上
}
{
    // PRのdescriptionに書いてあったこっちのほうがまだしも有用っぽい
    type FunctionPropertyNames<T> = { [K in keyof T]: T[K] extends Function ? K : never }[keyof T];
    type FunctionProperties<T> = Pick<T, FunctionPropertyNames<T>>;
    type NonFunctionPropertyNames<T> = { [K in keyof T]: T[K] extends Function ? never : K }[keyof T];
    type NonFunctionProperties<T> = Pick<T, NonFunctionPropertyNames<T>>;
}
