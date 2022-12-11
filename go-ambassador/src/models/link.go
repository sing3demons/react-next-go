package models

type Link struct {
	ID       uint
	Code     string
	UserId   uint
	User     User      `gorm:"foreignKey:UserId"`
	Products []Product `gorm:"many2many:link_products"`
	Orders   []Order   `gorm:"-"`
}
