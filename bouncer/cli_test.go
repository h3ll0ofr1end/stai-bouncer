package bouncer

import (
	"errors"
	"io/fs"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

const steffenHome = "/home/steffen"

func noHomeDir() (string, error) {
	return "", errors.New("no home dir")
}
func homeDir() (string, error) {
	return steffenHome, nil
}

func exists(_ string) error {
	return nil
}

func doesNotExist(_ string) error {
	return errors.New("does not exist")
}

func Test_defaultStaiExecutable(t *testing.T) {
	tests := []struct {
		name            string
		userHomeDirFunc func() (string, error)
		want            string
		wantErr         bool
	}{
		{
			name:            "home exists",
			userHomeDirFunc: homeDir,
			want:            steffenHome + "/" + DefaultStaiExecutableSuffix,
		},
		{
			name:            "home does not exist",
			userHomeDirFunc: noHomeDir,
			wantErr:         true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			getUserHomeDir = tt.userHomeDirFunc
			got, err := defaultStaiExecutable()
			if (err != nil) != tt.wantErr {
				t.Errorf("defaultStaiExecutable() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("defaultStaiExecutable() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_run(t *testing.T) {
	enforceStaiExecutable = func(staiExecutable string) error {
		return nil
	}

	tests := []struct {
		name    string
		args    []string
		homeDir func() (string, error)
		exists  func(string) error
		want    *Context
		wantErr string
	}{
		{
			name:    "ok with default stai",
			args:    []string{"stai-bouncer", "mars"},
			homeDir: homeDir,
			exists:  exists,
			want: &Context{
				StaiExecutable: steffenHome + "/" + DefaultStaiExecutableSuffix,
				Location:       "mars",
			}},
		{
			name:    "ok with custom stai",
			args:    []string{"stai-bouncer", "-e", "/home/something/else/stai", "elon on mars"},
			homeDir: homeDir,
			exists:  exists,
			want: &Context{
				StaiExecutable: "/home/something/else/stai",
				Location:       "elon on mars",
			}},
		{
			name:    "nok with no location",
			args:    []string{"stai-bouncer"},
			wantErr: "LOCATION is missing",
		},
		{
			name:    "nok with default home not exist",
			args:    []string{"stai-bouncer", "elon on mars"},
			homeDir: noHomeDir,
			wantErr: "no home dir",
		},
		{
			name:    "nok with custom stai executable not exist",
			args:    []string{"stai-bouncer", "elon on mars"},
			homeDir: homeDir,
			exists:  doesNotExist,
			wantErr: "does not exist",
		},
	}
	for _, tt := range tests {
		args = tt.args
		getUserHomeDir = tt.homeDir
		enforceStaiExecutable = tt.exists
		t.Run(tt.name, func(t *testing.T) {
			got, err := RunCli()
			if err != nil {
				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("Run() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Run() got = %v, want %v", got, tt.want)
			}
		})
	}
}

type DummyFileInfo struct {
	directory bool
}

func (d DummyFileInfo) Name() string {
	panic("implement me")
}

func (d DummyFileInfo) Size() int64 {
	panic("implement me")
}

func (d DummyFileInfo) Mode() fs.FileMode {
	panic("implement me")
}

func (d DummyFileInfo) ModTime() time.Time {
	panic("implement me")
}

func (d DummyFileInfo) Sys() interface{} {
	panic("implement me")
}

func (d DummyFileInfo) IsDir() bool {
	return d.directory
}

func fileExists(_ string) (os.FileInfo, error) {
	return DummyFileInfo{}, nil
}

func fileDoesNotExist(_ string) (os.FileInfo, error) {
	return nil, fs.ErrNotExist
}

func fileIsDirectory(_ string) (os.FileInfo, error) {
	return DummyFileInfo{
		directory: true,
	}, nil
}

func Test_enforceExists(t *testing.T) {
	tests := []struct {
		name           string
		staiExecutable string
		fileinfo       func(name string) (os.FileInfo, error)
		wantErr        string
	}{
		{name: "ok", staiExecutable: "stai", fileinfo: fileExists},
		{name: "nok file does not exist", staiExecutable: "stai", fileinfo: fileDoesNotExist, wantErr: "stai executable does not exist"},
		{name: "nok file is a directory", staiExecutable: "stai", fileinfo: fileIsDirectory, wantErr: "stai executable can not be a directory"},
	}
	for _, tt := range tests {
		getFileInfo = tt.fileinfo
		t.Run(tt.name, func(t *testing.T) {
			err := enforceExists(tt.staiExecutable)
			if err != nil {
				if tt.wantErr == "" {
					t.Errorf("enforceExists() error = %v, wantErr %v", err, tt.wantErr)
					return
				}

				if !strings.Contains(err.Error(), tt.wantErr) {
					t.Errorf("enforceExists() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
			}
		})
	}
}
