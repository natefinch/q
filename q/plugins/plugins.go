package plugins

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/BurntSushi/toml"

	"npf.io/q/q/log"
)

type Manifest struct {
	Name       string    // plugin name
	PluginPath string    `toml:"-"`
	Version    string    // plugin version
	Commands   []Command `toml:"Command"` // commands w/o context
	Contexts   []Context `toml:"Context"` // contexts
	Services   []Service `toml:"Service"` // services
}

type Command struct {
	Name  string // command name
	Short string // single line help text
	Long  string // multi-line help text
}

type Context struct {
	Name     string    // context name
	Plural   string    // plural version of context name
	Short    string    // single line help text
	Long     string    // multi-line help text
	Commands []Command `toml:"Command"` // commands that use this context
}

type Service struct {
	Name    string // name of the service
	Version string // version of the service API
}

func LoadManifests(dir string) ([]*Manifest, []error) {
	names, err := filepath.Glob(filepath.Join(dir, "*.toml"))
	if err != nil {
		return nil, []error{fmt.Errorf("error reading manifests from plugin directory: %v", err)}
	}

	// Map of manifest file name to plugin file name.
	// The manifest name should be <plugin_filename>.toml
	// We use this to look for orphaned manifests or plugins that haven't had
	// their manifest extracted yet.
	files := make(map[string]string, len(names))
	for _, name := range names {
		files[name] = ""
	}

	fs, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, []error{fmt.Errorf("can't read plugin directory: %v", err)}
	}

	wg := &sync.WaitGroup{}
	results := make(chan manifestResult)

	for _, f := range fs {
		name := f.Name()
		if _, ok := files[name]; ok {
			// manifest file
			continue
		}
		// not a manifest, should be a plugin, check for corresponding
		// manifest
		manifest := manifestForPlugin(name)
		if _, ok := files[manifest]; ok {
			// now that we know both exist, we can parse the manifest
			wg.Add(1)
			go readManifest(name, manifest, wg, results)
		}
		// file with no corresponding .toml manifest - must be a plugin,
		// run it to get the manifest.
		wg.Add(1)
		go getManifest(name, wg, results)
	}

	collres := make(chan collateResult)

	go collate(files, results, collres)

	wg.Wait()
	close(results)
	result := <-collres

	for manifest, plugin := range files {
		if plugin == "" {
			log.Verbose("Orphaned manifest found without plugin: %v", manifest)
		}
	}
	return result.mfests, result.errs
}

func manifestForPlugin(path string) string {
	return path + ".toml"
}

type manifestResult struct {
	mfest *Manifest
	err   error
}

type collateResult struct {
	mfests []*Manifest
	errs   []error
}

func getManifest(plugin string, wg *sync.WaitGroup, result chan<- manifestResult) {
	defer wg.Done()
	cmd := exec.Command(plugin, "manifest")
	out := &bytes.Buffer{}
	cmd.Stdout = out
	if err := cmd.Start(); err != nil {
		result <- manifestResult{err: fmt.Errorf("couldn't run plugin %q: %v", plugin, err)}
		return
	}
	done := make(chan error)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		if err != nil {
			result <- manifestResult{err: fmt.Errorf("error running plugin %q: %v", plugin, err)}
			return
		}
	case <-time.After(time.Millisecond * 500):
		_ = cmd.Process.Kill() // ignore any possible error... there's not much we can do about it.
		result <- manifestResult{err: fmt.Errorf("timed out waiting for plugin %q to output manifest", plugin)}
		return
	}

	mfest := Manifest{PluginPath: plugin}
	meta, err := toml.DecodeReader(out, &mfest)
	if err != nil {
		result <- manifestResult{err: fmt.Errorf("failed to decode manifest from %q: %v", plugin, err)}
		return
	}
	if len(meta.Undecoded()) > 0 {
		log.Verbose("Unexpected options in manifest for %q: %v", plugin, meta.Undecoded())
	}
	manifest := manifestForPlugin(plugin)
	if err := ioutil.WriteFile(manifest, out.Bytes(), 0644); err != nil {
		result <- manifestResult{err: fmt.Errorf("couldn't write manifest %q: %v", manifest, err)}
	}
	result <- manifestResult{mfest: &mfest}
}

func readManifest(plugin, manifest string, wg *sync.WaitGroup, results chan<- manifestResult) {
	defer wg.Done()

	m := &Manifest{PluginPath: plugin}
	meta, err := toml.DecodeFile(manifest, m)
	if err != nil {
		results <- manifestResult{err: fmt.Errorf("failed to decode manifest from %q: %v", manifest, err)}
		return
	}
	if len(meta.Undecoded()) > 0 {
		log.Verbose("Unexpected values in manifest %q: %v", manifest, meta.Undecoded())
	}
	results <- manifestResult{mfest: m}
}

func collate(mfests map[string]string, results <-chan manifestResult, resChan chan<- collateResult) {
	result := collateResult{}
	for r := range results {
		if r.err != nil {
			result.errs = append(result.errs, r.err)
			continue
		}
		mfests[manifestForPlugin(r.mfest.PluginPath)] = r.mfest.PluginPath
		result.mfests = append(result.mfests, r.mfest)
	}
	resChan <- result
}
