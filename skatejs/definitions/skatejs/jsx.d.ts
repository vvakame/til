export { }

interface HTMLProps<T> extends HTMLAttributes<T>, ClassAttributes<T> {
}

interface Attributes {
}

interface ClassAttributes<T> extends Attributes {
}

interface HTMLAttributes<T> extends DOMAttributes<T> {
}

interface DOMAttributes<T> {
    // Clipboard Events
    onCopy?: ClipboardEventHandler<T>;
    onCut?: ClipboardEventHandler<T>;
    onPaste?: ClipboardEventHandler<T>;

    // Composition Events
    onCompositionEnd?: CompositionEventHandler<T>;
    onCompositionStart?: CompositionEventHandler<T>;
    onCompositionUpdate?: CompositionEventHandler<T>;

    // Focus Events
    onFocus?: FocusEventHandler<T>;
    onBlur?: FocusEventHandler<T>;

    // Form Events
    onChange?: FormEventHandler<T>;
    onInput?: FormEventHandler<T>;
    onSubmit?: FormEventHandler<T>;

    // Image Events
    onLoad?: ReactEventHandler<T>;
    onError?: ReactEventHandler<T>; // also a Media Event

    // Keyboard Events
    onKeyDown?: KeyboardEventHandler<T>;
    onKeyPress?: KeyboardEventHandler<T>;
    onKeyUp?: KeyboardEventHandler<T>;

    // Media Events
    onAbort?: ReactEventHandler<T>;
    onCanPlay?: ReactEventHandler<T>;
    onCanPlayThrough?: ReactEventHandler<T>;
    onDurationChange?: ReactEventHandler<T>;
    onEmptied?: ReactEventHandler<T>;
    onEncrypted?: ReactEventHandler<T>;
    onEnded?: ReactEventHandler<T>;
    onLoadedData?: ReactEventHandler<T>;
    onLoadedMetadata?: ReactEventHandler<T>;
    onLoadStart?: ReactEventHandler<T>;
    onPause?: ReactEventHandler<T>;
    onPlay?: ReactEventHandler<T>;
    onPlaying?: ReactEventHandler<T>;
    onProgress?: ReactEventHandler<T>;
    onRateChange?: ReactEventHandler<T>;
    onSeeked?: ReactEventHandler<T>;
    onSeeking?: ReactEventHandler<T>;
    onStalled?: ReactEventHandler<T>;
    onSuspend?: ReactEventHandler<T>;
    onTimeUpdate?: ReactEventHandler<T>;
    onVolumeChange?: ReactEventHandler<T>;
    onWaiting?: ReactEventHandler<T>;

    // MouseEvents
    onClick?: MouseEventHandler<T>;
    onContextMenu?: MouseEventHandler<T>;
    onDoubleClick?: MouseEventHandler<T>;
    onDrag?: DragEventHandler<T>;
    onDragEnd?: DragEventHandler<T>;
    onDragEnter?: DragEventHandler<T>;
    onDragExit?: DragEventHandler<T>;
    onDragLeave?: DragEventHandler<T>;
    onDragOver?: DragEventHandler<T>;
    onDragStart?: DragEventHandler<T>;
    onDrop?: DragEventHandler<T>;
    onMouseDown?: MouseEventHandler<T>;
    onMouseEnter?: MouseEventHandler<T>;
    onMouseLeave?: MouseEventHandler<T>;
    onMouseMove?: MouseEventHandler<T>;
    onMouseOut?: MouseEventHandler<T>;
    onMouseOver?: MouseEventHandler<T>;
    onMouseUp?: MouseEventHandler<T>;

    // Selection Events
    onSelect?: ReactEventHandler<T>;

    // Touch Events
    onTouchCancel?: TouchEventHandler<T>;
    onTouchEnd?: TouchEventHandler<T>;
    onTouchMove?: TouchEventHandler<T>;
    onTouchStart?: TouchEventHandler<T>;

    // UI Events
    onScroll?: UIEventHandler<T>;

    // Wheel Events
    onWheel?: WheelEventHandler<T>;

    // Animation Events
    onAnimationStart?: AnimationEventHandler;
    onAnimationEnd?: AnimationEventHandler;
    onAnimationIteration?: AnimationEventHandler;

    // Transition Events
    onTransitionEnd?: TransitionEventHandler;
}

interface SyntheticEvent<T> {
    bubbles: boolean;
    currentTarget: EventTarget & T;
    cancelable: boolean;
    defaultPrevented: boolean;
    eventPhase: number;
    isTrusted: boolean;
    nativeEvent: Event;
    preventDefault(): void;
    isDefaultPrevented(): boolean;
    stopPropagation(): void;
    isPropagationStopped(): boolean;
    persist(): void;
    target: EventTarget;
    timeStamp: Date;
    type: string;
}

interface ClipboardEvent<T> extends SyntheticEvent<T> {
    clipboardData: DataTransfer;
}

interface CompositionEvent<T> extends SyntheticEvent<T> {
    data: string;
}

interface DragEvent<T> extends MouseEvent<T> {
    dataTransfer: DataTransfer;
}

interface FocusEvent<T> extends SyntheticEvent<T> {
    relatedTarget: EventTarget;
}

interface FormEvent<T> extends SyntheticEvent<T> {
}

interface KeyboardEvent<T> extends SyntheticEvent<T> {
    altKey: boolean;
    charCode: number;
    ctrlKey: boolean;
    getModifierState(key: string): boolean;
    key: string;
    keyCode: number;
    locale: string;
    location: number;
    metaKey: boolean;
    repeat: boolean;
    shiftKey: boolean;
    which: number;
}

interface MouseEvent<T> extends SyntheticEvent<T> {
    altKey: boolean;
    button: number;
    buttons: number;
    clientX: number;
    clientY: number;
    ctrlKey: boolean;
    getModifierState(key: string): boolean;
    metaKey: boolean;
    pageX: number;
    pageY: number;
    relatedTarget: EventTarget;
    screenX: number;
    screenY: number;
    shiftKey: boolean;
}

interface TouchEvent<T> extends SyntheticEvent<T> {
    altKey: boolean;
    changedTouches: TouchList;
    ctrlKey: boolean;
    getModifierState(key: string): boolean;
    metaKey: boolean;
    shiftKey: boolean;
    targetTouches: TouchList;
    touches: TouchList;
}

interface UIEvent<T> extends SyntheticEvent<T> {
    detail: number;
    // view: AbstractView;
}

interface WheelEvent<T> extends MouseEvent<T> {
    deltaMode: number;
    deltaX: number;
    deltaY: number;
    deltaZ: number;
}

interface AnimationEvent extends SyntheticEvent<{}> {
    animationName: string;
    pseudoElement: string;
    elapsedTime: number;
}

interface TransitionEvent extends SyntheticEvent<{}> {
    propertyName: string;
    pseudoElement: string;
    elapsedTime: number;
}

interface EventHandler<E extends SyntheticEvent<any>> {
    (event: E): void;
}

// TODO naming
type ReactEventHandler<T> = EventHandler<SyntheticEvent<T>>;

type ClipboardEventHandler<T> = EventHandler<ClipboardEvent<T>>;
type CompositionEventHandler<T> = EventHandler<CompositionEvent<T>>;
type DragEventHandler<T> = EventHandler<DragEvent<T>>;
type FocusEventHandler<T> = EventHandler<FocusEvent<T>>;
type FormEventHandler<T> = EventHandler<FormEvent<T>>;
type KeyboardEventHandler<T> = EventHandler<KeyboardEvent<T>>;
type MouseEventHandler<T> = EventHandler<MouseEvent<T>>;
type TouchEventHandler<T> = EventHandler<TouchEvent<T>>;
type UIEventHandler<T> = EventHandler<UIEvent<T>>;
type WheelEventHandler<T> = EventHandler<WheelEvent<T>>;
type AnimationEventHandler = EventHandler<AnimationEvent>;
type TransitionEventHandler = EventHandler<TransitionEvent>;


declare global {
    // https://www.typescriptlang.org/docs/handbook/jsx.html
    namespace JSX {
        interface ElementClass {
            render(elem: this): JSX.Element; // TODO temp
        }

        interface ElementAttributesProperty {
            '': any;
        }

        interface Element {
        }

        interface IntrinsicElements {
            a: HTMLProps<HTMLAnchorElement>;
            abbr: HTMLProps<HTMLElement>;
            address: HTMLProps<HTMLElement>;
            area: HTMLProps<HTMLAreaElement>;
            article: HTMLProps<HTMLElement>;
            aside: HTMLProps<HTMLElement>;
            audio: HTMLProps<HTMLAudioElement>;
            b: HTMLProps<HTMLElement>;
            base: HTMLProps<HTMLBaseElement>;
            bdi: HTMLProps<HTMLElement>;
            bdo: HTMLProps<HTMLElement>;
            big: HTMLProps<HTMLElement>;
            blockquote: HTMLProps<HTMLElement>;
            body: HTMLProps<HTMLBodyElement>;
            br: HTMLProps<HTMLBRElement>;
            button: HTMLProps<HTMLButtonElement>;
            canvas: HTMLProps<HTMLCanvasElement>;
            caption: HTMLProps<HTMLElement>;
            cite: HTMLProps<HTMLElement>;
            code: HTMLProps<HTMLElement>;
            col: HTMLProps<HTMLTableColElement>;
            colgroup: HTMLProps<HTMLTableColElement>;
            data: HTMLProps<HTMLElement>;
            datalist: HTMLProps<HTMLDataListElement>;
            dd: HTMLProps<HTMLElement>;
            del: HTMLProps<HTMLElement>;
            details: HTMLProps<HTMLElement>;
            dfn: HTMLProps<HTMLElement>;
            dialog: HTMLProps<HTMLElement>;
            div: HTMLProps<HTMLDivElement>;
            dl: HTMLProps<HTMLDListElement>;
            dt: HTMLProps<HTMLElement>;
            em: HTMLProps<HTMLElement>;
            embed: HTMLProps<HTMLEmbedElement>;
            fieldset: HTMLProps<HTMLFieldSetElement>;
            figcaption: HTMLProps<HTMLElement>;
            figure: HTMLProps<HTMLElement>;
            footer: HTMLProps<HTMLElement>;
            form: HTMLProps<HTMLFormElement>;
            h1: HTMLProps<HTMLHeadingElement>;
            h2: HTMLProps<HTMLHeadingElement>;
            h3: HTMLProps<HTMLHeadingElement>;
            h4: HTMLProps<HTMLHeadingElement>;
            h5: HTMLProps<HTMLHeadingElement>;
            h6: HTMLProps<HTMLHeadingElement>;
            head: HTMLProps<HTMLHeadElement>;
            header: HTMLProps<HTMLElement>;
            hgroup: HTMLProps<HTMLElement>;
            hr: HTMLProps<HTMLHRElement>;
            html: HTMLProps<HTMLHtmlElement>;
            i: HTMLProps<HTMLElement>;
            iframe: HTMLProps<HTMLIFrameElement>;
            img: HTMLProps<HTMLImageElement>;
            input: HTMLProps<HTMLInputElement>;
            ins: HTMLProps<HTMLModElement>;
            kbd: HTMLProps<HTMLElement>;
            keygen: HTMLProps<HTMLElement>;
            label: HTMLProps<HTMLLabelElement>;
            legend: HTMLProps<HTMLLegendElement>;
            li: HTMLProps<HTMLLIElement>;
            link: HTMLProps<HTMLLinkElement>;
            main: HTMLProps<HTMLElement>;
            map: HTMLProps<HTMLMapElement>;
            mark: HTMLProps<HTMLElement>;
            menu: HTMLProps<HTMLElement>;
            menuitem: HTMLProps<HTMLElement>;
            meta: HTMLProps<HTMLMetaElement>;
            meter: HTMLProps<HTMLElement>;
            nav: HTMLProps<HTMLElement>;
            noscript: HTMLProps<HTMLElement>;
            object: HTMLProps<HTMLObjectElement>;
            ol: HTMLProps<HTMLOListElement>;
            optgroup: HTMLProps<HTMLOptGroupElement>;
            option: HTMLProps<HTMLOptionElement>;
            output: HTMLProps<HTMLElement>;
            p: HTMLProps<HTMLParagraphElement>;
            param: HTMLProps<HTMLParamElement>;
            picture: HTMLProps<HTMLElement>;
            pre: HTMLProps<HTMLPreElement>;
            progress: HTMLProps<HTMLProgressElement>;
            q: HTMLProps<HTMLQuoteElement>;
            rp: HTMLProps<HTMLElement>;
            rt: HTMLProps<HTMLElement>;
            ruby: HTMLProps<HTMLElement>;
            s: HTMLProps<HTMLElement>;
            samp: HTMLProps<HTMLElement>;
            script: HTMLProps<HTMLElement>;
            section: HTMLProps<HTMLElement>;
            select: HTMLProps<HTMLSelectElement>;
            small: HTMLProps<HTMLElement>;
            source: HTMLProps<HTMLSourceElement>;
            span: HTMLProps<HTMLSpanElement>;
            strong: HTMLProps<HTMLElement>;
            style: HTMLProps<HTMLStyleElement>;
            sub: HTMLProps<HTMLElement>;
            summary: HTMLProps<HTMLElement>;
            sup: HTMLProps<HTMLElement>;
            table: HTMLProps<HTMLTableElement>;
            tbody: HTMLProps<HTMLTableSectionElement>;
            td: HTMLProps<HTMLTableDataCellElement>;
            textarea: HTMLProps<HTMLTextAreaElement>;
            tfoot: HTMLProps<HTMLTableSectionElement>;
            th: HTMLProps<HTMLTableHeaderCellElement>;
            thead: HTMLProps<HTMLTableSectionElement>;
            time: HTMLProps<HTMLElement>;
            title: HTMLProps<HTMLTitleElement>;
            tr: HTMLProps<HTMLTableRowElement>;
            track: HTMLProps<HTMLTrackElement>;
            u: HTMLProps<HTMLElement>;
            ul: HTMLProps<HTMLUListElement>;
            "var": HTMLProps<HTMLElement>;
            video: HTMLProps<HTMLVideoElement>;
            wbr: HTMLProps<HTMLElement>;
        }
    }
}
