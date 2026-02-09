package model

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// HasItem проверяет есть ли предмет у игрока
func (p *Player) HasItem(itemName string) bool {
	target := strings.ToLower(itemName)
	for _, item := range p.Inventory {
		if strings.ToLower(item.Name) == target {
			return true
		}
	}
	return false
}

// SetPassword сохраняет зашифрованный пароль
func (p *Player) SetPassword(plainPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)

	if err != nil {
		return err
	}

	p.Password = string(hash)
	return nil
}

// CheckPassword сравнивает зашифрованный пароль с паролем, введённым пользователем
func (p *Player) CheckPassword(plainPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(p.Password), []byte(plainPassword))

	return err == nil
}

// RenderItems выводит список предметов, которые сть у игрока
func RenderItems(items []Item, withStats bool) string {
	var totalCost int
	var finalString string

	for index, item := range items {
		finalString += fmt.Sprintf("%d. %s\n", index+1, item.Info())
		totalCost += item.Price
	}

	if withStats {
		finalString += fmt.Sprintln("---------------------------")
		finalString += fmt.Sprintf("ИТОГО: %d предметов на сумму %d золотых.\n\n", len(items), totalCost)
	}

	return finalString
}
