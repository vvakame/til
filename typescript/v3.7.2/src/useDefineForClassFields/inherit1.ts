class Base {
    set data(value: string) {
        console.log("data changed to " + value);
    }
}

class Derived1 extends Base {
    // ES仕様に準拠した振る舞い(useDefineForClassFieldsを有効にした場合)では、 set accessor が上書きされる
    // 結果、 data changed to ... の出力が行われなくなる
    // これを防ぐため、コンパイルエラーが発生するようになった
    // error TS2610: Class 'Base' defines instance member accessor 'data', but extended class 'Derived1' defines it as instance member property.
    // data = "foobar";

    constructor() {
        super();
        // エラーを解消するためには、コンストラクタ内で値を初期化すればよい
        this.data = "foobar";
    }
}

let obj = new Derived1();
obj.data = "fizzbuzz";
