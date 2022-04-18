package constants

import (
	"strings"

	"github.com/dropbox/godropbox/container/set"
)

const (
	DockerOrg = "pacur/"
)

var (
	Releases = [...]string{
		"archlinux",
		"amazonlinux-1",
		"amazonlinux-2",
		"fedora-35",
		"fedora-36",
		"centos-7",
		"centos-8",
		"debian-stretch",
		"debian-buster",
		"debian-bullseye",
		"oraclelinux-7",
		"oraclelinux-8",
		"ubuntu-trusty",
		"ubuntu-xenial",
		"ubuntu-bionic",
		"ubuntu-focal",
		"ubuntu-hirsute",
		"ubuntu-impish",
	}
	ReleasesMatch = map[string]string{
		"archlinux":       "",
		"amazonlinux-1":   ".amzn1.",
		"amazonlinux-2":   ".amzn2.",
		"fedora-35":       ".fc35.",
		"fedora-36":       ".fc36.",
		"centos-7":        ".el7.centos.",
		"centos-8":        ".el8.centos.",
		"debian-stretch":  ".stretch_",
		"debian-buster":   ".buster_",
		"debian-bullseye": ".bullseye_",
		"oraclelinux-7":   ".el7.oraclelinux.",
		"oraclelinux-8":   ".el8.oraclelinux.",
		"ubuntu-trusty":   ".trusty_",
		"ubuntu-xenial":   ".xenial_",
		"ubuntu-bionic":   ".bionic_",
		"ubuntu-focal":    ".focal_",
		"ubuntu-hirsute":  ".hirsute_",
		"ubuntu-impish":   ".impish_",
	}
	DistroPack = map[string]string{
		"archlinux":   "pacman",
		"amazonlinux": "redhat",
		"fedora":      "redhat",
		"centos":      "redhat",
		"debian":      "debian",
		"oraclelinux": "redhat",
		"ubuntu":      "debian",
	}
	Packagers = [...]string{
		"apt",
		"pacman",
		"yum",
	}

	ReleasesSet    = set.NewSet()
	Distros        = []string{}
	DistrosSet     = set.NewSet()
	DistroPackager = map[string]string{}
	PackagersSet   = set.NewSet()
	CleanPrevious  = false
)

func init() {
	for _, release := range Releases {
		ReleasesSet.Add(release)
		distro := strings.Split(release, "-")[0]
		Distros = append(Distros, distro)
		DistrosSet.Add(distro)
	}

	for _, distro := range Distros {
		packager := ""

		switch DistroPack[distro] {
		case "debian":
			packager = "apt"
		case "pacman":
			packager = "pacman"
		case "redhat":
			packager = "yum"
		default:
			panic("Failed to find packager for distro")
		}

		DistroPackager[distro] = packager
	}

	for _, packager := range Packagers {
		PackagersSet.Add(packager)
	}
}
