import { graphql, graphqlSync, print, introspectionQuery } from "graphql";
import { makeExecutableSchema, addMockFunctionsToSchema } from 'graphql-tools';
import { ApolloLink, Observable, Operation, NextLink, FetchResult } from "apollo-link";

import gql from "graphql-tag";
import { ApolloClient } from "apollo-client";
import { InMemoryCache } from "apollo-cache-inmemory";

const schemaString = `
type Query {
  node(id: ID!): Node
  comment(id: ID): Comment
}

type Subscription {
  commentAdded: Comment
}

interface Node {
  id: ID!
}

type Comment implements Node {
  id: ID!
  text: String!
}
`;

const schema = makeExecutableSchema({ typeDefs: schemaString });
addMockFunctionsToSchema({ schema });
const introspectionResult = graphqlSync(schema, introspectionQuery).data!;

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
    const def = query.definitions.find(def => def.kind === "OperationDefinition" && def.operation === "subscription");
    if (!def) {
      if (forward) {
        return forward(operation);
      }
      return null;
    }
    if (def.kind !== "OperationDefinition" || def.operation !== "subscription") {
      // for type narrowing
      throw new Error("unexpected state");
    }

    // 肩代わり開始

    return new Observable(observer => {
      this.idObserver.subscribe(id => {
        console.log("receiver id", id);
        (async () => {
          try {
            await delay(800);
            const sel = def.selectionSet.selections[0];
            let name = "default";
            if (sel.kind === "Field") {
              if (sel.alias) {
                name = sel.alias.value;
              } else {
                name = sel.name.value;
              }
              sel.selectionSet
            }
            const subscriptionType: string = introspectionResult.__schema.subscriptionType.name;
            const subscription = introspectionResult.__schema.types.find((t: any) => t.kind === "OBJECT" && t.name === subscriptionType);
            const resultType = subscription.fields.find((f: any) => f.name === (sel as any).name.value);

            const queryName = `Resolve_${operationName}`;
            const query = gql`
              query ${queryName} ($id: ID!) {
                ${name}: node(id: $id) {
                  ... on ${resultType.type.name} {
                    __typename
                  }
                }
              }
            `;
            (query as any).definitions[0].selectionSet.selections[0].selectionSet.selections[0].selectionSet = (sel as any).selectionSet;
            console.log(print(query));

            const variables = { id };
            const result = await graphql(schema, print(query), null, null, variables, queryName);
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
