package repository

import (
	"be-sagara-hackathon/src/modules/payment/model"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common"
	e "be-sagara-hackathon/src/utils/errors"
	"be-sagara-hackathon/src/utils/helper"
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"math"
	"strings"
	"time"
)

type InvoiceRepository interface {
	Save(invoice model.Invoice) error
	Update(invoiceID uint, invoice model.Invoice) error
	UpdatePaidAmount(invoiceID uint, invoice model.Invoice) error
	UpdateStatus(invoiceID uint, status, updateBy string) error
	UpdateApprove(invoiceID uint, invoice model.Invoice) error
	FindAll(
		filter model.FilterInvoice,
		pg *utils.PaginateQueryOffset,
	) (invoices []model.InvoiceLite, totalData, totalPage int64, err error)
	FindOne(invoiceID uint) (model.InvoiceFull, error)
	FindManyByEventID(eventID uint) ([]model.InvoiceFull, error)
	FindManyByParticipantID(participantID uint) ([]model.InvoiceFull, error)
	FindByParticipantIDAndEventID(participantID, eventID uint) (model.InvoiceFull, error)
	FindByInvoiceNumber(invNumber string) (model.InvoiceFull, error)
}

type InvoiceRepositoryImpl struct {
	DB *gorm.DB
}

func NewInvoiceRepository(db *gorm.DB) InvoiceRepository {
	return &InvoiceRepositoryImpl{DB: db}
}

func (repository *InvoiceRepositoryImpl) Save(invoice model.Invoice) error {
	if err := repository.DB.Create(&invoice).Error; err != nil {
		return err
	}
	return nil
}

func (repository *InvoiceRepositoryImpl) Update(invoiceID uint, invoice model.Invoice) error {
	query := `
		UPDATE invoices SET paid_amount=?, status=?, approved_at=?, approved_by=?, updated_by=?
		WHERE id=?
	`
	if err := repository.DB.Exec(
		query, invoice.PaidAmount, invoice.Status, invoice.ApprovedAt,
		invoice.ApprovedBy, invoice.UpdatedBy, invoiceID,
	).Error; err != nil {
		return err
	}
	return nil
}

func (repository *InvoiceRepositoryImpl) UpdatePaidAmount(invoiceID uint, invoice model.Invoice) error {
	result := repository.DB.Where("id=?", invoiceID).Updates(&model.Invoice{
		PaidAmount: invoice.PaidAmount,
		BaseEntity: common.BaseEntity{
			UpdatedBy: invoice.UpdatedBy,
		},
	})
	return result.Error
}

func (repository *InvoiceRepositoryImpl) UpdateStatus(invoiceID uint, status, updateBy string) error {
	result := repository.DB.Where("id=?", invoiceID).Updates(&model.Invoice{
		Status: status,
		BaseEntity: common.BaseEntity{
			UpdatedBy: updateBy,
		},
	})
	return result.Error
}

func (repository *InvoiceRepositoryImpl) UpdateApprove(invoiceID uint, invoice model.Invoice) error {
	result := repository.DB.Where("id=?", invoiceID).Updates(&model.Invoice{
		Status:     invoice.Status,
		ApprovedBy: invoice.ApprovedBy,
		ApprovedAt: helper.ReferTime(time.Now()),
		BaseEntity: common.BaseEntity{
			UpdatedBy: invoice.UpdatedBy,
		},
	})
	return result.Error
}

func (repository *InvoiceRepositoryImpl) FindAll(
	filter model.FilterInvoice,
	pg *utils.PaginateQueryOffset,
) (invoices []model.InvoiceLite, totalData, totalPage int64, err error) {
	where, whereVals := BuildFilterInvoice(filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.Table("invoices as inv").
		Select(`inv.id, inv.event_id, inv.invoice_number, inv.amount, inv.status, u.name as participant_name`).
		Joins("inner join participants p on p.id = inv.participant_id").
		Joins("inner join users u on u.id = p.user_id").
		Order(fmt.Sprintf("%s %s", pg.Order.Field, pg.Order.By)).
		Limit(pg.Limit).Offset(pg.Offset).
		Where(buildWhereQuery, whereVals...).
		Find(&invoices).Error; err != nil {
		return
	}

	totalData, err = repository.getTotalInvoice(&filter)
	if err != nil {
		return
	}

	if pg.Limit > 0 {
		totalPage = int64(math.Ceil(float64(totalData) / float64(pg.Limit)))
	} else {
		totalPage = 1
	}

	return
}

func (repository *InvoiceRepositoryImpl) getTotalInvoice(filter *model.FilterInvoice) (int64, error) {
	var (
		totalData int64
		err       error
	)

	where, whereVals := BuildFilterInvoice(*filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.Table("invoices as inv").
		Joins("inner join participants p on p.id = inv.participant_id").
		Joins("inner join users u on u.id = p.user_id").
		Where(buildWhereQuery, whereVals...).
		Count(&totalData).Error; err != nil {
		return 0, err
	}
	return totalData, nil
}

func BuildFilterInvoice(filter model.FilterInvoice) (where []string, whereVal []interface{}) {
	where = append(where, "inv.deleted_at IS NULL")

	if filter.Search != "" {
		filter.Search = strings.ToLower(filter.Search)
		where = append(where, "(LOWER(inv.invoice_number) LIKE @q OR LOWER(u.name) LIKE @q OR inv.amount LIKE @q)")
		whereVal = append(whereVal, sql.Named("q", "%"+filter.Search+"%"))
	}

	if filter.EventID != 0 {
		where = append(where, "inv.event_id = @event_id")
		whereVal = append(whereVal, sql.Named("event_id", filter.EventID))
	}

	if filter.Status != "" {
		where = append(where, "inv.status = @status")
		whereVal = append(whereVal, sql.Named("status", filter.Status))
	}

	return
}

func (repository *InvoiceRepositoryImpl) FindOne(invoiceID uint) (model.InvoiceFull, error) {
	query := `
		SELECT inv.id, inv.event_id, e.name as event_name, inv.invoice_number, inv.participant_id,
			u.name as participant_name, u.email as participant_email, u.phone_number as participant_phone,
			inv.status, inv.amount, inv.paid_amount, inv.approved_at, inv.approved_by
		FROM invoices inv
		INNER JOIN participants p on p.id = inv.participant_id
		INNER JOIN users u on u.id = p.user_id
		INNER JOIN events e on e.id = inv.event_id
		WHERE inv.id = ? LIMIT 1
	`
	var invoice model.InvoiceFull
	if err := repository.DB.Raw(query, invoiceID).Scan(&invoice).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return invoice, err
	}
	return invoice, nil
}

func (repository *InvoiceRepositoryImpl) FindManyByEventID(eventID uint) ([]model.InvoiceFull, error) {
	query := `
		SELECT inv.id, inv.event_id, e.name as event_name, inv.invoice_number, inv.participant_id,
			u.name as participant_name, u.email as participant_email, u.phone_number as participant_phone,
			inv.status, inv.amount, inv.paid_amount, inv.approved_at, inv.approved_by
		FROM invoices inv
		INNER JOIN participants p on p.id = inv.participant_id
		INNER JOIN users u on u.id = p.user_id
		INNER JOIN events e on e.id = inv.event_id
		WHERE inv.event_id = ?
	`
	var invoices []model.InvoiceFull
	if err := repository.DB.Raw(query, eventID).Scan(&invoices).Error; err != nil {
		return invoices, err
	}
	return invoices, nil
}

func (repository *InvoiceRepositoryImpl) FindManyByParticipantID(participantID uint) ([]model.InvoiceFull, error) {
	query := `
		SELECT inv.id, inv.event_id, e.name as event_name, inv.invoice_number, inv.participant_id,
			u.name as participant_name, u.email as participant_email, u.phone_number as participant_phone,
			inv.status, inv.amount, inv.paid_amount, inv.approved_at, inv.approved_by
		FROM invoices inv
		INNER JOIN participants p on p.id = inv.participant_id
		INNER JOIN users u on u.id = p.user_id
		INNER JOIN events e on e.id = inv.event_id
		WHERE inv.participant_id = ?
	`
	var invoices []model.InvoiceFull
	if err := repository.DB.Raw(query, participantID).Scan(&invoices).Error; err != nil {
		return invoices, err
	}
	return invoices, nil
}

func (repository *InvoiceRepositoryImpl) FindByParticipantIDAndEventID(participantID, eventID uint) (model.InvoiceFull, error) {
	query := `
		SELECT inv.id, inv.event_id, e.name as event_name, inv.invoice_number, inv.participant_id,
			u.name as participant_name, u.email as participant_email, u.phone_number as participant_phone,
			inv.status, inv.amount, inv.paid_amount, inv.approved_at, inv.approved_by
		FROM invoices inv
		INNER JOIN participants p on p.id = inv.participant_id
		INNER JOIN users u on u.id = p.user_id
		INNER JOIN events e on e.id = inv.event_id
		WHERE inv.participant_id = ? AND inv.event_id = ?
		LIMIT 1
	`
	var invoice model.InvoiceFull
	if err := repository.DB.Raw(query, participantID, eventID).Scan(&invoice).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return invoice, err
	}
	return invoice, nil
}

func (repository *InvoiceRepositoryImpl) FindByInvoiceNumber(invNumber string) (model.InvoiceFull, error) {
	query := `
		SELECT inv.id, inv.event_id, e.name as event_name, inv.invoice_number, inv.participant_id,
			u.name as participant_name, u.email as participant_email, u.phone_number as participant_phone,
			inv.status, inv.amount, inv.paid_amount, inv.approved_at, inv.approved_by
		FROM invoices inv
		INNER JOIN participants p on p.id = inv.participant_id
		INNER JOIN users u on u.id = p.user_id
		INNER JOIN events e on e.id = inv.event_id
		WHERE inv.invoice_number=? LIMIT 1
	`
	var invoice model.InvoiceFull
	if err := repository.DB.Raw(query, invNumber).Scan(&invoice).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return invoice, err
	}
	return invoice, nil
}
