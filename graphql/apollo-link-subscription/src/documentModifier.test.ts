import { createOperation } from "apollo-link";
import gql from "graphql-tag";

import { findSubscription } from "./documentModifier";

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
