package main

type Hosts struct {
	Hosts []*Host `yaml:"hosts"`
}

type Host struct {
	Name         string `yaml:"name"`
	HTTPEndpoint string `yaml:"http"`
}
