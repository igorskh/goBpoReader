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
	Path           string
	Content        []byte
	InitialContent []byte
	Values         map[string]string
	Sections       []string
}

func getSectionName(fullName string) string {
	dotPos := strings.Index(fullName, ".")
	if dotPos > -1 {
		return fullName[0:dotPos]
	}
	return fullName
}

func getKeyName(fullName string) string {
	dotPos := strings.Index(fullName, ".")
	if dotPos > -1 {
		return fullName[dotPos+1 : len(fullName)]
	}
	return fullName
}

func (c *BpoReader) isSectionExists(sectionName string) int {
	for i, s := range c.Sections {
		if sectionName == s {
			return i
		}
	}

	return -1
}

// WriteToFile writes Content to a file
func (c *BpoReader) WriteToFile(filePath string) error {
	err := ioutil.WriteFile(filePath, c.Content, 0755)
	if err != nil {
		return err
	}
	return nil
}

// GenerateContentClean generates Content from the KeyValue map without initial comments
func (c *BpoReader) GenerateContentClean() string {
	mapToWrite := c.Values
	newContent := ""

	for _, sectionName := range c.Sections {
		newContent += "[" + sectionName + "]\n"
		for k, v := range mapToWrite {
			if getSectionName(k) == sectionName {
				newContent += getKeyName(k) + " = " + v + "\n"
			}
		}
	}
	return newContent
}

// WriteToFileClean writes Content to a file without initial comments
func (c *BpoReader) WriteToFileClean(filePath string) error {
	err := ioutil.WriteFile(filePath, []byte(c.GenerateContentClean()), 0755)
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

		c.Sections = append(c.Sections, currentPrefix)

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
	c.InitialContent = dat
	c.parse()
	return nil
}

// SetString sets a string value
func (c *BpoReader) SetString(key string, value string) {
	sectionName := getSectionName(key)
	if c.isSectionExists(sectionName) == -1 {
		c.Sections = append(c.Sections, sectionName)
	}
	c.Values[key] = value
}

// GetString gets a string value of a key
func (c *BpoReader) GetString(key string) (*string, error) {
	if val, ok := c.Values[key]; ok {
		return &val, nil
	}
	return nil, errors.New("Key not found")
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
