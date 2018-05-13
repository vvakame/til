import { makeExecutableSchema } from "graphql-tools";

import fetch from "node-fetch";

const gql = String.raw;

// Construct a schema, using GraphQL schema language
const typeDefs = gql`
  type Query {
    events: [Event]
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

  type CircleExhibitInfo @cacheControl(maxAge: 60) {
    id: ID
    name: String
    place: String
    exhibitFee: Int
  }
`;

const apiBaseUrl = "https://techbookfest.org";

const resolvers = {
    Query: {
        events: async (_root: any, _args: any, _context: any) => {
            const resp = await fetch(`${apiBaseUrl}/api/event`);
            const json = await resp.json();
            return json.list;
        }
    },
    Event: {
        courses: (event: any) => {
            return event.eventExhibitCourses;
        },
        circles: async (event: any) => {
            const resp = await fetch(`${apiBaseUrl}/api/circle?eventID=${event.id}`);
            const json = await resp.json();
            return json.list;
        },
    },
};

// Required: Export the GraphQL.js schema object as "schema"
const schema = makeExecutableSchema({
    typeDefs,
    resolvers
});

export { schema };
