package models

import (
	"decept-defense/controllers/comm"
	"decept-defense/pkg/configs"
	"decept-defense/pkg/util"
	"fmt"
	"path"
	"strings"
)

type Token struct {
	ID          int64  `gorm:"primary_key;AUTO_INCREMENT;not null;unique;column:id" json:"id"`            //密签ID
	CreateTime  string `gorm:"not null"`                                                                  //创建时间
	Creator     string `gorm:"not null;size:256"`                                                         //创建用户
	TokenType   string `gorm:"not null;size:256" form:"TokenType" binding:"required"`                     //密签类型
	UploadPath  string `gorm:"null;size:256"`                                                             //上传路径
	FileName    string `gorm:"null;size:256"`                                                             //文件名称
	TokenName   string `gorm:"not null;unique"form:"TokenName" gorm:"unique;size:256" binding:"required"` //密签名称
	TokenData   string `gorm:"null" form:"TokenData"`                                                     //密签内容
	DefaultFlag bool   `json:"DefaultFlag" form:"DefaultFlag" gorm:"null, default:false"`
}

func (token *Token) CreateToken() error {
	if err := db.Create(token).Error; err != nil {
		return err
	}
	return nil
}

func (token *Token) CreateDefaultToken() error {
	var DefaultTokens = []Token{
		{TokenName: "token_pdf", FileName: "token.pdf", TokenType: "FILE", UploadPath: path.Join(util.WorkingPath(), configs.GetSetting().App.UploadPath, "token", "token_pdf", "token.pdf"), Creator: "default", CreateTime: util.GetCurrentTime(), DefaultFlag: true},
		{TokenName: "token_docx", FileName: "token.docx", TokenType: "FILE", UploadPath: path.Join(util.WorkingPath(), configs.GetSetting().App.UploadPath, "token", "token_docx", "token.docx"), Creator: "default", CreateTime: util.GetCurrentTime(), DefaultFlag: true},
		{TokenName: "token_pptx", FileName: "token.pptx", TokenType: "FILE", UploadPath: path.Join(util.WorkingPath(), configs.GetSetting().App.UploadPath, "token", "token_pptx", "token.pptx"), Creator: "default", CreateTime: util.GetCurrentTime(), DefaultFlag: true},
		{TokenName: "token_xlsx", FileName: "token.xlsx", TokenType: "FILE", UploadPath: path.Join(util.WorkingPath(), configs.GetSetting().App.UploadPath, "token", "token_xlsx", "token.xlsx"), Creator: "default", CreateTime: util.GetCurrentTime(), DefaultFlag: true},
	}
	for _, d := range DefaultTokens {
		p, _ := token.GetTokenByName(d.TokenName)
		if p != nil {
			continue
		}
		err := d.CreateToken()
		if err != nil {
			continue
		}
	}
	return nil
}

func (token *Token) GetToken(payload *comm.SelectPayload) (*[]comm.TokenSelectResultPayload, int64, error) {
	var ret []comm.TokenSelectResultPayload
	var count int64
	if util.CheckInjectionData(payload.Payload) {
		return nil, 0, nil
	}
	var p string = "%" + payload.Payload + "%"
	sql := fmt.Sprintf("select id, token_type, token_name, file_name, create_time, creator, default_flag from tokens where CONCAT(token_type, token_name, file_name, create_time, creator) LIKE '%s' order by create_time DESC", p)
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, 0, err
	}
	count = (int64)(len(ret))
	t := fmt.Sprintf("limit %d offset %d", payload.PageSize, (payload.PageNumber-1)*payload.PageSize)
	sql = strings.Join([]string{sql, t}, " ")
	if err := db.Raw(sql).Scan(&ret).Error; err != nil {
		return nil, 0, err
	}
	return &ret, count, nil
}

func (token *Token) GetTokenByID(id int64) (*Token, error) {
	var ret Token
	if err := db.Take(&ret, id).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (token *Token) GetTokenByName(name string) (*Token, error) {
	var ret Token
	if err := db.Where("token_name = ?", name).Take(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (token *Token) GetTokenNameList() ([]string, error) {
	var ret []string
	if err := db.Select("token_name").Find(token).Scan(&ret).Error; err != nil {
		return nil, err
	}
	return ret, nil
}

func (token *Token) DeleteTokenByID(id int64) error {
	if err := db.Delete(&Token{}, id).Error; err != nil {
		return err
	}
	return nil
}
