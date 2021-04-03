
package main

func wm_run(host *WmHost){
	host.wm_host_init()

	host.wm_host_init_log_file()
	defer host.wm_host_close_log_file()

	host.wm_host_select_input(host.root_window,
		XSubstructureNotifyMask |
		XButtonPressMask |
		XButtonReleaseMask |
		XPointerMotionMask |
		XSubstructureRedirectMask |
		XKeyPressMask |
		XKeyReleaseMask)
	
	host.wm_host_define_cursor(XCLeftPtr)

	host.config.client_drawable_range_border_width = 25
	host.config.client_grab_area_resize_width = 10
	host.config.client_button_width = 15
	host.config.client_button_margin_width = 5
	host.config.max_config_wait = 1.5

	host.wm_host_run()
}