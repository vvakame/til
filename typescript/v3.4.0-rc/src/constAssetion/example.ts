// 実用例 as const 無しだと array1 は { kind: string; language?: string[]; endpoints?: string[]; } 的な型になってしまう
let array1 = [
    { kind: "AppEngine", services: ["default", "worker"] },
    { kind: "Cloud Functions", endpoints: ["Hello", "Bye"] },
] as const;
for (let value of array1) {
    // 各要素の持つ値がはっきりしているのでtype narrowingで安全にアクセスできる
    if (value.kind === "AppEngine") {
        value.services.forEach(v => console.log(v));
    } else {
        value.endpoints.forEach(v => console.log(v));
    }
}

// 既存の何かの型にあわせるみたいなのもできる
type CloudService = { kind: "AppEngine"; services: readonly string[]; } | { kind: "Cloud Functions", endpoints: readonly string[] };
let services: ReadonlyArray<CloudService> = array1;

// 素直にこう書けばよくない？という説もなきにしもあらず
let array2: CloudService[] = [
    { kind: "AppEngine", services: ["default", "worker"] },
    { kind: "Cloud Functions", endpoints: ["Hello", "Bye"] },
];

export { }
