package history

import (
	"sort"
	"strconv"
	"time"

	"github.com/fibercrypto/fibercryptowallet/src/coin/skycoin/params"
	"github.com/fibercrypto/fibercryptowallet/src/util/logging"

	coin "github.com/fibercrypto/fibercryptowallet/src/coin/skycoin/models"
	"github.com/fibercrypto/fibercryptowallet/src/core"
	local "github.com/fibercrypto/fibercryptowallet/src/main"
	"github.com/fibercrypto/fibercryptowallet/src/models/address"
	"github.com/fibercrypto/fibercryptowallet/src/models/transactions"
	"github.com/fibercrypto/fibercryptowallet/src/util"
	qtCore "github.com/therecipe/qt/core"
)

var logHistoryManager = logging.MustGetLogger("modelsHistoryManager")

const (
	dateTimeFormatForGo  = "2006-01-02T15:04:05"
	dateTimeFormatForQML = "yyyy-MM-ddThh:mm:ss"
)

/*
	HistoryManager
	Represent the controller of history page and all the actions over this page
*/
type HistoryManager struct {
	qtCore.QObject
	filters []string
	_       func() `constructor:"init"`

	_         func() []*transactions.TransactionDetails `slot:"loadHistoryWithFilters"`
	_         func() []*transactions.TransactionDetails `slot:"loadHistory"`
	_         func(string)                              `slot:"addFilter"`
	_         func(string)                              `slot:"removeFilter"`
	walletEnv core.WalletEnv
}

func (hm *HistoryManager) init() {
	hm.ConnectLoadHistoryWithFilters(hm.loadHistoryWithFilters)
	hm.ConnectLoadHistory(hm.loadHistory)
	hm.ConnectAddFilter(hm.addFilter)
	hm.ConnectRemoveFilter(hm.removeFilter)
	altManager := local.LoadAltcoinManager()
	walletsEnvs := make([]core.WalletEnv, 0)
	for _, plug := range altManager.ListRegisteredPlugins() {
		walletsEnvs = append(walletsEnvs, plug.LoadWalletEnvs()...)
	}

	hm.walletEnv = walletsEnvs[0]
}

type ByDate []*transactions.TransactionDetails

func (a ByDate) Len() int {
	return len(a)
}
func (a ByDate) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByDate) Less(i, j int) bool {
	d1, _ := time.Parse(dateTimeFormatForGo, a[i].Date().ToString(dateTimeFormatForQML))
	d2, _ := time.Parse(dateTimeFormatForGo, a[j].Date().ToString(dateTimeFormatForQML))
	return d1.After(d2)
}
func (hm *HistoryManager) getTransactionsOfAddresses(filterAddresses []string) []*transactions.TransactionDetails {
	logHistoryManager.Info("Getting transactions of Addresses")
	addresses := hm.getAddressesWithWallets()

	var sent, internally bool
	var traspassedHoursIn, traspassedHoursOut, skyAmountIn, skyAmountOut uint64

	find := make(map[string]struct{}, len(filterAddresses))
	for _, addr := range filterAddresses {
		find[addr] = struct{}{}
	}
	txnFind := make(map[string]struct{})
	txns := make([]core.Transaction, 0)

	wltIterator := hm.walletEnv.GetWalletSet().ListWallets()
	if wltIterator == nil {
		logHistoryManager.WithError(nil).Warn("Couldn't get transactions of Addresses")
		return make([]*transactions.TransactionDetails, 0)
	}
	for wltIterator.Next() {
		addressIterator, err := wltIterator.Value().GetLoadedAddresses()
		if err != nil {
			logHistoryManager.Warn("Couldn't get address iterator")
			continue
		}
		for addressIterator.Next() {
			_, ok := find[addressIterator.Value().String()]
			if ok {
				txnsIterator := addressIterator.Value().GetCryptoAccount().ListTransactions()
				if txnsIterator == nil {
					logHistoryManager.Warn("Couldn't get transaction iterator")
					continue
				}
				for txnsIterator.Next() {
					_, ok2 := txnFind[txnsIterator.Value().GetId()]
					if !ok2 {
						txns = append(txns, txnsIterator.Value())
						txnFind[txnsIterator.Value().GetId()] = struct{}{}
					}
				}
			}

		}
	}

	txnsDetails := make([]*transactions.TransactionDetails, 0)
	for _, txn := range txns {
		traspassedHoursIn = 0
		traspassedHoursOut = 0
		skyAmountIn = 0
		skyAmountOut = 0
		internally = true
		sent = false
		txnDetails := transactions.NewTransactionDetails(nil)
		txnAddresses := address.NewAddressList(nil)
		inAddresses := make(map[string]struct{}, 0)
		inputs := address.NewAddressList(nil)
		outputs := address.NewAddressList(nil)
		txnIns := txn.GetInputs()

		for _, in := range txnIns {
			qIn := address.NewAddressDetails(nil)
			qIn.SetAddress(in.GetSpentOutput().GetAddress().String())
			skyUint64, err := in.GetCoins(params.SkycoinTicker)
			if err != nil {
				logHistoryManager.WithError(err).Warn("Couldn't get Skycoins balance")
				continue
			}
			accuracy, err := util.AltcoinQuotient(params.SkycoinTicker)
			if err != nil {
				logHistoryManager.WithError(err).Warn("Couldn't get Skycoins quotient")
				continue
			}
			skyFloat := float64(skyUint64) / float64(accuracy)
			qIn.SetAddressSky(strconv.FormatFloat(skyFloat, 'f', -1, 64))
			chUint64, err := in.GetCoins(params.CoinHoursTicker)
			if err != nil {
				logHistoryManager.WithError(err).Warn("Couldn't get Coin Hours balance")
				continue
			}
			accuracy, err = util.AltcoinQuotient(params.CoinHoursTicker)
			if err != nil {
				logHistoryManager.WithError(err).Warn("Couldn't get Coin Hours quotient")
				continue
			}
			qIn.SetAddressCoinHours(util.FormatCoins(chUint64, accuracy))
			inputs.AddAddress(qIn)
			_, ok := addresses[in.GetSpentOutput().GetAddress().String()]
			if ok {
				skyAmountOut += skyUint64
				sent = true
				_, ok := inAddresses[qIn.Address()]
				if !ok {
					txnAddresses.AddAddress(qIn)
					inAddresses[qIn.Address()] = struct{}{}
				}

			}
		}
		txnDetails.SetInputs(inputs)

		for _, out := range txn.GetOutputs() {
			sky, err := out.GetCoins(params.SkycoinTicker)
			if err != nil {
				logHistoryManager.WithError(err).Warn("Couldn't get Skycoins balance")
				continue
			}
			qOu := address.NewAddressDetails(nil)
			qOu.SetAddress(out.GetAddress().String())
			accuracy, err := util.AltcoinQuotient(params.SkycoinTicker)
			if err != nil {
				logHistoryManager.WithError(err).Warn("Couldn't get Skycoins quotient")
				continue
			}
			qOu.SetAddressSky(util.FormatCoins(sky, accuracy))
			val, err := out.GetCoins(params.CoinHoursTicker)
			if err != nil {
				logHistoryManager.WithError(err).Warn("Couldn't get Coin Hours balance")
				continue
			}
			accuracy, err = util.AltcoinQuotient(coin.CoinHour)
			if err != nil {
				logHistoryManager.WithError(err).Warn("Couldn't get Coin Hours quotient")
				continue
			}
			qOu.SetAddressCoinHours(util.FormatCoins(val, accuracy))
			outputs.AddAddress(qOu)
			if sent {

				if addresses[txn.GetInputs()[0].GetSpentOutput().GetAddress().String()] == addresses[out.GetAddress().String()] {
					skyAmountOut -= sky

				} else {
					internally = false
					val, err = out.GetCoins(params.CoinHoursTicker)
					if err != nil {
						logHistoryManager.WithError(err).Warn("Couldn't get Coin Hours send it")
						continue
					}
					traspassedHoursOut += val
				}
			} else {
				_, ok := find[out.GetAddress().String()]
				if ok {
					val, err = out.GetCoins(params.CoinHoursTicker)
					if err != nil {
						logHistoryManager.WithError(err).Warn("Couldn't get Coin Hours balance")
						continue
					}
					traspassedHoursIn += val
					skyAmountIn += sky

					_, ok := inAddresses[qOu.Address()]
					if !ok {
						txnAddresses.AddAddress(qOu)
						inAddresses[qOu.Address()] = struct{}{}
					}

				}

			}

		}
		txnDetails.SetOutputs(outputs)
		t := time.Unix(int64(txn.GetTimestamp()), 0)
		txnDetails.SetDate(qtCore.NewQDateTime3(qtCore.NewQDate3(t.Year(), int(t.Month()), t.Day()), qtCore.NewQTime3(t.Hour(), t.Minute(), 0, 0), qtCore.Qt__LocalTime))
		txnDetails.SetStatus(transactions.TransactionStatusPending)

		if txn.GetStatus() == core.TXN_STATUS_CONFIRMED {
			txnDetails.SetStatus(transactions.TransactionStatusConfirmed)
		}
		txnDetails.SetType(transactions.TransactionTypeReceive)
		if sent {
			txnDetails.SetType(transactions.TransactionTypeSend)
			if internally {
				txnDetails.SetType(transactions.TransactionTypeInternal)
			}
		}
		fee, err := txn.ComputeFee(params.CoinHoursTicker)
		if err != nil {
			logHistoryManager.WithError(err).Warn("Couldn't compute fee of the operation")
			continue
		}
		accuracy, err := util.AltcoinQuotient(coin.CoinHoursTicker)
		if err != nil {
			logHistoryManager.WithError(err).Warn("Couldn't get " + coin.CoinHoursTicker + " coins quotient")
		}
		txnDetails.SetHoursBurned(util.FormatCoins(fee, accuracy))

		switch txnDetails.Type() {
		case transactions.TransactionTypeReceive:
			{
				accuracy, err := util.AltcoinQuotient(coin.CoinHoursTicker)
				if err != nil {
					logHistoryManager.WithError(err).Warn("Couldn't get " + coin.CoinHoursTicker + " coins quotient")
				}
				txnDetails.SetHoursTraspassed(util.FormatCoins(traspassedHoursIn, accuracy))
				val := float64(skyAmountIn)
				accuracy, err = util.AltcoinQuotient(params.SkycoinTicker)
				if err != nil {
					logHistoryManager.WithError(err).Warn("Couldn't get Skycoins quotient")
					continue
				}
				val = val / float64(accuracy)
				txnDetails.SetAmount(strconv.FormatFloat(val, 'f', -1, 64))

			}
		case transactions.TransactionTypeInternal:
			{
				var traspassedHoursMoved, skyAmountMoved uint64
				traspassedHoursMoved = 0
				skyAmountMoved = 0
				ins := inputs.Addresses()
				inFind := make(map[string]struct{}, len(ins))
				for _, addr := range ins {
					inFind[addr.Address()] = struct{}{}
				}
				outs := outputs.Addresses()
				for _, addr := range outs {
					_, ok := inFind[addr.Address()]
					if !ok {
						hours, err := strconv.ParseUint(addr.AddressCoinHours(), 10, 64)
						if err != nil {
							logHistoryManager.WithError(err).Warn("Couldn't parse Coin Hours from address")
							continue
						}
						traspassedHoursMoved += hours
						skyFloat, err := strconv.ParseFloat(addr.AddressSky(), 64)
						if err != nil {
							logHistoryManager.WithError(err).Warn("Couldn't parse Skycoins from addresses")
							continue
						}
						accuracy, err := util.AltcoinQuotient(params.SkycoinTicker)
						if err != nil {
							logHistoryManager.WithError(err).Warn("Couldn't get Skycoins quotient")
							continue
						}
						sky := uint64(skyFloat * float64(accuracy))
						skyAmountMoved += sky
					}

				}
				accuracy, err := util.AltcoinQuotient(coin.CoinHoursTicker)
				if err != nil {
					logHistoryManager.WithError(err).Warn("Couldn't get " + coin.CoinHoursTicker + " coins quotient")
				}
				txnDetails.SetHoursTraspassed(util.FormatCoins(traspassedHoursMoved, accuracy))
				val := float64(skyAmountMoved)
				//FIXME: Error here is skipped
				accuracy, _ = util.AltcoinQuotient(params.SkycoinTicker)
				if err != nil {
					logHistoryManager.WithError(err).Warn("Couldn't get Skycoins quotient")
					continue
				}
				val = val / float64(accuracy)
				txnDetails.SetAmount(strconv.FormatFloat(val, 'f', -1, 64))

			}
		case transactions.TransactionTypeSend:
			{
				accuracy, err := util.AltcoinQuotient(coin.CoinHoursTicker)
				if err != nil {
					logHistoryManager.WithError(err).Warn("Couldn't get " + coin.CoinHoursTicker + " coins quotient")
				}
				txnDetails.SetHoursTraspassed(util.FormatCoins(traspassedHoursOut, accuracy))
				val := float64(skyAmountOut)
				accuracy, err = util.AltcoinQuotient(params.SkycoinTicker)
				if err != nil {
					logHistoryManager.WithError(err).Warn("Couldn't get Skycoins quotient")
					continue
				}
				val = val / float64(accuracy)
				txnDetails.SetAmount(strconv.FormatFloat(val, 'f', -1, 64))

			}
		}
		txnDetails.SetAddresses(txnAddresses)
		txnDetails.SetTransactionID(txn.GetId())

		txnsDetails = append(txnsDetails, txnDetails)

	}
	sort.Sort(ByDate(txnsDetails))
	return txnsDetails
}
func (hm *HistoryManager) loadHistoryWithFilters() []*transactions.TransactionDetails {
	logHistoryManager.Info("Loading history with some filters")
	filterAddresses := hm.filters
	return hm.getTransactionsOfAddresses(filterAddresses)

}

func (hm *HistoryManager) loadHistory() []*transactions.TransactionDetails {
	logHistoryManager.Info("Loading history")
	addresses := hm.getAddressesWithWallets()

	filterAddresses := make([]string, 0)
	for addr, _ := range addresses {
		filterAddresses = append(filterAddresses, addr)
	}

	return hm.getTransactionsOfAddresses(filterAddresses)

}

func (hm *HistoryManager) addFilter(addr string) {
	logHistoryManager.Info("Add filter")
	alreadyIs := false
	for _, filter := range hm.filters {
		if filter == addr {
			alreadyIs = true
			break
		}
	}
	if !alreadyIs {
		hm.filters = append(hm.filters, addr)
	}

}

func (hm *HistoryManager) removeFilter(addr string) {
	logHistoryManager.Info("Remove filter")

	for i := 0; i < len(hm.filters); i++ {
		if hm.filters[i] == addr {
			hm.filters = append(hm.filters[0:i], hm.filters[i+1:]...)
			break
		}
	}

}
func (hm *HistoryManager) getAddressesWithWallets() map[string]string {
	logHistoryManager.Info("Get Addresses with wallets")
	response := make(map[string]string, 0)
	it := hm.walletEnv.GetWalletSet().ListWallets()
	if it == nil {
		logHistoryManager.WithError(nil).Warn("Couldn't load addresses")
		return response
	}
	for it.Next() {
		wlt := it.Value()
		addresses, err := wlt.GetLoadedAddresses()
		if err != nil {
			logHistoryManager.WithError(err).Warn("Couldn't get loaded addresses")
			continue
		}
		for addresses.Next() {
			response[addresses.Value().String()] = wlt.GetId()
		}

	}

	return response
}
