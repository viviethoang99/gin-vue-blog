package model

import "gorm.io/gorm"

type OperationLog struct {
	Model

	OptModule string `gorm:"type:varchar(50);comment:Operation module" json:"opt_module"`
	OptType   string `gorm:"type:varchar(50);comment:Operation type" json:"opt_type"`
	OptMethod string `gorm:"type:varchar(100);comment:Operation method" json:"opt_method"`
	OptUrl    string `gorm:"type:varchar(255);comment:Operation URL" json:"opt_url"`
	OptDesc   string `gorm:"type:varchar(255);comment:Operation description" json:"opt_desc"`

	RequestParam  string `gorm:"type:longtext;comment:Request parameters" json:"request_param"`
	RequestMethod string `gorm:"type:longtext;comment:Request method" json:"request_method"`
	ResponseData  string `gorm:"type:longtext;comment:Response data" json:"response_data"`

	UserId    int    `gorm:"comment:User ID" json:"user_id"`
	Nickname  string `gorm:"type:varchar(50);comment:User nickname" json:"nickname"`
	IpAddress string `gorm:"type:varchar(255);comment:Operation IP" json:"ip_address"`
	IpSource  string `gorm:"type:varchar(255);comment:Operation address" json:"ip_source"`
}

func GetOperationLogList(db *gorm.DB, num, size int, keyword string) (data []OperationLog, total int64, err error) {
	db = db.Model(&OperationLog{})
	if keyword != "" {
		db = db.Where("opt_module LIKE ?", "%"+keyword+"%").
			Or("opt_desc LIKE ?", "%"+keyword+"%")
	}
	db.Count(&total)
	result := db.Order("created_at DESC").
		Scopes(Paginate(num, size)).
		Find(&data)
	return data, total, result.Error
}
