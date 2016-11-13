export { } from "./jsx";

export as namespace skate;

export var define: {
    (name: string): ComponentFactory<Component>;
    <T extends Element<U>, U>(name: string, definition: T): ComponentFactory<T & U & Component>;
    <T extends Component>(name: string, definition: ComponentFactory<T>): ComponentFactory<T>;
};

export interface Element<Proto> {
    prototype?: Proto;
    props?: any;
    attached?: (elem: this) => void;
    detached?: (elem: this) => void;
    render?: (elem: this) => void;
    updated?: (elem: this, prev: this) => void;
    attributeChanged?: (elem: this, data: { name: string; newValue: any; oldValue: any }) => void;
    observedAttributes?: string[];
}

export interface Prop {
}

type PropAttr<T> = PropAttrSimple<T> | PropAttrConvert<any, T>;

export interface PropAttrSimple<T> {
    attribute?: boolean | string;
    coerce?: (value: T) => T;
    default?: T;
    deserialize?: (value: any) => T;
    get?: (elem: any, data: { name: string; internalValue: T; }) => T;
    initial?: T | ((elem: any, data: { name: string; }) => T);
    serialize?: (value: T) => any;
    set?: (elem: any, data: { name: string; newValue: T; oldValue: T }) => void;
    created?: (elem: any) => void;
    update?: (elem: any, prevProps: any) => boolean | undefined;
}

export interface PropAttrConvert<S, D> {
    attribute?: boolean | string;
    coerce?: (value: S) => D;
    default?: (elem: any, data: S) => D;
    deserialize?: (value: any) => D;
    get?: (elem: any, data: { name: string; internalValue: D; }) => D;
    initial?: D | ((elem: any, data: { name: string; }) => D);
    serialize?: (value: D) => any;
    set?: (elem: any, data: { name: string; newValue: D; oldValue: D }) => void;
    created?: (elem: any) => void;
    update?: (elem: any, prevProps: any) => boolean | undefined;
}

export var prop: {
    number(p?: PropAttr<number>): Prop;
    string(p?: PropAttr<string>): Prop;
    boolean(p?: PropAttr<boolean>): Prop;
    array(p?: PropAttr<any[]>): Prop;
};

export interface ComponentFactory<T extends Component> {
    ["____skate_name"]?: string;

    new (): T;
    extend<T2 extends Element<U>, U>(defintion: T2): ComponentFactory<T & T2 & U>;
    updated?: (elem: T, prev: T) => boolean;
}

export class Component extends HTMLElement {
    static extend<T extends Element<U>, U>(defintion: T): Component & T & U;
    static updated(elem: Component, prev: Component): boolean;
}

export var h: HyperscriptBuilder;

export var emit: {
    (elem: Component, eventName: string, eventData?: any): void;
};

export var link: {
    (elem: Component, propSpec?: string): void;
};

export var props: {
    (elem: Component, props?: any): void;
};

export var symbols: {
    name: "____skate_name";
};

export var vdom: {
    builder(): HyperscriptBuilder;
    builder(...elements: any[]): any[];

    elementOpen(type: ComponentFactory<any>): void;
};

export interface HyperscriptBuilder {
    (type: string | ComponentFactory<any> | ((...props: any[]) => any), ...args: any[]): {};
}
