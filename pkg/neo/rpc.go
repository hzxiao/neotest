package neo

import (
	"fmt"
	"github.com/CityOfZion/neo-go/pkg/core/transaction"
	"github.com/CityOfZion/neo-go/pkg/util"
	"github.com/OTCGO/sea-server-go/pkg/neo"
	"github.com/hzxiao/goutil"
	"github.com/hzxiao/neotest/pkg/jsonrpc2"
	"sort"
	"strconv"
	"strings"
)

const (
	GasAssetHash = `602c79718b16e442de58778e148d0b1084e3b2dffd5de6b7b16cee7969282de7`
	NeoAssetHash = `c56f33fc6ecfcd0c225c4ab356fee59390af8560be0e930faebe74a6daff7c9b`
)

func Rpc(url, method string, params interface{}, result interface{}) error {
	r := &jsonrpc2.JRpcRequest{
		ID:     1,
		Method: method,
	}
	err := r.SetParams(params)
	if err != nil {
		return err
	}

	return jsonrpc2.Send(url, r, &result)
}

type References []goutil.Map

func (a References) Len() int           { return len(a) }
func (a References) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a References) Less(i, j int) bool { return a[i].GetFloat64("value") < a[j].GetFloat64("value") }

//getReference get unspent asset as input equal or large than given value
func getReference(asset, address string, value float64, node string) (float64, []*transaction.Input, error){
	var res goutil.Map
	err := Rpc(node, "getunspents", []string{address}, &res)
	if err != nil {
		return 0, nil, err
	}

	var all float64
	var reference []*transaction.Input
	for _, balance := range res.GetMapArray("balance") {
		if balance.GetString("asset_hash") != asset {
			continue
		}
		if value > balance.GetFloat64("amount") {
			return 0, nil, fmt.Errorf("not enough of %v", asset)
		}
		unspent := balance.GetMapArray("unspent")
		sort.Sort(References(unspent))
		for _, ref := range unspent {
			all += ref.GetFloat64("value")
			hash, _ := util.Uint256DecodeString(ref.GetString("txid"))
			reference = append(reference, &transaction.Input{
				PrevHash: hash,
				PrevIndex: uint16(ref.GetInt64("n")),
			})

			if all > value {
				return all, reference, nil
			}
		}
	}
	if all < value {
		return 0, nil, fmt.Errorf("not enough of %v", asset)
	}

	return all, reference, nil
}

//getAssetDecimals get asset decimals
func getAssetDecimals(node, asset string) (uint8, error) {
	asset = strings.TrimPrefix(asset, "0x")
	switch asset {
	case GasAssetHash:
		return 8, nil
	case NeoAssetHash:
		return 0, nil
	}

	if IsGlobalAsset(asset) {
		return getGlobalAssetDecimals(node, asset)
	} else {
		return getNep5AssetDecimals(node, asset)
	}
}

func getGlobalAssetDecimals(node string, asset string) (uint8, error) {
	var res goutil.Map
	err := Rpc(node, "getassetstate", []string{asset}, &res)
	if err != nil {
		return 0, err
	}
	return uint8(res.GetInt64("precision")), nil
}

func getNep5AssetDecimals(node string, contract string) (uint8, error) {
	v, success, err := rpcInvoke(node, []string{contract, "decimals"})
	if err != nil {
		return 0, fmt.Errorf("rpc invoke func(%v) fail(%v)", "decimals", err)
	}
	if !success {
		return 0, fmt.Errorf("rpc invoke fail: maybe doesn't have decimals func")
	}
	d, _ := strconv.Atoi(v.GetString("value"))
	return uint8(d), nil
}

func rpcInvoke(node string, params interface{}) (goutil.Map, bool, error) {
	var invokeFail = func(r goutil.Map) bool {
		if r == nil {
			return false
		}
		return strings.HasPrefix(r.GetString("state"), "FAULT")
	}

	var result = goutil.Map{}
	err := neo.Rpc(node, "invokefunction", params, &result)
	if err != nil {
		return nil, false, err
	}
	if invokeFail(result) {
		return nil, false, nil
	}

	return result.GetMapP("stack/0"), true, nil
}
