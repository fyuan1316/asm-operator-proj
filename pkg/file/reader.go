package file

import "io/ioutil"

func LoadFromPath(path string) (string, error) {
	bytes, err := ioutil.ReadFile(path)
	return string(bytes), err
}
