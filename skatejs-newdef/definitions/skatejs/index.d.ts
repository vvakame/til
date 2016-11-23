export { } from "./jsx";

export as namespace skate;

export interface OnUpdatedCallback {
    updatedCallback(previousProps: any): boolean | undefined | void;
}

export interface OnRenderCallback {
    renderCallback(): any;
}

export interface OnRenderedCallback {
    renderedCallback(): void;
}

export class Component extends HTMLElement {
    static readonly props: { [name: string]: PropAttr<any> };
    static readonly observedAttributes: string[];

    // Custom Elements v1
    connectedCallback(): void;
    disconnectedCallback(): void;
    attributeChangedCallback(name: string, oldValue: any, newValue: any): void;

    // SkateJS
    updated(prev: any): boolean;
}

export interface PropAttr<T> {
    attribute?: boolean | string;
    coerce?: (value: any) => T | null | undefined | void;
    default?: ((elem: any, data: { name: string; }) => T) | T;
    deserialize?: (value: any) => T;
    get?: (elem: any, data: { name: string; internalValue: T; }) => T;
    initial?: T | ((elem: any, data: { name: string; }) => T);
    serialize?: (value: T) => any;
    set?: (elem: any, data: { name: string; newValue: T; oldValue: T; }) => void;
    created?: (elem: any) => void;
    updated?: (elem: any, prevProps: any) => boolean;
}

export var define: {
    (name: string, ctor: Function): any;
    (ctor: Function): any;
};

export interface EmitOpts {
    bubbles?: boolean;
    cancelable?: boolean;
    composed?: boolean;
    detail?: any;
}
export function emit(elem: EventTarget, name: string, opts?: EmitOpts): void;

export function link(elem: Component, target?: string): (e: Event) => void;

type VDOMElementTName = string | typeof Component | typeof vdom.element | { id: string; };
type VDOMElementChild = Function | string | number;
type VDOMElementSet = VDOMElementChild | VDOMElementChild[];

export var h: typeof vdom.element;

export var vdom: {
    element(tname: VDOMElementTName, attrs: { key: any; statics: any; } & any, ...chren: VDOMElementSet[]): Component | any;
    element(tname: VDOMElementTName, ...chren: VDOMElementSet[]): Component | any;
    builder(): typeof vdom.element;
    builder(...tags: string[]): (typeof vdom.element)[];

    attr: Function;
    elementClose: Function;
    elementOpen: Function;
    elementOpenEnd: Function;
    elementOpenStart: Function;
    elementVoid: Function;
    text: Function;
};

export function ready(elem: Component, done: (c: Component) => void): void;

// DEPRECATED
// export var symbols: any;

export var prop: {
    create<T>(attr: PropAttr<T>): PropAttr<T> & ((attr: PropAttr<T>) => PropAttr<T>);

    number(attr?: PropAttr<number>): PropAttr<number>;
    boolean(attr?: PropAttr<boolean>): PropAttr<boolean>;
    string(attr?: PropAttr<string>): PropAttr<string>;
    array(attr?: PropAttr<any[]>): PropAttr<any[]>;
};

export function props(elem: Component, props?: any): void;
