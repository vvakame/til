import { serve } from "https://deno.land/x/std/http/server.ts";
const s = serve("0.0.0.0:8000");

async function main() {
  for await (const {req, res} of s) {
    res.respond({ body: new TextEncoder().encode("Hello World\n") });
  }
}

main();
