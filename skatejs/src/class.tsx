import * as skate from "skatejs";

const React = { createElement: skate.h };

skate.define("x-component2", class Component2 extends skate.Component {
    static get props() {
        return {
            myProp: { attribute: true }
        };
    }

    static render(elem: Component2) {
        return (
            <div>
                <form>
                </form>
                <p>Hello!</p>
            </div>
        );
    }
});
