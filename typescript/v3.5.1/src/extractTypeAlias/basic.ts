// { kind: string; name: string; } 部分を範囲選択して Extract to type alias
function report(cat: { kind: string; name: string; }): string {
    return `${cat.name}, ${cat.kind}`;
}
