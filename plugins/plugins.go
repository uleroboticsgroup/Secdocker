package plugins

import (
	"fmt"
	"os"
	PluginProcesser "plugin"

	"github.com/uleroboticsgroup/Secdocker/config"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

type Plugin interface {
	Process(string) bool
}

func ProcessPlugins(image string) bool {
	plugins := config.LoadConfig().Plugins

	result := true
	for _, pluginName := range plugins {
		// load module
		plug, err := PluginProcesser.Open("./plugins/" + pluginName + "/" + pluginName + ".so")
		checkErr(err)

		// 2. look up a symbol (an exported function or variable)
		symPlugin, err := plug.Lookup("Plugin")
		checkErr(err)

		// 3. Assert that loaded symbol is of a desired type
		var plugin Plugin
		plugin, ok := symPlugin.(Plugin)
		if !ok {
			fmt.Println("unexpected type from module symbol")
			os.Exit(1)
		}

		// 4. use the module
		if !plugin.Process(image) {
			fmt.Printf("Plugin %s exited with errors", pluginName)
			result = false
		}
	}

	return result
}
