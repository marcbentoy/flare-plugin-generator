package plugin

type PluginReleaseVersion struct {
	Major int
	Minor int
	Patch int
}

type PluginRelease struct {
	Name        string
	Package     string
	Description string
	Version     PluginReleaseVersion
}
