/**
 *  T（union types）から、Uで指定した型を除外したものを返す
 */
type Exclude<T, U> = T extends U ? never : T;

/**
 * T（union types）から、Uで指定した型に含まれるもののみを返す
 */
type Extract<T, U> = T extends U ? T : never;

/**
 * T（union types）から、null, undefinedを除外して返す
 * Exclude<T, null | undefined> と一緒（のはず）
 */
type NonNullable<T> = T extends null | undefined ? never : T;

/**
 * 関数の返り値の型を切り出す
 */
type ReturnType<T extends (...args: any[]) => any> = T extends (...args: any[]) => infer R ? R : any;

/**
 * コンストラクタを持つ型からインスタンスの型を切り出す
 */
type InstanceType<T extends new (...args: any[]) => any> = T extends new (...args: any[]) => infer R ? R : any;

{
    type A = Exclude<string | null, null | undefined>;
    type B = NonNullable<string | null>;
}

export { }
