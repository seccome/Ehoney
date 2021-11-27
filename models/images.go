package models

import (
	"decept-defense/controllers/comm"
	"fmt"
	"strings"
)

type Images struct {
	ID           int64  `gorm:"primary_key;AUTO_INCREMENT;unique;column:id" json:"id"`     //镜像ID
	ImageAddress string `json:"ImageAddress" form:"ImageAddress" gorm:"not null;size:256"` //镜像地址
	ImageName    string `json:"ImageName" form:"ImageName" gorm:"not null"`                //镜像名称
	ImagePort    int32  `json:"ImagePort" form:"ImagePort" gorm:"null"`                    //镜像端口
	ImageType    string `json:"ImageType" form:"ImageType" gorm:"null"`                    //镜像服务
	DefaultFlag  bool   `json:"DefaultFlag" form:"DefaultFlag" gorm:"null, default:false"` //默认属性
}

var DefaultImages = []Images{
	{ImageAddress: "47.96.71.197:90/ehoney/tomcat:v1", ImageName: "ehoney/tomcat", ImagePort: 8080, ImageType: "httpproxy", DefaultFlag: true},
	{ImageAddress: "47.96.71.197:90/ehoney/ssh:v1", ImageName: "ehoney/ssh", ImagePort: 22, ImageType: "sshproxy", DefaultFlag: true},
	{ImageAddress: "47.96.71.197:90/ehoney/mysql:v1", ImageName: "ehoney/mysql", ImagePort: 3306, ImageType: "mysqlproxy", DefaultFlag: true},
	{ImageAddress: "47.96.71.197:90/ehoney/redis:v1", ImageName: "ehoney/redis", ImagePort: 6379, ImageType: "redisproxy", DefaultFlag: true},
	{ImageAddress: "47.96.71.197:90/ehoney/telnet:v1", ImageName: "ehoney/telnet", ImagePort: 23, ImageType: "telnetproxy", DefaultFlag: true},
	{ImageAddress: "47.96.71.197:90/ehoney/smb:v1", ImageName: "ehoney/smb", ImagePort: 445, ImageType: "smbproxy", DefaultFlag: true}, // 由于smb client 必须连接 445 所以改4450 445 留给协议代理
	{ImageAddress: "47.96.71.197:90/ehoney/ftp:v1", ImageName: "ehoney/ftp", ImagePort: 21, ImageType: "ftpproxy", DefaultFlag: true},
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
	sql := fmt.Sprintf("select id, image_name, image_address, image_port, image_type, default_flag from images where CONCAT(image_name, image_address, image_port, image_type) LIKE '%s'", p)
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

func (image *Images) UpdateImageByID(id int64, payload comm.ImageUpdatePayload) error {
	if err := db.Model(image).Where("id = ?", id).Updates(Images{ImagePort: payload.ImagePort, ImageType: payload.ImageType}).Error; err != nil {
		return err
	}
	return nil
}

func (image *Images) GetImageByID(id int64) (*Images, error) {
	var ret Images
	if err := db.Where("id = ?", id).Take(&ret).Error; err != nil {
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

func (image *Images) GetPodImageList() (*[]string, error) {
	var ret []string
	if err := db.Model(image).Select("image_address").Where("image_type != '' AND image_port != 0 ").Find(&ret).Error; err != nil {
		return nil, err
	}
	return &ret, nil
}
