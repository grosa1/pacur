package source

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"github.com/dropbox/godropbox/errors"
	"github.com/pacur/pacur/utils"
	"hash"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	path = 0
	url  = 1
)

type Source struct {
	Hash   string
	Source string
	Output string
	Path   string
}

func (s *Source) getType() int {
	if strings.HasPrefix(s.Source, "http") {
		return url
	}
	return path
}

func (s *Source) getUrl() (err error) {
	name, err := utils.UrlFilename(s.Source)
	if err != nil {
		return
	}

	s.Path = filepath.Join(s.Output, name)

	exists, err := utils.Exists(s.Path)
	if err != nil {
		return
	}

	if !exists {
		err = utils.HttpGet(s.Source, s.Output)
		if err != nil {
			return
		}
	}

	return
}

func (s *Source) getPath() (err error) {
	return
}

func (s *Source) extract() (err error) {
	cmd := exec.Command("tar", "xfz", s.Path)
	cmd.Dir = s.Output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err = cmd.Run()
	if err != nil {
		err = &GetError{
			errors.Wrapf(err, "builder: Failed to extract source '%s'",
				s.Source),
		}
		return
	}

	return
}

func (s *Source) validate() (err error) {
	if strings.ToLower(s.Hash) == "skip" {
		return
	}

	file, err := os.Open(s.Path)
	if err != nil {
		err = &HashError{
			errors.Wrap(err, "source: Failed to open file for hash"),
		}
		return
	}
	defer file.Close()

	var hash hash.Hash
	switch len(s.Hash) {
	case 32:
		hash = md5.New()
	case 40:
		hash = sha1.New()
	case 64:
		hash = sha256.New()
	case 128:
		hash = sha512.New()
	default:
		err = &HashError{
			errors.Newf("source: Unknown hash type for hash '%s'", s.Hash),
		}
		return
	}

	_, err = io.Copy(hash, file)
	if err != nil {
		return
	}

	sum := hash.Sum([]byte{})

	hexSum := fmt.Sprintf("%x", sum)

	if hexSum != s.Hash {
		err = &HashError{
			errors.Newf("source: Hash verification failed for '%s'", s.Source),
		}
		return
	}

	return
}

func (s *Source) Get() (err error) {
	switch s.getType() {
	case url:
		err = s.getUrl()
	case path:
		err = s.getPath()
	default:
		panic("utils: Unknown type")
	}
	if err != nil {
		return
	}

	err = s.validate()
	if err != nil {
		return
	}

	err = s.extract()
	if err != nil {
		return
	}

	return
}
