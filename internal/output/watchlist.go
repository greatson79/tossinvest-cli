package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"

	"github.com/junghoonkye/tossinvest-cli/internal/domain"
)

func WriteWatchlist(w io.Writer, format Format, items []domain.WatchlistItem) error {
	switch format {
	case FormatJSON:
		encoder := json.NewEncoder(w)
		encoder.SetIndent("", "  ")
		return encoder.Encode(items)
	case FormatCSV:
		writer := csv.NewWriter(w)
		if err := writer.Write([]string{"group", "symbol", "name", "currency", "base", "last"}); err != nil {
			return err
		}
		for _, item := range items {
			if err := writer.Write([]string{
				item.Group,
				item.Symbol,
				item.Name,
				item.Currency,
				formatFloat(item.Base),
				formatFloat(item.Last),
			}); err != nil {
				return err
			}
		}
		writer.Flush()
		return writer.Error()
	case FormatTable:
		headers := []string{"그룹", "종목", "이름", "기준가", "현재가", "통화"}
		var rows [][]string
		for _, item := range items {
			rows = append(rows, []string{
				item.Group,
				item.Symbol,
				item.Name,
				formatKRW(item.Base),
				formatKRW(item.Last),
				item.Currency,
			})
		}
		return renderTable(w, headers, rows)
	default:
		return fmt.Errorf("unsupported output format: %s", format)
	}
}
