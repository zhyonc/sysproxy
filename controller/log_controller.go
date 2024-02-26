package controller

import (
	"bytes"
	"log"
)

type logController struct {
	Buffer *bytes.Buffer
}

func NewLogController() LogController {
	c := &logController{Buffer: new(bytes.Buffer)}
	log.SetOutput(c.Buffer)
	return c
}

func (c *logController) GetLogInfo() string {
	msg := c.Buffer.String()
	c.Buffer.Reset()
	return msg
}
