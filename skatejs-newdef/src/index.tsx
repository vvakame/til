import "skatejs-web-components";
import * as skate from "skatejs";

skate.define("x-test", class extends skate.Component {
    renderCallback() {
        return (
            <div>hello</div>
        );
    }
});
