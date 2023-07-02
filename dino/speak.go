package dino

import (
	"fmt"
	"strings"
)

type Actor struct {
	left  string
	right string
}
type Server struct{}

func NewDino() Actor {
	return Actor{
		left: `
%s
        \   __
         \ (_ \
             \ \_.----._
              \         \
               |  ) |  ) \___
               |_|--|_|'-.___\
`,
		right: `
%s
               __   /
              / _) /
     _.----._/ /
    /         /
 __/ (  | (  |
/__.-'|_|--|_|
`,
	}
}

func NewServer(name string) Actor {
	return Actor{
		left: `
%s
` + fmt.Sprintf(`        \   .-----------------.
         \  | =============== |
            | %-15s |
            |-----------------|
            | =============== |
            | ::::::::::::::: |
            ._________________.
`, name),
	}
}

func (a Actor) SayLeft(s string) {
	bubble := createTextBubble(s)

	fmt.Printf(a.left, bubble)
}

func (a Actor) SayRight(s string) {
	bubble := createTextBubble(s)

	fmt.Printf(a.right, bubble)
}

func createTextBubble(s string) string {
	w := NewWrapper(38)
	wrapped := w.splitString(s)
	header := " ________________________________________"
	banner := []string{header}
	for index, line := range wrapped {
		if index == 0 {
			banner = append(banner, fmt.Sprintf("/ %-38s \\", line))
		} else if index == len(wrapped)-1 {
			banner = append(banner, fmt.Sprintf("\\ %-38s /", line))
		} else {
			banner = append(banner, fmt.Sprintf("| %-38s |", line))
		}
	}
	banner = append(banner, header)
	bubble := strings.Join(banner, "\n")
	return bubble
}
