package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"awesomeProject/internal/model"
	"awesomeProject/internal/storage"
)

// handleMain отображает главную страницу
func handleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Добро пожаловать в РПГ игру!\nПерейдите на адрес '/status?name='свой ник'', чтобы увидеть свой актуальный статус.")
}

// handleStatus отображает страницу со статусом игрока, а именно его данные из таблицы players
func handleStatus(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	password := r.URL.Query().Get("password")
	if name == "" {
		http.Error(w, "Укажите имя!", http.StatusBadRequest)
	} else if password == "" {
		http.Error(w, "Укажите пароль!", http.StatusBadRequest)
	} else {
		p, err := storage.LoadGame(name)

		if err != nil {
			http.Error(w, "Игрок не найден", http.StatusNotFound)
			return
		}

		if !p.CheckPassword(password) {
			http.Error(w, "Неверный пароль", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(p)
	}
}

// handleBuy позволяет через метод POST передать название предмета, который ты хочешь купить
func handleBuy(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только POST запрос", http.StatusMethodNotAllowed)
		return
	}

	var req model.BuyRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Кривой JSON", http.StatusBadRequest)
		return
	}

	p, err := storage.LoadGame(req.Name)
	if err != nil {
		http.Error(w, "Игрок не найден", http.StatusNotFound)
		return
	}

	if !p.CheckPassword(req.Password) {
		http.Error(w, "Неверный пароль", http.StatusUnauthorized)
		return
	}

	var itemToBuy model.Item
	found := false
	for _, item := range model.GlobalShop {
		if item.Name == req.ItemName {
			itemToBuy = item
			found = true
			break
		}
	}

	if !found {
		http.Error(w, "Товар не найден в магазине", http.StatusNotFound)
		return
	}

	err = p.BuyItem(itemToBuy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	storage.SaveGame(p)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// handleCreate позволяет через метод POST передать данные для создания пользователя
func handleCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Только Post запрос", http.StatusMethodNotAllowed)
		return
	}

	var req model.CreatePlayer
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		http.Error(w, "Кривой JSON", http.StatusBadRequest)
		return
	}
	if storage.HasName(req.Name) {
		http.Error(w, "Игрок уже существует", http.StatusConflict)
		return
	}

	p := model.Player{req.Name, "", 1000, make(map[string]string), []model.Item{{"Деревянный меч", 20}}}
	err = p.SetPassword(req.Password)

	if err != nil {
		http.Error(w, "Ошибка при обработке пароля", http.StatusInternalServerError)
	}
	storage.SaveGame(p)
	w.WriteHeader(http.StatusCreated)
	fmt.Fprint(w, "Игрок успешно создан!")
}

// StartServer запускает все обработчики и поддерживает их на определённом порту
func StartServer(port string) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handleMain)
	mux.HandleFunc("/status", handleStatus)
	mux.HandleFunc("/buy", handleBuy)
	mux.HandleFunc("/create", handleCreate)

	handler := LoggerMiddleware(mux)

	fmt.Println("Сервер запущен на порту", port)

	http.ListenAndServe(port, handler)

}
