package main

// import "github.com/segmentio/go-loggly-search"
// import . "github.com/bitly/go-simplejson"
// import "strings"
import "flag"
import "fmt"
import "os"

//
// Usage information.
//

const usage = `
  Usage: loggly [options] [query]

  Options:

    --account <name>   account name
    --user <name>      account username
    --pass <word>      account password
    --size <count>     response event count [100]
    --from <time>      starting time [-24h]
    --to <time>        ending time [now]
`

//
// Command options.
//

var flags = flag.NewFlagSet("loggly", flag.ExitOnError)
var account = flags.String("account", "", "account name")
var user = flags.String("user", "", "account username")
var pass = flags.String("pass", "", "account password")
var size = flags.Int64("size", 100, "response event count")
var from = flags.String("from", "-24h", "starting time")
var to = flags.String("to", "now", "ending time")

//
// Print usage and exit.
//

func printUsage() {
	fmt.Println(usage)
	os.Exit(0)
}

func assert(ok bool, msg string) {
	if !ok {
		fmt.Printf("\n  Error: %s\n\n", msg)
		os.Exit(10)
	}
}

//
// Main.
//

func main() {
	flags.Usage = printUsage
	flags.Parse(os.Args[1:])

	assert(*account != "", "--account required")
	assert(*user != "", "--user required")
	assert(*pass != "", "--pass required")

}
