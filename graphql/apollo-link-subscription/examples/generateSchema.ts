import fs from "fs";
import path from "path";
import glob from "glob";
import { makeExecutableSchema } from "graphql-tools";
import { graphqlSync, introspectionQuery } from "graphql";

const typeDefs = glob
    .sync(path.resolve(__dirname, "./*.graphql"))
    .map(file => fs.readFileSync(file, { encoding: "utf8" }));

const graphqlSchemaObject = makeExecutableSchema({
    typeDefs: typeDefs,
    resolverValidationOptions: {
        requireResolversForResolveType: false,
    },
});

{
    const result = graphqlSync(graphqlSchemaObject, introspectionQuery).data;
    fs.writeFileSync(path.resolve(__dirname, "./schema.json"), JSON.stringify(result, null, 2));
}
