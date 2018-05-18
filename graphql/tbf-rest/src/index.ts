import express from "express";
import bodyParser from "body-parser";
import cookieParser from "cookie-parser";

import { graphqlExpress, graphiqlExpress } from "apollo-server-express";
import { ApolloEngine } from "apollo-engine";

import { schema } from "./schema";

const app = express();
app.use(cookieParser());

app.post(
    "/graphql",
    bodyParser.json(),
    graphqlExpress(req => {
        return {
            schema,
            tracing: true,
            cacheControl: true,
            context: {
                secrets: {
                },
            },
            rootValue: {
                authCookie: req!.cookies["user"],
            },
        };
    }),
);

const gql = String.raw;

app.get(
    "/graphiql",
    graphiqlExpress({
        endpointURL: "/graphql",
        query: gql`
        query UpcomingEvents {
            events {
              name
              circles {
                name
                products {
                  name
                  images {
                    url
                  }
                  contents {
                    fileName
                  }
                }
              }
            }
            searchCircle(eventID: "tbf04", first: 3, after: "") {
              pageInfo {
                startCursor
                endCursor
                hasNextPage
                hasPreviousPage
              }
              edges {
                cursor
                node {
                  id
                  name
                  products {
                    name
                  }
                }
              }
            }
          }
        `
    })
);

app.use(express.static("public"));

const PORT = process.env.PORT || 3000;

const engine = new ApolloEngine({
    apiKey: process.env.ENGINE_API_KEY,
    logging: {
        level: "WARN",
    },
    sessionAuth: {
        cookie: "user",
    },
    stores: [
        {
            name: "publicResponseCache",
            inMemory: {
                cacheSize: 10485760,
            },
        },
    ],
    queryCache: {
        publicFullQueryStore: "publicResponseCache"
    }
});

// Start the app
engine.listen(
    {
        port: PORT,
        expressApp: app
    },
    () => {
        console.log(`Go to http://localhost:${PORT}/graphiql to run queries!`);
    }
);
