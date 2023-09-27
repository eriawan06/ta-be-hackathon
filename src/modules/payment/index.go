package payment

import (
	"be-sagara-hackathon/src/modules/payment/controller"
	"be-sagara-hackathon/src/modules/payment/repository"
	"be-sagara-hackathon/src/modules/payment/service"
	ur "be-sagara-hackathon/src/modules/user/repository"
	"gorm.io/gorm"
)

var (
	paymentMethodRepository repository.PaymentMethodRepository
	paymentMethodService    service.PaymentMethodService
	paymentMethodController controller.PaymentMethodController
	invoiceRepository       repository.InvoiceRepository
	invoiceService          service.InvoiceService
	invoiceController       controller.InvoiceController
	paymentRepository       repository.PaymentRepository
	paymentService          service.PaymentService
	paymentController       controller.PaymentController
)

type Module interface {
	InitModule()
}

type ModuleImpl struct {
	DB *gorm.DB
}

func New(db *gorm.DB) Module {
	return &ModuleImpl{DB: db}
}

func (module ModuleImpl) InitModule() {
	participantRepository := ur.NewParticipantRepository(module.DB)

	paymentMethodRepository = repository.NewPaymentMethodRepository(module.DB)
	paymentMethodService = service.NewPaymentMethod(paymentMethodRepository)
	paymentMethodController = controller.NewPaymentMethodController(paymentMethodService)

	invoiceRepository = repository.NewInvoiceRepository(module.DB)
	invoiceService = service.NewInvoiceService(invoiceRepository)
	invoiceController = controller.NewInvoiceController(invoiceService)

	paymentRepository = repository.NewPaymentRepository(module.DB)
	paymentService = service.NewPaymentService(
		paymentRepository, invoiceRepository, participantRepository)
	paymentController = controller.NewPaymentController(paymentService)
}

func GetPaymentMethodController() controller.PaymentMethodController {
	return paymentMethodController
}

func GetInvoiceController() controller.InvoiceController {
	return invoiceController
}

func GetPaymentController() controller.PaymentController {
	return paymentController
}
