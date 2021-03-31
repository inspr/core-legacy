package tree

import (
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gitlab.inspr.dev/inspr/core/pkg/meta"
)

func CompareWithoutUUID(first, second interface{}) bool {
	opt := cmp.Comparer(func(a, b meta.Metadata) bool {
		a.UUID = ""
		b.UUID = ""
		return cmp.Equal(a, b)
	})
	return cmp.Equal(first, second, opt, cmp.AllowUnexported(MemoryManager{}))
}
func (m *MemoryManager) Equals(other *MemoryManager) bool {
	return CompareWithUUID(m.root, other.root) && CompareWithUUID(m.tree, other.tree)
}

func (m *MemoryManager) Diff(other *MemoryManager) string {
	return cmp.Diff(m.root, other.root, cmp.AllowUnexported(MemoryManager{})) + "\n" + cmp.Diff(m.tree, other.tree, cmp.AllowUnexported(MemoryManager{}))
}

func CompareWithUUID(first, second interface{}) bool {

	return cmp.Equal(first, second)
}

func ValidateUUID(uuid string) bool {
	valRegex := regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$")
	return valRegex.MatchString(uuid)
}

func recursiveValidateUUIDS(name string, app *meta.App, t *testing.T) {
	if !ValidateUUID(app.Meta.UUID) {
		t.Errorf("%s invalid UUID on %s, uuid=%v", name, app.Meta.Name, app.Meta.UUID)
	}
	for _, a := range app.Spec.Apps {

		recursiveValidateUUIDS(name, a, t)
	}
	for _, c := range app.Spec.Channels {
		if !ValidateUUID(c.Meta.UUID) {
			t.Errorf("%s invalid channel UUID on %s, uuid = %v", name, c.Meta.Name, c.Meta.UUID)
		}
	}

	for _, ct := range app.Spec.ChannelTypes {
		if !ValidateUUID(ct.Meta.UUID) {
			t.Errorf("%s invalid channel type UUID on %s, uuid = %v", name, ct.Meta.Name, ct.Meta.UUID)
		}
	}

	for _, a := range app.Spec.Aliases {
		if !ValidateUUID(a.Meta.UUID) {
			t.Errorf("%s invalid alias UUID on %s, uuid = %v", name, a.Meta.Name, a.Meta.UUID)
		}
	}

}
