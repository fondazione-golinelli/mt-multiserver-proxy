package proxy

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"plugin"
	"strings"
	"sync"
)

var pluginsOnce sync.Once

func BuildPlugin() error {
	if sourceDir := os.Getenv("PROXY_SOURCE_DIR"); sourceDir != "" {
		return buildPluginLocalSource(sourceDir)
	}
	version, err := Version()
	if err != nil {
		return err
	}

	if strings.Contains(version, "+dirty") {
		return buildPluginDev(version)
	}

	return buildPlugin(version)
}

func buildPlugin(version string) error {
	pathVer := "github.com/HimbeerserverDE/mt-multiserver-proxy@" + version

	if err := goCmd("get", pathVer); err != nil {
		return err
	}

	if err := goCmd("mod", "tidy"); err != nil {
		return err
	}
	if version != "(devel)" {
		output, err := exec.Command("go", "mod", "download", "-json", pathVer).Output()
		if err != nil {
			return err
		}
		var module struct{ Dir string }
		if err := json.Unmarshal(output, &module); err != nil {
			return err
		}
		if module.Dir != "" {
			return buildPluginLocalSource(module.Dir)
		}
	}
	return goCmd("build", "-buildmode=plugin")
}

func buildPluginDev(version string) error {
	sourceDir := Path()
	if _, err := os.Stat(sourceDir + "/go.mod"); err == nil {
		return buildPluginLocalSource(sourceDir)
	}
	replace := "-replace=github.com/HimbeerserverDE/mt-multiserver-proxy=" + Path()
	const dropReplace = "-dropreplace=github.com/HimbeerserverDE/mt-multiserver-proxy"
	if err := goCmd("mod", "edit", replace); err != nil {
		return err
	}
	if err := goCmd("mod", "tidy"); err != nil {
		return err
	}
	if err := goCmd("build", "-buildmode=plugin"); err != nil {
		return err
	}
	return goCmd("mod", "edit", dropReplace)
}

// Build with the same replacement and module version as the local image's
// binary build. A temporary modfile preserves the plugin's published pin.
func buildPluginLocalSource(sourceDir string) error {
	mod, err := os.ReadFile("go.mod")
	if err != nil {
		return err
	}
	file, err := os.CreateTemp(".", "proxy-local-*.mod")
	if err != nil {
		return err
	}
	modPath := file.Name()
	sumPath := strings.TrimSuffix(modPath, ".mod") + ".sum"
	defer os.Remove(modPath)
	defer os.Remove(sumPath)
	if _, err := file.Write(mod); err != nil {
		file.Close()
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	if sum, err := os.ReadFile("go.sum"); err == nil {
		if err := os.WriteFile(sumPath, sum, 0600); err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	modFlag := "-modfile=" + modPath
	replace := "-replace=github.com/HimbeerserverDE/mt-multiserver-proxy=" + sourceDir
	if err := goCmd("mod", "edit", modFlag, replace, "-require=github.com/HimbeerserverDE/mt-multiserver-proxy@v0.0.0"); err != nil {
		return err
	}
	if _, err := os.Stat(sourceDir + "/third_party/mt/go.mod"); err == nil {
		mtReplace := "-replace=github.com/HimbeerserverDE/mt=" + sourceDir + "/third_party/mt"
		if err := goCmd("mod", "edit", modFlag, mtReplace); err != nil {
			return err
		}
	} else if !os.IsNotExist(err) {
		return err
	}
	if err := goCmd("mod", "tidy", modFlag); err != nil {
		return err
	}
	return goCmd("build", modFlag, "-buildmode=plugin")
}

// LoadPlugins loads all plugins. Further calls are no-ops.
// This function is not intended to be called by plugins.
// It is meant for tools such as mt-auth-convert.
func LoadPlugins() {
	pluginsOnce.Do(openPlugins)
}

func openPlugins() {
	path := Path("plugins")
	os.Mkdir(path, 0777)

	dir, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}

	for _, pl := range dir {
		if pl.IsDir() && !Conf().NoAutoPlugins {
			plPath := path + "/" + pl.Name()

			wd, err := os.Getwd()
			if err != nil {
				log.Fatal(err)
			}

			if err := os.Chdir(plPath); err != nil {
				log.Fatal(err)
			}

			if err := BuildPlugin(); err != nil {
				log.Fatal(err)
			}

			if err := os.Chdir(wd); err != nil {
				log.Fatal(err)
			}

			_, err = plugin.Open(path + "/" + pl.Name() + "/" + pl.Name() + ".so")
			if err != nil {
				log.Print(err)
				continue
			}

			log.Println("load auto plugin", pl.Name())
		} else if !pl.IsDir() {
			_, err := plugin.Open(path + "/" + pl.Name())
			if err != nil {
				log.Print(err)
				continue
			}

			log.Println("load comp plugin", pl.Name())
		}
	}
}

func goCmd(args ...string) error {
	cmd := exec.Command("go", args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}
