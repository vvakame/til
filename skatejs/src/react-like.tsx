import * as skate from "skatejs";

const skatex = { createElement: skate.h };

export default skate.define("my-component", {
    render() {
        return <p>Hello!</p>;
    }
});
