package main

import (
	"os"
	"encoding/json"
	"io/ioutil"
)


func (host *WmHost) wm_json_read(){
	if len(os.Args) <= 1 { return }

	wm_debug_log(os.Args[1])
	data, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		wm_debug_log("Failed to load Setting file.")
	} else {
		err = json.Unmarshal(data, &host.json_setting)
		if err != nil{
			wm_debug_error("Failed to parse Setting file.")
		}
	}
	//wm_debug_log(host.json_setting.UserSetting.BackgroundPngFilePath)
}

func (host *WmHost) wm_json_apply_user_setting(){
	host.setting.background_file = host.json_setting.UserSetting.BackgroundPngFilePath
}