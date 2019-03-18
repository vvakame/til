import { DefinitionNode, OperationDefinitionNode, IntrospectionQuery, DocumentNode, visitWithTypeInfo, TypeInfo, buildClientSchema, visit, ASTNode } from "graphql";
import { Operation } from "apollo-link";
import gql from "graphql-tag";

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

export function subscriptionToQuery(introspectionResult: IntrospectionQuery, operation: Operation): { [fieldName: string]: DocumentNode; } | null {
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

    const result: { [fieldName: string]: DocumentNode; } = {};

    def.selectionSet.selections.forEach(sel => {
        // ここでの作戦
        // 1. 要求されたsubscriptionから結果の型を調べる
        if (!sel || sel.kind !== "Field") {
            return null;
        }
        const resultType = subscriptionType.fields.find(f => f.name === sel.name.value);
        if (!resultType || resultType.type.kind !== "OBJECT") {
            return null;
        }

        // 2. subscriptionで要求されたfieldの名前を調べてそれを真似する準備
        let fieldName;
        if (sel.alias) {
            fieldName = sel.alias.value;
        } else {
            fieldName = sel.name.value;
        }

        // 3. クエリを組み立てる
        //   フィールドの指定はsubscriptionからコピーするので適当
        const queryName = `Resolve_${operation.operationName}_${fieldName}`;
        const query: DocumentNode = gql`
            query ${queryName} ($id: ID!) {
                ${fieldName}: node(id: $id) {
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

        // 関連fragmentの抽出と移植
        let relatedFragments: string[] = [];
        function extractFragmentName(node: ASTNode) {
            visit(node, {
                enter(node) {
                    if (node.kind !== "FragmentSpread") {
                        return;
                    } else if (relatedFragments.some(fragmentName => fragmentName === node.name.value)) {
                        return;
                    }
                    relatedFragments = [...relatedFragments, node.name.value];
                }
            });
        }
        extractFragmentName(sel);
        debugger;
        while (true) {
            const before = relatedFragments.length;

            visit(operation.query, {
                enter(node) {
                    if (node.kind !== "FragmentDefinition") {
                        return;
                    } else if (relatedFragments.every(fragmentName => fragmentName !== node.name.value)) {
                        return;
                    }
                    extractFragmentName(node.selectionSet);
                }
            });

            const after = relatedFragments.length;
            if (before === after) {
                break;
            }
        }

        visit(operation.query, {
            enter(node) {
                if (node.kind !== "FragmentDefinition") {
                    return;
                } else if (relatedFragments.indexOf(node.name.value) === -1) {
                    return;
                }
                (query as any).definitions = [...query.definitions, node];
            }
        });

        result[fieldName] = query;
    });

    return result;
}
