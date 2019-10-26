interface User {
    isAdministrator(): boolean;
    notify(): void;
    doNotDisturb?(): boolean;
}

function sudo() {
    console.log("exec sudu!");
}

// function doAdminThingA(user: User) {
//     // エラーになる！それ絶対存在するプロパティだから常にtrueなんだけど、ホントは呼び出したかったんじゃないの？
//     // error TS2774: This condition will always return true since the function is always defined. Did you mean to call it instead?
//     if (user.isAdministrator) {
//         sudo();
//     } else {
//         throw new Error("User is not an admin");
//     }
// }

function doAdminThingB(user: User) {
    // 当然、呼び出している場合はエラーにならない
    if (user.isAdministrator()) {
        sudo();
    } else {
        throw new Error("User is not an admin");
    }
}

function doAdminThingC(user: User) {
    // わざとだよ！という場合は !! として真偽値に変換することで意図を伝えることができる
    if (!!user.isAdministrator) {
        sudo();
    } else {
        throw new Error("User is not an admin");
    }
}

function doAdminThingD(user: User) {
    if (user.notify) {
        // その後、呼び出すならOK
        user.notify();
    }
    if (user.doNotDisturb) {
        // doNotDisturb は optional なのでOK
        sudo(); // 現実的にはOKじゃないかもね！
    }
}

export { }
