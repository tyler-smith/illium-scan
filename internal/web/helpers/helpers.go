package helpers

import (
	"strconv"
	"strings"
	"time"

	"github.com/maniartech/gotime"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/tyler-smith/iexplorer/internal/db/models"
)

func FormatBlockHeight(height int) string {
	//if height == 0 {
	//	return "(genesis)"
	//}
	return FormatInt(height)
}

func FormatIDShort(id string) string {
	if id == "" {
		return "(none)"
	}
	return id[:8]
}

func FormatValidatorIDShort(id string) string {
	if id == "" {
		return "(none)"
	}
	return id[:8] + "..." + id[len(id)-8:]
}

func FormatInt(i int) string {
	return strconv.Itoa(i)
}

func FormatUint64(i uint64) string {
	return strconv.Itoa(int(i))
}

func FormatTimeRelative(timestamp int) string {
	return strings.ToLower(gotime.TimeAgo(time.Unix(int64(timestamp), 0), time.Now()))
}

func FormatAmount(amount uint64) string {
	return message.NewPrinter(language.English).Sprintf("%d", amount)
}

func FormatType(t models.TxType) string {
	switch t {
	case 0:
		return "Standard"
	case 1:
		return "Coinbase"
	case 2:
		return "Stake"
	case 3:
		return "Treasury"
	case 4:
		return "Mint"
	default:
		return "Unknown"
	}
}

func FormatLocktime(locktime int) string {
	if locktime == 0 {
		return "(none)"
	}
	return time.Unix(int64(locktime), 0).Format(time.DateTime)
}
