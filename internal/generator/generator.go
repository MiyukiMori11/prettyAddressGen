package generator

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/MiyukiMori11/prettyAddressGen/internal/workerPool"
)

type AddrCreator interface {
	Create() (address string, publicKey string, privateKey string)
}

type generator struct {
	network AddrCreator

	workerPool workerPool.WorkerPool
}

func New(network AddrCreator, wkPool workerPool.WorkerPool) (*generator, error) {
	if network == nil {
		return nil, errors.New("network can't be nil")
	}

	return &generator{
		network:    network,
		workerPool: wkPool,
	}, nil
}

func (g *generator) Generate(ctx context.Context, pattern string, workersNum, resultsNum int64) ([]addressInfo, error) {
	if resultsNum < 1 {
		return nil, errors.New("number of results must be > 0")
	}

	ctx, cancelFunc := context.WithCancel(ctx)
	defer cancelFunc()

	result := make([]addressInfo, 0, resultsNum)

	genFunc := func() interface{} {
		ctx, cancelFunc := context.WithCancel(ctx)
		defer cancelFunc()
		re := regexp.MustCompile(pattern)
	loop:
		for {
			select {
			case <-ctx.Done():
				break loop
			default:
				if addr, pub, private := g.network.Create(); re.MatchString(addr) {
					return addressInfo{address: addr, publicKey: pub, privateKey: private}
				}
			}

		}
		return nil
	}

	in := make(chan workerPool.TaskFunc)
	var resultCount int

	out, err := g.workerPool.Run(in, workersNum)
	if err != nil {
		return nil, fmt.Errorf("can't run worker pool: %w", err)
	}

	go func() {
	loop:
		for {
			select {
			case <-ctx.Done():
				close(in)
				break loop

			default:
				in <- genFunc
			}

		}

	}()

	for v := range out {
		if r, ok := v.(addressInfo); ok && resultCount < int(resultsNum) {
			result = append(result, r)
			resultCount++

		} else {
			cancelFunc()
		}

	}

	return result, nil
}
