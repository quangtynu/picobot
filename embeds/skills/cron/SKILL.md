---
name: cron
description: Schedule reminders and recurring tasks
---

# Cron

Use the `cron` tool to schedule reminders or recurring tasks.

## Actions

- `add` — schedule a new reminder
- `list` — show all pending jobs
- `cancel` — remove a job by name

## Examples

Set a reminder:

```
cron(action="add", name="break-reminder", message="Time to take a break!", delay="20m")
```

Longer delay:

```
cron(action="add", name="standup", message="Daily standup in 5 minutes", delay="1h")
```

List and cancel:

```
cron(action="list")
cron(action="cancel", name="break-reminder")
```

## Delay Format

Use Go duration strings:

| User says | Delay value |
|---|---|
| 2 minutes | `2m` |
| 1 hour | `1h` |
| 1 hour 30 minutes | `1h30m` |
| 30 seconds | `30s` |
| 1 day | `24h` |
