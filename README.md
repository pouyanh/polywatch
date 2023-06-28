# Poly Watch
File change watcher

# Contents
* [Features](#features)
* [Installation](#installation)
	* [Go](#go)
	* [AUR (Arch User Repository)](#aur-arch-user-repository)
* [Usage](#usage)
* [Configuration](#configuration)
* [Todo](#todo)
* [Related projects](#related-projects)
* [License](#license)

# Features

# Installation
## Go

```shell
go install -v github.com/pouyanh/polywatch/cmd/polywatch@latest
```

## AUR (Arch User Repository)
If you're using [Arch Linux][archlinux] install [polywatch][aur-polywatch] package from [AUR][aur]

```shell
yay -S polywatch
```

# Usage

# Configuration

# Todo
1. Support multiline commands
2. Support event filters like filters on operation scope
3. Consider kill timeout

# Related projects
* [fswatch][fswatch]
* [watcher][watcher]
* [fsnotify][fsnotify]

# License
This software is [licensed](LICENSE) under the [GPL v3 License][gpl]. © 2023 [Janstun][janstun]

[archlinux]: https://www.archlinux.org/
[aur-polywatch]: https://aur.archlinux.org/packages/polywatch
[aur]: https://wiki.archlinux.org/index.php/AUR
[fswatch]: https://github.com/codeskyblue/fswatch
[watcher]: https://github.com/radovskyb/watcher
[fsnotify]: https://github.com/fsnotify/fsnotify
[gpl]: https://www.gnu.org/licenses/gpl-3.0.en.html
[janstun]: http://janstun.com
