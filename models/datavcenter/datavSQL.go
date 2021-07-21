package datavcenter

import (
	"database/sql"
	"decept-defense/models/util"
	"decept-defense/models/util/comm"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"strings"
)

var (
	Dbhost     = beego.AppConfig.String("dbhost")
	Dbport     = beego.AppConfig.String("dbport")
	Dbuser     = beego.AppConfig.String("dbuser")
	Dbpassword = beego.AppConfig.String("dbpassword")
	Dbname     = beego.AppConfig.String("dbname")
)

func QueryAttackerInfoByIP(ip string) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT * from `attacker_info` WHERE ip=?", ip).Values(&maps)
	if err != nil {
		logs.Error("[SelectApplicationByAgentID] select event list error,%s", err)
	}
	return maps
}

func TryTransferIPToAgentIp(targetIP string) string {
	o := orm.NewOrm()
	var maps []orm.Params
	param := "%" + targetIP + "%"
	_, err := o.Raw("SELECT * from `servers` WHERE serverip like ?", param).Values(&maps)
	if err != nil {
		logs.Error("[TryTransferIPToAgentIp] select server list error,%s", err)
	}

	if len(maps) == 0 {
		return targetIP
	}

	for _, vMap := range maps {
		serverip := util.Strval(vMap["serverip"])
		if strings.Index(serverip, ",") > -1 {
			ips := strings.Split(serverip, ",")
			for _, ip := range ips {

				if ip == targetIP {
					return serverip
				}
			}
		} else {
			if serverip == targetIP {
				return serverip
			}
		}
	}

	return targetIP
}

func InsertAttackerInfo(sourceSite, account, ip, city string) (map[string]interface{}, string, int) {
	msg := "成功"
	var data map[string]interface{}
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("insert into attacker_info (source_site, account, ip, city) VALUES (?,?,?,?)", sourceSite, account, ip, city).Values(&maps)
	if err != nil {
		logs.Error("[InsertClamavData] insert servers error,%s", err)
		msg = "数据插入失败"
		return data, msg, comm.ErrorCode
	}
	return data, msg, comm.SuccessCode
}

func GetTopAttackMap() map[string]interface{} {
	sqlCon, err1 := sql.Open("mysql", Dbuser+":"+Dbpassword+"@tcp("+Dbhost+":"+Dbport+")/"+Dbname+"?charset=utf8&loc=Asia%2FShanghai")
	if err1 != nil {
		logs.Error("[SelectTopAttackIps]open mysql fail %s", err1)
	}
	var total int
	DbCon := sqlCon
	defer sqlCon.Close()
	sqlstr := "SELECT longitude,latitude,country,province,attackip,COUNT(attackip) AS attacksum FROM attacklog WHERE attackip !='' GROUP BY attackip ORDER BY attacktime desc LIMIT 0,30"

	rows, err := DbCon.Query(sqlstr)
	if err != nil {
		logs.Error("[GetTopAttackMap] select list error,%s", err)
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[GetTopAttackMap] rows.Columns() error,%s", err)
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var data map[string]interface{}

	o := orm.NewOrm()
	var maps []orm.Params
	_, err = o.Raw("SELECT * FROM `desipconf` ").Values(&maps)
	if err != nil {
		logs.Error("[SelectSignsById] select event list error,%s", err)
	}
	desip := "127.0.0.1"
	longitude := "局域网"
	latitude := "局域网"
	if len(maps) > 0 {
		desip = util.Strval(maps[0]["desip"])
		longitude = util.Strval(maps[0]["longitude"])
		latitude = util.Strval(maps[0]["latitude"])
	}

	list, count, err := util.GetTopAttackMapMysqlJson(rows, columns, total, values, scanArgs, util.Strval(desip), util.Strval(longitude), util.Strval(latitude))
	//fmt.Println("list:", list)
	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			logs.Error("[GetTopAttackMap] Unmarshal list error,%s", err)
		}
	}
	return data
}

func SelectTopAttackIps() map[string]interface{} {
	sqlCon, err1 := sql.Open("mysql", Dbuser+":"+Dbpassword+"@tcp("+Dbhost+":"+Dbport+")/"+Dbname+"?charset=utf8&loc=Asia%2FShanghai")
	if err1 != nil {
		logs.Error("[SelectTopAttackIps]open mysql fail %s", err1)
	}
	var total int
	DbCon := sqlCon
	defer sqlCon.Close()
	sqlstr := "SELECT attackip,COUNT(*) as ipcount FROM attacklog WHERE attackip !='' GROUP BY attackip ORDER BY COUNT(*) desc LIMIT 0,5"
	sqltotal := "SELECT COUNT(*) as ipsum FROM attacklog WHERE attackip !=''"
	sqltotal = sqltotal
	err := DbCon.QueryRow(sqltotal).Scan(&total)
	if err != nil {
		logs.Error("[SelectBaits] select total error,%s", err)
	}

	rows, err := DbCon.Query(sqlstr)
	if err != nil {
		logs.Error("[SelectTopAttackIps] select list error,%s", err)
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectTopAttackIps] rows.Columns() error,%s", err)
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var data map[string]interface{}

	list, count, err := util.GetTopAttackIpsMysqlJson(rows, columns, total, values, scanArgs)
	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			logs.Error("[SelectTopAttackIps] Unmarshal list error,%s", err)
		}
	}
	return data
}

func SelectTopSourceIps() map[string]interface{} {
	sqlCon, err1 := sql.Open("mysql", Dbuser+":"+Dbpassword+"@tcp("+Dbhost+":"+Dbport+")/"+Dbname+"?charset=utf8&loc=Asia%2FShanghai")
	if err1 != nil {
		logs.Error("[SelectTopSourceIps]open mysql fail %s", err1)
	}
	DbCon := sqlCon
	defer sqlCon.Close()
	sqlstr := "SELECT srchost,COUNT(*) as ipcount FROM attacklog WHERE srchost !='' GROUP BY srchost ORDER BY COUNT(*) desc"

	rows, err := DbCon.Query(sqlstr)
	if err != nil {
		logs.Error("[SelectTopSourceIps] select list error,%s", err)
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectTopSourceIps] rows.Columns() error,%s", err)
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var data map[string]interface{}

	list, count, err := util.GetTopSourceIpsMysqlJson(rows, columns, 10, values, scanArgs)
	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			logs.Error("[SelectTopSourceIps] Unmarshal list error,%s", err)
		}
	}
	return data
}

func SelectTopAreas() map[string]interface{} {
	sqlCon, err1 := sql.Open("mysql", Dbuser+":"+Dbpassword+"@tcp("+Dbhost+":"+Dbport+")/"+Dbname+"?charset=utf8&loc=Asia%2FShanghai")
	if err1 != nil {
		logs.Error("[SelectTopAreas]open mysql fail %s", err1)
	}
	DbCon := sqlCon
	var total int
	defer sqlCon.Close()
	sqltotal := "SELECT COUNT(*) FROM attacklog WHERE country !=''"
	sqlstr := "SELECT country,province,COUNT(province) as areacount FROM attacklog WHERE country !='' GROUP BY province ORDER BY areacount desc limit 0,5"

	err := DbCon.QueryRow(sqltotal).Scan(&total)
	if err != nil {
		logs.Error("[SelectHoneyInfos] select total error,%s", err)
	}
	rows, err := DbCon.Query(sqlstr)
	if err != nil {
		logs.Error("[SelectTopAreas] select list error,%s", err)
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectTopAreas] rows.Columns() error,%s", err)
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var data map[string]interface{}

	list, count, err := util.GetTopAreasMysqlJson(rows, columns, total, values, scanArgs)
	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			logs.Error("[SelectTopAreas] Unmarshal list error,%s", err)
		}
	}
	return data
}

func SelectTopAttackTypes() map[string]interface{} {
	sqlCon, err1 := sql.Open("mysql", Dbuser+":"+Dbpassword+"@tcp("+Dbhost+":"+Dbport+")/"+Dbname+"?charset=utf8&loc=Asia%2FShanghai")
	if err1 != nil {
		logs.Error("[SelectTopAttackTypes]open mysql fail %s", err1)
	}
	DbCon := sqlCon
	defer sqlCon.Close()
	sqlstr := "SELECT T1.honeypottype,COUNT(*) as typecount FROM attacklog T0 LEFT JOIN honeypotstype T1 ON T0.honeytypeid = T1.typeid WHERE T0.honeytypeid !='' GROUP BY T0.honeytypeid ORDER BY typecount desc"

	rows, err := DbCon.Query(sqlstr)
	if err != nil {
		logs.Error("[SelectTopAttackTypes] select list error,%s", err)
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectTopAttackTypes] rows.Columns() error,%s", err)
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var data map[string]interface{}
	list, count, err := util.GetTopAttackTypesMysqlJson(rows, columns, 10, values, scanArgs)
	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			logs.Error("[SelectTopAttackTypes] Unmarshal list error,%s", err)
		}
	}
	return data
}

func QueryTopologyNodes() map[string]interface{} {
	sqlCon, err1 := sql.Open("mysql", Dbuser+":"+Dbpassword+"@tcp("+Dbhost+":"+Dbport+")/"+Dbname+"?charset=utf8&loc=Asia%2FShanghai")
	if err1 != nil {
		logs.Error("[SelectTopAttackTypes]open mysql fail %s", err1)
	}
	DbCon := sqlCon
	defer sqlCon.Close()
	sqlstr := "SELECT T1.honeypottype,COUNT(*) as typecount FROM attacklog T0 LEFT JOIN honeypotstype T1 ON T0.honeytypeid = T1.typeid WHERE T0.honeytypeid !='' GROUP BY T0.honeytypeid ORDER BY typecount desc"

	rows, err := DbCon.Query(sqlstr)
	if err != nil {
		logs.Error("[SelectTopAttackTypes] select list error,%s", err)
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectTopAttackTypes] rows.Columns() error,%s", err)
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var data map[string]interface{}
	list, count, err := util.GetTopAttackTypesMysqlJson(rows, columns, 10, values, scanArgs)
	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			logs.Error("[SelectTopAttackTypes] Unmarshal list error,%s", err)
		}
	}
	return data
}

func QueryTopologyLines() map[string]interface{} {
	sqlCon, err1 := sql.Open("mysql", Dbuser+":"+Dbpassword+"@tcp("+Dbhost+":"+Dbport+")/"+Dbname+"?charset=utf8&loc=Asia%2FShanghai")
	if err1 != nil {
		logs.Error("[SelectTopAttackTypes]open mysql fail %s", err1)
	}
	DbCon := sqlCon

	defer sqlCon.Close()

	allHoneyPotSelectSql := "SELECT h.honeypotid, h.honeyip, concat(ht.honeypottype, \":\", h.honeyport) as hostname FROM honeypots h LEFT JOIN honeypotstype ht ON h.honeytypeid = ht.typeid where status = 1 AND ht.honeypottype is not NULL AND  h.honeyport is not null UNION ALL SELECT f.id, s.serverip, concat(s.servername, \":\", f.honeypotport) as hostname FROM fowards f LEFT JOIN servers s ON f.agentid = s.agentid WHERE f.status = 1"

	rows, err := DbCon.Query(allHoneyPotSelectSql)
	if err != nil {
		logs.Error("[SelectTopAttackTypes] select list error,%s", err)
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectTopAttackTypes] rows.Columns() error,%s", err)
	}
	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))

	for i := range values {
		scanArgs[i] = &values[i]
	}
	var data map[string]interface{}
	list, count, err := util.GetTopAttackTypesMysqlJson(rows, columns, 10, values, scanArgs)
	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			logs.Error("[SelectTopAttackTypes] Unmarshal list error,%s", err)
		}
	}
	return data
}
