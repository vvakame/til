import * as skate from "skatejs";

function checkAttr<T>(attr: skate.PropAttr<PropExampleComponent, T>) {
    return attr;
}

interface PropExampleProps {
    a: skate.PropAttr<PropExampleComponent, string>;
    b: skate.PropAttr<PropExampleComponent, string>;
    c: skate.PropAttr<PropExampleComponent, string>;
    d: skate.PropAttr<PropExampleComponent, string>;
    e: skate.PropAttr<PropExampleComponent, number>;
    f: skate.PropAttr<PropExampleComponent, string>;
    g: skate.PropAttr<PropExampleComponent, string>;
    h: skate.PropAttr<PropExampleComponent, string>;
    i: skate.PropAttr<PropExampleComponent, string>;
    j: skate.PropAttr<PropExampleComponent, string>;
    k: skate.PropAttr<PropExampleComponent, string>;
}

class PropExampleComponent extends skate.Component implements skate.OnRenderCallback {
    static get props(): PropExampleProps {
        return {
            a: {
                attribute: true,
            },
            b: {
                attribute: "b-modify",
            },
            c: {
                coerce: value => {
                    return value + "2";
                },
            },
            d: {
                default: "default",
            },
            e: {
                default: (elem, data) => {
                    return Math.random();
                },
            },
            f: {
                deserialize: value => {
                    return `${value}`;
                },
            },
            g: {
                get: (elem: any, data: any) => {
                    return `prefix_${data.internalValue}`;
                },
            },
            h: {
                initial: 'initial value',
            },
            i: {
                initial: (elem, data) => {
                    return 'initial value';
                },
            },
            j: {
                serialize: value => {
                    return `${1}`;
                },
            },
            k: {
                get: (elem: any, data: any) => {
                    return `prefix_${data.internalValue}`;
                },
                set: (elem, data) => {
                    console.log(elem, data);
                },
            },
        };
    }

    a: any;
    b: any;
    c: any;
    d: any;
    e: any;
    f: any;
    g: any;
    h: any;
    i: any;
    j: any;
    k: any;

    renderCallback(): any {
        return (
            <div>
                <div>a: <span>{this.a}</span></div>
                <div>b: <span>{this.b}</span></div>
                <div>c: <span>{this.c}</span></div>
                <div>d: <span>{this.d}</span></div>
                <div>e: <span>{this.e}</span></div>
                <div>f: <span>{this.f}</span></div>
                <div>g: <span>{this.g}</span></div>
                <div>h: <span>{this.h}</span></div>
                <div>i: <span>{this.i}</span></div>
                <div>j: <span>{this.j}</span></div>
                <div>k: <span>{this.k}</span></div>
            </div>
        );
    }
}
customElements.define("x-prop-example", PropExampleComponent);

class PropContainerComponent extends skate.Component implements skate.OnRenderCallback {
    renderCallback(): any {
        const anyProps1: any = {
            a: 1,
            b: 2,
            c: 3,
            d: 4,
            e: 5,
            f: 6,
            g: 7,
            h: 8,
            i: 9,
            j: 10,
            k: 11,
        };
        const anyProps2: any = {};
        return (
            <div>
                <PropExampleComponent {...anyProps1} />
                <PropExampleComponent {...anyProps2} />
            </div>
        );
    }
}
customElements.define("x-prop-container", PropContainerComponent);
