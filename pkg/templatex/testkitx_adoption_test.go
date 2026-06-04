package templatex

import (
	"strings"
	"testing"

	"github.com/ZoneCNH/testkitx/pkg/testkitx/assertx"
	"github.com/ZoneCNH/testkitx/pkg/testkitx/fixture"
)

func TestTestkitxFixtureAdoption(t *testing.T) {
	workspace := fixture.NewWorkspace(t, "github.com/ZoneCNH/xlib-standard/downstream-fixture")
	path := workspace.Write("README.md", []byte("ok\n"))

	if !strings.HasSuffix(path, "README.md") {
		t.Fatalf("expected fixture path to end with README.md, got %q", path)
	}
	assertx.Equal(t, "off", workspace.Env["GOWORK"])
}
