package metrics

import "strconv"

type Counter int64

func CounterFromString(value string) (Counter, error) {
	counter, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, ErrBadCounterValue
	}
	return Counter(counter), nil
}
