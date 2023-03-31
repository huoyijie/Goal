package cdn

type Resource struct {
	ID   uint   `goal:"<number>primary,hidden" gorm:"primaryKey"`
	Path string `goal:"<text>uploadTo=uploads,unique" gorm:"unique"`
}
