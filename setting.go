package stdcli

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func (e *Engine) SettingDelete(name string) error {
	file, err := e.settingFile(name)
	if err != nil {
		return err
	}

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return nil
	}

	if err := os.Remove(file); err != nil {
		return err
	}

	return nil
}

func (e *Engine) LocalSetting(name string) string {
	file := filepath.Join(e.localSettingDir(), name)

	if _, err := os.Stat(file); os.IsNotExist(err) {
		return ""
	}

	data, err := ioutil.ReadFile(file)
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(data))
}

func (e *Engine) SettingRead(name string) (string, error) {
	file, err := e.settingFile(name)
	if err != nil {
		return "", err
	}

	data, err := ioutil.ReadFile(file)
	if os.IsNotExist(err) {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(data)), nil
}

func (e *Engine) SettingWrite(name, value string) error {
	file, err := e.settingFile(name)
	if err != nil {
		return err
	}

	dir := filepath.Dir(file)

	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	if err := ioutil.WriteFile(file, []byte(value), 0600); err != nil {
		return err
	}

	return nil
}
