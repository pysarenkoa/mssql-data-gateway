// +build windows

package main

import (
	"fmt"
	"os"
	"strings"
)

func usage(errmsg string) {
	fmt.Fprintf(os.Stderr,
		"%s\n\n"+
			"usage: %s <command>\n"+
			"       where <command> is one of\n"+
			"       install, remove, debug, start, stop, pause or continue.\n",
		errmsg, os.Args[0])
	os.Exit(2)
}

// https://github.com/golang/sys/tree/master/windows/svc/example
// https://stackoverflow.com/questions/42774467/how-to-convert-sql-rows-to-typed-json-in-golang
// https://docs.microsoft.com/ru-ru/azure/azure-sql/database/connect-query-go
func main() {
	if len(os.Args) < 2 {
		usage("no command specified")
		return
	}

	cmd := strings.ToLower(os.Args[1])
	execCmd(cmd)
}
