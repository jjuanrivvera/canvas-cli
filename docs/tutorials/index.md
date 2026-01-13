# Tutorials

Step-by-step guides for common Canvas CLI workflows.

## Available Tutorials

<div class="grid cards" markdown>

-   :material-table-edit:{ .lg .middle } **Bulk Grading**

    ---

    Grade multiple submissions efficiently using CSV files

    [:octicons-arrow-right-24: Bulk Grading](bulk-grading.md)

-   :material-sync:{ .lg .middle } **Course Sync**

    ---

    Synchronize courses between Canvas instances

    [:octicons-arrow-right-24: Course Sync](course-sync.md)

-   :material-code-braces:{ .lg .middle } **Scripting & Automation**

    ---

    Automate tasks with shell scripts and pipelines

    [:octicons-arrow-right-24: Scripting](scripting.md)

</div>

## Quick Tips

!!! tip "Use JSON Output for Scripts"
    When writing scripts, always use `-o json` for reliable parsing:
    ```bash
    canvas courses list -o json | jq '.[].id'
    ```

!!! tip "Batch Operations"
    Canvas CLI supports concurrent batch operations for better performance:
    ```bash
    canvas submissions grade-batch --file grades.csv --course-id 123
    ```

!!! tip "Caching"
    Enable caching for repeated queries to reduce API calls:
    ```bash
    canvas courses list  # Cached for subsequent calls
    canvas courses list --no-cache  # Force fresh data
    ```
