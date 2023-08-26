package main

import (
	"fmt"
	"strconv"
	"time"

	"github.com/rivo/tview"
)

type Bag struct {
	amount int64
}

func (bag *Bag) deposit(n int64) {
	bag.amount += n
}

func (bag *Bag) withdraw(n int64) {
	bag.amount -= n
}

func (bag *Bag) getTax(rate float64) float64 {
	return float64(bag.amount) * rate
}

const (
	DEPOSIT  = 1
	WITHDRAW = 2
)

type Transaction struct {
	variant   int64
	amount    int64
	timestamp time.Time
}

func (t *Transaction) String() string {
	action := "Undefined"

	switch t.variant {
	case DEPOSIT:
		action = "Deposit"
	case WITHDRAW:
		action = "Withdraw"
	default:
		panic(fmt.Sprintf("Undefined action %d %d", t.variant, t.amount))
	}

	return fmt.Sprintf("%s: %s %d", t.timestamp.Format(time.UnixDate), action, t.amount)
}

func main() {
	const (
		LABEL_DEPOSIT  = "Deposit"
		LABEL_WITHDRAW = "Withdraw"
	)
	var transactions []Transaction
	transactions = make([]Transaction, 0)
	bag := Bag{amount: 0}

	app := tview.NewApplication()

	bagView := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	bagView.SetBorder(true).SetTitle("Bag").SetTitleAlign(tview.AlignLeft)
	taxView := tview.NewTextView().SetTextAlign(tview.AlignCenter)
	taxView.SetBorder(true).SetTitle("Tax").SetTitleAlign(tview.AlignLeft)

	transactionView := tview.NewList()
	transactionView.SetBorder(true).SetTitle("Transactions")

	form := tview.NewForm().
		AddInputField(LABEL_DEPOSIT, "", 0, tview.InputFieldInteger, nil).
		AddInputField(LABEL_WITHDRAW, "", 0, tview.InputFieldInteger, nil)

	form.AddButton("Submit", func() {
		depositInputView := form.GetFormItemByLabel(LABEL_DEPOSIT).(*tview.InputField)
		withdrawInputView := form.GetFormItemByLabel(LABEL_WITHDRAW).(*tview.InputField)

		depositAmount, err := strconv.ParseInt(depositInputView.GetText(), 10, 64)

		if err == nil {
			bag.deposit(depositAmount)
			transactions = append(transactions,
				Transaction{
					variant:   DEPOSIT,
					amount:    depositAmount,
					timestamp: time.Now(),
				})
			depositInputView.SetText("")
		}

		withdrawAmount, err := strconv.ParseInt(withdrawInputView.GetText(), 10, 64)

		if err == nil {
			bag.withdraw(withdrawAmount)
			transactions = append(transactions,
				Transaction{
					variant:   WITHDRAW,
					amount:    withdrawAmount,
					timestamp: time.Now(),
				})
			withdrawInputView.SetText("")
		}

		bagView.SetText(fmt.Sprint(bag.amount))
		taxView.SetText(fmt.Sprintf("%0.2f", bag.getTax(0.2)))
		transactionView.Clear()
		for i := 0; i < len(transactions); i++ {
			transactionView = transactionView.AddItem(transactions[i].String(), "", 0, nil)
		}
		app.SetFocus(depositInputView)
	})

	form.SetBorder(true).SetTitle("Action").SetTitleAlign(tview.AlignLeft)

	stateFlex := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(bagView, 0, 1, false).
		AddItem(taxView, 0, 1, false)

	stateFlex.SetBorder(true).SetTitle("State").SetTitleAlign(tview.AlignLeft)

	statusFlex := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(stateFlex, 0, 1, false).
		AddItem(transactionView, 0, 1, false)

	flex := tview.NewFlex().
		AddItem(form, 0, 1, true).
		AddItem(statusFlex, 0, 1, false)

	if err := app.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
