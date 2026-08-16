package handler

import (
	"manajemen-keuangan-api/model"
	"manajemen-keuangan-api/service"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	svcTransaction *service.TransactionService
}

func NewTransactionHandler(svcTransaction *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{svcTransaction: svcTransaction}
}

func getUserID(c *fiber.Ctx) uint {
	return c.Locals("userID").(uint)
}

// endpoint: /api/v1/transaction {post}
func (h *TransactionHandler) CreateTransaction(c *fiber.Ctx) error {
	userid := getUserID(c)
	var inputTransaction struct {
		BalanceID uint    `json:"balanceID"`
		Type      model.BalanceType  `json:"type"`
		Amount    int64        `json:"amount"`
		Category  model.CategoryType `json:"category"`
		Description string	`json:"description"`
	}

	if err := c.BodyParser(&inputTransaction); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": "Format request tidak sesuai!",
		})
	}

	trx, err := h.svcTransaction.CreateTransaction(c.Context(), userid, inputTransaction.BalanceID, inputTransaction.Type, inputTransaction.Amount, inputTransaction.Category, inputTransaction.Description)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"data": trx,
	})
}


// endpoint: /api/v1/transaction?type=&category=&start_date=&end_date=&page=&limit= {get}
func (h *TransactionHandler) GetAllTransactions(c *fiber.Ctx) error {
	userid := getUserID(c)
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "20"))

	typeTrx := model.BalanceType(c.Query("type"))
	categoryTrx := model.CategoryType(c.Query("category"))

	var startDate, endDate time.Time
	var err error
	    if s := c.Query("start_date"); s != "" {
        startDate, err = time.Parse("2006-01-02", s)
        if err != nil {
            return c.Status(400).JSON(fiber.Map{
                "success":  false,
                "messagge": "format start_date harus YYYY-MM-DD",
            })
        }
    }
    if e := c.Query("end_date"); e != "" {
        endDate, err = time.Parse("2006-01-02", e)
        if err != nil {
            return c.Status(400).JSON(fiber.Map{
                "success":  false,
                "messagge": "format end_date harus YYYY-MM-DD",
            })
        }
    }

	trxs, err := h.svcTransaction.FilterDateTransaction(c.Context(), uint(userid), int(page), int(limit), typeTrx, categoryTrx, startDate, endDate)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"data": trxs,
	})
}


// endpoint: /api/v1/transaction/:id {get}
func (h *TransactionHandler) GetTransactionByIDandUserID(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": "Format id tidak sesuai",
		})
	}

	userid := getUserID(c)
	

	trx, err := h.svcTransaction.GetTransactionByIDandUserID(c.Context(), uint(id), uint(userid))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"data": trx,
	})
}


// endpoint: /api/v1/transaction/summary/start_date=?&end_date? {get}
func (h *TransactionHandler) GetSummary(c *fiber.Ctx) error {
	userid := getUserID(c)

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" ||  endDateStr == "" {
		return c.Status(404).JSON(fiber.Map{
			"status": false,
			"messagge": "Parameter start_date dan end_date wajib diisi",
		})
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": err.Error(),
		})
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": err.Error(),			
		})
	}

	summary, err := h.svcTransaction.GetSummary(c.Context(), uint(userid), startDate, endDate)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"data": summary,
	})
 }


// endpoint: /api/v1/transaction/:id {put}
func (h *TransactionHandler) UpdateTransaction(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": "Format id tidak sesuai",
		})
	}

	userid := getUserID(c)
	
	var inputUpdateTrx struct {
		Type model.BalanceType `json:"type"`
		Amount int64 `json:"amount"`
		Category model.CategoryType `json:"category"`
		Description string `json:"description"`
 	}

	if err := c.BodyParser(&inputUpdateTrx); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": "Format request tidak valid!",
		})
	}

	UpdateTrx := &model.Transaction{
		ID: uint(id),
		UserID: userid,
		Type: inputUpdateTrx.Type,
		Amount: inputUpdateTrx.Amount,
		Category: inputUpdateTrx.Category,
		Description: inputUpdateTrx.Description,
	}

	if err := h.svcTransaction.UpdateTransaction(c.Context(), UpdateTrx); err != nil {
		return c.Status(500).JSON(fiber.Map{
			"success": false,
			"messagge": err.Error(),
		})	
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"messagge": "Transaction berhasil di update!",
	})

}

// enpoint: /api/v1/transaction/:id {delete}
func (h *TransactionHandler) DeleteTransaction(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"success": false,
			"messagge": "Format id tidak sesuai",
		})
	}

	userid := getUserID(c)
	if err = h.svcTransaction.DeleteTransaction(c.Context(), uint(id), userid); err != nil {
		return c.Status(200).JSON(fiber.Map{
			"success": false,
			"messagge": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"success": true,
		"messagge": "Data transaction berhasil dihapus!",
	})
}
