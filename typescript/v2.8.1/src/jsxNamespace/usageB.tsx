/** @jsx Lib.myjsx */
import Lib from "./lib";

export function render() {
    return <span></span>;
}

// "__myjsxBrandB" | "children" | "props" と評価される
type El = keyof ReturnType<typeof render>;
