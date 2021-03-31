package main

func (host *WmHost) wm_host_init(){
	host.display = wm_x11_open_display()
	host.root_window = wm_x11_get_root(host.display)
	host.client = []WmClient{}
	var clt WmClient
	host.client = append(host.client, clt)
}

func (host *WmHost) wm_host_init_log_file(){
	wm_debug_log_file_init(host.log_file)
}

func (host *WmHost) wm_host_close_log_file(){
	host.log_file.Close()
}

func (host *WmHost) wm_host_grab_button(window XWindowID){
	wm_x11_grab_button(host.display, window)
}

func (host *WmHost) wm_host_raise_window(window XWindowID){
	wm_x11_raise_window(host.display, window)
}

func (host *WmHost) wm_host_lower_window(window XWindowID){
	wm_x11_lower_window(host.display, window)
}

func (host *WmHost) wm_host_get_window_attributes(window XWindowID) XWindowAttributes{
	return wm_x11_get_window_attributes(host.display, window)
}

func (host *WmHost) wm_host_move_window(window XWindowID, x int, y int){
	wm_x11_move_window(host.display, window, x, y)
}

func (host *WmHost) wm_host_select_input(window XWindowID, mask CLong){
	wm_x11_select_input(host.display, window, mask)
}

func (host *WmHost) wm_host_resize_window(window XWindowID, w int, h int){
	wm_x11_resize_window(host.display, window, w, h)
}

func (host *WmHost) wm_host_resize_surface(surface *CairoSfc, w int, h int){
	wm_x11_resize_surface(surface, w, h)
}

func (host *WmHost) wm_host_resize_transparent(transparent WmTransparent, w int, h int){
	wm_x11_resize_window(host.display, transparent.window, w, h)
	wm_x11_resize_surface(transparent.surface, w, h)
}

func (host *WmHost) wm_host_setup_transparent(transparent *WmTransparent, parent XWindowID,
											  x int, y int, w int, h int,
											  ){
	wm_x11_create_transparent_window(host.display, parent, x, y, w, h, transparent)
}

func (host *WmHost) wm_host_remove_transparent(transparent WmTransparent){
	wm_x11_destroy_window(host.display, transparent.window)
	wm_x11_destroy_cairo_surface(host.display, transparent.surface)
}

func (host *WmHost) wm_host_map_window(window XWindowID){
	wm_x11_map_window(host.display, window)
}

func (host *WmHost) wm_host_unmap_window(window XWindowID){
	wm_x11_unmap_window(host.display, window)
}

func (host *WmHost) wm_host_draw_transparent(transparent WmTransparent){
	wm_x11_draw_transparent(host.display, transparent,
							host.config.client_drawable_range_border_width,
							host.config.client_grab_area_resize_width)
}

func (host *WmHost) wm_host_reparent_window(window XWindowID, parent XWindowID, x int, y int){
	wm_x11_reparent_window(host.display, window, parent ,x, y)
}

func (host *WmHost) wm_host_define_cursor(cursor int){
    wm_x11_define_cursor(host.display, host.root_window, cursor)
}

func (host *WmHost) wm_host_query_parent(window XWindowID) XWindowID{
	relation := wm_x11_query_tree(host.display, window)
	return relation.parent
}

func (host *WmHost) wm_host_check_n_of_queued_event() int{
	return wm_x11_check_n_of_queued_event(host.display)
}

func (host *WmHost) wm_host_set_focus_to_client(address WmClientAddress){
	clt := host.client[address]
	for i := address; i < len(host.client)-1; i++{
		host.client[i] = host.client[i+1]
	}
	host.client[len(host.client)-1] = clt
	host.wm_host_update_client_focus()
}

func (host *WmHost) wm_host_update_client_focus(){
	clt := host.client[len(host.client)-1]
	wm_x11_set_input_focus(host.display, clt.app)

	for i := 1; i < len(host.client)-1; i++{
		host.wm_client_raise_mask(i)
	}
	host.wm_client_raise_app(len(host.client)-1)
}

func (host *WmHost) wm_host_update_grab_mode(window XWindowID, bx int, by int){

	resize_area_width := host.config.client_grab_area_resize_width

	attr := host.wm_host_get_window_attributes(window)

	grab_rx := bx-int(attr.x)
	grab_ry := by-int(attr.y)
	
	host.grab_mode_1 = WM_RESIZE_MODE_NONE

	if grab_rx < resize_area_width  {
		host.grab_mode_1 = WM_RESIZE_MODE_LEFT
	}
	if grab_rx > int(attr.width)-resize_area_width {
		host.grab_mode_1 = WM_RESIZE_MODE_RIGHT
	}

	host.grab_mode_2 = WM_RESIZE_MODE_NONE

	if grab_ry < resize_area_width  {
		host.grab_mode_2 = WM_RESIZE_MODE_TOP
	}
	if grab_ry > int(attr.height)-resize_area_width {
		host.grab_mode_2 = WM_RESIZE_MODE_BOTTOM
	}
}

func (host *WmHost) wm_host_run(){
	for{
		host.event = wm_x11_peek_event(host.display)
		switch host.event.wm_event_get_type(){
		case XMapNotify:
			host.wm_event_loop_map_notify()
		case XUnmapNotify:
			host.wm_event_loop_unmap_notify()
		case XMapRequest:
			host.wm_event_loop_map_request()
		case XDestroyNotify:
			host.wm_event_loop_destroy_notify()
		case XConfigureNotify:
			host.wm_event_loop_configure_notify()
		case XConfigureRequest:
			host.wm_event_loop_configure_request()
		case XKeyPress:
			host.wm_event_loop_key_press()
		case XButtonPress:
			host.wm_event_loop_button_press()
		case XButtonRelease:
			host.wm_event_loop_button_release()
		case XMotionNotify:
			host.wm_event_loop_motion_notify()
		}
	}
}


