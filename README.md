# npmirror - A lazy NPM mirror

WIP - Things will change!
**npmirror** is a local NPM cache/mirror to make your NPM installs faster.

## Install
1. [Install Go](http://golang.org/doc/install)
2. Run `go install github.com/sebastianm/npmirror`

## Usage

Run `npmirror` at least with the `-external-addr` flag:
```
Usage of npmirror:
  -external-addr="": Required! External HTTP address (registry address)
  -http-addr="127.0.0.1:8023": httpd bind address
  -storage-file-dir="./npmirror": Storage directory for all cached NPM files
  -storage-type="file": Storage Type (only 'file' available in this version)
```

### NPM Registry Configuration
You have to tell NPM to use your hosted **npmirror**. Replace [EXTERNAL_ADDR] with the address provided in `-external-addr` flag`.

To temporarily set the NPM registry:
```
npm --registry [EXTERNAL_ADDR] install angular
```

To permanently set the NPM registry:
```
npm config set registry [EXTERNAL_ADDR]
```

## Author & License
Sebastian Müller ([SebastianM](https://github.com/SebastianM) - [@Sebamueller](https://twitter.com/sebamueller))
License: [MIT](LICENSE)

## TODO
1. Add tests
2. Make registry upstream configurable
3. Better logging mechanism