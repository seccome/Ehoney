package util

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"log"
	"strconv"
	"strings"
	"time"
	_ "time"
)

type BaitPolicyAll struct {
	Pbid        string    `json:"pbid"`
	Aid         string    `json:"aid"`
	Bid         string    `json:"bid"`
	ExecuteFile string    `json:"executeFile"`
	OptTime     time.Time `json:"optTime"`
	OptUser     string    `json:"optUser"`
	Enable      string    `json:"enable"`
	Address     string    `json:"address"`
	Md5         string    `json:"md5"`
	Dir         string    `json:"dir"`
	BaitType    string    `json:"baittype"`
}
type TransPolicyAll struct {
	Ptid       string    `json:"ptid"`
	Aid        string    `json:"aid"`
	ListenPort string    `json:"listenPort"`
	HoneyPort  string    `json:"honeyPort"`
	ServerType string    `json:"serverType"`
	Enable     string    `json:"enable"`
	HoneyIP    string    `json:"honeyIp"`
	OptUser    string    `json:"optUser"`
	OptTime    time.Time `json:"optTime"`
}
type Sum struct {
	Count string `json:"total"`
}

func BaitPolicyToJson(maps []orm.Params) []map[string]interface{} {
	var baitpolicy []BaitPolicyAll
	re, _ := json.Marshal(maps)
	err := json.Unmarshal(re, &baitpolicy)
	if err != nil {
		logs.Error("[SqlDataToMap] oplist unmarshal error,%s", err)
	}
	fmt.Println(baitpolicy)
	listMap := []map[string]interface{}{}
	for i := 0; i < len(baitpolicy); i++ {
		bids := baitpolicy[i].Bid
		bid, _ := strconv.Atoi(bids)
		agentId := baitpolicy[i].Aid
		operateTime := baitpolicy[i].OptTime.Format("2006-01-02 15:04:05")
		operateUser := baitpolicy[i].OptUser
		dir := baitpolicy[i].Dir
		enables := baitpolicy[i].Enable
		baittype := baitpolicy[i].BaitType
		var enable bool
		if enables == "1" {
			enable = true
		} else {
			enable = false
		}
		listMap = append(listMap, BaitPolicyListMap(agentId, bid, operateTime, operateUser, baittype, dir, enable))
	}
	return listMap
}

func TransPolicyToJson(maps []orm.Params) []map[string]interface{} {
	var transpolicy []TransPolicyAll
	re, _ := json.Marshal(maps)
	err := json.Unmarshal(re, &transpolicy)
	if err != nil {
		logs.Error("[SqlDataToMap] oplist unmarshal error,%s", err)
	}
	fmt.Println(transpolicy)
	listMap := []map[string]interface{}{}
	for i := 0; i < len(transpolicy); i++ {
		agentId := transpolicy[i].Aid
		optTime := transpolicy[i].OptTime.Format("2006-01-02 15:04:05")
		optUser := transpolicy[i].OptUser
		lsPort := transpolicy[i].ListenPort
		listenPort, _ := strconv.Atoi(lsPort)
		hPort := transpolicy[i].HoneyPort
		honeyPort, _ := strconv.Atoi(hPort)
		serverType := transpolicy[i].ServerType
		honeyIP := transpolicy[i].HoneyIP
		enables := transpolicy[i].Enable
		var enable bool
		if enables == "1" {
			enable = true
		} else {
			enable = false
		}
		listMap = append(listMap, TransPolicyListMap(agentId, listenPort, honeyPort, optTime, optUser, serverType, honeyIP, enable))
	}
	return listMap
}

func SqlDataSumToInt(maps []orm.Params) int {
	var sum []Sum
	re, _ := json.Marshal(maps)
	err := json.Unmarshal(re, &sum)
	if err != nil {
		logs.Error("[SqlDataSumToInt] oplist unmarshal error,%s", err)
		return -1
	}
	total, _ := strconv.Atoi(sum[0].Count)
	return total
}

type BaitPolicyByStatus struct {
	TaskId      string
	AgentId     string
	BaitId      string
	CreateTime  string
	OfflineTime string
	Status      int
	Data        string
	Md5         string
	Type        string
}

type TransPolicyByStatus struct {
	TaskId       string
	AgentId      string
	HoneyIP      string
	HoneyPotType string
	ForwardPort  int
	HoneyPotPort int
	CreateTime   string
	OfflineTime  string
	Status       int
	Type         string
	Path         string
}

//type BaitPolicyAll struct {
//	Pbid        string     	`json:"pbid"`
//	Aid			string    	`json:"aid"`
//	Bid			string	 	`json:"bid"`
//	ExecuteFile string		`json:"executeFile"`
//	OptTime		time.Time   `json:"optTime"`
//	OptUser		string		`json:"optUser"`
//	Enable		string		`json:"enable"`
//	Address		string		`json:"address"`
//	Md5			string		`json:"md5"`
//	Dir			string		`json:"dir"`
//	BaitType	string		`json:"baittype"`
//}
//type TransPolicyAll struct {
//	Ptid        string     	`json:"ptid"`
//	Aid			string    	`json:"aid"`
//	ListenPort	string		`json:"listenPort"`
//	HoneyPort	string		`json:"honeyPort"`
//	ServerType	string		`json:"serverType"`
//	Enable		string		`json:"enable"`
//	HoneyIP		string		`json:"honeyIp"`
//	OptUser		string		`json:"optUser"`
//	OptTime		time.Time	`json:"optTime"`
//}
//type Sum struct {
//	Count	string	`json:"total"`
//}

//func BaitPolicyToJson(maps []orm.Params) []map[string]interface{} {
//	var baitpolicy []BaitPolicyAll
//	re,_ := json.Marshal(maps)
//	err := json.Unmarshal(re,&baitpolicy)
//	if err != nil{
//		logs.Error("[SqlDataToMap] oplist unmarshal error,%s",err)
//	}
//	fmt.Println(baitpolicy)
//	listMap := []map[string]interface{}{}
//	for i:=0 ; i<len(baitpolicy);i++{
//		bids := baitpolicy[i].Bid
//		bid,_:= strconv.Atoi(bids)
//		agentId := baitpolicy[i].Aid
//		operateTime := baitpolicy[i].OptTime.Format("2006-01-02 15:04:05")
//		operateUser := baitpolicy[i].OptUser
//		dir := baitpolicy[i].Dir
//		enables := baitpolicy[i].Enable
//		baittype := baitpolicy[i].BaitType
//		var enable bool
//		if enables == "1"{
//			enable = true
//		}else{
//			enable = false
//		}
//		listMap = append(listMap,BaitPolicyListMap(agentId,bid,operateTime,operateUser,baittype,dir,enable))
//	}
//	return listMap
//}
//
//func TransPolicyToJson(maps []orm.Params) []map[string]interface{} {
//	var transpolicy []TransPolicyAll
//	re,_ := json.Marshal(maps)
//	err := json.Unmarshal(re,&transpolicy)
//	if err != nil{
//		logs.Error("[SqlDataToMap] oplist unmarshal error,%s",err)
//	}
//	fmt.Println(transpolicy)
//	listMap := []map[string]interface{}{}
//	for i:=0 ; i<len(transpolicy);i++{
//		agentId := transpolicy[i].Aid
//		optTime := transpolicy[i].OptTime.Format("2006-01-02 15:04:05")
//		optUser := transpolicy[i].OptUser
//		lsPort := transpolicy[i].ListenPort
//		listenPort,_:= strconv.Atoi(lsPort)
//		hPort := transpolicy[i].HoneyPort
//		honeyPort,_ := strconv.Atoi(hPort)
//		serverType := transpolicy[i].ServerType
//		honeyIP := transpolicy[i].HoneyIP
//		enables := transpolicy[i].Enable
//		var enable bool
//		if enables == "1"{
//			enable = true
//		}else{
//			enable = false
//		}
//		listMap = append(listMap,TransPolicyListMap(agentId,listenPort,honeyPort,optTime,optUser,serverType,honeyIP,enable))
//	}
//	return listMap
//}
//
//func SqlDataSumToInt(maps []orm.Params) int {
//	var sum []Sum
//	re,_ :=json.Marshal(maps)
//	err := json.Unmarshal(re,&sum)
//	if err != nil{
//		logs.Error("[SqlDataSumToInt] oplist unmarshal error,%s",err)
//		return -1
//	}
//	total,_ := strconv.Atoi(sum[0].Count)
//	return total
//}

func GetTopAttackTypesMysqlJson(rows *sql.Rows, columns []string, all int, values []sql.RawBytes, scanArgs []interface{}) (string, int, error) {
	list := "{\"list\":["
	count := 0

	for rows.Next() {
		count += 1
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Println("scan failed, err: ", err)
			return "", count, err
		}

		row := "{"
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			//site := ""
			columName := strings.ToLower(columns[i])
			cell := ""
			if columName == "honeypottype" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"attacktype":""`)
				} else {
					cell = fmt.Sprintf(`"attacktype":"%s"`, value)
				}
			} else if columName == "typecount" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"typecount":""`)
				} else {
					cell = fmt.Sprintf(`"typecount":"%s"`, value)
				}
			}
			row = row + cell + ","
		}
		row = row[0 : len(row)-1]
		row += "}"
		list = list + row + ","
	}
	if count != 0 {
		list = list[0 : len(list)-1]
	}
	list += "],"
	cell := fmt.Sprintf(`"total":%d`, all)
	list = list + cell
	list += "}"

	list = strings.Replace(list, "\\", "", -1)
	log.Println(list)
	return list, count, nil
}

func GetTopSourceIpsMysqlJson(rows *sql.Rows, columns []string, all int, values []sql.RawBytes, scanArgs []interface{}) (string, int, error) {
	list := "{\"list\":["
	count := 0

	for rows.Next() {
		count += 1
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Println("scan failed, err: ", err)
			return "", count, err
		}

		row := "{"
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			//site := ""
			columName := strings.ToLower(columns[i])
			cell := ""
			if columName == "srchost" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"srchost":""`)
				} else {
					cell = fmt.Sprintf(`"srchost":"%s"`, value)
				}
			} else if columName == "ipcount" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"ipcount":""`)
				} else {
					cell = fmt.Sprintf(`"ipcount":"%s"`, value)
				}
			}
			row = row + cell + ","
		}
		row = row[0 : len(row)-1]
		row += "}"
		list = list + row + ","
	}
	if count != 0 {
		list = list[0 : len(list)-1]
	}
	list += "],"
	cell := fmt.Sprintf(`"total":%d`, all)
	list = list + cell
	list += "}"

	list = strings.Replace(list, "\\", "", -1)
	log.Println(list)
	return list, count, nil
}

func GetTopAttackMapMysqlJson(rows *sql.Rows, columns []string, all int, values []sql.RawBytes, scanArgs []interface{}, ownip string, longitude string, latitude string) (string, int, error) {
	list := "{\"list\":["
	count := 0

	for rows.Next() {
		count += 1
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Println("scan failed, err: ", err)
			return "", count, err
		}

		row := "{"
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			//site := ""
			columName := strings.ToLower(columns[i])
			cell := ""
			if columName == "attackip" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"attackip":""`)
				} else {
					cell = fmt.Sprintf(`"attackip":"%s"`, value)
				}
			} else if columName == "province" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"province":""`)
				} else {
					cell = fmt.Sprintf(`"province":"%s"`, value)
				}
			} else if columName == "country" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"country":""`)
				} else {
					cell = fmt.Sprintf(`"country":"%s"`, value)
				}
			} else if columName == "longitude" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"longitude":""`)
				} else {
					cell = fmt.Sprintf(`"longitude":"%s"`, value)
				}
			} else if columName == "latitude" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"latitude":""`)
				} else {
					cell = fmt.Sprintf(`"latitude":"%s"`, value)
				}
			} else if columName == "attacksum" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"attacksum":""`)
				} else {
					cell = fmt.Sprintf(`"attacksum":"%s"`, value)
				}
			}
			row = row + cell + ","

		}
		row = row[0 : len(row)-1]
		row += "}"
		list = list + row + ","
	}
	if count != 0 {
		list = list[0 : len(list)-1]
	}
	list += "],"
	cell := fmt.Sprintf(`"total":%d`, all)
	list = list + cell + ","
	cell = fmt.Sprintf(`"desip":"%s"`, ownip)
	list = list + cell + ","
	cell = fmt.Sprintf(`"longitude":"%s"`, longitude)
	list = list + cell + ","
	cell = fmt.Sprintf(`"latitude":"%s"`, latitude)
	list = list + cell
	list += "}"

	list = strings.Replace(list, "\\", "", -1)
	return list, count, nil
}

func GetTopAttackIpsMysqlJson(rows *sql.Rows, columns []string, all int, values []sql.RawBytes, scanArgs []interface{}) (string, int, error) {
	list := "{\"list\":["
	count := 0

	for rows.Next() {
		count += 1
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Println("scan failed, err: ", err)
			return "", count, err
		}

		row := "{"
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			//site := ""
			columName := strings.ToLower(columns[i])
			cell := ""
			if columName == "attackip" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"attackip":""`)
				} else {
					cell = fmt.Sprintf(`"attackip":"%s"`, value)
				}
			} else if columName == "ipcount" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"ipcount":""`)
				} else {
					cell = fmt.Sprintf(`"ipcount":"%s"`, value)
					row = row + cell + ","
					ipcountnum, _ := strconv.Atoi(value)
					percentnum := GetPercent(ipcountnum, all)
					cell = fmt.Sprintf(`"ippercent":"%s"`, percentnum)
				}
			}
			row = row + cell + ","

		}
		row = row[0 : len(row)-1]
		row += "}"
		list = list + row + ","
	}
	if count != 0 {
		list = list[0 : len(list)-1]
	}
	list += "],"
	cell := fmt.Sprintf(`"total":%d`, all)
	list = list + cell
	list += "}"

	list = strings.Replace(list, "\\", "", -1)
	log.Println(list)
	return list, count, nil
}

func GetTopAreasMysqlJson(rows *sql.Rows, columns []string, all int, values []sql.RawBytes, scanArgs []interface{}) (string, int, error) {
	list := "{\"list\":["
	count := 0
	areacount := all
	for rows.Next() {
		count += 1
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Println("scan failed, err: ", err)
			return "", count, err
		}

		row := "{"
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			//site := ""
			columName := strings.ToLower(columns[i])
			cell := ""
			if columName == "country" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"country":""`)
				} else {
					cell = fmt.Sprintf(`"country":"%s"`, value)
				}
			} else if columName == "province" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"province":""`)
				} else {
					cell = fmt.Sprintf(`"province":"%s"`, value)
				}
			} else if columName == "areacount" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"areacount":""`)
				} else {
					ipcountnum, _ := strconv.Atoi(value)
					cell = fmt.Sprintf(`"%v":%v`, "areacount", ipcountnum)

					all = all - ipcountnum
				}
			}
			row = row + cell + ","

		}

		//cell := fmt.Sprintf(`"areacount":"%s"`, all)
		//row = row + cell + ","
		//cell = fmt.Sprintf(`"country":"%s"`, "其他")
		//row = row + cell + ","
		//cell = fmt.Sprintf(`"province":"%s"`, "其他")
		//row = row + cell + ","
		row = row[0 : len(row)-1]
		row += "}"
		list = list + row + ","
	}
	if count != 0 {
		list = list[0 : len(list)-1]
	}
	str2 := fmt.Sprintf("%d", all)
	list += ",{\"areacount\":" + str2 + ",\"country\":\"其他\",\"province\":\"其他\"}"
	list += "],"
	cell := fmt.Sprintf(`"total":%d`, areacount)
	list = list + cell
	list += "}"

	list = strings.Replace(list, "\\", "", -1)
	log.Println(list)
	return list, count, nil
}

func GetHoneyPotTransMysqlJson(rows *sql.Rows, columns []string, all int, values []sql.RawBytes, scanArgs []interface{}, pageNum int, pageSize int, totalPage int) (string, int, error) {
	list := "{\"list\":["
	count := 0

	for rows.Next() {
		count += 1
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Println("scan failed, err: ", err)
			return "", count, err
		}

		row := "{"
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}

			columName := strings.ToLower(columns[i])
			cell := ""

			if columName == "forwardport" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"forwardport":""`)
				} else {
					cell = fmt.Sprintf(`"forwardport":"%v"`, value)
				}
			} else if columName == "honeypottype" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"honeypottype":""`)
				} else {
					cell = fmt.Sprintf(`"honeypottype":"%v"`, value)
				}
			} else if columName == "taskid" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"taskid":""`)
				} else {
					cell = fmt.Sprintf(`"taskid":"%v"`, value)
				}
			} else if columName == "honeyip" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"honeyip":""`)
				} else {
					cell = fmt.Sprintf(`"honeyip":"%v"`, value)
				}
			} else if columName == "honeyport" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"honeyport":""`)
				} else {
					cell = fmt.Sprintf(`"honeyport":"%v"`, value)
				}

			} else if columName == "createtime" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"createTime":""`)
				} else {
					cell = fmt.Sprintf(`"createTime":"%v"`, value)
				}
			} else if columName == "offlinetime" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"offlinetime":""`)
				} else {
					cell = fmt.Sprintf(`"offlinetime":"%v"`, value)
				}
			} else if columName == "creator" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"creator":""`)
				} else {
					cell = fmt.Sprintf(`"%v":"%v"`, columName, value)
				}
			} else if columName == "serverip" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"serverip":""`)
				} else {
					cell = fmt.Sprintf(`"serverip":"%v"`, value)
				}
			} else if columName == "status" {
				status, _ := strconv.Atoi(value)
				cell = fmt.Sprintf(`"%v":%v`, columName, status)
			}

			row = row + cell + ","
		}
		row = row[0 : len(row)-1]
		row += "}"
		list = list + row + ","

	}
	if count != 0 {
		list = list[0 : len(list)-1]
	}
	list += "],"
	cell := fmt.Sprintf(`"total":%d`, all)
	totalpage := fmt.Sprintf(`"totalPage":%d`, totalPage)
	pagenum := fmt.Sprintf(`"pageNum":%d`, pageNum)
	pagesize := fmt.Sprintf(`"pageSize":%d`, pageSize)
	list = list + cell + "," + totalpage + "," + pagenum + "," + pagesize
	list += "}"

	list = strings.Replace(list, "\\", "", -1)
	//log.Println(list)
	return list, count, nil
}

func GetSignPolicyMysqlJson(rows *sql.Rows, columns []string, all int, values []sql.RawBytes, scanArgs []interface{}, pageNum int, pageSize int, totalPage int) (string, int, error) {
	list := "{\"list\":["
	count := 0

	for rows.Next() {
		count += 1
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Println("scan failed, err: ", err)
			return "", count, err
		}

		row := "{"
		var value string
		//signinfo := ""
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}

			columName := strings.ToLower(columns[i])
			cell := ""

			if columName == "data" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"data":""`)
				} else {
					cell = fmt.Sprintf(`"data":"%v"`, value)
					//signinfo = value
				}
			} else if columName == "taskid" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"taskId":""`)
				} else {
					cell = fmt.Sprintf(`"taskId":"%v"`, value)
				}
			} else if columName == "signname" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"signname":""`)
				} else {
					cell = fmt.Sprintf(`"signname":"%v"`, value)
				}
			} else if columName == "signinfo" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"signinfo":""`)
				} else {
					//signinfo = signinfo + "/" + value
					cell = fmt.Sprintf(`"signinfo":"%v"`, value)
				}
			} else if columName == "type" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"type":""`)
				} else {
					cell = fmt.Sprintf(`"type":"%v"`, value)
				}

			} else if columName == "createtime" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"createTime":""`)
				} else {
					cell = fmt.Sprintf(`"createTime":"%v"`, value)
				}
			} else if columName == "offlinetime" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"offlineTime":""`)
				} else {
					cell = fmt.Sprintf(`"offlineTime":"%v"`, value)
				}
			} else if columName == "creator" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"creator":""`)
				} else {
					cell = fmt.Sprintf(`"%v":"%v"`, columName, value)
				}
			} else if columName == "status" {
				status, _ := strconv.Atoi(value)
				cell = fmt.Sprintf(`"%v":%v`, columName, status)
			} else if columName == "tracecode" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"tracecode":""`)
				} else {
					cell = fmt.Sprintf(`"%v":"%v"`, columName, value)
				}
			}

			row = row + cell + ","
		}
		row = row[0 : len(row)-1]
		row += "}"
		list = list + row + ","

	}
	if count != 0 {
		list = list[0 : len(list)-1]
	}
	list += "],"
	cell := fmt.Sprintf(`"total":%d`, all)
	totalpage := fmt.Sprintf(`"totalPage":%d`, totalPage)
	pagenum := fmt.Sprintf(`"pageNum":%d`, pageNum)
	pagesize := fmt.Sprintf(`"pageSize":%d`, pageSize)
	list = list + cell + "," + totalpage + "," + pagenum + "," + pagesize
	list += "}"

	list = strings.Replace(list, "\\", "", -1)
	log.Println(list)
	return list, count, nil
}

func GetBaitPolicyMysqlJson(rows *sql.Rows, columns []string, all int, values []sql.RawBytes, scanArgs []interface{}, pageNum int, pageSize int, totalPage int) (string, int, error) {
	list := "{\"list\":["
	count := 0

	for rows.Next() {
		count += 1
		err := rows.Scan(scanArgs...)
		if err != nil {
			logs.Error("[GetBaitPolicyMysqlJson] list error:", err)
			return "", count, err
		}

		row := "{"
		var value string
		baitdata := ""
		baittype := ""
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}

			columName := strings.ToLower(columns[i])
			cell := ""

			if columName == "type" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"baitType":""`)
				} else {
					cell = fmt.Sprintf(`"baitType":"%v"`, value)
					baittype = value
				}
			} else if columName == "taskid" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"taskId":""`)
				} else {
					cell = fmt.Sprintf(`"taskId":"%v"`, value)
				}
			} else if columName == "baitname" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"baitname":""`)
				} else {
					cell = fmt.Sprintf(`"baitname":"%v"`, value)
				}
			} else if columName == "data" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"data":""`)
				} else {
					cell = fmt.Sprintf(`"data":"%v"`, value)
					if baittype == "history" {
						baitdata = Base64Decode(value)
						baitdata = strings.Replace(baitdata, "\r", "\\r", -1)
						baitdata = strings.Replace(baitdata, "\n", "\\n", -1)
						baitdata = strings.Replace(baitdata, "\r\n", "\\r\\n", -1)
						//baitdata = strings.Replace(baitdata, "\\r\\n", "", -1)
						baitdata = strings.Replace(baitdata, "\"", "\\\"", -1)
					} else {
						baitdata = value
					}
				}
			} else if columName == "baitinfo" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"baitInfo":"%v"`, baitdata)
				} else {
					if baittype == "history" {
						cell = fmt.Sprintf(`"baitInfo":"%v"`, Base64Decode(value))
					} else {
						cell = fmt.Sprintf(`"baitInfo":"%v"`, baitdata+"/"+value)
					}
				}
				baitdata = ""
				baittype = ""
			} else if columName == "createtime" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"createTime":""`)
				} else {
					cell = fmt.Sprintf(`"createTime":"%v"`, value)
				}
			} else if columName == "offlinetime" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"offlineTime":""`)
				} else {
					cell = fmt.Sprintf(`"offlineTime":"%v"`, value)
				}
			} else if columName == "creator" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"creator":""`)
				} else {
					cell = fmt.Sprintf(`"%v":"%v"`, columName, value)
				}
			} else if columName == "status" {
				status, _ := strconv.Atoi(value)
				cell = fmt.Sprintf(`"%v":%v`, columName, status)
			}

			row = row + cell + ","
		}
		row = row[0 : len(row)-1]
		row += "}"
		list = list + row + ","

	}
	if count != 0 {
		list = list[0 : len(list)-1]
	}
	list += "],"
	cell := fmt.Sprintf(`"total":%d`, all)
	totalpage := fmt.Sprintf(`"totalPage":%d`, totalPage)
	pagenum := fmt.Sprintf(`"pageNum":%d`, pageNum)
	pagesize := fmt.Sprintf(`"pageSize":%d`, pageSize)
	list = list + cell + "," + totalpage + "," + pagenum + "," + pagesize
	list += "}"

	list = strings.Replace(list, "\\", "", -1)
	log.Println(list)
	return list, count, nil
}

func GetTransPolicyMysqlJson(rows *sql.Rows, columns []string, all int, values []sql.RawBytes, scanArgs []interface{}, pageNum int, pageSize int, totalPage int) (string, int, error) {
	list := "{\"list\":["
	count := 0

	for rows.Next() {
		count += 1
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Println("scan failed, err: ", err)
			return "", count, err
		}

		row := "{"
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}

			columName := strings.ToLower(columns[i])
			cell := ""
			if columName == "forwardport" {
				forwardport, _ := strconv.Atoi(value)
				cell = fmt.Sprintf(`"forwardPort":%d`, forwardport)
			} else if columName == "honeypottype" {
				if value == "NULL" || value == "" {
					cell = fmt.Sprintf(`"honeyPotType":""`)
				} else {
					cell = fmt.Sprintf(`"honeyPotType":"%s"`, value)
				}
			} else if columName == "serverip" {
				if value == "NULL" || value == "" {
					cell = fmt.Sprintf(`"serverip":""`)
				} else {
					cell = fmt.Sprintf(`"serverip":"%s"`, value)
				}
			} else if columName == "honeyserverip" {
				if value == "NULL" || value == "" {
					cell = fmt.Sprintf(`"honeyserverip":""`)
				} else {
					cell = fmt.Sprintf(`"honeyserverip":"%s"`, value)
				}
			} else if columName == "honeypotport" {
				honeypotport, _ := strconv.Atoi(value)
				cell = fmt.Sprintf(`"honeyPotPort":%d`, honeypotport)
			} else if columName == "createtime" {
				if value == "NULL" || value == "" {
					cell = fmt.Sprintf(`"createTime":""`)
				} else {
					cell = fmt.Sprintf(`"createTime":"%s"`, value)
				}
			} else if columName == "offlinetime" {
				if value == "NULL" || value == "" {
					cell = fmt.Sprintf(`"offlineTime":""`)
				} else {
					cell = fmt.Sprintf(`"offlineTime":"%s"`, value)
				}
			} else if columName == "creator" {
				if value == "NULL" || value == "" {
					cell = fmt.Sprintf(`"creator":""`)
				} else {
					cell = fmt.Sprintf(`"%s":"%s"`, columName, value)
				}
			} else if columName == "taskid" {
				if value == "NULL" || value == "" {
					cell = fmt.Sprintf(`"taskId":""`)
				} else {
					cell = fmt.Sprintf(`"taskId":"%s"`, value)
				}
			} else if columName == "status" {
				status, _ := strconv.Atoi(value)
				cell = fmt.Sprintf(`"%s":%d`, columName, status)
			}

			row = row + cell + ","
		}
		row = row[0 : len(row)-1]
		row += "}"
		list = list + row + ","

	}
	if count != 0 {
		list = list[0 : len(list)-1]
	}
	list += "],"
	cell := fmt.Sprintf(`"total":%d`, all)
	totalpage := fmt.Sprintf(`"totalPage":%d`, totalPage)
	pagenum := fmt.Sprintf(`"pageNum":%d`, pageNum)
	pagesize := fmt.Sprintf(`"pageSize":%d`, pageSize)
	list = list + cell + "," + totalpage + "," + pagenum + "," + pagesize
	list += "}"

	list = strings.Replace(list, "\\", "", -1)
	fmt.Println(list)
	return list, count, nil
}

func GetHoneyPotsInfoMysqlJson(rows *sql.Rows, columns []string, all int, values []sql.RawBytes, scanArgs []interface{}) (string, int, error) {
	list := "{\"list\":["
	count := 0

	for rows.Next() {
		count += 1
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Println("scan failed, err: ", err)
			return "", count, err
		}

		row := "{"
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}

			columName := strings.ToLower(columns[i])
			cell := ""
			if columName == "forwardport" {
				forwardport, _ := strconv.Atoi(value)
				cell = fmt.Sprintf(`"honeyPotPort":%d`, forwardport)
			} else if columName == "servername" {
				if value == "NULL" || value == "" {
					cell = fmt.Sprintf(`"serverName":""`)
				} else {
					cell = fmt.Sprintf(`"serverName":"%s"`, value)
				}
			} else if columName == "serverip" {
				if value == "NULL" || value == "" {
					cell = fmt.Sprintf(`"serverIp":""`)
				} else {
					cell = fmt.Sprintf(`"serverIp":"%s"`, value)
				}
			} else if columName == "serverid" {
				if value == "NULL" || value == "" {
					cell = fmt.Sprintf(`"serverId":""`)
				} else {
					cell = fmt.Sprintf(`"serverId":"%s"`, value)
				}
			}

			row = row + cell + ","
		}
		row = row[0 : len(row)-1]
		row += "}"
		list = list + row + ","

	}
	if count != 0 {
		list = list[0 : len(list)-1]
	}
	list += "],"
	cell := fmt.Sprintf(`"total":%d`, all)
	list = list + cell
	list += "}"

	list = strings.Replace(list, "\\", "", -1)
	log.Println(list)
	return list, count, nil
}

func GetAttackLogListMysqlJson(rows *sql.Rows, columns []string, all int, values []sql.RawBytes, scanArgs []interface{}, pageNum int, pageSize int, totalPage int) (string, int, error) {
	list := "{\"list\":["
	count := 0

	for rows.Next() {
		count += 1
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Println("scan failed, err: ", err)
			return "", count, err
		}

		row := "{"
		location := ""
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			columName := strings.ToLower(columns[i])
			cell := ""
			if columName == "id" {
				id, _ := strconv.Atoi(value)
				cell = fmt.Sprintf(`"id":%d`, id)
			} else if columName == "srchost" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"srcHost":""`)
				} else {
					cell = fmt.Sprintf(`"srcHost":"%s"`, value)
				}
			} else if columName == "honeypottype" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"honeyPotType":""`)
				} else {
					cell = fmt.Sprintf(`"honeyPotType":"%s"`, value)
				}
			} else if columName == "serverip" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"serverip":""`)
				} else {
					cell = fmt.Sprintf(`"serverip":"%s"`, value)
				}
			} else if columName == "servername" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"servername":""`)
				} else {
					cell = fmt.Sprintf(`"servername":"%s"`, value)
				}
			} else if columName == "honeyip" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"honeyip":""`)
				} else {
					cell = fmt.Sprintf(`"honeyip":"%s"`, value)
				}
			} else if columName == "attacktime" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"attackTime":""`)
				} else {
					cell = fmt.Sprintf(`"attackTime":"%s"`, value)
				}
			} else if columName == "attackip" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"attackIP":""`)
				} else {
					cell = fmt.Sprintf(`"attackIP":"%s"`, value)
				}
			} else if columName == "probeip" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"probeIp":""`)
				} else {
					cell = fmt.Sprintf(`"probeIp":"%s"`, value)
				}
			} else if columName == "country" {
				if value != "" && value != "NULL" {
					location += value
				}
			} else if columName == "province" {
				sites := ""
				if value != "" && value != "NULL" {
					sites = fmt.Sprintf(`"location":"%s-%s"`, location, value)
				} else {
					if location != "" {
						sites = fmt.Sprintf(`"location":"%s"`, location)
					} else {
						sites = fmt.Sprintf(`"location":""`)
					}

				}
				row = row + sites + ","
			}
			if columName != "province" && columName != "country" {
				row = row + cell + ","
			}
		}
		row = row[0 : len(row)-1]
		row += "}"
		list = list + row + ","
	}
	if count != 0 {
		list = list[0 : len(list)-1]
	}
	list += "],"
	cell := fmt.Sprintf(`"total":%d`, all)
	totalpage := fmt.Sprintf(`"totalPage":%d`, totalPage)
	pagenum := fmt.Sprintf(`"pageNum":%d`, pageNum)
	pagesize := fmt.Sprintf(`"pageSize":%d`, pageSize)
	list = list + cell + "," + totalpage + "," + pagenum + "," + pagesize
	list += "}"

	list = strings.Replace(list, "\\", "", -1)
	return list, count, nil
}

func GetAttackLogDetailMysqlJson(rows *sql.Rows, columns []string, all int, values []sql.RawBytes, scanArgs []interface{}, pageNum int, pageSize int, totalPage int) (string, int, error) {
	list := "{\"list\":["
	count := 0

	for rows.Next() {
		count += 1
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Println("scan failed, err: ", err)
			return "", count, err
		}

		row := "{"
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}

			columName := strings.ToLower(columns[i])
			cell := ""
			if columName == "honeypottype" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"honeyPotType":""`)
				} else {
					cell = fmt.Sprintf(`"honeyPotType":"%s"`, value)
				}
			} else if columName == "serverip" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"honeyserverip":""`)
				} else {
					cell = fmt.Sprintf(`"honeyserverip":"%s"`, value)
				}
			} else if columName == "honeyip" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"honeyip":""`)
				} else {
					cell = fmt.Sprintf(`"honeyip":"%s"`, value)
				}
			} else if columName == "honeypotport" {
				honeypotport, _ := strconv.Atoi(value)
				cell = fmt.Sprintf(`"honeyPotPort":%d`, honeypotport)
			} else if columName == "srcport" {
				honeypotport, _ := strconv.Atoi(value)
				cell = fmt.Sprintf(`"srcPort":%d`, honeypotport)
			} else if columName == "attacktime" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"attackTime":""`)
				} else {
					cell = fmt.Sprintf(`"attackTime":"%s"`, value)
				}
			} else if columName == "eventdetail" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"eventDetail":""`)
				} else {
					value = strings.Replace(value, "\\r\\n", "", -1)
					value = strings.Replace(value, "\r\n", "", -1)
					value = strings.Replace(value, "\r", "", -1)
					value = strings.Replace(value, "\n", "", -1)
					value = strings.Replace(value, "\"", "\\\"", -1)
					value = strings.Replace(value, "\\", "\\\\", -1)
					value = DoStr(value)
					cell = fmt.Sprintf(`"eventDetail":"%s"`, value)
				}
			} else if columName == "logdata" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"logData":""`)
				} else {
					//value = strings.Replace(value, "\\\\r\\\\n", "", -1)
					//value = strings.Replace(value, "\\r\\n", "", -1)
					//value = strings.Replace(value, "\\r", "", -1)
					//value = strings.Replace(value, "\\n", "", -1)
					//value = strings.Replace(value, "\"", "\\\"", -1)
					//value = strings.Replace(value, "\\\\", "\\\\", -1)
					cell = fmt.Sprintf(`"logData":%s`, value)
				}
			} else if columName == "srchost" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"srcHost":""`)
				} else {
					cell = fmt.Sprintf(`"srcHost":"%s"`, value)
				}
			} else if columName == "attackip" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"attackIP":""`)
				} else {
					cell = fmt.Sprintf(`"attackIP":"%s"`, value)
				}
			}

			row = row + cell + ","
		}
		row = row[0 : len(row)-1]
		row += "}"
		list = list + row + ","

	}

	if count != 0 {
		list = list[0 : len(list)-1]
	}
	list += "],"
	cell := fmt.Sprintf(`"total":%d`, all)
	totalpage := fmt.Sprintf(`"totalPage":%d`, totalPage)
	pagenum := fmt.Sprintf(`"pageNum":%d`, pageNum)
	pagesize := fmt.Sprintf(`"pageSize":%d`, pageSize)
	list = list + cell + "," + totalpage + "," + pagenum + "," + pagesize
	list += "}"
	return list, count, nil
}

func GetConfigJson(rows *sql.Rows, columns []string, all int, values []sql.RawBytes, scanArgs []interface{}, pageNum int, pageSize int, totalPage int) (string, int, error) {
	list := "{\"list\":["
	count := 0

	for rows.Next() {
		count += 1
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Println("scan failed, err: ", err)
			return "", count, err
		}

		row := "{"
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}

			columName := strings.ToLower(columns[i])
			cell := ""
			if columName == "confname" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"confName":""`)
				} else {
					cell = fmt.Sprintf(`"confName":"%s"`, value)
				}
			} else if columName == "confvalue" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"confValue":""`)
				} else {
					cell = fmt.Sprintf(`"confValue":"%s"`, value)
				}
			} else if columName == "id" {
				id, _ := strconv.Atoi(value)
				cell = fmt.Sprintf(`"id":%d`, id)
			}

			row = row + cell + ","
		}
		row = row[0 : len(row)-1]
		row += "}"
		list = list + row + ","

	}

	if count != 0 {
		list = list[0 : len(list)-1]
	}
	list += "],"
	cell := fmt.Sprintf(`"total":%d`, all)
	totalpage := fmt.Sprintf(`"totalPage":%d`, totalPage)
	pagenum := fmt.Sprintf(`"pageNum":%d`, pageNum)
	pagesize := fmt.Sprintf(`"pageSize":%d`, pageSize)
	list = list + cell + "," + totalpage + "," + pagenum + "," + pagesize
	list += "}"

	list = strings.Replace(list, "\\", "", -1)
	log.Println(list)
	return list, count, nil
}

func GetApplicationClustersListMysqlJson(rows *sql.Rows, columns []string, all int, values []sql.RawBytes, scanArgs []interface{}, pageNum int, pageSize int, totalPage int) (string, int, error) {
	list := "{\"list\":["
	count := 0
	nowtime := time.Now().Unix()
	//regtime := 0

	for rows.Next() {
		count += 1
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Println("scan failed, err: ", err)
			return "", count, err
		}

		row := "{"
		var value string
		var chatime int64 = 0
		var heartbeattime int64 = 0
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}

			columName := strings.ToLower(columns[i])
			cell := ""
			if columName == "id" {
				id, _ := strconv.Atoi(value)
				cell = fmt.Sprintf(`"id":%d`, id)
			} else if columName == "servername" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"servername":""`)
				} else {
					cell = fmt.Sprintf(`"servername":"%s"`, value)
				}
			} else if columName == "serverip" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"serverip":""`)
				} else {
					cell = fmt.Sprintf(`"serverip":"%s"`, value)
				}
			} else if columName == "sys" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"sys":""`)
				} else {
					cell = fmt.Sprintf(`"sys":"%s"`, value)
				}
			} else if columName == "serverid" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"serverid":""`)
				} else {
					cell = fmt.Sprintf(`"serverid":"%s"`, value)
				}
			} else if columName == "agentid" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"agentid":""`)
				} else {
					cell = fmt.Sprintf(`"agentid":"%s"`, value)
				}
			} else if columName == "regtime" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"regtime":""`)
				} else {
					cell = fmt.Sprintf(`"regtime":"%s"`, value)
				}
			} else if columName == "heartbeattime" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"heartbeattime":""`)
				} else {
					cell = fmt.Sprintf(`"heartbeattime":"%s"`, value)
					heartbeattime, err = strconv.ParseInt(value, 10, 64)
					chatime = nowtime - heartbeattime
					if err != nil {
						logs.Error("GetApplicationClustersListMysqlJson regtime Parse to Int64 Error:", err)
					}
				}
			} else if columName == "status" {
				if chatime >= 300 {
					cell = fmt.Sprintf(`"status":%d`, 2)
				} else {
					cell = fmt.Sprintf(`"status":%d`, 1)
				}
			} else if columName == "vpcname" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"vpcname":""`)
				} else {
					cell = fmt.Sprintf(`"vpcname":"%s"`, value)
				}
			} else if columName == "vpsowner" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"vpsowner":""`)
				} else {
					cell = fmt.Sprintf(`"vpsowner":"%s"`, value)
				}
			}
			row = row + cell + ","
			//if columName == "attackip" {
			//	sites := fmt.Sprintf(`"location":"%s"`,site)
			//	row = row + cell + "," + sites + ","
			//}else {
			//	row = row + cell + ","
			//}
		}
		row = row[0 : len(row)-1]
		row += "}"
		list = list + row + ","
	}
	if count != 0 {
		list = list[0 : len(list)-1]
	}
	list += "],"
	cell := fmt.Sprintf(`"total":%d`, all)
	totalpage := fmt.Sprintf(`"totalPage":%d`, totalPage)
	pagenum := fmt.Sprintf(`"pageNum":%d`, pageNum)
	pagesize := fmt.Sprintf(`"pageSize":%d`, pageSize)
	list = list + cell + "," + totalpage + "," + pagenum + "," + pagesize
	list += "}"

	list = strings.Replace(list, "\\", "", -1)
	log.Println(list)
	return list, count, nil
}

func GetSignListMysqlJson(rows *sql.Rows, columns []string, all int, values []sql.RawBytes, scanArgs []interface{}, pageNum int, pageSize int, totalPage int) (string, int, error) {
	list := "{\"list\":["
	count := 0

	for rows.Next() {
		count += 1
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Println("scan failed, err: ", err)
			return "", count, err
		}

		row := "{"
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			//site := ""
			columName := strings.ToLower(columns[i])
			cell := ""
			if columName == "id" {
				id, _ := strconv.Atoi(value)
				cell = fmt.Sprintf(`"id":%d`, id)
			} else if columName == "signtype" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"signtype":""`)
				} else {
					cell = fmt.Sprintf(`"signtype":"%s"`, value)
				}
			} else if columName == "signname" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"signname":""`)
				} else {
					cell = fmt.Sprintf(`"signname":"%s"`, value)
				}
			} else if columName == "createtime" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"createtime":""`)
				} else {
					cell = fmt.Sprintf(`"createtime":"%s"`, value)
				}
			} else if columName == "creator" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"creator":""`)
				} else {
					cell = fmt.Sprintf(`"creator":"%s"`, value)
				}
			} else if columName == "signid" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"signid":""`)
				} else {
					cell = fmt.Sprintf(`"signid":"%s"`, value)
				}
			} else if columName == "signinfo" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"signinfo":""`)
				} else {
					cell = fmt.Sprintf(`"signinfo":"%s"`, value)
				}
			} else if columName == "signsystype" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"signsystype":""`)
				} else {
					cell = fmt.Sprintf(`"signsystype":"%s"`, value)
				}
			}
			row = row + cell + ","
		}
		row = row[0 : len(row)-1]
		row += "}"
		list = list + row + ","
	}
	if count != 0 {
		list = list[0 : len(list)-1]
	}
	list += "],"
	cell := fmt.Sprintf(`"total":%d`, all)
	totalpage := fmt.Sprintf(`"totalPage":%d`, totalPage)
	pagenum := fmt.Sprintf(`"pageNum":%d`, pageNum)
	pagesize := fmt.Sprintf(`"pageSize":%d`, pageSize)
	list = list + cell + "," + totalpage + "," + pagenum + "," + pagesize
	list += "}"

	list = strings.Replace(list, "\\", "", -1)
	log.Println(list)
	return list, count, nil
}

func GetSignListByTypeMysqlJson(rows *sql.Rows, columns []string, all int, values []sql.RawBytes, scanArgs []interface{}, pageNum int, pageSize int, totalPage int) (string, int, error) {
	list := "{\"list\":["
	count := 0

	for rows.Next() {
		count += 1
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Println("scan failed, err: ", err)
			return "", count, err
		}

		row := "{"
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			//site := ""
			columName := strings.ToLower(columns[i])
			cell := ""
			if columName == "signtype" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"signtype":""`)
				} else {
					cell = fmt.Sprintf(`"signtype":"%s"`, value)
				}
			} else if columName == "signinfo" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"signname":""`)
				} else {
					cell = fmt.Sprintf(`"signname":"%s"`, value)
				}
			} else if columName == "signid" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"signid":""`)
				} else {
					cell = fmt.Sprintf(`"signid":"%s"`, value)
				}
			}
			row = row + cell + ","
		}
		row = row[0 : len(row)-1]
		row += "}"
		list = list + row + ","
	}
	if count != 0 {
		list = list[0 : len(list)-1]
	}
	list += "],"
	cell := fmt.Sprintf(`"total":%d`, all)
	totalpage := fmt.Sprintf(`"totalPage":%d`, totalPage)
	pagenum := fmt.Sprintf(`"pageNum":%d`, pageNum)
	pagesize := fmt.Sprintf(`"pageSize":%d`, pageSize)
	list = list + cell + "," + totalpage + "," + pagenum + "," + pagesize
	list += "}"

	list = strings.Replace(list, "\\", "", -1)
	log.Println(list)
	return list, count, nil
}

func GetBaitListMysqlJson(rows *sql.Rows, columns []string, all int, values []sql.RawBytes, scanArgs []interface{}, pageNum int, pageSize int, totalPage int) (string, int, error) {
	list := "{\"list\":["
	count := 0

	for rows.Next() {
		count += 1
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Println("scan failed, err: ", err)
			return "", count, err
		}

		row := "{"
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			//site := ""
			columName := strings.ToLower(columns[i])
			cell := ""
			if columName == "id" {
				id, _ := strconv.Atoi(value)
				cell = fmt.Sprintf(`"id":%d`, id)
			} else if columName == "baittype" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"baittype":""`)
				} else {
					cell = fmt.Sprintf(`"baittype":"%s"`, value)
				}
			} else if columName == "baitname" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"baitname":""`)
				} else {
					cell = fmt.Sprintf(`"baitname":"%s"`, value)
				}
			} else if columName == "createtime" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"createtime":""`)
				} else {
					cell = fmt.Sprintf(`"createtime":"%s"`, value)
				}
			} else if columName == "creator" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"creator":""`)
				} else {
					cell = fmt.Sprintf(`"creator":"%s"`, value)
				}
			} else if columName == "baitid" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"baitid":""`)
				} else {
					cell = fmt.Sprintf(`"baitid":"%s"`, value)
				}
			} else if columName == "baitsystype" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"baitsystype":""`)
				} else {
					cell = fmt.Sprintf(`"baitsystype":"%s"`, value)
				}
			} else if columName == "systype" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"systype":""`)
				} else {
					cell = fmt.Sprintf(`"systype":"%s"`, value)
				}
			} else if columName == "baitinfo" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"baitinfo":""`)
				} else {
					cell = fmt.Sprintf(`"baitinfo":"%s"`, value)
				}
			} else if columName == "md5" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"md5":""`)
				} else {
					cell = fmt.Sprintf(`"md5":"%s"`, value)
				}
			}
			row = row + cell + ","
			//if columName == "attackip" {
			//	sites := fmt.Sprintf(`"location":"%s"`,site)
			//	row = row + cell + "," + sites + ","
			//}else {
			//	row = row + cell + ","
			//}
		}
		row = row[0 : len(row)-1]
		row += "}"
		list = list + row + ","
	}
	if count != 0 {
		list = list[0 : len(list)-1]
	}
	list += "],"
	cell := fmt.Sprintf(`"total":%d`, all)
	totalpage := fmt.Sprintf(`"totalPage":%d`, totalPage)
	pagenum := fmt.Sprintf(`"pageNum":%d`, pageNum)
	pagesize := fmt.Sprintf(`"pageSize":%d`, pageSize)
	list = list + cell + "," + totalpage + "," + pagenum + "," + pagesize
	list += "}"

	list = strings.Replace(list, "\\", "", -1)
	log.Println(list)
	return list, count, nil
}

func GetPodImageListMysqlJson(rows *sql.Rows, columns []string, all int, values []sql.RawBytes, scanArgs []interface{}, pageNum int, pageSize int, totalPage int) (string, int, error) {
	list := "{\"list\":["
	count := 0

	for rows.Next() {
		count += 1
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Println("scan failed, err: ", err)
			return "", count, err
		}

		row := "{"
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			//site := ""
			columName := strings.ToLower(columns[i])
			cell := ""
			if columName == "id" {
				id, _ := strconv.Atoi(value)
				cell = fmt.Sprintf(`"id":%d`, id)
				row = row + cell + ","
			} else if columName == "imageaddress" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"podimage":""`)
				} else {
					cell = fmt.Sprintf(`"podimage":"%s"`, value)
				}
				row = row + cell + ","
			} else if columName == "imagename" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"podname":""`)
				} else {
					cell = fmt.Sprintf(`"podname":"%s"`, value)
				}
				row = row + cell + ","
			} else if columName == "repository" {
				if value == "" || value == "NULL" {
					cell = ""
				} else {
					cell = fmt.Sprintf(`"repository":"%s"`, value)
					row = row + cell + ","
				}
			}

		}
		row = row[0 : len(row)-1]
		row += "}"
		list = list + row + ","
	}
	if count != 0 {
		list = list[0 : len(list)-1]
	}
	list += "],"
	cell := fmt.Sprintf(`"total":%d`, all)
	totalpage := fmt.Sprintf(`"totalPage":%d`, totalPage)
	pagenum := fmt.Sprintf(`"pageNum":%d`, pageNum)
	pagesize := fmt.Sprintf(`"pageSize":%d`, pageSize)
	list = list + cell + "," + totalpage + "," + pagenum + "," + pagesize
	list += "}"

	list = strings.Replace(list, "\\", "", -1)
	//log.Println(list)
	return list, count, nil
}

func GetHoneyClustersListMysqlJson(rows *sql.Rows, columns []string, all int, values []sql.RawBytes, scanArgs []interface{}, pageNum int, pageSize int, totalPage int) (string, int, error) {
	list := "{\"list\":["
	count := 0

	for rows.Next() {
		count += 1
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Println("scan failed, err: ", err)
			return "", count, err
		}

		row := "{"
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			//site := ""
			columName := strings.ToLower(columns[i])
			cell := ""
			if columName == "id" {
				id, _ := strconv.Atoi(value)
				cell = fmt.Sprintf(`"id":%d`, id)
			} else if columName == "servername" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"servername":""`)
				} else {
					cell = fmt.Sprintf(`"servername":"%s"`, value)
				}
			} else if columName == "serverip" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"serverip":""`)
				} else {
					cell = fmt.Sprintf(`"serverip":"%s"`, value)
				}
			} else if columName == "serverid" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"serverid":""`)
				} else {
					cell = fmt.Sprintf(`"serverid":"%s"`, value)
				}
			} else if columName == "status" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"status":""`)
				} else {
					status, _ := strconv.Atoi(value)
					cell = fmt.Sprintf(`"status":%d`, status)
				}
			} else if columName == "agentid" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"agentid":""`)
				} else {
					cell = fmt.Sprintf(`"agentid":"%s"`, value)
				}
			} else if columName == "regtime" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"regtime":""`)
				} else {
					cell = fmt.Sprintf(`"regtime":"%s"`, value)
				}
			} else if columName == "heartbeattime" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"heartbeattime":""`)
				} else {
					cell = fmt.Sprintf(`"heartbeattime":"%s"`, value)
				}
			}
			row = row + cell + ","
			//if columName == "attackip" {
			//	sites := fmt.Sprintf(`"location":"%s"`,site)
			//	row = row + cell + "," + sites + ","
			//}else {
			//	row = row + cell + ","
			//}
		}
		row = row[0 : len(row)-1]
		row += "}"
		list = list + row + ","
	}
	if count != 0 {
		list = list[0 : len(list)-1]
	}
	list += "],"
	cell := fmt.Sprintf(`"total":%d`, all)
	totalpage := fmt.Sprintf(`"totalPage":%d`, totalPage)
	pagenum := fmt.Sprintf(`"pageNum":%d`, pageNum)
	pagesize := fmt.Sprintf(`"pageSize":%d`, pageSize)
	list = list + cell + "," + totalpage + "," + pagenum + "," + pagesize
	list += "}"

	list = strings.Replace(list, "\\", "", -1)
	log.Println(list)
	return list, count, nil
}

func GetHoneyInfosListMysqlJson(rows *sql.Rows, columns []string, all int, values []sql.RawBytes, scanArgs []interface{}, pageNum int, pageSize int, totalPage int) (string, int, error) {
	list := "{\"list\":["
	count := 0

	for rows.Next() {
		count += 1
		err := rows.Scan(scanArgs...)
		if err != nil {
			logs.Error("[GetHoneyInfosListMysqlJson]s can failed, err: ", err)
			return "", count, err
		}

		row := "{"
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			//site := ""
			columName := strings.ToLower(columns[i])
			cell := ""
			if columName == "id" {
				id, _ := strconv.Atoi(value)
				cell = fmt.Sprintf(`"id":%d`, id)
			} else if columName == "serverid" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"serverid":""`)
				} else {
					cell = fmt.Sprintf(`"serverid":"%s"`, value)
				}
			} else if columName == "podname" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"podname":""`)
				} else {
					cell = fmt.Sprintf(`"podname":"%s"`, value)
				}
			} else if columName == "honeyport" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"honeyport":""`)
				} else {
					cell = fmt.Sprintf(`"honeyport":"%s"`, value)
				}
			} else if columName == "honeyimage" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"honeyimage":""`)
				} else {
					cell = fmt.Sprintf(`"honeyimage":"%s"`, value)
				}
			} else if columName == "honeynamespce" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"honeynamespce":""`)
				} else {
					cell = fmt.Sprintf(`"honeynamespce":"%s"`, value)
				}
			} else if columName == "honeyname" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"honeyname":""`)
				} else {
					cell = fmt.Sprintf(`"honeyname":"%s"`, value)
				}
			} else if columName == "honeytypeid" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"honeytypeid":""`)
				} else {
					cell = fmt.Sprintf(`"honeytypeid":"%s"`, value)
				}
			} else if columName == "honeypotid" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"honeypotid":""`)
				} else {
					cell = fmt.Sprintf(`"honeypotid":"%s"`, value)
				}
			} else if columName == "honeyip" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"honeyip":""`)
				} else {
					cell = fmt.Sprintf(`"honeyip":"%s"`, value)
				}
			} else if columName == "createtime" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"createtime":""`)
				} else {
					cell = fmt.Sprintf(`"createtime":"%s"`, value)
				}
			} else if columName == "creator" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"creator":""`)
				} else {
					cell = fmt.Sprintf(`"creator":"%s"`, value)
				}
			} else if columName == "sysid" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"sysid":""`)
				} else {
					cell = fmt.Sprintf(`"sysid":"%s"`, value)
				}
			} else if columName == "status" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"status":""`)
				} else {
					status, _ := strconv.Atoi(value)
					cell = fmt.Sprintf(`"status":%d`, status)
				}
			} else if columName == "serverip" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"serverip":""`)
				} else {
					cell = fmt.Sprintf(`"serverip":"%s"`, value)
				}
			} else if columName == "hostip" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"hostip":""`)
				} else {
					cell = fmt.Sprintf(`"hostip":"%s"`, value)
				}
			} else if columName == "servername" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"servername":""`)
				} else {
					cell = fmt.Sprintf(`"servername":"%s"`, value)
				}
			} else if columName == "offlinetime" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"offlinetime":""`)
				} else {
					cell = fmt.Sprintf(`"offlinetime":"%s"`, value)
				}
			} else if columName == "honeypottype" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"honeypottype":""`)
				} else {
					cell = fmt.Sprintf(`"honeypottype":"%s"`, value)
				}
			} else if columName == "systype" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"systype":""`)
				} else {
					cell = fmt.Sprintf(`"systype":"%s"`, value)
				}
			} else if columName == "agentid" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"agentid":""`)
				} else {
					cell = fmt.Sprintf(`"agentid":"%s"`, value)
				}
			}

			row = row + cell + ","
		}
		row = row[0 : len(row)-1]
		row += "}"
		list = list + row + ","
	}
	if count != 0 {
		list = list[0 : len(list)-1]
	}
	list += "],"
	cell := fmt.Sprintf(`"total":%d`, all)
	totalpage := fmt.Sprintf(`"totalPage":%d`, totalPage)
	pagenum := fmt.Sprintf(`"pageNum":%d`, pageNum)
	pagesize := fmt.Sprintf(`"pageSize":%d`, pageSize)
	list = list + cell + "," + totalpage + "," + pagenum + "," + pagesize
	list += "}"

	list = strings.Replace(list, "\\", "", -1)
	fmt.Println("list:", list)
	//log.Println(list)
	return list, count, nil
}

func GetHoneyImagesListMysqlJson(rows *sql.Rows, columns []string, all int, values []sql.RawBytes, scanArgs []interface{}, pageNum int, pageSize int, totalPage int) (string, int, error) {
	list := "{\"list\":["
	count := 0

	for rows.Next() {
		count += 1
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Println("scan failed, err: ", err)
			return "", count, err
		}

		row := "{"
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			//site := ""
			columName := strings.ToLower(columns[i])
			cell := ""
			if columName == "id" {
				id, _ := strconv.Atoi(value)
				cell = fmt.Sprintf(`"id":%d`, id)
			} else if columName == "imageaddress" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"imageaddress":""`)
				} else {
					cell = fmt.Sprintf(`"imageaddress":"%s"`, value)
				}
			} else if columName == "repository" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"repository":""`)
				} else {
					cell = fmt.Sprintf(`"repository":"%s"`, value)
				}
			} else if columName == "imagename" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"imagename":""`)
				} else {
					cell = fmt.Sprintf(`"imagename":"%s"`, value)
				}
			} else if columName == "imageport" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"imageport":null`)
				} else {
					port, _ := strconv.Atoi(value)
					cell = fmt.Sprintf(`"imageport":%d`, port)
				}
			} else if columName == "imagetype" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"imagetype":""`)
				} else {
					cell = fmt.Sprintf(`"imagetype":"%s"`, value)
				}
			} else if columName == "imageos" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"imageos":""`)
				} else {
					cell = fmt.Sprintf(`"imageos":"%s"`, value)
				}
			}
			row = row + cell + ","
		}
		row = row[0 : len(row)-1]
		row += "}"
		list = list + row + ","
	}
	if count != 0 {
		list = list[0 : len(list)-1]
	}
	list += "],"
	cell := fmt.Sprintf(`"total":%d`, all)
	totalpage := fmt.Sprintf(`"totalPage":%d`, totalPage)
	pagenum := fmt.Sprintf(`"pageNum":%d`, pageNum)
	pagesize := fmt.Sprintf(`"pageSize":%d`, pageSize)
	list = list + cell + "," + totalpage + "," + pagenum + "," + pagesize
	list += "}"
	list = strings.Replace(list, "\\", "", -1)
	return list, count, nil
}

func GetHoneyListMysqlJson(rows *sql.Rows, columns []string, all int, values []sql.RawBytes, scanArgs []interface{}, pageNum int, pageSize int, totalPage int) (string, int, error) {
	list := "{\"list\":["
	count := 0

	for rows.Next() {
		count += 1
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Println("scan failed, err: ", err)
			return "", count, err
		}

		row := "{"
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			//site := ""
			columName := strings.ToLower(columns[i])
			cell := ""
			if columName == "honeyname" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"honeyname":""`)
				} else {
					cell = fmt.Sprintf(`"honeyname":"%s"`, value)
				}
			} else if columName == "honeypotid" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"honeypotid":""`)
				} else {
					cell = fmt.Sprintf(`"honeypotid":"%s"`, value)
				}
			}

			row = row + cell + ","
		}
		row = row[0 : len(row)-1]
		row += "}"
		list = list + row + ","
	}
	if count != 0 {
		list = list[0 : len(list)-1]
	}
	list += "],"
	cell := fmt.Sprintf(`"total":%d`, all)
	totalpage := fmt.Sprintf(`"totalPage":%d`, totalPage)
	pagenum := fmt.Sprintf(`"pageNum":%d`, pageNum)
	pagesize := fmt.Sprintf(`"pageSize":%d`, pageSize)
	list = list + cell + "," + totalpage + "," + pagenum + "," + pagesize
	list += "}"

	list = strings.Replace(list, "\\", "", -1)
	//log.Println(list)
	return list, count, nil
}

func GetBaitPolicyCheckCreateStatusMysqlJson(rows *sql.Rows) ([]BaitPolicyByStatus, error) {
	var baitpolicylist []BaitPolicyByStatus

	for rows.Next() {
		var baitpolicy BaitPolicyByStatus
		err := rows.Scan(&baitpolicy.TaskId, &baitpolicy.AgentId, &baitpolicy.Status, &baitpolicy.CreateTime, &baitpolicy.Data, &baitpolicy.Md5, &baitpolicy.Type)
		if err != nil {
			logs.Error("[GetBaitPolicyCheckCreateStatusMysqlJson] scan failed, err: ", err)
			return nil, err
		}
		baitpolicylist = append(baitpolicylist, baitpolicy)
	}

	return baitpolicylist, nil
}
func GetBaitPolicyCheckOfflineStatusMysqlJson(rows *sql.Rows) ([]BaitPolicyByStatus, error) {
	var baitpolicylist []BaitPolicyByStatus

	for rows.Next() {
		var baitpolicy BaitPolicyByStatus
		err := rows.Scan(&baitpolicy.TaskId, &baitpolicy.AgentId, &baitpolicy.Status, &baitpolicy.OfflineTime, &baitpolicy.Data, &baitpolicy.Md5, &baitpolicy.Type)
		if err != nil {
			logs.Error("[GetBaitPolicyCheckOfflineStatusMysqlJson] scan failed, err: ", err)
			return nil, err
		}
		baitpolicylist = append(baitpolicylist, baitpolicy)
	}

	return baitpolicylist, nil
}

// 创建透明转发、蜜罐流量转发策略SQL to json
func GetTransPolicyCheckCreateStatusMysqlJson(rows *sql.Rows) ([]TransPolicyByStatus, error) {
	var transpolicylist []TransPolicyByStatus

	for rows.Next() {
		var transpolicy TransPolicyByStatus
		err := rows.Scan(&transpolicy.TaskId, &transpolicy.AgentId, &transpolicy.ForwardPort, &transpolicy.HoneyPotPort, &transpolicy.CreateTime, &transpolicy.Status, &transpolicy.HoneyIP, &transpolicy.HoneyPotType, &transpolicy.Type, &transpolicy.Path)
		if err != nil {
			logs.Error("[GetTransPolicyCheckCreateStatusMysqlJson] scan failed, err: ", err)
			return nil, err
		}
		transpolicylist = append(transpolicylist, transpolicy)
	}

	return transpolicylist, nil
}

// 下线透明转发、蜜罐流量转发策略SQL to json
func GetTransPolicyCheckOfflineStatusMysqlJson(rows *sql.Rows) ([]TransPolicyByStatus, error) {
	var transpolicylist []TransPolicyByStatus

	for rows.Next() {
		var transpolicy TransPolicyByStatus
		err := rows.Scan(&transpolicy.TaskId, &transpolicy.AgentId, &transpolicy.ForwardPort, &transpolicy.HoneyPotPort, &transpolicy.OfflineTime, &transpolicy.Status, &transpolicy.HoneyIP, &transpolicy.HoneyPotType, &transpolicy.Type, &transpolicy.Path)
		if err != nil {
			logs.Error("[GetTransPolicyCheckOfflineStatusMysqlJson] scan failed, err: ", err)
			return nil, err
		}
		transpolicylist = append(transpolicylist, transpolicy)
	}

	return transpolicylist, nil
}

func GetHoneySignsMsgMysqlJson(rows *sql.Rows, columns []string, all int, values []sql.RawBytes, scanArgs []interface{}, pageNum int, pageSize int, totalPage int) (string, int, error) {
	list := "{\"list\":["
	count := 0

	for rows.Next() {
		count += 1
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Println("scan failed, err: ", err)
			return "", count, err
		}
		row := "{"
		var value string
		signf := ""
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			//site := ""
			columName := strings.ToLower(columns[i])
			cell := ""

			if columName == "id" {
				id, _ := strconv.Atoi(value)
				cell = fmt.Sprintf(`"id":%d`, id)
			} else if columName == "ipcity" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"ipcity":""`)
				} else {
					cell = fmt.Sprintf(`"ipcity":"%s"`, value)
				}
			} else if columName == "ipcountry" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"ipcountry":""`)
				} else {
					cell = fmt.Sprintf(`"ipcountry":"%s"`, value)
				}
			} else if columName == "latitude" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"latitude":""`)
				} else {
					cell = fmt.Sprintf(`"latitude":"%s"`, value)
				}
			} else if columName == "longitude" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"longitude":""`)
				} else {
					cell = fmt.Sprintf(`"longitude":"%s"`, value)
				}
			} else if columName == "openip" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"openip":""`)
				} else {
					cell = fmt.Sprintf(`"openip":"%s"`, value)
				}
			} else if columName == "opentime" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"opentime":""`)
				} else {
					cell = fmt.Sprintf(`"opentime":%s`, value)
				}
			} else if columName == "tracecode" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"tracecode":""`)
				} else {
					cell = fmt.Sprintf(`"tracecode":"%s"`, value)
				}
			} else if columName == "tracefilename" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"tracefilename":""`)
				} else {
					cell = fmt.Sprintf(`"tracefilename":"%s"`, value)
				}
			} else if columName == "useragent" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"useragent":""`)
				} else {
					cell = fmt.Sprintf(`"useragent":"%s"`, value)
				}
			} else if columName == "signfilename" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"signfilename":""`)
				} else {
					signf = value
					cell = fmt.Sprintf(`"signfilename":"%s"`, value)
				}
			} else if columName == "signinfo" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"signinfo":""`)
				} else {
					cell = fmt.Sprintf(`"signinfo":"%s"`, value+"/"+signf)
				}
			}

			row = row + cell + ","
		}
		row = row[0 : len(row)-1]
		row += "}"
		list = list + row + ","
	}
	if count != 0 {
		list = list[0 : len(list)-1]
	}
	list += "],"
	cell := fmt.Sprintf(`"total":%d`, all)
	totalpage := fmt.Sprintf(`"totalPage":%d`, totalPage)
	pagenum := fmt.Sprintf(`"pageNum":%d`, pageNum)
	pagesize := fmt.Sprintf(`"pageSize":%d`, pageSize)
	list = list + cell + "," + totalpage + "," + pagenum + "," + pagesize
	list += "}"

	list = strings.Replace(list, "\\", "", -1)
	//log.Println(list)
	return list, count, nil
}

func GetHoneySignsListMysqlJson(rows *sql.Rows, columns []string, all int, values []sql.RawBytes, scanArgs []interface{}, pageNum int, pageSize int, totalPage int) (string, int, error) {
	list := "{\"list\":["
	count := 0

	for rows.Next() {
		count += 1
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Println("scan failed, err: ", err)
			return "", count, err
		}
		row := "{"
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			//site := ""
			columName := strings.ToLower(columns[i])
			cell := ""
			if columName == "id" {
				id, _ := strconv.Atoi(value)
				cell = fmt.Sprintf(`"id":%d`, id)
			} else if columName == "signinfo" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"signinfo":""`)
				} else {
					cell = fmt.Sprintf(`"signinfo":"%s"`, value)
				}
			} else if columName == "signfilename" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"signfilename":""`)
				} else {
					cell = fmt.Sprintf(`"signfilename":"%s"`, value)
				}
			} else if columName == "signid" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"signid":""`)
				} else {
					cell = fmt.Sprintf(`"signid":"%s"`, value)
				}
			} else if columName == "honeypotid" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"honeypotid":""`)
				} else {
					cell = fmt.Sprintf(`"honeypotid":"%s"`, value)
				}
			} else if columName == "createtime" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"createtime":""`)
				} else {
					cell = fmt.Sprintf(`"createtime":"%s"`, value)
				}
			} else if columName == "offlinetime" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"offlinetime":""`)
				} else {
					cell = fmt.Sprintf(`"offlinetime":"%s"`, value)
				}
			} else if columName == "creator" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"creator":""`)
				} else {
					cell = fmt.Sprintf(`"creator":"%s"`, value)
				}
			} else if columName == "status" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"status":""`)
				} else {
					status, _ := strconv.Atoi(value)
					cell = fmt.Sprintf(`"status":%d`, status)
				}
			} else if columName == "signname" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"signname":""`)
				} else {
					cell = fmt.Sprintf(`"signname":"%s"`, value)
				}
			} else if columName == "signtype" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"signtype":""`)
				} else {
					cell = fmt.Sprintf(`"signtype":"%s"`, value)
				}
			} else if columName == "taskid" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"taskid":""`)
				} else {
					cell = fmt.Sprintf(`"taskid":"%s"`, value)
				}
			} else if columName == "tracecode" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"tracecode":""`)
				} else {
					cell = fmt.Sprintf(`"tracecode":"%s"`, value)
				}
			}

			row = row + cell + ","
		}
		row = row[0 : len(row)-1]
		row += "}"
		list = list + row + ","
	}
	if count != 0 {
		list = list[0 : len(list)-1]
	}
	list += "],"
	cell := fmt.Sprintf(`"total":%d`, all)
	totalpage := fmt.Sprintf(`"totalPage":%d`, totalPage)
	pagenum := fmt.Sprintf(`"pageNum":%d`, pageNum)
	pagesize := fmt.Sprintf(`"pageSize":%d`, pageSize)
	list = list + cell + "," + totalpage + "," + pagenum + "," + pagesize
	list += "}"

	list = strings.Replace(list, "\\", "", -1)
	log.Println(list)
	return list, count, nil
}

func GetHoneyBaitsListMysqlJson(rows *sql.Rows, columns []string, all int, values []sql.RawBytes, scanArgs []interface{}, pageNum int, pageSize int, totalPage int) (string, int, error) {
	list := "{\"list\":["
	count := 0

	for rows.Next() {
		count += 1
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Println("scan failed, err: ", err)
			return "", count, err
		}
		row := "{"
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			//site := ""
			columName := strings.ToLower(columns[i])
			cell := ""
			if columName == "id" {
				id, _ := strconv.Atoi(value)
				cell = fmt.Sprintf(`"id":%d`, id)
			} else if columName == "baitinfo" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"baitinfo":""`)
				} else {
					cell = fmt.Sprintf(`"baitinfo":"%s"`, value)
				}
			} else if columName == "baitid" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"baitid":""`)
				} else {
					cell = fmt.Sprintf(`"baitid":"%s"`, value)
				}
			} else if columName == "honeypotid" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"honeypotid":""`)
				} else {
					cell = fmt.Sprintf(`"honeypotid":"%s"`, value)
				}
			} else if columName == "createtime" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"createtime":""`)
				} else {
					cell = fmt.Sprintf(`"createtime":"%s"`, value)
				}
			} else if columName == "offlinetime" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"offlinetime":""`)
				} else {
					cell = fmt.Sprintf(`"offlinetime":"%s"`, value)
				}
			} else if columName == "creator" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"creator":""`)
				} else {
					cell = fmt.Sprintf(`"creator":"%s"`, value)
				}
			} else if columName == "status" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"status":""`)
				} else {
					status, _ := strconv.Atoi(value)
					cell = fmt.Sprintf(`"status":%d`, status)
				}
			} else if columName == "baitname" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"baitname":""`)
				} else {
					cell = fmt.Sprintf(`"baitname":"%s"`, value)
				}
			} else if columName == "baittype" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"baittype":""`)
				} else {
					cell = fmt.Sprintf(`"baittype":"%s"`, value)
				}
			} else if columName == "data" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"data":""`)
				} else {
					cell = fmt.Sprintf(`"data":"%s"`, value)
				}
			} else if columName == "taskid" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"taskid":""`)
				} else {
					cell = fmt.Sprintf(`"taskid":"%s"`, value)
				}
			}

			row = row + cell + ","
			//if columName == "attackip" {
			//	sites := fmt.Sprintf(`"location":"%s"`,site)
			//	row = row + cell + "," + sites + ","
			//}else {
			//	row = row + cell + ","
			//}
		}
		row = row[0 : len(row)-1]
		row += "}"
		list = list + row + ","
	}
	if count != 0 {
		list = list[0 : len(list)-1]
	}
	list += "],"
	cell := fmt.Sprintf(`"total":%d`, all)
	totalpage := fmt.Sprintf(`"totalPage":%d`, totalPage)
	pagenum := fmt.Sprintf(`"pageNum":%d`, pageNum)
	pagesize := fmt.Sprintf(`"pageSize":%d`, pageSize)
	list = list + cell + "," + totalpage + "," + pagenum + "," + pagesize
	list += "}"

	list = strings.Replace(list, "\\", "", -1)
	log.Println(list)
	return list, count, nil
}

func GetClamavDataListMysqlJson(rows *sql.Rows, columns []string, all int, values []sql.RawBytes, scanArgs []interface{}, pageNum int, pageSize int, totalPage int) (string, int, error) {
	list := "{\"list\":["
	count := 0

	for rows.Next() {
		count += 1
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Println("scan failed, err: ", err)
			return "", count, err
		}

		row := "{"
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			//site := ""
			columName := strings.ToLower(columns[i])
			cell := ""
			if columName == "id" {
				id, _ := strconv.Atoi(value)
				cell = fmt.Sprintf(`"id":%d`, id)
			} else if columName == "filename" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"filename":""`)
				} else {
					cell = fmt.Sprintf(`"filename":"%s"`, value)
				}
			} else if columName == "virname" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"virname":""`)
				} else {
					cell = fmt.Sprintf(`"virname":"%s"`, value)
				}
			} else if columName == "honeypotip" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"honeypotip":""`)
				} else {
					cell = fmt.Sprintf(`"honeypotip":"%s"`, value)
				}
			} else if columName == "createtime" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"createtime":""`)
				} else {
					cell = fmt.Sprintf(`"createtime":"%s"`, value)
				}
			}
			row = row + cell + ","
		}
		row = row[0 : len(row)-1]
		row += "}"
		list = list + row + ","
	}
	if count != 0 {
		list = list[0 : len(list)-1]
	}
	list += "],"
	cell := fmt.Sprintf(`"total":%d`, all)
	totalpage := fmt.Sprintf(`"totalPage":%d`, totalPage)
	pagenum := fmt.Sprintf(`"pageNum":%d`, pageNum)
	pagesize := fmt.Sprintf(`"pageSize":%d`, pageSize)
	list = list + cell + "," + totalpage + "," + pagenum + "," + pagesize
	list += "}"

	list = strings.Replace(list, "\\", "", -1)

	return list, count, nil
}

func GetProtocolListMysqlJson(rows *sql.Rows, columns []string, all int, values []sql.RawBytes, scanArgs []interface{}, pageNum int, pageSize int, totalPage int) (string, int, error) {
	list := "{\"list\":["
	count := 0

	for rows.Next() {
		count += 1
		err := rows.Scan(scanArgs...)
		if err != nil {
			log.Println("scan failed, err: ", err)
			return "", count, err
		}

		row := "{"
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}
			columName := strings.ToLower(columns[i])
			cell := ""
			if columName == "id" {
				id, _ := strconv.Atoi(value)
				cell = fmt.Sprintf(`"id":%d`, id)
			} else if columName == "honeypottype" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"honeypottype":""`)
				} else {
					cell = fmt.Sprintf(`"honeypottype":"%s"`, value)
				}
			} else if columName == "typeid" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"typeid":""`)
				} else {
					cell = fmt.Sprintf(`"typeid":"%s"`, value)
				}
			} else if columName == "softpath" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"softpath":""`)
				} else {
					cell = fmt.Sprintf(`"softpath":"%s"`, value)
				}
			} else if columName == "createtime" {
				if value == "" || value == "NULL" {
					cell = fmt.Sprintf(`"createtime":""`)
				} else {
					cell = fmt.Sprintf(`"createtime":"%s"`, value)
				}
			}
			row = row + cell + ","
		}
		row = row[0 : len(row)-1]
		row += "}"
		list = list + row + ","
	}
	if count != 0 {
		list = list[0 : len(list)-1]
	}
	list += "],"
	cell := fmt.Sprintf(`"total":%d`, all)
	totalpage := fmt.Sprintf(`"totalPage":%d`, totalPage)
	pagenum := fmt.Sprintf(`"pageNum":%d`, pageNum)
	pagesize := fmt.Sprintf(`"pageSize":%d`, pageSize)
	list = list + cell + "," + totalpage + "," + pagenum + "," + pagesize
	list += "}"

	list = strings.Replace(list, "\\", "", -1)
	log.Println(list)
	return list, count, nil
}
