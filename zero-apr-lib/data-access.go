package zeroaprlib

import (
	"encoding/json"
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
	filename  string
	purchases []Purchase
}

func Connect(path string) (*DataAccess, error) {
	dataFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer dataFile.Close()

	allData, err := io.ReadAll(dataFile)
	if err != nil {
		return nil, err
	}

	var root dataRoot
	err = json.Unmarshal(allData, &root)
	if err != nil {
		return nil, err
	}

	da := DataAccess{path, root.Purchases}
	return &da, nil
}

func (da *DataAccess) Save() error {
	// Create root object that holds data and metadata
	root := dataRoot{da.purchases, "1.0"}

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
	return da.purchases
}

func (da *DataAccess) GetPurchase(name string) *Purchase {
	return da.findPurchaseByName(name)
}

func (da *DataAccess) PurchaseExists(name string) bool {
	return da.findPurchaseByName(name) != nil
}

func (da *DataAccess) AddPurchase(newPurchase Purchase) error {
	anyExist := da.findPurchaseByName(newPurchase.Name)
	if anyExist != nil {
		return DataError{"Purchase with name already exists"}
	}

	da.purchases = append(da.purchases, newPurchase)

	return nil
}

// Add a Payment to the purchase identified by the given name.
func (da *DataAccess) AddPayment(purchaseName string, newPayment Payment) error {
	// Must have a valid and existing purchase.
	pur := da.findPurchaseByName(purchaseName)
	if pur == nil {
		return DataError{"Purchase not found"}
	}

	pur.Payments = append(pur.Payments, newPayment)

	return nil
}

func (da *DataAccess) RemovePayment(purchaseName string, index int) error {
	// Must have a valid and existing purchase.
	pur := da.findPurchaseByName(purchaseName)
	if pur == nil {
		return DataError{"Purchase not found"}
	}

	if index >= len(pur.Payments) {
		return DataError{"Invalid index"}
	}

	pur.Payments = append(pur.Payments[:index], pur.Payments[index+1:]...)

	return nil
}

// Find a purchase with the given key/name
func (da *DataAccess) findPurchaseByName(name string) *Purchase {
	// Must use index here. Otherwise, the element will be copied, and results will not be pointing to the
	// expected address.
	for i := range da.purchases {
		if da.purchases[i].Name == name {
			return &da.purchases[i]
		}
	}

	return nil
}
