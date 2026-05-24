package models

type ModuleConfig struct {
	Name            string
	ProjectPath     string
	Architecture    Architecture
	StateManagement StateManagement
}
