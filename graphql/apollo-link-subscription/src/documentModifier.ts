import { DefinitionNode, OperationDefinitionNode } from "graphql";
import { Operation } from "apollo-link";

export function findSubscription(operation: Operation): OperationDefinitionNode | null {
    const { query, operationName } = operation;

    if (query.kind !== "Document") {
        return null;
    }

    let def: DefinitionNode | undefined;
    if (operationName) {
        def = query.definitions.find(def => {
            if (def.kind !== "OperationDefinition") {
                return false;
            }
            if (!def.name) {
                return false;
            }
            return def.name.value === operationName;
        });
    } else {
        def = query.definitions[0];
    }
    if (!def) {
        return null;
    }
    if (def.kind !== "OperationDefinition" || def.operation !== "subscription") {
        return null;
    }

    return def;
}

export function subscriptionToQuery(operation: Operation): Operation | null {

    const def = findSubscription(operation);
    if (!def) {
        return null;
    }

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


    return null;
}
