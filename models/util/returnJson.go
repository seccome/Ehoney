package util

func ReturnJsonsError(str string) map[string]interface{} {
	list := []map[string]interface{}{}
	data := map[string]interface{}{"status":"ERROR","msg":str,"data":list}
	return data
}

func ReturnJsonsSuccess(data map[string]interface{}) map[string]interface{} {
	datas := map[string]interface{}{"code":0,"message":"成功","data":data}
	return datas
}

func FormatJson(total int, listMap []map[string]interface{},pageSize int,pageNum int,totalPage int) map[string]interface{}{
	data := map[string]interface{}{"total":total,"totalPage":totalPage,"pageNum":pageNum,"pageSize":pageSize, "list":listMap}
	// 最终的集合
	datas := ReturnJsonsSuccess(data)
	return datas
}

func BaitPolicyListMap(agentId string,bid int,optTime string,optUser string,baittype string,dir string,enable bool) map[string]interface{}{
	listData := map[string]interface{}{"agentId":agentId,"bid":bid,"optTime":optTime,"optUser":optUser,"baitType":baittype,"dir":dir,"enable":enable}
	return listData
}

func TransPolicyListMap(agentId string,listenPort int,honeyPort int,optTime string,optUser string,serverType string,honeyIP string,enable bool) map[string]interface{}{
	listData := map[string]interface{}{"agentId":agentId,"listenPort":listenPort,"honeyPort":honeyPort,"optTime":optTime,"optUser":optUser,"serverType":serverType,"honeyIP":honeyIP,"enable":enable}
	return listData
}

func ReturnTransJson(listenPort int, serverType string, honeyIP string, honeyPort int, enable bool) map[string]interface{}{
	data := map[string]interface{}{"listenPort":listenPort,"serverType":serverType,"honeyIP":honeyIP,"honeyPort":honeyPort,"enable":enable}
	return data
}

func ReturnBaitJson(address string, executeFile string, md5 string) map[string]interface{}{
	data := map[string]interface{}{"address":address,"executeFile":executeFile,"md5":md5}
	return data
}