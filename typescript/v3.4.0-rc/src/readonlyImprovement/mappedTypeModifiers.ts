// mapped typesの挙動が変わった

// 今までもこれからも (string | undefined)[]
// Arrayに対してMapped Typesを適用すると各要素に適用された
type A0 = Partial<string[]>;

// Readonlyについては上記のルールは当てはまらなかった！が、今回から適用されるようになった
// これから readonly string[]
// これまで string[]
type A1 = Readonly<string[]>;

// -readonly による属性剥がしも同様
type Writable<T> = {
    -readonly [K in keyof T]: Writable<T[K]>;
}

// これから string[]
// これまで ReadonlyArray<any>
type A2 = Writable<ReadonlyArray<string>>;

export { }
