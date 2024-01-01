package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ParseError struct{}

func (p *ParseError) Error() string {
	return fmt.Sprint("Failed to parse date")
}

type Date struct {
	year, month, day int
}

func (d *Date) String() string {
	return fmt.Sprintf("%04d-%02d-%02d", d.year, d.month, d.day)
}

func (d *Date) Scan(src interface{}) error {
	fmt.Println(src)
	switch v := src.(type) {
	case string:
		slices := strings.Split(v, "-")
		if len(slices) != 3 {
			return fmt.Errorf("%T - Incorrect number of date components", src)
		}
		nums := []int{}
		for _, slice := range slices {
			num, err := strconv.Atoi(slice)
			if err != nil {
				return fmt.Errorf("%T - Unable to parse date component %T", src, slice)
			}
			nums = append(nums, num)
		}

		d.year, d.month, d.day = nums[0], nums[1], nums[2]
		return nil
	}
	return fmt.Errorf("%T - Invalid type", src)
}

func DateNow() Date {
	now := time.Now()
	date := Date{}
	date.FromTime(now.Date())
	return date
}

func ParseDate(str string) (Date, error) {
	slices := strings.Split(str, "-")
	if len(slices) != 3 {
		return Date{}, &ParseError{}
	}
	nums := []int{}
	for _, slice := range slices {
		num, err := strconv.Atoi(slice)
		if err != nil {
			return Date{}, &ParseError{}
		}
		nums = append(nums, num)
	}
	date := Date{}

	date.year, date.month, date.day = nums[0], nums[1], nums[2]
	return date, nil
}

func (d *Date) FromTime(year int, month time.Month, day int) {
	d.year, d.month, d.day = year, int(month), day
}
