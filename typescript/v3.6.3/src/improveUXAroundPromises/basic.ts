interface User {
    name: string;
    age: number;
    location: string;
}

declare function getUserData(): Promise<User>;
declare function displayUser(user: User): void;

async function f1() {
    // 普通のエラーと改善方法の提案が出る
    // error TS2345: Argument of type 'Promise<User>' is not assignable to parameter of type 'User'.
    //   Type 'Promise<User>' is missing the following properties from type 'User': name, age, location
    // 
    // `getUserData()` 部分に対して Did you forget to use 'await'?
    // displayUser(getUserData());

    // Quick fix を適用するとこうなる
    displayUser(await getUserData());
}

async function getCuteAnimals() {
    // error TS2339: Property 'json' does not exist on type 'Promise<Response>'.
    // 
    // `json` 部分に対して Did you forget to use 'await'?
    // fetch("https://reddit.com/r/aww.json").json();

    // Quick fix を適用するとこうなる
    (await fetch("https://reddit.com/r/aww.json")).json();
}
