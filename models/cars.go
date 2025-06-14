package models

type Cars struct {
	CarsID     uint   `gorm:"primaryKey;column:idCars" json:"idCars" form:"idCars"`
	NamaMobil  string `gorm:"not null;column:nama_mobil" json:"nama_mobil" form:"nama_mobil"`
	MerekID    *uint  `gorm:"column:idMerek_fk;type:int unsigned" json:"-" form:"-"`
	JenisID    *uint  `gorm:"column:idJenis_fk;type:int unsigned" json:"-" form:"-"`
	HorsePower uint   `gorm:"column:horse_power" json:"horse_power" form:"horse_power"`

	// Relationships
	Merek *Merek `gorm:"foreignKey:MerekID;references:idMerek" json:"merek"`
	Jenis *Jenis `gorm:"foreignKey:JenisID;references:idJenis" json:"jenis"`
}

type Merek struct {
	ID    uint   `gorm:"column:idMerek;primaryKey;autoIncrement;type:int unsigned" json:"id" form:"id"`
	Merek string `gorm:"column:merek" json:"merek" form:"merek"`
}

type Jenis struct {
	ID    uint   `gorm:"column:idJenis;primaryKey;autoIncrement;type:int unsigned" json:"id" form:"id"`
	Jenis string `gorm:"column:jenis" json:"jenis" form:"jenis"`
}
