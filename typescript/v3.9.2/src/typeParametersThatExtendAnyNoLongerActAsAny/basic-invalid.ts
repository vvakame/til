function foo<T extends any>(arg: T) {
    // v3.8 までは T は any として扱われてたのでエラーなし
    // v3.9 からは T は extends unknown と同等に振る舞いっぽ
    // error TS2339: Property 'notExists' does not exist on type 'T'.
    arg.notExists;
}

export {}
