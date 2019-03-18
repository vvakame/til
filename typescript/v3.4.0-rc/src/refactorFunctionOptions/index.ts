type Option = { a: number; b?: boolean; c?: string; };
function foo({ a, b, c = "foo" }: Option) {
    return { a, b, c };
}

export { }
