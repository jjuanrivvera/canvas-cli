# Canvas CLI landing

Standalone marketing site. The MkDocs documentation remains at
https://jjuanrivvera.github.io/canvas-cli/ and is not changed by this site.

## Preview

The static build requires Node.js 18+ and Python 3 for the preview server.
No npm dependencies are needed to build or serve the site.

```sh
cd landing
npm run build
npm run preview
```

Open http://localhost:4173. The build goes to `dist/`. Preview builds include
`noindex` so an undecided or temporary domain is not indexed. Serve the build,
not the source directory, to include generated metadata.

## Production

The production site is https://canvas-cli.jjuanrivvera.com/. To build it locally:

```sh
SITE_URL=https://canvas-cli.jjuanrivvera.com/ npm run build
```

Deploy `landing/dist/` to any static host. This produces the canonical URL,
Open Graph URL and social image metadata, `sitemap.xml`, `robots.txt`, and
SoftwareApplication structured data. A production build removes preview
`noindex`. All local assets use relative paths and support subpath hosting.
Use HTTPS for clipboard support. Copy falls back to selecting the command
when the Clipboard API is unavailable.

Set the host's build command to `npm run build`, root to `landing`, output to
`dist`, and environment variable `SITE_URL` to the final URL. No SPA rewrite
is required. Configure compression and short caching for HTML; assets use
stable filenames and must be revalidated rather than cached immutably.
Redirect alternate hostnames to the canonical hostname at the hosting layer.
Submit the generated sitemap to search engines after deployment.

## Content and design

- `index.html`: semantic content, navigation, static examples, FAQ.
- `assets/style.css`: responsive layout, focus states, reduced-motion support.
- `assets/app.js`: animated workflows, searchable command resources, interface
  and installation tabs, clipboard, scroll reveals, and motion controls.
- `assets/logo.svg`, `logo-icon.svg`, and `favicon.svg`: original project assets.
- `assets/social.png`: 1200 × 630 social preview.
- `build.mjs`: dependency-free static build and domain-dependent SEO metadata.

No external fonts, trackers, runtime packages, or API calls are loaded. The
terminal uses illustrative data; it does not connect to a Canvas account.
Core content, the default command examples, installation, and FAQ remain
readable without JS. Alternate demos and filtering require JavaScript.
Endpoint and command counts reflect the repository README; update them when
API coverage changes. Documentation links intentionally use the existing URL.

## Validation

Install development dependencies and a Playwright browser to run the browser suite:

```sh
npm ci
npx playwright install chromium
npm test
```

The six browser tests cover workflow animations and cancellation, no external
API requests, resource filtering and empty states, real documentation routes,
keyboard tabs, clipboard contents, FAQ, no-JS content, preview SEO metadata,
reduced motion, and overflow at 320–1440px. Axe checks WCAG 2 A/AA and 2.1 AA.
You can use an existing Chromium with `PLAYWRIGHT_CHROMIUM_EXECUTABLE`.

The logo orbit, API packets, cursor, and scroll reveals respect reduced motion.
The header motion button also pauses animations; command demos show their final
state immediately while paused. `/` focuses the resource search when you're not
already typing into a field. The landing presents animated walkthroughs with sample data, not an executable
playground. Workflow selection changes the illustration; there are no run
controls. Examples never execute shell commands or call Canvas.

To regenerate the social preview after a design change:

```sh
npm run social
npm run build
```

The social source is `tools/social.html`; its generator uses Playwright only in
development. Production contains a pre-rendered PNG and no browser dependencies.
Automated checks complement manual desktop/mobile visual review; they do not
constitute a complete accessibility certification.

The redesigned production build scored 100 in local mobile Lighthouse audits
for performance, accessibility, best practices, and SEO. Hosting may affect
these results. Domain metadata was checked with a placeholder; preview builds
remain non-indexable until the real SITE_URL is configured.

## Deployment pipeline

`feature/*` changes go through a PR to `develop`. The `Landing` GitHub Actions
workflow runs browser/accessibility tests and validates production metadata.
After merge to `develop`, it deploys that validated artifact to Netlify and
checks the live HTTPS site, assets, security headers, and canonical redirect.
The CLI's `main` release branch and MkDocs deployment remain separate.

- Netlify project: `canvas-cli-jjuanrivvera`
- Domain: `canvas-cli.jjuanrivvera.com` (Cloudflare delegates this subdomain to Netlify DNS)
- GitHub environment: `landing-production`, restricted to the `develop` branch
- Environment secret: `NETLIFY_AUTH_TOKEN`
- Environment variable: `NETLIFY_SITE_ID`
- Hosting configuration and security headers: root `netlify.toml`

Production deploys use the pinned Netlify CLI development dependency. Its
`sharp` override keeps the CLI's optional image tooling on a patched version;
the landing itself does not use image processing or server-side functions.

The workflow can be dispatched manually from `develop` for a redeploy. Netlify
retains deploy history for rollback. Draft deploys use a non-indexable build.
To verify the public deployment locally, run `node tools/check-live.mjs`.

Cloudflare holds four NS records for `canvas-cli.jjuanrivvera.com`, pointing to
`dns1.p07.nsone.net` through `dns4.p07.nsone.net`. Netlify manages the subdomain's
NETLIFY record and automatic TLS. The parent `jjuanrivvera.com` zone remains on
Cloudflare. This delegation replaced the initial CNAME after Netlify's external
DNS certificate provisioning stalled despite successful DNS verification.
