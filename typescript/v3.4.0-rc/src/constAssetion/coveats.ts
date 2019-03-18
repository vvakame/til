// 型アサーション自体は値の世界の住人なので型の世界で利用することはできない
// …というかその必要がないよね
// type A1 = { name: string; } as const;
// 次のように書けばよい
type A2 = Readonly<{ name: string; }>;

// こういうのもダメ
// error TS1355: A 'const' assertion can only be applied to a string, number, boolean, array, or object literal.
// let b1 = (Math.random() < 0.5 ? 0 : 1) as const;

// OK! b2 は 0 | 1 型になる
let b2 = Math.random() < 0.5 ?
    0 as const :
    1 as const;

export { }
