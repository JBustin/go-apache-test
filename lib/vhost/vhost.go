package vhost

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-apache-test/lib/cst"
	"github.com/go-apache-test/lib/filesystem"
	"github.com/go-apache-test/lib/test"
)

var increment int = 0

type Vhost struct {
	OriPath      string
	NewPath      string
	TestPath     string
	Alias        []string
	Suites       []test.Suite
	Skip         bool
	Headers      map[string]string
	Expectations test.Expectations
}

func New(vhostPath string) Vhost {
	increment += 1
	return Vhost{
		OriPath:      vhostPath,
		NewPath:      fmt.Sprintf("%v%v%v%v", cst.HttpdVhostsDir, increment, "-", filepath.Base(vhostPath)),
		TestPath:     fmt.Sprintf("%v%v", cst.DockerTestsDir, strings.Replace(filepath.Base(vhostPath), ".conf", ".json", 1)),
		Alias:        []string{},
		Suites:       []test.Suite{},
		Skip:         false,
		Headers:      map[string]string{},
		Expectations: test.Expectations{},
	}
}

func (v *Vhost) GetAlias(filer filesystem.Filer) error {
	if !filer.Exists(v.OriPath) {
		return fmt.Errorf("test file %v cannot be loaded", v.OriPath)
	}

	vhostContent, err := filer.Read(v.OriPath)
	if err != nil {
		return err
	}

	re := regexp.MustCompile(`ServerName ([^\n]*)`)
	serverName := string(re.Find([]byte(vhostContent)))
	serverName = strings.Replace(serverName, "ServerName ", "", -1)

	re = regexp.MustCompile(`ServerAlias ([^\n]*)`)
	serverAliasStr := string(re.Find([]byte(vhostContent)))
	serverAliasStr = strings.Replace(serverAliasStr, "ServerAlias ", "", -1)
	serverAlias := strings.Split(serverAliasStr, " ")

	v.Alias = append(serverAlias, serverName)

	return nil
}

func (v *Vhost) GetSuites(filer filesystem.Filer) error {
	if !filer.Exists(v.TestPath) {
		return fmt.Errorf("test file %v cannot be loaded", v.TestPath)
	}

	jsonContent, err := filer.Read(v.TestPath)
	if err != nil {
		return err
	}

	return json.Unmarshal(jsonContent, &v)
}

func (v Vhost) String() string {
	return fmt.Sprintf(`
	oriPath		%v
	NewPath		%v
	TestPath	%v
	suites		
%v
	skip		%v
	headers		%v
	expectations%v
`,
		v.OriPath,
		v.NewPath,
		v.TestPath,
		v.Suites,
		v.Skip,
		v.Headers,
		v.Expectations,
	)
}
