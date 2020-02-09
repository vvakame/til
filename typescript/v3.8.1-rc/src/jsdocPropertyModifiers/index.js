// @ts-check

class Foo {
    constructor() {
        /** @private */
        this.stuff = 100;
    }

    printStuff() {
        console.log(this.stuff);
    }
}

// error TS2341: Property 'stuff' is private and only accessible within class 'Foo'.
new Foo().stuff;
