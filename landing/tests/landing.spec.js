import { test, expect } from "@playwright/test";
import AxeBuilder from "@axe-core/playwright";
import { readFile } from "node:fs/promises";

test("workflows animate, cancel cleanly, and never make API calls", async ({
  page,
}) => {
  const requests = [];
  page.on("request", (request) => requests.push(request.url()));
  await page.goto("/");
  await expect(
    page.getByRole("link", { name: "See it in action" }),
  ).toBeVisible();
  await expect(page.locator("#run-demo, #replay")).toHaveCount(0);
  await expect(page.locator("body")).not.toContainText(/playground/i);
  await page.locator('[data-workflow="grades"]').click();
  await page.locator('[data-workflow="agents"]').click();
  await page.locator('[data-workflow="pipeline"]').click();
  await expect(page.locator("#demo-state")).toHaveText("Example shown");
  await expect(page.locator("#demo-command")).toContainText(
    "jq '.[] | {id, name}'",
  );
  await expect(page.locator("#demo-output")).toContainText("1042");
  await expect(page.locator("#workflow-docs")).toHaveAttribute(
    "href",
    /tutorials\/scripting\/$/,
  );
  await page.getByRole("button", { name: "Pause animations" }).click();
  for (const [workflow, command] of [
    ["grades", "bulk-grade"],
    ["agents", "mcp start"],
    ["courses", "--output json"],
  ]) {
    await page.locator(`[data-workflow="${workflow}"]`).click();
    await expect(page.locator("#demo-command")).toContainText(command);
    await expect(page.locator("#demo-state")).toHaveText("Example shown");
  }
  expect(
    requests.every((url) => url.startsWith("http://127.0.0.1:4174/")),
  ).toBe(true);
});

test("resource search provides actual docs and a useful empty state", async ({
  page,
}) => {
  await page.goto("/");
  await page.keyboard.press("/");
  await expect(page.locator("#resource-search")).toBeFocused();
  await page.locator("#resource-search").fill("rubric");
  await expect(page.locator("#resource-list a")).toHaveCount(1);
  await expect(page.locator("#resource-list a")).toHaveAttribute(
    "href",
    /commands\/canvas_rubrics\/$/,
  );
  await page.locator("#resource-search").fill("no-matching-resource");
  await expect(page.locator(".no-results")).toBeVisible();
  await expect(page.locator("#resource-list a")).toHaveText("All commands ↗");
});

test("interface and install tabs support keyboard navigation", async ({
  page,
  context,
}) => {
  await context.grantPermissions(["clipboard-read", "clipboard-write"]);
  await page.goto("/");
  await page.locator("#interface-human").focus();
  await page.keyboard.press("End");
  await expect(page.locator("#interface-agent")).toBeFocused();
  await expect(page.locator("#interface-panel")).toHaveAttribute(
    "aria-labelledby",
    "interface-agent",
  );
  await expect(page.locator("#interface-code")).toContainText(
    "canvas mcp start",
  );
  await page.locator("#tab-homebrew").focus();
  await page.keyboard.press("ArrowRight");
  await expect(page.locator("#tab-go")).toHaveAttribute(
    "aria-selected",
    "true",
  );
  await expect(page.locator("#install-command")).toContainText("go install");
  await page.locator("#copy-install").click();
  await expect(page.locator("#copy-status")).toContainText("Copied.");
  const clipboard = await page.evaluate(() => navigator.clipboard.readText());
  expect(clipboard).toBe(await page.locator("#install-command").textContent());
});

test("layout stays within the viewport and content works without JavaScript", async ({
  browser,
}) => {
  const page = await browser.newPage({ reducedMotion: "reduce" });
  await page.goto("http://127.0.0.1:4174/");
  for (const width of [320, 375, 600, 768, 1024, 1440]) {
    await page.setViewportSize({ width, height: 900 });
    expect(
      await page.evaluate(
        () => document.documentElement.scrollWidth <= innerWidth,
      ),
    ).toBe(true);
  }
  const staticPage = await browser.newPage({ javaScriptEnabled: false });
  await staticPage.goto("http://127.0.0.1:4174/");
  await expect(staticPage.locator("h1")).toBeVisible();
  await expect(staticPage.locator("#install-command")).toContainText(
    "brew install canvas-cli",
  );
  await staticPage.locator("summary").first().click();
  await expect(staticPage.locator("details").first()).toHaveAttribute(
    "open",
    "",
  );
  await page.close();
  await staticPage.close();
});

test("reduced motion renders immediately and accessibility checks pass", async ({
  page,
}) => {
  await page.emulateMedia({ reducedMotion: "reduce" });
  await page.goto("/");
  await expect(
    page.getByRole("button", { name: "Resume animations" }),
  ).toHaveAttribute("aria-pressed", "true");
  await page.locator('[data-workflow="agents"]').click();
  await expect(page.locator("#demo-state")).toHaveText("Example shown");
  await page.locator("#interface-agent").click();
  const results = await new AxeBuilder({ page })
    .withTags(["wcag2a", "wcag2aa", "wcag21aa"])
    .analyze();
  expect(results.violations).toEqual([]);
});

test("links target committed documentation and preview stays non-indexable", async ({
  page,
}) => {
  await page.goto("/");
  await expect(page.locator('meta[name="robots"]')).toHaveAttribute(
    "content",
    "noindex, follow",
  );
  const routes = await page
    .locator('a[href^="https://jjuanrivvera.github.io/canvas-cli/"]')
    .evaluateAll((links) =>
      links.map((link) =>
        new URL(link.href).pathname
          .replace("/canvas-cli/", "")
          .replace(/\/$/, ""),
      ),
    );
  for (const route of new Set(routes)) {
    const file = route ? `../docs/${route}.md` : "../docs/index.md";
    const alternate = `../docs/${route}/index.md`;
    const contents = await readFile(file, "utf8").catch(() =>
      readFile(alternate, "utf8"),
    );
    expect(contents.length).toBeGreaterThan(0);
  }
  const schema = JSON.parse(
    await page.locator('script[type="application/ld+json"]').textContent(),
  );
  expect(schema["@type"]).toBe("SoftwareApplication");
  expect(schema.name).toBe("Canvas CLI");
});
