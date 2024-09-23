package main

import (
	"bufio"
	"flare_plugin_generator/plugin"
	"flare_plugin_generator/zip"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("flare-plugin release generator")

	pluginRelease := plugin.PluginRelease{}
	promptPluginReleaseData(&pluginRelease)

	// feedback plugin release input
	fmt.Println(pluginRelease)

	err := generateFiles(pluginRelease)
	if err != nil {
		log.Println(err)
		return
	}

	zip.Zip(pluginRelease.Package, fmt.Sprintf("%s.zip", pluginRelease.Package))

	cleanup(pluginRelease.Package)
}

// prompt for the plugin release input data
func promptPluginReleaseData(pluginRelease *plugin.PluginRelease) {
	// TODO: add an options for inputting plugin data
	// either manually or through a file like data.txt

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Name: ")
	scanner.Scan()
	pluginRelease.Name = scanner.Text()
	if pluginRelease.Name == "" {
		pluginRelease.Name = "Sample Plugin"
	}

	fmt.Print("Package: ")
	scanner.Scan()
	pluginRelease.Package = scanner.Text()
	if pluginRelease.Package == "" {
		pluginRelease.Package = "com.sample.plugin"
	}

	fmt.Print("Description: ")
	scanner.Scan()
	pluginRelease.Description = scanner.Text()
	if pluginRelease.Description == "" {
		pluginRelease.Description = "A sample plugin"
	}

	fmt.Print("Version: ")
	scanner.Scan()
	if scanner.Text() == "" {
		parsedVersion, err := parseVersion("0.0.1")
		if err != nil {
			log.Println(err)
			return
		}
		pluginRelease.Version = parsedVersion
		return
	}
	parsedVersion, err := parseVersion(scanner.Text())
	if err != nil {
		log.Println(err)
		return
	}
	pluginRelease.Version = parsedVersion
}

// Parses the string version to int based on semantic versioning
func parseVersion(version string) (plugin.PluginReleaseVersion, error) {
	versions := strings.Split(version, ".")

	major, err := strconv.Atoi(versions[0])
	if err != nil {
		log.Println(err)
		return plugin.PluginReleaseVersion{}, err
	}

	minor, err := strconv.Atoi(versions[1])
	if err != nil {
		log.Println(err)
		return plugin.PluginReleaseVersion{}, err
	}

	patch, err := strconv.Atoi(versions[2])
	if err != nil {
		log.Println(err)
		return plugin.PluginReleaseVersion{}, err
	}

	return plugin.PluginReleaseVersion{
		Major: major,
		Minor: minor,
		Patch: patch,
	}, nil
}

func stringifyVersion(version plugin.PluginReleaseVersion) string {
	return fmt.Sprintf("%d.%d.%d", version.Major, version.Minor, version.Patch)
}

func generateFiles(pluginRelease plugin.PluginRelease) error {
	// generate dir
	err := os.Mkdir(pluginRelease.Package, 0755)
	if err != nil {
		log.Println(err)
		return err
	}

	// define plugin files
	pluginFiles := []plugin.PluginFile{
		{
			FileName:    "go.mod",
			FileContent: goModContent(pluginRelease.Package),
		},
		{
			FileName:    "main.go",
			FileContent: mainGoContent(),
		},
		{
			FileName:    "plugin.json",
			FileContent: pluginJsonContent(pluginRelease.Name, pluginRelease.Package, pluginRelease.Description, stringifyVersion(pluginRelease.Version)),
		},
		{
			FileName:    ".gitignore",
			FileContent: gitignoreContent(),
		},
	}

	// write plugin files
	for _, pf := range pluginFiles {
		err := os.WriteFile(filepath.Join(pluginRelease.Package, pf.FileName), []byte(pf.FileContent), 0644)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	fmt.Println("plugin files successfully written..")
	return nil
}

func goModContent(packageName string) string {
	return fmt.Sprintf(`
module %s

go 1.21.0
    `, packageName)
}

func mainGoContent() string {
	return `
package main

import (
    sdkplugin "sdk/api/plugin"
)

func main() {}

func Init(api sdkplugin.PluginApi) {
    // Your plugin code here
}
    `
}

func pluginJsonContent(pluginName, pluginPkg, pluginDesc, pluginVer string) string {
	return fmt.Sprintf(`
{
  "Name": "%s",
  "Package": "%s",
  "Description": "%s",
  "Version": "%s"
}
    `, pluginName, pluginPkg, pluginDesc, pluginVer)
}

func gitignoreContent() string {
	return `# Ignore main_mono.go
main_mono.go`
}

func cleanup(dir string) error {
	err := os.RemoveAll(dir)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
