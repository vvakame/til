interface Define {
    <T extends Component>(name: string, definition: T): T;
    <T extends PropAttrs>(name: string, definition: ComponentProp<T>): Component & T;
}

export let define: Define;
export let vdom: any;
export class Component extends HTMLElement {
    connectedCallback: any;
    disconnectedCallback: any;
    attributeChangedCallback: any;
}

type PropAttrs = { [key: string]: PropAttr }
interface ComponentProp<T extends PropAttrs> {
    extends?: string;
    type?: Types;
    props?: T;
    attached?: (elem: any) => void;
    detached?: (elem: any) => void;
    render?: (elem: any) => void;
}

interface PropAttr { // use generics
    attribute?: boolean | string;
    coerce?: (value: any) => any;
    default?: any | ((elem: any, data: { name: any; }) => any);
    deserialize?: (value: any) => any;
    get?: (elem: any, data: { name: any; internalValue: any; }) => any;
    initial?: any | ((elem: any, data: { name: any; }) => any);
    serialize?: (value: any) => any;
    set?: (elem: any, data: { name: any; newValue: any; oldValue: any; }) => void;
    created?: (elem: any) => void;
    updated?: (elem: any, prevProps: any) => boolean;
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
    string(attr?: PropAttr): PropAttr;
    array(attr?: PropAttr): PropAttr;
    number(attr?: PropAttr): PropAttr;
};

export let h: any; // TODO

declare enum Types {
    ATTR,
    CLASS,
}
export let types: typeof Types;
