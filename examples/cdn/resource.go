package cdn

type Resource struct {
	ID   uint   `gorm:"primaryKey"`
	Path string `goal:"uploadTo=uploads" gorm:"unique"`
}
