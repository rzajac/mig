package mig

// loader represents configuration loader.
type loader struct{}

// NewConfigLoader returns new configuration loader.
func NewConfigLoader() *loader {
    return &loader{}
}

// Load loads Configurer based on url.
//
// For now it only knows how to load local YAML configuration file.
// In the future it will be able to load configuration from different sources
// based on url.
func (l *loader) Load(url string) (Configurer, error) {
    return NewYAMLCfg(url)
}
