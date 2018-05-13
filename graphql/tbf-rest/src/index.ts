import express from "express";
import bodyParser from "body-parser";

import { graphqlExpress, graphiqlExpress } from "apollo-server-express";
import { ApolloEngine } from "apollo-engine";

import { schema } from "./schema";

const app = express();

app.post(
    "/graphql",
    bodyParser.json(),
    graphqlExpress({
        schema,
        tracing: true,
        cacheControl: true,
        context: {
            secrets: {
            },
        },
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
          id
          name
          courses {
            id
            name
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
    stores: [
        {
            name: "publicResponseCache",
            inMemory: {
                cacheSize: 10485760
            }
        }
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
