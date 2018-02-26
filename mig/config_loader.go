package mig

type loader struct{}

// NewConfigLoader returns new configuration loader.
func NewConfigLoader() Loader {
    return &loader{}
}

// Load loads configurator based on url.
// For now it only knows how to load local YAML configuration file.
// In the future it will be able to load configuration from different sources
// based on url.
func (l *loader) Load(url string) (Configurator, error) {
    return NewYAMLConfigurator(url)
}
