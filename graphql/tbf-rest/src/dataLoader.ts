import DataLoader from "dataloader";

import fetch from "node-fetch";

import { Event, ProductInfo, ProductInfoListResp, EventListResp, CircleExhibitInfo, CircleListResp, ProductContent, ProductContentListResp, ProductInfoBatchListReq, ProductInfoBatchListResp } from "./model";

const apiBaseUrl = "https://techbookfest.org"; // TODO 二重定義してる箇所がある

export interface Connection<T1 extends Node, T2 extends Edge<T1> = Edge<T1>> {
    pageInfo: PageInfo;
    edges: T2[];
}

export interface PageInfo {
    startCursor?: string;
    endCursor?: string;
    hasNextPage: boolean;
    hasPreviousPage: boolean;
}

export interface Edge<T extends Node> {
    cursor?: string;
    node: T;
}

export interface Node {
    id?: number | string; // NOTE 本当は ? ナシだけどコンパイル通すのめんどくさいので
}

function createCacheMap<K, V>(): DataLoader.CacheMap<K, V> {
    return new Map();
}

const emitLog = true;
let cnt = 0;

function log(message?: any, ...optionalParams: any[]) {
    cnt++;

    if (emitLog) {
        console.log(cnt, message, ...optionalParams);
    }
}

export function createEventLoader() {
    type Key = string;

    return new DataLoader<string, Event>(
        async keys => {
            log("event", keys.length);

            // TODO RedisとのBatchGet, BatchSet
            // TODO REST APIにBatchGet用のエンドポイントを生やす
            return Promise.all(keys.map(async id => {
                const resp = await fetch(`${apiBaseUrl}/api/event/${id}`);
                return await resp.json();
            }));
        },
        {
            maxBatchSize: 100,
            cache: true, // TODO 後でfalseにする Redisで管理するのでcacheMapは使わない キャッシュのライフタイム制御のため
            cacheKeyFn: (key: Key): string => JSON.stringify(key),
            cacheMap: createCacheMap(),
        },
    );
}

export function createEventQueryLoader(baseLoader: ReturnType<typeof createEventLoader>) {
    type Key = {
        all?: true,
        cursor?: string;
        limit?: number;
    };

    return new DataLoader<Key, Event[]>(
        async keys => {
            log("events", keys.length);

            // TODO RedisとのBatchGet, BatchSet
            // TODO REST APIにBatchGet用のエンドポイントを生やす
            return Promise.all(keys.map(async key => {
                if (key.all) {
                    let list: Event[] = [];
                    let cursor = "";
                    while (true) {
                        const json = await fetchEntity({
                            ...key,
                            ...{ cursor },
                        });
                        list = [...list, ...(json.list || [])];
                        if (json.cursor) {
                            cursor = json.cursor;
                        } else {
                            break;
                        }
                    }
                    makeEntityCache(list);
                    return list;
                }

                const json = await fetchEntity(key);
                makeEntityCache(json.list);
                return json.list || [];
            }));
        },
        {
            maxBatchSize: 100,
            cache: true, // TODO 後でfalseにする Redisで管理するのでcacheMapは使わない キャッシュのライフタイム制御のため
            cacheKeyFn: (key: Key): string => JSON.stringify(key),
            cacheMap: createCacheMap(),
        },
    );

    async function fetchEntity({ cursor, limit }: Key) {
        const resp = await fetch(`${apiBaseUrl}/api/event&cursor=${cursor || ""}&limit=${limit || 30}`);
        const json: EventListResp = await resp.json();
        return json;
    }
    function makeEntityCache(entities?: Event[]) {
        (entities || []).map(entity => {
            baseLoader.prime(entity.id!, entity);
        });
    }
}


export function createCircleLoader() {
    type Key = string;
    type Entity = CircleExhibitInfo;

    return new DataLoader<string, Entity>(
        async ids => {
            log("circle", ids.length);

            // TODO RedisとのBatchGet, BatchSet
            // TODO REST APIにBatchGet用のエンドポイントを生やす
            return Promise.all(ids.map(async id => {
                const resp = await fetch(`${apiBaseUrl}/api/circle/${id}`);
                const entity: Entity = await resp.json();
                return entity;
            }));
        },
        {
            maxBatchSize: 100,
            cache: true, // TODO 後でfalseにする Redisで管理するのでcacheMapは使わない キャッシュのライフタイム制御のため
            cacheKeyFn: (key: Key): string => JSON.stringify(key),
            cacheMap: createCacheMap(),
        },
    );
}

export function createCircleQueryLoader(baseLoader: ReturnType<typeof createCircleLoader>) {
    type Key = {
        all?: true,
        eventID?: string;
        cursor?: string;
        limit?: number;
    };
    type Entity = CircleExhibitInfo;
    type ListResp = CircleListResp;
    type PagedResp = Connection<Entity> & { nodes: Entity[]; };

    return new DataLoader<Key, PagedResp>(
        async keys => {
            log("circles", keys.length);

            // TODO RedisとのBatchGet, BatchSet
            // TODO REST APIにBatchGet用のエンドポイントを生やす
            return Promise.all(keys.map(async key => {
                if (key.all) {
                    let list: Entity[] = [];
                    let cursor = "";
                    while (true) {
                        const json = await fetchEntity({
                            ...key,
                            ...{ cursor },
                        });
                        list = [...list, ...(json.list || [])];
                        if (json.cursor) {
                            cursor = json.cursor;
                        } else {
                            break;
                        }
                    }
                    return makePagedResp(key, { list });
                }

                const json = await fetchEntity(key);
                makeEntityCache(json.list);
                return makePagedResp(key, json);
            }));
        },
        {
            maxBatchSize: 100,
            cache: true, // TODO 後でfalseにする Redisで管理するのでcacheMapは使わない キャッシュのライフタイム制御のため
            cacheKeyFn: (key: Key): string => JSON.stringify(key),
            cacheMap: createCacheMap(),
        },
    );

    async function fetchEntity({ eventID, cursor, limit }: Key) {
        const resp = await fetch(`${apiBaseUrl}/api/circle?eventID=${eventID}&cursor=${cursor || ""}&limit=${limit || 100}`);
        const json: ListResp = await resp.json();
        return json;
    }
    function makeEntityCache(entities?: Entity[]) {
        (entities || []).map(entity => baseLoader.prime(entity.id!, entity));
    }
    function makePagedResp(key: Key, json: ListResp) {
        const result: PagedResp = {
            pageInfo: {
                endCursor: json.cursor,
                hasNextPage: !!json.cursor,
                hasPreviousPage: !!key.cursor,
            },
            nodes: (json.list || []),
            edges: (json.list || []).map(node => ({ node })),
        };
        if (result.edges.length !== 0) {
            result.edges[result.edges.length - 1].cursor = json.cursor;
        }
        return result;
    }
}

export function createProductInfoLoader() {
    type Key = string;

    return new DataLoader<string, ProductInfo>(
        async ids => {
            log("productInfo", ids.length);

            // TODO RedisとのBatchGet, BatchSet
            // TODO REST APIにBatchGet用のエンドポイントを生やす
            return Promise.all(ids.map(async id => {
                const resp = await fetch(`${apiBaseUrl}/api/product/${id}`);
                const entity: ProductInfo = await resp.json();
                return entity;
            }));
        },
        {
            maxBatchSize: 100,
            cache: true, // TODO 後でfalseにする Redisで管理するのでcacheMapは使わない キャッシュのライフタイム制御のため
            cacheKeyFn: (key: Key): string => JSON.stringify(key),
            cacheMap: createCacheMap(),
        },
    );
}

export function createProductInfoQueryLoader(baseLoader: ReturnType<typeof createProductInfoLoader>) {
    type Key = {
        all?: true,
        circleExhibitInfoID?: string;
        cursor?: string;
        limit?: number;
    };
    type Entity = ProductInfo;
    type ListResp = ProductInfoListResp;
    type PagedResp = Connection<Entity> & { nodes: Entity[]; };

    return new DataLoader<Key, PagedResp>(
        async keys => {
            log("productInfos", keys.length);

            const req: ProductInfoBatchListReq = {
                requestList: keys
                    .map(key => ({ circleExhibitInfoID: key.circleExhibitInfoID! })),
                visibility: "site",
            };
            const resp = await fetch(`${apiBaseUrl}/api/product/batch-list`, {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify(req),
            });
            const json: ProductInfoBatchListResp = await resp.json();

            return (json.list || []).map((list, idx): PagedResp => {
                makeEntityCache(list);
                return makePagedResp(keys[idx], { list });
            });
        },
        {
            maxBatchSize: 100,
            cache: true, // TODO 後でfalseにする Redisで管理するのでcacheMapは使わない キャッシュのライフタイム制御のため
            cacheKeyFn: (key: Key): string => JSON.stringify(key),
            cacheMap: createCacheMap(),
        },

    );

    function makeEntityCache(entities?: Entity[]) {
        (entities || []).map(entity => baseLoader.prime(entity.id!, entity));
    }
    function makePagedResp(key: Key, json: ListResp) {
        const result: PagedResp = {
            pageInfo: {
                endCursor: json.cursor,
                hasNextPage: !!json.cursor,
                hasPreviousPage: !!key.cursor,
            },
            nodes: (json.list || []),
            edges: (json.list || []).map(node => ({ node })),
        };
        if (result.edges.length !== 0) {
            result.edges[result.edges.length - 1].cursor = json.cursor;
        }
        return result;
    }
}

export function createProductContentLoader() {
    type Key = string;
    type Entity = ProductContent;

    return new DataLoader<string, Entity>(
        async ids => {
            log("productContent", ids.length);

            // TODO RedisとのBatchGet, BatchSet
            // TODO REST APIにBatchGet用のエンドポイントを生やす
            return Promise.all(ids.map(async id => {
                const resp = await fetch(`${apiBaseUrl}/api/product/${id}`);
                const entity: Entity = await resp.json();
                return entity;
            }));
        },
        {
            maxBatchSize: 100,
            cache: true, // TODO 後でfalseにする Redisで管理するのでcacheMapは使わない キャッシュのライフタイム制御のため
            cacheKeyFn: (key: Key): string => JSON.stringify(key),
            cacheMap: createCacheMap(),
        },
    );
}

export function createProductContentQueryLoader(baseLoader: ReturnType<typeof createProductContentLoader>) {
    type Key = {
        all?: true,
        productInfoID?: string;
        cursor?: string;
        limit?: number;
    };
    type Entity = ProductContent;
    type ListResp = ProductContentListResp;
    type PagedResp = Connection<Entity> & { nodes: Entity[]; };

    return new DataLoader<Key, PagedResp>(
        async keys => {
            log("productContents", keys.length);

            // TODO RedisとのBatchGet, BatchSet
            // TODO REST APIにBatchGet用のエンドポイントを生やす
            return Promise.all(keys.map(async key => {
                if (key.all) {
                    let list: Entity[] = [];
                    let cursor = "";
                    while (true) {
                        const json = await fetchEntity({
                            ...key,
                            ...{ cursor },
                        });
                        list = [...list, ...(json.list || [])];
                        if (json.cursor) {
                            cursor = json.cursor;
                        } else {
                            break;
                        }
                    }
                    return makePagedResp(key, { list });
                }

                const json = await fetchEntity(key);
                makeEntityCache(json.list);
                return makePagedResp(key, json);
            }));
        },
        {
            maxBatchSize: 100,
            cache: true, // TODO 後でfalseにする Redisで管理するのでcacheMapは使わない キャッシュのライフタイム制御のため
            cacheKeyFn: (key: Key): string => JSON.stringify(key),
            cacheMap: createCacheMap(),
        },
    );

    async function fetchEntity({ productInfoID, cursor, limit }: Key) {
        const resp = await fetch(`${apiBaseUrl}/api/productcontent?productInfoID=${productInfoID}&cursor=${cursor}&limit=${limit || 10}`);
        const json: ListResp = await resp.json();
        return json;
    }
    function makeEntityCache(entities?: Entity[]) {
        (entities || []).map(entity => baseLoader.prime(entity.id!, entity));
    }
    function makePagedResp(key: Key, json: ListResp) {
        const result: PagedResp = {
            pageInfo: {
                endCursor: json.cursor,
                hasNextPage: !!json.cursor,
                hasPreviousPage: !!key.cursor,
            },
            nodes: (json.list || []),
            edges: (json.list || []).map(node => ({ node })),
        };
        if (result.edges.length !== 0) {
            result.edges[result.edges.length - 1].cursor = json.cursor;
        }
        return result;
    }
}
