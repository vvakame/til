import React from "react";

export function render() {
    return <span></span>;
}

// "type" | "props" | "key" と評価される
type El = keyof ReturnType<typeof render>;
