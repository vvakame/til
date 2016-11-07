import * as skate from "skatejs";

const skatex = { createElement: skate.h };

export default skate.define("x-component2", class Component2 extends skate.Component {
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

    render(elem: this): JSX.Element {
        return (
            <div>
                <form>
                </form>
                <p><span>{this.myProp1}</span> & <span>{this.myProp2}</span></p>
                <button onClick={e => this.click()}>ぼたん</button>
            </div>
        );
    }

    static get props(): { [key: string]: skate.ComponentProp<any>; } {
        return {
            myProp1: skate.prop.number<Component2>({
                attribute: true,
                default(elem, data) {
                    return 7;
                },
            }),
            // myProp2: { attribute: true },
        }
    }

    static render(elem: any) {
        return elem.render(elem);
    }
});
