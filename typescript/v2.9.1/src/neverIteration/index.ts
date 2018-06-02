const neverVar: never = (() => { throw new Error() })();

// error TS2407: The right-hand side of a 'for...in' statement must be of type 'any', an object type or a type parameter, but here has type 'never'.
for (let v in neverVar) {
    console.log(v);
}

// error TS2488: Type 'never' must have a '[Symbol.iterator]()' method that returns an iterator.
for (let v of neverVar) {
    console.log(v);
}
