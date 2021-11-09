package validate

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// 常用tag介绍：
// ne：不等于参数值，例如ne=5；
// gt：大于参数值，例如gt=5；
// gte：大于等于参数值，例如gte=50；
// lt：小于参数值，例如lt=50；
// lte：小于等于参数值，例如lte=50；
// oneof：只能是列举出的值其中一个，这些值必须是数值或字符串，以空格分隔，如果字符串中有空格，将字符串用单引号包围，例如oneof=male female。
// eq：等于参数值，注意与len不同。对于字符串，eq约束字符串本身的值，而len约束字符串长度。例如eq=10；
// len：等于参数值，例如len=10；
// max：小于等于参数值，例如max=10；
// min：大于等于参数值，例如min=10

type Validate struct {
	Validate *validator.Validate
}

func NewValidate() *Validate {
	return &Validate{
		Validate: validator.New(),
	}
}

func (v *Validate) Verify(obj interface{}) error {
	var errMsg []string
	if err := v.Validate.Struct(obj); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			errMsg = append(errMsg, fmt.Sprintf("%s:%s", e.Field(), e.Tag()))
		}
	}
	if len(errMsg) > 0 {
		return fmt.Errorf("%s", strings.Join(errMsg, ","))
	}
	return nil
}
