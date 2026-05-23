// Package main provides the cronaudit command-line interface.
//
// cronaudit parses and validates crontab files, producing human-readable
// summaries of each scheduled entry. It reads from a file path argument
// or from standard input when no argument is given.
//
// Usage:
//
//	cronaudit [flags] [file]
//
// Flags:
//
//	-format string
//	      output format: text, markdown, json (default "text")
//	-invalid
//	      show only invalid entries
//	-valid
//	      show only valid entries
//
// Examples:
//
//	cronaudit /etc/crontab
//	cronaudit -format json /var/spool/cron/root
//	cronaudit -invalid < crontab.txt
//	cat /etc/crontab | cronaudit -format markdown
package main
