import type { Foo } from "./exports";

// 型のみなので当然継承できずに怒られる
// error TS1361: 'Foo' cannot be used as a value because it was imported using 'import type'.
export class Bar extends Foo {
    sayHello() {
        return "Hello, bar!";
    }
}
