package model

type Step struct {
	Name   string            `yaml:"name"`
	Uses   string            `yaml:"uses"`
	With   map[string]string `yaml:"with"`
	RunsOn string            `yaml:"runs-on"`
	Run    string            `yaml:"run"`
}
