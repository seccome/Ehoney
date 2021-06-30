package clamavcenter

import (
	"database/sql"
	"decept-defense/models/util"
	"decept-defense/models/util/comm"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

var (
	Dbhost     = beego.AppConfig.String("dbhost")
	Dbport     = beego.AppConfig.String("dbport")
	Dbuser     = beego.AppConfig.String("dbuser")
	Dbpassword = beego.AppConfig.String("dbpassword")
	Dbname     = beego.AppConfig.String("dbname")
)

func InsertClamavData(filename string, virname string, createtime int64, honeypotip string) (map[string]interface{}, string, int) {
	msg := "成功"
	var data map[string]interface{}
	o := orm.NewOrm()
	var maps []orm.Params
	_, err := o.Raw("insert into clamavdate (filename,virname,createtime,honeypotip) VALUES (?,?,?,?)", filename, virname, createtime, honeypotip).Values(&maps)
	if err != nil {
		logs.Error("[InsertClamavData] insert servers error,%s", err)
		msg = "数据插入失败"
		return data, msg, comm.ErrorCode
	}
	return data, msg, comm.SuccessCode
}

func SelectClamavData(filename string, virname string, honeypotip string, starttime string, endtime string, pageSize int, pageNum int) map[string]interface{} {
	sqlCon, err1 := sql.Open("mysql", Dbuser+":"+Dbpassword+"@tcp("+Dbhost+":"+Dbport+")/"+Dbname+"?charset=utf8&loc=Asia%2FShanghai")
	if err1 != nil {
		logs.Error("[SelectClamavData]open mysql fail %s", err1)
	}
	DbCon := sqlCon
	var total int
	defer sqlCon.Close()
	var condition string
	sqltotal := "select count(1) from `clamavdate` where 1=1"
	sqlstr := "select * from `clamavdate` where 1=1"
	var a []interface{}
	if filename != "" {
		condition += " and filename like ?"
		a = append(a, "%"+filename+"%")
	}
	if virname != "" {
		condition += " and virname like ? "
		a = append(a, "%"+virname+"%")
	}
	if honeypotip != "" {
		condition += " and honeypotip=?"
		a = append(a, honeypotip)
	}
	if starttime != "" && endtime != "" {
		condition += " and createtime between ? and ?"
		a = append(a, starttime)
		a = append(a, endtime)
	}
	sqltotal = sqltotal + condition
	err := DbCon.QueryRow(sqltotal, a...).Scan(&total)
	if err != nil {
		logs.Error("[SelectClamavData] select total error,%s", err)
	}
	totalpage, pagenum, pagesize := util.JudgePage(total, pageNum, pageSize)
	offset := (pagenum - 1) * pagesize
	condition += "  limit ?,?"
	a = append(a, offset)
	a = append(a, pagesize)

	sqlstr = sqlstr + condition
	rows, err := DbCon.Query(sqlstr, a...)
	if err != nil {
		logs.Error("[SelectClamavData] select list error,%s", err)
	}
	columns, err := rows.Columns()
	if err != nil {
		logs.Error("[SelectClamavData] rows.Columns() error,%s", err)
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	var data map[string]interface{}

	list, count, err := util.GetClamavDataListMysqlJson(rows, columns, total, values, scanArgs, pagenum, pagesize, totalpage)

	if count > 0 {
		err = json.Unmarshal([]byte(list), &data)
		if err != nil {
			logs.Error("[SelectClamavData] Unmarshal list error,%s", err)
		}
	}
	return data
}
