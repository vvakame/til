import { makeExecutableSchema } from "graphql-tools";

import {
    Event,
    CircleExhibitInfo,
    ProductInfo,
} from "./model";

import {
    Connection,
    createEventLoader,
    createEventQueryLoader,
    createCircleLoader,
    createCircleQueryLoader,
    createProductInfoLoader,
    createProductInfoQueryLoader,
    createProductContentLoader,
    createProductContentQueryLoader,
} from "./dataLoader";

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
    description: String
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

const eventLoader = createEventLoader();
const eventQueryLoader = createEventQueryLoader(eventLoader);

const circleLoader = createCircleLoader();
const circleQueryLoader = createCircleQueryLoader(circleLoader);

const productLoader = createProductInfoLoader();
const productQueryLoader = createProductInfoQueryLoader(productLoader);

const productContentLoader = createProductContentLoader();
const productContentQueryLoader = createProductContentQueryLoader(productContentLoader);

const resolvers = {
    Query: {
        events: (_root: any, _args: any, _context: any) => {
            return eventQueryLoader.load({ all: true });
        },
        event: (_root: any, args: any) => {
            return eventLoader.load(args.id);
        },
        searchCircle: (_root: any, args: any): Promise<Connection<CircleExhibitInfo> & { nodes: CircleExhibitInfo[]; }> => {
            return circleQueryLoader.load({
                eventID: args.eventID,
                cursor: args.after,
                limit: args.first,
            })
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
            const resp = await circleQueryLoader.load({
                all: true,
                eventID: event.id,
            });
            return resp.nodes;
        },
    },
    CircleExhibitInfo: {
        products: async (circle: CircleExhibitInfo) => {
            const resp = await productQueryLoader.load({
                all: true,
                circleExhibitInfoID: circle.id!,
            });
            return resp.nodes;
        }
    },
    ProductInfo: {
        contents: async (productInfo: ProductInfo) => {
            const resp = await productContentQueryLoader.load({
                all: true,
                productInfoID: productInfo.id,
            });
            return resp.nodes;
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
