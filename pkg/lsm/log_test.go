package lsm

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAppendLog_WriteString(t *testing.T) {
	assert := assert.New(t)

	filepath := "./test.log"

	appendLog, err := NewAppendLog(filepath)
	assert.Nil(err)

	for i := 0; i < 5; i++ {
		err = appendLog.WriteString("hello\n")
		assert.Nil(err)
	}

	data, err := ioutil.ReadFile(filepath)
	assert.Nil(err)
	parts := strings.Split(string(data), "\n")

	assert.Equal(len(parts), 6)

	err = appendLog.Clear()
	assert.Nil(err)

	for i := 0; i < 2; i++ {
		err = appendLog.WriteString("pedro\n")
		assert.Nil(err)
	}

	data, err = ioutil.ReadFile(filepath)
	assert.Nil(err)
	parts = strings.Split(string(data), "\n")

	assert.Equal(len(parts), 6)

	err = os.Remove(filepath)
	assert.Nil(err)
}
