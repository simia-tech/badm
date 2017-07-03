package badm

import (
	"fmt"
	"plugin"
)

// ListPlugins prints a list of plugins.
func ListPlugins(configurationPath string) error {
	c, err := readConfiguration(configurationPath)
	if err != nil {
		return fmt.Errorf("read configuration: %v", err)
	}

	p, err := openPlugins(c)
	if err != nil {
		return fmt.Errorf("open plugins: %v", err)
	}

	if err := p.list(); err != nil {
		return fmt.Errorf("list: %v", err)
	}

	return nil
}

// AddPlugin adds the provided plugin.
func AddPlugin(configurationPath, path string) error {
	return updateConfiguration(configurationPath, func(c *Configuration) error {
		c.Plugins = append(c.Plugins, &Plugin{Path: path})
		return nil
	})
}

// RegisterTypes registers all types from all plugins.
func RegisterTypes(configurationPath string) error {
	c, err := readConfiguration(configurationPath)
	if err != nil {
		return fmt.Errorf("read configuration: %v", err)
	}

	p, err := openPlugins(c)
	if err != nil {
		return fmt.Errorf("open plugins: %v", err)
	}

	if err := p.registerTypes(); err != nil {
		return fmt.Errorf("load types: %v", err)
	}

	return nil
}

type plugins struct {
	plugins []*plugin.Plugin
}

func openPlugins(c *Configuration) (*plugins, error) {
	ps := []*plugin.Plugin{}
	for _, pc := range c.Plugins {
		p, err := plugin.Open(pc.Path)
		if err != nil {
			return nil, fmt.Errorf("could not open plugin %s: %v", pc.Path, err)
		}
		ps = append(ps, p)
	}
	return &plugins{
		plugins: ps,
	}, nil
}

func (p *plugins) registerTypes() error {
	for _, plugin := range p.plugins {
		m, err := plugin.Lookup("RegisterTypes")
		if err != nil {
			return nil
		}
		registerTypes, ok := m.(func() error)
		if !ok {
			return nil
		}
		if err := registerTypes(); err != nil {
			return fmt.Errorf("could not register types: %v", err)
		}
	}
	return nil
}

func (p *plugins) list() error {
	for _, plugin := range p.plugins {
		m, err := plugin.Lookup("Name")
		if err != nil {
			return fmt.Errorf("could not find method Name in plugin: %v", err)
		}
		name := m.(func() string)
		fmt.Println(name())
	}
	return nil
}
