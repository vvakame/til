import { graphqlSync, introspectionQuery, IntrospectionQuery } from "graphql";
import { makeExecutableSchema } from "graphql-tools";
import { createOperation } from "apollo-link";
import gql from "graphql-tag";

import { findSubscription, subscriptionToQuery } from "./documentModifier";

describe('example', () => {
    test("make gql tag generated object is equals", () => {
        const a = gql`query { id }`;
        const b = gql`query { id }`;
        expect(a).toBe(b);
    });
});

test('createOperation', () => {
    const query = gql`
        query A { id }
    `;
    const operation = createOperation({}, query);
    expect(operation.operationName).toBeUndefined();
});

describe('findSubscription', () => {
    test('gets false from query', () => {
        const operation = createOperation({}, {
            query: gql`
                query {
                    viewer {
                        id
                    }
                }
            `,
        })
        const result = findSubscription(operation);
        expect(result).toBeFalsy();
    });
    test('gets true from subscription', () => {
        const operation = createOperation({}, {
            query: gql`
                subscription {
                    viewer {
                        id
                    }
                }
            `,
        })
        const result = findSubscription(operation);
        expect(result).toBeTruthy();
    });

    test('pick query from multiple operations', () => {
        const operation = createOperation({}, {
            query: gql`
                query A { id }
                subscription B { id }
            `,
        });
        operation.operationName = "A";
        const result = findSubscription(operation);
        expect(result).toBeFalsy();
    });
    test('pick subscription from multiple operations', () => {
        const operation = createOperation({}, {
            query: gql`
                query A { id }
                subscription B { id }
            `,
        });
        operation.operationName = "B";
        const result = findSubscription(operation);
        expect(result).toBeTruthy();
    });
});

describe('subscriptionToQuery', () => {
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

    const schema = makeExecutableSchema({ typeDefs: schemaString, resolverValidationOptions: { requireResolversForResolveType: false } });
    const introspectionResult = graphqlSync(schema, introspectionQuery).data! as any as IntrospectionQuery;

    test("", () => {
        const operation = createOperation({}, {
            query: gql`
                subscription Foo {
                    commentAdded {
                        id
                        text
                    }
                }
            `,
        });
        operation.operationName = "Foo";
        const query = subscriptionToQuery(introspectionResult, operation);
        expect(query).toBeDefined();
        const expected = gql`
            query Resolve_Foo ($id: ID!) {
                commentAdded: node(id: $id) {
                    ... on Comment {
                        id
                        text
                    }
                }
            }
        `;
        delete (query as any).loc;
        delete expected.loc;
        expect(query).toEqual(expected);
    });
});
