import { DefinitionNode, OperationDefinitionNode, IntrospectionQuery, DocumentNode, SelectionNode } from "graphql";
import { Operation } from "apollo-link";
import gql from "graphql-tag";
import { nullLiteral } from "@babel/types";

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

export function subscriptionToQuery(introspectionResult: IntrospectionQuery, operation: Operation): DocumentNode | null {
    if (!introspectionResult.__schema.subscriptionType) {
        return null;
    }
    const subscriptionTypeName = introspectionResult.__schema.subscriptionType.name;
    const subscriptionType = introspectionResult.__schema.types.find(t => t.kind === "OBJECT" && t.name === subscriptionTypeName);
    if (!subscriptionType || subscriptionType.kind !== "OBJECT") {
        // type narrowing
        return null;
    }

    const def = findSubscription(operation);
    if (!def) {
        return null;
    }

    // ここでの作戦
    // 1. 要求されたsubscriptionから結果の型を調べる
    const sel = def.selectionSet.selections[0];
    if (!sel || sel.kind !== "Field") {
        return null;
    }
    const resultType = subscriptionType.fields.find(f => f.name === sel.name.value);
    if (!resultType || resultType.type.kind !== "OBJECT") {
        return null;
    }

    // 2. subscriptionで要求されたfieldの名前を調べてそれを真似する準備
    let name;
    if (sel.alias) {
        name = sel.alias.value;
    } else {
        name = sel.name.value;
    }

    // 3. クエリを組み立てる
    //   フィールドの指定はsubscriptionからコピーするので適当
    const queryName = `Resolve_${operation.operationName}`;
    const query: DocumentNode = gql`
      query ${queryName} ($id: ID!) {
        ${name}: node(id: $id) {
          ... on ${resultType.type.name} {
            __typename
          }
        }
      }
    `;

    // 4. クエリのfield部分をsubscriptionのものと差し替える
    {
        let def = query.definitions[0];
        if (!def || def.kind !== "OperationDefinition") {
            return null;
        }
        let field = def.selectionSet.selections[0];
        if (!field || field.kind !== "Field" || !field.selectionSet) {
            return null;
        }
        let inline = field.selectionSet.selections[0]
        if (!inline || inline.kind !== "InlineFragment") {
            return null;
        }

        (inline.selectionSet as any) = sel.selectionSet;
    }

    return query;
}
