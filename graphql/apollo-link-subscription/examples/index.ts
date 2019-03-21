import { graphql, graphqlSync, print, introspectionQuery, IntrospectionQuery } from "graphql";
import { makeExecutableSchema, addMockFunctionsToSchema } from 'graphql-tools';
import { ApolloLink, Observable, Operation, NextLink, FetchResult } from "apollo-link";

import gql from "graphql-tag";
import { ApolloClient } from "apollo-client";
import { InMemoryCache, NormalizedCacheObject } from "apollo-cache-inmemory";
import { SubscriptionCovertLink, EntityInfo } from "../lib/";

import fs from "fs";
import path from "path";
import { CommentAdded } from "./graphql/CommentAdded";

const schemaString = fs.readFileSync(path.resolve(__dirname, "./schema.graphql"), { encoding: "utf8" });

const schema = makeExecutableSchema({ typeDefs: schemaString });
addMockFunctionsToSchema({ schema });
const introspectionResult = graphqlSync(schema, introspectionQuery).data! as any as IntrospectionQuery;

function delay(ms: number) {
  return new Promise(resolve => {
    setTimeout(() => {
      resolve();
    }, ms);
  });
}

export const link = new ApolloLink(operation => {
  return new Observable(observer => {
    const { query, operationName, variables } = operation;

    (async () => {
      try {
        await delay(800);
        const result = await graphql(schema, print(query), null, null, variables, operationName);
        observer.next(result);
        observer.complete();
      } catch (e) {
        observer.error(e);
      }
    })();
  });
});

function idReceiver(): Observable<string> {
  return new Observable<string>(observer => {
    let counter = 1;
    setInterval(() => {
      observer.next(`Comment:${counter}`);
      counter++;
    }, 1000);
  });
}

const client: ApolloClient<NormalizedCacheObject> = new ApolloClient({
  cache: new InMemoryCache(),
  link: ApolloLink.from([
    new SubscriptionCovertLink({
      schema: introspectionResult,
      observableConstructor: (originalOperation, alternativeQueries) => {
        return idReceiver().map((id): EntityInfo => {
          return {
            fieldName: Object.keys(alternativeQueries)[0],
            id,
            action: "insert",
          };
        });
      },
      executeQuery: ({ query, variables }) => {
        return client.query({
          query,
          variables,
        });
      }
    }),
    link,
  ]),
});

client
  .subscribe<{ data: CommentAdded; }>({
    query: gql`
      subscription CommentAdded {
        commentAdded {
          id
          text
        }
      }
    `,
  }).subscribe(v => {
    console.log("client.subscribe", v.data);
  });
