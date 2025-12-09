package skill_test

import (
	"fmt"
	"testing"

	"github.com/Galdoba/cepheus/internal/domain/core/values/skill"
)

func TestFromDescription(t *testing.T) {
	sk, err := skill.FromDescription("must use Admin or Broker Skill for this roll")
	fmt.Println(sk, err)
}
