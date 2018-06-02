import { hello } from "./foo";

type fooType = typeof import("./foo");

function f(fn: typeof import("./foo").hello) {
    fn("import()");
}

const fn1: typeof import("./foo").hello = name => `Hi, ${name}`;
fn1("import()");

const fn2: fooType["hello"] = name => `Hi, ${name}`;
fn2("import()");

const fn3: typeof hello = name => `Hi, ${name}`;
fn3("import()");


const data1: import("./foo").Data = {
    id: "foo",
    content: "bar",
};

type Data = import("./foo").Data;
interface Data1 extends Data { }
// この書き方はinterfaceのsyntax的にダメ
// interface Data2 extends import("./foo").Data { }
