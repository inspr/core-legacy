package core

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	meta "gitlab.inspr.dev/inspr/core/pkg/meta"
	"google.golang.org/grpc"
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

	// The address of a cluster operator
	clusterOperatorAddrress       string
	clusterOperatorCancelFunction context.CancelFunc
	clusterOperator               meta.NodeOperatorClient
	clusterOperatorContext        context.Context
}

// InsprDTree defines the interface to operate inspr's memory objects
type InsprDTree interface {
	CloseMemoryConnections()

	CreateApp(scopeName string, app *meta.App) error
	DeleteApp(appName string) error
	UpdateApp(app *meta.App) error

	CreateChannel(scopeName string, channel *meta.Channel) error
	DeleteChannel(channelName string) error
	CreateAliasChannel(channelFrom string, channelTo string) error
	DeleteAliasChannel(channelFrom string) error
	getStruct() *insprdMemory
}

func newRoot() *meta.App {
	app := meta.App{}
	app.Metadata.Name = "inspr"
	app.Metadata.Parent = ""
	app.Spec.Apps = make([]*meta.App, 0)
	app.Spec.Channels = make([]*meta.Channel, 0)
	return &app
}

func (m *insprdMemory) instantiateOperator() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(m.clusterOperatorAddrress,
		grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	m.clusterOperator = meta.NewNodeOperatorClient(conn)

	m.clusterOperatorContext, m.clusterOperatorCancelFunction =
		context.WithTimeout(context.Background(), time.Second)
}

func (m *insprdMemory) CloseMemoryConnections() {
	m.clusterOperatorCancelFunction()
}

// NewInsprDTree instanciate a new Inspr Tree object
func NewInsprDTree() InsprDTree {
	tree := insprdMemory{
		//aliasChannel:  make(map[string]string),
		apps:          make(map[string]*meta.App),
		channels:      make(map[string]*meta.Channel),
		channelScopes: make(map[string]string),
		appScopes:     make(map[string]string),
		tree:          newRoot(),
	}
	tree.apps["inspr"] = tree.tree
	tree.appScopes["inspr"] = ""
	tree.stopSubProcess = false
	tree.clusterOperatorAddrress = "localhost:50000"
	tree.instantiateOperator()
	return &tree
}

func (m *insprdMemory) CreateApp(scopeName string, app *meta.App) error {
	m.Lock()
	appName := app.Metadata.Name
	m.apps[appName] = app
	m.appScopes[appName] = scopeName
	m.apps[scopeName].Spec.Apps =
		append(m.apps[scopeName].Spec.Apps, app)
	// Code to create the app into the cluster
	nodeReply, err := m.clusterOperator.CreateNode(m.clusterOperatorContext, app.Spec.Node)
	if err != nil {
		m.Unlock()
		return err
	} else if nodeReply.Error != "" {
		m.Unlock()
		return errors.New(nodeReply.Error)
	}

	m.Unlock()
	return nil
}

func (m *insprdMemory) DeleteApp(appName string) error {
	m.Lock()
	if _, ok := m.apps[appName]; !ok {
		m.Unlock()
		return nil
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
		key := m.apps[appName].Spec.Channels[index].Metadata.Name
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
	nodeDescription := meta.NodeDescription{
		NodeDescription: appName,
	}
	nodeReply, err := m.clusterOperator.DeleteNode(m.clusterOperatorContext, &nodeDescription)
	if err != nil {
		m.Unlock()
		return err
	} else if nodeReply.Error != "" {
		m.Unlock()
		return errors.New(nodeReply.Error)
	}
	m.Unlock()
	return nil
}

func (m *insprdMemory) UpdateApp(app *meta.App) error {
	m.Lock()
	m.Unlock()
	err := m.UpdateApp(app)
	if err != nil {
		m.Unlock()
		return err
	}
	m.Unlock()
	return err
}

func (m *insprdMemory) CreateChannel(scopeName string, channel *meta.Channel) error {
	m.Lock()
	channelName := channel.Metadata.Name
	m.channelScopes[channelName] = scopeName
	m.apps[scopeName].Spec.Channels =
		append(m.apps[scopeName].Spec.Channels, channel)
	m.channels[channelName] = channel
	// TODO: code to create the channel on message Broker
	m.Unlock()
	return nil
}

func (m *insprdMemory) DeleteChannel(channelName string) error {
	m.Lock()
	delete(m.channelScopes, channelName)
	delete(m.channels, channelName)
	// TODO: code to delete the channel on message Broker
	m.Unlock()
	return nil
}

func (m *insprdMemory) CreateAliasChannel(channelFrom string, channelTo string) error {
	m.Lock()
	// TODO: create a alias just if the app use the channel.
	// At this point, we do not know the list of channels used by an app
	for k := range m.apps {
		// --------------------------------------------------------------------
		// Update all alias to the target reference:
		for index := 0; index < len(m.apps[k].Spec.ChannelAlias.Reference); index++ {
			if m.apps[k].Spec.ChannelAlias.Reference[index] == channelFrom {
				m.apps[k].Spec.ChannelAlias.Reference =
					append(m.apps[k].Spec.ChannelAlias.Reference[:index],
						m.apps[k].Spec.ChannelAlias.Reference[index+1:]...)

				m.apps[k].Spec.ChannelAlias.Target =
					append(m.apps[k].Spec.ChannelAlias.Target[:index],
						m.apps[k].Spec.ChannelAlias.Target[index+1:]...)
				index = index - 1
			}
		}
		m.apps[k].Spec.ChannelAlias.Reference =
			append(m.apps[k].Spec.ChannelAlias.Reference, channelFrom)
		m.apps[k].Spec.ChannelAlias.Target =
			append(m.apps[k].Spec.ChannelAlias.Target, channelTo)
	}
	m.Unlock()
	return nil
}

func (m *insprdMemory) DeleteAliasChannel(channelFrom string) error {
	m.Lock()
	// TODO: create a alias just if the app use the channel.
	// At this point, we do not know the list of channels used by an app
	for k := range m.apps {
		// --------------------------------------------------------------------
		// Update all alias to the target reference:
		for index := 0; index < len(m.apps[k].Spec.ChannelAlias.Reference); index++ {
			if m.apps[k].Spec.ChannelAlias.Reference[index] == channelFrom {
				m.apps[k].Spec.ChannelAlias.Reference =
					append(m.apps[k].Spec.ChannelAlias.Reference[:index],
						m.apps[k].Spec.ChannelAlias.Reference[index+1:]...)

				m.apps[k].Spec.ChannelAlias.Target =
					append(m.apps[k].Spec.ChannelAlias.Target[:index],
						m.apps[k].Spec.ChannelAlias.Target[index+1:]...)
				index = index - 1
			}
		}
		m.apps[k].Spec.ChannelAlias.Reference =
			append(m.apps[k].Spec.ChannelAlias.Reference, channelFrom)
		m.apps[k].Spec.ChannelAlias.Target =
			append(m.apps[k].Spec.ChannelAlias.Target, channelFrom)
	}
	m.Unlock()
	return nil
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
			nodes, err := m.clusterOperator.UpdateNodeStatus(m.clusterOperatorContext, &meta.Stub{})
			if err != nil {
				return
			}

			for index := 0; index < len(nodes.Node); index++ {
				key := nodes.Node[index].Metadata.Name
				value := nodes.Status[index]
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
						// Check if the name contain the app that should be updated.
						// In affirmative case, update it with the patent's status.
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
