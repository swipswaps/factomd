package ipfs

import (
	"fmt"
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/modules/ipfs/entryfs"
	"github.com/FactomProject/factomd/modules/worker"
	shell "github.com/ipfs/go-ipfs-api"
	"os"
	"strings"
)

// KLUDGE: fs mount location is hardcoded
var mountPath = "/tmp/fct"

// start Filesystem mount
func Start(w *worker.Thread, s interfaces.IState) {

	// KLUDGE: path is hardcoded
	server, err := entryfs.Mount(&entryfs.Storage{DB: s.GetDB(), Path: mountPath})
	if err != nil {
		panic(err) // FIXME
	}

	w.Spawn("entryFS", func(w *worker.Thread) {
		w.OnRun(server.Serve) // REVIEW: will this shutdown properly?
		w.OnExit(func() {
			server.Unmount()
		})
	})
}

func _main() {
	// Where your local node is running on localhost:5001
	sh := shell.NewShell("localhost:5001")
	shell.AddOpts()
	sh.Add()
	cid, err := sh.Add(strings.NewReader("hello world!"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %s", err)
		os.Exit(1)
	}
	fmt.Printf("added %s", cid)
}
