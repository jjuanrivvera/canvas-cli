# User Guide

Learn how to configure and use Canvas CLI effectively.

## Topics

<div class="grid cards" markdown>

-   :material-cog:{ .lg .middle } **Configuration**

    ---

    Configure Canvas CLI for your environment, including multiple instances

    [:octicons-arrow-right-24: Configuration](configuration.md)

-   :material-format-list-bulleted:{ .lg .middle } **Output Formats**

    ---

    Learn about the different output formats available

    [:octicons-arrow-right-24: Output Formats](output-formats.md)

-   :material-console:{ .lg .middle } **Shell Completion**

    ---

    Set up tab completion for your shell

    [:octicons-arrow-right-24: Shell Completion](shell-completion.md)

</div>

## Key Concepts

### Course Context

Most commands require a course context. You can specify this with the `--course-id` flag:

```bash
canvas assignments list --course-id 123
```

### Output Formats

Canvas CLI supports multiple output formats:

- **table** (default) - Human-readable table format
- **json** - JSON for scripting and automation
- **yaml** - YAML format
- **csv** - CSV for spreadsheet import

### Caching

Canvas CLI caches API responses for better performance. You can:

- Disable caching with `--no-cache`
- Clear the cache with `canvas cache clear`
- View cache statistics with `canvas cache stats`

### Masquerading

Administrators can act as other users with the `--as-user` flag:

```bash
canvas courses list --as-user 456
```

This requires masquerading permissions in Canvas.
