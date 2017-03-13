package main

const (
	DefaultMainService   = "app"
	DefaultVolumeService = "volume"
)

type Config struct {
	ProjectName string `json:"project_name" valid:"required"`
	MainService string `json:"main_service"`
	DataService string `json:"data_service"`
}
