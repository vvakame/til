import { define, h } from "skatejs";

const React = { createElement: h };

define("my-component", {
    render() {
        return <p>Hello!</p>;
    }
});
