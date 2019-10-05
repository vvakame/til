import { Foo } from "./basic";

class FooImpl implements Foo {
    get x(): number {
        throw new Error("Method not implemented.");
    }    set x(val: number) {
        throw new Error("Method not implemented.");
    }

}
