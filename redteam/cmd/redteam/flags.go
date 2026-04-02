package main

import "strings"

// stringSlice implements flag.Value for repeatable -attacks-dir flags.
type stringSlice []string

func (s *stringSlice) String() string { return strings.Join(*s, ",") }

func (s *stringSlice) Set(v string) error {
	*s = append(*s, v)
	return nil
}
