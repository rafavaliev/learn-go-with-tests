package mock

import (
	"fmt"
	"io"
	"os"
	"time"
)

const finalWord = "Go!"

type Sleeper interface {
	Sleep()
}

type DefaultSleeper struct{}

type ConfigurableSleeper struct {
	duration time.Duration
	Sleep    func(time.Duration)
}

func (d *DefaultSleeper) Sleep() {
	time.Sleep(1 * time.Second)
}


func Countdown(writer io.Writer, sleeper Sleeper, counter int) {
	for counter > 0 {
		sleeper.Sleep()
		fmt.Fprintln(writer, counter)
		counter--
	}
	sleeper.Sleep()
	fmt.Fprint(writer, finalWord)
}

func main() {
	Countdown(os.Stdout, nil, 0)
}
