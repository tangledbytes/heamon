package utils

import (
	"os"

	"github.com/bwmarrin/snowflake"
	"github.com/sirupsen/logrus"
)

var node *snowflake.Node

func init() {
	var err error

	node, err = snowflake.NewNode(1)
	if err != nil {
		logrus.Error("failed to create id generator: ", err)
		os.Exit(1)
	}
}

// GenerateID returns a 64 bit id
func GenerateID() int64 {
	return node.Generate().Int64()
}
