function tuple<TS extends any[]>(...args: TS): TS {
    return args;
}
// [ number, boolean, string | Date ] が得られる
let t = tuple(1, true, true ? "str" : new Date());
