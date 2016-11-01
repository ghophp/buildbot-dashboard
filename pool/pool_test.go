package pool

import (
	"fmt"
	"testing"

	gc "gopkg.in/check.v1"
)

var _ = gc.Suite(&PoolSuite{})

type PoolSuite struct{}

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { gc.TestingT(t) }

func (s *PoolSuite) TestRequestPool(c *gc.C) {
	requestPool := NewRequestPool()

	listener1 := make(chan string)
	listener2 := make(chan string)
	listener3 := make(chan string)

	requestPool.Fetch("http://www.blueprintcss.org/tests/parts/sample.html", listener1)
	requestPool.Fetch("http://www.blueprintcss.org/tests/parts/sample.html", listener2)
	requestPool.Fetch("http://www.blueprintcss.org/tests/parts/sample.html", listener3)

	select {
	case resp := <-listener1:
		fmt.Println(fmt.Sprintf("resp 1 size %d", len(resp)))
	}

	select {
	case resp := <-listener2:
		fmt.Println(fmt.Sprintf("resp 1 size %d", len(resp)))
	}

	select {
	case resp := <-listener3:
		fmt.Println(fmt.Sprintf("resp 1 size %d", len(resp)))
	}
}
