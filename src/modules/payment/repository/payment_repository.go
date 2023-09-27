package repository

import (
	"be-sagara-hackathon/src/modules/payment/model"
	um "be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common"
	"be-sagara-hackathon/src/utils/constants"
	e "be-sagara-hackathon/src/utils/errors"
	"database/sql"
	"fmt"
	"gorm.io/gorm"
	"math"
	"strings"
)

type PaymentRepository interface {
	Save(payment model.Payment) error
	Update(paymentID uint, payment model.Payment) error
	FindAll(
		filter model.FilterPayment,
		pg *utils.PaginateQueryOffset,
	) (payments []model.PaymentLite, totalData, totalPage int64, err error)
	FindOne(paymentID uint) (payment model.PaymentDetail, err error)
	FindUnprocessedByInvoiceID(invID uint) (payment model.Payment, err error)
	FindManyByInvoiceID(invID uint) (payments []model.PaymentLite, err error)
}

type PaymentRepositoryImpl struct {
	DB *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &PaymentRepositoryImpl{DB: db}
}

func (repository *PaymentRepositoryImpl) Save(payment model.Payment) error {
	tx := repository.DB.Begin()

	if err := tx.Model(&model.Payment{}).Omit("Invoice").Create(&payment).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&model.Invoice{}).Where("id=?", payment.InvoiceID).Updates(&model.Invoice{
		Status: constants.InvoiceProcessing,
		BaseEntity: common.BaseEntity{
			UpdatedAt: payment.UpdatedAt,
			UpdatedBy: payment.UpdatedBy,
		},
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&um.Participant{}).
		Where("id=?", payment.Invoice.ParticipantID).
		Updates(&um.Participant{
			PaymentStatus: constants.InvoiceProcessing,
			BaseEntity: common.BaseEntity{
				UpdatedAt: payment.UpdatedAt,
				UpdatedBy: payment.UpdatedBy,
			},
		}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (repository *PaymentRepositoryImpl) Update(paymentID uint, payment model.Payment) error {
	tx := repository.DB.Begin()
	if err := tx.Model(&model.Payment{}).Where("id=?", paymentID).Updates(map[string]interface{}{
		"amount":     payment.Amount,
		"note":       payment.Note,
		"status":     payment.Status,
		"proceed_at": payment.ProceedAt,
		"proceed_by": payment.ProceedBy,
		"updated_at": payment.UpdatedAt,
		"updated_by": payment.UpdatedBy,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&model.Invoice{}).Where("id=?", payment.InvoiceID).Updates(map[string]interface{}{
		"paid_amount": payment.Invoice.PaidAmount,
		"status":      payment.Invoice.Status,
		"approved_at": payment.Invoice.ApprovedAt,
		"approved_by": payment.Invoice.ApprovedBy,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Model(&um.Participant{}).Where("id=?", payment.Invoice.ParticipantID).Updates(map[string]interface{}{
		"payment_status": payment.Invoice.Status,
		"updated_at":     payment.UpdatedAt,
		"updated_by":     payment.UpdatedBy,
	}).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func (repository *PaymentRepositoryImpl) FindAll(
	filter model.FilterPayment,
	pg *utils.PaginateQueryOffset,
) (payments []model.PaymentLite, totalData, totalPage int64, err error) {
	where, whereVals := BuildFilterPayment(filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.Table("payments as py").
		Select(`py.id, inv.invoice_number, py.payment_type, pym.name as payment_method, py.created_at, py.status`).
		Joins("inner join invoices inv on inv.id = py.invoice_id").
		Joins("left join payment_methods pym on pym.id = py.payment_method_id").
		Order(fmt.Sprintf("%s %s", pg.Order.Field, pg.Order.By)).
		Limit(pg.Limit).Offset(pg.Offset).
		Where(buildWhereQuery, whereVals...).
		Find(&payments).Error; err != nil {
		return
	}

	totalData, err = repository.getTotalPayment(&filter)
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

func (repository *PaymentRepositoryImpl) getTotalPayment(filter *model.FilterPayment) (int64, error) {
	var (
		totalData int64
		err       error
	)

	where, whereVals := BuildFilterPayment(*filter)

	var buildWhereQuery string
	if where != nil {
		buildWhereQuery = strings.Join(where, " AND ")
	}

	if err = repository.DB.Table("payments as py").
		Joins("inner join invoices inv on inv.id = py.invoice_id").
		Joins("left join payment_methods pym on pym.id = py.payment_method_id").
		Where(buildWhereQuery, whereVals...).
		Count(&totalData).Error; err != nil {
		return 0, err
	}
	return totalData, nil
}

func BuildFilterPayment(filter model.FilterPayment) (where []string, whereVal []interface{}) {
	where = append(where, "inv.deleted_at IS NULL")

	if filter.Search != "" {
		filter.Search = strings.ToLower(filter.Search)
		where = append(where, "(LOWER(inv.invoice_number) LIKE @q OR LOWER(pym.name) LIKE @q)")
		whereVal = append(whereVal, sql.Named("q", "%"+filter.Search+"%"))
	}

	if filter.CreatedAt != "" {
		where = append(where, "date(py.created_at) = @created_at")
		whereVal = append(whereVal, sql.Named("created_at", filter.CreatedAt))
	}

	if filter.Status != "" {
		where = append(where, "py.status = @status")
		whereVal = append(whereVal, sql.Named("status", filter.Status))
	}

	if filter.PaymentType != "" {
		where = append(where, "py.payment_type = @type")
		whereVal = append(whereVal, sql.Named("type", filter.PaymentType))
	}

	return
}

func (repository *PaymentRepositoryImpl) FindOne(paymentID uint) (payment model.PaymentDetail, err error) {
	query := `
		SELECT pay.id, inv.event_id, pay.invoice_id, inv.invoice_number, inv.participant_id,
			u.name as participant_name, pay.payment_type, pay.payment_method_id, pm.name as payment_method_name,
			pay.created_at, pay.status, pay.bank_name, pay.bank_account_name, pay.bank_account_number, pay.evidence,
			pay.amount, pay.proceed_at, pay.proceed_by, pay.note
		FROM payments pay
		LEFT JOIN payment_methods pm on pm.id = pay.payment_method_id
		INNER JOIN invoices inv on inv.id = pay.invoice_id
		INNER JOIN participants p on p.id = inv.participant_id
		INNER JOIN users u on u.id = p.user_id
		WHERE pay.id = ?
		LIMIT 1
	`
	if err = repository.DB.Raw(query, paymentID).Scan(&payment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}

func (repository *PaymentRepositoryImpl) FindUnprocessedByInvoiceID(invID uint) (payment model.Payment, err error) {
	if err = repository.DB.
		Where("invoice_id=? AND proceed_at IS NULL", invID).
		First(&payment).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = e.ErrDataNotFound
		}
		return
	}
	return
}

func (repository *PaymentRepositoryImpl) FindManyByInvoiceID(invID uint) (payments []model.PaymentLite, err error) {
	if err = repository.DB.Table("payments as py").
		Select(`py.id, inv.invoice_number, py.payment_type, pym.name as payment_method, py.created_at, py.status`).
		Joins("inner join invoices inv on inv.id = py.invoice_id").
		Joins("left join payment_methods pym on pym.id = py.payment_method_id").
		Where("py.invoice_id=?", invID).
		Find(&payments).Error; err != nil {
		return
	}
	return
}
