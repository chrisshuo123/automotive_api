-- MySQL dump 10.13  Distrib 8.0.31, for Win64 (x86_64)
--
-- Host: localhost    Database: automotive_api
-- ------------------------------------------------------
-- Server version	8.0.31

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `cars`
--

DROP TABLE IF EXISTS `cars`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `cars` (
  `idCars` int NOT NULL AUTO_INCREMENT,
  `nama_mobil` longtext NOT NULL,
  `idMerek_fk` int DEFAULT NULL,
  `idJenis_fk` int DEFAULT NULL,
  `horse_power` int DEFAULT NULL,
  PRIMARY KEY (`idCars`),
  KEY `idMerek_fk` (`idMerek_fk`),
  KEY `idJenis_fk` (`idJenis_fk`),
  CONSTRAINT `idJenis_fk` FOREIGN KEY (`idJenis_fk`) REFERENCES `jenis` (`idJenis`),
  CONSTRAINT `idMerek_fk` FOREIGN KEY (`idMerek_fk`) REFERENCES `merek` (`idMerek`)
) ENGINE=InnoDB AUTO_INCREMENT=47 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cars`
--

LOCK TABLES `cars` WRITE;
/*!40000 ALTER TABLE `cars` DISABLE KEYS */;
INSERT INTO `cars` VALUES (1,'Toyota Yaris',1,1,106),(2,'Honda Jazz',2,1,118),(3,'Suzuki Swift',3,1,89),(4,'Mitsubishi Mirage',4,1,76),(5,'Toyota Avanza',1,2,103),(6,'Honda Mobilio',2,2,118),(7,'Suzuki Ertiga',3,2,103),(8,'Mitsubishi Expander',4,6,104),(9,'Toyota RAV4',1,3,203),(10,'Honda CR-V',2,3,190),(11,'Suzuki Jimny',3,3,101),(12,'Mitsubishi Outlander',4,3,181),(13,'Toyota Alphard',1,4,275),(14,'Honda Odyssey',2,4,212),(15,'Suzuki Every',3,5,63),(16,'Mitsubishi Delica',4,5,147),(17,'Toyota Sienta',1,6,109),(18,'Honda Freed',2,6,129),(19,'Suzuki Spacia',3,6,52),(20,'Mitsubishi Expander',4,6,104);
/*!40000 ALTER TABLE `cars` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `jenis`
--

DROP TABLE IF EXISTS `jenis`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `jenis` (
  `idJenis` int NOT NULL AUTO_INCREMENT,
  `jenis` longtext,
  `nama` longtext,
  PRIMARY KEY (`idJenis`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `jenis`
--

LOCK TABLES `jenis` WRITE;
/*!40000 ALTER TABLE `jenis` DISABLE KEYS */;
INSERT INTO `jenis` VALUES (1,'Hatchback',NULL),(2,'MPV',NULL),(3,'SUV',NULL),(4,'Minivan',NULL),(5,'Van',NULL),(6,'Low-MPV',NULL);
/*!40000 ALTER TABLE `jenis` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `merek`
--

DROP TABLE IF EXISTS `merek`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `merek` (
  `idMerek` int NOT NULL AUTO_INCREMENT,
  `merek` longtext,
  PRIMARY KEY (`idMerek`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=latin1;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `merek`
--

LOCK TABLES `merek` WRITE;
/*!40000 ALTER TABLE `merek` DISABLE KEYS */;
INSERT INTO `merek` VALUES (1,'Toyota'),(2,'Honda'),(3,'Suzuki'),(4,'Mitsubishi');
/*!40000 ALTER TABLE `merek` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2025-06-30 20:30:12
