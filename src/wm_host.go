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
/*
func (host *WmHost) wm_host_resize_transparent(transparent WmTransparent, w int, h int){
	wm_x11_resize_window(host.display, transparent.window, w, h)
	wm_x11_resize_surface(transparent.surface, w, h)
}
*/

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
	wm_x11_draw_transparent(host.display, transparent, host.config, 0, "")
}

func (host *WmHost) wm_host_draw_client(address WmClientAddress){
	clt := host.client[address]
	wm_x11_draw_transparent(host.display, clt.mask, host.config, host.mask_button, clt.title)
}

func (host *WmHost) wm_host_reparent_window(window XWindowID, parent XWindowID, x int, y int){
	wm_x11_reparent_window(host.display, window, parent ,x, y)
}

func (host *WmHost) wm_host_define_cursor(cursor int){
	if host.cursor == cursor { return }
	host.cursor = cursor
    wm_x11_define_cursor(host.display, host.root_window, cursor)
}

func (host *WmHost) wm_host_query_parent(window XWindowID) XWindowID{
	relation := wm_x11_query_tree(host.display, window)
	return relation.parent
}

func (host *WmHost) wm_host_send_delete_event(window XWindowID){
	wm_x11_send_delete_event(host.display, window)
}

func (host *WmHost) wm_host_check_n_of_queued_event() int{
	return wm_x11_check_n_of_queued_event(host.display)
}

func (host *WmHost) wm_host_set_focus_to_client(address WmClientAddress){

	host.wm_client_raise_mask(len(host.client)-1)

	clt := host.client[address]
	for i := address; i < len(host.client)-1; i++{
		host.client[i] = host.client[i+1]
	}
	host.client[len(host.client)-1] = clt

	host.wm_client_raise_app(len(host.client)-1)

}

func (host *WmHost) wm_host_restack_clients(){
	for i := 1; i < len(host.client)-1; i++{
		host.wm_client_raise_mask(i)
	}
	host.wm_client_raise_app(len(host.client)-1)
}

func (host *WmHost) wm_host_get_size_hints(window XWindowID) XSizeHints{
	return wm_x11_get_size_hints(host.display, window)
}

func (host *WmHost) wm_host_get_window_title(window XWindowID) string{
	return wm_x11_get_window_title(host.display, window)
}

func (host *WmHost) wm_host_intern_atom(name string) XAtom{
	return wm_x11_intern_atom(host.display, name)
}

func (host *WmHost) wm_host_update_grab_mode(point_x int, point_y int, mask_x int, mask_y int, mask_w int, mask_h int){

	resize_area_width := host.config.client_grab_area_resize_width

	grab_rx := point_x-mask_x
	grab_ry := point_y-mask_y
	
	host.grab_mode_1 = WM_RESIZE_MODE_NONE

	if grab_rx < resize_area_width  {
		host.grab_mode_1 = WM_RESIZE_MODE_LEFT
	}
	if grab_rx > mask_w-resize_area_width {
		host.grab_mode_1 = WM_RESIZE_MODE_RIGHT
	}

	host.grab_mode_2 = WM_RESIZE_MODE_NONE

	if grab_ry < resize_area_width  {
		host.grab_mode_2 = WM_RESIZE_MODE_TOP
	}
	if grab_ry > mask_h-resize_area_width {
		host.grab_mode_2 = WM_RESIZE_MODE_BOTTOM
	}
}

func (host *WmHost) wm_host_update_button_mode(point_x int, point_y int, mask_x int, mask_y int, mask_w int, mask_h int){
	var in_rect = func(px, py, x, y, w, h int) bool{
		return (px >= x) && (px <= x+w) && (py >= y) && (py <= y+h)
	}

	host.mask_button = WM_BUTTON_NONE

	border_width := host.config.client_drawable_range_border_width
	button_width := host.config.client_button_width
	button_margin_width := host.config.client_button_margin_width

	for i := 1; i <= 3; i++ {
        
        button_x := mask_w - border_width - (button_width + button_margin_width)*i;
		button_y := border_width - button_width - button_margin_width;

		if in_rect(point_x-mask_x, point_y-mask_y, button_x, button_y, button_width, button_width){
			switch i {
			case 1:
				host.mask_button = WM_BUTTON_EXIT
			case 2:
				host.mask_button = WM_BUTTON_MINIMIZE
			case 3:
				host.mask_button = WM_BUTTON_MAXIMIZE
			}
			break
		}
	}
}

func (host *WmHost) wm_host_update_cursor(){
	if host.grab_mode_1 == WM_RESIZE_MODE_NONE && host.grab_mode_2 == WM_RESIZE_MODE_NONE {
		host.wm_host_define_cursor(XCLeftPtr)
	}
	mode_l := host.grab_mode_1 == WM_RESIZE_MODE_LEFT
	mode_r := host.grab_mode_1 == WM_RESIZE_MODE_RIGHT
	mode_t := host.grab_mode_2 == WM_RESIZE_MODE_TOP
	mode_b := host.grab_mode_2 == WM_RESIZE_MODE_BOTTOM

	if mode_l && mode_t { host.wm_host_define_cursor(XCSideTL);return }
	if mode_r && mode_t { host.wm_host_define_cursor(XCSideTR);return }
	if mode_l && mode_b { host.wm_host_define_cursor(XCSideBL);return }
	if mode_r && mode_b { host.wm_host_define_cursor(XCSideBR);return }

	if mode_t { host.wm_host_define_cursor(XCSideT);return } 
	if mode_b { host.wm_host_define_cursor(XCSideB);return } 
	if mode_l { host.wm_host_define_cursor(XCSideL);return } 
	if mode_r { host.wm_host_define_cursor(XCSideR);return } 
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
		case XEnterNotify:
			host.wm_event_loop_enter_notify()
		case XLeaveNotify:
			host.wm_event_loop_leave_notify()
		case XPropertyNotify:
			host.wm_event_loop_property_notify()
		}
	}
}


