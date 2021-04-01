package main

import (
	"fmt"
	"github.com/petrjahoda/database"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"math/rand"
	"time"
)

const version = "2021.1.3.30"

const config = "user=postgres password=pj79.. dbname=system host=localhost port=5432 sslmode=disable"

func main() {
	end := time.Now().AddDate(0, 0, -1)
	beginning := end.AddDate(0, -3, 0)

	fmt.Println(beginning)
	fmt.Println(end)
	db, err := gorm.Open(postgres.Open(config), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	if err != nil {
	}
	var devices []database.Device
	db.Where("device_type_id = 1").Find(&devices)
	for _, device := range devices {
		timer := time.Now()
		insertTime := beginning
		var analogPort database.DevicePort
		db.Where("device_id = ? and device_port_type_id = 2", device.ID).Find(&analogPort)
		var digitalPort database.DevicePort
		db.Where("device_id = ? and device_port_type_id = 1", device.ID).Find(&digitalPort)
		var analogRecordsToInsert []database.DevicePortAnalogRecord
		var digitalRecordsToInsert []database.DevicePortDigitalRecord
		for insertTime.Before(end) {
			generatedState := rand.Intn(3)
			generatedDuration := rand.Intn(180)
			finalTime := insertTime.Add(time.Duration(generatedDuration) * time.Minute)
			if generatedState == 1 { // Production
				for insertTime.Before(finalTime) {
					var digitalData database.DevicePortDigitalRecord
					digitalData.DateTime = insertTime
					digitalData.DevicePortID = int(digitalPort.ID)
					digitalData.Data = 1
					digitalRecordsToInsert = append(digitalRecordsToInsert, digitalData)
					var digitalData2 database.DevicePortDigitalRecord
					digitalData2.DateTime = insertTime.Add(1 * time.Second)
					digitalData2.DevicePortID = int(digitalPort.ID)
					digitalData2.Data = 0
					digitalRecordsToInsert = append(digitalRecordsToInsert, digitalData2)
					var analogData database.DevicePortAnalogRecord
					analogData.DateTime = insertTime
					analogData.DevicePortID = int(analogPort.ID)
					analogData.Data = float32(rand.Intn(100))
					analogRecordsToInsert = append(analogRecordsToInsert, analogData)
					insertTime = insertTime.Add(10 * time.Second)
				}
			} else if generatedState == 2 { // Downtime
				for insertTime.Before(finalTime) {
					var analogData database.DevicePortAnalogRecord
					analogData.DateTime = insertTime
					analogData.DevicePortID = int(analogPort.ID)
					analogData.Data = float32(rand.Intn(20))
					analogRecordsToInsert = append(analogRecordsToInsert, analogData)
					insertTime = insertTime.Add(10 * time.Second)
				}

			} else {
				insertTime = insertTime.Add(time.Duration(generatedDuration) * time.Minute)
			}
			if len(analogRecordsToInsert) > 1000 {
				db.Clauses(clause.OnConflict{DoNothing: true}).Create(&analogRecordsToInsert)
				analogRecordsToInsert = nil
			}
			if len(digitalRecordsToInsert) > 1000 {
				db.Clauses(clause.OnConflict{DoNothing: true}).Create(&digitalRecordsToInsert)
				digitalRecordsToInsert = nil
			}

		}
		fmt.Println(device.Name + " finished in " + time.Since(timer).String())
	}

}
