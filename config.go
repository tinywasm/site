package site

var (
	config = &Config{
		CacheSize:    3,
		DefaultRoute: "home",
		OutputDir:    "./public",
		DevMode:      false,
	}
)

type Config struct {
	CacheSize    int
	DefaultRoute string
	OutputDir    string
	DevMode      bool
}

// SetCacheSize configures module cache size (default: 3)
func SetCacheSize(size int) {
	config.CacheSize = size
}

// SetDefaultRoute configures default route (default: "home")
func SetDefaultRoute(route string) {
	config.DefaultRoute = route
}

// SetOutputDir configures the output directory for assets (default: "./public")
func SetOutputDir(dir string) {
	config.OutputDir = dir
}

// SetDevMode configures development mode (default: false)
func SetDevMode(enabled bool) {
	config.DevMode = enabled
	handler.DevMode = enabled
}
