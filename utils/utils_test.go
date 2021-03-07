package utils

import (
	"github.com/ethereum/go-ethereum/common"
	"math/big"
	"testing"
)

var infura = "https://ropsten.infura.io/v3/1b65e3621097432db8199a8f6252dc53"
var master_chief_contract = common.HexToAddress("0xf11C813041477A55be888C989d690307dFB902f0")

func TestGetPoolInfo(t *testing.T) {

	length, err := GetPoolLength(infura, master_chief_contract)
	if err != nil {
		//t.Fatal(err)
		t.Log(err)
	}
	t.Log(length)

	pools, err := GetPoolInfo(infura, master_chief_contract, big.NewInt(0))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(pools)

	total, err := GetTotalSupply(infura, common.HexToAddress("0xd40aC65B1771dA5cFC5fE43dA2854d2ff659A562"))
	if err != nil {
		t.Fatal(err)
	}
	t.Log(total)
}

func TestGetDepositAndWithLogs(t *testing.T) {
	//logs, err := GetDepositAndWithLogs(infura, big.NewInt(0), common.HexToAddress("0xbf861002a10f1dd32b62f37cd0794574c4494193"))
	logs, err := GetDepositAndWithLogs(infura, big.NewInt(9009852),
		[]common.Address{common.HexToAddress("0xdcd0dd7569f84567c5f95d4298805851491be171"),
			common.HexToAddress("0x3cd6b049ecb9deb6c1521359c445cfdf2eeb52ec")})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(len(logs))
	for _, l := range logs {
		t.Log(l)
	}
}

func TestGetTokenPrice(t *testing.T) {
	price, err := GetPitayaPrice()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(price)
	p,err:=GetPriceFromCoinw("hc")
	if err!=nil{
		t.Fatal(err)
	}
	t.Log(p)
	//GetTokenBalance()
}

func TestGetLpPrice(t *testing.T) {
	GetLpPrice("0x4d5ef58aac27d99935e5b6b4a6778ff292059991")
}