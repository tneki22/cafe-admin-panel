package main

import (
	"backend/pkg/vars"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	minPassword int
	maxPassword int
	minUsername int
	maxUsername int

	server  *echo.Echo
	r       *echo.Group
	address string

	uc Usecase
}

func NewServer(ip string, port int, minPassword, maxPassword, minUsername, maxUsername int, secret string, uc Usecase) *Server {
	api := Server{
		minPassword: minPassword,
		maxPassword: maxPassword,
		minUsername: minUsername,
		maxUsername: maxUsername,
		uc:          uc,
	}

	api.server = echo.New()
	api.server.Use(middleware.Logger())
	api.server.Use(middleware.Recover())
	api.server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"},
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAuthorization},
	}))

	api.r = api.server.Group("/api")

	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(vars.JWTClaims)
		},
		SigningKey: []byte(secret),
	}

	api.r.Use(echojwt.WithConfig(config))
	api.r.GET("", api.AuthMiddleware)
	api.RegisterRoutes()
	api.address = fmt.Sprintf("%s:%d", ip, port)

	return &api
}
func (api *Server) RegisterRoutes() {
	api.server.POST("/api/register", api.Register)
	api.server.POST("/api/login", api.Login)
	api.r.GET("/profile", api.AuthMiddleware)
	api.server.GET("/api/menu", api.GetMenu)
	api.server.DELETE("/api/menu/:id", api.DeleteMenuItem)
	api.server.PUT("/api/menu/:id", api.UpdateMenuItem)
	api.server.POST("/api/menu", api.AddMenuItem) // Новый маршрут
	api.server.POST("/api/orders", api.AddOrder)  // Новый маршрут
	api.server.GET("/api/orders", api.GetOrders)  // Новый маршрут
	api.server.PUT("/api/orders/:id/status", api.UpdateOrderStatus)
	api.server.GET("/api/revenue", api.GetRevenue)
	api.server.GET("/api/order_counts", api.GetOrderCounts) // Новый маршрут

}
func (api *Server) Run() {
	api.server.Logger.Fatal(api.server.Start(api.address))
}

func (srv *Server) Register(c echo.Context) error {
	user := struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := c.Bind(&user)
	if err != nil {
		log.Printf("Error binding data: %v", err)
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if user.Name == "" || user.Password == "" || user.Email == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	if len(user.Name) < srv.minUsername || len(user.Name) > srv.maxUsername {
		return echo.NewHTTPError(http.StatusUnauthorized, "Username should be "+strconv.Itoa(srv.minUsername)+"-"+strconv.Itoa(srv.maxUsername)+" length")
	}

	if len(user.Password) < srv.minPassword || len(user.Password) > srv.maxPassword {
		return echo.NewHTTPError(http.StatusUnauthorized, "Password should be "+strconv.Itoa(srv.minPassword)+"-"+strconv.Itoa(srv.maxPassword)+" length")
	}

	err = srv.uc.Register(user.Name, user.Email, user.Password)
	if err != nil {
		log.Printf("Error registering user: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to create account")
	}

	return c.JSON(http.StatusOK, "OK!")
}

func (srv *Server) Login(c echo.Context) error {
	user := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := c.Bind(&user)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	if user.Email == "" || user.Password == "" {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid credentials")
	}

	token, err := srv.uc.Authenticate(user.Email, user.Password)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, echo.Map{"token": token})
}

func (srv *Server) AuthMiddleware(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	if user == nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Token is missing or invalid"})
	}
	claims := user.Claims.(*vars.JWTClaims)
	log.Printf("Claims: %v", claims)
	username := claims.Username
	log.Printf("Username: %v", username)
	return c.JSON(http.StatusOK, map[string]string{"message": "Welcome " + username})
}

func (srv *Server) GetMenu(c echo.Context) error {
	menuItems, err := srv.uc.GetMenuItems()
	if err != nil {
		log.Printf("Error fetching menu items: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch menu items")
	}

	return c.JSON(http.StatusOK, menuItems)
}

func (srv *Server) DeleteMenuItem(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Неверный ID")
	}

	err = srv.uc.DeleteMenuItem(id)
	if err != nil {
		log.Printf("Error deleting menu item: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Не удалось удалить элемент меню")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Элемент меню удален"})
}

func (srv *Server) UpdateMenuItem(c echo.Context) error {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Неверный ID")
	}

	var input struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Неверные данные")
	}

	item := MenuItem{
		ID:          id,
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
	}

	updatedItem, err := srv.uc.UpdateMenuItem(item)
	if err != nil {
		log.Printf("Error updating menu item: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Не удалось обновить элемент меню")
	}

	return c.JSON(http.StatusOK, updatedItem)
}
func (srv *Server) AddMenuItem(c echo.Context) error {
	var input struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
	}
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Неверные данные")
	}

	item := MenuItem{
		Name:        input.Name,
		Description: input.Description,
		Price:       input.Price,
	}

	newItem, err := srv.uc.AddMenuItem(item)
	if err != nil {
		log.Printf("Error adding menu item: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Не удалось добавить элемент меню")
	}

	return c.JSON(http.StatusOK, newItem)
}
func (srv *Server) AddOrder(c echo.Context) error {
	var input struct {
		Items []struct {
			MenuItemId int `json:"menuItemId"`
			Quantity   int `json:"quantity"`
		} `json:"items"`
	}

	// Привязка входящих данных
	if err := c.Bind(&input); err != nil {
		log.Printf("Error binding order data: %v", err)
		return echo.NewHTTPError(http.StatusBadRequest, "Неверные данные заказа")
	}

	// Проверка наличия элементов
	if len(input.Items) == 0 {
		return echo.NewHTTPError(http.StatusBadRequest, "Заказ должен содержать хотя бы один товар")
	}

	orderItems := make([]OrderItem, len(input.Items))
	for i, item := range input.Items {
		if item.MenuItemId <= 0 || item.Quantity <= 0 {
			return echo.NewHTTPError(http.StatusBadRequest, "Некорректные данные товара в заказе")
		}
		orderItems[i] = OrderItem{
			MenuItemId: item.MenuItemId,
			Quantity:   item.Quantity,
		}
	}

	newOrder, err := srv.uc.AddOrder(orderItems)
	if err != nil {
		log.Printf("Error adding order: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Не удалось добавить заказ")
	}

	return c.JSON(http.StatusOK, newOrder)
}
func (srv *Server) GetOrders(c echo.Context) error {
	orders, err := srv.uc.GetOrders()
	if err != nil {
		log.Printf("Error fetching orders: %v", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "Не удалось загрузить заказы")
	}

	return c.JSON(http.StatusOK, orders)
}
func (srv *Server) UpdateOrderStatus(c echo.Context) error {
	idParam := c.Param("id")
	orderID, err := strconv.Atoi(idParam)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid order ID")
	}

	status := struct {
		Status string `json:"status"`
	}{}
	if err := c.Bind(&status); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid status data")
	}

	if status.Status != "Выполнен" && status.Status != "Отменен" {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid status value")
	}

	err = srv.uc.UpdateOrderStatus(orderID, status.Status)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to update order status")
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Статус заказа обновлен"})
}
func (srv *Server) GetRevenue(c echo.Context) error {
	period := c.QueryParam("period")
	if period != "day" && period != "week" && period != "month" && period != "year" {
		return echo.NewHTTPError(http.StatusBadRequest, "Недопустимый параметр period")
	}

	revenueData, err := srv.uc.GetRevenue(period)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Ошибка получения данных выручки")
	}

	return c.JSON(http.StatusOK, revenueData)
}
func (srv *Server) GetOrderCounts(c echo.Context) error {
	period := c.QueryParam("period")
	if period == "" {
		period = "day" // По умолчанию "day"
	}
	if period != "day" && period != "week" && period != "month" && period != "year" {
		return echo.NewHTTPError(http.StatusBadRequest, "Недопустимый параметр period")
	}

	orderCounts, err := srv.uc.GetOrderCounts(period)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Ошибка получения данных количества заказов")
	}

	return c.JSON(http.StatusOK, orderCounts)
}
