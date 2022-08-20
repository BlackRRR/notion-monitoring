package model

import (
	"encoding/json"
	"fmt"
	"os"
)

var AdminSettings *Admin

type Admin struct {
	AdminID map[int64]*AdminUser `json:"admin_id"`
}

type AdminUser struct {
	Name string `json:"name"`
}

func DownloadAdminSettings() {
	var settings *Admin
	data, err := os.ReadFile("./assets/admin.json")
	if err != nil {
		fmt.Println(err)
	}

	err = json.Unmarshal(data, &settings)
	if err != nil {
		fmt.Println(err)
	}

	Validate(settings)

	AdminSettings = settings
	UploadAdminSettings()
}

func UploadAdminSettings() {
	data, err := json.MarshalIndent(AdminSettings, "", "  ")
	if err != nil {
		panic(err)
	}

	if err = os.WriteFile("./assets/admin.json", data, 0600); err != nil {
		panic(err)
	}
}

func Validate(settings *Admin) {
	if settings == nil {
		settings = &Admin{}
	}

	if settings.AdminID == nil {
		settings.AdminID = make(map[int64]*AdminUser)
	}
}
