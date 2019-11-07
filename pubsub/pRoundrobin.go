package pubsub

import "time"

// PubRoundRobin only sends events to 1 subscriber on a round robin basis.
type PubRoundRobin struct {
	*PubThreaded
	next int
}

func NewPubRoundRobin(buffer int) *PubRoundRobin {
	p := new(PubRoundRobin)
	p.PubThreaded = NewPubThreaded(buffer)
	p.PubThreaded.ChangeWriteHandle(p.write)

	return p
}

func (p *PubRoundRobin) write(o interface{}) {
	for len(p.Subscribers) == 0 {
		// TODO: This isn't the best way to handle this.
		// 		Someone can unsub after we exit this for too.
		time.Sleep(100 * time.Millisecond)
	}
	p.Subscribers[p.next%len(p.Subscribers)].write(o)
	p.next++
}

func (p *PubRoundRobin) Start() {
	p.PubThreaded.Start()
}

func (p *PubRoundRobin) Publish(path string) *PubRoundRobin {
	globalPublish(path, p)
	return p
}
