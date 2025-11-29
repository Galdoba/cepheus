package value

import (
	"fmt"
	"testing"

	"github.com/Galdoba/cepheus/internal/domain/core/values/modifier"
)

func TestNewSkillValue(t *testing.T) {
	sk := NewSkillValue()
	sk.SetBase(3)
	sk.SetModifier(modifier.NewModifier("Injury", "Broken Leg", -2))
	fmt.Println(sk.Value())
}
