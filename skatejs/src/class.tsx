import * as skate from "skatejs";

const React = { createElement: skate.h };

skate.define("x-component2", class extends skate.Component {
    myProp1: number;
    myProp2 = 0;

    constructor() {
        super();
        (window as any).xComponent2 = this;
    }

    click() {
        this.myProp1 += 1;
        this.myProp2 += 1;
    }

    render() {
        return (
            <div>
                <form>
                </form>
                <p><span>{this.myProp1}</span> & <span>{this.myProp2}</span></p>
                <button onClick={e => this.click()}>ぼたん</button>
            </div>
        );
    }

    static get props() {
        return {
            myProp1: skate.prop.number({
                attribute: true,
                default: (elem: any, data: any) => {
                    return 7;
                },
            }),
            // myProp2: { attribute: true },
        }
    }

    static render(elem: any) {
        return elem.render();
    }
});
