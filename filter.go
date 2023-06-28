package polywatch

import (
	"os"
	"regexp"

	"github.com/radovskyb/watcher"
)

func fileFilterRegex(include bool, patterns ...string) watcher.FilterFileHookFunc {
	rr := make([]*regexp.Regexp, len(patterns))
	for k, pattern := range patterns {
		rr[k] = regexp.MustCompile(pattern)
	}

	yes, no := includeInFileFilter(include)

	return func(info os.FileInfo, _ string) error {
		filename := info.Name()

		// Match
		for _, r := range rr {
			if r.MatchString(filename) {
				return yes
			}
		}

		// No match.
		return no
	}
}

func fileFilterList(include bool, list ...string) watcher.FilterFileHookFunc {
	files := make(map[string]bool)
	for _, name := range list {
		files[name] = true
	}

	yes, no := includeInFileFilter(include)

	return func(info os.FileInfo, _ string) error {
		filename := info.Name()

		if _, ok := files[filename]; ok {
			return yes
		}

		return no
	}
}

func includeInFileFilter(include bool) (yes error, no error) {
	yes, no = nil, watcher.ErrSkip
	if !include {
		yes, no = no, yes
	}

	return
}
