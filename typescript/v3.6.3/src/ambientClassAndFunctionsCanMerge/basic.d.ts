// こういう書き方ができるようになった
export declare function Point2D(x: number, y: number): Point2D;
export declare class Point2D {
    x: number;
    y: number;
    constructor(x: number, y: number);
}


// declare var Date: DateConstructor;
// interface DateConstructor {
//     new(): Date;
//     (): string;
// }

export declare function Date(): string;
export declare class Date {
    constructor();
}
