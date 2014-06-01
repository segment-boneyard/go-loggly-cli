package main

import "github.com/segmentio/go-loggly-search"
import . "github.com/bitly/go-simplejson"
import "strings"
import "flag"
import "fmt"
import "os"

//
// Usage information.
//

const usage = `
  Usage: loggly [options] [query...]

  Options:

    --account <name>   account name
    --user <name>      account username
    --pass <word>      account password
    --size <count>     response event count [100]
    --from <time>      starting time [-24h]
    --to <time>        ending time [now]
    --path <str>       output json fields in <path>
    --count            output total event count
`

//
// Command options.
//

var flags = flag.NewFlagSet("loggly", flag.ExitOnError)
var count = flags.Bool("count", false, "")
var account = flags.String("account", "", "")
var user = flags.String("user", "", "")
var pass = flags.String("pass", "", "")
var size = flags.Int("size", 100, "")
var from = flags.String("from", "-24h", "")
var path = flags.String("path", "", "")
var to = flags.String("to", "now", "")

//
// Print usage and exit.
//

func printUsage() {
	fmt.Println(usage)
	os.Exit(0)
}

//
// Assert with msg.
//

func assert(ok bool, msg string) {
	if !ok {
		fmt.Printf("\n  Error: %s\n\n", msg)
		os.Exit(1)
	}
}

//
// Check error.
//

func check(err error) {
	if err != nil {
		fmt.Printf("\n  Error: %s\n\n", err)
		os.Exit(1)
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

	args := flags.Args()
	query := strings.Join(args, " ")

	c := search.New(*account, *user, *pass)

	// --count
	if *count {
		res, err := c.Query(query).Size(1).From(*from).To(*to).Fetch()
		check(err)
		fmt.Println(res.Total)
		os.Exit(0)
	}

	res, err := c.Query(query).Size(*size).From(*from).To(*to).Fetch()
	check(err)

	// --path
	if *path != "" {
		outputPath(res.Events, *path)
		os.Exit(0)
	}

	outputJson(res.Events)
}

//
// Output as json.
//

func outputJson(events []interface{}) {
	for _, event := range events {
		msg := event.(map[string]interface{})["logmsg"].(string)
		fmt.Println(msg)
	}
}

//
// Output path as json.
//

func outputPath(events []interface{}, path string) {
	for _, event := range events {
		msg := event.(map[string]interface{})["logmsg"].(string)

		obj, err := NewJson([]byte(msg))
		check(err)

		b, err := obj.GetPath(path).Encode()
		check(err)

		fmt.Println(string(b))
	}
}

// func output(event interface{}) {
// 	msg := event.(map[string]interface{})["logmsg"].(string)
// 	obj, err := NewJson([]byte(msg))
// 	check(err)

// 	fmt.Println()
// 	for k, v := range obj.MustMap() {
// 		fmt.Printf("  \033[36m%14s\033[0m \033[90m:\033[0m %s\n", k, v)
// 	}
// }
