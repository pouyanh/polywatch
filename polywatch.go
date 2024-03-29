package polywatch

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"

	"github.com/radovskyb/watcher"
	"github.com/zmwangx/debounce"

	"github.com/pouyanh/polywatch/config"
)

var (
	ErrUnsupportedFilter = errors.New("filter not supported")
)

func Start() error {
	ctx, stop := context.WithCancel(context.Background())
	defer stop()

	cfg := config.MustLoad()
	wg := sync.WaitGroup{}
	for _, cw := range cfg.Watchers {
		w, err := newPolyWatcher(cw)
		if err != nil {
			return err
		}

		wg.Add(1)
		go func() {
			defer wg.Done()

			if err := w.watch(ctx); err != nil {
				stop()
			}
		}()
	}

	bindSignals(stop, syscall.SIGTERM, syscall.SIGINT)

	wg.Wait()

	return nil
}

func bindSignals(fn func(), ss ...os.Signal) {
	ntfy := make(chan os.Signal, 1)
	signal.Notify(ntfy, ss...)

	go func() {
		<-ntfy
		fn()
	}()
}

type polyWatcher struct {
	cfg config.Watcher

	w   *watcher.Watcher
	lg  *log.Logger
	cmd *exec.Cmd
}

func newPolyWatcher(cfg config.Watcher) (*polyWatcher, error) {
	lg := log.New(os.Stderr, fmt.Sprintf("poly-watcher[%s]: ", cfg.Name), log.LstdFlags)

	w := watcher.New()
	for _, wf := range cfg.Watch.Files {
		path := filepath.Clean(os.ExpandEnv(wf.Path))

		var err error
		if wf.Recursive {
			err = w.AddRecursive(path)
		} else {
			err = w.Add(path)
		}

		if err != nil {
			return nil, err
		}
	}

	for _, wf := range cfg.Watch.Filters {
		switch wf.On {
		case config.WatchFilterScopeFilename:
			switch wf.Type {
			case config.WatchFilterTypeRegex:
				w.AddFilterHook(fileFilterRegex(wf.Include, wf.List...))

			case config.WatchFilterTypeList:
				w.AddFilterHook(fileFilterList(wf.Include, wf.List...))

			default:
				return nil, ErrUnsupportedFilter
			}

		default:
			// todo: event filters
			lg.Printf(`filter scope "%s" not supported yet`, wf.On)
		}
	}

	pw := &polyWatcher{
		cfg: cfg,

		w:  w,
		lg: lg,
	}

	pw.renewCommand()

	return pw, nil
}

func (pw *polyWatcher) renewCommand() {
	// todo: support multiline command

	rawcmd := strings.Split(pw.cfg.Command.Shell, " ")
	rawcmd = append(rawcmd, pw.cfg.Command.Exec)

	cmd := exec.Command(rawcmd[0], rawcmd[1:]...)
	cmd.Env = pw.cfg.Command.Env
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	pw.cmd = cmd
}

func (pw *polyWatcher) watch(ctx context.Context) error {
	for path, f := range pw.w.WatchedFiles() {
		pw.lg.Printf("%s: %s\n", path, f.Name())
	}

	chErr := make(chan error)
	defer close(chErr)

	uh := pw.updateHandler()
	go func() {
		for {
			select {
			case e := <-pw.w.Event:
				pw.lg.Printf("event received: %+v\n", e)
				err := uh(ctx, e)
				if err != nil {
					pw.lg.Printf("error occurred during handling update: %s\n", err)
				}
			case err := <-pw.w.Error:
				chErr <- err
			case <-pw.w.Closed:
				return
			}
		}
	}()

	go func() {
		go func() {
			pw.w.Wait()

			pw.lg.Println("started")
		}()

		pw.lg.Println("starting...")
		chErr <- pw.w.Start(pw.cfg.Watch.Interval)
	}()

	defer pw.w.Close()
	defer func() { _ = pw.kill(ctx) }()
	select {
	case err := <-chErr:
		pw.lg.Printf("error occurred during watch: %s", err)
		return err

	case <-ctx.Done():
		pw.lg.Println("stopping...")
	}

	return nil
}

type updateHandler func(ctx context.Context, event watcher.Event) error

func (pw *polyWatcher) updateHandler() updateHandler {
	var uh func(uu ...update) error
	switch pw.cfg.RateLimit.Strategy {
	case config.RateLimitStrategyDebounce:
		uh, _ = debounce.DebounceWithCustomSignature(pw.handleUpdate, pw.cfg.RateLimit.Wait)

	case config.RateLimitStrategyThrottle:
		uh, _ = debounce.ThrottleWithCustomSignature(pw.handleUpdate, pw.cfg.RateLimit.Wait)

	//case config.RateLimitStrategyAudit:
	//case config.RateLimitStrategySample:

	default:
		uh = pw.handleUpdate
	}

	return func(ctx context.Context, event watcher.Event) error {
		return uh(update{
			ctx:   ctx,
			event: event,
		})
	}
}

type update struct {
	ctx   context.Context
	event watcher.Event
}

func (pw *polyWatcher) handleUpdate(uu ...update) error {
	u := uu[0]

	return pw._handleUpdate(u.ctx, u.event)
}

func (pw *polyWatcher) _handleUpdate(ctx context.Context, event watcher.Event) error {
	pw.lg.Println("updating...")

	err := pw.kill(ctx)
	if err != nil {
		pw.lg.Printf("unable to kill previous command: %s", err)
	}

	return pw.cmd.Start()
}

func (pw *polyWatcher) kill(ctx context.Context) error {
	if pw.cmd.Process == nil {
		return nil
	}

	if pw.cmd.Process.Pid < 1 {
		return nil
	}

	pw.lg.Printf("killing previous command pid(%d) with sig(%s)", pw.cmd.Process.Pid, pw.cfg.Kill.Signal)

	// todo: apply kill timeout
	group, err := os.FindProcess(-1 * pw.cmd.Process.Pid)
	if err == nil {
		err = group.Signal(pw.cfg.Kill.Signal)
		if err != nil {
			pw.lg.Printf("unable to send signal to previous command: %s", err)
		}
	}

	defer pw.renewCommand()

	return handleWaitError(pw.cmd.Wait(), pw.cfg.Kill.Signal)
}

func handleWaitError(err error, sig os.Signal) error {
	if e, ok := err.(*exec.ExitError); ok {
		status := e.ProcessState.Sys().(syscall.WaitStatus)
		if status.Signaled() {
			// TODO: seperate windows and posix functionality and compare relevant types
			if status.Signal().String() == sig.String() {
				return nil
			}
		}
	}

	return err
}
