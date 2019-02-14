package numberRange

import "strconv"
import "fmt"

/*
 * ported generator from bezmax's java code from
 * https://stackoverflow.com/questions/33512037/a-regular-expression-generator-for-number-ranges
 * to golang
 */

type numberRange struct {
	start int
	end   int
}

func GetRegex(start, end int) string {
	left := leftBounds(start, end)
	lastLeft := left[len(left)-1]
	left = left[:len(left)-1]

	right := rightBounds(lastLeft.getStart(), end)
	firstRight := right[0]
	right = right[1:]

	merged := make([]numberRange, 0, 100)
	merged = append(merged, left...)

	if !lastLeft.overlaps(firstRight) {
		merged = append(merged, lastLeft)
		merged = append(merged, firstRight)
	} else {
		merged = append(merged, joinRange(lastLeft, firstRight))
	}

	merged = append(merged, right...)

	var finalStr string
	for i := 0; i < len(merged); i++ {
		finalStr += "_" + merged[i].toRegex()
		if ( i+1 < len(merged)) {
			finalStr += "|"
		}
	}
	return finalStr
}

func leftBounds(start, end int) []numberRange {
	listRange := make([]numberRange, 0, 100)
	for start < end {
		neueRange := fromStart(start)
		start = neueRange.getEnd() + 1
		listRange = append(listRange, neueRange)
	}
	return listRange
}

func rightBounds(start, end int) []numberRange {
	listRange := make([]numberRange, 0, 100)
	for start < end {
		neueRange := fromEnd(end)
		end = neueRange.getStart() - 1
		listRange = append(listRange, neueRange)
	}

	listRange = reverse(listRange)
	return listRange
}

func (n numberRange) getStart() int {
	return n.start
}

func (n numberRange) getEnd() int {
	return n.end
}

func (n numberRange) overlaps(n2 numberRange) bool {
	return n.end > n2.start && n2.end > n.start
}

func (n numberRange) String() string {
	return fmt.Sprintf("Range{start=%d, end=%d}", n.start, n.end)
}

func (n numberRange) toRegex() string {
	var finalStr string
	startStr := []rune(strconv.Itoa(n.start))
	endStr := []rune(strconv.Itoa(n.end))

	for index, _ := range startStr {
		if startStr[index] == endStr[index] {
			finalStr += fmt.Sprintf("%c", startStr[index])
		} else {
			finalStr += "[" + fmt.Sprintf("%c", startStr[index]) + "-" + fmt.Sprintf("%c", endStr[index]) + "]"
		}
	}
	return finalStr
}

func fromEnd(pEnd int) numberRange {
	chars := []rune(strconv.Itoa(pEnd))
	for i := len(chars) - 1; i >= 0; i-- {
		if chars[i] == '9' {
			chars[i] = '0'
		} else {
			chars[i] = '0'
			break
		}
	}

	var finalStr string
	for index, _ := range chars {
		finalStr += fmt.Sprintf("%c", chars[index])
	}

	newStart, err := strconv.Atoi(finalStr)
	if err != nil {
		panic(err)
	}

	return NewRange(newStart, pEnd)
}

func fromStart(pStart int) numberRange {
	chars := []rune(strconv.Itoa(pStart))
	for i := len(chars) - 1; i >= 0; i-- {
		if chars[i] == '0' {
			chars[i] = '9'
		} else {
			chars[i] = '9'
			break
		}
	}

	var finalStr string
	for index, _ := range chars {
		finalStr += fmt.Sprintf("%c", chars[index])
	}

	newEnd, err := strconv.Atoi(finalStr)
	if err != nil {
		panic(err)
	}

	return NewRange(pStart, newEnd)
}

func joinRange(n1, n2 numberRange) numberRange {
	return NewRange(n1.getStart(), n2.getEnd())
}

func NewRange(start, end int) numberRange {
	return numberRange{start, end}
}

func reverse(list []numberRange) []numberRange {
	reversedList := make([]numberRange, len(list), len(list))

	for index, _ := range list {
		reversedList[len(list)-1-index] = list[index]
	}

	return reversedList
}