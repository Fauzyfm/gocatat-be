package handler

import (
	"manajemen-keuangan-api/model"
	"manajemen-keuangan-api/service"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type BalanceHandler struct {
	svcBalance *service.BalanceService
}

func NewBalanceHandler(svcBalance *service.BalanceService) *BalanceHandler {
	return &BalanceHandler{svcBalance: svcBalance}
}


//endpoint: /api/v1/balance {post}
func (h *BalanceHandler) CreateBalance(c *fiber.Ctx) error {
	userID := getUserID(c)
	var inputBalance struct {
		Wallet string `json:"wallet"`
		Type model.BalanceType `json:"type"`
	}

	if err := c.BodyParser(&inputBalance); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": "Format request tidak valid!",
		})
	}

	balance, err := h.svcBalance.CreateBalance(c.Context(), userID, inputBalance.Wallet, inputBalance.Type)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"data": balance,
	})

}


// endpoint: /api/v1/balance {get}
func (h *BalanceHandler) GetAllBalanceByUserID(c *fiber.Ctx) error {
	userid := getUserID(c)

	balances, err := h.svcBalance.GetAllBalanceByUserID(c.Context(), uint(userid))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"data": balances,
	})

}


// endpoint: /api/v1/balance/:id {get}
func (h *BalanceHandler) GetBalanceByID(c *fiber.Ctx) error {
	userid := getUserID(c)

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": "Format id tidak sesuai",
		})
	}


	balance, err := h.svcBalance.GetBalanceByID(c.Context(), uint(id), userid)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"data": balance,
	})

}


// endpoint: /api/v1/balance/:id {put}
func (h *BalanceHandler) UpdateBalance(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": "Format id tidak sesuai",
		})
	}
	
	userid := getUserID(c)

	var inputBalance struct {
		Wallet string `json:"wallet"`
		Type model.BalanceType `json:"type"`
		UpdateAt time.Time `json:"updateAt"`
	}

	if err := c.BodyParser(&inputBalance); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": "Format request tidak valid!",
		})
	}

	NewBalanceUpdate := &model.Balance{
		ID: uint(id),
		UserID: userid,
		Wallet: inputBalance.Wallet,
		Type: inputBalance.Type,
		UpdateAt: time.Now(),
	}

	if err := h.svcBalance.UpdateBalance(c.Context(), NewBalanceUpdate); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"messagge": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"messagge": "Balance berhasil diupdate!",
	})
}


// endpoint: /api/v1/balance/:id {delete}
func (h *BalanceHandler) DeleteBalance(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": "Format id tidak sesuai",
		})
	}

	userid := getUserID(c)

	if err := h.svcBalance.DeleteBalance(c.Context(), uint(id), uint(userid)); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"messagge": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"messagge": "balance berhasil dihapus!",
	})
}
