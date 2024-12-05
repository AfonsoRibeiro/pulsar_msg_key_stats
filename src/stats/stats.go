package stats

import "sort"

type stats struct {
	group_size uint64

	len  uint64
	keys map[string]uint64
}

func getStats(bucketSize uint64) *stats {
	return &stats{
		group_size: bucketSize,
		len:        0,
		keys:       make(map[string]uint64),
	}
}

func (s *stats) addKey(key string) {
	s.keys[key]++
	s.len++
}

func (s *stats) check() (uint64, uint64, float64, float64) {
	var sum uint64 = 0
	var min uint64 = 0
	var max uint64 = 0
	var mean float64 = 0
	var median float64 = 0

	if s.len == s.group_size {

		// Collect all values for sorting (used for median calculation)
		values := make([]uint64, 0, len(s.keys))
		for _, value := range s.keys {
			sum += value
			values = append(values, value)
		}
		mean = float64(sum) / float64(len(s.keys))

		// Sort values for median calculation
		sort.Slice(values, func(i, j int) bool { return values[i] < values[j] })

		min = values[0]
		max = values[len(values)-1]

		// Calculate median
		n := len(values)
		if n%2 == 0 {
			// Even number of elements: median is the average of the two middle values
			median = float64(values[n/2-1]+values[n/2]) / 2.0
		} else {
			// Odd number of elements: median is the middle value
			median = float64(values[n/2])
		}

		s.len = 0
		s.keys = make(map[string]uint64)
	}
	return min, max, mean, median
}
