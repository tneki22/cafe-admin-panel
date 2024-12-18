package main

import (
	"backend/pkg/vars"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func (u *Usecase) Authenticate(email, password string) (string, error) {
	exist, err := u.p.CheckUserByEmail(email)
	if !exist {
		return "", errors.New("user not found")
	}
	if err != nil {
		return "", err
	}

	name, hashedPassword, err := u.p.GetUsernameAndHashedPassword(email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	return u.jp.GenerateToken(name)
}

func (u *Usecase) ValidateJWT(token string) (*vars.JWTClaims, error) {
	return u.jp.ValidateToken(token)
}

func (u *Usecase) Register(name, email, password string) error {
	exist, err := u.p.CheckUserByEmail(email)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("user already exists")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}
	return u.p.CreateUser(name, email, string(hashedPassword))
}

type Usecase struct {
	defaultMsg string

	p  Provider
	jp JWTProvider
}

func NewUsecase(defaultMsg string, p Provider, jp JWTProvider) *Usecase {
	return &Usecase{
		defaultMsg: defaultMsg,
		p:          p,
		jp:         jp,
	}
}
func (u *Usecase) GetMenuItems() ([]MenuItem, error) {
	return u.p.FetchMenuItems()
}

type MenuItem struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	CreatedAt   string  `json:"created_at"`
}

func (u *Usecase) DeleteMenuItem(id int) error {
	return u.p.DeleteMenuItem(id)
}

func (u *Usecase) UpdateMenuItem(item MenuItem) (MenuItem, error) {
	return u.p.UpdateMenuItem(item)
}
func (u *Usecase) AddMenuItem(item MenuItem) (MenuItem, error) {
	return u.p.AddMenuItem(item)
}

type OrderItem struct {
	MenuItemId int `json:"menuItemId"`
	Quantity   int `json:"quantity"`
}

type Order struct {
	ID        int         `json:"id"`
	Total     float64     `json:"total"`
	Status    string      `json:"status"`
	CreatedAt string      `json:"created_at"`
	Items     []OrderItem `json:"items"`
}

func (u *Usecase) AddOrder(items []OrderItem) (Order, error) {
	return u.p.AddOrder(items)
}
func (u *Usecase) GetOrders() ([]Order, error) {
	return u.p.FetchOrders()
}
func (u *Usecase) UpdateOrderStatus(orderID int, status string) error {
	return u.p.UpdateOrderStatus(orderID, status)
}

type RevenueData struct {
	TimeUnit string  `json:"time_unit"`
	Total    float64 `json:"total"`
}

func (u *Usecase) GetRevenue(period string) ([]RevenueData, error) {
	return u.p.FetchRevenue(period)
}

type OrderCountData struct {
	TimeUnit string `json:"time_unit"`
	Count    int    `json:"count"`
}

// Метод для получения количества заказов
func (u *Usecase) GetOrderCounts(period string) ([]OrderCountData, error) {
	return u.p.FetchOrderCounts(period)
}
