package util

import (
	"github.com/bwmarrin/snowflake"
)

const (
	nodeNum                  = 1
	sequenceIdStep           = 12
	sequenceIdStepMask int64 = -1 ^ (-1 << sequenceIdStep)
)

var idGenerator *IdGenerator

type IdGenerator struct {
	node *snowflake.Node
}

func InitIdGenerator() {
	node, err := snowflake.NewNode(nodeNum)
	if err != nil {
		panic(err)
		return
	}
	idGenerator = &IdGenerator{node: node}
}

// +-----------------------------------------+
// |  41 Bit Timestamp |  12 Bit Sequence Id |
// +-----------------------------------------+

func GenerateId() int64 {
	// 雪花算法原ID
	origId := idGenerator.node.Generate()
	// 自制ID
	id := int64(0)
	id |= origId.Step()                   // 12 Bit Sequence Id
	id |= origId.Time() << sequenceIdStep // 41 Bit Timestamp
	return id
}
