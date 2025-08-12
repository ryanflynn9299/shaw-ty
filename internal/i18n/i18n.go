package i18n

import (
	"embed"
	"io/fs"
	"path/filepath"
	"strings"
	"sync"

	"gopkg.in/yaml.v3"
)

var (
	//go:embed locales/*.yaml
	localesFS embed.FS

	loadOnce      sync.Once
	loadErr       error
	catalogs      = map[string]map[string]string{}
	defaultLocale = "en"
)

// Load parses embedded locale YAML files into memory.
// Safe to call multiple times; work happens only once.
func Load() error {
	loadOnce.Do(func() {
		matches, err := fs.Glob(localesFS, "locales/*.yaml")
		if err != nil {
			loadErr = err
			return
		}
		for _, path := range matches {
			data, err := localesFS.ReadFile(path)
			if err != nil {
				loadErr = err
				return
			}
			var raw map[string]interface{}
			if err := yaml.Unmarshal(data, &raw); err != nil {
				loadErr = err
				return
			}
			flat := map[string]string{}
			flatten("", raw, flat)

			base := filepath.Base(path)                 // e.g. en.yaml
			locale := strings.TrimSuffix(base, ".yaml") // e.g. en
			catalogs[normalizeLocale(locale)] = flat
		}
	})
	return loadErr
}

// T returns the localized text for key in locale.
// Falls back to defaultLocale; if missing, returns the key itself.
func T(locale, key string) string {
	if err := Load(); err != nil {
		return key
	}
	if c, ok := catalogs[normalizeLocale(locale)]; ok {
		if msg, ok := c[key]; ok {
			return msg
		}
	}
	if c, ok := catalogs[defaultLocale]; ok {
		if msg, ok := c[key]; ok {
			return msg
		}
	}
	return key
}

// FromAcceptLanguage extracts a locale like "en" from an Accept-Language header.
// If empty or unsupported, returns defaultLocale.
func FromAcceptLanguage(header string) string {
	lang := strings.TrimSpace(header)
	if lang == "" {
		return defaultLocale
	}
	first := strings.Split(lang, ",")[0]
	first = strings.Split(first, ";")[0]
	base := strings.ToLower(strings.Split(first, "-")[0])
	if base == "" {
		return defaultLocale
	}
	return base
}

// SetDefaultLocale sets the default locale used for fallback.
func SetDefaultLocale(l string) {
	defaultLocale = normalizeLocale(l)
}

func normalizeLocale(l string) string {
	ll := strings.ToLower(strings.TrimSpace(l))
	if ll == "" {
		return "en"
	}
	return strings.Split(ll, "-")[0]
}

func flatten(prefix string, in map[string]interface{}, out map[string]string) {
	for k, v := range in {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}
		switch val := v.(type) {
		case string:
			out[key] = val
		case map[string]interface{}:
			flatten(key, val, out)
		}
	}
}
