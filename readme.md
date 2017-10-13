## Ignore list implementation with the go lang
For more information read description in the [source file](ignore/IgnoreList.go)

## Examples
```go
list := ignore.NewListFromFile("my-ignores")
list.AddPattern("folder1/A")
list.AddPattern("folder1/B")
list.AddPattern("[my-tag] folder1/C")
list.AddPattern("not folder1/D")
list.AddPattern("folder2/*")
list.AddPattern("!folder2/E")
list.AddPattern("!folder2/*ex")

if list.IsIgnored("folder1/A") {
    // do something
}

result, tag := list.IsIgnoredEx("folder1/C")
if result {
    // do something
    fmt.Println(tag)
}
```

## Installation
With go
```
go get github.com/steptosky/go-ignorelist/ignore
```
With [go dep](https://github.com/golang/dep)
```
dep ensure -add github.com/steptosky/go-ignorelist/ignore
```

## Dev Installation
Getting dependencies for testing
```
dep ensure
```