class Foo {
    str: string;
    date: Date;
}

// Add 'undefined' type to property 'str' とかを実行するとこうなる
class FooA {
    str: string | undefined;
    date: Date | undefined;
}

// Add definite assignment to property 'str: string' とかを実行するとこうなる
class FooB {
    str!: string;
    date!: Date;
}

// Add initializer to property 'str' とかを実行するとこうなる
// Date の初期値は自明ではないので使えない…
class FooC {
    str: string = "";
    date: Date;
}
