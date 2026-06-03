package utils

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/nullt3r/udpx/pkg/probes"
)

type Options struct {
	Arg_t      string
	Arg_tf     string
	Arg_o      string
	Arg_c      int
	Arg_nr     bool
	Arg_st     int
	Arg_sp     bool
	Arg_s      string
	Arg_q           bool
	Arg_src_ip      string
	Arg_exclude     string
	Arg_excludefile string
}

func ParseOptions() *Options {
	opts := &Options{}
	flag.StringVar(&opts.Arg_t, "t", "", "IPs/hostnames/CIDRs/octet-ranges to scan")
	flag.StringVar(&opts.Arg_tf, "tf", "", "File containing IPs/hostnames/CIDRs/octet-ranges to scan")
	flag.StringVar(&opts.Arg_o, "o", "", "Output file to write results")
	flag.StringVar(&opts.Arg_s, "s", "", fmt.Sprintf("Scan only for a specific service, one of: %s", probes.GetProbeNames()))
	flag.IntVar(&opts.Arg_c, "c", 32, "Maximum number of concurrent connections")
	flag.BoolVar(&opts.Arg_nr, "nr", false, "Do not randomize addresses")
	flag.IntVar(&opts.Arg_st, "w", 500, "Maximum time to wait for a response (socket timeout) in ms")
	flag.BoolVar(&opts.Arg_sp, "sp", false, "Show received packets (only first 32 bytes)")
	flag.BoolVar(&opts.Arg_q, "q", false, "Quiet mode: suppress banner and progress log lines (results still emitted)")
	flag.StringVar(&opts.Arg_src_ip, "src-ip", "", "Source IP to bind probes to (must be assigned to a local interface; overrides the kernel's default route-table pick)")
	flag.StringVar(&opts.Arg_exclude, "exclude", "", "Comma-separated list of targets (IPs/hostnames/CIDRs/octet-ranges) to exclude from the scan")
	flag.StringVar(&opts.Arg_excludefile, "excludefile", "", "Path to a file with targets (IPs/hostnames/CIDRs/octet-ranges) to exclude, separated by newlines, spaces, or tabs; '#' starts an end-of-line comment")

	// Custom usage: long flags (name length > 2) are shown with a leading '--',
	// short flags keep '-'. Both prefixes are accepted at parse time — this is
	// purely a documentation convenience so '--src-ip', '--exclude' and
	// '--excludefile' read with a double-dash long-option style.
	flag.Usage = func() {
		out := flag.CommandLine.Output()
		fmt.Fprintf(out, "Usage of %s:\n", os.Args[0])
		flag.VisitAll(func(f *flag.Flag) {
			prefix := "-"
			if len(f.Name) > 2 {
				prefix = "--"
			}
			name, usage := flag.UnquoteUsage(f)
			line := fmt.Sprintf("  %s%s", prefix, f.Name)
			if name != "" {
				line += " " + name
			}
			fmt.Fprintln(out, line)
			desc := strings.ReplaceAll(usage, "\n", "\n    \t")
			// Append "(default ...)" for non-zero defaults, mirroring Go's
			// stdlib PrintDefaults. Strings are quoted, everything else bare.
			if !isZeroValue(f, f.DefValue) {
				if g, ok := f.Value.(flag.Getter); ok {
					if _, isStr := g.Get().(string); isStr {
						desc += fmt.Sprintf(" (default %q)", f.DefValue)
					} else {
						desc += fmt.Sprintf(" (default %s)", f.DefValue)
					}
				} else {
					desc += fmt.Sprintf(" (default %s)", f.DefValue)
				}
			}
			fmt.Fprintln(out, "    \t"+desc)
		})
	}

	flag.Parse()

	return opts
}

// isZeroValue mirrors the stdlib helper of the same name in flag/flag.go:
// it returns true if `value` is the zero value for the flag's underlying type.
// Used by the custom Usage above so we only annotate genuinely defaulted flags.
func isZeroValue(f *flag.Flag, value string) bool {
	// Best-effort: most flag types' zero stringification matches one of these.
	switch value {
	case "", "0", "false":
		return true
	}
	return false
}
