import "skatejs-web-components";
import * as skate from "skatejs";

const skatex = { createElement: skate.h };

import "./es5";
import ClassLike from "./class";

export default skate.define("x-app", class App extends skate.Component {
    render() {
        // for https://github.com/Microsoft/TypeScript/issues/7004
        // 型チェックが行われると大量にエラーが出るのでanyで殺す 堅牢さに影響はないはず
        const anyProps: any = {};
        return (
            <div>
                <ClassLike myProp1={100} {...anyProps}></ClassLike>
            </div>
        );
    }

    static render(elem: any) {
        return elem.render();
    }
});
