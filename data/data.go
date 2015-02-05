package data

import (
	"github.com/trenker/boxserver/log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var data *Data
var prefix string
var findComponents *regexp.Regexp

func init() {
	divider := `\` + (string)(os.PathSeparator)
	validKey := `[a-z0-9][a-z0-9_\-]*[a-z0-9]`
	providers := []string{
		string(Virtualbox),
		string(Vmware),
		string(Docker),
		string(Hyperv),
	}
	allowedBoxes := `(` + strings.Join(providers, "|") + ")"
	validVersion := `[0-9]+\.[0-9]+\.[0-9]+`

	findComponents = regexp.MustCompile(
		`^` +
			validKey +
			divider +
			validKey +
			divider +
			validVersion +
			divider +
			allowedBoxes +
			`\.box$`)
}

func readFile(path string, info os.FileInfo, err error) error {

	if err != nil {
		return err
	}

	AddFromPath(path, info)

	return nil
}

func AddFromPath(path string, info os.FileInfo) {

	if !info.IsDir() && len(path) > len(prefix) {

		path = strings.TrimPrefix(path, prefix)

		if findComponents.MatchString(path) {

			parts := strings.Split(strings.TrimSuffix(path, ".box"), "/")

			log.Debug("Found box %s", parts)

			providerName, _ := ProviderByName(parts[3])

			var p *Project
			var b *Box
			var v *Version

			log.Debug("Append %s", parts)

			p = data.getProject(parts[0])

			if p == nil {
				p = data.addProject(Project{Name: parts[0], Boxes: make([]Box, 0)})
			}

			b = p.getBox(parts[1])

			if b == nil {
				b = p.addBox(Box{Name: parts[1], Versions: make([]Version, 0)})
			}

			v = b.getVersion(parts[2])

			if v == nil {
				v = b.addVersion(Version{Version: parts[2], Providers: make([]VagrantProvider, 0)})
			}

			v.addProvider(providerName)
		}
	}
}

func LoadData() {
	data = new(Data)
	data.Projects = make([]Project, 0)

	filepath.Walk(basePath, readFile)
}

func Initialize(basePath string) *Data {
	prefix = basePath + string(os.PathSeparator)

	LoadData()

	return data
}
