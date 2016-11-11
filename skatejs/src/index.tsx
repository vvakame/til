import "skatejs-web-components";
import * as skate from "skatejs";

const skatex = { createElement: skate.h };

import "./es5";
import ClassLile from "./class";

export default skate.define("x-app", class App extends skate.Component {
    render() {
        return (
            <div>
                <ClassLile></ClassLile>
            </div>
        );
    }

    static render(elem: any) {
        return elem.render();
    }
});
