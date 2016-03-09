package main

import (
	"fmt"
	"log"
	"os"

	"github.com/codegangsta/cli"
	"github.com/coreos/go-semver/semver"
	"github.com/fatih/color"
	"github.com/octoblu/go-simple-etcd-client/etcdclient"
	De "github.com/tj/go-debug"
)

var debug = De.Debug("etcd-watch-key:main")

func main() {
	app := cli.NewApp()
	app.Name = "etcd-watch-key"
	app.Version = version()
	app.Action = run
	app.ArgsUsage = "<key>"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "etcd-uri, e",
			EnvVar: "ETCD_WATCH_KEY_ETCD_URI",
			Usage:  "Etcd uri to watch",
		},
	}
	app.Run(os.Args)
}

func run(context *cli.Context) {
	etcdURI, watchKey := getOpts(context)

	client, err := etcdclient.Dial(etcdURI)
	if err != nil {
		log.Fatalln("Error on etcdclient.Dial", err.Error())
	}

	err = client.WatchRecursive(watchKey, func(key, value string) {
		fmt.Printf("key: %v, value: %v\n", key, value)
		os.Exit(0)
	})
	if err != nil {
		log.Fatalln("Error on client.WatchRecursive", err.Error())
	}
}

func getOpts(context *cli.Context) (string, string) {
	etcdURI := context.String("etcd-uri")
	key := context.Args().First()

	if etcdURI == "" || key == "" {
		cli.ShowAppHelp(context)

		if etcdURI == "" {
			color.Red("  Missing required flag --etcd-uri or ETCD_WATCH_KEY_ETCD_URI")
		}
		if key == "" {
			color.Red("  Missing required argument <key>")
		}
		os.Exit(1)
	}

	return etcdURI, key
}

func version() string {
	version, err := semver.NewVersion(VERSION)
	if err != nil {
		errorMessage := fmt.Sprintf("Error with version number: %v", VERSION)
		log.Panicln(errorMessage, err.Error())
	}
	return version.String()
}
