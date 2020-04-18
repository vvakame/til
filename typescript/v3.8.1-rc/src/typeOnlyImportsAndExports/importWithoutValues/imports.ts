// importsNotUsedAsValues が remove の場合、消える
// importsNotUsedAsValues が preserve の場合、 require("./exports"); などが残る
// importsNotUsedAsValues が error の場合、怒られる
//   error TS1371: This import is never used as a value and must use 'import type' because the 'importsNotUsedAsValues' is set to 'error'.
import { A } from "./exports";
// ↓ならOK
// import type { A } from "./exports";

function foo(): typeof A {
    return "a";
}
