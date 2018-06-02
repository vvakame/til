import { Component } from "react";

interface Props<T> {
    data: T;
}

class MyComponent<T> extends Component<Props<T>> {
    render() {
        return <div>{this.props.data}</div>;
    }
}

<MyComponent<number> data={12} />
