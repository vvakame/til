const maxValue = 100;

// start から end の範囲を Extract to function in Global scope のrefactorを行うとする
/*start*/
for (let i = 0; i <= maxValue; i++) {
    // First get the squared value.
    let square = i ** 2;

    // Now print the squared value.
    console.log(square);
}
/*end*/

// 今まで 改行のみ の行は尊重されず、潰されてしまっていた
function newFunction1() {
    for (let i = 0; i <= maxValue; i++) {
        // First get the squared value.
        let square = i ** 2;
        // Now print the squared value.
        console.log(square);
    }
}

// 3.9 からちゃんと元の改行が保持されるようになった
function newFunction2() {
    for (let i = 0; i <= maxValue; i++) {
        // First get the squared value.
        let square = i ** 2;

        // Now print the squared value.
        console.log(square);
    }
}
