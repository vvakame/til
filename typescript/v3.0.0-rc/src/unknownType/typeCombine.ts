// unknownとのintersection typeは相手側に吸収される
// null
type U00 = unknown & null;
// undefined
type U01 = unknown & undefined;
// string
type U02 = unknown & string;
// any
type U03 = unknown & any;

// unknownとのunion typeはunknownになる
// 全部unknown
type U10 = unknown | null;
type U11 = unknown | undefined;
type U12 = unknown | string;
type U13 = unknown | any;

type T30<T> = unknown extends T ? true : false;  // Deferred
type T31<T> = T extends unknown ? true : false;  // Deferred (so it distributes)
type T32<T> = never extends T ? true : false;  // true
type T33<T> = T extends never ? true : false;  // Deferred
type T30D = T30<string>;
type T31D = T31<string>;
type T32D = T32<string>;
type T33D = T33<string>;
