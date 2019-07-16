import fetch from "node-fetch";

const endpoint = "https://api.github.com/graphql";
const authToken = process.env["GITHUB_TOKEN"];

const author = "vvakame";
const ignoreOrgs = [
    "techbookfest",
    "TechBooster",
];
const start = new Date("2019-07-06T09:00:00Z");
const end = new Date("2019-07-13T09:00:00Z");

// https://developer.github.com/v4/explorer/
const query = `
{
    search(first: 100, query: "author:${author}", type: ISSUE) {
      nodes {
        __typename
        ... on Issue {
          id
          number
          title
          body
          createdAt
          closed
          closedAt
          url
          repository {
            owner {
              id
              login
            }
            name
          }
        }
        ... on PullRequest {
          id
          number
          title
          body
          createdAt
          closed
          closedAt
          url
          repository {
            owner {
              id
              login
            }
            name
          }
        }
      }
    }
  }
`

async function exec() {
    const resp = await fetch(endpoint, {
        method: "POST",
        headers: {
            Authorization: `bearer ${authToken}`,
        },
        body: `{"query":${JSON.stringify(query)}}`,
    });
    if (resp.status !== 200) {
        throw new Error(`error, ${resp.status} ${await resp.text()}`);
    }
    const data = await resp.json();

    const text = data.data.search.nodes
        .filter((v: any) => ignoreOrgs.indexOf(v.repository.owner.login) === -1)
        .filter((v: any) => {
            const createdAt = new Date(v.createdAt);
            return start.getTime() <= createdAt.getTime() && createdAt.getTime() < end.getTime();
        })
        .map((v: any) => `* ${v.title} ${v.createdAt}\n    * ${v.url}`).join("\n");
    console.log(text);
}

exec().catch(err => console.error(err));
