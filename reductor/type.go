package reductor

import (
	"fmt"
	"pdfimporter/domain"
	"sync"

	"go.uber.org/zap"
)

type ModelList map[domain.Model]domain.Modeler

type Reductor struct {
	mutex        sync.RWMutex
	loger        *zap.SugaredLogger
	models       ModelList
	outStateChan chan domain.Model
}

var once sync.Once
var instance *Reductor

// создаем singleton без начальной модели
func New(logger *zap.SugaredLogger) error {
	if logger == nil {
		return fmt.Errorf("reductor: logger is nil")
	}
	once.Do(func() {
		instance = &Reductor{
			loger:  logger,
			models: make(ModelList),
		}
	})
	return nil
}

func Instance() *Reductor {
	if instance == nil {
		panic("reductor instance is nil")
	}
	return instance
}

// если ли в запомненных моделях данная
func (rdc *Reductor) IsExistModel(model domain.Model) bool {
	rdc.mutex.RLock()
	defer rdc.mutex.RUnlock()
	if _, ok := rdc.models[model]; ok {
		return true
	}
	return false
}

// прописываем канал по которому будем ждать уведомления об обновлении модели
func (rdc *Reductor) SetOutChanState(out chan domain.Model) {
	rdc.mutex.Lock()
	defer rdc.mutex.Unlock()
	rdc.outStateChan = out
}
