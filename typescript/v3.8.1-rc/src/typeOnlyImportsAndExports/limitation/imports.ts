// typeのみの場合、default importとnamed importを両方同時にすることはできない
//   ↓ は現時点ではエラーになる type import は A と B 両方？ A のみ？ パット見意味が一意にならない
// import type A, { B } from "./exports";

// 分割すればOK
import type A from "./exports";
import type { B } from "./exports";
