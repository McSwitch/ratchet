package stages

import (
	"bufio"
	"io"

	"github.com/DailyBurn/ratchet/data"
	"github.com/DailyBurn/ratchet/util"
)

// IoReader is a simple pipeline stage that wraps an io.Reader.
// It will typically be used as a starting stage in a pipeline.
type IoReader struct {
	Reader     io.Reader
	LineByLine bool // defaults to true
	BufferSize int
}

// NewIoReader returns a new IoReader wrapping the given io.Reader object
func NewIoReader(reader io.Reader) *IoReader {
	return &IoReader{Reader: reader, LineByLine: true, BufferSize: 1024}
}

func (r *IoReader) ProcessData(d data.JSON, outputChan chan data.JSON, killChan chan error) {
	r.ForEachData(killChan, func(d data.JSON) {
		outputChan <- d
	})
}

func (r *IoReader) Finish(outputChan chan data.JSON, killChan chan error) {
	close(outputChan)
}

func (r *IoReader) ForEachData(killChan chan error, foo func(d data.JSON)) {
	if r.LineByLine {
		r.scanLines(killChan, foo)
	} else {
		r.bufferedRead(killChan, foo)
	}
}

func (r *IoReader) scanLines(killChan chan error, forEach func(d data.JSON)) {
	scanner := bufio.NewScanner(r.Reader)
	for scanner.Scan() {
		forEach(data.JSON(scanner.Text()))
	}
	err := scanner.Err()
	util.KillPipelineIfErr(err, killChan)
}

func (r *IoReader) bufferedRead(killChan chan error, forEach func(d data.JSON)) {
	reader := bufio.NewReader(r.Reader)
	d := make([]byte, r.BufferSize)
	for {
		n, err := reader.Read(d)
		if err != nil && err != io.EOF {
			killChan <- err
		}
		if n == 0 {
			break
		}
		forEach(d)
	}
}

func (r *IoReader) String() string {
	return "IoReader"
}
