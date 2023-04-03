package cdn

type Resource struct {
	ID   uint   `goal:"<number>primary,sortable,desc" gorm:"primaryKey"`
	File string `binding:"required" goal:"<file>uploadTo=uploads,postonly,unique" gorm:"unique"`
}
