package models

type Cars struct {
	CarsID     uint   `gorm:"primaryKey;column:idCars" json:"idCars" form:"idCars"`
	NamaMobil  string `gorm:"not null;column:nama_mobil" json:"nama_mobil" form:"nama_mobil"`
	MerekID    *uint  `gorm:"column:idMerek_fk" json:"-" form:"-"`
	JenisID    *uint  `gorm:"column:idJenis_fk" json:"-" form:"-"`
	HorsePower uint   `gorm:"column:horse_power" json:"horse_power" form:"horse_power"`

	// Relationships
	Merek *Merek `gorm:"foreignKey:idMerek_fk;references:idMerek" json:"merek"`
	Jenis *Jenis `gorm:"foreignKey:idJenis_fk;references:idJenis" json:"jenis"`
}

type Merek struct {
	ID   uint   `gorm:"column:idMerek;primaryKey" json:"id" form:"id"`
	Nama string `gorm:"column:merek" json:"merek" form:"merek"`
}

/* GORM Error Mitigation on Reading 'Merek' table as 'Mereks' */
func (Merek) TableName() string {
	return "merek" // Explicitly tells GORM to use "merek" not "mereks"
}

type Jenis struct {
	ID   uint   `gorm:"column:idJenis;primaryKey" json:"id" form:"id"`
	Nama string `gorm:"column:jenis" json:"jenis" form:"jenis"`
}
