import { mkdir, readFile, writeFile, cp, rm } from "node:fs/promises";
import { fileURLToPath } from "node:url";
import path from "node:path";

const root = path.dirname(fileURLToPath(import.meta.url));
const output = path.join(root, "dist");
const configuredUrl = process.env.SITE_URL;
// GA4 only ships on an indexable production build: a preview or PR build would
// otherwise send hits from a throwaway URL into the same property.
const analyticsId = (process.env.GOOGLE_ANALYTICS_KEY ?? "").trim();
if (analyticsId && !/^G-[A-Z0-9]{4,}$/.test(analyticsId)) {
  throw new Error(
    `GOOGLE_ANALYTICS_KEY must look like a GA4 measurement ID (G-XXXXXXX), got ${JSON.stringify(analyticsId)}.`,
  );
}
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
// Google Search Console verifies ownership by fetching this file from the site
// root, so it has to be copied verbatim, not templated.
const searchConsoleFile = "google9631057cc493be1e.html";
await cp(
  path.join(root, searchConsoleFile),
  path.join(output, searchConsoleFile),
);
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
// The FAQ answers already live in the page; deriving the schema from that markup
// keeps the two from drifting apart the way a hand-written copy would.
const stripTags = (value) =>
  value
    .replace(/<span aria-hidden="true">.*?<\/span>/gs, "")
    .replace(/<[^>]+>/g, " ")
    .replace(/&amp;/g, "&")
    .replace(/&quot;/g, '"')
    .replace(/&#39;/g, "'")
    .replace(/&lt;/g, "<")
    .replace(/&gt;/g, ">")
    .replace(/\s+/g, " ")
    .trim();
const faqSection = html.match(/<div class="faq-list">(.*?)<\/div>/s);
const faqEntries = [
  ...(faqSection?.[1] ?? "").matchAll(
    /<details>\s*<summary>(.*?)<\/summary>\s*<p>(.*?)<\/p>\s*<\/details>/gs,
  ),
].map(([, question, answer]) => ({
  "@type": "Question",
  name: stripTags(question),
  acceptedAnswer: { "@type": "Answer", text: stripTags(answer) },
}));
if (faqEntries.length === 0) {
  throw new Error(
    "No FAQ entries found in index.html — the FAQPage schema would ship empty.",
  );
}
const faqSchema = {
  "@context": "https://schema.org",
  "@type": "FAQPage",
  mainEntity: faqEntries,
  ...(siteUrl ? { url: siteUrl } : {}),
};

let metadata = `<script type="application/ld+json">${JSON.stringify(schema).replaceAll("<", "\\u003c")}</script>`;
metadata += `\n<script type="application/ld+json">${JSON.stringify(faqSchema).replaceAll("<", "\\u003c")}</script>`;
// og:title/description already cover X, but naming them explicitly stops a card
// from silently falling back to whatever a scraper decides to infer.
metadata +=
  '\n<meta property="og:site_name" content="Canvas CLI">' +
  '\n<meta property="og:locale" content="en_US">' +
  '\n<meta name="twitter:title" content="Canvas CLI — Your entire LMS. Under command.">' +
  '\n<meta name="twitter:description" content="One binary. Human commands, automated pipelines, AI tools. Put the Canvas API to work from your terminal.">' +
  '\n<meta name="author" content="Juan Rivera">';
if (analyticsId && siteUrl) {
  metadata +=
    `\n<script async src="https://www.googletagmanager.com/gtag/js?id=${escape(analyticsId)}"></script>` +
    "\n<script>window.dataLayer=window.dataLayer||[];function gtag(){dataLayer.push(arguments);}" +
    `gtag('js', new Date());gtag('config', '${escape(analyticsId)}');</script>`;
} else if (analyticsId) {
  console.log(
    "GOOGLE_ANALYTICS_KEY set but SITE_URL is not: skipping GA4 on this preview build.",
  );
}
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
