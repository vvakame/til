import { makeExecutableSchema } from "graphql-tools";

import { fetch } from "./fetch";

import { Event, EventListResp, CircleExhibitInfo, ProductInfoListResp, ProductInfo, ProductContentListResp, CircleListResp, ProductContent } from "./model";

const gql = String.raw;

// Construct a schema, using GraphQL schema language
const typeDefs = gql`
  type Query {
    events: [Event]
    event(id: String!, visibility: Visibillity = Site): Event
    searchCircle(
      eventID: String!
      first: Int!
      after: String
    ): CircleExhibitInfoSearchResultItemConnection!
  }

  # type Mutation {
  # }

  enum Visibillity {
    SITE
    STAFF
  }

  interface Connection {
    pageInfo: PageInfo
    edges: [Edge]
  }
  
  interface Edge {
    cursor: String
    node: Node!
  }
  
  interface Node {
    id: ID!
  }
  
  type PageInfo {
    startCursor: String
    endCursor: String
    hasNextPage: Boolean!
    hasPreviousPage: Boolean!
  }

  type CircleExhibitInfoSearchResultItemConnection implements Connection {
    pageInfo: PageInfo!
    edges: [CircleExhibitInfoSearchResultItemEdge]
    nodes: [CircleExhibitInfo]
  }

  type CircleExhibitInfoSearchResultItemEdge implements Edge {
    cursor: String
    node: CircleExhibitInfo!
  }

  type Event @cacheControl(maxAge: 60) {
    id: ID
    name: String
    place: String
    courses: [EventExhibitCourse]
    circles: [CircleExhibitInfo]
  }

  type EventExhibitCourse @cacheControl(maxAge: 60) {
    id: ID
    name: String
    place: String
    exhibitFee: Int
  }

  type CircleExhibitInfo implements Node @cacheControl(maxAge: 60) {
    id: ID!
    name: String
    place: String
    exhibitFee: Int
    products: [ProductInfo]
  }

  type ProductInfo @cacheControl(maxAge: 60) {
    id: ID
    name: String
    firstAppearanceEventName: String
    firstAtTechBookFest: Boolean
    price: Int
    relatedURLs: [String]
    type: String # TODO enum
    images: [Image]
    contents: [ProductContent]
  }

  type ProductContent @cacheControl(maxAge: 60) {
    id: ID
    contentType: String
    fileName: String
    fileSize: Int
  }

  type Image @cacheControl(maxAge: 60) {
    id: ID
    fileSize: Int
    url: String
    height: Int
    Width: Int
  }
`;

const apiBaseUrl = "https://techbookfest.org";

interface Connection<T1 extends Node, T2 extends Edge<T1> = Edge<T1>> {
    pageInfo: PageInfo;
    edges: T2[];
}

interface PageInfo {
    startCursor?: string;
    endCursor?: string;
    hasNextPage: boolean;
    hasPreviousPage: boolean;
}

interface Edge<T extends Node> {
    cursor?: string;
    node: T;
}

interface Node {
    id?: number | string; // NOTE 本当は ? ナシだけどコンパイル通すのめんどくさいので
}

const resolvers = {
    Query: {
        events: async (_root: any, _args: any, _context: any) => {
            let list: Event[] = [];
            let cursor: string | undefined;
            while (true) {
                const resp = await fetch(`${apiBaseUrl}/api/event&cursor=${cursor || ""}`);
                const json: EventListResp = await resp.json();
                list = [...list, ...(json.list || [])];
                if (json.cursor) {
                    cursor = json.cursor;
                } else {
                    break;
                }
            }
            return list;
        },
        event: async (_root: any, args: any) => {
            const resp = await fetch(`${apiBaseUrl}/api/event/${args.id}`);
            const json: Event = await resp.json();
            return json;
        },
        searchCircle: async (_root: any, args: any): Promise<Connection<CircleExhibitInfo> & { nodes: CircleExhibitInfo[]; }> => {
            const resp = await fetch(`${apiBaseUrl}/api/circle?eventID=${args.eventID}&cursor=${args.after || ""}&limit=${args.first || 100}`);
            const json: CircleListResp = await resp.json();
            const result: Connection<CircleExhibitInfo> & { nodes: CircleExhibitInfo[]; } = {
                pageInfo: {
                    endCursor: json.cursor,
                    hasNextPage: !!json.cursor,
                    hasPreviousPage: !!args.after,
                },
                nodes: (json.list || []),
                edges: (json.list || []).map(node => ({ node })),
            };
            if (result.edges.length !== 0) {
                result.edges[result.edges.length - 1].cursor = json.cursor;
            }
            return result;
        },
    },
    // Mutation: {
    // },
    Visibillity: {
        SITE: "site",
        STAFF: "staff",
    },
    Event: {
        courses: (event: Event) => {
            return event.eventExhibitCourses;
        },
        circles: async (event: Event) => {
            let list: CircleExhibitInfo[] = [];
            let cursor: string | undefined;
            while (true) {
                const resp = await fetch(`${apiBaseUrl}/api/circle?eventID=${event.id}&cursor=${cursor || ""}`);
                const json: CircleListResp = await resp.json();
                list = [...list, ...(json.list || [])];
                if (json.cursor) {
                    cursor = json.cursor;
                } else {
                    break;
                }
            }
            return list;
        },
    },
    CircleExhibitInfo: {
        products: async (circle: CircleExhibitInfo) => {
            let list: ProductInfo[] = [];
            let cursor: string | undefined;
            while (true) {
                const resp = await fetch(`${apiBaseUrl}/api/product?circleExhibitInfoID=${circle.id}&cursor=${cursor || ""}`);
                const json: ProductInfoListResp = await resp.json();
                list = [...list, ...(json.list || [])];
                if (json.cursor) {
                    cursor = json.cursor;
                } else {
                    break;
                }
            }
            return list;
        }
    },
    ProductInfo: {
        contents: async (productInfo: ProductInfo) => {
            let list: ProductContent[] = [];
            let cursor: string | undefined;
            while (true) {
                const resp = await fetch(`${apiBaseUrl}/api/productcontent?productInfoID=${productInfo.id}&cursor=${cursor || ""}`);
                const json: ProductContentListResp = await resp.json();
                list = [...list, ...(json.list || [])];
                if (json.cursor) {
                    cursor = json.cursor;
                } else {
                    break;
                }
            }
            return list;
        },
    },
};

// Required: Export the GraphQL.js schema object as "schema"
const schema = makeExecutableSchema({
    typeDefs,
    resolvers,
    logger: { log: e => console.log(e) },
});

export { schema };
