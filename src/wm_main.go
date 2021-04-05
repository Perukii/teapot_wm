
package main

func main(){
	var host_src WmHost
	(&host_src).wm_main_run()
}

func (host *WmHost) wm_main_run(){
	host.wm_host_init()
	host.wm_host_init_log_file()
	defer host.wm_host_close_log_file()

	host.wm_json_read()
	host.wm_json_apply_user_setting()

	host.wm_host_select_input(host.root_window,
		XSubstructureNotifyMask |
		XButtonPressMask |
		XButtonReleaseMask |
		XPointerMotionMask |
		XSubstructureRedirectMask |
		XKeyPressMask |
		XKeyReleaseMask)
	
	host.wm_host_define_cursor(XCLeftPtr)

	host.setting.client_border_overall_width= 25
	host.setting.client_border_shadow_width = 10
	host.setting.client_button_width = 15
	host.setting.client_button_margin_width = 5
	host.setting.client_text_margin_width = 5
	host.setting.max_resize_process_wait = 1.5

	host.wm_background_map()
	host.wm_host_run()
}