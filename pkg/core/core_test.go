package core

import (
	"testing"

	meta "gitlab.inspr.dev/inspr/core/pkg/meta"
	"gitlab.inspr.dev/inspr/core/pkg/test"
)

func Test_newRoot(t *testing.T) {
	suite := test.NewSuite("Test_newRoot", t)
	var root *meta.App

	suite.BeforeEach(func() {
		root = newRoot()
	})

	suite.NewCase("Without error",
		func(it *testing.T) {
			name := root.Name
			test.AssertEquals(name, "inspr", it)
		})

	suite.RunAll()
}

func TestNewInsprDTree(t *testing.T) {
	suite := test.NewSuite("TestNewInsprDTree", t)
	tree := NewInsprDTree()

	suite.NewCase("Without error",
		func(it *testing.T) {
			test.AssertNotNil(tree, it)
		})

	suite.RunAll()
}

func Test_insprdMemory_CreateApp(t *testing.T) {
	suite := test.NewSuite("Test_insprdMemory_CreateApp", t)
	tree := NewInsprDTree()
	var app meta.App

	suite.BeforeEach(func() {
		app.Metadata.Name = "test"
		app.Parent = "inspr"
	})

	suite.NewCase("Without error",
		func(it *testing.T) {
			tree.CreateApp("test", &app)
			test.AssertNotNil(tree, it)

			object := tree.getStruct()
			test.AssertEquals(object.apps["test"].Metadata.Name, "test", it)
			test.AssertEquals(object.apps["test"].Metadata.Parent, "inspr", it)
		})

	suite.RunAll()
}

func Test_insprdMemory_DeleteApp(t *testing.T) {
	suite := test.NewSuite("Test_insprdMemory_DeleteApp", t)
	tree := NewInsprDTree()
	var app meta.App

	suite.BeforeEach(func() {
		app.Metadata.Name = "test"
		app.Metadata.Parent = "inspr"
	})

	suite.NewCase("With app",
		func(it *testing.T) {
			tree.CreateApp("inspr", &app)
			tree.DeleteApp("test")
			test.AssertNotNil(tree, it)

			object := tree.getStruct()
			_, ok := object.apps["test"]
			test.AssertEquals(ok, false, it)

			_, ok = object.appScopes["test"]
			test.AssertEquals(ok, false, it)
		})

	suite.NewCase("With sub app",
		func(it *testing.T) {
			tree.CreateApp("inspr", &app)

			var app2 meta.App
			app2.Metadata.Name = "test-2"
			app2.Metadata.Parent = "test"

			tree.CreateApp("test", &app2)

			tree.DeleteApp("test")
			test.AssertNotNil(tree, it)

			object := tree.getStruct()
			_, ok := object.apps["test"]
			test.AssertEquals(ok, false, it)

			_, ok = object.appScopes["test"]
			test.AssertEquals(ok, false, it)
		})

	suite.NewCase("With sub channel",
		func(it *testing.T) {
			tree.CreateApp("inspr", &app)

			var app2 meta.App
			app2.Metadata.Name = "test-2"
			app2.Metadata.Parent = "test"

			var channel meta.Channel
			channel.Meta.Name = "chan-test"
			tree.CreateChannel("test", &channel)

			tree.CreateApp("test", &app2)

			tree.DeleteApp("test")

			object := tree.getStruct()
			_, ok := object.apps["test"]
			test.AssertEquals(ok, false, it)

			_, ok = object.appScopes["test"]
			test.AssertEquals(ok, false, it)

			_, ok = object.channelScopes["chan-test"]
			test.AssertEquals(ok, false, it)
		})

	suite.RunAll()
}

func Test_insprdMemory_UpdateApp(t *testing.T) {
	suite := test.NewSuite("Test_insprdMemory_UpdateApp", t)
	tree := NewInsprDTree()
	var app meta.App

	suite.BeforeEach(func() {
		app.Metadata.Name = "test"
		app.Metadata.Parent = "inspr"
		app.Spec.Apps = make([]*meta.App, 0)
		tree.CreateApp("inspr", &app)
	})

	suite.NewCase("Without error",
		func(it *testing.T) {
			var app2 meta.App
			app2.Metadata.Name = "test"
			app2.Metadata.Parent = "inspr"
			app2.Metadata.Reference = "ref01"

			tree.UpdateApp(&app2)
			test.AssertNotNil(tree, it)

			object := tree.getStruct()
			test.AssertEquals(object.apps["test"].Metadata.Reference, "ref01", it)
		})

	suite.RunAll()
}

func Test_insprdMemory_recoverChannelTarget(t *testing.T) {
	suite := test.NewSuite("Test_insprdMemory_recoverChannelTarget", t)
	tree := NewInsprDTree()

	suite.BeforeEach(func() {
		tree.getStruct().aliasChannel["test"] = "test2"
		tree.getStruct().aliasChannel["test2"] = "test3"
		tree.getStruct().aliasChannel["test3"] = "test4"
	})

	suite.NewCase("Without error",
		func(it *testing.T) {
			object := tree.getStruct()

			target := object.recoverChannelTarget("test")
			test.AssertEquals(target, "test4", it)

			target = object.recoverChannelTarget("test2")
			test.AssertEquals(target, "test4", it)

			target = object.recoverChannelTarget("test3")
			test.AssertEquals(target, "test4", it)
		})

	suite.NewCase("Without error",
		func(it *testing.T) {
			tree.getStruct().aliasChannel["test4"] = ""
			object := tree.getStruct()

			target := object.recoverChannelTarget("test")
			test.AssertEquals(target, "test4", it)

			target = object.recoverChannelTarget("test2")
			test.AssertEquals(target, "test4", it)

			target = object.recoverChannelTarget("test3")
			test.AssertEquals(target, "test4", it)
		})

	suite.RunAll()
}

func Test_insprdMemory_CreateAliasChannel(t *testing.T) {
	suite := test.NewSuite("Test_insprdMemory_CreateAliasChannel", t)
	tree := NewInsprDTree()
	var app meta.App

	suite.BeforeEach(func() {
		app.Metadata.Name = "test"
		app.Metadata.Parent = "inspr"
		app.Spec.Apps = make([]*meta.App, 0)
		tree.CreateApp("inspr", &app)

		var channel meta.Channel
		channel.Meta.Name = "chan-test"
		tree.CreateChannel("test", &channel)

		var channel2 meta.Channel
		channel2.Meta.Name = "chan-test2"
		tree.CreateChannel("test", &channel2)

	})

	suite.NewCase("Without error",
		func(it *testing.T) {
			tree.CreateAliasChannel("chan-test", "chan-test2")

			object := tree.getStruct()
			test.AssertEquals(len(object.apps["test"].Spec.Channels), 1, it)
			test.AssertEquals(object.apps["test"].Spec.Channels[0].Meta.Name, "chan-test2", it)
		})

	suite.RunAll()
}

func Test_insprdMemory_DeleteAliasChannel(t *testing.T) {
	suite := test.NewSuite("Test_insprdMemory_DeleteAliasChannel", t)
	tree := NewInsprDTree()
	var app meta.App

	suite.BeforeEach(func() {
		app.Metadata.Name = "test"
		app.Metadata.Parent = "inspr"
		app.Spec.Apps = make([]*meta.App, 0)
		tree.CreateApp("inspr", &app)

		var channel meta.Channel
		channel.Meta.Name = "chan-test"
		tree.CreateChannel("test", &channel)

		var channel2 meta.Channel
		channel2.Meta.Name = "chan-test2"
		tree.CreateChannel("test", &channel2)

	})

	suite.NewCase("Without error",
		func(it *testing.T) {
			tree.CreateAliasChannel("chan-test", "chan-test2")
			tree.DeleteAliasChannel("chan-test")

			object := tree.getStruct()
			_, ok := object.aliasChannel["chan-test"]
			test.AssertEquals(ok, false, it)
		})

	suite.RunAll()
}
