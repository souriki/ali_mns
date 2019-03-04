package ali_mns

import (
	"testing"
)

func TestCheckQPS(t *testing.T) {
	qm := NewQPSMonitor(3, 10)
	for {
		qm.Pulse()
		if qps := qm.QPS(); qps > 10 {
			break
		}
	}

	qm.checkQPS()
}
