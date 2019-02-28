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
- [ ] Set values
- [ ] Write to file
- [ ] Tests