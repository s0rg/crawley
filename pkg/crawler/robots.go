package crawler

import "errors"

// ErrActionUnknown is returned when requested name unknown.
var ErrActionUnknown = errors.New("unknown action")

// RobotsAction is a action for robots.txt.
type RobotsAction byte

const (
	// RobotsIgnore ignores robots.txt completly.
	RobotsIgnore RobotsAction = 0
	// RobotsCrawl crawls urls from robots.txt, ignoring its rules.
	RobotsCrawl RobotsAction = 1
	// RobotsRespect same as above, but respects given rules.
	RobotsRespect RobotsAction = 2
)

// ParseAction parses robots action from string.
func ParseAction(s string) (a RobotsAction, err error) {
	switch s {
	case "ignore":
		a = RobotsIgnore
	case "crawl":
		a = RobotsCrawl
	case "respect":
		a = RobotsRespect
	default:
		err = ErrActionUnknown

		return
	}

	return a, nil
}
