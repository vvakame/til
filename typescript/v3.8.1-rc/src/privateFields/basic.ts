// error TS18028: Private identifiers are only available when targeting ECMAScript 2015 and higher.

class Person {
    #name: string
    private fav?: string;
    // 名前空間が分かれてるので同名のプロパティが存在してもよい
    name: string;

    // 他の可視性制御の修飾子と併用したらダメ
    // error TS18010: An accessibility modifier cannot be used with a private identifier.
    // public #other1?: string;
    // private #other2?: string;
    // protected #other3?: string;
    // これもダメ abstract class で定義しても次のエラーが出る
    // error TS18019: 'abstract' modifier cannot be used with a private identifier
    // abstract #other4?: string;

    // これはOKらしい
    readonly #other5?: string;

    constructor(name: string, fav?: string) {
        this.#name = name;
        this.name = name;
        this.fav = fav;
    }

    greet() {
        console.log(`Hello, my name is ${this.#name}!${this.fav ? ` I like ${this.fav}.` : ""}`);
    }
}

class Person2 extends Person {
    // 継承しても親から引き継がない 独立している
    #name: string;

    constructor(name: string) {
        super(name + "👍");
        this.#name = name;
    }

    greet2() {
        console.log(`Hi, I'm ${this.#name}.`);
    }
}

let papyrus = new Person("Papyrus", "🐾");
// Hello, my name is Jeremy Bearimy! I like 🐾. と表示される
papyrus.greet();

let sans = new Person2("Sans");
// Hello, my name is Sans👍! と表示される
sans.greet();
// Hi, I'm Sans. と表示される
sans.greet2();

// クラスの外側からはアクセスできない
// error TS18013: Property '#name' is not accessible outside class 'Person' because it has a private identifier.
// papyrus.#name

// 型上にはちゃんと存在しているので、該当のprivate fieldがないとエラーになる
// 実質、Person型かそれを継承したクラスのインスタンス縛りになる
// error TS2741: Property '#name' is missing in type '{ fav: string; name: string; greet(): void; }' but required in type 'Person'.
// let mimic1: Person = {
//     fav: "",
//     name: "",
//     greet() { },
// };

// これはOK
let mimic2: Person = sans;

// esnextで .js 出せばわかるけど #name 的なのは宣言必須
