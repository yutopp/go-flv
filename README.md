# go-flv

[![CircleCI](https://circleci.com/gh/yutopp/go-flv.svg?style=svg)](https://circleci.com/gh/yutopp/go-flv)
[![Coverage Status](https://coveralls.io/repos/github/yutopp/go-flv/badge.svg)](https://coveralls.io/github/yutopp/go-flv)
[![GoDoc](https://godoc.org/github.com/yutopp/go-flv?status.svg)](http://godoc.org/github.com/yutopp/go-flv)
[![Go Report Card](https://goreportcard.com/badge/github.com/yutopp/go-flv)](https://goreportcard.com/report/github.com/yutopp/go-flv)
[![license](https://img.shields.io/github/license/yutopp/go-flv.svg)](https://github.com/yutopp/go-flv/blob/master/LICENSE_1_0.txt)

FLV decoder/encoder library written in Go.

- [x] decoder
  - [x] header
  - [x] body
  - [x] tags
    - [x] flv
    - [x] audio
    - [x] video
    - [x] data
- [x] encoder
  - [x] header
  - [x] body
  - [x] tags
    - [x] flv
    - [x] audio
    - [x] video
    - [x] data
  
## Installation

```
go get github.com/yutopp/go-flv
```

## Examples

- [yutopp/go-flv-examples](https://github.com/yutopp/go-flv-examples)

## License

[Boost Software License - Version 1.0](./LICENSE_1_0.txt)

## References

- [spec](https://wwwimages2.adobe.com/content/dam/acom/en/devnet/flv/video_file_format_spec_v10.pdf)
  - The FLV File Format
