# ordx

[![build-img]][build-url]
[![pkg-img]][pkg-url]
[![reportcard-img]][reportcard-url]
[![coverage-img]][coverage-url]
[![version-img]][version-url]

## Features

* Simple API.
* Clean and tested code.
* Dependency-free.

## Install

Go version 1.25+

```
go get github.com/cristalhq/ordx
```

## Example

```
less := func(a, b int) bool { return a < b }
cmp := ordx.AsCmp(less)

fmt.Println(cmp(1, 2))
fmt.Println(cmp(2, 1))
fmt.Println(cmp(2, 2))

// Output:
// -1
// 1
// 0

cmp := ordx.RankCmp([]string{"low", "medium", "high"})

fmt.Println(cmp("low", "high"))
fmt.Println(cmp("high", "low"))
fmt.Println(cmp("medium", "medium"))

// Output:
// -1
// 1
// 0
```

See examples: [example_test.go](https://github.com/cristalhq/ordx/blob/main/example_test.go).

## Documentation

See [these docs][pkg-url].

## License

[MIT License](LICENSE).

[build-img]: https://github.com/cristalhq/ordx/workflows/build/badge.svg
[build-url]: https://github.com/cristalhq/ordx/actions
[pkg-img]: https://pkg.go.dev/badge/cristalhq/ordx
[pkg-url]: https://pkg.go.dev/github.com/cristalhq/ordx
[reportcard-img]: https://goreportcard.com/badge/cristalhq/ordx
[reportcard-url]: https://goreportcard.com/report/cristalhq/ordx
[coverage-img]: https://codecov.io/gh/cristalhq/ordx/branch/main/graph/badge.svg
[coverage-url]: https://codecov.io/gh/cristalhq/ordx
[version-img]: https://img.shields.io/github/v/release/cristalhq/ordx
[version-url]: https://github.com/cristalhq/ordx/releases
