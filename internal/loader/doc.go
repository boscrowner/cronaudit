// Package loader handles reading crontab content from various input sources
// such as files on disk or standard input.
//
// Typical usage:
//
//	src, err := loader.FromFile("/etc/crontab")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer loader.Close(src)
//
//	lines, err := loader.Lines(src)
//	if err != nil {
//		log.Fatal(err)
//	}
//
// The lines returned are then passed directly to the parser package for
// validation and summarisation.
package loader
