// from http://w3c.github.io/webcomponents/spec/shadow/
// ShadowDOM v1

export { }

declare global {
    interface DocumentOrShadowRoot {
        getSelection(): Selection;
        elementFromPoint(x: number, y: number): Element | null;
        elementsFromPoint(x: number, y: number): Element[];
        caretPositionFromPoint(x: number, y: number): CaretPosition | null;
        readonly activeElement: Element | null;
        readonly styleSheets: StyleSheetList | null;
    }

    interface Document extends DocumentOrShadowRoot {
    }
    interface ShadowRoot extends DocumentOrShadowRoot {
    }

    interface CaretPosition { }

    interface ShadowRoot extends DocumentFragment {
        readonly host: Element;
        innerHTML: string;
    }

    interface Element {
        attachShadow(shadowRootInitDict: ShadowRootInit): ShadowRoot;
        readonly assignedSlot: HTMLSlotElement | null;
        slot: string;
        readonly shadowRoot: ShadowRoot | null;
    }

    interface ShadowRootInit {
        mode: ShadowRootMode;
        delegatesFocus?: boolean; // default false
    }

    type ShadowRootMode = "open" | "closed";

    interface Text {
        readonly assignedSlot: HTMLSlotElement | null;
    }

    interface HTMLSlotElement extends HTMLElement {
        name: string;
        assignedNodes(options?: AssignedNodesOptions): Node[];
    }

    interface AssignedNodesOptions {
        flatten?: boolean; // default false
    }

    interface EventInit {
        scoped?: boolean; // default false
    }

    interface Event {
        deepPath(): EventTarget[];
        readonly scoped: boolean;
    }

    interface Document {
        createElement(tagName: "slot"): HTMLSlotElement;
    }
}
