function tag<T>(strs: TemplateStringsArray, ...args: T[]) {
    console.log(strs, ...args);
}

// パラメタが全部 number なのでOK
tag<number>`hoge${1}${2}`;
// 2つ目のパラメタが number なのでエラー
// tag<string>`hoge${"A"}${2}`;
// 1つ目と2つ目のパラメタが一致してないのでエラー
// tag`${true}${1}`;


declare function styledComponent<Props>(strs: TemplateStringsArray): Component<Props>;

styledComponent<MyProps>`
    font-size: 1.5em;
`

interface Component<T> {
}

interface MyProps {
    name: string;
    age: number;
}
