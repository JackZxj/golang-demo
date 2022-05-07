package main

import (
	"log"

	semver "github.com/Masterminds/semver/v3"
)

func main() {
	raw := []string{"1.2.3", "1.0", "1.3", "1.14.2", "v1.20.1-20", "v1.20.2", "v1.20.1-alpha1"}
	c, err := semver.NewConstraint(">=1.16.0-0")
	if err != nil {
		panic(err)
	}
	for _, r := range raw {
		v, err := semver.NewVersion(r)
		if err != nil {
			log.Panicf("Error parsing version: %s", err)
		}
		log.Println(v, c.Check(v))
	}
}
