type Result<T> = { error?: undefined, value: T } | { error: Error };
// 昔もこんな感じだったらまぁイケた
// type Result<T> = { error: false, value: T } | { error: true };

function test(x: Result<number>) {
    if (!x.error) {
        // 今まではこのやり方でtype narrowingできなかった
        // error TS2339: Property 'value' does not exist on type 'Result<number>'.
        //   Property 'value' does not exist on type '{ error: Error; }'.
        // error が true | false だったりとunit typeである必要があった…
        // だがしかし今はこれでイケる！
        x.value;  // number
    }
    else {
        x.error.message;  // string
    }
}

test({ value: 10 });
test({ error: new Error("boom") });
