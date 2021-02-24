package cmd

// InsprOptions values for the cli
var InsprOptions struct {
	SampleFlagValue string
	// Scope recieves the dApp scope/context from the cli
	Scope string
	// AppliedFileStructure recieves the folder/file to be applied
	AppliedFileStructure string
	// AppliedFolderStructure recieves the folder/file to be applied
	AppliedFolderStructure string
	// DryRun defines if given command will be a dry run or not
	DryRun bool
	// Update defines if Apply is going to create a new app or update an existing one
	Update bool
}
