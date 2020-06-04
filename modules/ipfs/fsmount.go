package ipfs

import (
	"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/modules/ipfs/entryfs"
	"github.com/FactomProject/factomd/modules/worker"
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
