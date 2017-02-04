fn main() {
    let str = test("Hello, world".to_string());
    println!("{}", str);
}

fn test(arg: String) -> String {
    return arg + "!";
}
