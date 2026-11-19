package win32

import (
	"fmt"

	"golang.org/x/sys/windows"
)

type Client struct {
	kernel32 *windows.LazyDLL
	shlwapi  *windows.LazyDLL
}

func NewClient() *Client {
	return &Client{
		kernel32: windows.NewLazySystemDLL("kernel32.dll"),
		shlwapi:  windows.NewLazySystemDLL("shlwapi.dll"),
	}
}

func (c Client) EnableDebugPrivilege() error {
	name, err := windows.UTF16PtrFromString("SeDebugPrivilege")
	if err != nil {
		return fmt.Errorf("win32: кодирование имени привилегии: %w", err)
	}

	var luid windows.LUID
	if err := windows.LookupPrivilegeValue(nil, name, &luid); err != nil {
		return fmt.Errorf("win32: поиск SeDebugPrivilege: %w", err)
	}

	var token windows.Token
	if err := windows.OpenProcessToken(
		windows.CurrentProcess(),
		windows.TOKEN_ADJUST_PRIVILEGES|windows.TOKEN_QUERY,
		&token,
	); err != nil {
		return fmt.Errorf("win32: открытие токена процесса: %w", err)
	}
	defer func() { _ = windows.CloseHandle(windows.Handle(token)) }()

	privileges := &windows.Tokenprivileges{
		PrivilegeCount: 1,
		Privileges: [1]windows.LUIDAndAttributes{
			{Luid: luid, Attributes: windows.SE_PRIVILEGE_ENABLED},
		},
	}

	if err = windows.AdjustTokenPrivileges(token, false, privileges, 0, nil, nil); err != nil {
		return fmt.Errorf("win32: добавление привилегии в токен: %w", err)
	}

	if err = windows.GetLastError(); err != nil {
		if errno, ok := err.(windows.Errno); ok {
			switch errno {
			case windows.ERROR_PRIVILEGE_NOT_HELD:
				return fmt.Errorf("win32: пользователю не выдана привилегия SeDebugPrivilege (нужно право «Debug programs»): %w", err)
			case windows.ERROR_NOT_ALL_ASSIGNED:
				return fmt.Errorf("win32: не все привилегии в списке были включены: %w", err)
			}
		}
		return fmt.Errorf("win32: включение SeDebugPrivilege: %w", err)
	}

	return nil
}
