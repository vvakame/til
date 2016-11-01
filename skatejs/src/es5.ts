import * as skate from "skatejs";

const sym = Symbol();

skate.define('x-counter', {
    props: {
        // By declaring the property an attribute, we can now pass an initial value
        // for the count as part of the HTML.
        count: skate.prop.number({ attribute: true })
    },
    attached(elem) {
        // We use a symbol so we don't pollute the element's namespace.
        elem[sym] = setInterval(() => ++elem.count, 1000);
    },
    detached(elem) {
        // If we didn't clean up after ourselves, we'd continue to render
        // unnecessarily.
        clearInterval(elem[sym]);
    },
    render(elem) {
        // By separating the strings (and not using template literals or string
        // concatenation) it ensures the strings are diffed indepenedently. If
        // you select "Count" with your mouse, it will not deselect whenr endered.
        return skate.h('div', 'Count ', elem.count);
    }
});
