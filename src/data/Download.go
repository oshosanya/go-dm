package data

import (
	"time"
)

type Download struct {
	ID               uint              `json:id,gorm:"primary_key;AUTO_INCREMENT;unique_index"`
	DownloadRoutines []DownloadRoutine `gorm:"foreignkey:DownloadRefer"`
	Url              string            `json:url`
	FileName         string            `json:file_name`
	ContentLength    int               `json:content_length`
	BytesDownloaded  int               `json:bytes_downloaded`
	Done             bool              `json:done`
	CreatedAt        int64             `json:created_at`
	UpdatedAt        int64             `json:updated_at`
}

func (d *Download) BeforeSave() {
	d.CreatedAt = time.Now().Unix()
}

func (d *Download) BeforeUpdate() {
	d.UpdatedAt = time.Now().Unix()
}

func (d *Download) Update() {
	db := getDB()
	db.Save(&d)
}

func CreateDownload(url string, fileName string, contentLength int) Download {
	db := getDB()
	downloadModel := Download{
		Url:             url,
		FileName:        fileName,
		ContentLength:   contentLength,
		BytesDownloaded: int(0),
		Done:            false,
	}

	db.Create(&downloadModel)

	return downloadModel
}

func GetIncompleteDownloads() []Download {
	var downloads []Download
	db := getDB()
	db.Where("done = ?", 0).Order("created_at desc").Find(&downloads)
	return downloads
}

func GetDownloadByID(id uint) Download {
	var download Download
	db := getDB()
	db.Where("id = ?", id).First(&download)
	return download
}

func GetAllDownloads() []Download {
	db := getDB()
	var downloads []Download
	db.Find(&downloads)
	return downloads
}
