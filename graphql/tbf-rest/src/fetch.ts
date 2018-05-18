import * as nodeFetch from "node-fetch";

import redis from "redis";

import { promisify } from "util";

const redisCli = redis.createClient();

let set = promisify(redisCli.set);
set = set.bind(redisCli);
let get = promisify(redisCli.get);
get = get.bind(redisCli);

export async function fetch(url: string | nodeFetch.Request, init?: nodeFetch.RequestInit): Promise<nodeFetch.Response> {
    const cacheKey = "graphql-" + url as string;

    const cached = await get(cacheKey); // TODO
    if (cached) {
        console.log(`from cache: ${cacheKey}`);
        const obj = JSON.parse(cached);
        return new nodeFetch.Response(obj.body, obj.response);
    }

    console.log(`from network: ${cacheKey}`);

    const resp = await nodeFetch.default(url, init);
    const bkResp = resp.clone();

    const text = await bkResp.text();

    const body: nodeFetch.BodyInit = text;
    const response: nodeFetch.ResponseInit = {
        status: bkResp.status,
        headers: bkResp.headers,
    };

    await set(cacheKey, JSON.stringify({
        body,
        response,
    }))

    return resp;
}
