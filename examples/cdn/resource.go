package cdn

type Resource struct {
	ID   uint   `goal:"<number>primary,hidden" gorm:"primaryKey"`
	File string `binding:"required" goal:"<file>uploadTo=uploads,postonly,unique" gorm:"unique"`
}
