/** @jsx myjsx */
import { myjsx } from "./lib";

export function render() {
    return <span></span>;
}

// "__myjsxBrandA" | "children" | "props" と評価される
type El = keyof ReturnType<typeof render>;
