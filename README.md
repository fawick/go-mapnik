go-mapnik
=========

Go bindings for mapnik (http://www.mapnik.org or
http://github.com/mapnik/mapnik)

These bindings rely on http://github.com/fawick/mapnik-c-api, which was forked
from http://github.com/springmeyer/mapnik-c-api. (Pull request pending)

Installation
-----------

### Linux / MacOS

1. Install Mapnik build environmnent,
	- e.g. Ubuntu: `apt-get install libmapnik-dev`
2. Download the sources, either by 
- `git clone` the repository to $GOPATH/src/fawick/go-mapnik 

	OR

- `go get -d github.com/fawick/go-mapnik`
3. `cd mapnik` and run the configuration script `./configure.bash`. 
   That script will setup the correct paths for including Mapnik headers and
   linking against the Mapnik shared library, as well as download the Mapnik C
   API source 
4. `go install` the bindings in `./mapnik` or `cd example && go run example.go`



### Windows

On Windows, go-mapnik is restricted to GOARCH=386 for the moment, as the SDK is
build with MSVC 32-bit and the C api must be build with a compatible
compiler, before the Go bindings can be linked against them.

So, make sure `go version` reports `windows/386` as your toolchain. If you
happen to have the `windows/amd64` toolchain, either download Go binaries for
`windows/386` or put a MinGW 32-bit compiler in your path and rebuild the Go
binaries on your own.


1. Install Mapnik Windows 32bit SDK from http://www.mapnik.org/download
2. Download the sources, either by 
- `git clone` the repository to $GOPATH/src/fawick/go-mapnik 

	OR

- `go get -d github.com/fawick/go-mapnik`
3. Open a MSVC command prompt and run `configure.bat` to compile a C library
   that can be used by Go/CGO/GCC later (sources will be downloaded
   automatically)
4. `go install` the bindings in `./mapnik` or `cd example && go run example.go`
    (Make sure your %PATH% environment variable contains the path of
    `mapnik.dll`.)

Usage
-----

See `example/example.go` for a short usage example.

