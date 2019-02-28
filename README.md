# Go BPO Reader
A class for Boost parameters option files parser for Golang

Any value of the configuration file is stored as a string.

## Example
```
cfgManager := bporeader.BpoReader{}
err := cfgManager.ReadFromFile("/home/user/config.conf")
if err != nil {
    fmt.Println(err.Error())
}
fmt.Println(cfgManager.GetString("section.value"))
```

Writing to a file
```
cfgManager.SetString("section.value", "new_value")
cfgManager.WriteToFileClean("test.conf") // writes without keeping initial comments
```

## TODO
- [x] Set string values
- [ ] Set values of different types
- [x] Write to file
- [ ] Write to a file keeping initial comments
- [ ] Remove keys
- [ ] Tests