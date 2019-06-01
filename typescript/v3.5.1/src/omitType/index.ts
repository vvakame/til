type Cat = {
    kind: string;
    name: string;
    weight: number;
};

// T1 = { kind: string; name: string; }
type T1 = Omit<Cat, "weight" | "eyeColoer">;

type OmitStrict<T, K extends keyof T> = Pick<T, Exclude<keyof T, K>>;

// error TS2344: Type '"weight" | "eyeColoer"' does not satisfy the constraint '"weight" | "kind" | "name"'.
//   Type '"eyeColoer"' is not assignable to type '"weight" | "kind" | "name"'.
// type T2 = OmitStrict<A, "weight" | "eyeColoer">;


export { }
