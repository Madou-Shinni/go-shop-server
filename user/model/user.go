package model

type User struct {
	Model
	Mobile   string     `gorm:"index:idx_mobile;unique;type:varchar(11);not null;" json:"mobile,omitempty"`
	Password string     `gorm:"type:varchar(100);not null;" json:"password,omitempty"`
	NickName string     `gorm:"type:varchar(20);" json:"nickName,omitempty"`
	Birthday *LocalTime `gorm:"type:datetime;" json:"birthday,omitempty"`
	Gender   string     `gorm:"default:male;type:varchar(6);comment:male男,female女" json:"gender,omitempty"`
	Role     *int       `json:"role" gorm:"default:1;type:int;comment:1普通用户,2管理员" json:"role,omitempty"`
}
