package filters

const (
	ProgramIdRadiumAmm = "675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8"
)

func FilterInstructionByProgramId(programId string) bool {
	if programId != ProgramIdRadiumAmm {
		return true
	}

	return false
}
