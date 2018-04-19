package main

import (
	"flag"
	"sync"

	"github.com/ailidani/paxi"
	"github.com/ailidani/paxi/atomic"
	"github.com/ailidani/paxi/blockchain"
	"github.com/ailidani/paxi/dynamo"
	"github.com/ailidani/paxi/epaxos"
	"github.com/ailidani/paxi/kpaxos"
	"github.com/ailidani/paxi/log"
	"github.com/ailidani/paxi/paxos"
	"github.com/ailidani/paxi/paxos_group"
	"github.com/ailidani/paxi/ppaxos"
	"github.com/ailidani/paxi/vpaxos"
	"github.com/ailidani/paxi/wankeeper"
	"github.com/ailidani/paxi/wpaxos"
)

var id = flag.String("id", "", "ID in format of Zone.Node.")
var simulation = flag.Bool("sim", false, "simulation mode")

var master = flag.String("master", "", "Master address.")

func replica(id paxi.ID) {
	if *master != "" {
		paxi.ConnectToMaster(*master, false, id)
	}

	log.Infof("node %v starting...", id)

	switch paxi.GetConfig().Algorithm {

	case "paxos":
		paxos.NewReplica(id).Run()

	case "vpaxos":
		vpaxos.NewReplica(id).Run()

	case "wpaxos":
		wpaxos.NewReplica(id).Run()

	case "wankeeper":
		wankeeper.NewReplica(id).Run()

	case "epaxos":
		epaxos.NewReplica(id).Run()

	case "kpaxos":
		kpaxos.NewReplica(id).Run()

	case "paxos_groups":
		paxos_group.NewReplica(id).Run()

	case "atomic":
		atomic.NewReplica(id).Run()

	case "ppaxos":
		ppaxos.NewReplica(id).Run()

	case "dynamo":
		dynamo.NewReplica(id).Run()

	case "blockchain":
		blockchain.NewMiner(id).Run()

	default:
		panic("Unknown algorithm.")
	}
}

func main() {
	paxi.Init()

	if *simulation {
		var wg sync.WaitGroup
		wg.Add(1)
		paxi.Simulation()
		for id := range paxi.GetConfig().Addrs {
			n := id
			go replica(n)
		}
		wg.Wait()
	} else {
		replica(paxi.ID(*id))
	}
}
