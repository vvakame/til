import { graphql, print } from "graphql";
import { makeExecutableSchema, addMockFunctionsToSchema } from 'graphql-tools';
import { ApolloLink, Observable, Operation, NextLink, FetchResult } from "apollo-link";

import gql from "graphql-tag";
import { ApolloClient } from "apollo-client";
import { InMemoryCache } from "apollo-cache-inmemory";

const schemaString = `
type Query {
  comment(id: ID): Comment
}

type Subscription {
  commentAdded: Comment
}

type Comment {
  id: ID!
  text: String!
}
`;

const schema = makeExecutableSchema({ typeDefs: schemaString });
addMockFunctionsToSchema({ schema });


function idReceiver(): Observable<string> {
  return new Observable<string>(observer => {
    let counter = 1;
    setInterval(() => {
      observer.next(`Comment:${counter}`);
      counter++;
    }, 1000);
  });
}

class SubscribeLink extends ApolloLink {
  constructor(private idObserver: Observable<string>) {
    super();
  }
  request(operation: Operation, forward?: NextLink): Observable<FetchResult> | null {
    const { query, operationName, variables } = operation;

    // アイディア
    // 1. subscriptionの要求を肩代わりする
    // 2. 何らかの接続先よりadd/update/deleteされたIDの一覧を受け取る
    // 3. queryを使って結果を得る
    // 4. queryで得た結果をsubscriptionの結果として返す

    // 1. subscriptionを見つける
    if (query.kind !== "Document") {
      if (forward) {
        return forward(operation);
      }
      return null;
    }
    const subscriptionDef = query.definitions.find(def => def.kind === "OperationDefinition" && def.operation === "subscription");
    if (!subscriptionDef) {
      if (forward) {
        return forward(operation);
      }
      return null;
    }

    // 肩代わり開始

    return new Observable(observer => {
      this.idObserver.subscribe(id => {
        console.log("receiver id", id);
        (async () => {
          try {
            await delay(800);
            const query = gql`
              query AlternativeGet($id: ID!) {
                commentAdded: comment(id: $id) {
                  # ここsubscriptionのselectionで代替したい
                  id
                  text
                }
              }
            `;
            const variables = { id };
            const result = await graphql(schema, print(query), null, null, variables, "AlternativeGet");
            console.log("alternative result", JSON.stringify(result, null, 2));
            observer.next(result);
            observer.complete();
          } catch (e) {
            observer.error(e);
          }
        })();
      });
    });
  }
}

function delay(ms: number) {
  return new Promise(resolve => {
    setTimeout(() => {
      resolve();
    }, ms);
  });
}

const client = new ApolloClient({
  cache: new InMemoryCache(),
  link: new SubscribeLink(idReceiver()),
});
client
  .subscribe({
    query: gql`
      subscription CommentAdded {
        commentAdded {
          id
          text
        }
      }
    `,
  }).subscribe(v => {
    console.log("client.subscribe", v);
  });
