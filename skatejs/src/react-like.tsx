import { define, h } from "skatejs";

const React = { createElement: h };

export default define("my-component", {
    render() {
        return <p>Hello!</p>;
    }
});
