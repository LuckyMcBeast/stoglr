package main

import (
	"flag"
	"fmt"
	"github.com/LuckyMcBeast/stoglr/server"
	"github.com/LuckyMcBeast/stoglr/server/datastore"
)

var (
	port         *string
	help         *bool
	version      *bool
	versionValue string
	buildTime    string
)

func init() {
	port = flag.String("p", "8080", "The port to listen on")
	help = flag.Bool("help", false, "Show this help (exclusive with other options)")
	version = flag.Bool("v", false, "Show version (exclusive with other options)")
	flag.Usage = func() {
		fmt.Println("Usage: stoglr [options]")
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}
	if *version {
		fmt.Println("stoglr version", versionValue)
		fmt.Println("build", buildTime)
		return
	}
	printStartText()
	server.NewToggleServer(
		*port,
		datastore.NewRuntimeDatastore()).Start()
}

func printStartText() {
	print(`
             .XXx.                            x$$;         
             .$$X                             +$$;         
 .$$$$$$$$. +$$$$$+  ;$$$$$$$X.   ;$$$$$+$$+  x$$; .X$$;$$X
.X$$.  .$$X  .$$X   +$$X   :$$$  +$$x  .$$$+  +$$;  X$$X.  
 ;$$$$$X+.   .$$X  .$$$     x$$:.X$$.   :$$+  x$$;  X$$.   
   ..x$$$$$  .$$X. .$$$     x$$: $$$.   ;$$+  +$$; .X$$.   
.$$X    X$$. .$$X   ;$$X.  ;$$$  +$$X. .$$$+  x$$;  X$$.   
 :$$$$$$$$.   X$$$+  :$$$$$$$+    ;$$$$$x$$+  +$$;  X$$.   
                                        ;$$;               
                                 :$$$XX$$$X
`)
	println("The           Simple           Feature           Toggler")
	println()
}
