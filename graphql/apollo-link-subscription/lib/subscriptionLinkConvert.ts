import { IntrospectionQuery, DocumentNode, print } from "graphql";

import { ApolloClient, ApolloQueryResult } from "apollo-client";
import { ApolloLink, Observable, Operation, NextLink, FetchResult } from "apollo-link";

import { subscriptionToQuery, AlternativeQueries } from "./documentModifier";

export type Action = "insert" | "update" | "delete";

export type EntityInfo = {
    fieldName: string;
    id: string;
    action: Action;
};

export type Options = {
    schema: IntrospectionQuery;
    observableConstructor: (originalOperation: Operation, alternativeQueries: AlternativeQueries) => Observable<EntityInfo>
    executeQuery?: (options: {
        query: DocumentNode; variables: { id: string; };
    }) => Promise<ApolloQueryResult<any>>;
};

export class SubscriptionCovertLink extends ApolloLink {

    private _schema: IntrospectionQuery;
    private _observableConstructor: Options["observableConstructor"];
    private _executeQuery?: Options["executeQuery"];

    constructor({ schema, observableConstructor, executeQuery }: Options) {
        super();
        this._schema = schema;
        this._observableConstructor = observableConstructor;
        this._executeQuery = executeQuery;
    }

    request(operation: Operation, forward?: NextLink): Observable<FetchResult> | null {
        const queries = subscriptionToQuery(this._schema, operation);
        if (!queries) {
            if (forward) {
                return forward(operation);
            }
            console.warn("forward link is undefined");
            return null;
        }

        const entityObserver = this._observableConstructor(operation, queries);

        return new Observable(observer => {
            entityObserver.subscribe(entity => {
                const query = queries[entity.fieldName];
                if (!query) {
                    throw new Error(`unknown fieldName: ${entity.fieldName}`);
                }

                let executeQuery = this._executeQuery;
                if (!executeQuery) {
                    const client: ApolloClient<never> = operation.getContext().client;
                    if (client instanceof ApolloClient === false) {
                        console.error(`client can't fine`, operation.getContext());
                        throw new Error(`client can't find`);
                    }

                    executeQuery = ({ query, variables }) => {
                        return client.query({
                            query,
                            variables,
                        });
                    };
                }

                executeQuery({
                    query,
                    variables: { id: entity.id },
                })
                    .then(result => {
                        observer.next(result);
                    })
                    .catch(e => {
                        observer.error(e);
                    });
            });
        });
    }
}
