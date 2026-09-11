import assert from "node:assert/strict";
const base = "https://canvas-cli.jjuanrivvera.com/";
const response = await fetch(base, { signal: AbortSignal.timeout(30000) });
assert.equal(response.status, 200);
assert.equal(response.url, base);
const html = await response.text();
assert(html.includes(`rel="canonical" href="${base}"`));
assert(!html.includes("noindex"));
assert.equal(response.headers.get("x-content-type-options"), "nosniff");
assert(
  response.headers
    .get("content-security-policy")
    ?.includes("frame-ancestors 'none'"),
);
for (const path of [
  "assets/style.css",
  "assets/app.js",
  "assets/logo.svg",
  "assets/social.png",
  "sitemap.xml",
  "robots.txt",
]) {
  const asset = await fetch(new URL(path, base), {
    signal: AbortSignal.timeout(30000),
  });
  assert.equal(asset.status, 200, path);
  if (path === "sitemap.xml" || path === "robots.txt")
    assert((await asset.text()).includes(base), path);
  else await asset.arrayBuffer();
}
const alias = await fetch("https://canvas-cli-jjuanrivvera.netlify.app/", {
  redirect: "manual",
  signal: AbortSignal.timeout(30000),
});
assert.equal(alias.status, 301);
assert.equal(alias.headers.get("location"), base);
console.log(
  "Public HTTPS, assets, headers, canonical metadata, and Netlify domain redirect verified.",
);
