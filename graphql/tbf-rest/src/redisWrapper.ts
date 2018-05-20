import redis from "redis";

import { promisify } from "util";

const cli = redis.createClient();

let set = promisify(cli.set);
set = set.bind(cli);

let get = promisify(cli.get);
get = get.bind(cli);

type RedisMSet = (...kvs: string[]) => string;
let mset: RedisMSet = promisify(cli.mset) as any;
mset = mset.bind(cli);

type RedisMGet = <R>(keys: string[]) => (R | null)[];
let mget: RedisMGet = promisify(cli.mget) as any;
mget = mget.bind(cli);

export const redisCli = {
    set,
    get,
    mset,
    mget,
};
