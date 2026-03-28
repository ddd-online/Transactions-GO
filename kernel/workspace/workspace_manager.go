package workspace

import (
	"sync"

	"github.com/sirupsen/logrus"
)

const (
	ErrOpenedWorkspaceNotFound = "未打开工作空间"
)

var Manager = &WsManager{}

type WsManager struct {
	workspace *Workspace
	lock      sync.Mutex
}

func (wm *WsManager) OpenWorkspace(directory string) error {
	wm.lock.Lock()
	defer wm.lock.Unlock()

	ws, err := NewWorkspace(directory)
	if err != nil {
		logrus.Infof("打开工作空间失败 %v", err)
		return err
	}

	wm.workspace = ws
	return nil
}

func (wm *WsManager) OpenedWorkspace() *Workspace {
	wm.lock.Lock()
	defer wm.lock.Unlock()

	return wm.workspace
}

func (wm *WsManager) Close() {
	wm.lock.Lock()
	defer wm.lock.Unlock()

	if wm.workspace != nil {
		wm.workspace.Close()
	}
}
