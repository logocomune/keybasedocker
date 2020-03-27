package message

import (
	"strings"
	"time"

	lru "github.com/hashicorp/golang-lru"
)

const (
	cacheSize     = 128
	dockerWebhook = "docker-webhook"
)

type Formatter struct {
	cache    *lru.Cache
	nodeName string
}

//NewFormatter Initialize new formatter
func NewFormatter(nodeName string) (*Formatter, error) {
	cache, err := lru.New(cacheSize)
	if err != nil {
		return nil, err
	}

	return &Formatter{
		cache:    cache,
		nodeName: nodeName,
	}, nil
}

//EventPacker Receive a list of docker events and group them by container id
func (f *Formatter) EventPacker(events []Event) map[string]EventsGroup {
	aggrMsgs := make(map[string]EventsGroup)

	for _, e := range events {
		id := e.Container.ID

		if e.MetaData.Type == dockerWebhook {
			id = dockerWebhook
			aggrMsgs[id] = EventsGroup{
				StartupInfo: e.StartupInfo,
			}

			continue
		}

		eventsGroup := aggrMsgs[id]
		eventsGroup.Events = append(eventsGroup.Events, e)

		if e.MetaData.Type == "container" {
			eventsGroup.Container = e.Container
			f.cache.Add(id, e.Container)
		}

		aggrMsgs[id] = eventsGroup
	}

	return aggrMsgs
}

//EventsToStr Create a keybase message string for a group of container events
func (f *Formatter) EventsToStr(eventsGroup EventsGroup) (string, bool) {
	s := startupMessage(f.nodeName, eventsGroup.StartupInfo)
	if len(eventsGroup.Events) == 0 && s == "" {
		return "", false
	}

	nEvents := len(eventsGroup.Events)

	containerID := eventsGroup.ID
	containerImage := eventsGroup.Image

	containerName := eventsGroup.Name
	if containerName == "" {
		//If container name is missing... try to use cache value
		if c, ok := f.cache.Get(containerID); ok {
			container := c.(Container)
			containerName = container.Name
			containerImage = container.Image
		}
	}

	for idx, event := range eventsGroup.Events {
		if idx == 0 {
			s += titleMessage(containerName, containerImage, f.nodeName, event.Time)
		}

		switch event.MetaData.Type {
		case "container":
			s += containerMessage(event.MetaData, event.Container, event.ContainerStatus)
		case "volume":
			s += volumeMessage(event.MetaData, event.Volume)

		case "network":
			s += networkMessage(event.MetaData, event.Network)
		}

		if idx != nEvents-1 {
			s += "> \n"
		} else {
			s += footerMessage(containerID)
		}
	}

	return s, true
}

func startupMessage(nodeName string, info StartupInfo) string {
	if info.DockerVersion == "" {
		return ""
	}

	msg := startupTitle1 + startupTitle2

	msg = strings.Replace(msg, "__DOCKER_VERSION__", info.DockerVersion, -1)
	msg = strings.Replace(msg, "__DOCKER_API_VERSION__", info.APIVersion, -1)
	msg = strings.Replace(msg, "__OS__", info.Os, -1)
	msg = strings.Replace(msg, "__KERNEL_VERSION__", info.KernelVersion, -1)
	msg = strings.Replace(msg, "__NEW_LINE__", "\n", -1)
	msg = strings.Replace(msg, "__TAB__", "\t", -1)
	hostnameReplacer := ""

	if nodeName != "" {
		hostnameReplacer = nodeName
	}

	msg = strings.Replace(msg, "__NODE_NAME__", hostnameReplacer, -1)

	return msg
}

func titleMessage(name string, image string, nodeName string, t time.Time) string {
	time := t.Format(time.RFC3339)
	msg := groupTitle
	msg = strings.Replace(msg, "__IMAGE__", image, -1)
	msg = strings.Replace(msg, "__NAME__", name, -1)
	msg = strings.Replace(msg, "__TIME__", time, -1)
	msg = strings.Replace(msg, "__NEW_LINE__", "\n", -1)
	hostnameReplacer := ""

	if nodeName != "" {
		hostnameReplacer = "*@*_" + nodeName + "_ "
	}

	msg = strings.Replace(msg, "__NODE_NAME__", hostnameReplacer, -1)

	return msg
}

func footerMessage(id string) string {
	msg := groupFooter
	msg = strings.Replace(msg, "__ID__", id, -1)
	msg = strings.Replace(msg, "__NEW_LINE__", "\n", -1)

	return msg
}

func containerMessage(meta MetaData, eContainer Container, eStatus ContainerStatus) string {
	msg := ""

	switch meta.Action {
	case "kill":
		msg = containerKill

	case "die":
		msg = containerDie
	default:
		msg = containerDefault
	}

	instanceID := eContainer.ID
	action := meta.Action
	image := eContainer.Image
	name := eContainer.Name
	time := meta.Time.Format(time.RFC3339)
	exitCode := eStatus.ExitCode
	signal := eStatus.Signal

	msg = strings.Replace(msg, "__ID__", instanceID, -1)
	msg = strings.Replace(msg, "__ACTION__", strings.Title(action), -1)
	msg = strings.Replace(msg, "__IMAGE__", image, -1)
	msg = strings.Replace(msg, "__NAME__", name, -1)
	msg = strings.Replace(msg, "__TIME__", time, -1)
	msg = strings.Replace(msg, "__EXIT_CODE__", exitCode, -1)
	msg = strings.Replace(msg, "__SIGNAL__", signal, -1)
	msg = strings.Replace(msg, "__NEW_LINE__", "\n", -1)

	return msg
}

func volumeMessage(meta MetaData, volume Volume) string {
	id := volume.ID
	dest := volume.Destination
	action := meta.Action

	msg := ""

	switch action {
	case "mount":
		msg = volumeMount

	case "unmount":
		msg = volumeUnmount
	}

	msg = strings.Replace(msg, "__ACTION__", strings.Title(action), -1)
	msg = strings.Replace(msg, "__VOLUME_ID__", id, -1)
	msg = strings.Replace(msg, "__VOLUME_DESTINATION__", dest, -1)
	msg = strings.Replace(msg, "__NEW_LINE__", "\n", -1)

	return msg
}

func networkMessage(meta MetaData, network Network) string {
	msg := networkDefault
	name := network.Name
	id := network.ID
	action := meta.Action

	msg = strings.Replace(msg, "__ACTION__", strings.Title(action), -1)
	msg = strings.Replace(msg, "__NETWORK_ID__", id, -1)
	msg = strings.Replace(msg, "__NETWORK_NAME__", name, -1)
	msg = strings.Replace(msg, "__NEW_LINE__", "\n", -1)

	return msg
}
