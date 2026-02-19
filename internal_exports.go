package site

// TestResetHandler resets the global handler state for testing.
// For testing purposes only.
func TestResetHandler() {
	handler.registeredModules = nil
	handler.DevMode = false
}

// TestIsDevMode returns the current DevMode state of the handler.
// For testing purposes only.
func TestIsDevMode() bool {
	return handler.DevMode
}

// TestGetConfig returns the global configuration.
// For testing purposes only.
func TestGetConfig() *Config {
	return config
}

// TestParseRoute exposes the internal parseRoute function for testing.
// For testing purposes only.
func TestParseRoute(hash string) (module string, params []string) {
	return parseRoute(hash)
}

// TestGetModules returns the list of registered modules.
// For testing purposes only.
func TestGetModules() []*registeredModule {
	return handler.registeredModules
}
