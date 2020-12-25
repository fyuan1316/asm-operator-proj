package task

import (
	"github.com/fyuan1316/asm-operator/pkg/oprlib/manage"
	"github.com/fyuan1316/asm-operator/pkg/task/mock"
)

func GetDeployStages() [][]manage.ExecuteItem {
	return mock.GetDeployStages()
}
