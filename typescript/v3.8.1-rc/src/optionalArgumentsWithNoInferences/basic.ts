function foo(f: () => void) {
}

// error TS7006: Parameter 'param' implicitly has an 'any' type.
foo((param?) => {
});
