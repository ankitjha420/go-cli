package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"strconv"
)

// generic type for all operations
type statsFunc func(data []float64) float64

func sum(data []float64) float64 {
	sum := 0.0
	for _, v := range data {
		sum += v
	}
	return sum
}

func avg(data []float64) float64 {
	return sum(data) / float64(len(data))
}

func colMin(data []float64) float64 {
	res := math.MaxFloat64
	for i := 0; i < len(data); i++ {
		if res > data[i] {
			res = data[i]
		}
	}

	return res
}

func colMax(data []float64) float64 {
	res := -math.MaxFloat64
	for i := 0; i < len(data); i++ {
		if res < data[i] {
			res = data[i]
		}
	}

	return res
}

func csv2float(r io.Reader, column int) ([]float64, error) {
	cr := csv.NewReader(r)
	cr.ReuseRecord = true

	// because columns in csv are not zero indexed like arrays in programming languages, so it has to be decremented ->
	column--

	var data []float64

	for i := 0; ; i++ { // INFINITE LOOP
		row, err := cr.Read()
		if err == io.EOF {
			break // end of file, can quit the loop
		}
		if err != nil {
			return nil, fmt.Errorf("cannot read data from this file %w", err)
		}

		if i == 0 {
			continue
		}

		if len(row) <= column {
			return nil, fmt.Errorf("%w: File has only %d columns", ErrInvalidColumn, len(row))
		}

		v, err := strconv.ParseFloat(row[column], 64)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrNotNumber, err)
		}

		data = append(data, v)
	}

	return data, nil
}
