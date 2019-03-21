import { print, graphqlSync, introspectionQuery, IntrospectionQuery } from "graphql";
import { makeExecutableSchema } from "graphql-tools";
import gql, { resetCaches } from "graphql-tag";

import { createOperation } from "apollo-link";

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
            favorite(id: ID): Favorite
        }

        type Subscription {
            commentAdded: Comment
            favoriteAdded: Favorite
        }

        interface Node {
            id: ID!
        }

        type Comment implements Node {
            id: ID!
            text: String!
        }

        type Favorite implements Node {
            id: ID!
            comment: Comment!
        }
    `;

    const schema = makeExecutableSchema({ typeDefs: schemaString, resolverValidationOptions: { requireResolversForResolveType: false } });
    const introspectionResult = graphqlSync(schema, introspectionQuery).data! as any as IntrospectionQuery;
    afterEach(() => {
        resetCaches();
    })

    test("simple case", () => {
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
        const result = subscriptionToQuery(introspectionResult, operation);
        expect(result).toBeDefined();
        const query = result!["commentAdded"];
        expect(query).toBeDefined();
        const expected = gql`
            query Resolve_Foo_commentAdded ($id: ID!) {
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
        expect(print(query)).toEqual(print(expected));
        expect(query).toEqual(expected);
    });

    test("contains fragment", () => {
        const operation = createOperation({}, {
            query: gql`
                subscription Foo {
                    commentAdded {
                        id
                        ...CommentFragment
                    }
                }

                fragment CommentFragment on Comment {
                    id
                    text
                }
            `,
        });
        operation.operationName = "Foo";
        const result = subscriptionToQuery(introspectionResult, operation);
        expect(result).toBeDefined();
        const query = result!["commentAdded"];
        const expected = gql`
            query Resolve_Foo_commentAdded ($id: ID!) {
                commentAdded: node(id: $id) {
                    ... on Comment {
                        id
                        ...CommentFragment
                    }
                }
            }

            fragment CommentFragment on Comment {
                id
                text
            }
        `;
        delete (query as any).loc;
        delete expected.loc;
        expect(print(query)).toEqual(print(expected));
        expect(query).toEqual(expected);
    });

    test("multiple operations", () => {
        const operation = createOperation({}, {
            query: gql`
                subscription FooA {
                    a: commentAdded {
                        ...CommentFragmentA
                    }
                }

                subscription FooB {
                    b: commentAdded {
                        ...CommentFragmentB
                    }
                }

                fragment CommentFragmentA on Comment {
                    id
                }
                fragment CommentFragmentB on Comment {
                    text
                }
            `,
        });
        operation.operationName = "FooB";
        const result = subscriptionToQuery(introspectionResult, operation);
        expect(result).toBeDefined();
        const query = result!["b"];
        expect(query).toBeDefined();
        const expected = gql`
            query Resolve_FooB_b ($id: ID!) {
                b: node(id: $id) {
                    ... on Comment {
                        ...CommentFragmentB
                    }
                }
            }

            fragment CommentFragmentB on Comment {
                text
            }
        `;
        delete (query as any).loc;
        delete expected.loc;
        expect(print(query)).toEqual(print(expected));
        expect(query).toEqual(expected);
    });
    test("multiple fields", () => {
        const operation = createOperation({}, {
            query: gql`
                subscription Foo {
                    commentAdded {
                        ...CommentFragment
                    }
                    favoriteAdded {
                        ...FavoriteFragment
                    }
                }

                fragment CommentFragment on Comment {
                    id
                }
                fragment FavoriteFragment on Favorite {
                    id
                    comment {
                        ...CommentFragment
                    }
                }
            `,
        });
        operation.operationName = "Foo";
        const result = subscriptionToQuery(introspectionResult, operation);
        expect(result).toBeDefined();
        {
            const query = result!["commentAdded"];
            expect(query).toBeDefined();
            const expected = gql`
                query Resolve_Foo_commentAdded ($id: ID!) {
                    commentAdded: node(id: $id) {
                        ... on Comment {
                            ...CommentFragment
                        }
                    }
                }

                fragment CommentFragment on Comment {
                    id
                }
            `;
            delete (query as any).loc;
            delete expected.loc;
            expect(print(query)).toEqual(print(expected));
            expect(query).toEqual(expected);
        }
        {
            const query = result!["favoriteAdded"];
            expect(query).toBeDefined();
            const expected = gql`
                query Resolve_Foo_favoriteAdded ($id: ID!) {
                    favoriteAdded: node(id: $id) {
                        ... on Favorite {
                            ...FavoriteFragment
                        }
                    }
                }

                fragment CommentFragment on Comment {
                    id
                }
                fragment FavoriteFragment on Favorite {
                    id
                    comment {
                        ...CommentFragment
                    }
                }
            `;
            delete (query as any).loc;
            delete expected.loc;
            expect(print(query)).toEqual(print(expected));
            expect(query).toEqual(expected);
        }
    });
    test("multiple fields with same type", () => {
        const operation = createOperation({}, {
            query: gql`
                subscription Foo {
                    a: commentAdded {
                        ...CommentFragmentA
                    }
                    b: commentAdded {
                        ...CommentFragmentB
                    }
                }

                fragment CommentFragmentA on Comment {
                    id
                }
                fragment CommentFragmentB on Comment {
                    text
                }
            `,
        });
        operation.operationName = "Foo";
        const result = subscriptionToQuery(introspectionResult, operation);
        expect(result).toBeDefined();
        {
            const query = result!["a"];
            expect(query).toBeDefined();
            const expected = gql`
                query Resolve_Foo_a ($id: ID!) {
                    a: node(id: $id) {
                        ... on Comment {
                            ...CommentFragmentA
                        }
                    }
                }

                fragment CommentFragmentA on Comment {
                    id
                }
            `;
            delete (query as any).loc;
            delete expected.loc;
            expect(print(query)).toEqual(print(expected));
            expect(query).toEqual(expected);
        }
        {
            const query = result!["b"];
            expect(query).toBeDefined();
            const expected = gql`
                query Resolve_Foo_b ($id: ID!) {
                    b: node(id: $id) {
                        ... on Comment {
                            ...CommentFragmentB
                        }
                    }
                }

                fragment CommentFragmentB on Comment {
                    text
                }
            `;
            delete (query as any).loc;
            delete expected.loc;
            expect(print(query)).toEqual(print(expected));
            expect(query).toEqual(expected);
        }
    });
});
