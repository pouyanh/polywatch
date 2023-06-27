package viper

import (
	"strings"
	"syscall"
	"time"

	"github.com/spf13/viper"

	"github.com/pouyanh/polywatch/config"
)

func init() {
	config.Register(newConfigurator())
}

type configurator struct {
	*viper.Viper
}

func newConfigurator() *configurator {
	v := viper.New()
	v.SetConfigName(".pw")
	v.AddConfigPath(".")

	return &configurator{
		Viper: v,
	}
}

func (cfr configurator) Load() (*config.Config, error) {
	if err := cfr.ReadInConfig(); err != nil {
		return nil, err
	}

	cfg := new(Config)
	if err := cfr.Unmarshal(cfg); err != nil {
		return nil, err
	}

	out := cfg.decode()

	return &out, nil
}

type Config struct {
	Watchers []Watcher `mapstructure:"watchers"`
}

func (cfg Config) decode() config.Config {
	dst := config.Config{}
	for _, w := range cfg.Watchers {
		dst.Watchers = append(dst.Watchers, w.decode())
	}

	return dst
}

type Watcher struct {
	Name      string    `mapstructure:"name"`
	Watch     Watch     `mapstructure:"watch"`
	RateLimit RateLimit `mapstructure:"rateLimit"`
	Kill      Kill      `mapstructure:"kill"`
	Command   string    `mapstructure:"cmd"`
}

func (w Watcher) decode() config.Watcher {
	dst := config.DefaultWatcher
	dst.Name = override(w.Name, dst.Name, testStringZero)
	dst.Watch = w.Watch.decode()
	dst.RateLimit = w.RateLimit.decode()
	dst.Kill = w.Kill.decode()
	dst.Command = override(w.Command, dst.Command, testStringZero)

	return dst
}

type Watch struct {
	Method  config.WatchMethod `mapstructure:"method"`
	Files   []WatchFile        `mapstructure:"files"`
	Filters []WatchFilter      `mapstructure:"filters"`
}

func (w Watch) decode() config.Watch {
	dst := config.DefaultWatch
	dst.Method = config.WatchMethod(override(string(w.Method), string(dst.Method), testStringZero))
	for _, f := range w.Files {
		dst.Files = append(dst.Files, f.decode())
	}
	for _, f := range w.Filters {
		dst.Filters = append(dst.Filters, f.decode())
	}

	return dst
}

type WatchFile struct {
	Path      string `mapstructure:"path"`
	Recursive *bool  `mapstructure:"recursive"`
}

func (wf WatchFile) decode() config.WatchFile {
	dst := config.DefaultWatchFile
	dst.Path = override(wf.Path, dst.Path, testStringZero)
	dst.Recursive = *override(wf.Recursive, &dst.Recursive, testNil[bool])

	return dst
}

type WatchFilter struct {
	On      config.WatchFilterScope `mapstructure:"on"`
	Include *bool                   `mapstructure:"include"`
	Type    config.WatchFilterType  `mapstructure:"type"`
	List    []string                `mapstructure:"list"`
}

func (wf WatchFilter) decode() config.WatchFilter {
	dst := config.DefaultWatchFilter
	dst.On = config.WatchFilterScope(override(string(wf.On), string(dst.On), testStringZero))
	dst.Include = *override(wf.Include, &dst.Include, testNil[bool])
	dst.Type = config.WatchFilterType(override(string(wf.Type), string(dst.Type), testStringZero))
	dst.List = wf.List

	return dst
}

type RateLimit struct {
	Strategy config.RateLimitStrategy `mapstructure:"strategy"`
	Wait     *time.Duration           `mapstructure:"wait"`
}

func (rl RateLimit) decode() config.RateLimit {
	dst := config.DefaultRateLimit
	dst.Strategy = config.RateLimitStrategy(override(string(rl.Strategy), string(dst.Strategy), testStringZero))
	dst.Wait = *override(rl.Wait, &dst.Wait, testNil[time.Duration])

	return dst
}

type Kill struct {
	Signal  string         `mapstructure:"signal"`
	Timeout *time.Duration `mapstructure:"timeout"`
}

func (k Kill) decode() config.Kill {
	dst := config.DefaultKill
	dst.Signal = syscall.Signal(override(int(signalFromName(k.Signal)), int(dst.Signal.(syscall.Signal)), testIntZero))
	dst.Timeout = *override(k.Timeout, &dst.Timeout, testNil[time.Duration])

	return dst
}

func override[T any](v, def T, testZero func(v T) bool) T {
	if testZero(v) {
		return def
	}

	return v
}

func testStringZero(v string) bool {
	return len(strings.TrimSpace(v)) == 0
}

func testNil[T any](v *T) bool {
	return v == nil
}

func testIntZero(v int) bool {
	return v == 0
}

// generated using list of syscall.SIG* and replacing by: SIG(\w+)\s*=\s*.+$ ==> "$1": syscall.SIG$1,
var signals = map[string]syscall.Signal{
	"ABRT":   syscall.SIGABRT,
	"ALRM":   syscall.SIGALRM,
	"BUS":    syscall.SIGBUS,
	"CHLD":   syscall.SIGCHLD,
	"CLD":    syscall.SIGCLD,
	"CONT":   syscall.SIGCONT,
	"FPE":    syscall.SIGFPE,
	"HUP":    syscall.SIGHUP,
	"ILL":    syscall.SIGILL,
	"INT":    syscall.SIGINT,
	"IO":     syscall.SIGIO,
	"IOT":    syscall.SIGIOT,
	"KILL":   syscall.SIGKILL,
	"PIPE":   syscall.SIGPIPE,
	"POLL":   syscall.SIGPOLL,
	"PROF":   syscall.SIGPROF,
	"PWR":    syscall.SIGPWR,
	"QUIT":   syscall.SIGQUIT,
	"SEGV":   syscall.SIGSEGV,
	"STKFLT": syscall.SIGSTKFLT,
	"STOP":   syscall.SIGSTOP,
	"SYS":    syscall.SIGSYS,
	"TERM":   syscall.SIGTERM,
	"TRAP":   syscall.SIGTRAP,
	"TSTP":   syscall.SIGTSTP,
	"TTIN":   syscall.SIGTTIN,
	"TTOU":   syscall.SIGTTOU,
	"UNUSED": syscall.SIGUNUSED,
	"URG":    syscall.SIGURG,
	"USR1":   syscall.SIGUSR1,
	"USR2":   syscall.SIGUSR2,
	"VTALRM": syscall.SIGVTALRM,
	"WINCH":  syscall.SIGWINCH,
	"XCPU":   syscall.SIGXCPU,
	"XFSZ":   syscall.SIGXFSZ,
}

func signalFromName(name string) syscall.Signal {
	if s, ok := signals[strings.ToUpper(name)]; ok {
		return s
	}

	return syscall.Signal(0)
}
