package models

type Architecture int

const (
	ArchClean Architecture = iota
	ArchMVC
	ArchMVVM
)

func (a Architecture) String() string {
	switch a {
	case ArchClean:
		return "CLEAN Architecture"
	case ArchMVC:
		return "MVC"
	case ArchMVVM:
		return "MVVM"
	default:
		return "Unknown"
	}
}

type StateManagement int

const (
	StateBLoC StateManagement = iota
	StateRiverpod
	StateSignals
	StateProvider
	StateGetX
)

func (s StateManagement) String() string {
	switch s {
	case StateBLoC:
		return "Flutter BLoC"
	case StateRiverpod:
		return "Riverpod"
	case StateSignals:
		return "Signals"
	case StateProvider:
		return "Provider"
	case StateGetX:
		return "GetX"
	default:
		return "Unknown"
	}
}

type ProjectConfig struct {
	Name            string
	BundleID        string
	Architecture    Architecture
	StateManagement StateManagement
	Path            string
}
