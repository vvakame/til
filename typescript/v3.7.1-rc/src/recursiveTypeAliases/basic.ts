// TypeScript v3.6 までは直接自分自身を参照するような再帰構造は書けなかった
// error TS2456: Type alias 'Json' circularly references itself.
// TypeScript v3.7 以降は大丈夫
type Json =
    | string
    | number
    | boolean
    | null
    | { [property: string]: Json }
    | Json[];

let obj1: Json = 1;
let obj2: Json = "string";
let obj3: Json = {};
let obj4: Json = [];
let obj5: Json = {
    foo: [],
    bar: true,
};

{ // TypeScript v3.6 までは補助となるinterfaceとかの定義が必要だった
    type Json =
        | string
        | number
        | boolean
        | null
        | JsonObject
        | JsonArray;
    type JsonObject = {
        [property: string]: Json;
    };
    interface JsonArray extends Array<Json> { }
}
