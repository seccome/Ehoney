package honeycluster

import (
	"database/sql"
	"decept-defense/models/util"
	"decept-defense/models/util/comm"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"time"
)

var (
	Dbhost     = beego.AppConfig.String("dbhost")
	Dbport     = beego.AppConfig.String("dbport")
	Dbuser     = beego.AppConfig.String("dbuser")
	Dbpassword = beego.AppConfig.String("dbpassword")
	Dbname     = beego.AppConfig.String("dbname")
)

func CheckAdminLogin(username string, password string) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("select * from `e_admin` where uname=? and upass=?", username, password).Values(&maps)
	if err != nil {
		logs.Error("[CheckAdminLogin] select event list error,%s", err)
		return maps
	}
	return maps
}

func UpdateAdminLoginStatus(username string, password string) {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("update e_admin set status=2 where uname=? and upass=?", username, password).Values(&maps)
	if err != nil {
		logs.Error("[UpdateAdminLoginStatus] update AdminLoginStatus policy error,%s", err)
	}
}

//func CheckAdminLogin(username string, password string) bool {
//	result := false
//	sqlCon, err1 := sql.Open("mysql", Dbuser+":"+Dbpassword+"@tcp("+Dbhost+":"+Dbport+")/"+Dbname+"?charset=utf8&loc=Asia%2FShanghai")
//	if err1 != nil {
//		logs.Error("[CheckAdminLogin]open mysql fail %s", err1)
//	}
//	DbCon := sqlCon
//	defer sqlCon.Close()
//	var condition string
//	sqlstr := "select * from `e_admin` where 1=1"
//	var a []interface{}
//	if username != "" {
//		condition += " and uname=?"
//		a = append(a, username)
//	}
//	if password != "" {
//		condition += " and upass=?"
//		a = append(a, password)
//	}
//	rows, err := DbCon.Query(sqlstr+condition, a...)
//	if err != nil {
//		logs.Error("[CheckAdminLogin] select list error,%s", err)
//	}
//	if rows.Next() {
//		result = true
//	}
//	return result
//}

func CheckLoginErrNo(username string) bool {
	result := false
	sqlCon, err1 := sql.Open("mysql", Dbuser+":"+Dbpassword+"@tcp("+Dbhost+":"+Dbport+")/"+Dbname+"?charset=utf8&loc=Asia%2FShanghai")
	if err1 != nil {
		logs.Error("[CheckLoginErrNo]open mysql fail %s", err1)
	}
	DbCon := sqlCon
	defer sqlCon.Close()
	var condition string
	sqlstr := "select * from `e_admin` where errno>=5"
	var a []interface{}
	if username != "" {
		condition += " and uname=?"
		a = append(a, username)
	}
	rows, err := DbCon.Query(sqlstr+condition, a...)
	if err != nil {
		logs.Error("[CheckLoginErrNo] select list error,%s", err)
	}
	if rows.Next() {
		result = true
	}
	return result

}

func UpdateErrNo(username string) {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("update e_admin set errno=errno+1 where uname=?", username).Values(&maps)
	if err != nil {
		logs.Error("[UpdateErrNo] error,%s", err)
	}
}

func InsertApplicationCluster(servername string, serverip string, serverid string, agentid string) (map[string]interface{}, string, int) {
	msg := "成功"
	var data map[string]interface{}
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("insert into honeypotservers (servername, serverip, serverid, agentid, status ) VALUES (?,?,?,?,1) ", servername, serverip, serverid, agentid).Values(&maps)
	if err != nil {
		logs.Error("[InsertCreateApplicationCluster] insert servers error,%s", err)
		msg = "数据插入失败"
		return data, msg, comm.ErrorCode
	} else {
		msg = "成功"
	}
	return data, msg, comm.SuccessCode
}

func SelectApplicationClusters(serverip string, servername string, vpcname string, status int, pageSize int, pageNum int) map[string]interface{} {
	sqlCon, err1 := sql.Open("mysql", Dbuser+":"+Dbpassword+"@tcp("+Dbhost+":"+Dbport+")/"+Dbname+"?charset=utf8&loc=Asia%2FShanghai")
	if err1 != nil {
		logs.Error("[SelectApplicationClusters]open mysql fail %s", err1)
	}
	DbCon := sqlCon
	var total int
	defer sqlCon.Close()
	var condition string
	sqltotal := "select count(1) from `servers` where 1=1"
	sqlstr := "select * from `servers` where 1=1"
	var a []interface{}
	if serverip != "" {
		condition += " and serverip like ?"
		a = append(a, "%"+serverip+"%")
	}
	if servername != "" {
		condition += " and servername like ? "
		a = append(a, "%"+servername+"%")
	}
	if status == 2 {
		nowtime := time.Now().Unix()
		nowtimestr := util.Strval(nowtime - 300)
		condition += " and heartbeattime <= " + nowtimestr
	}
	if status == 1 {
		nowtime := time.Now().Unix()
		nowtimestr := util.Strval(nowtime - 300)
		condition += " and heartbeattime > " + nowtimestr
	}
	if vpcname != "" {
		condition += " and vpcname=?"
		a = append(a, vpcname)
	}

	sqltotal = sqltotal + condition
	err := DbCon.QueryRow(sqltotal, a...).Scan(&total)
	if err != nil {
		logs.Error("[SelectApplicationClusters] select total error,%s", err)
	}
	totalpage, pagenum, pagesize := util.JudgePage(total, pageNum, pageSize)
	offset := (pagenum - 1) * pagesize
	condition += "  limit ?,?"
	a = append(a, offset)
	a = append(a, pagesize)
	sqlstr = sqlstr + condition
	//fmt.Println("sqlstr:", sqlstr)
	//fmt.Println("a:", a)
	rows, err := DbCon.Query(sqlstr, a...)
	if err != nil {
		logs.Error("[SelectApplicationClusters] select list error,%s", err)
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectApplicationClusters] rows.Columns() error,%s", err)
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var data map[string]interface{}

	list, count, err := util.GetApplicationClustersListMysqlJson(rows, columns, total, values, scanArgs, pagenum, pagesize, totalpage)

	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			logs.Error("[SelectApplicationClusters] Unmarshal list error,%s", err)
		}
	}
	return data
}

func SelectApplicationLists(nowtimestr string) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT DISTINCT T0.agentid AS serveragentid, T0.*,COUNT(T1.agentid) AS servercount from `servers` T0 LEFT JOIN `fowards` T1 ON T0.agentid = T1.agentid WHERE T0.heartbeattime > ? GROUP BY T0.agentid ORDER BY T0.id DESC", nowtimestr).Values(&maps)
	if err != nil {
		logs.Error("[SelectApplicationLists] select event list error,%s", err)
	}
	return maps
}

func SelectBaitsById(baitid string) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT a.*,b.systype FROM `baits` a LEFT JOIN `systemtype` b ON a.baitsystype = b.sysid where a.baitid=?", baitid).Values(&maps)
	if err != nil {
		logs.Error("[SelectBaitsById] select event list error,%s", err)
	}
	return maps
}

func SelectSignsById(signid string) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT * FROM `signs` where signid=?", signid).Values(&maps)
	if err != nil {
		logs.Error("[SelectSignsById] select event list error,%s", err)
	}
	return maps
}

func SelectApplicationSignsById(taskid string) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT a.*,b.signname,b.signinfo as signfilename FROM `server_sign` a LEFT JOIN `signs` b ON a.signid = b.signid where a.taskid=?", taskid).Values(&maps)
	if err != nil {
		logs.Error("[SelectApplicationSignsById] select event list error,%s", err)

	}
	return maps
}

func SelectApplicationBaitsById(taskid string) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT a.*,b.baitname,b.baitinfo as baitfilename FROM `server_bait` a LEFT JOIN `baits` b ON a.baitid = b.baitid where a.taskid=?", taskid).Values(&maps)
	if err != nil {
		logs.Error("[SelectApplicationBaitsById] select event list error,%s", err)

	}
	return maps
}

func SelectSigns(signname string, signid string, signtype string, creator string, starttime string, endtime string, pageSize int, pageNum int) map[string]interface{} {

	sqlCon, err1 := sql.Open("mysql", Dbuser+":"+Dbpassword+"@tcp("+Dbhost+":"+Dbport+")/"+Dbname+"?charset=utf8&loc=Asia%2FShanghai")
	if err1 != nil {
		logs.Error("[SelectSigns]open mysql fail %s", err1)
	}
	DbCon := sqlCon
	var total int
	defer sqlCon.Close()
	var condition string
	sqltotal := "select count(1) FROM `signs` where 1=1"
	sqlstr := "SELECT * FROM `signs` where 1=1"
	var a []interface{}
	if signid != "" {
		condition += " and signid=?"
		a = append(a, signid)
	}
	if signname != "" {
		condition += " and signname=?"
		a = append(a, signname)
	}
	if signtype != "" {
		condition += " and signtype=?"
		a = append(a, signtype)
	}
	if creator != "" {
		condition += " and creator=?"
		a = append(a, creator)
	}
	if starttime != "" && endtime != "" {
		condition += " and createtime between ? and ?"
		a = append(a, starttime)
		a = append(a, endtime)
	}

	sqltotal = sqltotal + condition
	err := DbCon.QueryRow(sqltotal, a...).Scan(&total)
	if err != nil {
		logs.Error("[SelectSigns] select total error,%s", err)
	}
	totalpage, pagenum, pagesize := util.JudgePage(total, pageNum, pageSize)
	offset := (pagenum - 1) * pagesize
	condition += " order by id desc  limit ?,?"
	a = append(a, offset)
	a = append(a, pagesize)

	rows, err := DbCon.Query(sqlstr+condition, a...)
	if err != nil {
		logs.Error("[SelectSigns] select list error,%s", err)
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectSigns] rows.Columns() error,%s", err)
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var data map[string]interface{}

	list, count, err := util.GetSignListMysqlJson(rows, columns, total, values, scanArgs, pagenum, pagesize, totalpage)
	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			logs.Error("[SelectSigns] Unmarshal list error,%s", err)
		}
	}
	return data
}

func SelectBaits(baitid string, baittype string, baitsystype string, creator string, starttime string, endtime string, pageSize int, pageNum int) map[string]interface{} {

	sqlCon, err1 := sql.Open("mysql", Dbuser+":"+Dbpassword+"@tcp("+Dbhost+":"+Dbport+")/"+Dbname+"?charset=utf8&loc=Asia%2FShanghai")
	if err1 != nil {
		logs.Error("[SelectBaits]open mysql fail %s", err1)
	}
	DbCon := sqlCon
	var total int
	defer sqlCon.Close()
	var condition string
	sqltotal := "select count(1) FROM `baits` a LEFT JOIN `systemtype` b ON a.baitsystype = b.sysid where 1=1"
	sqlstr := "SELECT a.*,b.systype FROM `baits` a LEFT JOIN `systemtype` b ON a.baitsystype = b.sysid where 1=1"
	var a []interface{}
	if baitid != "" {
		condition += " and a.baitid=?"
		a = append(a, baitid)
	}
	if baittype != "" {
		condition += " and a.baittype=?"
		a = append(a, baittype)
	}
	if creator != "" {
		condition += " and a.creator=?"
		a = append(a, creator)
	}
	if baitsystype != "" {
		condition += " and a.baitsystype=?"
		a = append(a, baitsystype)
	}
	if starttime != "" && endtime != "" {
		condition += " and a.createtime between ? and ?"
		a = append(a, starttime)
		a = append(a, endtime)
	}

	sqltotal = sqltotal + condition
	err := DbCon.QueryRow(sqltotal, a...).Scan(&total)
	if err != nil {
		logs.Error("[SelectBaits] select total error,%s", err)
	}
	totalpage, pagenum, pagesize := util.JudgePage(total, pageNum, pageSize)
	offset := (pagenum - 1) * pagesize
	condition += " order by a.id desc limit ?,?"
	a = append(a, offset)
	a = append(a, pagesize)
	rows, err := DbCon.Query(sqlstr+condition, a...)
	if err != nil {
		logs.Error("[SelectBaits] select list error,%s", err)
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectBaits] rows.Columns() error,%s", err)
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var data map[string]interface{}

	list, count, err := util.GetBaitListMysqlJson(rows, columns, total, values, scanArgs, pagenum, pagesize, totalpage)
	//list, count, err := util.GetHoneyClustersListMysqlJson(rows, columns, total, values, scanArgs, pagenum, pagesize, totalpage)
	//fmt.Println("list:", list)
	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			logs.Error("[SelectBaits] Unmarshal list error,%s", err)
		}
	}
	return data
}

func SelectBaitsByType(baitid string, baittype string, baitsystype string, creator string, pageSize int, pageNum int) map[string]interface{} {

	sqlCon, err1 := sql.Open("mysql", Dbuser+":"+Dbpassword+"@tcp("+Dbhost+":"+Dbport+")/"+Dbname+"?charset=utf8&loc=Asia%2FShanghai")
	if err1 != nil {
		logs.Error("[SelectBaits]open mysql fail %s", err1)
	}
	DbCon := sqlCon
	var total int
	defer sqlCon.Close()
	var condition string
	sqltotal := "select count(1) FROM `baits` a LEFT JOIN `systemtype` b ON a.baitsystype = b.sysid where 1=1"
	sqlstr := "SELECT a.*,b.systype FROM `baits` a LEFT JOIN `systemtype` b ON a.baitsystype = b.sysid where 1=1"
	var a []interface{}
	if baitid != "" {
		condition += " and baitid=?"
		a = append(a, baitid)
	}
	if baittype != "" {
		condition += " and baittype=?"
		a = append(a, baittype)
	}
	if creator != "" {
		condition += " and creator=?"
		a = append(a, creator)
	}
	if baitsystype != "" {
		condition += " and baitsystype=?"
		a = append(a, baitsystype)
	}

	sqltotal = sqltotal + condition
	err := DbCon.QueryRow(sqltotal, a...).Scan(&total)
	if err != nil {
		logs.Error("[SelectBaits] select total error,%s", err)
	}
	condition += " order by id desc"
	rows, err := DbCon.Query(sqlstr+condition, a...)
	if err != nil {
		logs.Error("[SelectBaits] select list error,%s", err)
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectBaits] rows.Columns() error,%s", err)
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var data map[string]interface{}

	list, count, err := util.GetBaitListMysqlJson(rows, columns, total, values, scanArgs, 0, 0, 0)
	//list, count, err := util.GetHoneyClustersListMysqlJson(rows, columns, total, values, scanArgs, pagenum, pagesize, totalpage)
	//fmt.Println("list:", list)
	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			logs.Error("[SelectBaits] Unmarshal list error,%s", err)
		}
	}
	return data
}

func SelectSignsByType(signid string, signtype string, creator string) map[string]interface{} {

	sqlCon, err1 := sql.Open("mysql", Dbuser+":"+Dbpassword+"@tcp("+Dbhost+":"+Dbport+")/"+Dbname+"?charset=utf8&loc=Asia%2FShanghai")
	if err1 != nil {
		logs.Error("[SelectSignsByType]open mysql fail %s", err1)
	}
	DbCon := sqlCon
	var total int
	defer sqlCon.Close()
	var condition string
	sqltotal := "select count(1) FROM `signs` where 1=1"
	sqlstr := "SELECT signtype,signinfo,signid FROM `signs` where 1=1"
	var a []interface{}
	if signid != "" {
		condition += " and signid=?"
		a = append(a, signid)
	}

	if creator != "" {
		condition += " and creator=?"
		a = append(a, creator)
	}
	if signtype != "" {
		condition += " and signtype=?"
		a = append(a, signtype)
	}

	sqltotal = sqltotal + condition
	err := DbCon.QueryRow(sqltotal, a...).Scan(&total)
	if err != nil {
		logs.Error("[SelectSignsByType] select total error,%s", err)
	}
	condition += " order by id desc"
	rows, err := DbCon.Query(sqlstr+condition, a...)
	if err != nil {
		logs.Error("[SelectSignsByType] select list error,%s", err)
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectSignsByType] rows.Columns() error,%s", err)
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var data map[string]interface{}
	list, count, err := util.GetSignListByTypeMysqlJson(rows, columns, total, values, scanArgs, 0, 0, 0)
	fmt.Println("list:", list)
	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			logs.Error("[SelectSignsByType] Unmarshal list error,%s", err)
		}
	}
	return data
}

func SelectHoneyClusters(clusterip string, clustername string, clusterstatus int, pageSize int, pageNum int) map[string]interface{} {

	sqlCon, err1 := sql.Open("mysql", Dbuser+":"+Dbpassword+"@tcp("+Dbhost+":"+Dbport+")/"+Dbname+"?charset=utf8&loc=Asia%2FShanghai")
	if err1 != nil {
		logs.Error("[SelectHoneyClusters]open mysql fail %s", err1)
	}
	DbCon := sqlCon
	var total int
	defer sqlCon.Close()
	var condition string
	sqltotal := "select count(1) from `honeypotservers` where 1=1"
	sqlstr := "select * from `honeypotservers` where 1=1"
	var a []interface{}
	if clusterip != "" {
		condition += " and serverip=?"
		a = append(a, clusterip)
	}
	if clustername != "" {
		condition += " and servername=?"
		a = append(a, clustername)
	}
	if clusterstatus != 0 {
		condition += " and status=?"
		a = append(a, clusterstatus)
	}

	sqltotal = sqltotal + condition
	err := DbCon.QueryRow(sqltotal, a...).Scan(&total)
	if err != nil {
		logs.Error("[SelectHoneyClusters] select total error,%s", err)
	}
	totalpage, pagenum, pagesize := util.JudgePage(total, pageNum, pageSize)
	offset := (pagenum - 1) * pagesize
	condition += "  limit ?,?"
	a = append(a, offset)
	a = append(a, pagesize)

	rows, err := DbCon.Query(sqlstr+condition, a...)
	if err != nil {
		logs.Error("[SelectHoneyClusters] select list error,%s", err)
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectHoneyClusters] rows.Columns() error,%s", err)
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var data map[string]interface{}

	list, count, err := util.GetHoneyClustersListMysqlJson(rows, columns, total, values, scanArgs, pagenum, pagesize, totalpage)
	fmt.Println("list:", list)
	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			logs.Error("[SelectHoneyClusters] Unmarshal list error,%s", err)
		}
	}
	return data
}

func GetPodImage(podname string, pageSize int, pageNum int) map[string]interface{} {
	sqlCon, err1 := sql.Open("mysql", Dbuser+":"+Dbpassword+"@tcp("+Dbhost+":"+Dbport+")/"+Dbname+"?charset=utf8&loc=Asia%2FShanghai")
	if err1 != nil {
		logs.Error("[GetPodImage]open mysql fail %s", err1)
	}
	DbCon := sqlCon
	var total int
	defer sqlCon.Close()
	var condition string
	sqltotal := "select count(1) from `podimage` where 1=1 and imageport IS NOT NULL "
	sqlstr := "select * from `podimage` where 1=1 and imageport IS NOT NULL "
	var a []interface{}
	if podname != "" {
		condition += " and imagename=?"
		a = append(a, podname)
	}

	sqltotal = sqltotal + condition
	err := DbCon.QueryRow(sqltotal, a...).Scan(&total)
	if err != nil {
		logs.Error("[GetPodImage] select total error,%s", err)
	}
	pageNum = 1
	pageSize = 99
	totalpage, pagenum, pagesize := util.JudgePage(total, pageNum, pageSize)
	offset := (pagenum - 1) * pagesize
	condition += "  limit ?,?"
	a = append(a, offset)
	a = append(a, pagesize)

	rows, err := DbCon.Query(sqlstr+condition, a...)
	if err != nil {
		logs.Error("[GetPodImage] select list error,%s", err)
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[GetPodImage] rows.Columns() error,%s", err)
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var data map[string]interface{}
	list, count, err := util.GetPodImageListMysqlJson(rows, columns, total, values, scanArgs, pagenum, pagesize, totalpage)
	fmt.Println("list:", list)
	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			logs.Error("[GetPodImage] Unmarshal list error,%s", err)
		}
	}
	return data
}

/**
更新蜜罐状态
*/
func UpdatePod(podname string, podid string, offlinetime int64) {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("update honeypots set status=2,offlinetime=? where honeypotid=?", offlinetime, podid).Values(&maps)
	if err != nil {
		logs.Error("[UpdatePod] insert bait policy error,%s", err)
	}
}

/**
更新诱饵状态
*/
func UpdateApplicationBaits(podname string, podid string, offlinetime int64) {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("update server_bait set status=2,offlinetime=? where honeypotid=?", offlinetime, podid).Values(&maps)
	if err != nil {
		logs.Error("[UpdatePod] insert bait policy error,%s", err)
	}
}

/**
添加新增蜜罐信息
*/
func InsertPodinfo(podname string, Honeytypeid string, honeypotid string, honeyip string, honeyport int, createtime int64, creator string, agentid string, sysid string) (map[string]interface{}, string, int) {
	msg := "成功"
	var data map[string]interface{}
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("insert into honeypots (honeyname, honeytypeid, honeypotid, honeyip, honeyport, createtime, creator, agentid, sysid, status ) VALUES (?,?,?,?,?,?,?,?,?,3)", podname, Honeytypeid, honeypotid, honeyip, honeyport, createtime, creator, agentid, sysid).Values(&maps)
	if err != nil {
		logs.Error("[InsertHoneyBait] insert servers error,%s", err)
		msg = "数据插入失败"
		return data, msg, comm.ErrorCode
	} else {
		msg = "成功"
	}
	return data, msg, comm.SuccessCode
}

/**
添加新增蜜罐镜像信息
*/
func InsertPodImage(podname string, podurl string) (map[string]interface{}, string, int) {
	msg := "成功"
	var data map[string]interface{}
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("insert into podimage (imageaddress, imagename) VALUES (?,?)", podurl, podname).Values(&maps)
	if err != nil {
		logs.Error("[InsertPodImage] insert servers error,%s", err)
		msg = "数据插入失败"
		return data, msg, comm.ErrorCode
	} else {
		msg = "成功"
	}
	return data, msg, comm.SuccessCode
}

/**
添加新增蜜罐信息
*/

func FreshPodInfo1(podsourcename string, podname string, honeynamespce string, podip string, poduid string, podport int32, podsimage string, pstatus int) (map[string]interface{}, string, int) {
	msg := "成功"
	var data map[string]interface{}
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("insert into honeypots (honeyname, podname, honeynamespce, honeypotid, honeyip, honeyport, honeyimage, status ) VALUES (?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE podname=?, honeynamespce=?, honeypotid=?, honeyip=?, honeyport=?, honeyimage=?, status=? ", podsourcename, podname, honeynamespce, poduid, podip, podport, podsimage, pstatus, podname, honeynamespce, poduid, podip, podport, podsimage, pstatus).Values(&maps)
	if err != nil {
		logs.Error("[FreshPodInfo] insert servers error,%s", err)
		msg = "数据插入失败"
		return data, msg, comm.ErrorCode
	} else {
		msg = "成功"
	}
	return data, msg, comm.SuccessCode
}

func FreshPodInfo(podsourcename string, podname string, honeynamespce string, podip string, poduid string, podsimage string, pstatus int) (map[string]interface{}, string, int) {
	msg := "成功"
	var data map[string]interface{}
	o := orm.NewOrm()
	var maps []orm.Params
	//_, err := o.Raw("insert into honeypots (honeyname, podname, honeynamespce, honeypotid, honeyip, honeyport, honeyimage, status ) VALUES (?,?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE podname=?, honeynamespce=?, honeypotid=?, honeyip=?, honeyport=?, honeyimage=?, status=? ", podsourcename, podname, honeynamespce, poduid, podip, podport, podsimage, pstatus, podname, honeynamespce, poduid, podip, podport, podsimage, pstatus).Values(&maps)
	_, err := o.Raw("UPDATE honeypots set honeyname=?,podname=?,honeynamespce=?,honeypotid=?, honeyip=?, honeyimage=?, status=? WHERE honeyname=? and isNull(offlinetime) ", podsourcename, podname, honeynamespce, poduid, podip, podsimage, pstatus, podsourcename).Values(&maps)
	if err != nil {
		logs.Error("[FreshPodInfo] insert servers error,%s", err)
		msg = "数据插入失败"
		return data, msg, comm.ErrorCode
	} else {
		msg = "成功"
	}
	return data, msg, comm.SuccessCode
}

/**
插入诱饵
*/
func InsertBait(baitname string, createtime int64, creator string, baitid string, baitsystype string, baitinfo string, filemd5 string) (map[string]interface{}, string, int) {
	msg := "成功"
	var data map[string]interface{}
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("insert into baits (baittype ,baitname, createtime, creator, baitid, baitsystype, baitinfo , md5) VALUES ('file',?,?,?,?,?,?,?)", baitname, createtime, creator, baitid, baitsystype, baitinfo, filemd5).Values(&maps)
	if err != nil {
		logs.Error("[InsertBait] insert servers error,%s", err)
		msg = "数据插入失败"
		return data, msg, comm.ErrorCode
	} else {
		msg = "成功"
	}
	return data, msg, comm.SuccessCode
}

/**
插入密签
*/
func InsertSign(signtype string, signname string, createtime int64, creator string, signid string, signsystype string, signinfo string) (map[string]interface{}, string, int) {
	msg := "成功"
	var data map[string]interface{}
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("insert into signs (signtype ,signname,createtime, creator, signid, signsystype, signinfo) VALUES (?,?,?,?,?,?,?)", signtype, signname, createtime, creator, signid, signsystype, signinfo).Values(&maps)
	if err != nil {
		logs.Error("[InsertSign] insert servers error,%s", err)
		msg = "数据插入失败"
		return data, msg, comm.ErrorCode
	} else {
		msg = "成功"
	}
	return data, msg, comm.SuccessCode
}

/**
删除诱饵
*/
func DeleteBaitById(baitid string) (map[string]interface{}, string, int) {
	msg := "成功"
	o := orm.NewOrm()
	var maps []orm.Params
	var data map[string]interface{}
	_, err := o.Raw("DELETE FROM baits where baitid=?", baitid).Values(&maps)
	if err != nil {
		logs.Error("[DeleteBaitById] delete config error,%s", err)
		msg = "数据删除失败"
		return data, msg, comm.ErrorCode
	}
	return data, msg, comm.SuccessCode
}

/**
删除密签
*/
func DeleteSignById(signid string) (map[string]interface{}, string, int) {
	msg := "成功"
	o := orm.NewOrm()
	var maps []orm.Params
	var data map[string]interface{}
	_, err := o.Raw("DELETE FROM signs where signid=?", signid).Values(&maps)
	if err != nil {
		logs.Error("[DeleteSignById] delete config error,%s", err)
		msg = "数据删除失败"
		return data, msg, comm.ErrorCode
	}
	return data, msg, comm.SuccessCode
}

func InsertHoneyBait(createbaitstatus int, taskid string, baitid string, baitinfo string, honeypotid string, createtime int64, creator string) (map[string]interface{}, string, int) {
	msg := "成功"
	var data map[string]interface{}
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("insert into honeypot_bait (taskid, baitid, baitinfo, honeypotid,createtime, creator,status ) VALUES (?,?,?,?,?,?,?)", taskid, baitid, baitinfo, honeypotid, createtime, creator, createbaitstatus).Values(&maps)
	if err != nil {
		logs.Error("[InsertHoneyBait] insert servers error,%s", err)
		msg = "数据插入失败"
		return data, msg, comm.ErrorCode
	} else {
		msg = "成功"
	}
	return data, msg, comm.SuccessCode
}

func InsertHonetSign(createsignstatus int, taskid string, signid string, signinfo string, honeypotid string, createtime int64, creator string, tracecode string) (map[string]interface{}, string, int) {
	msg := "成功"
	var data map[string]interface{}
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("insert into honeypot_sign (taskid, signid, signinfo, honeypotid,createtime, creator, tracecode,status) VALUES (?,?,?,?,?,?,?,?)", taskid, signid, signinfo, honeypotid, createtime, creator, tracecode, createsignstatus).Values(&maps)
	if err != nil {
		logs.Error("[InsertHonetSign] insert servers error,%s", err)
		msg = "数据插入失败"
		return data, msg, comm.ErrorCode
	} else {
		msg = "成功"
	}
	return data, msg, comm.SuccessCode
}

func UpdateHoneyBait(deletebaitstatus int, taskid string, offtime int64) {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("update honeypot_bait set status=?, offlinetime=? where taskid=?", deletebaitstatus, offtime, taskid).Values(&maps)
	if err != nil {
		logs.Error("[UpdateHoneyBait] update bait policy error,%s", err)
	}
}

func UpdateHoneySign(deletesignstatus int, taskid string, offtime int64) {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("update honeypot_sign set status=?, offlinetime=? where taskid=?", deletesignstatus, offtime, taskid).Values(&maps)
	if err != nil {
		logs.Error("[UpdateHoneySign] update sign policy error,%s", err)
	}
}

func UpdateHoneyImageById(id int, imagetype string, imageos string, imageport int) error {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("update podimage set imagetype=?, imageos=?, imageport=? where id=?", imagetype, imageos, imageport, id).Values(&maps)
	if err != nil {
		logs.Error("[UpdateHoneyImageById] update image policy error,%s", err)
		return err
	} else {
		return nil
	}
}

func SelectHoneyInfoById(honeypotid string) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT * FROM `honeypots`  where honeypotid=?", honeypotid).Values(&maps)
	if err != nil {
		logs.Error("[SelectHoneyInfoById] select event list error,%s", err)

	}
	return maps
}

func SelectHoneyInfoByIp(honeyip string) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT * FROM `honeypots`  where honeyip=?", honeyip).Values(&maps)
	if err != nil {
		logs.Error("[SelectHoneyInfoByIp] select event list error,%s", err)

	}
	return maps
}

func SelectHoneyInfoByHoneyId(honeypotid string) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT * FROM `honeypots`  where honeyip=?", honeypotid).Values(&maps)
	if err != nil {
		logs.Error("[SelectHoneyInfoByHoneyId] select event list error,%s", err)

	}
	return maps
}

func SelectHoneyPotByIp(honeyip string) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT T0.*,T1.servername,T1.serverip FROM `honeypots` T0 LEFT JOIN honeypotservers T1 ON T0.agentid = T1.agentid where T0.honeyip=? AND T0.`status` =1 ORDER BY id DESC LIMIT 0,1", honeyip).Values(&maps)
	if err != nil {
		logs.Error("[SelectHoneyPotByIp] select event list error,%s", err)

	}
	return maps
}

func SelectHoneyBaitsById(taskid string) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT a.*,c.honeyname,c.podname,b.baitname,b.baitinfo as baitfilename FROM `honeypot_bait` a LEFT JOIN `baits` b ON a.baitid = b.baitid LEFT JOIN honeypots c ON a.honeypotid=c.honeypotid where a.taskid=?", taskid).Values(&maps)
	if err != nil {
		logs.Error("[SelectHoneyBaitsById] select event list error,%s", err)

	}
	return maps
}

func SelectHoneypotServer() []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT * FROM `honeypotservers`").Values(&maps)
	if err != nil {
		logs.Error("[SelectHoneypotServer] select event list error,%s", err)

	}
	return maps
}

func SelectHoneypotTypeByID(typeID string) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT softpath FROM `honeypotstype` where typeid=?", typeID).Values(&maps)
	if err != nil {
		logs.Error("[selectHoneySignById] select event list error,%s", err)

	}
	return maps
}

func SelectHoneySignById(taskid string) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT a.*,c.honeyname,c.podname,b.signname,b.signinfo as signfilename FROM `honeypot_sign` a LEFT JOIN `signs` b ON a.signid = b.signid LEFT JOIN honeypots c ON a.honeypotid=c.honeypotid where a.taskid=?", taskid).Values(&maps)
	if err != nil {
		logs.Error("[selectHoneySignById] select event list error,%s", err)

	}
	return maps
}

func SelectHoneyInfos(serverid string, sysid string, honeytypeid string, honeyip string, honeyname string, starttime string, endtime string, status int, creator string, pageSize int, pageNum int) map[string]interface{} {
	sqlCon, err1 := sql.Open("mysql", Dbuser+":"+Dbpassword+"@tcp("+Dbhost+":"+Dbport+")/"+Dbname+"?charset=utf8&loc=Asia%2FShanghai")
	if err1 != nil {
		logs.Error("[SelectHoneyInfos]open mysql fail %s", err1)
	}
	DbCon := sqlCon
	var total int
	defer sqlCon.Close()
	var condition string
	sqltotal := "select COUNT(1) from `honeypots` a LEFT JOIN `honeypotservers` b  ON a.agentid = b.agentid LEFT JOIN honeypotstype c on a.honeytypeid=c.typeid LEFT JOIN systemtype d ON a.sysid=d.sysid where 1=1"
	sqlstr := "select a.*,b.serverip,b.servername,c.honeypottype,d.systype from `honeypots` a LEFT JOIN `honeypotservers` b  ON a.agentid = b.agentid LEFT JOIN honeypotstype c on a.honeytypeid=c.typeid LEFT JOIN systemtype d ON a.sysid=d.sysid where 1=1"
	var a []interface{}
	if sysid != "" {
		condition += " and a.sysid=?"
		a = append(a, sysid)
	}
	if serverid != "" {
		condition += " and a.serverid=?"
		a = append(a, serverid)
	}
	if honeytypeid != "" {
		condition += " and a.honeytypeid=?"
		a = append(a, honeytypeid)
	}
	if honeyip != "" {
		condition += " and a.honeyip=?"
		a = append(a, honeyip)
	}
	if honeyname != "" {
		condition += " and a.honeyname like ? "
		a = append(a, "%"+honeyname+"%")
	}
	if starttime != "" && endtime != "" {
		condition += " and a.createtime between ? and ?"
		a = append(a, starttime)
		a = append(a, endtime)
	}
	if creator != "" {
		condition += " and a.creator=?"
		a = append(a, creator)
	}
	if status != 0 {
		condition += " and a.status=?"
		a = append(a, status)
	}

	sqltotal = sqltotal + condition
	err := DbCon.QueryRow(sqltotal, a...).Scan(&total)
	if err != nil {
		logs.Error("[SelectHoneyInfos] select total error,%s", err)
	}
	totalpage, pagenum, pagesize := util.JudgePage(total, pageNum, pageSize)
	offset := (pagenum - 1) * pagesize
	condition += " ORDER BY a.createtime DESC"
	condition += " limit ?,?"
	a = append(a, offset)
	a = append(a, pagesize)
	rows, err := DbCon.Query(sqlstr+condition, a...)
	if err != nil {
		logs.Error("[SelectHoneyClusters] select list error,%s", err)
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectHoneyClusters] rows.Columns() error,%s", err)
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var data map[string]interface{}

	list, count, err := util.GetHoneyInfosListMysqlJson(rows, columns, total, values, scanArgs, pagenum, pagesize, totalpage)

	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			logs.Error("[SelectHoneyInfos] Unmarshal list error,%s", err)
		}
	}
	return data
}

func SelectHoneyList(serverid string, honeytypeid string, honeyip string, honeyname string, starttime string, endtime string, status int, creator string, pageSize int, pageNum int) map[string]interface{} {
	sqlCon, err1 := sql.Open("mysql", Dbuser+":"+Dbpassword+"@tcp("+Dbhost+":"+Dbport+")/"+Dbname+"?charset=utf8&loc=Asia%2FShanghai")
	if err1 != nil {
		logs.Error("[SelectHoneyList]open mysql fail %s", err1)
	}
	DbCon := sqlCon
	var total int
	defer sqlCon.Close()
	var condition string
	sqltotal := "select COUNT(1) from `honeypots` where status=1"
	//sqltotal := "select COUNT(1) from `honeypots` a LEFT JOIN `honeypotservers` b  ON a.serverid = b.serverid LEFT JOIN honeypotstype c on a.honeytypeid=c.typeid LEFT JOIN systemtype d ON a.sysid=d.sysid where 1=1"
	sqlstr := "select honeyname,honeypotid from `honeypots` where status=1"
	var a []interface{}
	if serverid != "" {
		condition += " and serverid=?"
		a = append(a, serverid)
	}
	if honeytypeid != "" {
		condition += " and honeytypeid=?"
		a = append(a, honeytypeid)
	}
	if honeyip != "" {
		condition += " and honeyip=?"
		a = append(a, honeyip)
	}
	if honeyname != "" {
		condition += " and honeyname=?"
		a = append(a, honeyname)
	}

	sqltotal = sqltotal + condition
	err := DbCon.QueryRow(sqltotal, a...).Scan(&total)
	if err != nil {
		logs.Error("[SelectHoneyList] select total error,%s", err)
	}
	rows, err := DbCon.Query(sqlstr+condition, a...)
	if err != nil {
		logs.Error("[SelectHoneyList] select list error,%s", err)
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectHoneyList] rows.Columns() error,%s", err)
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var data map[string]interface{}

	list, count, err := util.GetHoneyListMysqlJson(rows, columns, total, values, scanArgs, 0, 0, 0)

	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			logs.Error("[SelectHoneyInfos] Unmarshal list error,%s", err)
		}
	}
	return data
}

/**
透明转发蜜罐列表选择
*/
func SelectHoneyForwardListForTrans() []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("select T1.honeyname,T1.honeypotid,T1.honeyip,T1.honeyport from `honeyfowards` T0 left join `honeypots` T1 ON T0.honeypotid = T1.honeypotid WHERE T1.`status`=1 GROUP BY T1.honeypotid ORDER BY T0.id DESC").Values(&maps)
	if err != nil {
		logs.Error("[SelectHoneyForwardListForTrans] select event list error,%s", err)
	}
	return maps
}

/**
协议转发蜜罐列表选择
*/
func SelectHoneyListForTrans(honeytypeid string) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT honeyname,honeypotid,honeyip,honeyport FROM honeypots WHERE `status`=1 AND honeytypeid=? ORDER BY id DESC", honeytypeid).Values(&maps)
	if err != nil {
		logs.Error("[SelectHoneyListForTrans] select event list error,%s", err)
	}
	return maps
}

func SelectHoneyTransPortsListForTrans(honeypotid string) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT T0.forwardport from `honeyfowards` T0 left join `honeypots` T1 ON T0.honeypotid = T1.honeypotid WHERE 1=1 AND T0.`status` = 1 and T0.honeypotid=?", honeypotid).Values(&maps)
	if err != nil {
		logs.Error("[SelectHoneyTransPortsListForTrans] select event list error,%s", err)
	}
	return maps
}

func SelectHoneyForwardList(serverid string, honeytypeid string, honeyip string, honeyname string, starttime string, endtime string, status int, creator string, pageSize int, pageNum int) map[string]interface{} {
	sqlCon, err1 := sql.Open("mysql", Dbuser+":"+Dbpassword+"@tcp("+Dbhost+":"+Dbport+")/"+Dbname+"?charset=utf8&loc=Asia%2FShanghai")
	if err1 != nil {
		logs.Error("[SelectHoneyList]open mysql fail %s", err1)
	}
	DbCon := sqlCon
	var total int
	defer sqlCon.Close()
	var condition string
	sqltotal := "select COUNT(distinct T1.honeyname) AS honeysum from `honeyfowards` T0 LEFT JOIN honeypots T1 ON T0.honeypotid = T1.honeypotid where T0.status=1 and T1.`status`=1 and ISNULL(T0.offlinetime) "
	sqlstr := "SELECT distinct T1.honeyname,T0.honeypotid from `honeyfowards` T0 LEFT JOIN honeypots T1 ON T0.honeypotid = T1.honeypotid where T0.status=1 and T1.`status`=1 and ISNULL(T0.offlinetime)"
	var a []interface{}
	if serverid != "" {
		condition += " and T1.serverid=?"
		a = append(a, serverid)
	}
	if honeytypeid != "" {
		condition += " and T1.honeytypeid=?"
		a = append(a, honeytypeid)
	}
	if honeyip != "" {
		condition += " and T1.honeyip=?"
		a = append(a, honeyip)
	}
	if honeyname != "" {
		condition += " and T1.honeyname=?"
		a = append(a, honeyname)
	}
	condition = condition + " ORDER BY T1.honeyname"

	sqltotal = sqltotal + condition
	err := DbCon.QueryRow(sqltotal, a...).Scan(&total)
	if err != nil {
		logs.Error("[SelectHoneyList] select total error,%s", err)
	}
	rows, err := DbCon.Query(sqlstr+condition, a...)
	if err != nil {
		logs.Error("[SelectHoneyList] select list error,%s", err)
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectHoneyList] rows.Columns() error,%s", err)
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var data map[string]interface{}

	list, count, err := util.GetHoneyListMysqlJson(rows, columns, total, values, scanArgs, 0, 0, 0)

	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			logs.Error("[SelectHoneyInfos] Unmarshal list error,%s", err)
		}
	}
	return data
}

func SelectHoneyBaits(honeypotid string, baittype string, creator string, status int, starttime string, endtime string, pageSize int, pageNum int) map[string]interface{} {
	sqlCon, err1 := sql.Open("mysql", Dbuser+":"+Dbpassword+"@tcp("+Dbhost+":"+Dbport+")/"+Dbname+"?charset=utf8&loc=Asia%2FShanghai")
	if err1 != nil {
		logs.Error("[SelectHoneyInfos]open mysql fail %s", err1)
	}
	DbCon := sqlCon
	var total int
	defer sqlCon.Close()
	var condition string
	sqltotal := "select count(1) from `honeypot_bait` a LEFT join `baits` b ON a.baitid = b.baitid  where 1=1"
	sqlstr := "SELECT a.*,b.baitname,b.baittype from `honeypot_bait` a LEFT join `baits` b ON a.baitid = b.baitid where 1=1"
	var a []interface{}
	if honeypotid != "" {
		condition += " and a.honeypotid=?"
		a = append(a, honeypotid)
	}
	if baittype != "" {
		condition += " and b.baittype=?"
		a = append(a, baittype)
	}
	if status != 0 {
		condition += " and a.status=?"
		a = append(a, status)
	}
	if creator != "" {
		condition += " and a.creator=?"
		a = append(a, creator)
	}
	if starttime != "" && endtime != "" {
		condition += " and a.createtime between ? and ?"
		a = append(a, starttime)
		a = append(a, endtime)
	}
	sqltotal = sqltotal + condition
	err := DbCon.QueryRow(sqltotal, a...).Scan(&total)
	if err != nil {
		logs.Error("[SelectHoneyBaits] select total error,%s", err)
	}
	totalpage, pagenum, pagesize := util.JudgePage(total, pageNum, pageSize)
	offset := (pagenum - 1) * pagesize
	condition += "  limit ?,?"
	a = append(a, offset)
	a = append(a, pagesize)
	rows, err := DbCon.Query(sqlstr+condition, a...)
	if err != nil {
		logs.Error("[SelectHoneyBaits] select list error,%s", err)
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectHoneyBaits] rows.Columns() error,%s", err)
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var data map[string]interface{}

	//list, count, err := util.GetHoneyClustersListMysqlJson(rows, columns, total, values, scanArgs ,pagenum, pagesize, totalpage)
	list, count, err := util.GetHoneyBaitsListMysqlJson(rows, columns, total, values, scanArgs, pagenum, pagesize, totalpage)
	fmt.Println("list:", list)
	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			logs.Error("[SelectHoneyBaits] Unmarshal list error,%s", err)
		}
	}
	return data
}

func SelectHoneySigns(honeypotid string, signtype string, creator string, status int, starttime string, endtime string, pageSize int, pageNum int) map[string]interface{} {
	sqlCon, err1 := sql.Open("mysql", Dbuser+":"+Dbpassword+"@tcp("+Dbhost+":"+Dbport+")/"+Dbname+"?charset=utf8&loc=Asia%2FShanghai")
	if err1 != nil {
		logs.Error("[SelectHoneyInfos]open mysql fail %s", err1)
	}
	DbCon := sqlCon
	var total int
	defer sqlCon.Close()
	var condition string
	sqltotal := "select count(1) from `honeypot_sign` a LEFT join `signs` b ON a.signid = b.signid where 1=1"
	sqlstr := "SELECT a.*,b.signname,b.signtype,b.signinfo AS signfilename from `honeypot_sign` a LEFT join `signs` b ON a.signid = b.signid where 1=1"
	var a []interface{}
	if honeypotid != "" {
		condition += " and a.honeypotid=?"
		a = append(a, honeypotid)
	}
	if signtype != "" {
		condition += " and b.signtype=?"
		a = append(a, signtype)
	}
	if creator != "" {
		condition += " and a.creator=?"
		a = append(a, creator)
	}
	if status != 0 {
		condition += " and a.status=?"
		a = append(a, status)
	}
	if starttime != "" && endtime != "" {
		condition += " and a.createtime between ? and ?"
		a = append(a, starttime)
		a = append(a, endtime)
	}
	sqltotal = sqltotal + condition
	err := DbCon.QueryRow(sqltotal, a...).Scan(&total)
	if err != nil {
		logs.Error("[SelectHoneySigns] select total error,%s", err)
	}
	totalpage, pagenum, pagesize := util.JudgePage(total, pageNum, pageSize)
	offset := (pagenum - 1) * pagesize
	condition += "  limit ?,?"
	a = append(a, offset)
	a = append(a, pagesize)
	rows, err := DbCon.Query(sqlstr+condition, a...)
	if err != nil {
		logs.Error("[SelectHoneySigns] select list error,%s", err)
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectHoneySigns] rows.Columns() error,%s", err)
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var data map[string]interface{}

	list, count, err := util.GetHoneySignsListMysqlJson(rows, columns, total, values, scanArgs, pagenum, pagesize, totalpage)
	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			logs.Error("[SelectHoneySigns] Unmarshal list error,%s", err)
		}
	}
	return data
}

func ApplicationSelectSignMsg(tracecode string, listmap string, starttime string, endtime string, pageNum int, pageSize int) map[string]interface{} {
	sqlCon, err1 := sql.Open("mysql", Dbuser+":"+Dbpassword+"@tcp("+Dbhost+":"+Dbport+")/"+Dbname+"?charset=utf8&loc=Asia%2FShanghai")
	if err1 != nil {
		logs.Error("[ApplicationSelectSignMsg]open mysql fail %s", err1)
	}
	DbCon := sqlCon
	var total int
	defer sqlCon.Close()
	var condition string
	sqltotal := "select count(*) from signtracemsg T0 LEFT JOIN server_sign T1 ON T0.tracecode=T1.tracecode LEFT JOIN signs T2 ON T2.signid=T1.signid where 1=1"
	sqlstr := "select T0.*,T1.signinfo,T2.signinfo as signfilename from signtracemsg T0 LEFT JOIN server_sign T1 ON T0.tracecode=T1.tracecode LEFT JOIN signs T2 ON T2.signid=T1.signid where 1=1"
	var a []interface{}
	if tracecode != "" {
		condition += " and T0.tracecode=?"
		a = append(a, tracecode)
	}
	totalpage := 0
	pagenum := 0
	pagesize := 0
	if listmap != "true" {
		if starttime != "" && endtime != "" {
			condition += " and opentime between ? and ?"
			a = append(a, starttime)
			a = append(a, endtime)
		}
		sqltotal = sqltotal + condition
		err := DbCon.QueryRow(sqltotal, a...).Scan(&total)
		if err != nil {
			logs.Error("[ApplicationSelectSignMsg] select total error,%s", err)
		}
		totalpage, pagenum, pagesize = util.JudgePage(total, pageNum, pageSize)
		offset := (pagenum - 1) * pagesize
		condition += "  limit ?,?"
		a = append(a, offset)
		a = append(a, pagesize)
	}
	rows, err := DbCon.Query(sqlstr+condition, a...)
	if err != nil {
		logs.Error("[ApplicationSelectSignMsg] select list error,%s", err)
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[ApplicationSelectSignMsg] rows.Columns() error,%s", err)
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var data map[string]interface{}

	list, count, err := util.GetHoneySignsMsgMysqlJson(rows, columns, total, values, scanArgs, pagenum, pagesize, totalpage)
	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			logs.Error("[HoneySelectSignMsg] Unmarshal list error,%s", err)
		}
	}
	return data
}

func HoneySelectSignMsg(tracecode string, listmap string, starttime string, endtime string, pageSize int, pageNum int) map[string]interface{} {
	sqlCon, err1 := sql.Open("mysql", Dbuser+":"+Dbpassword+"@tcp("+Dbhost+":"+Dbport+")/"+Dbname+"?charset=utf8&loc=Asia%2FShanghai")
	if err1 != nil {
		logs.Error("[SelectHoneyInfos]open mysql fail %s", err1)
	}
	DbCon := sqlCon
	var total int
	defer sqlCon.Close()
	var condition string
	sqltotal := "select count(*) from signtracemsg T0 LEFT JOIN honeypot_sign T1 ON T0.tracecode=T1.tracecode LEFT JOIN signs T2 ON T1.signid=T2.signid where 1=1"
	sqlstr := "select T0.*,T1.signinfo,T2.signinfo as signfilename from signtracemsg T0 LEFT JOIN honeypot_sign T1 ON T0.tracecode=T1.tracecode LEFT JOIN signs T2 ON T1.signid=T2.signid where 1=1"
	var a []interface{}
	if tracecode != "" {
		condition += " and T0.tracecode=?"
		a = append(a, tracecode)
	}
	totalpage := 0
	pagenum := 0
	pagesize := 0
	if listmap != "true" {
		if starttime != "" && endtime != "" {
			condition += " and opentime between ? and ?"
			a = append(a, starttime)
			a = append(a, endtime)
		}
		sqltotal = sqltotal + condition
		err := DbCon.QueryRow(sqltotal, a...).Scan(&total)
		if err != nil {
			logs.Error("[HoneySelectSignMsg] select total error,%s", err)
		}
		totalpage, pagenum, pagesize = util.JudgePage(total, pageNum, pageSize)
		offset := (pagenum - 1) * pagesize
		condition += "  limit ?,?"
		a = append(a, offset)
		a = append(a, pagesize)
	}
	rows, err := DbCon.Query(sqlstr+condition, a...)
	if err != nil {
		logs.Error("[HoneySelectSignMsg] select list error,%s", err)
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[HoneySelectSignMsg] rows.Columns() error,%s", err)
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var data map[string]interface{}

	list, count, err := util.GetHoneySignsMsgMysqlJson(rows, columns, total, values, scanArgs, pagenum, pagesize, totalpage)
	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			logs.Error("[HoneySelectSignMsg] Unmarshal list error,%s", err)
		}
	}
	return data
}

//func SelectApplicationsBaits(agengtid string, baittype string, creator string, starttime string, endtime string, pageSize int, pageNum int) map[string]interface{} {
//	sqlCon, err1 := sql.Open("mysql", Dbuser+":"+Dbpassword+"@tcp("+Dbhost+":"+Dbport+")/"+Dbname+"?charset=utf8&loc=Asia%2FShanghai")
//	if err1 != nil {
//		logs.Error("[SelectApplicationsBaits]open mysql fail %s", err1)
//	}
//	DbCon := sqlCon
//	var total int
//	defer sqlCon.Close()
//	var condition string
//	sqltotal := "select count(1) from `server_bait` a LEFT join `baits` b ON a.baitid = b.baitid  where 1=1"
//	sqlstr := "select a.*,b.baitname,b.baittype from `server_bait` a LEFT join `baits` b ON a.baitid = b.baitid where 1=1"
//	var a []interface{}
//	if agengtid != "" {
//		condition += " and a.agengtid=?"
//		a = append(a, agengtid)
//	}
//	if baittype != "" {
//		condition += " and b.baittype=?"
//		a = append(a, baittype)
//	}
//	if creator != "" {
//		condition += " and a.creator=?"
//		a = append(a, creator)
//	}
//	if starttime != "" && endtime != "" {
//		condition += " and a.createtime between ? and ?"
//		a = append(a, starttime)
//		a = append(a, endtime)
//	}
//	sqltotal = sqltotal + condition
//	fmt.Println("sqltotal:", sqltotal)
//	err := DbCon.QueryRow(sqltotal, a...).Scan(&total)
//	if err != nil {
//		logs.Error("[SelectApplicationsBaits] select total error,%s", err)
//	}
//	totalpage, pagenum, pagesize := util.JudgePage(total, pageNum, pageSize)
//	fmt.Println("sqlstr:", sqlstr+condition)
//	rows, err := DbCon.Query(sqlstr+condition, a...)
//	if err != nil {
//		logs.Error("[SelectApplicationsBaits] select list error,%s", err)
//	}
//	columns, err := rows.Columns()
//	if err != nil {
//		logs.Error("[SelectApplicationsBaits] rows.Columns() error,%s", err)
//	}
//	values := make([]sql.RawBytes, len(columns))
//	scanArgs := make([]interface{}, len(values))
//	for i := range values {
//		scanArgs[i] = &values[i]
//	}
//	var data map[string]interface{}
//
//	//list, count, err := util.GetApplicationsClustersListMysqlJson(rows, columns, total, values, scanArgs, pagenum, pagesize, totalpage)
//	list, count, err := util.GetApplicationsClustersListMysqlJson(rows, columns, total, values, scanArgs, pagenum, pagesize, totalpage)
//
//	if count > 0 {
//		err = json.Unmarshal([]byte(list), &data)
//		if err != nil {
//			logs.Error("[SelectApplicationsBaits] Unmarshal list error,%s", err)
//		}
//	}
//	return data
//}

func InsertApplication(ecsname string, ecsip string, ecsid string, status int, vpcname string, vpsuser string, agentid string) (map[string]interface{}, string, int) {
	msg := "成功"
	var data map[string]interface{}
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("insert into servers (servername,serverip,serverid,status,agentid,vpcname,vpsowner) VALUES (?,?,?,?,?,?,?) ON DUPLICATE KEY UPDATE serverip=?,serverid=?,agentid=?, status=?", ecsname, ecsip, ecsid, status, agentid, vpcname, vpsuser, ecsip, ecsid, agentid, status).Values(&maps)
	if err != nil {
		logs.Error("[InsertApplication] insert servers error,%s", err)
		msg = "数据插入失败"
		return data, msg, comm.ErrorCode
	}
	return data, msg, comm.SuccessCode
}

func ServerHeartBeatAct(agentid string, status int, ips string, servername string, timenow int64) (map[string]interface{}, string, int) {
	msg := "成功"
	var data map[string]interface{}
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("insert into servers (servername,serverip,status,agentid,regtime,heartbeattime) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE serverip=?, status=?, heartbeattime=?", servername, ips, status, agentid, timenow, timenow, ips, status, timenow).Values(&maps)
	if err != nil {
		logs.Error("[ServerHeartBeatAct]  ServerHeartBeatAct error,%s", err)
		msg = "数据插入更新失败"
		return data, msg, comm.ErrorCode
	}
	return data, msg, comm.SuccessCode
}

func HoneyServerHeartBeatAct(agentid string, status int, ips string, servername string, timenow int64) (map[string]interface{}, string, int) {
	msg := "成功"
	var data map[string]interface{}
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("insert into honeypotservers (servername,serverip,status,agentid,regtime,heartbeattime) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE serverip=?, status=?, heartbeattime=?", servername, ips, status, agentid, timenow, timenow, ips, status, timenow).Values(&maps)
	if err != nil {
		logs.Error("[HoneyServerHeartBeatAct]  ServerHeartBeatAct error,%s", err)
		msg = "数据插入更新失败"
		return data, msg, comm.ErrorCode
	}
	return data, msg, comm.SuccessCode
}

func InsertHarborInfo(harborid string, harborurl string, userName string, password string, projectName string, createtime int64) (map[string]interface{}, string, int) {
	msg := "成功"
	var data map[string]interface{}
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("insert into harborinfo (harborid, harborurl,username,password,projectname,createtime) VALUES (?,?,?,?,?,?) ON DUPLICATE KEY UPDATE createtime=?", harborid, harborurl, userName, password, projectName, createtime, createtime).Values(&maps)
	if err != nil {
		logs.Error("[InsertHarborInfo]  InsertHarborInfo error,%s", err)
		msg = "请确认是否已经存在harbor唯一记录"
		return data, msg, comm.ErrorCode
	}
	return data, msg, comm.SuccessCode
}

func InsertTraceInfo(traceid string, tracehost string, createtime int64) (map[string]interface{}, string, int) {
	msg := "成功"
	var data map[string]interface{}
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("INSERT INTO conf_traceinfo (traceid, tracehost,createtime) VALUES (?,?,?) ON DUPLICATE KEY UPDATE tracehost=?,createtime=?", traceid, tracehost, createtime, tracehost, createtime).Values(&maps)
	if err != nil {
		logs.Error("[InsertTraceInfo] error,%s", err)
		msg = "更新成功"
		return data, msg, comm.ErrorCode
	}
	return data, msg, comm.SuccessCode
}

func SelectTraceInfo() []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("select * from `conf_traceinfo` limit 0,1").Values(&maps)
	if err != nil {
		logs.Error("[SelectTraceInfo] select list error,%s", err)
	}
	return maps
}

func InsertRedisInfo(redisid string, redisurl string, redisport string, password string, createtime int64) (map[string]interface{}, string, int) {
	msg := "成功"
	var data map[string]interface{}
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("insert into conf_redisinfo (redisid, redisurl,redisport,password,createtime) VALUES (?,?,?,?,?)", redisid, redisurl, redisport, password, createtime).Values(&maps)
	if err != nil {
		logs.Error("[InsertRedisInfo] error,%s", err)
		msg = "请确认是否已经存在redis唯一记录"
		return data, msg, comm.ErrorCode
	}
	return data, msg, comm.SuccessCode
}

func SelectRedisInfo() []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("select * from `conf_redisinfo` limit 0,1").Values(&maps)
	if err != nil {
		logs.Error("[SelectRedisInfo] select list error,%s", err)
	}
	return maps
}

func DeleteRedisInfo(redisid string) (map[string]interface{}, string, int) {
	msg := "成功"
	o := orm.NewOrm()
	var maps []orm.Params
	var data map[string]interface{}
	_, err := o.Raw("DELETE FROM conf_redisinfo where redisid=?", redisid).Values(&maps)
	if err != nil {
		logs.Error("[DeleteRedisInfo] delete redisinfo error,%s", err)
		msg = "redis记录删除失败"
		return data, msg, comm.ErrorCode
	}
	return data, msg, comm.SuccessCode
}

func SelectHarborInfo() []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("select * from `harborinfo` limit 0,1").Values(&maps)
	if err != nil {
		logs.Error("[SelectHarborInfo] select list error,%s", err)
	}
	return maps
}

func InsertProtocol(protocolname string, typeid string, softpath string, createtime int64) (map[string]interface{}, string, int) {
	msg := "成功"
	var data map[string]interface{}
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("insert into honeypotstype (honeypottype ,typeid, softpath,createtime) VALUES (?,?,?,?)", protocolname, typeid, softpath, createtime).Values(&maps)
	if err != nil {
		logs.Error("[InsertProtocol] insert servers error,%s", err)
		msg = "数据插入失败"
		return data, msg, comm.ErrorCode
	} else {
		msg = "成功"
	}
	return data, msg, comm.SuccessCode
}

func SelectHoneyImageById(imageid int) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT T0.id,T0.imageaddress,T0.imagename,T0.imageport,T1.typeid AS imagetype, T2.sysid AS imageos FROM podimage T0 LEFT JOIN honeypotstype T1 ON T0.imagetype = T1.typeid LEFT JOIN systemtype T2 ON T0.imageos = T2.sysid WHERE T0.id=?", imageid).Values(&maps)
	if err != nil {
		logs.Error("[SelectHoneyImageById] select event list error,%s", err)
	}
	return maps
}

func SelectHoneyImage(imagesname string, imagetype string, imageos string, pageSize int, pageNum int) map[string]interface{} {
	sqlCon, err1 := sql.Open("mysql", Dbuser+":"+Dbpassword+"@tcp("+Dbhost+":"+Dbport+")/"+Dbname+"?charset=utf8&loc=Asia%2FShanghai")
	if err1 != nil {
		logs.Error("[SelectHoneyInfos]open mysql fail %s", err1)
	}
	DbCon := sqlCon
	var total int
	defer sqlCon.Close()
	var condition string
	sqltotal := "SELECT COUNT(1) FROM podimage T0 LEFT JOIN honeypotstype T1 ON T0.imagetype = T1.typeid LEFT JOIN systemtype T2 ON T0.imageos = T2.sysid WHERE 1=1 "
	sqlstr := "SELECT T0.id,T0.imageaddress,T0.imagename,T0.imageport,T1.honeypottype AS imagetype, T2.systype AS imageos FROM podimage T0 LEFT JOIN honeypotstype T1 ON T0.imagetype = T1.typeid LEFT JOIN systemtype T2 ON T0.imageos = T2.sysid WHERE 1=1 "
	var a []interface{}
	if imagesname != "" {
		condition += " and imagename  like ?"
		a = append(a, "%"+imagesname+"%")

	}
	if imagetype != "" {
		condition += " and imagetype =?"
		a = append(a, imagetype)
	}
	if imageos != "" {
		condition += " and imageos =?"
		a = append(a, imageos)
	}

	sqltotal = sqltotal + condition
	err := DbCon.QueryRow(sqltotal, a...).Scan(&total)
	if err != nil {
		logs.Error("[SelectHoneyImage] select total error,%s", err)
	}
	totalpage, pagenum, pagesize := util.JudgePage(total, pageNum, pageSize)
	offset := (pagenum - 1) * pagesize
	condition += " ORDER BY id DESC"
	condition += " limit ?,?"
	a = append(a, offset)
	a = append(a, pagesize)
	rows, err := DbCon.Query(sqlstr+condition, a...)
	if err != nil {
		logs.Error("[SelectHoneyImage] select list error,%s", err)
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectHoneyImage] rows.Columns() error,%s", err)
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var data map[string]interface{}

	list, count, err := util.GetHoneyImagesListMysqlJson(rows, columns, total, values, scanArgs, pagenum, pagesize, totalpage)

	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			logs.Error("[SelectHoneyImage] Unmarshal list error,%s", err)
		}
	}
	return data
}

func SelectProtocol(pageSize int, pageNum int) map[string]interface{} {
	sqlCon, err1 := sql.Open("mysql", Dbuser+":"+Dbpassword+"@tcp("+Dbhost+":"+Dbport+")/"+Dbname+"?charset=utf8&loc=Asia%2FShanghai")
	if err1 != nil {
		logs.Error("[SelectProtocol]open mysql fail %s", err1)
	}
	DbCon := sqlCon
	var total int
	defer sqlCon.Close()
	var condition string
	sqltotal := "select count(1) from `honeypotstype` where 1=1"
	sqlstr := "select * from `honeypotstype` where 1=1"
	var a []interface{}

	sqltotal = sqltotal + condition
	err := DbCon.QueryRow(sqltotal, a...).Scan(&total)
	if err != nil {
		logs.Error("[SelectProtocol] select total error,%s", err)
	}
	totalpage, pagenum, pagesize := util.JudgePage(total, pageNum, pageSize)
	offset := (pagenum - 1) * pagesize
	condition += "  limit ?,?"
	a = append(a, offset)
	a = append(a, pagesize)
	sqlstr = sqlstr + condition

	rows, err := DbCon.Query(sqlstr, a...)
	if err != nil {
		logs.Error("[SelectProtocol] select list error,%s", err)
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectProtocol] rows.Columns() error,%s", err)
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var data map[string]interface{}

	list, count, err := util.GetProtocolListMysqlJson(rows, columns, total, values, scanArgs, pagenum, pagesize, totalpage)

	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			logs.Error("[SelectProtocol] Unmarshal list error,%s", err)
		}
	}
	return data
}

func DeleteProtocol(protocolid string) (map[string]interface{}, string, int) {
	msg := "成功"
	o := orm.NewOrm()
	var maps []orm.Params
	var data map[string]interface{}
	_, err := o.Raw("DELETE FROM honeypotstype where typeid=?", protocolid).Values(&maps)
	if err != nil {
		logs.Error("[DeleteProtocol] delete protocol error,%s", err)
		msg = "数据删除失败"
		return data, msg, comm.ErrorCode
	}
	return data, msg, comm.SuccessCode
}

func DeleteHarborInfo(harborid string) (map[string]interface{}, string, int) {
	msg := "成功"
	o := orm.NewOrm()
	var maps []orm.Params
	var data map[string]interface{}
	_, err := o.Raw("DELETE FROM harborinfo where harborid=?", harborid).Values(&maps)
	if err != nil {
		logs.Error("[DeleteHarborInfo] delete protocol error,%s", err)
		msg = "数据删除失败"
		return data, msg, comm.ErrorCode
	}
	return data, msg, comm.SuccessCode
}

func InsertSSHInfo(sshkey string, agentid string) (map[string]interface{}, string, int) {
	msg := "成功"
	var data map[string]interface{}
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("insert into honeyserverconfig (honeyserverid,serversshkey) VALUES (?,?)", agentid, sshkey).Values(&maps)
	if err != nil {
		logs.Error("[InsertSSHInfo] insert servers error,%s", err)
		msg = "SSHKey数据插入失败"
		return data, msg, comm.ErrorCode
	}
	return data, msg, comm.SuccessCode
}

func SelectApplicationByAgentID(agentid string) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT * from `servers` WHERE agentid=?", agentid).Values(&maps)
	if err != nil {
		logs.Error("[SelectApplicationByAgentID] select event list error,%s", err)
	}
	return maps
}

func SelectHoneyInfoByTransInfo(destaddr string, destport int) []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT T0.*,T2.honeytypeid from honeyfowards T0 LEFT JOIN honeypotservers T1 ON T1.agentid = T0.agentid LEFT JOIN honeypots T2 ON T0.honeypotid = T2.honeypotid WHERE T1.serverip = ? AND T0.forwardport = ?", destaddr, destport).Values(&maps)
	if err != nil {
		logs.Error("[SelectHoneyInfoByTransInfo] select event list error,%s", err)
	}
	return maps
}

func GetTraceHost() []orm.Params {
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("SELECT * FROM conf_traceinfo limit 0,1").Values(&maps)
	if err != nil {
		logs.Error("[GetTraceHost] select event list error,%s", err)
	}
	return maps
}
