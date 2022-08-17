package main

import "flag"

type Flags struct {
	Address *string
	Port    *int
	Path    *string

	disableAuth *bool
	user        *string
	pass        *string

	Version *bool
}

func (f *Flags) init() {
	f.Address = flag.String("addr", "", "Address.")
	f.Port = flag.Int("port", 9111, "Port.")
	f.Path = flag.String("path", "/metrics", "Path.")

	f.disableAuth = flag.Bool("disable-auth", false, "Disable basic authentication.")
	f.user = flag.String("user", "admin", "Username.")
	f.pass = flag.String("pass", "admin", "Password.")

	f.Version = flag.Bool("version", false, "Show version.")

	flag.Parse()
}
