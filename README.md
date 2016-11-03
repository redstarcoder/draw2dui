draw2dui
======
[![GoDoc](https://godoc.org/github.com/redstarcoder/draw2dui?status.svg)](https://godoc.org/github.com/redstarcoder/draw2dui)

**\*\*NOTICE\*\*** This library is very young, and *may* be subject to API breaking changes. It should stabilize within the next few weeks. **\*\*NOTICE\*\***

Package draw2dui offers useful tools for drawing and handling UIs in Golang using [draw2d](https://github.com/llgcode/draw2d) with OpenGL.

Installation
---------------

Install [golang](http://golang.org/doc/install). To install or update the package draw2dui on your system, run:

```
go get -u github.com/redstarcoder/draw2dui
```

Building
---------------

These are instructions for building on Linux. This library should be able to compile on whatever you can make gl and glfw compile on though.

**Target: Linux**

```
go build
```

**Target: Windows**

First install mingw, then run something like this:

```
CC=x86_64-w64-mingw32-gcc CGO_ENABLED=1 GOOS=windows go build
```


Acknowledgments
---------------

* [redstarcoder](https://github.com/redstarcoder) wrote this library.
* [Laurent Le Goff](https://github.com/llgcode) wrote draw2d.

