package main

import (
	"stoglr/server"
	"stoglr/server/datastore"
)

func main() {
	printStartText()
	server.NewToggleServer(
		"8080",
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
