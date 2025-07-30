package main

type DynamicConfig struct {
	HTTP HttpConfig `yaml:"http"`
}

type HttpConfig struct {
	Routers  map[string]RouterConfig  `yaml:"routers"`
	Services map[string]ServiceConfig `yaml:"services"`
}

type RouterConfig struct {
	Rule        string   `yaml:"rule"`
	Service     string   `yaml:"service"`
	EntryPoints []string `yaml:"entryPoints"`
}

type ServiceConfig struct {
	LoadBalancer LoadBalancerConfig `yaml:"loadBalancer"`
}

type LoadBalancerConfig struct {
	Servers []ServerConfig `yaml:"servers"`
}

type ServerConfig struct {
	URL string `yaml:"url"`
}
