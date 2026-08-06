//go:build windows

package watcher

import (
	"fmt"
	"runtime"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	WM_CREATE            = 0x0001
	WM_DESTROY           = 0x0002
	WM_CLOSE             = 0x0010
	WM_QUERYENDSESSION   = 0x0011
	WM_ENDSESSION        = 0x0016
	WM_POWERBROADCAST    = 0x0218
	WM_WTSSESSION_CHANGE = 0x02B1

	PBT_APMSUSPEND         = 0x0004
	PBT_APMRESUMESUSPEND   = 0x0007
	PBT_APMRESUMEAUTOMATIC = 0x0012

	WTS_CONSOLE_CONNECT    = 0x1
	WTS_CONSOLE_DISCONNECT = 0x2
	WTS_SESSION_LOCK       = 0x7
	WTS_SESSION_UNLOCK     = 0x8

	NOTIFY_FOR_THIS_SESSION = 0
)

var (
	user32   = windows.NewLazySystemDLL("user32.dll")
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")
	wtsapi32 = windows.NewLazySystemDLL("wtsapi32.dll")

	procRegisterClassExW = user32.NewProc("RegisterClassExW")
	procCreateWindowExW  = user32.NewProc("CreateWindowExW")
	procDefWindowProcW   = user32.NewProc("DefWindowProcW")
	procGetMessageW      = user32.NewProc("GetMessageW")
	procTranslateMessage = user32.NewProc("TranslateMessage")
	procDispatchMessageW = user32.NewProc("DispatchMessageW")
	procPostQuitMessage  = user32.NewProc("PostQuitMessage")
	procDestroyWindow    = user32.NewProc("DestroyWindow")

	procGetModuleHandleW = kernel32.NewProc("GetModuleHandleW")

	procWTSRegisterSessionNotification   = wtsapi32.NewProc("WTSRegisterSessionNotification")
	procWTSUnregisterSessionNotification = wtsapi32.NewProc("WTSUnregisterSessionNotification")
)

type WNDCLASSEXW struct {
	Size       uint32
	Style      uint32
	WndProc    uintptr
	ClsExtra   int32
	WndExtra   int32
	Instance   windows.Handle
	Icon       windows.Handle
	Cursor     windows.Handle
	Background windows.Handle
	MenuName   *uint16
	ClassName  *uint16
	IconSm     windows.Handle
}

type MSG struct {
	Hwnd    windows.Handle
	Message uint32
	WParam  uintptr
	LParam  uintptr
	Time    uint32
	Pt      struct{ X, Y int32 }
}

type WindowsWatcher struct {
	callbacks EventCallbacks
	hwnd      windows.Handle
	mu        sync.Mutex
	lastWake  time.Time
}

var globalWatcher *WindowsWatcher

func New(callbacks EventCallbacks) Watcher {
	w := &WindowsWatcher{
		callbacks: callbacks,
	}
	globalWatcher = w
	return w
}

func (w *WindowsWatcher) Start() error {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	instance, _, _ := procGetModuleHandleW.Call(0)
	if instance == 0 {
		return fmt.Errorf("failed to get module handle")
	}

	className, err := windows.UTF16PtrFromString("PingSystemMessageWindowClass")
	if err != nil {
		return err
	}

	windowName, err := windows.UTF16PtrFromString("PingSystemMessageWindow")
	if err != nil {
		return err
	}

	wndProcCallback := syscall.NewCallback(wndProc)

	wc := WNDCLASSEXW{
		Size:      uint32(unsafe.Sizeof(WNDCLASSEXW{})),
		WndProc:   wndProcCallback,
		Instance:  windows.Handle(instance),
		ClassName: className,
	}

	procRegisterClassExW.Call(uintptr(unsafe.Pointer(&wc)))

	// Create a top-level hidden window (hWndParent = 0) so broadcast messages
	// like WM_POWERBROADCAST and WM_WTSSESSION_CHANGE are delivered properly.
	hwnd, _, _ := procCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowName)),
		0, 0, 0, 0, 0,
		0,
		0,
		instance,
		0,
	)

	if hwnd == 0 {
		return fmt.Errorf("failed to create Win32 message window")
	}

	w.hwnd = windows.Handle(hwnd)

	// Register for Windows Session Lock/Unlock notifications
	procWTSRegisterSessionNotification.Call(uintptr(w.hwnd), NOTIFY_FOR_THIS_SESSION)

	var msg MSG
	for {
		res, _, _ := procGetMessageW.Call(
			uintptr(unsafe.Pointer(&msg)),
			0, 0, 0,
		)

		if int32(res) <= 0 {
			break
		}

		procTranslateMessage.Call(uintptr(unsafe.Pointer(&msg)))
		procDispatchMessageW.Call(uintptr(unsafe.Pointer(&msg)))
	}

	return nil
}

func (w *WindowsWatcher) Stop() {
	if w.hwnd != 0 {
		procWTSUnregisterSessionNotification.Call(uintptr(w.hwnd))
		procPostQuitMessage.Call(0)
		procDestroyWindow.Call(uintptr(w.hwnd))
	}
}

func wndProc(hwnd windows.Handle, msg uint32, wParam, lParam uintptr) uintptr {
	if globalWatcher == nil {
		ret, _, _ := procDefWindowProcW.Call(uintptr(hwnd), uintptr(msg), wParam, lParam)
		return ret
	}

	switch msg {
	case WM_POWERBROADCAST:
		switch wParam {
		case PBT_APMSUSPEND:
			if globalWatcher.callbacks.OnSleep != nil {
				globalWatcher.callbacks.OnSleep()
			}
		case PBT_APMRESUMEAUTOMATIC, PBT_APMRESUMESUSPEND:
			globalWatcher.mu.Lock()
			now := time.Now()
			if now.Sub(globalWatcher.lastWake) > 2*time.Second {
				globalWatcher.lastWake = now
				globalWatcher.mu.Unlock()
				if globalWatcher.callbacks.OnWake != nil {
					globalWatcher.callbacks.OnWake()
				}
			} else {
				globalWatcher.mu.Unlock()
			}
		}
		return 1

	case WM_WTSSESSION_CHANGE:
		switch wParam {
		case WTS_SESSION_LOCK:
			if globalWatcher.callbacks.OnLock != nil {
				globalWatcher.callbacks.OnLock()
			}
		case WTS_SESSION_UNLOCK:
			if globalWatcher.callbacks.OnUnlock != nil {
				globalWatcher.callbacks.OnUnlock()
			}
		}
		return 0

	case WM_QUERYENDSESSION:
		if globalWatcher.callbacks.OnShutdown != nil {
			globalWatcher.callbacks.OnShutdown()
		}
		return 1 // Allow session end

	case WM_ENDSESSION:
		if wParam != 0 {
			if globalWatcher.callbacks.OnShutdown != nil {
				globalWatcher.callbacks.OnShutdown()
			}
		}
		return 0

	case WM_DESTROY:
		procPostQuitMessage.Call(0)
		return 0
	}

	ret, _, _ := procDefWindowProcW.Call(uintptr(hwnd), uintptr(msg), wParam, lParam)
	return ret
}
