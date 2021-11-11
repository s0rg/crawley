package crawler

import (
	"errors"
	"strings"
)

// ErrUnknownPolicy is returned when requested name unknown.
var ErrUnknownPolicy = errors.New("unknown policy")

// RobotsPolicy is a action for robots.txt.
type RobotsPolicy byte

const (
	// RobotsIgnore ignores robots.txt completly.
	RobotsIgnore RobotsPolicy = 0
	// RobotsCrawl crawls urls from robots.txt, ignoring its rules.
	RobotsCrawl RobotsPolicy = 1
	// RobotsRespect same as above, but respects given rules.
	RobotsRespect RobotsPolicy = 2
)

// ParseRobotsPolicy parses robots policy from string.
func ParseRobotsPolicy(s string) (a RobotsPolicy, err error) {
	switch strings.ToLower(s) {
	case "ignore":
		a = RobotsIgnore
	case "crawl":
		a = RobotsCrawl
	case "respect":
		a = RobotsRespect
	default:
		err = ErrUnknownPolicy

		return
	}

	return a, nil
}
