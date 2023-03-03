package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func main() {
	if lsfiles, err := exec.Command("git", "ls-files").Output(); err != nil {
		panic(err)
	} else {
		var wg sync.WaitGroup
		var files sync.Map
		for _, file := range strings.Split(string(lsfiles), "\n") {
			wg.Add(1)
			go func(file string) {
				log := exec.Command("git", "log", "-n 1", "--pretty=%h", "--", file)
				hash, _ := log.Output()
				files.Store(file, strings.TrimSpace(string(hash)))
				wg.Done()
			}(file)
		}
		wg.Wait()

		output := make(map[string]string)

		files.Range(func(key, value any) bool {
			output[key.(string)] = value.(string)
			return true
		})

		var raw []byte
		if raw, err = json.Marshal(output); err != nil {
			panic(err)
		} else if err = os.WriteFile("manifest.json", raw, 0); err != nil {
			panic(err)
		}

	}
}
