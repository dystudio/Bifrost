package server

import (
	"encoding/json"
	"fmt"
	"github.com/brokercap/Bifrost/config"
	pluginDriver "github.com/brokercap/Bifrost/plugin/driver"
	"github.com/brokercap/Bifrost/server/filequeue"
)

func GetFileQueue(dbName,SchemaName,tableName,ToServerID string) string  {
	return config.DataDir + "/filequeue/" + dbName + "/"+ SchemaName + "/" + tableName + "/" + ToServerID
}

//初始化文件队列
func (This *ToServer)  InitFileQueue(dbName,SchemaName,tableName string) *ToServer {
	if This.fileQueueObj == nil {
		This.fileQueueObj = filequeue.NewQueue( GetFileQueue(dbName,SchemaName,tableName,fmt.Sprint(This.ToServerID)) )
	}
	return This
}

//将数据刷到磁盘队列中
func (This *ToServer)  AppendToFileQueue(data *pluginDriver.PluginDataType) (error) {
	v,err := json.Marshal(data)
	if err != nil{
		return err
	}
	return This.fileQueueObj.AppendBytes(v)
}

//从磁盘队列中取出最前面一条数据
func (This *ToServer) PopFileQueue() (*pluginDriver.PluginDataType,error) {
	v,err := This.fileQueueObj.Pop()
	if err == nil && v == nil{
		return  nil,nil
	}
	if err != nil{
		return nil,err
	}
	var data pluginDriver.PluginDataType
	err = json.Unmarshal(v,&data)
	if err != nil {
		return nil,err
	}
	return &data,nil
}

//从磁盘队列中取出最后面一条数据
func (This *ToServer) ReadLastFromFileQueue() (*pluginDriver.PluginDataType,error) {
	v,err := This.fileQueueObj.ReadLast()
	if err == nil && v == nil{
		return  nil,nil
	}
	if err != nil{
		return nil,err
	}
	var data pluginDriver.PluginDataType
	err = json.Unmarshal(v,&data)
	if err != nil {
		return nil,err
	}
	return &data,nil
}