package web

import (
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
)

type params struct {
	pageSize  int
	updatedLT time.Time
	updatedGT time.Time
}

func parseParams(c *fiber.Ctx) (*params, error) {
	pageSize := c.QueryInt("page_size", 100)
	if pageSize < 0 || pageSize > 1000 {
		return nil, fmt.Errorf("%w: invalid page_size %d", fiber.ErrBadRequest, pageSize)
	}

	updatedLT, err := parseISO8601Datetime(c.Query("updated__lt"))
	if err != nil {
		return nil, fmt.Errorf("%w: invalid updated__lt %q", fiber.ErrBadRequest, c.Query("updated__lt"))
	}

	updatedGT, err := parseISO8601Datetime(c.Query("updated__gt"))
	if err != nil {
		return nil, fmt.Errorf("%w: invalid updated__gt %q", fiber.ErrBadRequest, c.Query("updated__gt"))
	}

	return &params{
		pageSize:  pageSize,
		updatedLT: updatedLT,
		updatedGT: updatedGT,
	}, nil
}
