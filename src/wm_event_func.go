package main

import "unsafe"

func (event *XEvent) wm_event_get_type() int {
	return wm_x11_get_type_of_event(event)
}

func (event *XEvent) wm_event_get_pointer() unsafe.Pointer{
	return unsafe.Pointer(event)
}
