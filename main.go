package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/jezek/xgb"
	"github.com/jezek/xgb/randr"
	"github.com/jezek/xgb/xproto"
)

const defaultAutorandrPath = "/usr/bin/autorandr"

func main() {
	var autorandr string
	flag.StringVar(&autorandr, "path", defaultAutorandrPath, "Path to the autorandr binary")
	flag.StringVar(&autorandr, "p", defaultAutorandrPath, "Path to the autorandr binary")

	flag.Parse()

	stat, err := os.Stat(autorandr)
	if err != nil {
		log.Fatal(err)
	}

	if (stat.Mode().Perm() & 0111) == 0 {
		log.Fatal(fmt.Errorf("%s is not executable", autorandr))
	} else if stat.IsDir() {
		log.Fatal(fmt.Errorf("%s is a directory", autorandr))
	}

	X, err := xgb.NewConn()
	if err != nil {
		log.Fatal(err)
	}

	err = randr.Init(X)
	if err != nil {
		log.Fatal(err)
	}

	root := xproto.Setup(X).DefaultScreen(X).Root

	err = randr.SelectInputChecked(X, root, randr.NotifyMaskScreenChange).Check()
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Autorandr launcher started, waiting for events...")

	for {
		ev, err := X.WaitForEvent()
		if err != nil {
			log.Fatal(err)
		}
		switch eventType := ev.(type) {
		case randr.ScreenChangeNotifyEvent:
			output, err := runAutorandr(autorandr)
			if err != nil {
				log.Printf("Error calling autorandr %s", err)
				continue
			}
			log.Println(output)
		default:
			log.Printf("Unsupported event type: %v\n", eventType)
		}
	}
}

func runAutorandr(autorandr string) (string, error) {
	cmd := []string{autorandr, "--change", "--default", "default"}

	output, err := exec.Command(cmd[0], cmd[1:]...).Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}
