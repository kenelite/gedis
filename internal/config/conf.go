package config

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

type Config struct {
	sections map[string]map[string]string
}

// Load reads and parses the config file at path.
func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := &Config{
		sections: make(map[string]map[string]string),
	}

	scanner := bufio.NewScanner(file)
	var currentSection string

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, ";") {
			continue
		}

		// Section header
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			sectionName := strings.ToLower(line[1 : len(line)-1])
			if _, exists := cfg.sections[sectionName]; !exists {
				cfg.sections[sectionName] = make(map[string]string)
			}
			currentSection = sectionName
			continue
		}

		if currentSection == "" {
			return nil, errors.New("key-value pair found outside of a section")
		}

		// key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			return nil, errors.New("invalid line: " + line)
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		cfg.sections[currentSection][key] = value
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Get returns the value for key in section. If not found, returns empty string.
func (c *Config) Get(section, key string) string {
	if sec, ok := c.sections[strings.ToLower(section)]; ok {
		if val, ok := sec[key]; ok {
			return val
		}
	}
	return ""
}

// GetSection returns the map of key-value pairs in the given section.
func (c *Config) GetSection(section string) map[string]string {
	if sec, ok := c.sections[strings.ToLower(section)]; ok {
		return sec
	}
	return nil
}
