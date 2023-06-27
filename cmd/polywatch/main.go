package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"regexp"
	"time"

	"github.com/radovskyb/watcher"
	"github.com/zmwangx/debounce"

	"github.com/pouyanh/polywatch/config"
	_ "github.com/pouyanh/polywatch/config/viper"
)

func main() {
	cfg := config.MustLoad()
	_ = cfg

	w := watcher.New()
	w.AddFilterHook(watcher.RegexFilterHook(regexp.MustCompile(`\.go$`), false))

	debounced, _ := debounce.DebounceWithCustomSignature(update, time.Millisecond*100)

	go func() {
		for {
			select {
			case e := <-w.Event:
				fmt.Println(e)
				_ = debounced(e)
			case err := <-w.Error:
				log.Fatalln(err)
			case <-w.Closed:
				return
			}
		}
	}()

	if err := w.AddRecursive("."); err != nil {
		log.Fatalln(err)
	}

	for path, f := range w.WatchedFiles() {
		fmt.Printf("%s: %s\n", path, f.Name())
	}

	go func() {
		w.Wait()

		fmt.Println("watcher started")
	}()

	if err := w.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
}

func update(ee ...watcher.Event) error {
	fmt.Println("updating...")

	cmd := exec.Command(
		"dlv",
		"debug",
		"--headless",
		"-l",
		":2345",
		"--api-version=2",
		"--accept-multiclient",
		"--log",
		"--continue",
		"./cmd/api",
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		log.Println(err)

		return err
	}

	return nil
}
