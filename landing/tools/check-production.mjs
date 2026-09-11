import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
const dist = new URL("../dist/", import.meta.url);
const expected = "https://canvas-cli.jjuanrivvera.com/";
const html = await readFile(new URL("index.html", dist), "utf8");
assert(
  html.includes(`rel="canonical" href="${expected}"`),
  "Canonical URL must use the production domain",
);
assert(!html.includes("noindex"), "Production must be indexable");
assert(
  html.includes(`content="${expected}assets/social.png"`),
  "Social image must use the production domain",
);
const schema = JSON.parse(
  html.match(/<script type="application\/ld\+json">(.*?)<\/script>/s)[1],
);
assert.equal(schema.url, expected);
assert(
  (await readFile(new URL("sitemap.xml", dist), "utf8")).includes(
    `<loc>${expected}</loc>`,
  ),
);
assert(
  (await readFile(new URL("robots.txt", dist), "utf8")).includes(
    `Sitemap: ${expected}sitemap.xml`,
  ),
);
assert((await readFile(new URL("assets/social.png", dist))).length > 0);
console.log(
  "Production canonical, social image, structured data, sitemap, and robots verified.",
);
