package service

import (
	"be-sagara-hackathon/src/modules/auth/mapper"
	"be-sagara-hackathon/src/modules/auth/model"
	"be-sagara-hackathon/src/modules/auth/repository"
	evm "be-sagara-hackathon/src/modules/event/model"
	evr "be-sagara-hackathon/src/modules/event/repository"
	pym "be-sagara-hackathon/src/modules/payment/model"
	um "be-sagara-hackathon/src/modules/user/model"
	ur "be-sagara-hackathon/src/modules/user/repository"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common"
	"be-sagara-hackathon/src/utils/constants"
	"be-sagara-hackathon/src/utils/email"
	e "be-sagara-hackathon/src/utils/errors"
	"be-sagara-hackathon/src/utils/helper"
	"be-sagara-hackathon/src/utils/oauth"
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"google.golang.org/api/idtoken"
	"gorm.io/gorm"
)

type AuthService interface {
	Register(request model.RegisterRequest) error
	Login(request model.LoginRequest) (response model.AuthResponse, err error)
	GoogleOauth(request oauth.GoogleUserResult) (response model.AuthResponse, err error)
	RegisterByGoogle(request model.RegisterByGoogleRequest) (model.AuthResponse, error)
	LoginByGoogle(request model.LoginByGoogleRequest) (model.AuthResponse, error)
	VerifyEmail(request model.VerifyEmailRequest) error
	SendVerificationCode(request model.SendVerificationCodeRequest) error
	ValidateVerificationCode(request model.ValidateVerificationCodeRequest) error
	ForgotPassword(request model.ForgotPasswordRequest) error
}

type AuthServiceImpl struct {
	AuthRepository       repository.AuthRepository
	VerificationCodeRepo repository.VerificationCodeRepository
	UserRepo             ur.UserRepository
	ParticipantRepo      ur.ParticipantRepository
	EventRepository      evr.EventRepository
	UserRoleRepo         ur.UserRoleRepository
	EventParticipantRepo evr.EventParticipantRepository
}

func NewAuthService(
	authRepository repository.AuthRepository,
	verificationCodeRepo repository.VerificationCodeRepository,
	userRepo ur.UserRepository,
	participantRepo ur.ParticipantRepository,
	eventRepo evr.EventRepository,
	userRoleRepo ur.UserRoleRepository,
	eventParticipantRepo evr.EventParticipantRepository,
) AuthService {
	return &AuthServiceImpl{
		AuthRepository:       authRepository,
		VerificationCodeRepo: verificationCodeRepo,
		UserRepo:             userRepo,
		ParticipantRepo:      participantRepo,
		EventRepository:      eventRepo,
		UserRoleRepo:         userRoleRepo,
		EventParticipantRepo: eventParticipantRepo,
	}
}

func (service *AuthServiceImpl) Register(request model.RegisterRequest) error {
	if request.Password != request.ConfirmPassword {
		return e.ErrConfirmPasswordNotSame
	}

	//Find latest running event
	latestEvent, err := service.EventRepository.FindLatest()
	if err != nil {
		return err
	}
	if latestEvent.Status != constants.EventRunning {
		return e.ErrEventNotRunning
	}

	//Find user role participant
	role, err := service.UserRoleRepo.FindByName(constants.UserParticipant)
	if err != nil {
		return err
	}

	hashed, err := utils.HashPassword(request.Password)
	if err != nil {
		return err
	}

	request.Password = hashed
	verificationCode := utils.GenerateUuid()
	registerModel := model.RegisterModel{
		User: um.User{
			Name:        request.FullName,
			Email:       request.Email,
			PhoneNumber: &request.PhoneNumber,
			Password:    &request.Password,
			UserRoleID:  role.ID,
			AuthType:    constants.AuthTypeRegular,
			BaseEntity: common.BaseEntity{
				CreatedBy: "self",
				UpdatedBy: "self",
			},
		},
		LatestEventID: latestEvent.ID,
		Invoice: pym.Invoice{
			EventID:       latestEvent.ID,
			InvoiceNumber: utils.GenerateInvoiceNumber(),
			Status:        constants.InvoiceUnpaid,
			Amount:        latestEvent.RegFee,
			BaseEntity: common.BaseEntity{
				CreatedBy: "self",
				UpdatedBy: "self",
			},
		},
		Verification: &model.VerificationCode{
			Code: verificationCode,
			Type: constants.VerifCodeEmail,
			BaseEntity: common.BaseEntity{
				CreatedBy: "self",
				UpdatedBy: "self",
			},
		},
	}

	newUser, err := service.AuthRepository.Register(registerModel)
	if err != nil {
		return err
	}

	// TODO: email
	// - send email verification [OK]
	// - use goroutine/task queue/background task
	url := fmt.Sprintf("%s%s", os.Getenv("BASE_FE_URL"), os.Getenv("VERIFY_EMAIL_REDIRECT_URL"))
	templateData := email.TemplateData{
		Name:        newUser.Name,
		Link:        fmt.Sprintf("%s%s", url, verificationCode),
		SenderEmail: "sagarahackathon@gmail.com",
		Title:       constants.EmailSubjectVerifyEmail,
		Type:        constants.VerifCodeEmail,
	}

	r := email.NewRequest([]string{newUser.Email}, constants.EmailSubjectVerifyEmail, "")
	if err = r.ParseTemplate("./src/utils/email/template_email.html", templateData); err == nil {
		ok, err2 := r.SendEmail()
		if !ok || err2 != nil {
			return e.ErrFailedSendEmail
		}
	} else {
		return e.ErrFailedParseEmailTemplate
	}

	return nil
}

func (service *AuthServiceImpl) Login(request model.LoginRequest) (response model.AuthResponse, err error) {
	user, err := service.UserRepo.FindByEmail(request.Email)
	if err != nil {
		if err == e.ErrEmailNotRegistered {
			err = e.ErrWrongLoginCredential
		}
		return
	}

	if user.AuthType != constants.AuthTypeRegular {
		err = e.ErrWrongAuthMethod
		return
	}

	isPasswordValid := utils.CheckPasswordHash(request.Password, helper.DereferString(user.Password))
	if !isPasswordValid {
		err = e.ErrWrongLoginCredential
		return
	}

	if !user.IsActive {
		err = e.ErrUserIsNotActivated
		return
	}

	latestEvent, err := service.EventRepository.FindLatest()
	if err != nil {
		return
	}

	_, err = service.EventParticipantRepo.FindOneByEventIDAndParticipantID(latestEvent.ID, user.Participant.ID)
	if err != nil && err != e.ErrDataNotFound {
		return
	}

	if err != nil && err == e.ErrDataNotFound {
		evp := evm.EventParticipant{
			BaseEntity: common.BaseEntity{
				CreatedAt: time.Now(),
				CreatedBy: "self",
				UpdatedAt: time.Now(),
				UpdatedBy: "self",
			},
			EventID:       latestEvent.ID,
			ParticipantID: user.Participant.ID,
		}
		inv := pym.Invoice{
			BaseEntity: common.BaseEntity{
				CreatedAt: time.Now(),
				CreatedBy: "self",
				UpdatedAt: time.Now(),
				UpdatedBy: "self",
			},
			EventID:       latestEvent.ID,
			ParticipantID: user.Participant.ID,
			InvoiceNumber: utils.GenerateInvoiceNumber(),
			Amount:        latestEvent.RegFee,
			Status:        constants.InvoiceUnpaid,
		}
		if err = service.ParticipantRepo.UpdateRegistrationStatus(
			user.Participant.ID, false, &evp, &inv,
		); err != nil {
			return
		}
		user.Participant.IsRegistered = false
		user.Participant.PaymentStatus = constants.InvoiceUnpaid
	}

	// Generate Token
	token, err := utils.GenerateToken(user)
	if err != nil {
		return
	}

	response = model.AuthResponse{
		Token: token,
		User: um.UserResponse{
			Id:                      user.ID,
			FullName:                user.Name,
			Email:                   user.Email,
			RoleId:                  user.UserRoleID,
			RoleName:                user.UserRole.Name,
			IsRegistrationCompleted: user.Participant.IsRegistered,
			PaymentStatus:           user.Participant.PaymentStatus,
		},
	}
	return
}

func (service *AuthServiceImpl) GoogleOauth(request oauth.GoogleUserResult) (response model.AuthResponse, err error) {
	user, err := service.UserRepo.FindByEmail(request.Email)
	if err != nil && err != e.ErrEmailNotRegistered {
		return
	}

	if user.ID != 0 && user.AuthType != constants.AuthTypeGoogle {
		err = e.ErrWrongAuthMethod
		return
	}

	latestEvent, err := service.EventRepository.FindLatest()
	if err != nil {
		return
	}

	if user.ID == 0 {
		if latestEvent.Status != constants.EventRunning {
			err = e.ErrEventNotRunning
			return
		}

		//Find user role participant
		role, err2 := service.UserRoleRepo.FindByName(constants.UserParticipant)
		if err2 != nil {
			err = err2
			return
		}

		var avatar *string
		if request.Picture != "" {
			avatar = &request.Picture
		}

		registerModel := model.RegisterModel{
			User: um.User{
				BaseEntity: common.BaseEntity{
					CreatedBy: "self",
					UpdatedBy: "self",
				},
				UserRoleID: role.ID,
				Name:       request.Name,
				Email:      request.Email,
				Avatar:     avatar,
				AuthType:   constants.AuthTypeGoogle,
				IsActive:   true,
			},
			LatestEventID: latestEvent.ID,
			Invoice: pym.Invoice{
				EventID:       latestEvent.ID,
				InvoiceNumber: utils.GenerateInvoiceNumber(),
				Status:        constants.InvoiceUnpaid,
				Amount:        latestEvent.RegFee,
				BaseEntity: common.BaseEntity{
					CreatedBy: "self",
					UpdatedBy: "self",
				},
			},
		}

		_, err = service.AuthRepository.Register(registerModel)
		if err != nil {
			return
		}

		user, err = service.UserRepo.FindByEmail(request.Email)
		if err != nil {
			return
		}
	} else {
		_, err = service.EventParticipantRepo.FindOneByEventIDAndParticipantID(latestEvent.ID, user.Participant.ID)
		if err != nil && err != e.ErrDataNotFound {
			return
		}

		if err != nil && err == e.ErrDataNotFound {
			evp := evm.EventParticipant{
				BaseEntity: common.BaseEntity{
					CreatedAt: time.Now(),
					CreatedBy: "self",
					UpdatedAt: time.Now(),
					UpdatedBy: "self",
				},
				EventID:       latestEvent.ID,
				ParticipantID: user.Participant.ID,
			}
			inv := pym.Invoice{
				BaseEntity: common.BaseEntity{
					CreatedAt: time.Now(),
					CreatedBy: "self",
					UpdatedAt: time.Now(),
					UpdatedBy: "self",
				},
				EventID:       latestEvent.ID,
				ParticipantID: user.Participant.ID,
				InvoiceNumber: utils.GenerateInvoiceNumber(),
				Amount:        latestEvent.RegFee,
				Status:        constants.InvoiceUnpaid,
			}
			if err = service.ParticipantRepo.UpdateRegistrationStatus(
				user.Participant.ID, false, &evp, &inv,
			); err != nil {
				return
			}
			user.Participant.IsRegistered = false
			user.Participant.PaymentStatus = constants.InvoiceUnpaid
		}
	}

	// Generate Token
	token, err := utils.GenerateToken(user)
	if err != nil {
		return
	}

	response = model.AuthResponse{
		Token: token,
		User: um.UserResponse{
			Id:                      user.ID,
			FullName:                user.Name,
			Email:                   user.Email,
			RoleId:                  user.UserRoleID,
			RoleName:                user.UserRole.Name,
			IsRegistrationCompleted: user.Participant.IsRegistered,
			PaymentStatus:           user.Participant.PaymentStatus,
		},
	}
	return
}

func (service *AuthServiceImpl) RegisterByGoogle(request model.RegisterByGoogleRequest) (model.AuthResponse, error) {
	// Validate Token
	payload, err := idtoken.Validate(context.Background(), request.IdToken, os.Getenv("GOOGLE_CLIENT_ID"))

	// Check if there is error
	if err != nil {
		return model.AuthResponse{}, err
	}

	// Catch User Email
	userEmail := fmt.Sprintf("%v", payload.Claims["email"])

	// Register New User
	user := mapper.RegisterByGoogleRequestToUser(request, userEmail)
	err = service.UserRepo.Save(user)
	if err != nil {
		if errors.Is(err, e.ErrDuplicateKey) {
			return model.AuthResponse{}, e.ErrEmailAlreadyExists
		}
		return model.AuthResponse{}, err
	}

	// Find New User
	newUser, err := service.UserRepo.FindByEmail(userEmail)
	if err != nil {
		return model.AuthResponse{}, err
	}

	//TODO:
	// - use db transaction. rollback when create participant failed
	participant := um.Participant{
		UserID: newUser.ID,
	}
	errCreateParticipant := service.ParticipantRepo.Save(participant)
	if errCreateParticipant != nil {
		return model.AuthResponse{}, err
	}

	// Generate Token
	token, err := utils.GenerateToken(user)
	if err != nil {
		return model.AuthResponse{}, err
	}

	return model.AuthResponse{
		Token: token,
		User: um.UserResponse{
			Id:       user.ID,
			FullName: user.Name,
			Email:    user.Email,
		},
	}, nil
}

func (service *AuthServiceImpl) LoginByGoogle(request model.LoginByGoogleRequest) (model.AuthResponse, error) {
	// Validate Token
	payload, err := idtoken.Validate(context.Background(), request.IdToken, os.Getenv("GOOGLE_CLIENT_ID"))
	if err != nil {
		return model.AuthResponse{}, err
	}

	// Catch User Email
	userEmail := fmt.Sprintf("%v", payload.Claims["email"])

	// Check User on Database
	user, err := service.UserRepo.FindByEmail(userEmail)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return model.AuthResponse{}, errors.New("user not found")
		}
		return model.AuthResponse{}, err
	}

	// check user's active status
	if !user.IsActive {
		return model.AuthResponse{}, e.ErrUserIsNotActivated
	}

	// Generate Token
	token, err := utils.GenerateToken(user)
	if err != nil {
		return model.AuthResponse{}, err
	}

	// Return nil Error
	return model.AuthResponse{
		Token: token,
		User: um.UserResponse{
			Id:       user.ID,
			FullName: user.Name,
			Email:    user.Email,
		},
	}, nil
}

func (service *AuthServiceImpl) VerifyEmail(request model.VerifyEmailRequest) error {
	verifCode, err := service.VerificationCodeRepo.FindByCode(request.Code)
	if err != nil {
		if errors.Is(err, e.ErrDataNotFound) {
			return e.ErrEmailAlreadyVerified
		}
		return err
	}

	if err = service.VerificationCodeRepo.Delete(verifCode.ID, verifCode.UserID); err != nil {
		return err
	}

	return nil
}

func (service *AuthServiceImpl) SendVerificationCode(request model.SendVerificationCodeRequest) error {
	user, err := service.UserRepo.FindByEmail(request.Email)
	if err != nil {
		if err == e.ErrEmailNotRegistered {
			return e.ErrWrongLoginCredential
		}
		return err
	}

	// auth type regular only
	if request.Type == constants.VerifCodeResetPassword && user.AuthType != constants.AuthTypeRegular {
		return e.ErrCantResetPassword
	}

	//check if verification code exists
	existingVerifCode, err := service.VerificationCodeRepo.FindByUserIdAndCodeType(user.ID, request.Type)
	if err != nil && !errors.Is(err, e.ErrDataNotFound) {
		return err
	}

	code := existingVerifCode.Code
	if errors.Is(err, e.ErrDataNotFound) {
		code = utils.GenerateUuid()
		verifCode := model.VerificationCode{
			UserID: user.ID,
			Code:   code,
			Type:   request.Type,
			BaseEntity: common.BaseEntity{
				CreatedBy: user.Email,
				UpdatedBy: user.Email,
			},
		}
		err = service.VerificationCodeRepo.Save(verifCode)
		if err != nil {
			return err
		}
	}

	//TODO:
	// - send email verification [OK]
	// - use goroutine/task queue

	templateData := email.TemplateData{
		Name:        user.Name,
		SenderEmail: "sagarahackathon@gmail.com",
		Type:        request.Type,
	}

	var emailSubject string
	if request.Type == constants.VerifCodeEmail {
		emailSubject = constants.EmailSubjectVerifyEmail
		url := fmt.Sprintf("%s%s", os.Getenv("BASE_FE_URL"), os.Getenv("VERIFY_EMAIL_REDIRECT_URL"))
		templateData.Link = fmt.Sprintf("%s%s", url, code)
	} else if request.Type == constants.VerifCodeResetPassword {
		emailSubject = constants.EmailSubjectResetPassword

		var url string
		if user.UserRole.Name == constants.UserParticipant {
			url = fmt.Sprintf("%s%s", os.Getenv("BASE_FE_URL"), os.Getenv("FORGET_PASSWORD_REDIRECT_URL"))
		} else {
			url = fmt.Sprintf("%s%s", os.Getenv("BASE_FE_ADMIN_URL"), os.Getenv("FORGET_PASSWORD_REDIRECT_URL"))
		}
		templateData.Link = fmt.Sprintf("%s%s", url, code)
	}

	templateData.Title = emailSubject
	r := email.NewRequest([]string{user.Email}, emailSubject, "")
	if err = r.ParseTemplate("./src/utils/email/template_email.html", templateData); err == nil {
		ok, err2 := r.SendEmail()
		if !ok || err2 != nil {
			fmt.Println(err2)
			return e.ErrFailedSendEmail
		}
	} else {
		return e.ErrFailedParseEmailTemplate
	}

	return nil
}

func (service *AuthServiceImpl) ValidateVerificationCode(request model.ValidateVerificationCodeRequest) error {
	verifCode, err := service.VerificationCodeRepo.FindByCode(request.Code)
	if err != nil {
		if errors.Is(err, e.ErrDataNotFound) {
			return e.ErrInvalidVerificationCode
		}
		return err
	}

	//TODO: check code expiration [TEST]
	if verifCode.ExpireAt != nil && verifCode.ExpireAt.After(time.Now()) {
		return e.ErrInvalidVerificationCode
	}

	return nil
}

func (service *AuthServiceImpl) ForgotPassword(request model.ForgotPasswordRequest) error {
	//check verification code
	verifCode, err := service.VerificationCodeRepo.FindByCode(request.Code)
	if err != nil {
		if errors.Is(err, e.ErrDataNotFound) {
			return e.ErrInvalidVerificationCode
		}
		return err
	}

	if verifCode.Type != constants.VerifCodeResetPassword {
		return e.ErrInvalidVerificationCode
	}

	//hash password
	hashed, err := utils.HashPassword(request.NewPassword)
	if err != nil {
		return err
	}

	//update password
	err = service.UserRepo.UpdatePassword(verifCode.UserID, hashed, request.Code)
	if err != nil {
		return err
	}

	return nil
}
