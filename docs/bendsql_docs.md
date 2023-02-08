## What is bendsql?

bendsql is a handy command line interface tool to work with Databend Cloud smoothly and efficiently as with a web browser. Use the tool as an alternative to manage your warehouses and run SQL queries.


## Install bendsql

### brew install

If you're on MacOS, install bendsql using 'brew install':
```shell
brew tap databendcloud/homebrew-tap && brew install bendsql
```

### binary install

Go to the [Releases](https://github.com/databendcloud/bendsql/releases/latest) page and find a binary package for your platform.

### go install

```shell
go install github.com/databendcloud/bendsql/cmd/bendsql@latest
```

## Connect bendsql to Databend Cloud

### Sign in to Your Databend Cloud Account

```shell
bendsql cloud login
```

![](https://tva3.sinaimg.cn/large/005UfcOkly8h78cbw42jcj30z80b0aat.jpg)

If you don't have an account yet, create one in Databend Cloud.
Signing into your Databend Cloud account in bendsql requires your organization's information. You can press Enter to select the default organization during the sign-in process and then change it afterwards with the command ` bendsql cloud configure --org <your_org> `.

### Manage Warehouses

bendsql provides a bunch of commands to work with warehouses in Databend Cloud.

```shell
USAGE
  bendsql cloud warehouse cmd [flags]
CORE COMMANDS
  create:      Create a warehouse
  delete:      Delete a warehouse
  ls:          show warehouse list
  resume:      Resume a warehouse
  status:      show warehouse status
  suspend:     Suspend a warehouse
INHERITED FLAGS
  --help   Show help for command
LEARN MORE
  Use 'bendsql cloud <command> <subcommand> --help' for more information about a command.
```


## Work with bendsql

### Run SQL Queries

You can run SQL queries with bendsql. Specify a large warehouse for queries that need more computing resources.

```shell
echo 'YOURSQL;' | bendsql query --warehouse YOURWAREHOUSAE
```

### Run Interactive Shell

You can get an interractive database shell powered by [usql](https://github.com/xo/usql) with bendsql.

```shell
bendsql query
```

### Do More with bendsql

Type `bendsql -h` and discover more useful commands to make your work easier.
