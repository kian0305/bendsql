# bendsql

Work seamlessly with [Databend Cloud](https://app.databend.com/) from the command line.


## Installation

### brew
For macOS, user can use `brew tap databendcloud/homebrew-tap && brew install bendsql` to install.

### go install
`go install github.com/databendcloud/bendsql/cmd/bendsql@v0.0.2`, then the `bendsql` will be installed in `$GOPATH/bin`.

### From source code
Clone the repo and exec `make install`, it will install `bendsql` in `/usr/local/bin`

```shell
make install
```

### Binary

Visit `https://github.com/databendcloud/bendsql/releases/latest
` and download binary package according your arch.

Then `bendsql -h`, enjoy it!
