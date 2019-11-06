interface Animal { animalStuff: any }
interface Dog extends Animal { dogStuff: any }

class AnimalHouse {
    resident?: Animal;
    constructor(animal: Animal) {
        this.resident = animal;
    }
}

class DogHouse extends AnimalHouse {
    // useDefineForClassFields を使っている場合、値が undefined になる！
    // super 呼ぶ → AnimalHouse で resident がセットされる → DogHouse のfield initialize が走る → undefined に再設定される！
    // というのを防ぐため、コンパイルエラーになる
    // error TS2612: Property 'resident' will overwrite the base property in 'AnimalHouse'. If this is intentional, add an initializer. Otherwise, add a 'declare' modifier or remove the redundant declaration.
    // resident?: Dog;

    // 解消方法
    // declare を追加するとコード生成に関与しなくなる
    declare resident?: Dog;

    constructor(dog: Dog) {
        // super に Animal な値を渡すこともできてしまう
        // resident?: Dog が正しくなるかどうかはプログラマ次第…
        super(dog);

        // declare 使うのやめて、素直に自分で初期化したほうがいいかもね
        // this.resident = dog;
    }
}

let obj = new DogHouse({animalStuff: 1, dogStuff: 2});
// undefined と表示される useDefineForClassFields を有効にしている場合
console.log(obj.resident);

export {}
