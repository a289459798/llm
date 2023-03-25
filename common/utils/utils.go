package utils

import (
	"fmt"
	"os"
	"sync/atomic"
	"time"
)

func GenerateOrderNo() string {
	var num int64
	s := time.Now().Format("20060102150405")
	m := time.Now().UnixNano()/1e6 - time.Now().UnixNano()/1e9*1e3
	ms := sup(m, 3)
	p := os.Getpid() % 1000
	ps := sup(int64(p), 3)
	i := atomic.AddInt64(&num, 1)
	r := i % 10000
	rs := sup(r, 4)
	n := fmt.Sprintf("%s%s%s%s", s, ms, ps, rs)
	return n
}

func sup(i int64, n int) string {
	m := fmt.Sprintf("%d", i)
	for len(m) < n {
		m = fmt.Sprintf("0%s", m)
	}
	return m
}
