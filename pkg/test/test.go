package test

import (
	"reflect"
	"testing"
)

type testSuite struct {
	name       string
	beforeAll  func()
	afterAll   func()
	beforeEach func()
	afterEach  func()
	testCases  map[string]func(t *testing.T)
	t          *testing.T
}

// Suite represents a test suite
type Suite interface {
	BeforeAll(func())
	BeforeEach(func())
	AfterEach(func())
	RunAll()
	NewCase(name string, f func(t *testing.T))
}

// NewSuite represents a new suite
func NewSuite(name string, t *testing.T) Suite {
	return &testSuite{
		t:         t,
		testCases: make(map[string]func(t *testing.T)),
	}
}

func (ts *testSuite) BeforeAll(beforeAll func()) {
	ts.beforeAll = beforeAll
}

func (ts *testSuite) BeforeEach(beforeEach func()) {
	ts.beforeEach = beforeEach
}
func (ts *testSuite) AfterEach(afterEach func()) {
	ts.afterEach = afterEach
}

func (ts *testSuite) NewCase(name string, c func(*testing.T)) {
	ts.testCases[name] = c
}

func (ts *testSuite) RunAll() {
	if ts.beforeAll != nil {
		ts.beforeAll()
	}
	for name, c := range ts.testCases {
		n := name
		cas := c
		if ts.beforeEach != nil {
			ts.beforeEach()
		}
		ts.t.Run(n, cas)
		if ts.afterEach != nil {
			ts.afterEach()
		}
	}
}

// AssertDeepEquals asserts that two interfaces are equal
func AssertDeepEquals(i1, i2 interface{}, t *testing.T) {
	if !reflect.DeepEqual(i1, i2) {
		t.Errorf("expected %#v to be equal to %#v", i1, i2)
	}
}

// AssertNotDeepEquals asserts that two interfaces are equal
func AssertNotDeepEquals(i1, i2 interface{}, t *testing.T) {
	if reflect.DeepEqual(i1, i2) {
		t.Errorf("expected %#v to be not equal to %#v", i1, i2)
	}
}

// AssertEquals asserts that two interfaces are equal
func AssertEquals(i1, i2 interface{}, t *testing.T) {
	if i1 != i2 {
		t.Errorf("expected %#v to be equal to %#v", i1, i2)
	}
}

// AssertNotEquals asserts that two interfaces are equal
func AssertNotEquals(i1, i2 interface{}, t *testing.T) {
	if i1 == i2 {
		t.Errorf("expected %#v to be not equal to %#v", i1, i2)
	}
}

// AssertNotNil asserts that an interface is not nil
func AssertNotNil(i1 interface{}, t *testing.T) {
	if i1 == nil {
		t.Errorf("expected %#v to be not nil", i1)
	}
}

// AssertNil asserts that an interface is not nil
func AssertNil(i1 interface{}, t *testing.T) {
	if i1 != nil {
		t.Errorf("expected %#v to be nil", i1)
	}
}
