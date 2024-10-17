// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package raydium_amm

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// MonitorStep is the `monitorStep` instruction.
type MonitorStep struct {
	PlanOrderLimit   *uint16
	PlaceOrderLimit  *uint16
	CancelOrderLimit *uint16

	// [0] = [] tokenProgram
	//
	// [1] = [] rent
	//
	// [2] = [] clock
	//
	// [3] = [WRITE] amm
	//
	// [4] = [] ammAuthority
	//
	// [5] = [WRITE] ammOpenOrders
	//
	// [6] = [WRITE] ammTargetOrders
	//
	// [7] = [WRITE] poolCoinTokenAccount
	//
	// [8] = [WRITE] poolPcTokenAccount
	//
	// [9] = [WRITE] poolWithdrawQueue
	//
	// [10] = [] serumProgram
	//
	// [11] = [WRITE] serumMarket
	//
	// [12] = [WRITE] serumCoinVaultAccount
	//
	// [13] = [WRITE] serumPcVaultAccount
	//
	// [14] = [] serumVaultSigner
	//
	// [15] = [WRITE] serumReqQ
	//
	// [16] = [WRITE] serumEventQ
	//
	// [17] = [WRITE] serumBids
	//
	// [18] = [WRITE] serumAsks
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewMonitorStepInstructionBuilder creates a new `MonitorStep` instruction builder.
func NewMonitorStepInstructionBuilder() *MonitorStep {
	nd := &MonitorStep{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 19),
	}
	return nd
}

// SetPlanOrderLimit sets the "planOrderLimit" parameter.
func (inst *MonitorStep) SetPlanOrderLimit(planOrderLimit uint16) *MonitorStep {
	inst.PlanOrderLimit = &planOrderLimit
	return inst
}

// SetPlaceOrderLimit sets the "placeOrderLimit" parameter.
func (inst *MonitorStep) SetPlaceOrderLimit(placeOrderLimit uint16) *MonitorStep {
	inst.PlaceOrderLimit = &placeOrderLimit
	return inst
}

// SetCancelOrderLimit sets the "cancelOrderLimit" parameter.
func (inst *MonitorStep) SetCancelOrderLimit(cancelOrderLimit uint16) *MonitorStep {
	inst.CancelOrderLimit = &cancelOrderLimit
	return inst
}

// SetTokenProgramAccount sets the "tokenProgram" account.
func (inst *MonitorStep) SetTokenProgramAccount(tokenProgram ag_solanago.PublicKey) *MonitorStep {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(tokenProgram)
	return inst
}

// GetTokenProgramAccount gets the "tokenProgram" account.
func (inst *MonitorStep) GetTokenProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetRentAccount sets the "rent" account.
func (inst *MonitorStep) SetRentAccount(rent ag_solanago.PublicKey) *MonitorStep {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(rent)
	return inst
}

// GetRentAccount gets the "rent" account.
func (inst *MonitorStep) GetRentAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetClockAccount sets the "clock" account.
func (inst *MonitorStep) SetClockAccount(clock ag_solanago.PublicKey) *MonitorStep {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(clock)
	return inst
}

// GetClockAccount gets the "clock" account.
func (inst *MonitorStep) GetClockAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetAmmAccount sets the "amm" account.
func (inst *MonitorStep) SetAmmAccount(amm ag_solanago.PublicKey) *MonitorStep {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(amm).WRITE()
	return inst
}

// GetAmmAccount gets the "amm" account.
func (inst *MonitorStep) GetAmmAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

// SetAmmAuthorityAccount sets the "ammAuthority" account.
func (inst *MonitorStep) SetAmmAuthorityAccount(ammAuthority ag_solanago.PublicKey) *MonitorStep {
	inst.AccountMetaSlice[4] = ag_solanago.Meta(ammAuthority)
	return inst
}

// GetAmmAuthorityAccount gets the "ammAuthority" account.
func (inst *MonitorStep) GetAmmAuthorityAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(4)
}

// SetAmmOpenOrdersAccount sets the "ammOpenOrders" account.
func (inst *MonitorStep) SetAmmOpenOrdersAccount(ammOpenOrders ag_solanago.PublicKey) *MonitorStep {
	inst.AccountMetaSlice[5] = ag_solanago.Meta(ammOpenOrders).WRITE()
	return inst
}

// GetAmmOpenOrdersAccount gets the "ammOpenOrders" account.
func (inst *MonitorStep) GetAmmOpenOrdersAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(5)
}

// SetAmmTargetOrdersAccount sets the "ammTargetOrders" account.
func (inst *MonitorStep) SetAmmTargetOrdersAccount(ammTargetOrders ag_solanago.PublicKey) *MonitorStep {
	inst.AccountMetaSlice[6] = ag_solanago.Meta(ammTargetOrders).WRITE()
	return inst
}

// GetAmmTargetOrdersAccount gets the "ammTargetOrders" account.
func (inst *MonitorStep) GetAmmTargetOrdersAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(6)
}

// SetPoolCoinTokenAccountAccount sets the "poolCoinTokenAccount" account.
func (inst *MonitorStep) SetPoolCoinTokenAccountAccount(poolCoinTokenAccount ag_solanago.PublicKey) *MonitorStep {
	inst.AccountMetaSlice[7] = ag_solanago.Meta(poolCoinTokenAccount).WRITE()
	return inst
}

// GetPoolCoinTokenAccountAccount gets the "poolCoinTokenAccount" account.
func (inst *MonitorStep) GetPoolCoinTokenAccountAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(7)
}

// SetPoolPcTokenAccountAccount sets the "poolPcTokenAccount" account.
func (inst *MonitorStep) SetPoolPcTokenAccountAccount(poolPcTokenAccount ag_solanago.PublicKey) *MonitorStep {
	inst.AccountMetaSlice[8] = ag_solanago.Meta(poolPcTokenAccount).WRITE()
	return inst
}

// GetPoolPcTokenAccountAccount gets the "poolPcTokenAccount" account.
func (inst *MonitorStep) GetPoolPcTokenAccountAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(8)
}

// SetPoolWithdrawQueueAccount sets the "poolWithdrawQueue" account.
func (inst *MonitorStep) SetPoolWithdrawQueueAccount(poolWithdrawQueue ag_solanago.PublicKey) *MonitorStep {
	inst.AccountMetaSlice[9] = ag_solanago.Meta(poolWithdrawQueue).WRITE()
	return inst
}

// GetPoolWithdrawQueueAccount gets the "poolWithdrawQueue" account.
func (inst *MonitorStep) GetPoolWithdrawQueueAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(9)
}

// SetSerumProgramAccount sets the "serumProgram" account.
func (inst *MonitorStep) SetSerumProgramAccount(serumProgram ag_solanago.PublicKey) *MonitorStep {
	inst.AccountMetaSlice[10] = ag_solanago.Meta(serumProgram)
	return inst
}

// GetSerumProgramAccount gets the "serumProgram" account.
func (inst *MonitorStep) GetSerumProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(10)
}

// SetSerumMarketAccount sets the "serumMarket" account.
func (inst *MonitorStep) SetSerumMarketAccount(serumMarket ag_solanago.PublicKey) *MonitorStep {
	inst.AccountMetaSlice[11] = ag_solanago.Meta(serumMarket).WRITE()
	return inst
}

// GetSerumMarketAccount gets the "serumMarket" account.
func (inst *MonitorStep) GetSerumMarketAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(11)
}

// SetSerumCoinVaultAccountAccount sets the "serumCoinVaultAccount" account.
func (inst *MonitorStep) SetSerumCoinVaultAccountAccount(serumCoinVaultAccount ag_solanago.PublicKey) *MonitorStep {
	inst.AccountMetaSlice[12] = ag_solanago.Meta(serumCoinVaultAccount).WRITE()
	return inst
}

// GetSerumCoinVaultAccountAccount gets the "serumCoinVaultAccount" account.
func (inst *MonitorStep) GetSerumCoinVaultAccountAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(12)
}

// SetSerumPcVaultAccountAccount sets the "serumPcVaultAccount" account.
func (inst *MonitorStep) SetSerumPcVaultAccountAccount(serumPcVaultAccount ag_solanago.PublicKey) *MonitorStep {
	inst.AccountMetaSlice[13] = ag_solanago.Meta(serumPcVaultAccount).WRITE()
	return inst
}

// GetSerumPcVaultAccountAccount gets the "serumPcVaultAccount" account.
func (inst *MonitorStep) GetSerumPcVaultAccountAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(13)
}

// SetSerumVaultSignerAccount sets the "serumVaultSigner" account.
func (inst *MonitorStep) SetSerumVaultSignerAccount(serumVaultSigner ag_solanago.PublicKey) *MonitorStep {
	inst.AccountMetaSlice[14] = ag_solanago.Meta(serumVaultSigner)
	return inst
}

// GetSerumVaultSignerAccount gets the "serumVaultSigner" account.
func (inst *MonitorStep) GetSerumVaultSignerAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(14)
}

// SetSerumReqQAccount sets the "serumReqQ" account.
func (inst *MonitorStep) SetSerumReqQAccount(serumReqQ ag_solanago.PublicKey) *MonitorStep {
	inst.AccountMetaSlice[15] = ag_solanago.Meta(serumReqQ).WRITE()
	return inst
}

// GetSerumReqQAccount gets the "serumReqQ" account.
func (inst *MonitorStep) GetSerumReqQAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(15)
}

// SetSerumEventQAccount sets the "serumEventQ" account.
func (inst *MonitorStep) SetSerumEventQAccount(serumEventQ ag_solanago.PublicKey) *MonitorStep {
	inst.AccountMetaSlice[16] = ag_solanago.Meta(serumEventQ).WRITE()
	return inst
}

// GetSerumEventQAccount gets the "serumEventQ" account.
func (inst *MonitorStep) GetSerumEventQAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(16)
}

// SetSerumBidsAccount sets the "serumBids" account.
func (inst *MonitorStep) SetSerumBidsAccount(serumBids ag_solanago.PublicKey) *MonitorStep {
	inst.AccountMetaSlice[17] = ag_solanago.Meta(serumBids).WRITE()
	return inst
}

// GetSerumBidsAccount gets the "serumBids" account.
func (inst *MonitorStep) GetSerumBidsAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(17)
}

// SetSerumAsksAccount sets the "serumAsks" account.
func (inst *MonitorStep) SetSerumAsksAccount(serumAsks ag_solanago.PublicKey) *MonitorStep {
	inst.AccountMetaSlice[18] = ag_solanago.Meta(serumAsks).WRITE()
	return inst
}

// GetSerumAsksAccount gets the "serumAsks" account.
func (inst *MonitorStep) GetSerumAsksAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(18)
}

func (inst MonitorStep) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: ag_binary.TypeIDFromUint8(Instruction_MonitorStep),
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst MonitorStep) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *MonitorStep) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.PlanOrderLimit == nil {
			return errors.New("PlanOrderLimit parameter is not set")
		}
		if inst.PlaceOrderLimit == nil {
			return errors.New("PlaceOrderLimit parameter is not set")
		}
		if inst.CancelOrderLimit == nil {
			return errors.New("CancelOrderLimit parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.TokenProgram is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.Rent is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.Clock is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.Amm is not set")
		}
		if inst.AccountMetaSlice[4] == nil {
			return errors.New("accounts.AmmAuthority is not set")
		}
		if inst.AccountMetaSlice[5] == nil {
			return errors.New("accounts.AmmOpenOrders is not set")
		}
		if inst.AccountMetaSlice[6] == nil {
			return errors.New("accounts.AmmTargetOrders is not set")
		}
		if inst.AccountMetaSlice[7] == nil {
			return errors.New("accounts.PoolCoinTokenAccount is not set")
		}
		if inst.AccountMetaSlice[8] == nil {
			return errors.New("accounts.PoolPcTokenAccount is not set")
		}
		if inst.AccountMetaSlice[9] == nil {
			return errors.New("accounts.PoolWithdrawQueue is not set")
		}
		if inst.AccountMetaSlice[10] == nil {
			return errors.New("accounts.SerumProgram is not set")
		}
		if inst.AccountMetaSlice[11] == nil {
			return errors.New("accounts.SerumMarket is not set")
		}
		if inst.AccountMetaSlice[12] == nil {
			return errors.New("accounts.SerumCoinVaultAccount is not set")
		}
		if inst.AccountMetaSlice[13] == nil {
			return errors.New("accounts.SerumPcVaultAccount is not set")
		}
		if inst.AccountMetaSlice[14] == nil {
			return errors.New("accounts.SerumVaultSigner is not set")
		}
		if inst.AccountMetaSlice[15] == nil {
			return errors.New("accounts.SerumReqQ is not set")
		}
		if inst.AccountMetaSlice[16] == nil {
			return errors.New("accounts.SerumEventQ is not set")
		}
		if inst.AccountMetaSlice[17] == nil {
			return errors.New("accounts.SerumBids is not set")
		}
		if inst.AccountMetaSlice[18] == nil {
			return errors.New("accounts.SerumAsks is not set")
		}
	}
	return nil
}

func (inst *MonitorStep) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("MonitorStep")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=3]").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("  PlanOrderLimit", *inst.PlanOrderLimit))
						paramsBranch.Child(ag_format.Param(" PlaceOrderLimit", *inst.PlaceOrderLimit))
						paramsBranch.Child(ag_format.Param("CancelOrderLimit", *inst.CancelOrderLimit))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=19]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("     tokenProgram", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("             rent", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("            clock", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta("              amm", inst.AccountMetaSlice.Get(3)))
						accountsBranch.Child(ag_format.Meta("     ammAuthority", inst.AccountMetaSlice.Get(4)))
						accountsBranch.Child(ag_format.Meta("    ammOpenOrders", inst.AccountMetaSlice.Get(5)))
						accountsBranch.Child(ag_format.Meta("  ammTargetOrders", inst.AccountMetaSlice.Get(6)))
						accountsBranch.Child(ag_format.Meta("    poolCoinToken", inst.AccountMetaSlice.Get(7)))
						accountsBranch.Child(ag_format.Meta("      poolPcToken", inst.AccountMetaSlice.Get(8)))
						accountsBranch.Child(ag_format.Meta("poolWithdrawQueue", inst.AccountMetaSlice.Get(9)))
						accountsBranch.Child(ag_format.Meta("     serumProgram", inst.AccountMetaSlice.Get(10)))
						accountsBranch.Child(ag_format.Meta("      serumMarket", inst.AccountMetaSlice.Get(11)))
						accountsBranch.Child(ag_format.Meta("   serumCoinVault", inst.AccountMetaSlice.Get(12)))
						accountsBranch.Child(ag_format.Meta("     serumPcVault", inst.AccountMetaSlice.Get(13)))
						accountsBranch.Child(ag_format.Meta(" serumVaultSigner", inst.AccountMetaSlice.Get(14)))
						accountsBranch.Child(ag_format.Meta("        serumReqQ", inst.AccountMetaSlice.Get(15)))
						accountsBranch.Child(ag_format.Meta("      serumEventQ", inst.AccountMetaSlice.Get(16)))
						accountsBranch.Child(ag_format.Meta("        serumBids", inst.AccountMetaSlice.Get(17)))
						accountsBranch.Child(ag_format.Meta("        serumAsks", inst.AccountMetaSlice.Get(18)))
					})
				})
		})
}

func (obj MonitorStep) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `PlanOrderLimit` param:
	err = encoder.Encode(obj.PlanOrderLimit)
	if err != nil {
		return err
	}
	// Serialize `PlaceOrderLimit` param:
	err = encoder.Encode(obj.PlaceOrderLimit)
	if err != nil {
		return err
	}
	// Serialize `CancelOrderLimit` param:
	err = encoder.Encode(obj.CancelOrderLimit)
	if err != nil {
		return err
	}
	return nil
}
func (obj *MonitorStep) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `PlanOrderLimit`:
	err = decoder.Decode(&obj.PlanOrderLimit)
	if err != nil {
		return err
	}
	// Deserialize `PlaceOrderLimit`:
	err = decoder.Decode(&obj.PlaceOrderLimit)
	if err != nil {
		return err
	}
	// Deserialize `CancelOrderLimit`:
	err = decoder.Decode(&obj.CancelOrderLimit)
	if err != nil {
		return err
	}
	return nil
}

// NewMonitorStepInstruction declares a new MonitorStep instruction with the provided parameters and accounts.
func NewMonitorStepInstruction(
	// Parameters:
	planOrderLimit uint16,
	placeOrderLimit uint16,
	cancelOrderLimit uint16,
	// Accounts:
	tokenProgram ag_solanago.PublicKey,
	rent ag_solanago.PublicKey,
	clock ag_solanago.PublicKey,
	amm ag_solanago.PublicKey,
	ammAuthority ag_solanago.PublicKey,
	ammOpenOrders ag_solanago.PublicKey,
	ammTargetOrders ag_solanago.PublicKey,
	poolCoinTokenAccount ag_solanago.PublicKey,
	poolPcTokenAccount ag_solanago.PublicKey,
	poolWithdrawQueue ag_solanago.PublicKey,
	serumProgram ag_solanago.PublicKey,
	serumMarket ag_solanago.PublicKey,
	serumCoinVaultAccount ag_solanago.PublicKey,
	serumPcVaultAccount ag_solanago.PublicKey,
	serumVaultSigner ag_solanago.PublicKey,
	serumReqQ ag_solanago.PublicKey,
	serumEventQ ag_solanago.PublicKey,
	serumBids ag_solanago.PublicKey,
	serumAsks ag_solanago.PublicKey) *MonitorStep {
	return NewMonitorStepInstructionBuilder().
		SetPlanOrderLimit(planOrderLimit).
		SetPlaceOrderLimit(placeOrderLimit).
		SetCancelOrderLimit(cancelOrderLimit).
		SetTokenProgramAccount(tokenProgram).
		SetRentAccount(rent).
		SetClockAccount(clock).
		SetAmmAccount(amm).
		SetAmmAuthorityAccount(ammAuthority).
		SetAmmOpenOrdersAccount(ammOpenOrders).
		SetAmmTargetOrdersAccount(ammTargetOrders).
		SetPoolCoinTokenAccountAccount(poolCoinTokenAccount).
		SetPoolPcTokenAccountAccount(poolPcTokenAccount).
		SetPoolWithdrawQueueAccount(poolWithdrawQueue).
		SetSerumProgramAccount(serumProgram).
		SetSerumMarketAccount(serumMarket).
		SetSerumCoinVaultAccountAccount(serumCoinVaultAccount).
		SetSerumPcVaultAccountAccount(serumPcVaultAccount).
		SetSerumVaultSignerAccount(serumVaultSigner).
		SetSerumReqQAccount(serumReqQ).
		SetSerumEventQAccount(serumEventQ).
		SetSerumBidsAccount(serumBids).
		SetSerumAsksAccount(serumAsks)
}
