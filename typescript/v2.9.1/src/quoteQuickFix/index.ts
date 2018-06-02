import http from 'http';
import https from "https";

console.log('single quote');
console.log("double quote");


interface I {
    method(): string;
}

class C implements I {
    method(): string {
        throw new Error("Method not implemented.");
    }
}

hello

console.log(http, https);
