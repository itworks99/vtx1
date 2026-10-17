// Command vtxsim is the VTX1 ISA-level simulator. It executes a memory image
// produced by the assembler (vtxasm -f memh) and reports the dynamic metrics
// used for architecture evaluation.
//
//	vtxsim [-entry 0x1000] [-max N] [-dump 0x2000:16] [-trace] [-quiet] program.memh
package main

import (
	"flag"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/kvany/vtx1/assembler/sim"
)

func parseNum(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		return strconv.ParseInt(s[2:], 16, 64)
	}
	return strconv.ParseInt(s, 10, 64)
}

func main() {
	entry := flag.String("entry", "", "entry point word address (default: first loaded address)")
	maxSteps := flag.Uint64("max", 10_000_000, "maximum instructions to execute")
	dump := flag.String("dump", "", "after the run, dump a memory region: ADDR:LEN (e.g. 0x2000:16)")
	trace := flag.Bool("trace", false, "print each executed instruction")
	quiet := flag.Bool("quiet", false, "suppress the statistics report")
	flag.Parse()

	if flag.NArg() != 1 {
		fmt.Fprintln(os.Stderr, "usage: vtxsim [options] program.memh")
		flag.PrintDefaults()
		os.Exit(1)
	}

	f, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
	defer f.Close()

	m := sim.New()
	first, err := m.LoadMemh(f)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error loading image:", err)
		os.Exit(1)
	}
	m.PC = first
	if *entry != "" {
		e, err := parseNum(*entry)
		if err != nil {
			fmt.Fprintln(os.Stderr, "bad -entry:", err)
			os.Exit(1)
		}
		m.PC = uint32(e)
	}
	if *trace {
		m.Trace = os.Stderr
	}

	runErr := m.Run(*maxSteps)
	if runErr != nil {
		fmt.Fprintln(os.Stderr, "trap:", runErr)
	}

	if !*quiet {
		s := m.Stats
		fmt.Printf("instructions: %d\n", s.Instructions)
		cats := make([]string, 0, len(s.ByCategory))
		for c := range s.ByCategory {
			cats = append(cats, c)
		}
		sort.Strings(cats)
		for _, c := range cats {
			fmt.Printf("  %-15s %d\n", c+":", s.ByCategory[c])
		}
		fmt.Printf("loads: %d  stores: %d\n", s.Loads, s.Stores)
		fmt.Printf("branches: %d (taken: %d)\n", s.Branches, s.BranchesTaken)
	}

	if *dump != "" {
		parts := strings.SplitN(*dump, ":", 2)
		if len(parts) != 2 {
			fmt.Fprintln(os.Stderr, "bad -dump, want ADDR:LEN")
			os.Exit(1)
		}
		addr, err1 := parseNum(parts[0])
		length, err2 := parseNum(parts[1])
		if err1 != nil || err2 != nil {
			fmt.Fprintln(os.Stderr, "bad -dump, want ADDR:LEN")
			os.Exit(1)
		}
		var ascii strings.Builder
		fmt.Printf("memory at 0x%X:\n", addr)
		for i := int64(0); i < length; i++ {
			v := m.ReadWord(uint32(addr + i))
			fmt.Printf("  0x%04X: %d\n", addr+i, v)
			if v >= 32 && v < 127 {
				ascii.WriteByte(byte(v))
			} else if v == 0 {
				break
			} else {
				ascii.WriteByte('.')
			}
		}
		if ascii.Len() > 0 {
			fmt.Printf("  as text: %q\n", ascii.String())
		}
	}

	if runErr != nil {
		os.Exit(2)
	}
}
