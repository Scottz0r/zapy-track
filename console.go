package main

import (
	"bufio"
	"fmt"
	zeroaprlib "go-zero-apr-mgr/zero-apr-lib"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

const dateFormat = "2006-01-02"

func consoleLoop(da *zeroaprlib.DataAccess) error {
	fmt.Println("Enter command: exit show new pay detail")

	for {
		fmt.Print("~>")
		command, err := readNextString()
		if err != nil {
			return err
		}

		splitCommand := strings.Split(command, " ")

		if len(splitCommand) == 0 {
			fmt.Println("Invalid command")
		} else {
			switch splitCommand[0] {
			case "exit":
				return nil
			case "show":
				showSummary(da)
			case "new":
				err := addPurchase(da)
				if err != nil {
					fmt.Println(err)
				}
			case "pay":
				if len(splitCommand) < 2 {
					fmt.Println(`"pay" needs a name`)
				} else {
					err := addPayment(da, splitCommand[1])
					if err != nil {
						fmt.Println(err)
					}
				}
			case "detail":
				if len(splitCommand) < 2 {
					fmt.Println(`"detail" needs a name`)
				} else {
					showDetails(da, splitCommand[1])
				}
			default:
				fmt.Println("Unknown command")
			}
		}
	}
}

// Show an overall summary.
func showSummary(da *zeroaprlib.DataAccess) {
	allPurchases := da.GetAllPurchases()

	fmt.Printf(
		"%15s %15s %15s %15s %15s\n",
		"Name",
		"Remain",
		"Paid",
		"Remain Pmts",
		"Pmt Amt")

	dashString := "---------------"
	fmt.Println(dashString, dashString, dashString, dashString, dashString)

	for i := range allPurchases {
		summary := zeroaprlib.CalculatePurchaseSummary(&allPurchases[i])

		p := &allPurchases[i]
		fmt.Printf(
			"%15s %15.2f %15.2f %15d %15.2f\n",
			p.Name,
			convertFromCurrency(summary.RemainingAmount),
			convertFromCurrency(summary.TotalPaid),
			summary.RemainingNumPayments,
			convertFromCurrency(summary.PeriodPaymentAmount))
	}
}

// Add a purchase.
func addPurchase(da *zeroaprlib.DataAccess) error {
	fmt.Print("Name: ")
	name, err := readNextString()
	if err != nil {
		return err
	}

	fmt.Print("Account: ")
	account, err := readNextString()
	if err != nil {
		return err
	}

	fmt.Print("Memo: ")
	memo, err := readNextString()
	if err != nil {
		return err
	}

	fmt.Print("Date (YYYY-MM-DD): ")
	parsedDate, err := readNextDate()
	if err != nil {
		return err
	}

	fmt.Print("Amount: ")
	amountF, err := readNextFloat()
	if err != nil {
		return err
	}
	amountInt := convertToCurrency(amountF)

	fmt.Print("Num Payments: ")
	numPaymentInt, err := readNextInteger()
	if err != nil {
		return err
	}

	if da.PurchaseExists(name) {
		return fmt.Errorf("a purchase with name '%s' already exists", name)
	}

	newPurchase := zeroaprlib.Purchase{
		Name:              name,
		Account:           account,
		Memo:              memo,
		Date:              parsedDate,
		InitialAmount:     amountInt,
		WantedNumPayments: numPaymentInt,
	}

	err = da.AddPurchase(newPurchase)
	if err != nil {
		return err
	}

	fmt.Println("Purchase added successfully")
	return nil
}

// Add a payment.
func addPayment(da *zeroaprlib.DataAccess, name string) error {
	if !da.PurchaseExists(name) {
		return fmt.Errorf("a purchase named '%s' does not exist", name)
	}

	fmt.Print("Date (YYYY-MM-DD): ")
	parsedDate, err := readNextDate()
	if err != nil {
		return err
	}

	fmt.Print("Amount: ")
	amountF, err := readNextFloat()
	if err != nil {
		return err
	}
	amountInt := convertToCurrency(amountF)

	newPayment := zeroaprlib.Payment{
		PaymentDate:   parsedDate,
		PaymentAmount: amountInt,
	}

	err = da.AddPayment(name, newPayment)
	if err != nil {
		return err
	}

	return nil
}

// Show details of a single purchase. Gracefully handles the "not found" case.
func showDetails(da *zeroaprlib.DataAccess, name string) {
	purchase := da.GetPurchase(name)
	if purchase == nil {
		fmt.Println("Not found")
		return
	}

	fmt.Println("Name: ", purchase.Name)
	fmt.Println("Account: ", purchase.Account)
	fmt.Println("Memo: ", purchase.Memo)
	fmt.Println("Date: ", purchase.Date.Format(dateFormat))
	fmt.Println("Amount: ", convertFromCurrency(purchase.InitialAmount))
	fmt.Println("Wanted Payments: ", purchase.WantedNumPayments)

	calcSummary := zeroaprlib.CalculatePurchaseSummary(purchase)

	fmt.Println("Total Paid: ", convertFromCurrency(calcSummary.TotalPaid))
	fmt.Println("Remaining Amount: ", convertFromCurrency(calcSummary.RemainingAmount))
	fmt.Println("Remaining Pmts: ", calcSummary.RemainingNumPayments)
	fmt.Println("Pmts/period: ", calcSummary.PeriodPaymentAmount)

	fmt.Println("Payments:")
	for i := range purchase.Payments {
		dateFormatted := purchase.Payments[i].PaymentDate.Format(dateFormat)
		amountFormatted := convertFromCurrency(purchase.Payments[i].PaymentAmount)
		fmt.Printf("%10s %10.2f\n", dateFormatted, amountFormatted)
	}
}

// Small helper to convert an integer in smallest denomination into a decimal.
func convertFromCurrency(value int) float64 {
	return float64(value) / 100.0
}

func convertToCurrency(value float64) int {
	return int(math.Round(value * 100.0))
}

// Helper method to read the next string from the console.
func readNextString() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	value, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(value), nil
}

// Helper method to read the next integer from the console.
func readNextInteger() (int, error) {
	reader := bufio.NewReader(os.Stdin)

	valueString, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	valueString = strings.TrimSpace(valueString)

	valueInt, err := strconv.ParseInt(valueString, 10, 32)
	if err != nil {
		return 0, err
	}

	// Reduce the integer to native, which is going to be 32-bit (or more) for targets of this program.
	reduced32 := int(valueInt)
	return reduced32, nil
}

// Helper method to read the next decimal/float from the console.
func readNextFloat() (float64, error) {
	reader := bufio.NewReader(os.Stdin)

	valueString, err := reader.ReadString('\n')
	if err != nil {
		return 0, err
	}

	valueString = strings.TrimSpace(valueString)

	amountF, err := strconv.ParseFloat(valueString, 64)
	if err != nil {
		return 0, err
	}

	return amountF, nil
}

// Helper method to read the next date value from the console.
func readNextDate() (time.Time, error) {
	reader := bufio.NewReader(os.Stdin)

	valueString, err := reader.ReadString('\n')
	if err != nil {
		return time.Time{}, err
	}

	valueString = strings.TrimSpace((valueString))

	parsedTime, err := time.Parse("2006-01-02", valueString)
	if err != nil {
		return time.Time{}, err
	}

	return parsedTime, nil
}
