package environment

import "strings"

// IsInInputChannel - checks if a word exists in the insprEnv.InputChannels
func (insprEnv *InsprEnvVars) IsInInputChannel(word string, separator string) bool {
	channels := strings.Split(insprEnv.InputChannels, separator)
	for _, c := range channels {
		if word == c {
			return true
		}
	}
	return false
}

// IsInOutputChannel - checks if a word exists in the insprEnv.OutputChannels
func (insprEnv *InsprEnvVars) IsInOutputChannel(word string, separator string) bool {
	channels := strings.Split(insprEnv.OutputChannels, separator)
	for _, c := range channels {
		if word == c {
			return true
		}
	}
	return false
}
