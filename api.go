package noah_node_go_api

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/noah-blockchain/go-amino"
	"github.com/noah-blockchain/noah-explorer-tools/models"
	"github.com/noah-blockchain/noah-node-go-api/responses"
	"github.com/valyala/fasthttp"
)

type NoahNodeApi struct {
	link   string
	client *fasthttp.Client
	cdc    *amino.Codec

	fallbackRetries int
	fallbackTimeout time.Duration
}

func New(link string) *NoahNodeApi {
	// Initialization
	cdc := amino.NewCodec()

	return &NoahNodeApi{
		link: link,
		client: &fasthttp.Client{
			Name:                "Explorer Extender API",
			MaxIdleConnDuration: 5,
		},
		cdc: cdc,
	}
}

func NewWithFallbackRetries(link string, fallbackRetries int, fallbackTimeout time.Duration) *NoahNodeApi {
	// Initialization
	cdc := amino.NewCodec()

	return &NoahNodeApi{
		link: link,
		client: &fasthttp.Client{
			Name:                "Explorer Extender API",
			MaxIdleConnDuration: 5,
		},
		cdc:             cdc,
		fallbackRetries: fallbackRetries,
		fallbackTimeout: fallbackTimeout,
	}
}

func (api *NoahNodeApi) SetLink(link string) {
	api.link = link
}

func (api *NoahNodeApi) GetLink() string {
	return api.link
}

func (api *NoahNodeApi) GetStatus() (*responses.StatusResponse, error) {
	response := responses.StatusResponse{}
	link := api.link + `/status`
	err := api.getJson(link, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

func (api *NoahNodeApi) GetBlock(height uint64) (*responses.BlockResponse, error) {
	response := responses.BlockResponse{}
	link := api.link + `/block?height=` + fmt.Sprint(height)
	err := api.getJson(link, &response)
	if err != nil {
		return nil, err
	}
	if response.Result.TxCount != "0" {
		for i, tx := range response.Result.Transactions {
			switch tx.Type {
			case models.TxTypeSend:
				var txData = models.SendTxData{}
				err = api.cdc.UnmarshalJSON(tx.Data, &txData)
				response.Result.Transactions[i].IData = txData
			case models.TxTypeSellCoin:
				var txData = models.SellCoinTxData{}
				err = api.cdc.UnmarshalJSON(tx.Data, &txData)
				response.Result.Transactions[i].IData = txData
			case models.TxTypeSellAllCoin:
				var txData = models.SellAllCoinTxData{}
				err = api.cdc.UnmarshalJSON(tx.Data, &txData)
				response.Result.Transactions[i].IData = txData
			case models.TxTypeBuyCoin:
				var txData = models.BuyCoinTxData{}
				err = api.cdc.UnmarshalJSON(tx.Data, &txData)
				response.Result.Transactions[i].IData = txData
			case models.TxTypeCreateCoin:
				var txData = models.CreateCoinTxData{}
				err = api.cdc.UnmarshalJSON(tx.Data, &txData)
				response.Result.Transactions[i].IData = txData
			case models.TxTypeDeclareCandidacy:
				var txData = models.DeclareCandidacyTxData{}
				err = api.cdc.UnmarshalJSON(tx.Data, &txData)
				response.Result.Transactions[i].IData = txData
			case models.TxTypeDelegate:
				var txData = models.DelegateTxData{}
				err = api.cdc.UnmarshalJSON(tx.Data, &txData)
				response.Result.Transactions[i].IData = txData
			case models.TxTypeUnbound:
				var txData = models.UnbondTxData{}
				err = api.cdc.UnmarshalJSON(tx.Data, &txData)
				response.Result.Transactions[i].IData = txData
			case models.TxTypeRedeemCheck:
				var txData = models.RedeemCheckTxData{}
				err = api.cdc.UnmarshalJSON(tx.Data, &txData)
				response.Result.Transactions[i].IData = txData
			case models.TxTypeSetCandidateOnline:
				var txData = models.SetCandidateTxData{}
				err = api.cdc.UnmarshalJSON(tx.Data, &txData)
				response.Result.Transactions[i].IData = txData
			case models.TxTypeSetCandidateOffline:
				var txData = models.SetCandidateTxData{}
				err = api.cdc.UnmarshalJSON(tx.Data, &txData)
				response.Result.Transactions[i].IData = txData
			case models.TxTypeMultiSig:
				var txData = models.CreateMultisigTxData{}
				err = api.cdc.UnmarshalJSON(tx.Data, &txData)
				response.Result.Transactions[i].IData = txData
			case models.TxTypeMultiSend:
				var txData = models.MultiSendTxData{}
				err = api.cdc.UnmarshalJSON(tx.Data, &txData)
				response.Result.Transactions[i].IData = txData
			case models.TxTypeEditCandidate:
				var txData = models.EditCandidateTxData{}
				err = api.cdc.UnmarshalJSON(tx.Data, &txData)
				response.Result.Transactions[i].IData = txData
			}
		}
	}
	return &response, err
}

func (api *NoahNodeApi) GetBlockEvents(height uint64) (*responses.EventsResponse, error) {
	response := responses.EventsResponse{}
	link := api.link + `/events?height=` + fmt.Sprint(height)
	err := api.getJson(link, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

func (api *NoahNodeApi) GetBlockValidators(height uint64) (*responses.ValidatorsResponse, error) {
	response := responses.ValidatorsResponse{}
	link := api.link + `/validators?height=` + fmt.Sprint(height)
	err := api.getJson(link, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

func (api *NoahNodeApi) GetCandidate(pubKey string, height uint64) (*responses.CandidateResponse, error) {
	response := responses.CandidateResponse{}
	link := api.link + `/candidate?pubkey=` + pubKey + `&height=` + strconv.Itoa(int(height))
	err := api.getJson(link, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

func (api *NoahNodeApi) GetCandidates(height uint64, stakes bool) (*responses.BlockCandidatesResponse, error) {
	response := responses.BlockCandidatesResponse{}
	link := api.link + `/candidates?height=` + strconv.Itoa(int(height))

	if stakes {
		link += `&include_stakes=true`
	}

	err := api.getJson(link, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

func (api *NoahNodeApi) GetCoinInfo(symbol string) (*responses.CoinInfoResponse, error) {
	response := responses.CoinInfoResponse{}
	link := api.link + `/coin_info?symbol=` + symbol
	err := api.getJson(link, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

func (api *NoahNodeApi) GetAddress(address string) (*responses.AddressResponse, error) {
	response := responses.AddressResponse{}
	link := api.link + `/address?address=` + strings.Title(address)
	err := api.getJson(link, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

func (api *NoahNodeApi) GetAddresses(addresses []string, height uint64) (*responses.BalancesResponse, error) {
	response := responses.BalancesResponse{}
	queryStr := "[" + strings.Join(addresses, ",") + "]"
	link := api.link + `/addresses?addresses=` + queryStr + `&height=` + strconv.Itoa(int(height))
	err := api.getJson(link, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

func (api *NoahNodeApi) GetEstimateTx(tx string) (*responses.EstimateTxResponse, error) {
	response := responses.EstimateTxResponse{}
	link := api.link + `/estimate_tx_commission?tx=` + tx
	err := api.getJson(link, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

func (api *NoahNodeApi) GetEstimateCoinBuy(coinToSell string, coinToBuy string, value string) (*responses.EstimateCoinBuyResponse, error) {
	response := responses.EstimateCoinBuyResponse{}
	link := api.link + `/estimate_coin_buy?coin_to_sell=` + coinToSell + `&coin_to_buy=` + coinToBuy + `&value_to_buy=` + value
	err := api.getJson(link, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

func (api *NoahNodeApi) GetEstimateCoinSell(coinToSell string, coinToBuy string, value string, height uint64) (*responses.EstimateCoinSellResponse, error) {
	response := responses.EstimateCoinSellResponse{}
	link := api.link + `/estimate_coin_sell?coin_to_sell=` + coinToSell + `&coin_to_buy=` + coinToBuy + `&value_to_sell=` + value + `&height=` + fmt.Sprint(height)
	err := api.getJson(link, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

func (api *NoahNodeApi) GetEstimateCoinSellAll(coinToSell string, coinToBuy string, value string, gasPrice string) (*responses.EstimateCoinSellAllResponse, error) {
	response := responses.EstimateCoinSellAllResponse{}
	link := api.link + `/estimate_coin_sell_all?coin_to_sell=` + coinToSell + `&coin_to_buy=` + coinToBuy + `&value_to_sell=` + value + `&gas_price=` + gasPrice
	err := api.getJson(link, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

func (api *NoahNodeApi) GetMinGasPrice() (*responses.GasResponse, error) {
	response := responses.GasResponse{}
	link := api.link + `/min_gas_price`
	err := api.getJson(link, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

func (api *NoahNodeApi) PushTransaction(tx string) (*responses.SendTransactionResponse, error) {
	response := responses.SendTransactionResponse{}
	link := api.link + `/send_transaction?tx=0x` + tx
	err := api.getJson(link, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

func (api *NoahNodeApi) GetTransactionsByQuery(query string) (*responses.TransactionsResponse, error) {
	response := responses.TransactionsResponse{}
	link := fmt.Sprintf(api.link+"/transactions?query=%s", url.QueryEscape(query))
	err := api.getJson(link, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

func (api *NoahNodeApi) GetTransaction(hash string) (*responses.TransactionResponse, error) {
	response := responses.TransactionResponse{}
	link := api.link + `/transaction?hash=` + hash
	err := api.getJson(link, &response)
	if err != nil {
		return nil, err
	}
	return &response, err
}

func (api *NoahNodeApi) getJson(url string, target interface{}) error {
	retries := 0

	var err error
	for retries <= api.fallbackRetries {
		err = (func() error {
			_, body, err := api.client.Get(nil, url)
			if err != nil {
				return err
			}

			err = api.cdc.UnmarshalJSON(body, target)

			return err
		})()
		if err == nil {
			return nil
		}
		retries++

		time.Sleep(api.fallbackTimeout)
	}

	return err
}
