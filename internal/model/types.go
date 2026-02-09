package model

import (
	"errors"
	"fmt"
)

type Item struct {
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type Player struct {
	Name      string            `json:"name"`
	Password  string            `json:"password"`
	Gold      int               `json:"gold"`
	Equipment map[string]string `json:"equipment"`
	Inventory []Item            `json:"inventory"`
}

type BuyRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
	ItemName string `json:"item_name"`
}

type CreatePlayer struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// Info выводит информацию о предмете
func (i Item) Info() string {
	return fmt.Sprintf("%s (цена: %d)", i.Name, i.Price)
}

// Info выводит информацию об игроке
func (p Player) Info() string {
	result := fmt.Sprintf("Золото: %d\n", p.Gold)
	result += "Надето на героя:\n"

	if len(p.Equipment) > 0 {
		for k, v := range p.Equipment {
			result += fmt.Sprintf("- %s: %s\n", k, v)
		}
	} else {
		result += "Ничего не надето.\n"
	}

	return result
}

// Pay проверяет возможность оплаты игроком
func (p *Player) Pay(amount int) bool {
	if amount <= p.Gold {
		p.Gold -= amount
		return true
	}
	return false
}

// BuyItem добавляет купленный предмет в инвентарь и списывает золото
func (p *Player) BuyItem(i Item) error {
	ifCanBuy := p.Pay(i.Price)

	if ifCanBuy {
		p.Inventory = append(p.Inventory, i)
		return nil
	} else {
		return errors.New("недостаточно золота")
	}

}
