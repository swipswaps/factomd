package entryfs

/*
dstore := datastore.NewMapDatastore()
bstore := blockstore.NewBlockstore(dstore)
bserv := blockservice.New(bstore, offline.Exchange(bstore))
dserv := merkledag.NewDAGService(bserv)

n := &core.IpfsNode{
	Blockstore: blockstore.NewGCBlockstore(bstore, blockstore.NewGCLocker()),
	Pinning:    pin.NewPinner(dstore, dserv, dserv),
	DAG:        dserv,
}

r, err := os.Open(filename)
if err != nil {
	return err
}

// key is the hash from ipfs add.
key, err := coreunix.Add(n, r)
if err != nil {
	return err
}

*/
