package shutdown

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/kklab-com/goth-kklogger"
)

var once = sync.Once{}
var sig = make(chan os.Signal, 1)
var regs []func(sig os.Signal)
var opLock sync.Mutex

func _Init() {
	once.Do(func() {
		signal.Notify(sig, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGHUP)
		go func() {
			s := <-sig
			opLock.Lock()
			defer opLock.Unlock()
			msg := fmt.Sprintf("SIGNAL: %s, SHUTDOWN CATCH", s.String())
			kklogger.InfoJ("Shutdown", msg)
			println(msg)
			for i, f := range regs {
				f(s)
				kklogger.InfoJ("Shutdown", fmt.Sprintf("Task %d done.", i))
			}

			kklogger.InfoJ("Shutdown", "DONE")
			println("DONE")
		}()
	})
}

func InvokeLast(f func(sig os.Signal)) {
	_Init()
	opLock.Lock()
	defer opLock.Unlock()
	regs = append(regs, f)
}

func InvokeFirst(f func(sig os.Signal)) {
	_Init()
	opLock.Lock()
	defer opLock.Unlock()
	regs = append([]func(sig os.Signal){f}, regs...)
}
