package core

/*
import (
	"strings"
	"sync"
	"time"

	k8soperator "gitlab.inspr.dev/inspr/core/cmd/operator/k8s/k8soperator"
	meta "gitlab.inspr.dev/inspr/core/pkg/meta"
)

type appStatus struct {
	// Associate an app with its status
	status map[string]string

	// A mutex to support parallel implementations
	sync.Mutex
}

//---------------------------------------------------------
//  This data structure is projected to support unique
//  names
//---------------------------------------------------------
type insprdMemory struct {
	// An empty string represents a channel that exists in memory
	// Example:    ch01 -> "" represents that ch01 exists in memory
	// Example:    ch02 -> ch01 represents that ch02 is alias to ch01
	aliasChannel map[string]string

	// An identifier that maps App's ID to apps.
	apps     map[string]*meta.App
	channels map[string]*meta.Channel

	// Define the scopes of channels
	// Channel's ID to its app's parent ID
	channelScopes map[string]string

	// Define the scopes of channels
	// Apps's ID to its app's parent ID
	appScopes map[string]string

	// Define the root of the tree
	tree *meta.App

	// A mutex to support parallel implementations
	sync.Mutex

	// Associate an app with its status
	monitor appStatus

	// A flag to kill all sub process when something get wrong
	stopSubProcess bool

	// Kind of app status
	// Valid kinds are:
	statusKind string
}

// InsprDTree defines the interface to operate inspr's memory objects
type InsprDTree interface {
	CreateApp(scopeName string, app *meta.App)
	DeleteApp(appName string)
	UpdateApp(app *meta.App)
	CreateChannel(scopeName string, channel *meta.Channel)
	DeleteChannel(channelName string)
	CreateAliasChannel(channelFrom string, channelTo string)
	DeleteAliasChannel(channelFrom string)
	getStruct() *insprdMemory
}

func newRoot() *meta.App {
	app := meta.App{}
	app.Metadata.Name = "inspr"
	app.Parent = ""
	app.Spec.Apps = make([]*meta.App, 0)
	app.Spec.Channels = make([]*meta.Channel, 0)
	return &app
}

// NewInsprDTree instanciate a new Inspr Tree object
func NewInsprDTree() InsprDTree {
	tree := insprdMemory{
		aliasChannel:  make(map[string]string),
		apps:          make(map[string]*meta.App),
		channels:      make(map[string]*meta.Channel),
		channelScopes: make(map[string]string),
		appScopes:     make(map[string]string),
		tree:          newRoot(),
	}
	tree.apps["inspr"] = tree.tree
	tree.appScopes["inspr"] = ""
	tree.stopSubProcess = false
	return &tree
}

func (m *insprdMemory) CreateApp(scopeName string, app *meta.App) {
	m.Lock()
	appName := app.Metadata.Name
	m.apps[appName] = app
	m.appScopes[appName] = scopeName
	m.apps[scopeName].Spec.Apps =
		append(m.apps[scopeName].Spec.Apps, app)
	// TODO: code to create the app on kubernetes
	m.Unlock()
}

func (m *insprdMemory) DeleteApp(appName string) {
	m.Lock()
	if _, ok := m.apps[appName]; !ok {
		m.Unlock()
		return
	}
	// Delete the app dependencies for apps
	for index := 0; index < len(m.apps[appName].Spec.Apps); index++ {
		if m.apps[appName].Spec.Apps[index] == nil {
			continue
		}

		key := m.apps[appName].Spec.Apps[index].Metadata.Name

		size := len(m.apps[appName].Spec.Apps) - 1
		m.apps[appName].Spec.Apps =
			append(m.apps[appName].Spec.Apps[:size],
				m.apps[appName].Spec.Apps[size+1:]...)

		m.Unlock()
		m.DeleteApp(key)
		m.Lock()
	}
	// Delete the app dependencies for channels
	for index := 0; index < len(m.apps[appName].Spec.Channels); index++ {
		key := m.apps[appName].Spec.Channels[index].Meta.Name
		m.Unlock()
		m.DeleteChannel(key)
		m.Lock()
		copy(m.apps[appName].Spec.Channels[index:],
			m.apps[appName].Spec.Channels[index+1:])
		m.apps[appName].Spec.Channels =
			m.apps[appName].Spec.Channels[:len(m.apps[appName].Spec.Channels)-1]
	}

	// Delete static references to an app
	_, ok := m.apps[appName]
	if ok {
		delete(m.apps, appName)
	}
	_, ok = m.appScopes[appName]
	if ok {
		delete(m.appScopes, appName)
	}
	// TODO: code to delete the app on kubernetes
	m.Unlock()
}

func (m *insprdMemory) UpdateApp(app *meta.App) {
	name := app.Metadata.Name
	m.Lock()
	scope := app.Metadata.Parent
	m.Unlock()
	m.DeleteApp(name)
	m.CreateApp(scope, app)
}

func (m *insprdMemory) CreateChannel(scopeName string, channel *meta.Channel) {
	m.Lock()
	channelName := channel.Meta.Name
	m.channelScopes[channelName] = scopeName
	m.apps[scopeName].Spec.Channels =
		append(m.apps[scopeName].Spec.Channels, channel)
	m.channels[channelName] = channel
	// TODO: code to create the channel on kubernetes
	m.Unlock()
}

func (m *insprdMemory) DeleteChannel(channelName string) {
	m.Lock()
	delete(m.channelScopes, channelName)
	delete(m.channels, channelName)
	// TODO: code to delete the channel on kubernetes
	m.Unlock()
}

func (m *insprdMemory) recoverChannelTarget(channelTo string) string {
	m.Lock()
	key := channelTo
	for {
		if _, ok := m.aliasChannel[key]; !ok {
			break
		}
		if m.aliasChannel[key] != "" {
			key = m.aliasChannel[key]
		} else {
			break
		}
	}
	m.Unlock()
	return key
}

func (m *insprdMemory) CreateAliasChannel(channelFrom string, channelTo string) {
	m.Lock()
	m.aliasChannel[channelFrom] = channelTo
	m.Unlock()
	channelTarget := m.recoverChannelTarget(channelTo)
	m.Lock()
	for k := range m.apps {
		// --------------------------------------------------------------------
		// Delete the unused Channel:
		// Kube.delete(m.apps[k].Spec.Channels[channelFrom])

		// TODO: Add code to propagatete the change to kubernetes nodes

		// --------------------------------------------------------------------
		//
		//
		// --------------------------------------------------------------------
		// Update all alias to the target reference:
		for index := range m.apps[k].Spec.Channels {
			if m.apps[k].Spec.Channels[index].Meta.Name == channelFrom {
				m.apps[k].Spec.Channels[index] = m.channels[channelTarget]
			}
		}
		// TODO: Add code to propagatete the change to kubernetes nodes
		// --------------------------------------------------------------------

	}

	for k := range m.apps {
		// Delete duplicated channels:
		count := 0
		for index := range m.apps[k].Spec.Channels {
			if m.apps[k].Spec.Channels[index].Meta.Name == channelTarget {
				if count >= 1 {
					copy(m.apps[k].Spec.Channels[index:],
						m.apps[k].Spec.Channels[index+1:])
					m.apps[k].Spec.Channels =
						m.apps[k].Spec.Channels[:len(m.apps[k].Spec.Channels)-1]
				}
				count++
			}
		}
	}
	m.Unlock()
}

func (m *insprdMemory) DeleteAliasChannel(channelFrom string) {
	m.Lock()
	delete(m.aliasChannel, channelFrom)
	m.Unlock()
}

func (m *insprdMemory) MonitoreApps() {
	go func() {
		for {
			m.Lock()
			stop := m.stopSubProcess
			m.Unlock()
			if stop == true {
				return
			}
			m.monitor.Lock()
			var toUpdate []string
			for key, value := range k8soperator.UpdateNodeStatus() {
				m.monitor.status[key] = value
				if value != "running" {
					if m.statusKind == "hard" {
						toUpdate = append(toUpdate, key)
					}
				}
			}
			if m.statusKind == "hard" {
				for i := 0; i < len(toUpdate); i++ {
					for key := range m.monitor.status {
						if strings.Contains(key, toUpdate[i]) {
							m.monitor.status[key] = m.monitor.status[toUpdate[i]]
						}
					}
				}
			}
			m.monitor.Unlock()
			time.Sleep(500)
		}
	}()
}

func (m *insprdMemory) getStruct() *insprdMemory {
	return m
}
*/
