# Go BPO Reader
A class for Boost parameters option files parser for Golang

## Example
```
cfgManager := bporeader.BpoReader{}
err := cfgManager.ReadFromFile("/home/user/config.conf")
if err != nil {
    fmt.Println(err.Error())
}
fmt.Println(cfgManager.GetString("section.value"))
```

## TODO
- [x] Set string values
- [ ] Set values of different types
- [x] Write to file
- [ ] Write to a file keeping initial comments
- [ ] Remove keys
- [ ] Tests