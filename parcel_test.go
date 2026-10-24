package main

import (
//	"fmt"
	"database/sql"
	"math/rand"
	"testing"
	"time"
"github.com/stretchr/testify/assert"
//	"github.com/stretchr/testify/require"
)

var (
	// randSource источник псевдо случайных чисел.
	// Для повышения уникальности в качестве seed
	// используется текущее время в unix формате (в виде числа)
	randSource = rand.NewSource(time.Now().UnixNano())
	// randRange использует randSource для генерации случайных чисел
	randRange = rand.New(randSource)
)

// getTestParcel возвращает тестовую посылку
func getTestParcel() Parcel {
	return Parcel{
		Client:    1000,
		Status:    ParcelStatusRegistered,
		Address:   "test",
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
}

// TestAddGetDelete проверяет добавление, получение и удаление посылки
func TestAddGetDelete(t *testing.T) {
	// prepare
	db, err := sql.Open("sqlite", "tracker.db") 
	if err != nil {
		t.Fatalf("не удалось подключиться к базе данных: %v", err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	id, err := store.Add(parcel)
	if err != nil {
		t.Fatalf("ошибка при добавлении посылки: %v", err)
	}
	if id == 0 {
		t.Fatal("ожидался ненулевой идентификатор (ID) посылки")
	}
	parcel.Number = id // сохраняем сгенерированный ID в исходную структуру для последующего сравнения

	// get
	// получите только что добавленную посылку, убедитесь в отсутствии ошибки
	storedParcel, err := store.Get(id)
	if err != nil {
		t.Fatalf("ошибка при получении посылки: %v", err)
	}

	// проверьте, что значения всех полей в полученном объекте совпадают со значениями полей в переменной parcel
	if storedParcel.Number != parcel.Number ||
		storedParcel.Client != parcel.Client ||
		storedParcel.Status != parcel.Status ||
		storedParcel.Address != parcel.Address ||
		storedParcel.CreatedAt != parcel.CreatedAt {
		t.Errorf("полученные данные не совпадают с ожидаемыми.\nПолучено: %+v\nОжидалось: %+v", storedParcel, parcel)
	}

	// delete
	// удалите добавленную посылку, убедитесь в отсутствии ошибки
	err = store.Delete(id)
	if err != nil {
		t.Fatalf("ошибка при удалении посылки: %v", err)
	}

	// проверьте, что посылку больше нельзя получить из БД
	_, err = store.Get(id)
	if err == nil {
		t.Error("ожидалась ошибка при получении удаленной посылки, но ошибка не возвращена")
	}
}

// TestSetAddress проверяет обновление адреса
func TestSetAddress(t *testing.T) {
	// prepare
	
	db, err := sql.Open("sqlite", "tracker.db") 
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	store := NewParcelStore(db)
	parcel := getTestParcel()
	

	// add
	// Добавляем новую посылку в БД
//	initialAddress := "old test address"

	id, err := store.Add(parcel)
	if err != nil {
		t.Fatalf("ошибка при добавлении посылки: %v", err)
	}
	if id == 0 {
		t.Fatal("ожидался ненулевой идентификатор (ID) посылки")
	}


	// set address
	// Обновляем адрес, убеждаемся в отсутствии ошибки
	newAddress := "new test address"
	_, err = db.Exec("UPDATE parcel SET address = ? WHERE number = ?", newAddress, id)
	if err != nil {
		t.Fatalf("failed to update address: %v", err)
	}

	// check
	// Получаем добавленную посылку
	var p Parcel
	err = db.QueryRow("SELECT number, address FROM parcel WHERE number = ?", id).Scan(&p.Number, &p.Address)
	if err != nil {
		t.Fatalf("failed to select parcel: %v", err)
	}

	// Убеждаемся, что адрес обновился
	if p.Address != newAddress {
		t.Errorf("expected address %q, got %q", newAddress, p.Address)
	}
}

// TestSetStatus проверяет обновление статуса
func TestSetStatus(t *testing.T) {
	// prepare
	// Настраиваем подключение к тестовой БД (передавайте параметры через конфиг или ENV)
	db, err := sql.Open("sqlite", "tracker.db") 
	if err != nil {
		t.Fatalf("не удалось подключиться к БД: %v", err)
	}
	defer db.Close()

	// Убедимся, что соединение живое
	if err := db.Ping(); err != nil {
		t.Fatalf("БД недоступна: %v", err)
	}

	// Данные для теста
	
	targetStatus := "delivered"

	// add
	// Добавляем новую посылку в БД, возвращая её ID
	store := NewParcelStore(db)
	parcel := getTestParcel()

	// add
	// добавьте новую посылку в БД, убедитесь в отсутствии ошибки и наличии идентификатора
	id, err := store.Add(parcel)
	if err != nil {
		t.Fatalf("ошибка при добавлении посылки: %v", err)
	}
	if id == 0 {
		t.Fatal("ожидался ненулевой идентификатор (ID) посылки")
	}

	// set status
var actualParcel Parcel
	//selectQuery := `SELECT number, status  FROM parcel WHERE number = $1`
	//err = db.QueryRow(selectQuery, id).Scan(&actualParcel.Number, &actualParcel.Status)
//oldstatus:=actualParcel.Status
	// Обновляем статус посылки по её ID


	updateQuery := `UPDATE parcel SET status = $1  WHERE number = $2`
	_, err = db.Exec(updateQuery, targetStatus, id)
	if err != nil {
		t.Fatalf("не удалось обновить статус посылки: %v", err)
	}

	// check
	// Получаем добавленную посылку и убеждаемся, что статус обновился
	
	selectQuery := `SELECT number, status  FROM parcel WHERE number = $1`
	err = db.QueryRow(selectQuery, id).Scan(&actualParcel.Number, &actualParcel.Status)
	if err != nil {
		t.Fatalf("не удалось получить посылку для проверки: %v", err)
	}

	// Проверяем соответствие статуса
	if actualParcel.Status != targetStatus {
		t.Errorf("ожидался статус %q, но получен %q", targetStatus, actualParcel.Status)
	}
}

// TestGetByClient проверяет получение посылок по идентификатору клиента
func TestGetByClient(t *testing.T) {
	
	db, err := sql.Open("sqlite", "tracker.db") 
	assert.NoError(t, err)
	defer db.Close()

	store := NewParcelStore(db) 
	parcels := []Parcel{
		getTestParcel(),
		getTestParcel(),
		getTestParcel(),
	}

	parcelMap := map[int]Parcel{}

	
	client := randRange.Intn(10_000_000)
	parcels[0].Client = client
	parcels[1].Client = client
	parcels[2].Client = client

	
	for i := 0; i < len(parcels); i++ {
		
		id, err := store.Add(parcels[i])
		assert.NoError(t, err)
		assert.NotEmpty(t, id)

		
		parcels[i].Number = id

		
		parcelMap[id] = parcels[i]
	}

	
	storedParcels, err := store.GetByClient(client)
	
	assert.NoError(t, err)
	
	assert.Len(t, storedParcels, len(parcels))

	
	for _, parcel := range storedParcels {
		
		mappedParcel, exists := parcelMap[parcel.Number]
		assert.True(t, exists)

		
		assert.Equal(t, mappedParcel.Client, parcel.Client)
		assert.Equal(t, mappedParcel.Status, parcel.Status)
		assert.Equal(t, mappedParcel.Address, parcel.Address)
		assert.Equal(t, mappedParcel.CreatedAt, parcel.CreatedAt)
	}
}
