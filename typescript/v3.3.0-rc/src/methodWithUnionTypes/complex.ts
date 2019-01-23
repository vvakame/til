type Fruit = "apple" | "orange";
type Color = "red" | "orange";

interface FruitConstructor {
    new(arg: Fruit): { fruit: string };
}
interface ColorConstructor {
    new(arg: Color): { color: Color };
}

declare let Ctor: FruitConstructor | ColorConstructor;
let obj = new Ctor("orange");

if (isFruit(obj)) {
    console.log(obj.fruit);
}

function isFruit(obj: any): obj is { fruit: Fruit } {
    return !!obj.fruit;
}

export { }
