# cronaudit

> Lightweight utility to parse and validate crontabs with human-readable summaries

---

## Installation

```bash
go install github.com/yourusername/cronaudit@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/cronaudit.git
cd cronaudit
go build -o cronaudit .
```

---

## Usage

Validate and summarize a crontab file:

```bash
cronaudit /etc/crontab
```

Pipe directly from crontab output:

```bash
crontab -l | cronaudit -
```

**Example output:**

```
✔  */5 * * * *        → Every 5 minutes
✔  0 2 * * 1          → At 02:00 AM, every Monday
✔  30 8 1 * *         → At 08:30 AM, on the 1st of every month
✘  60 25 * * *        → Invalid: hour value 25 out of range (0-23)
```

### Flags

| Flag        | Description                          |
|-------------|--------------------------------------|
| `-f`        | Path to crontab file                 |
| `-strict`   | Exit with non-zero code on any error |
| `-json`     | Output results as JSON               |

---

## License

[MIT](LICENSE) © 2024 yourusername