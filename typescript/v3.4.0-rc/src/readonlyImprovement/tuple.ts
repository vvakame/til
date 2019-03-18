let readonlyTuple1: readonly [string, string] = ["a", "b"];
// ↓の書き方は前からできたけど readonlyTuple2[0] = "b" とかから守ってくれなかった（後述）
let readonlyTuple2: Readonly<[string, string]> = ["a", "b"];

// 両方NG!
// error TS2540: Cannot assign to '0' because it is a read-only property.
// readonlyTuple1[0] = "a";
// readonlyTuple2[0] = "b";

export { }
