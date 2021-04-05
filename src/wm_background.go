package main

func (host *WmHost) wm_background_map(){
	rattr := host.wm_host_get_window_attributes(host.root_window)

	bw := int(rattr.width)
	bh := int(rattr.height)

	host.wm_host_setup_transparent(&host.background, host.root_window,
								   int(rattr.x),
								   int(rattr.y),
								   bw,
								   bh)

	host.wm_host_map_window(host.background.window)

	host.wm_host_draw_background(bw, bh)
}