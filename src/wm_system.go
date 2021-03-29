
package main

func wm_run(host *WmHost){
	host.wm_host_init_x11_element()
	host.wm_host_init_log_file()
	defer host.wm_host_close_log_file()
	host.wm_host_grab_button(host.root_window)



	host.wm_host_run()
}