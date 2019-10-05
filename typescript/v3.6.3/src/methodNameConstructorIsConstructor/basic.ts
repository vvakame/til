class A {
    // 普通のコンストラクタ
    constructor() {
        console.log("A");
    }
}

class B {
    // コンストラクタと認識されるようになった
    "constructor"() {
        console.log("B");
    }

    // 2つ定義することになるのでエラーになる
    // error TS2392: Multiple constructor implementations are not allowed.
    // constructor() {}
}

class C {
    // computed propertyの場合コンストラクタにはならない
    ["constructor"]() {
        console.log("C");
    }

    // 重複定義にはならないのでエラーにならない
    constructor() {}
}

// A と表示
new A();
// B と表示
new B();
// なにも表示されない
new C(); // .constructor() が生えてる
