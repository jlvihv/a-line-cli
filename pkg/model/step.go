package model

type Step struct {
	Name   string            `yaml:"name" json:"name"`
	Id     string            `yaml:"id" json:"id"`
	Uses   string            `yaml:"uses" json:"uses"`
	With   map[string]string `yaml:"with" json:"with"`
	RunsOn string            `yaml:"runs-on" json:"runsOn"`
	Run    string            `yaml:"run" json:"run"`
}
