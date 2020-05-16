let numberArray = [1, 9, 2, 8, 3, 7, 4, 6, 5];

// sort(compareFn?: (a: T, b: T) => number): this; という定義
// 次の書き方だと sort に渡した関数が値を返していない(のでコンパイルエラーになる)
// numberArray.sort((a, b) => { a - b });

// Quick Fix が追加されているので試してみよう！

// Add a return statement を選ぶと次のように変形される
numberArray.sort((a, b) => { return a - b; });

// Remove block body braces を選ぶと次のように変形される
numberArray.sort((a, b) => a - b);
