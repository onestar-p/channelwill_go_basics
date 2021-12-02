package id_generate

import (
	"strconv"

	uuid "github.com/satori/go.uuid"
	"github.com/sony/sonyflake"

	"channelwill_go_basics/global"
)

type IDGenerateInterfaced interface {
	GetID() (string, error)
}

type IDGenerateType int

// 生成ID类型
const (
	Snowflake IDGenerateType = iota // 雪花ID
	GenUuid                         // UUID
)

// @params IDGenerateType 生成ID类型
func NewIDGenerate(IDGenerateType IDGenerateType) IDGenerateInterfaced {
	if IDGenerateType == Snowflake {
		return NewSnowFlake(func() (uint16, error) {
			return global.ApplicationConfig.MachineId, nil
		})
	}

	if IDGenerateType == GenUuid {
		return &GenUUID{}
	}
	return nil
}

// Generate ID by UUID
type GenUUID struct {
}

func (u *GenUUID) GetID() (string, error) {
	id := uuid.NewV4()
	return id.String(), nil

}

// Generate ID by SnowFlake
type SnowFlake struct {
	snowFlake *sonyflake.Sonyflake
}

// @params machineIDFun 获取本机的机器ID匿名函数
func NewSnowFlake(machineIDFun func() (uint16, error)) *SnowFlake {
	st := sonyflake.Settings{}
	// machineID是个回调函数
	st.MachineID = machineIDFun
	return &SnowFlake{
		snowFlake: sonyflake.NewSonyflake(st),
	}
}

func (s *SnowFlake) GetID() (string, error) {
	id, err := s.snowFlake.NextID()
	idStr := strconv.Itoa(int(id))
	return idStr, err
}
