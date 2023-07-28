# PolyWatch
File change watcher which gets triggerred by events like file addition, contents update, ownership change or removal and
run commands subsequently.
It is designed to help in setting up a development environment for some golang services which are all run using live code.
By using PolyWatch source code can get recompiled & started automatically on each service code file changes.
See [hotreload](examples/hotreload) for a working example.

# Contents
* [Features](#features)
* [Installation](#installation)
	* [Go](#go)
	* [AUR (Arch User Repository)](#aur-arch-user-repository)
* [Usage](#usage)
* [Configuration](#configuration)
* [Contributors](#contributors)
* [Todo](#todo)
* [Related projects](#related-projects)
* [License](#license)

# Features
* Configurable using a single config file; Supports JSON, TOML, YAML, HCL, INI files
* Watch multiple directories recursively or non-recursively
* Concurrent watchers which run independently having their own settings & command
* Inclusive & Exclusive file group **filters** using _regular expressions_ or list
* Rate limit using different strategies like _debounce_ and _throttle_
* Configurable kill **signal**; In fact running command can do a graceful shutdown, restart or reload due to the signal

# Installation
## Go
If you're a Gopher and want to install PolyWatch locally in your GOPATH:

```shell
go install -v github.com/pouyanh/polywatch/cmd/polywatch@latest
```

## AUR (Arch User Repository)
If you're using [Arch Linux][archlinux] install [PolyWatch][aur-polywatch] package from [AUR][wiki-aur]

```shell
yay -S polywatch
```

# Usage
Create the config file, run `polywatch` & it runs the command(s) on start and then whenever changes happen

# Configuration
Configuration is done using a single file named `.pw` with these extensions: `json`, `toml`, `yml`, `yaml`, `hcl` & `ini`.
It have to be located in **current working directory**.
[Example](pw.example.yml) yaml config file can be a proper starting point.

Configuration consists of multiple independent watchers:
```yaml
watchers:
	- name: "watcher 1"
	- name: "watcher 2"
	- name: "watcher 3"
```

Each watcher have 5 configuration sections:
## Name Config
Name is a single string field. It's just a label for the watcher
```yaml
name: "watcher name"
```

## Watch Config
Watch as it's expected contains files & directories watching settings
* method: Defines watching mechanism. Currently just supports `polling` method that watches for file changes
in fixed intervals
* interval: When method is `polling`, it sets interval between each watch.
* files: Array of [WatchFile](#watchfile). Matching files get appended together; they get combined by logical OR.
* filters: Array of [WatchFilter](#watchfilter). Each candidate file have to pass all filters' tests in order to
make notification.

### WatchFile

### WatchFilter

## RateLimit Config
## Kill Config
## Command Config

# Contributors
Thanks to [Saman Koushki][gh-saman3d] for enabling multiline commands & improving process management.

# Todo
1. Implement fsnotify watch method
2. Support event filters like filters on operation scope
3. Consider kill timeout
4. Add wildcard filter type

# Related projects
* [fswatch][fswatch]: Command line tool to watch file changes using fsnotify
* PolyWatch uses [watcher][watcher] which is a library that can watch file changes by polling mechanism, but it's not configurable
* [fsnotify][fsnotify]: A cross-platform library to work with filesystem notifications

# License
This software is [licensed](LICENSE) under the [GPL v3 License][gpl]. © 2023 [Janstun][janstun]

[archlinux]: https://www.archlinux.org/
[aur-polywatch]: https://aur.archlinux.org/packages/polywatch
[wiki-aur]: https://wiki.archlinux.org/index.php/AUR
[fswatch]: https://github.com/codeskyblue/fswatch
[watcher]: https://github.com/radovskyb/watcher
[fsnotify]: https://github.com/fsnotify/fsnotify
[gpl]: https://www.gnu.org/licenses/gpl-3.0.en.html
[janstun]: http://janstun.com
[gh-saman3d]: https://github.com/saman3d
