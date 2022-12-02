package telebotCalendar

import (
	"fmt"
	tele "gopkg.in/telebot.v3"
	"strconv"
	"time"
)

const BTN_PREV = "<"
const BTN_NEXT = ">"

func GenerateCalendar(year int, month time.Month) (InlineKeyboard [][]tele.InlineButton) {
	keyboard := [][]tele.InlineButton{}

	keyboard = addMonthYearRow(year, month, keyboard)
	keyboard = addDaysNamesRow(keyboard)
	keyboard = generateMonth(year, int(month), keyboard)
	keyboard = addSpecialButtons(keyboard)
	return keyboard
}

func HandlerPrevButton(year int, month time.Month) (InlineKeyboard [][]tele.InlineButton, result int, chattime time.Month) {
	if month != 1 {
		month--
	} else {
		month = 12
		year--
	}
	return GenerateCalendar(year, month), year, month
}

func HandlerNextButton(year int, month time.Month) (InlineKeyboard [][]tele.InlineButton, result int, chattime time.Month) {
	if month != 12 {
		month++
	} else {
		year++
	}
	return GenerateCalendar(year, month), year, month
}

func addMonthYearRow(year int, month time.Month, keyboard [][]tele.InlineButton) (InlineKeyboard [][]tele.InlineButton) {
	var row []tele.InlineButton
	btn := tele.InlineButton{Text: fmt.Sprintf("%s %v", month, year), Data: "1"}
	row = append(row, btn)
	keyboard = append(keyboard, row)
	return keyboard
}

func addDaysNamesRow(keyboard [][]tele.InlineButton) (InlineKeyboard [][]tele.InlineButton) {
	days := [7]string{"Mo", "Tu", "We", "Th", "Fr", "Sa", "Su"}
	var rowDays []tele.InlineButton
	for _, day := range days {
		btn := tele.InlineButton{Text: day, Data: day}
		rowDays = append(rowDays, btn)
	}
	keyboard = append(keyboard, rowDays)
	return keyboard
}

func generateMonth(year int, month int, keyboard [][]tele.InlineButton) (InlineKeyboard [][]tele.InlineButton) {
	firstDay := date(year, month, 0)
	amountDaysInMonth := date(year, month+1, 0).Day()

	weekday := int(firstDay.Weekday())
	rowDays := []tele.InlineButton{}
	for i := 1; i <= weekday; i++ {
		btn := tele.InlineButton{Text: " ", Data: string(i)}
		rowDays = append(rowDays, btn)
	}

	amountWeek := weekday
	for i := 1; i <= amountDaysInMonth; i++ {
		if amountWeek == 7 {
			keyboard = append(keyboard, rowDays)
			amountWeek = 0
			rowDays = []tele.InlineButton{}
		}

		day := strconv.Itoa(i)
		if len(day) == 1 {
			day = fmt.Sprintf("0%v", day)
		}
		monthStr := strconv.Itoa(month)
		if len(monthStr) == 1 {
			monthStr = fmt.Sprintf("0%v", monthStr)
		}

		btnText := fmt.Sprintf("%v", i)
		if time.Now().Day() == i {
			btnText = fmt.Sprintf("%v!", i)
		}
		btn := tele.InlineButton{Text: btnText, Data: fmt.Sprintf("%v.%v.%v", year, monthStr, day)}
		rowDays = append(rowDays, btn)
		amountWeek++
	}
	for i := 1; i <= 7-amountWeek; i++ {
		btn := tele.InlineButton{Text: " ", Data: string(i)}
		rowDays = append(rowDays, btn)
	}

	keyboard = append(keyboard, rowDays)

	return keyboard
}

func date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
func addSpecialButtons(keyboard [][]tele.InlineButton) (InlineKeyboard [][]tele.InlineButton) {
	var rowDays []tele.InlineButton
	btnPrev := tele.InlineButton{Text: BTN_PREV, Data: BTN_PREV}
	btnNext := tele.InlineButton{Text: BTN_NEXT, Data: BTN_NEXT}
	rowDays = append(rowDays, btnPrev, btnNext)
	keyboard = append(keyboard, rowDays)
	return keyboard
}
