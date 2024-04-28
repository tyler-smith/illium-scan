package helpers

import (
	"strconv"
	"time"

	"github.com/maniartech/gotime"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
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

func FormatInt(i int) string {
	return strconv.Itoa(i)
}

func FormatUint64(i uint64) string {
	return strconv.Itoa(int(i))
}

func FormatTimeRelative(timestamp int) string {
	return gotime.TimeAgo(time.Unix(int64(timestamp), 0), time.Now())
}

func FormatAmount(amount uint64) string {
	return message.NewPrinter(language.English).Sprintf("%d", amount)
}
