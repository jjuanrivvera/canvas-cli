"use strict";
const reducedMotion = window.matchMedia("(prefers-reduced-motion: reduce)");
let motionPaused = reducedMotion.matches;
const motionToggle = document.querySelector(".motion-toggle");
function updateMotion() {
  document.documentElement.classList.toggle("motion-paused", motionPaused);
  motionToggle.setAttribute("aria-pressed", String(motionPaused));
  motionToggle.setAttribute(
    "aria-label",
    motionPaused ? "Resume animations" : "Pause animations",
  );
  motionToggle.title = motionPaused ? "Resume animations" : "Pause animations";
  motionToggle.firstElementChild.textContent = motionPaused ? "▷" : "Ⅱ";
}
updateMotion();
motionToggle.addEventListener("click", () => {
  motionPaused = !motionPaused;
  updateMotion();
  if (motionPaused) finishDemo();
});
reducedMotion.addEventListener("change", (event) => {
  motionPaused = event.matches;
  updateMotion();
  if (motionPaused) finishDemo();
});
if ("IntersectionObserver" in window) {
  document.documentElement.classList.add("js-motion");
  const observer = new IntersectionObserver(
    (entries) => {
      entries.forEach((entry) => {
        if (!entry.isIntersecting) return;
        entry.target.classList.add("visible");
        observer.unobserve(entry.target);
      });
    },
    { threshold: 0.08 },
  );
  document
    .querySelectorAll(".reveal")
    .forEach((element) => observer.observe(element));
}
let scrollPending = false;
function updateScroll() {
  const available = document.documentElement.scrollHeight - innerHeight;
  const position = available > 0 ? scrollY / available : 0;
  document.querySelector(".scroll-progress").style.transform =
    `scaleX(${position})`;
  scrollPending = false;
}
window.addEventListener(
  "scroll",
  () => {
    if (scrollPending) return;
    scrollPending = true;
    requestAnimationFrame(updateScroll);
  },
  { passive: true },
);
window.addEventListener("resize", updateScroll);
updateScroll();

const docs = "https://jjuanrivvera.github.io/canvas-cli/";
const demos = {
  courses: {
    comment: "# Every course you can access. One command.",
    command: "canvas courses list --output json",
    context: "READ / COURSES",
    output:
      '[\n  { "id": 1042, "name": "Design Fundamentals",\n    "workflow_state": "available" },\n  { "id": 1043, "name": "Creative Coding",\n    "workflow_state": "available" }\n]',
    result: "Structured output. Ready for whatever comes next.",
    insight:
      "Stop copying data out of the browser. Get consistent, structured output you can filter, save, and reuse.",
    link: "user-guide/output-formats/",
  },
  grades: {
    comment: "# A whole class of grades. One CSV import.",
    command:
      "canvas submissions bulk-grade \\\n  --course-id 123 --csv-file grades.csv",
    context: "WRITE / SUBMISSIONS",
    output:
      "# grades.csv — sample input\nuser_id,assignment_id,score,comment\n1001,456,85,Good work\n1002,456,92,Excellent analysis\n1003,456,88,Strong argument\n\n# The CLI applies each row to its submission.",
    result: "Grade offline. Import together. Verify in Canvas.",
    insight:
      "Keep your spreadsheet workflow. Apply scores and comments from a CSV instead of opening each submission individually.",
    link: "tutorials/bulk-grading/",
  },
  pipeline: {
    comment: "# Your LMS data belongs in your toolchain.",
    command: "canvas courses list --output json \\\n  | jq '.[] | {id, name}'",
    context: "PIPE / JSON → JQ",
    output:
      '{\n  "id": 1042,\n  "name": "Design Fundamentals"\n}\n{\n  "id": 1043,\n  "name": "Creative Coding"\n}',
    result: "No scraping. No custom HTTP client. Just a pipe.",
    insight:
      "Compose Canvas operations with the tools you already use. Filter with jq, export to files, or build repeatable shell workflows.",
    link: "tutorials/scripting/",
  },
  agents: {
    comment: "# Same binary. A typed tool interface for your agent.",
    command: "canvas mcp start",
    context: "STDIO / MCP SERVER",
    output:
      "# Configure your AI client to launch this command.\n# CLI flags become typed tool parameters.\n# Tool results use structured JSON.\n\n# Optional: generate host-specific safety rules\ncanvas agent guard --host claude-code\n\n# STDIO is reserved for MCP protocol messages.",
    result: "Your Canvas workflows, callable by your assistant.",
    insight:
      "Connect an MCP-compatible assistant to the Canvas command surface. Add the bundled skill and configure guard rules for your agent host.",
    link: "user-guide/mcp/",
  },
};
const commandElement = document.getElementById("demo-command");
const outputElement = document.getElementById("demo-output");
const demoState = document.getElementById("demo-state");
let currentDemo = "courses";
let demoTimer;
let demoRevision = 0;
function highlightedCode(value) {
  const fragment = document.createDocumentFragment();
  const tokens = value.split(/("(?:[^"\\]|\\.)*"|\b\d+\b|#[^\n]*)/g);
  tokens.forEach((token, index) => {
    const span = document.createElement("span");
    span.textContent = token;
    if (token.startsWith('"'))
      span.className = /^\s*:/.test(tokens[index + 1] || "")
        ? "syntax-purple"
        : "syntax-green";
    else if (/^\d+$/.test(token)) span.className = "syntax-amber";
    else if (token.startsWith("#")) span.className = "syntax-dim";
    fragment.append(span);
  });
  return fragment;
}
function renderOutput() {
  const demo = demos[currentDemo];
  const pre = document.createElement("pre");
  const code = document.createElement("code");
  code.append(highlightedCode(demo.output));
  pre.append(code);
  const result = document.createElement("p");
  result.className = "terminal-result";
  const tick = document.createElement("span");
  tick.textContent = "✓";
  result.append(tick, document.createTextNode(demo.result));
  outputElement.replaceChildren(pre, result);
  outputElement.classList.remove("waiting");
  demoState.textContent = "Example shown";
}
function finishDemo() {
  clearTimeout(demoTimer);
  demoRevision++;
  commandElement.textContent = demos[currentDemo].command;
  renderOutput();
}
function showWorkflow() {
  clearTimeout(demoTimer);
  const revision = ++demoRevision;
  const demo = demos[currentDemo];
  document.getElementById("terminal-comment").textContent = demo.comment;
  document.getElementById("console-context").textContent = demo.context;
  document.getElementById("workflow-insight").textContent = demo.insight;
  document.getElementById("workflow-docs").href = docs + demo.link;
  if (motionPaused || reducedMotion.matches) {
    commandElement.textContent = demo.command;
    renderOutput();
    return;
  }
  demoState.textContent = "Showing example…";
  outputElement.classList.add("waiting");
  commandElement.textContent = "";
  let character = 0;
  function typeNext() {
    if (revision !== demoRevision) return;
    character += 2;
    commandElement.textContent = demo.command.slice(0, character);
    if (character < demo.command.length) demoTimer = setTimeout(typeNext, 20);
    else demoTimer = setTimeout(renderOutput, 180);
  }
  typeNext();
}
document.querySelectorAll("[data-workflow]").forEach((button) => {
  button.addEventListener("click", () => {
    currentDemo = button.dataset.workflow;
    document.querySelectorAll("[data-workflow]").forEach((item) => {
      item.classList.toggle("active", item === button);
      item.setAttribute("aria-pressed", String(item === button));
    });
    showWorkflow();
  });
});
// Announce completion, rather than every keystroke in the typing animation.
demoState.setAttribute("role", "status");

const resources = [
  "courses",
  "assignments",
  "submissions",
  "enrollments",
  "modules",
  "pages",
  "quizzes",
  "users",
  "rubrics",
  "outcomes",
  "analytics",
  "sis-imports",
  "groups",
  "discussions",
  "files",
  "sections",
  "accounts",
  "blueprint",
  "calendar",
  "roles",
  "content-migrations",
  "grading-periods",
  "announcements",
  "webhook",
];
const resourceSearch = document.getElementById("resource-search");
const resourceList = document.getElementById("resource-list");
function renderResources() {
  const query = resourceSearch.value.trim().toLowerCase();
  const matches = resources.filter((name) => name.includes(query));
  const visible = query ? matches : matches.slice(0, 12);
  resourceList.replaceChildren();
  visible.forEach((name) => {
    const link = document.createElement("a");
    link.href = `${docs}commands/canvas_${name}/`;
    link.append(document.createTextNode(name));
    const arrow = document.createElement("span");
    arrow.textContent = "↗";
    arrow.setAttribute("aria-hidden", "true");
    link.append(arrow);
    resourceList.append(link);
  });
  if (!matches.length) {
    const message = document.createElement("p");
    message.className = "no-results";
    message.textContent =
      "No matches in this selection. Explore all 93 groups in the command reference.";
    resourceList.append(message);
    const link = document.createElement("a");
    link.href = docs + "commands/";
    link.textContent = "All commands ↗";
    resourceList.append(link);
  }
  document.getElementById("resource-count").textContent = query
    ? `${matches.length} matching resources in this curated selection`
    : "12 featured resources · Search 24 common groups · 93 groups in the docs";
}
resourceSearch.addEventListener("input", renderResources);
renderResources();
document.addEventListener("keydown", (event) => {
  if (
    event.key !== "/" ||
    event.ctrlKey ||
    event.metaKey ||
    event.altKey ||
    event.target.closest("input, textarea, select, [contenteditable]")
  )
    return;
  event.preventDefault();
  resourceSearch.focus();
});

function setupTabs(selector, select) {
  const tabs = [...document.querySelectorAll(selector)];
  function activate(tab) {
    tabs.forEach((item) => {
      const selected = item === tab;
      item.setAttribute("aria-selected", String(selected));
      item.tabIndex = selected ? 0 : -1;
    });
    select(tab);
  }
  tabs.forEach((tab, index) => {
    tab.addEventListener("click", () => activate(tab));
    tab.addEventListener("keydown", (event) => {
      let next;
      if (event.key === "ArrowRight") next = (index + 1) % tabs.length;
      if (event.key === "ArrowLeft")
        next = (index - 1 + tabs.length) % tabs.length;
      if (event.key === "Home") next = 0;
      if (event.key === "End") next = tabs.length - 1;
      if (next === undefined) return;
      event.preventDefault();
      activate(tabs[next]);
      tabs[next].focus();
    });
  });
}
const interfaces = {
  human: {
    eyebrow: "YOUR TERMINAL, SUPERCHARGED",
    title: "Less remembering.\nMore executing.",
    description:
      "Explore Canvas in the interactive shell, with command history and completion. Set a course context once and stop repeating the same flags.",
    link: "user-guide/context/",
    label: "Explore context management ↗",
    file: "session.sh",
    code: "# Start an interactive session\ncanvas repl\n\n# Set your working course\ncanvas context set course 123\n\n# Use it across commands\ncanvas assignments list",
  },
  script: {
    eyebrow: "REPEATABLE BEATS REPETITIVE",
    title: "Your shell is\nthe integration layer.",
    description:
      "Export structured data, target named instances, and compose commands with standard shell tools. Make a one-off task a workflow you can run again.",
    link: "tutorials/scripting/",
    label: "Build an automation ↗",
    file: "export-courses.sh",
    code: "# Export a snapshot from a named instance\ncanvas courses list \\\n  --instance production \\\n  --output json > courses.json\n\n# Work with the result in your own pipeline\njq '.[] | {id, name}' courses.json",
  },
  agent: {
    eyebrow: "THE SAME COMMANDS, NOW CALLABLE",
    title: "An API for your\nAI workflow.",
    description:
      "Run Canvas CLI as an MCP server. Command flags become typed tool parameters; results come back as JSON. Add an agent skill and generate guard rules for your host.",
    link: "user-guide/mcp/",
    label: "Configure your MCP client ↗",
    file: "agent-setup.sh",
    code: "# Launch through your MCP client\ncanvas mcp start\n\n# Install CLI know-how for your agent\ncanvas skills install\n\n# Preview guard rules before writing them\ncanvas agent guard --host claude-code",
  },
};
setupTabs("[data-interface]", (tab) => {
  const selected = interfaces[tab.dataset.interface];
  const panel = document.getElementById("interface-panel");
  panel.setAttribute("aria-labelledby", tab.id);
  document.getElementById("interface-eyebrow").textContent = selected.eyebrow;
  const title = document.getElementById("interface-title");
  const lines = selected.title.split("\n");
  title.replaceChildren(
    document.createTextNode(lines[0]),
    document.createElement("br"),
    document.createTextNode(lines[1]),
  );
  document.getElementById("interface-description").textContent =
    selected.description;
  const link = document.getElementById("interface-link");
  link.href = docs + selected.link;
  link.textContent = selected.label;
  document.getElementById("interface-file").textContent = selected.file;
  document.getElementById("interface-code").textContent = selected.code;
  panel.classList.remove("changing");
  requestAnimationFrame(() => panel.classList.add("changing"));
});
const commands = {
  homebrew: "brew tap jjuanrivvera/canvas-cli\nbrew install canvas-cli",
  go: "go install github.com/jjuanrivvera/canvas-cli/cmd/canvas@latest",
  docker: "docker run --rm ghcr.io/jjuanrivvera/canvas-cli:latest version",
  windows:
    "scoop install https://raw.githubusercontent.com/jjuanrivvera/scoop-canvas-cli/main/canvas-cli.json",
};
const installCommand = document.getElementById("install-command");
const copyStatus = document.getElementById("copy-status");
const installationNotes = {
  homebrew: "For macOS and Linux with Homebrew installed.",
  go: "Requires Go 1.25+ for MCP support. Add your Go bin directory to PATH.",
  docker:
    "Runs the container to show its version. For API calls, pass CANVAS_URL and CANVAS_TOKEN with Docker -e flags; host commands below apply to a native install.",
  windows: "Requires Scoop on Windows. Run in PowerShell.",
};
setupTabs("[data-install]", (tab) => {
  installCommand.textContent = commands[tab.dataset.install];
  document.getElementById("install-note").textContent =
    installationNotes[tab.dataset.install];
  document
    .getElementById("install-code")
    .setAttribute("aria-labelledby", tab.id);
  copyStatus.textContent = "";
});
document.getElementById("copy-install").addEventListener("click", async () => {
  try {
    await navigator.clipboard.writeText(installCommand.textContent);
    copyStatus.textContent = "Copied. Paste it into your terminal.";
  } catch {
    const selection = window.getSelection();
    const range = document.createRange();
    range.selectNodeContents(installCommand);
    selection.removeAllRanges();
    selection.addRange(range);
    copyStatus.textContent = "Command selected. Press Ctrl+C or ⌘C to copy.";
  }
});

const workbench = document.querySelector(".workbench");
let pointerFrame;
workbench.addEventListener("pointermove", (event) => {
  if (motionPaused || event.pointerType === "touch") return;
  cancelAnimationFrame(pointerFrame);
  pointerFrame = requestAnimationFrame(() => {
    const bounds = workbench.getBoundingClientRect();
    workbench.style.setProperty(
      "--pointer-x",
      `${event.clientX - bounds.left}px`,
    );
    workbench.style.setProperty(
      "--pointer-y",
      `${event.clientY - bounds.top}px`,
    );
  });
});
