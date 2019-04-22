package data

import (
	"sync"
	"time"

	"github.com/oshosanya/go-dm/src/util"
)

type DownloadRoutine struct {
	ID              uint     `gorm:"primary_key"`
	Download        Download `gorm:"foreignkey:DownloadRefer"`
	DownloadID      uint
	Url             string
	FileName        string
	StartRange      int
	EndRange        int
	BytesDownloaded int
	Done            bool
	Status          string `gorm:"default:'ready'"`
	CreatedAt       int64
	UpdatedAt       int64
}

var mu sync.Mutex

func (d *DownloadRoutine) BeforeSave() {
	d.CreatedAt = time.Now().Unix()
}

func (d *DownloadRoutine) BeforeUpdate() {
	d.UpdatedAt = time.Now().Unix()
}

func (d *DownloadRoutine) Update() {
	db := getDB()
	db.Save(&d)
}

func CreateDownloadRoutine(downloadModel *Download, routineDef util.RoutineDefinition) {
	db := getDB()
	downloadRoutine := DownloadRoutine{
		DownloadID:      downloadModel.ID,
		Url:             downloadModel.Url,
		FileName:        routineDef.FileName,
		StartRange:      routineDef.StartRange,
		EndRange:        routineDef.EndRange,
		BytesDownloaded: int(0),
		Done:            false,
		Status:          "ready",
	}

	db.Create(&downloadRoutine)
}

func GetDownloadRoutines(downloadModel *Download) []DownloadRoutine {
	var downloadRoutines []DownloadRoutine
	db := getDB()
	db.Where(&DownloadRoutine{DownloadID: downloadModel.ID, Status: "ready"}).Order("created_at desc").Find(&downloadRoutines)
	return downloadRoutines
}

func GetAllDownloadRoutines(downloadModel *Download) []DownloadRoutine {
	mu.Lock()
	defer mu.Unlock()
	var downloadRoutines []DownloadRoutine
	db := getDB()
	db.Where(&DownloadRoutine{DownloadID: downloadModel.ID}).Order("id asc").Find(&downloadRoutines)
	return downloadRoutines
}

func GetDownloadRoutineById(id uint) DownloadRoutine {
	var downloadRoutine DownloadRoutine
	db := getDB()
	db.Where("id = ?", id).First(&downloadRoutine)
	return downloadRoutine
}

func (d *DownloadRoutine) GetSiblings() []DownloadRoutine {
	db := getDB()
	var downloadRoutine []DownloadRoutine
	db.Model(&d).Related(&downloadRoutine)
	return downloadRoutine
}

func GetLimitedRoutines(limit int) []DownloadRoutine {
	var downloadRoutines []DownloadRoutine
	db := getDB()
	db.Limit(limit).Where(&DownloadRoutine{Status: "ready"}).Order("created_at asc").Find(&downloadRoutines)
	return downloadRoutines
}
