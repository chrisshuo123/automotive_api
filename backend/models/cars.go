package models

type Cars struct {
	CarsID     uint   `gorm:"primaryKey;column:idcars" json:"idcars" form:"idcars"`
	NamaMobil  string `gorm:"not null;column:nama_mobil" json:"nama_mobil" form:"nama_mobil"`
	MerekID    *uint  `gorm:"column:idmerek_fk" json:"idmerek_fk" form:"idmerek_fk"`
	JenisID    *uint  `gorm:"column:idjenis_fk" json:"idjenis_fk" form:"idjenis_fk"`
	HorsePower uint   `gorm:"column:horse_power" json:"horse_power" form:"horse_power"`
	StatusID   *uint  `gorm:"column:idstatus_fk" json:"idstatus_fk" form:"idstatus_fk"`
	ImageCar   string `gorm:"column:nama_foto" json:"nama_foto" form:"nama_foto"`

	// Relationships
	Merek  *Merek  `gorm:"foreignKey:idmerek_fk;references:idmerek" json:"merek"`
	Jenis  *Jenis  `gorm:"foreignKey:idjenis_fk;references:idjenis" json:"jenis"`
	Status *Status `gorm:"foreignKey:idstatus_fk;references:idstatus" json:"status"`
}

type Merek struct {
	ID   uint   `gorm:"column:idmerek;primaryKey" json:"idmerek" form:"id"`
	Nama string `gorm:"column:namamerek" json:"namamerek" form:"namamerek"`
}

// GORM to find Postgree's Tablename:
func (Cars) TableName() string {
	return "cars"
}
func (Jenis) TableName() string {
	return "jenis"
}
func (Status) TableName() string {
	return "status"
}

/* GORM Error Mitigation on Reading 'Merek' table as 'Mereks' */
func (Merek) TableName() string {
	return "merek" // Explicitly tells GORM to use "merek" not "mereks"
}

type Jenis struct {
	ID   uint   `gorm:"column:idjenis;primaryKey" json:"idjenis" form:"id"`
	Nama string `gorm:"column:namajenis" json:"namajenis" form:"namajenis"`
}

type Status struct {
	ID   uint   `gorm:"column:idstatus;primaryKey" json:"idstatus" form:"id"`
	Nama string `gorm:"column:namastatus" json:"namastatus" form:"namastatus"`
}
