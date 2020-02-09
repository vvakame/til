// importsNotUsedAsValues が preserve でも exports.d.ts なら当然消える
import type { A } from "./exports";

const a: typeof A = "foo";
