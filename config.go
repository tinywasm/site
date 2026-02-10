package site

var (
	config = &Config{
		CacheSize:    3,
		DefaultRoute: "home",
	}
)

type Config struct {
	CacheSize    int
	DefaultRoute string
}

// SetCacheSize configures module cache size (default: 3)
func SetCacheSize(size int) {
	config.CacheSize = size
}

// SetDefaultRoute configures default route (default: "home")
func SetDefaultRoute(route string) {
	config.DefaultRoute = route
}
