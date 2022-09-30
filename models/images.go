package models

import (
	"decept-defense/controllers/comm"
	"decept-defense/pkg/util"
	"fmt"
	"strings"
)

type Images struct {
	Id           int64  `gorm:"primary_key;AUTO_INCREMENT;unique;column:id" json:"id"`     //镜像ID
	ImageId      string `json:"ImageId" form:"ImageId" gorm:"unique;not null;;size:64"`    //镜像ID
	ImageAddress string `json:"ImageAddress" form:"ImageAddress" gorm:"not null;size:128"` //镜像地址
	ImageName    string `json:"ImageName" form:"ImageName" gorm:"not null;size:128"`       //镜像名称
	ImagePort    int32  `json:"ImagePort" form:"ImagePort" gorm:"null"`                    //镜像端口
	ProtocolType string `json:"ProtocolType" form:"ProtocolType" gorm:"null;size:128"`     //镜像服务
	Label        string `json:"Label" form:"Label" gorm:"null;size:128"`                   //标签
	RepositoryId string `json:"RepositoryId" form:"RepositoryId" gorm:"null;size:128"`     //仓库Id 后续为了自定义镜像准备
	CreateTime   int64  `form:"CreateTime" json:"CreateTime"`
	DefaultFlag  bool   `json:"DefaultFlag" form:"DefaultFlag" gorm:"null, default:false"` //默认属性
}

// dckr_pat_d5BJ8k3y1QHVGPmkg9NMbGUVhJw  docker access token R-Only
// dckr_pat_qeZ5sDei_wkl8nWxil8s_xvQR-Q docker access token R-W-D

var DefaultImages = []Images{
	{ImageAddress: "ehoney/tomcat:v1", ImageId: util.GenerateId(), ImageName: "ehoney/tomcat", ImagePort: 8080, ProtocolType: "httpproxy", Label: "Default", CreateTime: util.GetCurrentIntTime(), DefaultFlag: true},
	{ImageAddress: "ehoney/ssh:v1", ImageId: util.GenerateId(), ImageName: "ehoney/ssh", ImagePort: 22, ProtocolType: "sshproxy", Label: "Default", CreateTime: util.GetCurrentIntTime(), DefaultFlag: true},
	{ImageAddress: "ehoney/mysql:v1", ImageId: util.GenerateId(), ImageName: "ehoney/mysql", ImagePort: 3306, ProtocolType: "mysqlproxy", Label: "Default", CreateTime: util.GetCurrentIntTime(), DefaultFlag: true},
	{ImageAddress: "ehoney/redis:v1", ImageId: util.GenerateId(), ImageName: "ehoney/redis", ImagePort: 6379, ProtocolType: "redisproxy", Label: "Default", CreateTime: util.GetCurrentIntTime(), DefaultFlag: true},
	{ImageAddress: "ehoney/telnet:v1", ImageId: util.GenerateId(), ImageName: "ehoney/telnet", ImagePort: 23, ProtocolType: "telnetproxy", Label: "Default", CreateTime: util.GetCurrentIntTime(), DefaultFlag: true},
	{ImageAddress: "ehoney/smb:v1", ImageId: util.GenerateId(), ImageName: "ehoney/smb", ImagePort: 445, ProtocolType: "smbproxy", Label: "Default", CreateTime: util.GetCurrentIntTime(), DefaultFlag: true}, // 由于smb client 必须连接 445 所以改4450 445 留给协议代理
	{ImageAddress: "ehoney/ftp:v1", ImageId: util.GenerateId(), ImageName: "ehoney/ftp", ImagePort: 21, ProtocolType: "ftpproxy", Label: "Default", CreateTime: util.GetCurrentIntTime(), DefaultFlag: true},
}

func (image *Images) CreateImage() error {
	ret, _ := image.GetImageByAddress(image.ImageAddress)
	if ret == nil {
		result := db.Create(image)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

func (image *Images) CreateDefaultImage() error {
	for _, d := range DefaultImages {
		err := d.CreateImage()
		if err != nil {
			continue
		}
	}
	return nil
}

func (image *Images) GetImage(payload *comm.SelectPayload) (*[]Images, int64, error) {
	var ret []Images
	var count int64
	var p = "%" + payload.Payload + "%"
	sql := fmt.Sprintf("select * from images where CONCAT(image_name, image_address, image_port, protocol_type) LIKE '%s'", p)
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

func (image *Images) UpdateImageByID(id string, payload comm.ImageUpdatePayload) error {
	if err := db.Model(image).Where("image_id = ?", id).Updates(Images{ImagePort: payload.ImagePort, ProtocolType: payload.ImageType, Label: payload.Label}).Error; err != nil {
		return err
	}
	return nil
}

func (image *Images) GetImageByID(id string) (*Images, error) {
	var ret Images
	if err := db.Where("image_id = ?", id).Take(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (image *Images) GetImageByAddress(address string) (*Images, error) {
	var ret Images
	if err := db.Take(&ret, "image_address = ?", address).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}

func (image *Images) DeleteImages() error {
	if err := db.Where("1 = 1").Delete(&Images{}).Error; err != nil {
		return err
	}
	return nil
}

func (image *Images) DeleteImageById(id string) error {
	if err := db.Where("image_id= ?", id).Delete(&Images{}).Error; err != nil {
		return err
	}
	return nil
}

func (image *Images) GetPodImageList() (*[]string, error) {
	var ret []string
	if err := db.Model(image).Select("image_address").Where("protocol_type != '' AND image_port != 0 ").Find(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}
