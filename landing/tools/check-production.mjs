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
const schemas = [
  ...html.matchAll(/<script type="application\/ld\+json">(.*?)<\/script>/gs),
].map(([, json]) => JSON.parse(json));
const schema = schemas.find(
  (entry) => entry["@type"] === "SoftwareApplication",
);
assert(schema, "SoftwareApplication structured data must be present");
assert.equal(schema.url, expected);

const faq = schemas.find((entry) => entry["@type"] === "FAQPage");
assert(faq, "FAQPage structured data must be present");
assert(
  faq.mainEntity.length > 0,
  "FAQPage must carry the questions rendered on the page",
);
for (const entry of faq.mainEntity) {
  assert(
    entry.name && entry.acceptedAnswer?.text,
    "Every FAQ entry needs a question and an answer",
  );
  assert(
    !entry.name.includes("<") && !entry.acceptedAnswer.text.includes("<"),
    "FAQ schema must carry plain text, not markup",
  );
}

// Search Console re-fetches this file periodically; losing it silently unverifies
// the property, and nothing else in the build would notice.
assert(
  (
    await readFile(new URL("google9631057cc493be1e.html", dist), "utf8")
  ).includes("google-site-verification"),
  "Search Console verification file must ship at the site root",
);

// GA4 is opt-in through the environment: assert whichever state the build was
// asked for, so a missing variable can never quietly turn analytics off.
const analyticsId = (process.env.GOOGLE_ANALYTICS_KEY ?? "").trim();
if (analyticsId) {
  assert(
    html.includes(`gtag/js?id=${analyticsId}`),
    "GOOGLE_ANALYTICS_KEY was set, so the production build must load GA4",
  );
} else {
  assert(
    !html.includes("googletagmanager.com"),
    "No GOOGLE_ANALYTICS_KEY was set, so no analytics should be embedded",
  );
}
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
  "Production canonical, social image, structured data, FAQ schema, Search Console file, analytics, sitemap, and robots verified.",
);
