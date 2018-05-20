import DataLoader from "dataloader";

import fetch from "node-fetch";

import { Event, ProductInfo, ProductInfoListResp, EventListResp, CircleExhibitInfo, CircleListResp, ProductContent, ProductContentListResp, ProductInfoBatchListReq, ProductInfoBatchListResp } from "./model";
import { redisCli } from "./redisWrapper";

const apiBaseUrl = "https://techbookfest.org";

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

const emitLog = true;
let cnt = 0;

function log(message?: any, ...optionalParams: any[]) {
    cnt++;

    if (emitLog) {
        console.log(cnt, message, ...optionalParams);
    }
}

type CustomLoaderOptions<K, V> = {
    cacheKeyPrefix: string;
    loaderOptions?: DataLoader.Options<K, V>;
    fetch: (entityInfos: { idx: number; key: K, entity?: V; }[]) => Promise<{ idx: number; key: K, entity: V; }[]>;
};

class CustomLoader<Key, Entity> implements DataLoader<Key, Entity> {
    private cacheKeyPrefix: string;
    private loader: DataLoader<Key, Entity>
    private fetch: (entityInfos: { idx: number; key: Key, entity?: Entity; }[]) => Promise<{ idx: number; key: Key, entity: Entity; }[]>;

    constructor(opts: CustomLoaderOptions<Key, Entity>) {
        this.cacheKeyPrefix = opts.cacheKeyPrefix;
        this.loader = new DataLoader<Key, Entity>(
            this.batchLoad.bind(this),
            opts.loaderOptions,
        );
        this.fetch = opts.fetch;
    }

    cacheKey(key: Key): string {
        return `${this.cacheKeyPrefix}-${JSON.stringify(key)}`;
    }

    async redisBatchGet(...keys: Key[]): Promise<(Entity | undefined)[]> {
        const entities = await redisCli.mget<string>(keys.map(key => this.cacheKey(key)));
        return entities.map(entity => entity ? JSON.parse(entity) : void 0);
    }

    async redisBatchSet(...kvs: { key: Key; entity: Entity; }[]): Promise<string> {
        const redisCacheEntities = kvs
            .map(ks => {
                return [
                    this.cacheKey(ks.key),
                    JSON.stringify(ks.entity),
                ];
            })
            .reduce((p, c) => p.concat(c), []);
        return await redisCli.mset(...redisCacheEntities);
    }

    async batchLoad(keys: Key[]): Promise<Array<Entity | Error>> {
        log(`${this.cacheKeyPrefix}`, keys.length);

        const result: Entity[] = [];

        const entities = await this.redisBatchGet(...keys);
        log(`${this.cacheKeyPrefix} on cache`, this.cacheKeyPrefix, entities);

        const missing = entities
            .map((entity, idx) => {
                if (entity) {
                    result[idx] = entity;
                }
                return {
                    idx,
                    key: keys[idx],
                    entity,
                };
            })
            .filter(entityInfo => !entityInfo.entity);
        log(`${this.cacheKeyPrefix} missing on cache`, missing);

        if (missing.length !== 0) {
            const foundEntities = await this.fetch(missing);
            foundEntities.forEach(({ idx, entity }) => {
                result[idx] = entity;
            });
            log(`${this.cacheKeyPrefix} found from network`, foundEntities.map(({ key }) => key));

            const msetResult = await this.redisBatchSet(...foundEntities);
            log(`${this.cacheKeyPrefix} mset`, msetResult);
        }

        return result;
    }

    load(key: Key): Promise<Entity> {
        return this.loader.load(key);
    }
    loadMany(keys: Key[]): Promise<Entity[]> {
        return this.loader.loadMany(keys);
    }
    clear(key: Key) {
        // TODO Redisのキャッシュ
        this.loader.clear(key);
        return this;
    }
    clearAll() {
        // TODO Redisのキャッシュ
        this.loader.clearAll();
        return this;
    }
    prime(key: Key, value: Entity) {
        this.redisBatchSet({ key, entity: value });
        this.loader.prime(key, value);
        return this;
    }
}

export function createEventLoader() {
    type Key = string;
    type Entity = Event;

    return new CustomLoader<Key, Entity>({
        cacheKeyPrefix: "Event-Single",
        loaderOptions: {
            maxBatchSize: 100,
            cache: false,
        },
        fetch: async entityInfos => {
            return await Promise.all(entityInfos.map(async entityInfo => {
                const resp = await fetch(`${apiBaseUrl}/api/event/${entityInfo.key}`);
                const json: Entity = await resp.json();
                return {
                    ...entityInfo,
                    entity: json,
                };
            }));
        },
    });
}

export function createEventQueryLoader(baseLoader: ReturnType<typeof createEventLoader>) {
    type Key = {
        all?: true,
        cursor?: string;
        limit?: number;
    };
    type Entity = EventListResp;

    return new CustomLoader<Key, Entity>({
        cacheKeyPrefix: "Event-Query",
        loaderOptions: {
            maxBatchSize: 100,
            cache: false,
        },
        fetch: async entityInfos => {
            return await Promise.all(entityInfos.map(async entityInfo => {
                const { cursor, limit } = entityInfo.key;
                const resp = await fetch(`${apiBaseUrl}/api/event&cursor=${cursor || ""}&limit=${limit || 100}`);
                const json: EventListResp = await resp.json();
                (json.list || []).forEach(entity => {
                    baseLoader.prime(entity.id!, entity);
                });
                return {
                    ...entityInfo,
                    entity: json,
                };
            }));
        },
    });
}

export function createCircleLoader() {
    type Key = string;
    type Entity = CircleExhibitInfo;

    return new CustomLoader<Key, Entity>({
        cacheKeyPrefix: "CircleExhibitInfo-Single",
        loaderOptions: {
            maxBatchSize: 100,
            cache: false,
        },
        fetch: async entityInfos => {
            return await Promise.all(entityInfos.map(async entityInfo => {
                const resp = await fetch(`${apiBaseUrl}/api/circle/${entityInfo.key}`);
                const json: Entity = await resp.json();
                return {
                    ...entityInfo,
                    entity: json,
                };
            }));
        },
    });
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

    return new CustomLoader<Key, PagedResp>({
        cacheKeyPrefix: "CircleExhibitInfo-Query",
        loaderOptions: {
            maxBatchSize: 100,
            cache: false,
        },
        fetch: async entityInfos => {
            return await Promise.all(entityInfos.map(async entityInfo => {
                const key = entityInfo.key;

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
                    makeEntityCache(list);
                    return {
                        ...entityInfo,
                        entity: makePagedResp(key, { list }),
                    };
                }

                const json = await fetchEntity(key);
                makeEntityCache(json.list);
                return {
                    ...entityInfo,
                    entity: makePagedResp(key, json),
                };
            }));
        },
    });

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
    type Entity = ProductInfo;

    return new CustomLoader<Key, Entity>({
        cacheKeyPrefix: "ProductInfo-Single",
        loaderOptions: {
            maxBatchSize: 100,
            cache: false,
        },
        fetch: async entityInfos => {
            return await Promise.all(entityInfos.map(async entityInfo => {
                const resp = await fetch(`${apiBaseUrl}/api/product/${entityInfo.key}`);
                const json: Entity = await resp.json();
                return {
                    ...entityInfo,
                    entity: json,
                };
            }));
        },
    });
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

    return new CustomLoader<Key, PagedResp>({
        cacheKeyPrefix: "ProductInfo-Query",
        loaderOptions: {
            maxBatchSize: 100,
            cache: false,
        },
        fetch: async entityInfos => {
            const keys = entityInfos.map(entityInfo => entityInfo.key);

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

            return (json.list || []).map((list, idx) => {
                makeEntityCache(list);
                return {
                    ...entityInfos[idx],
                    entity: makePagedResp(keys[idx], { list }),
                };
            });
        },
    });

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

    return new CustomLoader<Key, Entity>({
        cacheKeyPrefix: "productContent-Single",
        loaderOptions: {
            maxBatchSize: 100,
            cache: false,
        },
        fetch: async entityInfos => {
            return await Promise.all(entityInfos.map(async entityInfo => {
                const resp = await fetch(`${apiBaseUrl}/api/product/${entityInfo.key}`);
                const json: Entity = await resp.json();
                return {
                    ...entityInfo,
                    entity: json,
                };
            }));
        },
    });
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

    return new CustomLoader<Key, PagedResp>({
        cacheKeyPrefix: "productContent-Query",
        loaderOptions: {
            maxBatchSize: 100,
            cache: false,
        },
        fetch: async entityInfos => {
            return await Promise.all(entityInfos.map(async entityInfo => {
                const key = entityInfo.key;

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
                    makeEntityCache(list);
                    return {
                        ...entityInfo,
                        entity: makePagedResp(key, { list }),
                    };
                }

                const json = await fetchEntity(key);
                makeEntityCache(json.list);
                return {
                    ...entityInfo,
                    entity: makePagedResp(key, json),
                };
            }));
        },
    });

    async function fetchEntity({ productInfoID, cursor, limit }: Key) {
        const resp = await fetch(`${apiBaseUrl}/api/productcontent?productInfoID=${productInfoID}&cursor=${cursor}&limit=${limit || 100}`);
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
