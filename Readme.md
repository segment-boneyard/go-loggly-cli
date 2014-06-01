
# Loggly CLI

  Loggly search command-line tool.

## Installation

## Usage

```

  Usage: loggly [options] [query...]

  Options:

    --account <name>   account name
    --user <name>      account username
    --pass <word>      account password
    --size <count>     response event count [100]
    --from <time>      starting time [-24h]
    --to <time>        ending time [now]
    --json             output json array of events
    --count            output total event count

```

## Setup

 Loggly's search API requires basic auth credentials, so you _must_ pass
 the `--acount`, `--user`, and `--pass` flags. To make this less annoying
 I suggest creating an alias:

```sh
alias logs='loggly --account segment --user tj --pass something'
```

 This is a great place to stick personal defaults as well. Since flags are clobbered
 if defined multiple times you can define whatever defaults you'd like here, while
 still changing them via `log`:

```sh
alias logs='loggly --account segment --user tj --pass something --size 5'
```

## License

 MIT