package message

import "time"

type EventsGroup struct {
	Container
	StartupInfo
	Events []Event
}
type Event struct {
	MetaData
	StartupInfo
	Container
	ContainerStatus
	Volume
	Network
}

type MetaData struct {
	Type   string
	Action string
	Time   time.Time
}

type StartupInfo struct {
	DockerVersion string
	APIVersion    string
	Os            string
	KernelVersion string
}

type Container struct {
	ID    string
	Name  string
	Image string
}

func (c Container) IsEmpty() bool {
	return c.ID == ""
}

type ContainerStatus struct {
	Signal   string
	ExitCode string
}

type Volume struct {
	ID          string
	Destination string
	Propagation string
}

type Network struct {
	ID   string
	Name string
	Type string
}
