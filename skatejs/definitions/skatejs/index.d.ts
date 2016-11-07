export {} from "./jsx";

export let define: {
    <T extends typeof Component>(name: string, definition: T): T;
    <T extends PropAttrs<Component & T>>(name: string, definition: ComponentProp<any>): {new (): Component & T};
};
export let vdom: any;
export class Component extends HTMLElement {
    connectedCallback: any;
    disconnectedCallback: any;
    attributeChangedCallback: any;
}

type PropAttrs<El> = { [key: string]: PropAttr<El, any> }
export interface ComponentProp<T extends PropAttrs<T & Component>> {
    extends?: string;
    type?: Types;
    props?: T;
    attached?: (elem: T & Component) => void;
    detached?: (elem: T & Component) => void;
    render?: (elem: T & Component) => void;
}

export interface PropAttr<El, Prop> { // use generics
    attribute?: boolean | string;
    coerce?: (value: any) => Prop;
    default?: ((elem: El, data: { name: string; }) => Prop) | Prop;
    deserialize?: (value: any) => Prop;
    get?: (elem: El, data: { name: string; internalValue: Prop; }) => Prop;
    initial?: Prop | ((elem: El, data: { name: string; }) => Prop);
    serialize?: (value: Prop) => any;
    set?: (elem: El, data: { name: string; newValue: Prop; oldValue: Prop; }) => void;
    created?: (elem: El) => void;
    updated?: (elem: El, prevProps: PropAttrs<El>) => boolean;
}

export interface OnCreated {
    created(elem: this): void;
}

export interface OnAttached {
    attached(elem: this): void;
}

export interface OnDetached {
    detached(elem: this): void;
}

export let prop: {
    string<El>(attr?: PropAttr<El, string>): PropAttr<El, string>;
    array<El, T>(attr?: PropAttr<El, T[]>): PropAttr<El, T[]>;
    number<El>(attr?: PropAttr<El, number>): PropAttr<El, number>;
    boolean<El>(attr?: PropAttr<El, number>): PropAttr<El, number>;
};

export let h: any; // TODO

declare enum Types {
    ATTR,
    CLASS,
}
export let types: typeof Types;

export let symbols: {
    name: Symbol;
    shadowRoot: Symbol;
}
