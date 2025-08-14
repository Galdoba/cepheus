package interaction

import (
	"github.com/charmbracelet/huh"
)

type PromptBuilder struct {
	input      *huh.Input
	inputValue string
}

func GetInput(title string, opts ...PromptOption) (string, error) {
	result := ""
	prompt := PromptBuilder{
		input: huh.NewInput().Title(title).Value(&result),
	}
	for _, modify := range opts {
		modify(&prompt)
	}
	if err := prompt.input.Run(); err != nil {
		return "", err
	}
	return result, nil
}

type PromptOption func(*PromptBuilder)

func WithInitialPrompt(initial string) PromptOption {
	return func(pb *PromptBuilder) {
		pb.input = pb.input.Prompt(initial)
	}
}
func WithInitialPlaceholder(initial string) PromptOption {
	return func(pb *PromptBuilder) {
		pb.input = pb.input.Placeholder(initial)
	}
}

func WithValidator(validationFunc func(string) error) PromptOption {
	return func(pb *PromptBuilder) {
		pb.input = pb.input.Validate(validationFunc)
	}
}
func WithDescription(descr string) PromptOption {
	return func(pb *PromptBuilder) {
		pb.input = pb.input.Description(descr)
	}
}
