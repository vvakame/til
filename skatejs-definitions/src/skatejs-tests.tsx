// <reference types="skatejs" />

{
    skate.define('x-hello', {
        props: {
            name: { attribute: true }
        },
        render(elem: any) {
            return skate.h('div', `Hello, ${elem.name}`);
        }
    });
}

{
    const sym = Symbol();

    skate.define('x-counter', {
        props: {
            // By declaring the property an attribute, we can now pass an initial value
            // for the count as part of the HTML.
            count: skate.prop.number({ attribute: true })
        },
        attached(elem: any) {
            // We use a symbol so we don't pollute the element's namespace.
            elem[sym] = setInterval(() => ++elem.count, 1000);
        },
        detached(elem: any) {
            // If we didn't clean up after ourselves, we'd continue to render
            // unnecessarily.
            clearInterval(elem[sym]);
        },
        render(elem: any) {
            // By separating the strings (and not using template literals or string
            // concatenation) it ensures the strings are diffed indepenedently. If
            // you select "Count" with your mouse, it will not deselect whenr endered.
            return skate.h('div', 'Count ', elem.count);
        }
    });
}

const {define, Component} = skate;

{
    const MyComponent = define('my-component', class extends Component { });
}
{
    const MyComponent = define('my-component', {});
}
{
    const MyComponent = define('my-component', Component.extend({}));
}
{
    const MyComponent1 = define('my-component-1', {});
    const MyComponent2 = define('my-component-2', MyComponent1.extend({}));
}
{
    const MyComponent = define('my-component', {});
    const myElement = new MyComponent();
}

{
    skate.define('my-component', {
        prototype: {
            get someProperty() { return null; },
            set someProperty(v: any) { },
            someMethod() { },
        }
    });
}

{
    let props: any;
    skate.define('my-component', {
        props: { ...props }
    });
}

{
    // TODO try again after https://github.com/Microsoft/TypeScript/pull/12114 merged
    const Elem: any = skate.define('x-component', {
        props: {
            str: skate.prop.string(),
            arr: skate.prop.array(),
        },
        render() {

        },
    });

    const elem = new Elem();

    // Re-renders:
    elem.str = 'updated';

    // Will not re-render:
    elem.arr.push('something');

    // Will re-render:
    elem.arr = elem.arr.concat('something');
}

{
    skate.define('my-component', class extends skate.Component {
        static updated(elem: any, prev: any) {
            // You can reuse the original check if you want as part of your new check.
            // You could also call it directly if not extending: skate.Component().
            // TODO return super.updated(elem, prev) && myCustomCheck(elem, prev);
            return true;
        }
    });
}

{
    skate.define('my-component', {
        props: {
            name: skate.prop.string()
        },
        updated(elem, prev: any) {
            if (prev.name !== elem.name) {
                skate.emit(elem, 'name-changed', { detail: prev });
            }
        }
    });
}

{
    skate.define('my-component', {
        render(elem: any) {
            return skate.h('p', `My name is ${elem.tagName}.`);
        }
    });
}

{
    skate.define('my-component', {
        render(elem) {
            return [
                skate.h('paragraph 1'),
                skate.h('paragraph 2'),
            ];
        }
    });
}

{
    skate.define('my-component', {
        attached(elem) { }
    });
}

{
    skate.define('my-component', {
        detached(elem) { }
    });
}

{
    skate.define('my-component', {
        attributeChanged(elem, data) {
            if (data.oldValue === undefined) {
                // created
            } else if (data.newValue === undefined) {
                // removed
            } else {
                // updated
            }
        }
    });
}

{
    skate.define('my-component', {
        observedAttributes: ['some-attribute'],
        attributeChanged() { }
    });
}

{
    skate.define('my-component', {
        props: {
            someAttribute: { attribute: true }
        }
    });
}

{
    skate.define('x-tabs', {
        render(elem) {
            return skate.h('x-tab', { onSelect: () => { } });
        }
    });

    const {emit} = skate;
    skate.define('x-tab', {
        render(elem: any) {
            return skate.h('a', { onClick: () => emit(elem, 'select') });
        }
    });
}

{
    let elem: skate.Component = null as any;
    skate.emit(elem, 'event', {
        composed: false,
        bubbles: false,
        cancelable: false
    });
}

{
    let elem: skate.Component = null as any;
    skate.emit(elem, 'event', {
        detail: {
            data: 'my-data'
        }
    });
}

{ // TODO documentation issue?
    skate.define('my-input', /* ?? function () */ {
        props: {
            value: { attribute: true }
        },
        render(elem: any) {
            return skate.h('input', { onChange: skate.link(elem), type: 'text' });
        }
    });
}

{
    let elem: skate.Component = null as any;
    skate.h('input', { name: 'someValue', onChange: skate.link(elem), type: 'text' });
}

{
    let elem: skate.Component = null as any;
    skate.link(elem, 'someValue');
}

{
    let elem: skate.Component = null as any;
    skate.link(elem, 'obj.someValue')
}

{
    let elem: skate.Component = null as any;
    skate.h('input', { name: 'someValue', onChange: skate.link(elem, 'obj.'), type: 'text' });
}

{ // TODO documentation issue?
    let elem: skate.Component = null as any;
    const linkage = skate.link(elem, 'obj.');
    skate.h('input', { name: 'someValue1', onChange: linkage, type: 'text' });
    skate.h('input', { name: 'someValue2', onChange: linkage, type: 'checkbox' });
    skate.h('input', { name: 'someValue3', onChange: linkage, type: 'radio' });
    skate.h('select', { name: 'someValue4', onChange: linkage },
        skate.h('option', { value: 2 }, 'Option 2'), // ;
        skate.h('option', { value: 1 }, 'Option 1'), // ;
    );
}

{
    skate.prop.boolean();
}

{
    skate.prop.boolean({
        coerce() {
            // coerce it differently than the default way
            return true;
        },
        set() {
            // do something when set
        }
    });
}

{
    // TODO try again after https://github.com/Microsoft/TypeScript/pull/12114 merged
    const { define, props } = skate;

    const Elem = define('my-element', {
        props: {
            prop1: null
        }
    });
    const elem = new Elem();

    // Set any property you want.
    props(elem, {
        prop1: 'value 1',
        prop2: 'value 2'
    });

    // Only returns props you've defined on your component.
    // { prop1: 'value 1' }
    props(elem);
}

{
    const MyComponent1 = skate.define('my-component', {});

    // my-component
    console.log(MyComponent1[skate.symbols.name]);

    // If re-registering in HMR...
    const MyComponent2 = skate.define('my-component', {});

    // my-component-1
    console.log(MyComponent2[skate.symbols.name]);
}

{ // TODO is this exists?
    skate.define('my-component', {
        render() {
            return skate.h('p', 'test');
        },
        ready(elem: skate.Component) {
            // #shadow-root
            //   <p>test</p>
            // TODO elem[skate.symbols.shadowRoot];
        }
    });
}

{
    skate.define('my-component', {
        render() {
            return skate.h('p', { style: { fontWeight: 'bold' } }, 'Hello!');
        }
    });
}

{
    const skatex = { createElement: skate.h }; // for --reactNamespace skatex
    define('my-component', {
        render() {
            return <p>Hello!</p>;
        }
    });
}

{
    const {vdom} = skate;
    const h = vdom.builder();
    define('my-component', {
        render() {
            return h('div', { id: 'test', }, h('p', 'test'));
        }
    });
}

{
    const [div, p] = skate.vdom.builder('div', 'p');
    define('my-component', {
        render() {
            return div({ id: 'mydiv' }, p('test'));
        }
    });
}

{
    const MyElement = skate.define('my-element');

    // Renders <my-element />
    skate.h(MyElement);
}

{
    const MyElement = skate.define('my-element');
    skate.vdom.elementOpen(MyElement);
}

{
    const MyElement = () => skate.h('div', 'Hello, World!');

    // Renders <div>Hello, World!</div>
    skate.h(MyElement);
}

{
    const MyElement = (props: any) => skate.h('div', `Hello, ${props.name}!`);

    // Renders <div>Hello, Bob!</div>
    skate.h(MyElement, { name: 'Bob' });
}

{
    const MyElement = (props: any, chren: any) => skate.h('div', 'Hello, ', chren, '!');

    // Renders <div>Hello, Mary!</div>
    skate.h(MyElement, 'Mary');
}

{
    skate.h('ul',
        skate.h('li', { key: 0 }),
        skate.h('li', { key: 1 }),
    );
}

{
    const onClick = console.log;
    skate.h('button', { onClick });
    skate.h('button', { 'on-click': onClick });
}

{
    const onClick = console.log;
    skate.h('button', { onclick: onClick });
}

{
    skate.define('x-element', {
        created(elem: any) {
            elem.addEventListener('change', elem.handleChange);
        },
        prototype: {
            handleChange(e: any) {
                // `this` is the element.
                // The event is passed as the only argument.
            }
        }
    });
}

{
    const ref = (button: HTMLButtonElement) => button.addEventListener('click', console.log);
    skate.h('button', { ref });
}

{
    const ref = console.log;
    skate.define('my-element', {
        render() {
            return skate.h('div', { ref });
        }
    });
}

{
    skate.define('my-element', {
        render() {
            const ref = console.log;
            return skate.h('div', { ref });
        }
    });
}

{
    skate.h('div', { ref: (e: HTMLElement) => (e.innerHTML = '<p>oh no you didn\'t</p>'), skip: true });
}

{
    skate.h('div', { statics: ['attr1', 'prop2'] });
}

{
    skate.define('my-component', {
        props: {
            // Links the `name` property to the `name` attribute.
            name: { attribute: true }
        }
    });
}

{
    skate.define('my-component', {
        props: {
            values: {
                attribute: true,
                deserialize(val: any) {
                    return val.split(',');
                },
                serialize(val: any) {
                    return val.join(',');
                }
            }
        }
    });
}

{
    const scoped = function scoped(elem: any) { }

    skate.define('x-element', {
        created(elem: any) {
            scoped(elem);
            elem._privateButNotReally();
        },
        prototype: {
            _privateButNotReally() { }
        }
    });
}

{
    const sym = Symbol();

    skate.define('x-element', {
        created(elem: any) {
            elem[sym]();
        },
        prototype: {
            [sym]() { }
        }
    });
}

{
    const map = new WeakMap();

    skate.define('x-element', {
        created(elem: any) {
            map.set(elem, 'some data');
        },
        render(elem) {
            // Renders: "<div>some data</div>"
            return skate.h('div', map.get(elem));
        }
    });
}

{
    const skatex = { createElement: skate.h }; // for --reactNamespace skatex

    skate.define('x-form', {
        render() {
            return (
                <form>
                    <slot />
                </form>
            );
        }
    });

    skate.define('x-button', {
        render() {
            return (
                <button>
                    <slot />
                </button>
            );
        }
    });
}

{
    const skatex = { createElement: skate.h }; // for --reactNamespace skatex

    const onclick = function onclick(e: any) {
        if (e.target.getAttribute('type') === 'submit') {
            // do something submitty
        }
    }

    skate.define('x-form', {
        render() {
            return (<form onClick={onclick} ></form>);
        }
    });
}

{
    skate.define('x-component', {
        updated(elem: any, prev: any) {
            // Notify any listeners that the component updated. At this point the
            // listener can update the component's props without fear that this will
            // cause recursion - because it's prevented internally - and it will
            // proceed past this point with the updated props.
            const canRender = skate.emit(elem, 'updated', { detail: prev });

            // This can be custom, or just reuse the default implementation. Since we
            // emitted the event and listeners had a chance to update the component,
            // this will get called with the updated state.
            return canRender && skate.Component.updated(elem, prev);
        }
    });
}

{
    skate.define('x-component', {
        render() {
            return [
                skate.h('style', '.my-class { display: block; }'),
                skate.h('div', { class: 'my-class' }),
            ];
        }
    });
}

{
    skate.define('x-component', class extends skate.Component {
        static created(elem: any) { }
        static attached(elem: any) { }
        static detached(elem: any) { }
    });
}

{
    skate.define('x-component', class extends skate.Component {
        static get props() {
            return {
                myProp: { attribute: true }
            };
        }
    });
}

{
    skate.define('x-component', class extends skate.Component {
        static props = {
            myProp: { attribute: true }
        }
    });
}

{
    skate.define('x-component', class extends skate.Component {
        myProp1 = 'some value'
        get myProp2() { return 'another value'; }
        myMethod() { }
    });
}
