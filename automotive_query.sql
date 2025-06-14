use automotive_api;
select * from cars;

create table cars (
	idCars int(10) primary key not null auto_increment,
    nama_mobil varchar(100) not null,
    idMerek_fk int(10),
    idJenis_fk int(10),
    horse_power int(10)
);

alter table cars
drop column bio;

alter table cars
	add column merek varchar(100) after nama_mobil,
    add column horse_power int(10) after jenis;
alter table cars
	modify column nama_mobil varchar(100) not null;

create table merek (
	idMerek int(10) primary key auto_increment,
	merek varchar(100) not null
);

insert merek (merek)
	value('Toyota'), ('Honda'), ('Suzuki'), ('Mitsubishi');

create table jenis (
	idJenis int(10) primary key auto_increment,
    jenis varchar(100)
);
alter table jenis
	modify column jenis varchar(100) not null;

insert jenis(jenis)
	values('Hatchback'), ('MPV'), ('SUV'), ('Minivan'), ('Van');
insert jenis(jenis)
	values('Low-MPV');

alter table cars
	drop column merek,
    drop column jenis,
    add column idMerek_id int(10) not null after nama_mobil,
    add column idJenis_id int(10) not null after idMerek_id;

alter table cars
	add constraint idMerek_fk foreign key (idMerek_fk)
    references merek(idMerek),
    add constraint idJenis_fk foreign key (idJenis_fk)
    references jenis(idJenis);
    
alter table cars
	drop constraint idMerek_id;
alter table cars
	drop constraint idJenis_id;
    
alter table cars
	rename column idMerek_id to idMerek_fk,
    rename column idJenis_id to idJenis_fk;

SHOW CREATE TABLE cars;
SELECT CONSTRAINT_NAME 
FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS
WHERE TABLE_NAME = 'cars';
SHOW ENGINE INNODB STATUS;

describe cars;
describe merek;
describe jenis;

select * from cars;
select * from merek;
/* List Merek:
1 - Toyota
2 - Honda
3 - Suzuki
4 - Mitsubishi */

select * from jenis;
/* List Jenis:
1 - Hatchback
2 - MPV
3 - SUV
4 - Minivan
5 - Van 
6 - Low-MPV */

alter table cars
	rename column id to idCars;
alter table cars
    modify column idCars int(10) auto_increment;

insert cars(nama_mobil, idMerek_fk, idJenis_fk, horse_power)
	values
		('Toyota Yaris', 1, 1, 106),
        ('Honda Jazz', 2, 1, 118),
        ('Suzuki Swift', 3, 1, 89),
        ('Mitsubishi Mirage', 4, 1, 76),
        ('Toyota Avanza', 1, 2, 103),
        ('Honda Mobilio', 2, 2, 118),
		('Suzuki Ertiga', 3, 2, 103),
        ('Mitsubishi Expander', 4, 2, 104),
        ('Toyota RAV4', 1, 3, 203),
        ('Honda CR-V', 2, 3, 190),
        ('Suzuki Jimny', 3, 3, 101),
        ('Mitsubishi Outlander', 4, 3, 181),
        ('Toyota Alphard', 1, 4, 275),
        ('Honda Odyssey', 2, 4, 212),
        ('Suzuki Every', 3, 5, 63),
        ('Mitsubishi Delica', 4, 5, 147),
        ('Toyota Sienta', 1, 6, 109),
        ('Honda Freed', 2, 6, 129),
        ('Suzuki Spacia', 3, 6, 52),
        ('Mitsubishi Expander', 4, 6, 104);    
    
    select * from cars;
    describe cars;
    
/* Integrate the Table outputs */
select
	c.idCars,
    c.nama_mobil,
    m.merek,
    j.jenis,
    c.horse_power
from
	cars c
left join
	merek m on c.idMerek_fk = m.idMerek
left join
	jenis j on c.idJenis_fk = j.idJenis;
    
select * from cars;

/* Merek Table Error: GORM di Golang Failed To Load Merek Table DB.
Perlu melakukan pengecekan pada Constraintnya */
-- Check if the foreign key constraints exist
SELECT CONSTRAINT_NAME 
FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS 
WHERE TABLE_NAME = 'cars' AND CONSTRAINT_TYPE = 'FOREIGN_KEY';

SELECT CONSTRAINT_NAME 
FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS
WHERE TABLE_NAME = 'cars';

-- Check if the referenced records exist
SELECT idMerek FROM merek WHERE idMerek IN (1,2,3,4);
    
/* The FK Constraint are correct, but when you build the FK Constraint,
It just makes the GORM to be Confused.  With that, you need to drop the FK
Constraint, and recreate with the new one. */
ALTER TABLE cars DROP FOREIGN KEY idMerek_fk;
ALTER TABLE cars DROP FOREIGN KEY idJenis_fk;

ALTER TABLE cars MODIFY COLUMN idMerek_fk INT(10) UNSIGNED NULL;
ALTER TABLE cars MODIFY COLUMN idJenis_fk INT(10) UNSIGNED NULL;

SELECT COLUMN_NAME, COLUMN_TYPE 
FROM INFORMATION_SCHEMA.COLUMNS 
WHERE TABLE_NAME IN ('cars', 'merek') 
AND COLUMN_NAME IN ('idMerek_fk', 'idMerek');

use automotive_api;
-- Create a stored procedure to safely drop constraints
DELIMITER //
CREATE PROCEDURE safe_drop_fk(IN table_name VARCHAR(64), IN fk_name VARCHAR(64))
BEGIN
    DECLARE fk_exists INT;
    
    -- Check if foreign key exists
    SELECT COUNT(*) INTO fk_exists
    FROM INFORMATION_SCHEMA.REFERENTIAL_CONSTRAINTS
    WHERE CONSTRAINT_NAME = fk_name AND TABLE_NAME = table_name;
    
    -- Drop if exists
    IF fk_exists > 0 THEN
        SET @sql = CONCAT('ALTER TABLE ', table_name, ' DROP FOREIGN KEY ', fk_name);
        PREPARE stmt FROM @sql;
        EXECUTE stmt;
        DEALLOCATE PREPARE stmt;
    END IF;
END//
DELIMITER ;

-- Use the procedure to drop your FKs
CALL safe_drop_fk('cars', 'idMerek_fk');
CALL safe_drop_fk('cars', 'idJenis_fk');

-- Then proceed with your modifications
ALTER TABLE merek MODIFY idMerek INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;
ALTER TABLE cars MODIFY idMerek_fk INT(10) UNSIGNED NULL;

ALTER TABLE jenis MODIFY idJenis INT(10) UNSIGNED NOT NULL AUTO_INCREMENT;
ALTER TABLE cars MODIFY idJenis_fk INT(10) UNSIGNED NULL;

-- Recreate constraints
ALTER TABLE cars ADD CONSTRAINT idMerek_fk
FOREIGN KEY (idMerek_fk) REFERENCES merek(idMerek);
/*ON DELETE SET NULL ON UPDATE CASCADE;*/

ALTER TABLE cars ADD CONSTRAINT idJenis_fk
FOREIGN KEY (idJenis_fk) REFERENCES jenis(idJenis);
/*ON DELETE SET NULL ON UPDATE CASCADE;*/

-- Clean up
DROP PROCEDURE IF EXISTS safe_drop_fk;

alter table cars
add constraint idMerek_fk
foreign key (idMerek_fk) references merek(idMerek)
ON DELETE SET NULL ON UPDATE CASCADE;

alter table cars
add constraint idJenis_fk
foreign key (idJenis_fk) references jenis(idJenis)
ON DELETE SET NULL ON UPDATE CASCADE;

use automotive_api;
select * from cars;