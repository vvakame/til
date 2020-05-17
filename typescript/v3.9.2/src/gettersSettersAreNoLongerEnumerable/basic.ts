class A {
    get a() {
        return "a";
    }
    set a(v: string) {
        console.log(v);
    }
}

export {}

// v3.8 での出力
// var A = /** @class */ (function () {
//     function A() {
//     }
//     Object.defineProperty(A.prototype, "a", {
//         get: function () {
//             return "a";
//         },
//         set: function (v) {
//             console.log(v);
//         },
//         enumerable: true,
//         configurable: true
//     });
//     return A;
// }());

// v3.9 での出力
// var A = /** @class */ (function () {
//     function A() {
//     }
//     Object.defineProperty(A.prototype, "a", {
//         get: function () {
//             return "a";
//         },
//         set: function (v) {
//             console.log(v);
//         },
//         enumerable: false,
//         configurable: true
//     });
//     return A;
// }());
