package types

import (
	"bufio"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdkTypes "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client"
	authClient "github.com/cosmos/cosmos-sdk/x/auth/client"
	"github.com/spf13/cobra"
	"net/http"
)

type Transaction interface {
	GetModuleName() string
	GetName() string
	Command(*codec.Codec) *cobra.Command
	HandleMessage(sdkTypes.Context, TransactionKeeper, sdkTypes.Msg) (*sdkTypes.Result, error)
	RESTRequestHandler(context.CLIContext) http.HandlerFunc
	RegisterCodec(*codec.Codec)
}

type transaction struct {
	ModuleName                  string
	Name                        string
	CLICommand                  CLICommand
	Codec                       func(*codec.Codec)
	TransactionRequestPrototype func() TransactionRequest
}

var _ Transaction = (*transaction)(nil)

func (transaction transaction) GetModuleName() string { return transaction.ModuleName }
func (transaction transaction) GetName() string       { return transaction.Name }
func (transaction transaction) Command(codec *codec.Codec) *cobra.Command {
	runE := func(command *cobra.Command, args []string) error {
		bufioReader := bufio.NewReader(command.InOrStdin())
		transactionBuilder := auth.NewTxBuilderFromCLI(bufioReader).WithTxEncoder(authClient.GetTxEncoder(codec))
		cliContext := context.NewCLIContextWithInput(bufioReader).WithCodec(codec)

		request := transaction.TransactionRequestPrototype().FromCLI(transaction.CLICommand, cliContext)

		msg := request.MakeMsg()
		if Error := msg.ValidateBasic(); Error != nil {
			return Error
		}

		return authClient.GenerateOrBroadcastMsgs(cliContext, transactionBuilder, []sdkTypes.Msg{msg})
	}
	return transaction.CLICommand.CreateCommand(runE)
}

func (transaction transaction) HandleMessage(context sdkTypes.Context, transactionKeeper TransactionKeeper, message sdkTypes.Msg) (*sdkTypes.Result, error) {

	if Error := transactionKeeper.Transact(context, message); Error != nil {
		return nil, Error
	}

	context.EventManager().EmitEvent(
		sdkTypes.NewEvent(
			sdkTypes.EventTypeMessage,
			sdkTypes.NewAttribute(sdkTypes.AttributeKeyModule, transaction.ModuleName),
		),
	)

	return &sdkTypes.Result{Events: context.EventManager().ABCIEvents()}, nil
}

func (transaction transaction) RESTRequestHandler(cliContext context.CLIContext) http.HandlerFunc {
	return func(responseWriter http.ResponseWriter, httpRequest *http.Request) {
		request := transaction.TransactionRequestPrototype()
		if !rest.ReadRESTReq(responseWriter, httpRequest, cliContext.Codec, &request) {
			rest.WriteErrorResponse(responseWriter, http.StatusBadRequest, "")
			return
		}

		baseReq := request.GetBaseReq()
		msg := request.MakeMsg()

		baseReq = baseReq.Sanitize()
		if !baseReq.ValidateBasic(responseWriter) {
			rest.WriteErrorResponse(responseWriter, http.StatusBadRequest, "")
			return
		}

		Error := msg.ValidateBasic()
		if Error != nil {
			rest.WriteErrorResponse(responseWriter, http.StatusBadRequest, Error.Error())
			return
		}
		client.WriteGenerateStdTxResponse(responseWriter, cliContext, baseReq, []sdkTypes.Msg{msg})
	}
}

func (transaction transaction) RegisterCodec(codec *codec.Codec) {
	transaction.Codec(codec)
}

func NewTransaction(module string, name string, short string, long string, transactionRequestPrototype func() TransactionRequest, registerCodec func(*codec.Codec), flagList []CLIFlag) Transaction {
	return &transaction{
		ModuleName:                  module,
		Name:                        name,
		CLICommand:                  NewCLICommand(name, short, long, flagList),
		Codec:                       registerCodec,
		TransactionRequestPrototype: transactionRequestPrototype,
	}
}