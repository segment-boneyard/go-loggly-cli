package main

import "github.com/segmentio/go-loggly-search"
import . "github.com/segmentio/go-simplejson"
import "github.com/jehiah/go-strftime"
import "strings"
import "flag"
import "time"
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
    --json             output json array of events
    --count            output total event count
`

//
// Command options.
//

var flags = flag.NewFlagSet("loggly", flag.ExitOnError)
var count = flags.Bool("count", false, "")
var json = flags.Bool("json", false, "")
var account = flags.String("account", "", "")
var user = flags.String("user", "", "")
var pass = flags.String("pass", "", "")
var size = flags.Int("size", 100, "")
var from = flags.String("from", "-24h", "")
var to = flags.String("to", "now", "")

//
// Colors.
//

var colors = map[string]string{
	"debug":     "90",
	"info":      "90",
	"notice":    "33",
	"warning":   "33",
	"critical":  "31",
	"alert":     "31;1",
	"emergency": "31;1",
}

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

	// --json
	if *json {
		outputJson(res.Events)
		os.Exit(0)
	}

	// formatted
	output(res.Events)
}

//
// Output as json.
//

func outputJson(events []interface{}) {
	fmt.Println("[")

	l := len(events)

	for i, event := range events {
		msg := event.(map[string]interface{})["logmsg"].(string)
		if i < l-1 {
			fmt.Printf("  %s,\n", msg)
		} else {
			fmt.Printf("  %s\n", msg)
		}
	}

	fmt.Println("]")
}

func timeFromUnix(ms int64) time.Time {
	return time.Unix(0, ms*int64(time.Millisecond))
}

//
// Formatted output.
//

func output(events []interface{}) {
	for _, event := range events {
		msg := event.(map[string]interface{})["logmsg"].(string)
		obj, err := NewJson([]byte(msg))
		check(err)

		host := obj.Get("hostname").MustString()
		level := obj.Get("level").MustString()
		ts := timeFromUnix(obj.Get("timestamp").MustInt64())
		t := obj.Get("type").MustString()
		c := colors[level]

		obj.Del("hostname")
		// obj.Del("level")
		// obj.Del("timestamp")
		// obj.Del("type")

		json, err := obj.EncodePretty()
		check(err)

		date := strftime.Format("%m-%d %I:%M:%S %p", ts)
		level = strings.ToUpper(level)
		fmt.Printf("\n%s \033["+c+"m%s\033[0m - \033[36m%s \033[90m(%s)\033[0m\n", date, level, t, host)
		fmt.Printf("\n%s\n", string(json))
	}

	fmt.Println()
}
