package zeroaprlib

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

type DataError struct {
	What string
}

func (err DataError) Error() string {
	return err.What
}

type dataRoot struct {
	Purchases []Purchase
	Version   string
}

type Payment struct {
	PaymentDate   time.Time
	PaymentAmount int
}

type Purchase struct {
	Name              string
	Account           string
	Memo              string
	Date              time.Time
	InitialAmount     int
	WantedNumPayments int
	Payments          []Payment
}

type DataAccess struct {
	filename string
}

func Connect(path string) (*DataAccess, error) {
	da := DataAccess{path}
	return &da, nil
}

func (da *DataAccess) save(root *dataRoot) error {
	rawBytes, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return err
	}

	dataFile, err := os.Create(da.filename)
	if err != nil {
		return err
	}
	defer dataFile.Close()

	_, err = dataFile.Write(rawBytes)
	if err != nil {
		return err
	}

	return nil
}

func (da *DataAccess) GetAllPurchases() []Purchase {
	loaded, err := da.loadFile()
	if err != nil {
		// TODO: error handling.
		fmt.Println("Error:", err.Error())
	}

	return loaded.Purchases
}

func (da *DataAccess) GetPurchase(name string) *Purchase {
	root, err := da.loadFile()
	if err != nil {
		return nil
	}

	return findPurchaseByName(root, name)
}

func (da *DataAccess) PurchaseExists(name string) bool {
	root, err := da.loadFile()
	if err != nil {
		return false
	}

	return findPurchaseByName(root, name) != nil
}

func (da *DataAccess) AddPurchase(newPurchase Purchase) error {
	root, err := da.loadFile()
	if err != nil {
		return err
	}

	anyExist := findPurchaseByName(root, newPurchase.Name)
	if anyExist != nil {
		return DataError{"Purchase with name already exists"}
	}

	root.Purchases = append(root.Purchases, newPurchase)

	err = da.save(root)
	if err != nil {
		return err
	}

	return nil
}

// Add a Payment to the purchase identified by the given name.
func (da *DataAccess) AddPayment(purchaseName string, newPayment Payment) error {
	root, err := da.loadFile()
	if err != nil {
		return err
	}

	// Must have a valid and existing purchase.
	pur := findPurchaseByName(root, purchaseName)
	if pur == nil {
		return DataError{"Purchase not found"}
	}

	pur.Payments = append(pur.Payments, newPayment)

	err = da.save(root)
	if err != nil {
		return err
	}

	return nil
}

func (da *DataAccess) RemovePayment(purchaseName string, index int) error {
	root, err := da.loadFile()
	if err != nil {
		return err
	}

	// Must have a valid and existing purchase.
	pur := findPurchaseByName(root, purchaseName)
	if pur == nil {
		return DataError{"Purchase not found"}
	}

	if index >= len(pur.Payments) {
		return DataError{"Invalid index"}
	}

	pur.Payments = append(pur.Payments[:index], pur.Payments[index+1:]...)

	da.save(root)

	return nil
}

// Find a purchase with the given key/name
func findPurchaseByName(data *dataRoot, name string) *Purchase {
	// Must use index here. Otherwise, the element will be copied, and results will not be pointing to the
	// expected address.
	for i := range data.Purchases {
		if data.Purchases[i].Name == name {
			return &data.Purchases[i]
		}
	}

	return nil
}

func (da *DataAccess) loadFile() (*dataRoot, error) {
	// Handle case where the file does not already exist.
	if _, err := os.Stat(da.filename); errors.Is(err, os.ErrNotExist) {
		root := dataRoot{Version: "1.0"}
		return &root, nil
	} else {
		reader, err := os.Open(da.filename)
		if err != nil {
			return nil, err
		}
		defer reader.Close()

		allData, err := io.ReadAll(reader)
		if err != nil {
			return nil, err
		}

		var root dataRoot
		err = json.Unmarshal(allData, &root)
		if err != nil {
			return nil, err
		}

		return &root, nil
	}
}
