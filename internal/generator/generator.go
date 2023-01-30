package generator

import (
	"errors"
	"regexp"

	"github.com/MiyukiMori11/prettyAddressGen/internal/workerPool"
)

type addrCreator interface {
	Create() (address string, publicKey string, privateKey string)
}

type addressInfo struct {
	address    string
	publicKey  string
	privateKey string
}

type generator struct {
	network addrCreator

	workerPool workerPool.WorkerPool
}

func New(network addrCreator, wkPool workerPool.WorkerPool) (*generator, error) {
	if network == nil {
		return nil, errors.New("network can't be nil")
	}

	return &generator{
		network:    network,
		workerPool: wkPool,
	}, nil
}

func (g *generator) Generate(pattern string, workersNum, resultsNum int64) ([]addressInfo, error) {
	if resultsNum < 1 {
		return nil, errors.New("number of results must be > 0")
	}

	result := make([]addressInfo, 0, resultsNum)

	g.workerPool.AddWorkers(workersNum)

	genFunc := func() interface{} {
		re := regexp.MustCompile(pattern)
		for {
			if addr, pub, private := g.network.Create(); re.MatchString(addr) {
				return addressInfo{address: addr, publicKey: pub, privateKey: private}
			}
		}
	}

	in, out := g.workerPool.RunAll()
	var resultCount int

	go func() {
		for {
			in <- genFunc
		}

	}()

	for v := range out {
		if r, ok := v.(addressInfo); ok && resultCount < int(resultsNum) {
			result = append(result, r)
			resultCount++

		} else {
			g.workerPool.CloseAll()
			break
		}

	}

	return result, nil
}
