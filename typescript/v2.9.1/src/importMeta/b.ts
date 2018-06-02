declare global {
    interface ImportMeta {
        bar: number;
    }
}

// a.ts で定義されてる
import.meta.foo;
import.meta.bar;
// これはエラー
// import.meta.notExist;
