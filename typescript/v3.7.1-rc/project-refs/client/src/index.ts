import { hello, addSuffix } from "../../shared/src/";

const str1 = hello("TypeScript");
const helloWithExclamention = addSuffix(hello, "!");
const str2 = helloWithExclamention("TypeScript");

console.log(str1);
console.log(str2);
