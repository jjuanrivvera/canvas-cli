import { mkdir, readFile, writeFile, cp, rm } from "node:fs/promises";
import { fileURLToPath } from "node:url";
import path from "node:path";

const root = path.dirname(fileURLToPath(import.meta.url));
const output = path.join(root, "dist");
const configuredUrl = process.env.SITE_URL;
let siteUrl;
if (configuredUrl) {
  const parsed = new URL(configuredUrl);
  if (
    !["https:", "http:"].includes(parsed.protocol) ||
    parsed.username ||
    parsed.password ||
    parsed.search ||
    parsed.hash
  ) {
    throw new Error(
      "SITE_URL must be a public HTTP(S) URL without credentials, query, or fragment.",
    );
  }
  siteUrl = parsed.href.replace(/\/?$/, "/");
}
await rm(output, { recursive: true, force: true });
await mkdir(output, { recursive: true });
await cp(path.join(root, "assets"), path.join(output, "assets"), {
  recursive: true,
});
let html = await readFile(path.join(root, "index.html"), "utf8");
const escape = (value) =>
  value
    .replaceAll("&", "&amp;")
    .replaceAll('"', "&quot;")
    .replaceAll("<", "&lt;")
    .replaceAll(">", "&gt;");
const schema = {
  "@context": "https://schema.org",
  "@type": "SoftwareApplication",
  name: "Canvas CLI",
  applicationCategory: "DeveloperApplication",
  operatingSystem: "macOS, Linux, Windows",
  description:
    "An open-source command-line interface for Canvas LMS with bulk operations, multiple output formats, and an MCP server.",
  license: "https://opensource.org/license/mit",
  downloadUrl: "https://github.com/jjuanrivvera/canvas-cli/releases",
  offers: { "@type": "Offer", price: "0", priceCurrency: "USD" },
  author: {
    "@type": "Person",
    name: "Juan Rivera",
    url: "https://github.com/jjuanrivvera",
  },
  ...(siteUrl ? { url: siteUrl } : {}),
};
let metadata = `<script type="application/ld+json">${JSON.stringify(schema).replaceAll("<", "\\u003c")}</script>`;
if (siteUrl) {
  metadata += `\n<link rel="canonical" href="${escape(siteUrl)}">\n<meta property="og:url" content="${escape(siteUrl)}">\n<meta property="og:image" content="${escape(siteUrl)}assets/social.png">\n<meta property="og:image:width" content="1200">\n<meta property="og:image:height" content="630">\n<meta property="og:image:alt" content="Canvas CLI. Your entire LMS. Under command. One binary for humans, scripts, and AI agents.">`;
  await writeFile(
    path.join(output, "sitemap.xml"),
    `<?xml version="1.0" encoding="UTF-8"?>\n<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9"><url><loc>${escape(siteUrl)}</loc></url></urlset>\n`,
  );
  await writeFile(
    path.join(output, "robots.txt"),
    `User-agent: *\nAllow: /\nSitemap: ${siteUrl}sitemap.xml\n`,
  );
} else {
  metadata += '\n<meta name="robots" content="noindex, follow">';
  await writeFile(path.join(output, "robots.txt"), "User-agent: *\nAllow: /\n");
  console.log(
    "Preview build: noindex enabled. Set SITE_URL for an indexable production build with canonical URL and sitemap.",
  );
}
html = html.replace(
  "<!-- Deployment metadata is injected by build.mjs when SITE_URL is provided. -->",
  metadata,
);
await writeFile(path.join(output, "index.html"), html);
console.log("Built landing/dist.");
