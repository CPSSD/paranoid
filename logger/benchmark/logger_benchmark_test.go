package logger

import (
	"bytes"
	. "github.com/cpssd/paranoid/logger"
	syslog "log"
	"os"
	"strconv"
	"testing"
)

func BenchmarkLoging(b *testing.B) {
	log := New("testPackage", "testComponent", os.DevNull)
	err := log.SetOutput(STDERR)
	if err != nil {
		syslog.Fatal("Failed to set logger output:", err)
	}
	log.SetLogLevel(ERROR)
	for n := 0; n < b.N; n++ {
		str := strconv.Itoa(n)
		testString := "test" + str

		var b bytes.Buffer
		log.AddAdditionalWriter(&b)
		log.Info(testString)
	}
}
