package main


/*
#cgo pkg-config: x11
#include "./c_wm_x_access.h"
*/
import "C"

type (
	XEvent = C.XEvent
	XDisplay = C.Display
	XWindowID = C.Window
	XWindowAttributes = C.XWindowAttributes
	XKeyEvent = C.XKeyEvent
	XButtonEvent = C.XButtonEvent
	XMotionEvent = C.XMotionEvent
)

const (
	XNone = C.None
	XKeyPress = int(C.KeyPress)
	XButtonPress = int(C.ButtonPress)
	XButtonRelease = int(C.ButtonRelease)
	XMotionNotify = int(C.MotionNotify)
)

func wm_x11_open_display() *XDisplay{
	return C.XOpenDisplay(nil)
}

func wm_x11_get_root(display *XDisplay) XWindowID{
	return C.XDefaultRootWindow(display)
}

func wm_x11_peek_event(display *XDisplay) XEvent{
	var event XEvent
	C.XNextEvent(display, &event)
	return event
}

func wm_x11_grab_button(display *XDisplay, window XWindowID){
    C.XGrabButton(display, C.AnyButton, C.Mod1Mask, window, C.True,
		C.ButtonPressMask|C.ButtonReleaseMask|C.PointerMotionMask,
		C.GrabModeAsync, C.GrabModeAsync, C.None, C.None)
}

func wm_x11_get_type_of_event(event *XEvent) int {
	return int(C.c_wm_x11_get_type_of_event(*event))
}

func wm_x11_raise_window(display *XDisplay, window XWindowID){
	C.XRaiseWindow(display, window)
}

func wm_x11_get_window_attributes(display *XDisplay, window XWindowID) XWindowAttributes{
	var attr XWindowAttributes
	C.XGetWindowAttributes(display, window, &attr)
	return attr
}

func wm_x11_move_window(display *XDisplay, window XWindowID, x int, y int){
	C.XMoveWindow(display, window, C.int(x), C.int(y))
}

func wm_x11_resize_window(display *XDisplay, window XWindowID, w int, h int){
	if w < 1 { w = 1 }
	if h < 1 { h = 1 }
	C.XResizeWindow(display, window, C.uint(w), C.uint(h))
}