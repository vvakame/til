declare global {
    interface ImportMeta {
        foo: string;
    }
}

import.meta.foo;
// これは定義されてないのでエラー
// import.meta.notExist;
