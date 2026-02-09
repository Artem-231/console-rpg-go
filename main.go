package main

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"awesomeProject/internal/model"
	"awesomeProject/internal/server"
	"awesomeProject/internal/storage"

	"github.com/joho/godotenv"
)

// main инициализирует настройки, подключение к БД и запускает игровой цикл.
func main() {
	rand.Seed(time.Now().UnixNano())

	if err := godotenv.Load(".env"); err != nil {
		fmt.Println("Внимание: .env файл не найден, используются стандартные настройки")
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	storage.InitDB()

	var (
		mainChoice  int8
		inputBuffer string
	)

	for {
	serverLoop:
		for {
			fmt.Printf("1. Играть в консоли\n2. Запустить веб-сервер\n3. Выйти из приложения\n> ")
			if _, err := fmt.Scan(&mainChoice); err != nil {
				fmt.Println("Введите цифру!")
				fmt.Scan(&inputBuffer)
				continue
			}

			switch mainChoice {
			case 1:
				break serverLoop
			case 2:
				server.StartServer(port)
				break serverLoop
			case 3:
				return
			default:
				fmt.Println("Нет такой команды!")
			}
		}

	gameLoop:
		for {
			var (
				choice   int8
				name     string
				password string
				player   model.Player
			)

		authLoop:
			for {
				fmt.Printf("1. Войти\n2. Регистрация\n> ")
				if _, err := fmt.Scan(&choice); err != nil {
					fmt.Println("Введите цифру!")
					fmt.Scan(&inputBuffer)
					continue
				}

				switch choice {
				case 1:
				loginInputLoop:
					for {
						fmt.Printf("Введите имя (или 'Назад')\n> ")
						fmt.Scan(&name)

						if strings.ToLower(name) == "назад" {
							break loginInputLoop
						}

						if len(name) == 0 {
							fmt.Println("Укажите имя!")
							continue
						}
						if !storage.HasName(name) {
							fmt.Println("Такого игрока не существует!")
							continue
						}

						for {
							fmt.Printf("Введите пароль (или 'Назад')\n> ")
							fmt.Scan(&password)

							if strings.ToLower(password) == "назад" {
								break loginInputLoop
							}

							loadedPlayer, err := storage.LoadGame(name)
							if err != nil {
								fmt.Println("Ошибка чтения данных игрока.")
								break loginInputLoop
							}

							if !loadedPlayer.CheckPassword(password) {
								fmt.Println("Неверный пароль!")
							} else {
								player = loadedPlayer
								fmt.Println(">>> Успешный вход!")
								break authLoop
							}
						}
					}

				case 2:
				regInputLoop:
					for {
						fmt.Printf("Введите имя для регистрации (или 'Назад')\n> ")
						fmt.Scan(&name)

						if strings.ToLower(name) == "назад" {
							break regInputLoop
						}

						if len(name) == 0 {
							fmt.Println("Имя не может быть пустым!")
						} else if storage.HasName(name) {
							fmt.Println("Такое имя уже занято")
						} else {
							for {
								fmt.Printf("Придумайте пароль (мин. 4 символа)\n> ")
								fmt.Scan(&password)

								if strings.ToLower(password) == "назад" {
									break regInputLoop
								}

								if len(password) < 4 {
									fmt.Println("Пароль слишком короткий!")
									continue
								}

								player = model.Player{
									Name:      name,
									Gold:      1000,
									Equipment: make(map[string]string),
									Inventory: []model.Item{{"Деревянный меч", 20}},
								}

								if err := player.SetPassword(password); err != nil {
									fmt.Println("Ошибка обработки пароля")
									return
								}
								storage.SaveGame(player)
								fmt.Println(">>> Регистрация успешна! Вы вошли в игру.")
								break authLoop
							}
						}
					}

				default:
					fmt.Println("Нет такого варианта.")
				}
			}

			// Основной игровой цикл
			for {
				fmt.Println("\n=========================")
				fmt.Println(player.Info())
				fmt.Println("1. Рюкзак\n2. Надеть снаряжение\n3. Магазин\n4. Лес\n5. Выход в главное меню")
				fmt.Printf("> ")
				fmt.Scan(&choice)

				switch choice {
				case 1:
					for {
						var exitCmd string
						fmt.Printf("Ваш инвентарь:\n%s", model.RenderItems(player.Inventory, true))
						fmt.Printf("Введите 'Назад' для возврата\n> ")
						fmt.Scan(&exitCmd)
						if strings.ToLower(exitCmd) == "назад" {
							break
						}
					}

				case 2:
					var slot int8
					for {
						fmt.Printf("Куда надеть предмет?\n1. Голова\n2. Тело\n3. Руки\n4. Назад\n> ")
						if _, err := fmt.Scan(&slot); err != nil {
							fmt.Scan(&inputBuffer)
							continue
						}

						if slot == 4 {
							break
						}

						var slotName string
						switch slot {
						case 1:
							slotName = "Голова"
						case 2:
							slotName = "Тело"
						case 3:
							slotName = "Руки"
						default:
							fmt.Println("Неверный слот")
							continue
						}

						fmt.Printf("Ваши вещи:\n%s", model.RenderItems(player.Inventory, false))
						exitNum := int8(len(player.Inventory)) + 1
						fmt.Printf("%d. Отмена\nВыберите номер предмета > ", exitNum)

						var itemIdx int8
						if _, err := fmt.Scan(&itemIdx); err != nil {
							fmt.Scan(&inputBuffer)
							continue
						}

						if itemIdx == exitNum {
							continue
						}

						realIdx := itemIdx - 1
						if realIdx < 0 || int(realIdx) >= len(player.Inventory) {
							fmt.Println("Нет предмета с таким номером!")
						} else {
							itemName := player.Inventory[realIdx].Name
							player.Equipment[slotName] = itemName
							fmt.Printf(">>> Вы надели: %s (%s)\n", itemName, slotName)
							storage.SaveGame(player)
						}
					}

				case 3:
				shopLoop:
					for {
						var shopChoice int8
						fmt.Printf("Магазин:\n1. Купить\n2. Продать\n3. Назад\n> ")
						if _, err := fmt.Scan(&shopChoice); err != nil {
							fmt.Scan(&inputBuffer)
							continue
						}

						switch shopChoice {
						case 1:
							var itemNum int
							fmt.Println("\nТовары:")
							fmt.Print(model.RenderItems(model.GlobalShop, false))
							fmt.Printf("Золото: %d\nВведите номер товара (0 для выхода)\n> ", player.Gold)

							if _, err := fmt.Scan(&itemNum); err != nil {
								fmt.Scan(&inputBuffer)
								continue
							}
							if itemNum == 0 {
								continue
							}

							idx := itemNum - 1
							if idx < 0 || idx >= len(model.GlobalShop) {
								fmt.Println("Нет такого товара")
								continue
							}

							if err := player.BuyItem(model.GlobalShop[idx]); err != nil {
								fmt.Println("Ошибка:", err)
							} else {
								fmt.Println("Куплено:", model.GlobalShop[idx].Name)
								storage.SaveGame(player)
							}

						case 2:
							if len(player.Inventory) == 0 {
								fmt.Println("Инвентарь пуст.")
								continue
							}
							fmt.Print(model.RenderItems(player.Inventory, false))
							fmt.Printf("Золото: %d\nЧто продать? (номер, 0 для выхода)\n> ", player.Gold)

							var sellNum int
							if _, err := fmt.Scan(&sellNum); err != nil {
								fmt.Scan(&inputBuffer)
								continue
							}
							if sellNum == 0 {
								continue
							}

							idx := sellNum - 1
							if idx < 0 || idx >= len(player.Inventory) {
								fmt.Println("Неверный номер")
							} else {
								item := player.Inventory[idx]
								player.Inventory = append(player.Inventory[:idx], player.Inventory[idx+1:]...)
								player.Gold += item.Price
								fmt.Printf("Продано: %s (+%d золота)\n", item.Name, item.Price)
								storage.SaveGame(player)
							}

						case 3:
							break shopLoop
						}
					}

				case 4:
					var riskChoice string
					fmt.Printf("Лес. Опасно. Испытать удачу? (да/нет)\n> ")
					fmt.Scan(&riskChoice)

					if strings.ToLower(riskChoice) == "да" {
						dice := rand.Intn(100) + 1
						goldFound := rand.Intn(50) + 10

						if dice <= 30 {
							player.Gold += goldFound
							fmt.Printf("Удача! Вы нашли %d золота.\n", goldFound)
						} else if dice <= 50 {
							// TODO: В будущем вынести логику шансов в отдельный пакет
							hasSword := false
							for _, item := range player.Inventory {
								if strings.Contains(strings.ToLower(item.Name), "меч") {
									hasSword = true
									break
								}
							}

							wolfDice := rand.Intn(100) + 1

							if hasSword && wolfDice <= 80 {
								fmt.Println("Волки напали, но вы отбились мечом!")
							} else {
								loss := goldFound
								if player.Gold < loss {
									loss = player.Gold
								}
								player.Gold -= loss
								fmt.Printf("Волки покусали вас! Вы потеряли %d золота при бегстве.\n", loss)
							}
						} else {
							fmt.Println("Вы погуляли по лесу, но ничего не нашли.")
						}
						storage.SaveGame(player)
					}

				case 5:
					break gameLoop

				default:
					fmt.Println("Неверная команда")
				}
			}
		}
	}
}
