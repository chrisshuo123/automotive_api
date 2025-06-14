package configs

import "gorm.io/gorm"

func FixForeignKeys(db *gorm.DB) error {
	// Disable foreign key checks temporarily
	db.Exec("SET FOREIGN_KEY_CHECKS=0")

	var count int
	db.Raw(`
		SELECT COUNT(*)
		FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS
		WHERE TABLE_NAME = 'cars'
		AND CONSTRAINT_TYPE = 'FOREIGN KEY'
	`).Scan(&count)

	if count == 0 {
		db.Exec(`
			ALTER TABLE cars
			ADD CONSTRAINT idMerek_fk
			FOREIGN KEY (idMerek_fk) REFERENCES merek(idMerek)
			ON DELETE SET NULL
		`)
		db.Exec(`
			ALTER TABLE cars 
            ADD CONSTRAINT idJenis_fk 
            FOREIGN KEY (idJenis_fk) REFERENCES jenis(idJenis)
            ON DELETE SET NULL
		`)
	}

	// Re-enable foreign key checks
	db.Exec("SET FOREIGN_KEY_CHECKS=1")
	return nil
}
