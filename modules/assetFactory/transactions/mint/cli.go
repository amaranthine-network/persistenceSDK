package mint

import (
	"bufio"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authClient "github.com/cosmos/cosmos-sdk/x/auth/client"
	"github.com/persistenceOne/persistenceSDK/modules/assetFactory/constants"
	"github.com/persistenceOne/persistenceSDK/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"strconv"
)

func TransactionCommand(codec *codec.Codec) *cobra.Command {

	command := &cobra.Command{
		Use:   constants.MintTransaction,
		Short: "Create and sign a transaction to mint an asset",
		Long:  "",
		RunE: func(command *cobra.Command, args []string) error {
			bufioReader := bufio.NewReader(command.InOrStdin())
			transactionBuilder := auth.NewTxBuilderFromCLI(bufioReader).WithTxEncoder(authClient.GetTxEncoder(codec))
			cliContext := context.NewCLIContextWithInput(bufioReader).WithCodec(codec)

			basePropertyList := make([]types.BaseProperty, 0)
			for i := 0; i <= constants.MaxTraitCount; i++ {
				if viper.GetString(constants.TraitID+strconv.Itoa(i)) != "" {
					basePropertyList = append(basePropertyList,
						types.BaseProperty{
							BaseID:   types.BaseID{IDString: viper.GetString(constants.TraitID + strconv.Itoa(i))},
							BaseFact: types.BaseFact{BaseBytes: []byte(viper.GetString(constants.Property + strconv.Itoa(i)))},
						})
				}
			}
			baseProperties := types.BaseProperties{
				BasePropertyList: basePropertyList,
			}

			message := Message{
				From:             cliContext.GetFromAddress(),
				ChainID:          types.BaseID{IDString: viper.GetString(constants.ChainID)},
				MaintainersID:    types.BaseID{IDString: viper.GetString(constants.MaintainersID)},
				ClassificationID: types.BaseID{IDString: viper.GetString(constants.ClassificationID)},
				Properties:       &baseProperties,
				Lock:             types.BaseHeight{Height: viper.GetInt(constants.Lock)},
				Burn:             types.BaseHeight{Height: viper.GetInt(constants.Burn)},
			}

			if Error := message.ValidateBasic(); Error != nil {
				return Error
			}

			return authClient.GenerateOrBroadcastMsgs(cliContext, transactionBuilder, []sdkTypes.Msg{message})
		},
	}
	command.Flags().String(constants.ChainID, "", "ChainID")
	command.Flags().String(constants.MaintainersID, "", "MaintainersID")
	command.Flags().String(constants.ClassificationID, "", "ClassificationID")
	for i := 0; i <= constants.MaxTraitCount; i++ {
		command.Flags().String(constants.TraitID+strconv.Itoa(i), "", "traitID")
		command.Flags().String(constants.Property+strconv.Itoa(i), "", "property")
	}
	command.Flags().Int(constants.Lock, -1, "Lock")
	command.Flags().Int(constants.Burn, -1, "Burn")
	return flags.PostCommands(command)[0]
}

func NewTransactionCommand(codecMarshaler codec.Marshaler, txGenerator tx.Generator, accountRetriever tx.AccountRetriever) *cobra.Command {

	command := &cobra.Command{
		Use:   constants.MintTransaction,
		Short: "Create and sign a transaction to mint an asset",
		Long:  "",
		RunE: func(command *cobra.Command, args []string) error {
			bufioReader := bufio.NewReader(command.InOrStdin())
			cliContext := context.NewCLIContextWithInputAndFrom(bufioReader, args[0]).WithMarshaler(codecMarshaler)
			txFactory := tx.NewFactoryFromCLI(bufioReader).WithTxGenerator(txGenerator).WithAccountRetriever(accountRetriever)

			var basePropertyList []types.BaseProperty
			for i := 0; i <= constants.MaxTraitCount; i++ {
				if viper.GetString(viper.GetString(constants.TraitID+strconv.Itoa(i))) != "" {
					basePropertyList = append(basePropertyList,
						types.BaseProperty{
							BaseID:   types.BaseID{IDString: viper.GetString(constants.TraitID + strconv.Itoa(i))},
							BaseFact: types.BaseFact{BaseBytes: []byte(viper.GetString(constants.Property + strconv.Itoa(i)))},
						})
				}
			}
			baseProperties := types.BaseProperties{BasePropertyList: basePropertyList}
			message := Message{
				From:             cliContext.GetFromAddress(),
				ChainID:          types.BaseID{IDString: viper.GetString(constants.ChainID)},
				MaintainersID:    types.BaseID{IDString: viper.GetString(constants.MaintainersID)},
				ClassificationID: types.BaseID{IDString: viper.GetString(constants.ClassificationID)},
				Properties:       &baseProperties,
				Lock:             types.BaseHeight{Height: viper.GetInt(constants.Lock)},
				Burn:             types.BaseHeight{Height: viper.GetInt(constants.Burn)},
			}

			if Error := message.ValidateBasic(); Error != nil {
				return Error
			}

			return tx.GenerateOrBroadcastTx(cliContext, txFactory, message)
		},
	}
	command.Flags().String(constants.ChainID, "", "ChainID")
	command.Flags().String(constants.MaintainersID, "", "MaintainersID")
	command.Flags().String(constants.ClassificationID, "", "ClassificationID")
	for i := 0; i <= constants.MaxTraitCount; i++ {
		command.Flags().String(constants.TraitID+strconv.Itoa(i), "", "traitID")
		command.Flags().String(constants.Property+strconv.Itoa(i), "", "property")
	}
	command.Flags().Int(constants.Lock, -1, "Lock")
	command.Flags().Int(constants.Burn, -1, "Burn")
	return flags.PostCommands(command)[0]
}
