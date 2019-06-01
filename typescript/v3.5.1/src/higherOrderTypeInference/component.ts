type ComponentClass<P> = new (props: P) => Component<P>;
declare class Component<P> {
    props: P;
    constructor(props: P);
}

type Route = { route?: string[]; };
declare function myHoC<P>(C: ComponentClass<P>): ComponentClass<P & Route>;

type NestedProps<T> = { foo: number, stuff: T };
declare class GenericComponent<T> extends Component<NestedProps<T>> {
}

// MyComponent の型は
//   new <T>(props: NestedProps<T> & Route) => Component<NestedProps<T> & Route>
//     となる
// TypeScript v3.4 では
//   ComponentClass<NestedProps<{}> & Route>
//     だった
const MyComponent = myHoC(GenericComponent);
const c1 = new MyComponent({
    foo: 42,
    stuff: {
        name: "bar",
    },
});

c1.props.foo;
c1.props.stuff;
// TypeScript 3.5 ではちゃんと型が推論できる
// 3.4では error TS2339: Property 'name' does not exist on type '{}'.
c1.props.stuff.name;
c1.props.route;

export {};
