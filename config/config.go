package config

import (
	"errors"
	"os"
	"syscall"
	"time"
)

var (
	ErrConfiguratorAlreadyRegistered = errors.New("a configurator has been already registered")
	ErrNoRegisteredConfigurator      = errors.New("no configurator has been registered")
)

const (
	DefaultWatchMethod   = WatchMethodPolling
	DefaultWatchInterval = 100 * time.Millisecond

	DefaultWatchFileRecursive bool = true

	DefaultWatchFilterScope        = WatchFilterScopeFilename
	DefaultWatchFilterInclude bool = true
	DefaultWatchFilterType         = WatchFilterTypeRegex

	DefaultRateLimitStrategy               = RateLimitStrategyNone
	DefaultRateLimitWait     time.Duration = 0

	DefaultKillSignal                = syscall.SIGTERM
	DefaultKillTimeout time.Duration = 0
)

var (
	DefaultWatcher = Watcher{
		Name:      "",
		Watch:     DefaultWatch,
		RateLimit: DefaultRateLimit,
		Kill:      DefaultKill,
		Command:   DefaultCommand,
	}

	DefaultCommand = Command{
		Shell: "/bin/sh -c",
		Path:  ".",
		Env:   os.Environ(),
		Exec:  "",
	}

	DefaultWatch = Watch{
		Method:   DefaultWatchMethod,
		Interval: DefaultWatchInterval,
		Files:    nil,
		Filters:  nil,
	}

	DefaultWatchFile = WatchFile{
		Path:      "",
		Recursive: DefaultWatchFileRecursive,
	}

	DefaultWatchFilter = WatchFilter{
		On:      DefaultWatchFilterScope,
		Include: DefaultWatchFilterInclude,
		Type:    DefaultWatchFilterType,
		List:    nil,
	}

	DefaultRateLimit = RateLimit{
		Strategy: DefaultRateLimitStrategy,
		Wait:     DefaultRateLimitWait,
	}

	DefaultKill = Kill{
		Signal:  DefaultKillSignal,
		Timeout: DefaultKillTimeout,
	}
)

type Config struct {
	Watchers []Watcher `json:"watchers"`
}

type Watcher struct {
	Name      string    `json:"name"`
	Watch     Watch     `json:"watch"`
	RateLimit RateLimit `json:"rateLimit"`
	Kill      Kill      `json:"kill"`
	Command   Command   `json:"cmd"`
}

type Command struct {
	Shell string   `json:"shell"`
	Exec  string   `json:"exec"`
	Path  string   `json:"path"`
	Env   []string `json:"env"`
}

type Watch struct {
	Method   WatchMethod   `json:"method"`
	Interval time.Duration `json:"interval"`
	Files    []WatchFile   `json:"files"`
	Filters  []WatchFilter `json:"filters"`
}

type WatchMethod string

const (
	WatchMethodPolling  WatchMethod = "polling"
	WatchMethodFsnotify WatchMethod = "fsnotify"
)

type WatchFile struct {
	Path      string `json:"path"`
	Recursive bool   `json:"recursive"`
}

type WatchFilter struct {
	On      WatchFilterScope `json:"on"`
	Include bool             `json:"include"`
	Type    WatchFilterType  `json:"type"`
	List    []string         `json:"list"`
}

type WatchFilterScope string

const (
	WatchFilterScopeFilename  WatchFilterScope = "filename"
	WatchFilterScopeOperation WatchFilterScope = "operation"
)

type WatchFilterType string

const (
	WatchFilterTypeRegex WatchFilterType = "regex"
	WatchFilterTypeList  WatchFilterType = "list"
)

type RateLimit struct {
	Strategy RateLimitStrategy `json:"strategy"`
	Wait     time.Duration     `json:"wait"`
}

type RateLimitStrategy string

const (
	RateLimitStrategyNone     RateLimitStrategy = "none"
	RateLimitStrategyDebounce RateLimitStrategy = "debounce"
	RateLimitStrategyThrottle RateLimitStrategy = "throttle"
	RateLimitStrategyAudit    RateLimitStrategy = "audit"
	RateLimitStrategySample   RateLimitStrategy = "sample"
)

type Kill struct {
	Signal  os.Signal     `json:"signal"`
	Timeout time.Duration `json:"timeout"`
}

type Configurator interface {
	Load() (*Config, error)
}

var (
	configurator Configurator
	cfgCache     *Config
)

func Register(cfr Configurator) {
	if configurator != nil {
		panic(ErrConfiguratorAlreadyRegistered)
	}

	configurator = cfr
}

func Load() (*Config, error) {
	if configurator == nil {
		return nil, ErrNoRegisteredConfigurator
	}

	if cfgCache == nil {
		cfg, err := configurator.Load()
		if err != nil {
			return nil, err
		}

		cfgCache = cfg
	}

	return cfgCache, nil
}

func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}

	return cfg
}
