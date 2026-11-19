package usecase

import (
	"changeme/internal/domain"
	"changeme/internal/dto"
	"context"
	"fmt"
	"strings"
)

type win32Client interface {
	EnableDebugPrivilege() error
}

type wmiClient interface {
	GetAllProcesses() ([]dto.ProcessShort, error)
	IsServiceProcess(pid int64) (bool, error)
	GetServiceProcessDisplayName(pid int64) (string, error)
}

type processRepository interface {
	Save(ctx context.Context, process *domain.Process) error
}

type InitialStateCollector struct {
	win32Client       win32Client
	wmiClient         wmiClient
	processRepository processRepository
}

func NewInitialStateCollector(win32Client win32Client, wmiClient wmiClient, processRepository processRepository) *InitialStateCollector {
	return &InitialStateCollector{
		win32Client:       win32Client,
		wmiClient:         wmiClient,
		processRepository: processRepository,
	}
}

func (i InitialStateCollector) CollectProcesses() error {
	// Добавляем привилегию отладки текущему процессу
	err := i.win32Client.EnableDebugPrivilege()
	if err != nil {
		return err
	}

	// Получаем часть информации обо всех процессах в системе
	processesShort, err := i.wmiClient.GetAllProcesses()
	if err != nil {
		return err
	}

	// Обогащаем процесс до полной сущности
	processes := make([]domain.Process, len(processesShort))
	for index, processShort := range processesShort {
		// Определяем тип процесса
		var processType string
		if processShort.PID == 0 || processShort.PID == 4 {
			processType = domain.ProcessTypeCore
		}
		if processShort.ProcessName == "Secure System" || processShort.ProcessName == "Registry" {
			processType = domain.ProcessTypeCore
		}
		if processType == "" {
			lowerPath := strings.ToLower(processShort.Path)
			if strings.Contains(lowerPath, `\program files\windowsapps\`) ||
				strings.Contains(lowerPath, `\windows\systemapps\`) ||
				strings.Contains(lowerPath, `\windows\immersivecontrolpanel\`) ||
				strings.Contains(lowerPath, `\windows\apps\`) {
				processType = domain.ProcessTypeUWP
			}
		}
		if processType == "" {
			isServiceProcess, err := i.wmiClient.IsServiceProcess(processShort.PID)
			if err != nil {
				return err
			}
			if isServiceProcess {
				processType = domain.ProcessTypeService
			}
		}
		if processType == "" {
			processType = domain.ProcessTypeCommon
		}

		// Определяем читаемое имя
		var processDisplayName string
		if processType == domain.ProcessTypeCore {
			switch processShort.PID {
			case 0:
				processDisplayName = "Бездействие системы"
			case 4:
				processDisplayName = "Система"
			}
			switch processShort.ProcessName {
			case "Secure System":
				processDisplayName = "Изолированная среда выполнения ядра"
			case "Registry":
				processDisplayName = "Сжатое хранилище реестра windows"
			}
		}
		if processType == domain.ProcessTypeUWP {
			systemApps := map[string]string{
				"MicrosoftWindows.Client.CBS_cw5n1h2txyewy\\CrossDeviceResume.exe":     "Восстановление на устройствах",
				"MicrosoftWindows.Client.CBS_cw5n1h2txyewy\\FESearchHost.exe":          "Поиск в Проводнике",
				"MicrosoftWindows.Client.CBS_cw5n1h2txyewy\\SearchHost.exe":            "Поиск",
				"MicrosoftWindows.Client.CBS_cw5n1h2txyewy\\TextInputHost.exe":         "Сенсорная клавиатура",
				"MicrosoftWindows.Client.CBS_cw5n1h2txyewy\\VisualAssistExe.exe":       "Визуальная помощь",
				"MicrosoftWindows.Client.CBS_cw5n1h2txyewy\\WebExperienceHostApp.exe":  "Edge WebView2",
				"MicrosoftWindows.Client.CBS_cw5n1h2txyewy\\WindowsBackupClient.exe":   "Резервное копирование",
				"MicrosoftWindows.Client.Core_cw5n1h2txyewy\\ScreenClippingHost.exe":   "Ножницы",
				"MicrosoftWindows.Client.Core_cw5n1h2txyewy\\LogonWebHostProduct.exe":  "Экран входа",
				"MicrosoftWindows.Client.Core_cw5n1h2txyewy\\SnapLauncherProtocol.exe": "Сплит-экран",
			}
			lowerPath := strings.ToLower(processShort.Path)
			for key, name := range systemApps {
				if strings.Contains(lowerPath, strings.ToLower(key)) {
					processDisplayName = name
				}
			}
		}
		//if processType == domain.ProcessTypeService {
		//	processDisplayName, err = i.wmiClient.GetServiceProcessDisplayName(processShort.PID)
		//	if err != nil {
		//		return err
		//	}
		//}
		//if processType == domain.ProcessTypeCommon {
		//
		//}

		// ToDo Определяем статус

		//ToDo Определяем иконку

		processes[index] = domain.Process{
			PID:         processShort.PID,
			ProcessName: processShort.ProcessName,
			DisplayName: processDisplayName,
			Path:        processShort.Path,
			Status:      "",
			Type:        processType,
			Icon:        "",
			CreatedAt:   processShort.CreatedAt,
		}

		fmt.Print(processes[index], "\n")
	}

	return nil
}
