package crawler

import (
	"errors"
	"strings"
)

const (
	// DefaultRobotsPolicy is a default policy name for robots handling.
	DefaultRobotsPolicy = "ignore"
	// DefaultDirsPolicy is a default policy name for non-resource URLs.
	DefaultDirsPolicy = "show"
)

// ErrUnknownPolicy is returned when requested policy unknown.
var ErrUnknownPolicy = errors.New("unknown policy")

// RobotsPolicy is a policy for robots.txt.
type RobotsPolicy byte

const (
	// RobotsIgnore ignores robots.txt completly.
	RobotsIgnore RobotsPolicy = 0
	// RobotsCrawl crawls urls from robots.txt, ignoring its rules.
	RobotsCrawl RobotsPolicy = 1
	// RobotsRespect same as above, but respects given rules.
	RobotsRespect RobotsPolicy = 2
)

// DirsPolicy is a policy for non-resorce urls.
type DirsPolicy byte

const (
	// DirsShow show directories.
	DirsShow DirsPolicy = 0
	// DirsHide hide directories from output.
	DirsHide DirsPolicy = 1
	// DirsOnly show only directories in output.
	DirsOnly DirsPolicy = 2
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

// ParseDirsPolicy parses dirs policy from string.
func ParseDirsPolicy(s string) (p DirsPolicy, err error) {
	switch strings.ToLower(s) {
	case "show":
		p = DirsShow
	case "hide":
		p = DirsHide
	case "only":
		p = DirsOnly
	default:
		err = ErrUnknownPolicy

		return
	}

	return p, nil
}
