package gobporeader

import (
	"errors"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

// BpoReader is a Config manager class
type BpoReader struct {
	Path    string
	Content []byte
	Values  map[string]string
}

// WriteToFile writes Content to a file
func (c *BpoReader) WriteToFile(filePath string) error {
	err := ioutil.WriteFile(filePath, c.Content, 0755)
	if err != nil {
		return err
	}
	return nil
}

func (c *BpoReader) parse() error {
	var re = regexp.MustCompile(`(?m)(#.*)$`)
	c.Content = re.ReplaceAll(c.Content, []byte(""))
	re = regexp.MustCompile(`(?m)^(\s*)$`)
	c.Content = re.ReplaceAll(c.Content, []byte(""))
	c.Values = make(map[string]string)

	re = regexp.MustCompile(`(?m)^\s*(\[.+\])\s*$`)
	res := re.FindAllIndex(c.Content, -1)
	for i := 0; i < len(res)-1; i++ {
		sName := regexp.MustCompile(`\[(.+)\]`).FindSubmatch(c.Content[res[i][0]:res[i][1]])
		currentPrefix := string(sName[1])
		currentSlice := c.Content[res[i][0]:res[i+1][1]]

		reKeyVals := regexp.MustCompile(`(?m)^(.+)\s*=\s*(.+)$`)
		keyValsRes := reKeyVals.FindAllIndex(currentSlice, -1)
		for j := 0; j < len(keyValsRes); j++ {
			keyValsSubmatch := regexp.MustCompile(`(?m)^(.+)\s*=\s*(.+)$`).FindSubmatch(currentSlice[keyValsRes[j][0]:keyValsRes[j][1]])
			if len(keyValsSubmatch) < 3 {
				continue
			}
			currentKey := strings.Trim(currentPrefix+"."+string(keyValsSubmatch[1]), " \n")
			currentVal := strings.Trim(string(keyValsSubmatch[2]), " \n")
			c.Values[currentKey] = currentVal
		}
	}

	return nil
}

// ReadFromFile reads Content from a file
func (c *BpoReader) ReadFromFile(path string) error {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	c.Path = path
	c.Content = dat
	c.parse()
	return nil
}

// GetString gets a string value of a key
func (c *BpoReader) GetString(key string) string {
	if val, ok := c.Values[key]; ok {
		return val
	}
	return ""
}

// GetInt gets an integer value of a key
func (c *BpoReader) GetInt(key string) (*int64, error) {
	if val, ok := c.Values[key]; ok {
		intVal, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil, err
		}
		return &intVal, nil
	}
	return nil, errors.New("Key not found")
}
