package wmi

import (
	"changeme/internal/dto"
	"fmt"
	"time"

	"github.com/yusufpapurcu/wmi"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (c Client) GetAllProcesses() ([]dto.ProcessShort, error) {
	ch := make(chan struct {
		processesShort []dto.ProcessShort
		err            error
	})
	go func() {
		var processes []struct {
			ProcessId      int64
			Name           string
			ExecutablePath string
			CreationDate   time.Time
		}
		err := wmi.Query("SELECT ProcessId, Name, ExecutablePath, CreationDate FROM Win32_Process", &processes)
		if err != nil {
			err = fmt.Errorf("wmi: получение коллекции процессов: %w", err)
		}
		processesShort := make([]dto.ProcessShort, len(processes))
		for index, process := range processes {
			processesShort[index] = dto.ProcessShort{
				PID:         process.ProcessId,
				ProcessName: process.Name,
				Path:        process.ExecutablePath,
				CreatedAt:   process.CreationDate,
			}
		}
		ch <- struct {
			processesShort []dto.ProcessShort
			err            error
		}{
			processesShort: processesShort,
			err:            err,
		}
	}()
	result := <-ch
	if result.err != nil {
		return nil, result.err
	}

	return result.processesShort, nil
}

func (c Client) IsServiceProcess(pid int64) (bool, error) {
	ch := make(chan struct {
		isService bool
		err       error
	})
	go func() {
		var serviceID []struct {
			ProcessId int64
		}
		query := fmt.Sprintf("SELECT ProcessId FROM Win32_Service WHERE ProcessId = %d", pid)
		err := wmi.Query(query, &serviceID)
		ch <- struct {
			isService bool
			err       error
		}{
			isService: len(serviceID) > 0,
			err:       err,
		}
	}()
	result := <-ch
	if result.err != nil {
		return false, fmt.Errorf("wmi: определение службы: %w", result.err)
	}

	return result.isService, nil
}

func (c Client) GetServiceProcessDisplayName(pid int64) (string, error) {

	return "", nil
}
