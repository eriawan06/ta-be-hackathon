package service

import (
	"be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/modules/user/repository"
	"be-sagara-hackathon/src/utils"
	"be-sagara-hackathon/src/utils/common/builder"
	"be-sagara-hackathon/src/utils/constants"
	e "be-sagara-hackathon/src/utils/errors"
	"be-sagara-hackathon/src/utils/helper"
	"context"
)

type ParticipantService interface {
	GetProfile(ctx context.Context) (participant model.Participant, err error)
	UpdateFull(ctx context.Context, request model.UpdateParticipantFull) (model.Participant, error)
	UpdateProfileAndLocation(ctx context.Context, request model.UpdateParticipantProfileRequest) (model.Participant, error)
	UpdateEducation(ctx context.Context, request model.UpdateParticipantEducationRequest) (model.Participant, error)
	UpdatePreference(ctx context.Context, request model.UpdateParticipantPreferenceRequest) (model.Participant, error)
	UpdateAccount(ctx context.Context, request model.UpdateParticipantAccountRequest) (model.Participant, error)
	CompleteRegistration(ctx context.Context) error
	GetList(
		filter model.FilterUser,
		pg *utils.PaginateQueryOffset,
	) (response model.ListParticipantResponse, err error)
	GetDetail(id uint) (participant model.Participant, err error)
	GetListParticipantSearch(
		filter model.FilterParticipantSearch,
		pg *utils.PaginateQueryOffset,
	) (response model.ListParticipantSearchResponse, err error)
}

type ParticipantServiceImpl struct {
	Repository     repository.ParticipantRepository
	UserRepo       repository.UserRepository
	RoleRepository repository.UserRoleRepository
}

func NewParticipantService(
	repository repository.ParticipantRepository,
	userRepo repository.UserRepository,
	roleRepo repository.UserRoleRepository,
) ParticipantService {
	return &ParticipantServiceImpl{
		Repository:     repository,
		UserRepo:       userRepo,
		RoleRepository: roleRepo,
	}
}

func (service *ParticipantServiceImpl) GetProfile(ctx context.Context) (participant model.Participant, err error) {
	authenticatedUser := ctx.Value("user").(model.User)
	if participant, err = service.Repository.FindByEmail(authenticatedUser.Email); err != nil {
		return
	}
	return
}

func (service *ParticipantServiceImpl) UpdateFull(ctx context.Context, request model.UpdateParticipantFull) (model.Participant, error) {
	authenticatedUser := ctx.Value("user").(model.User)
	participant, err := service.Repository.FindByEmail(authenticatedUser.Email)
	if err != nil {
		return participant, err
	}

	participant.BaseEntity = builder.BuildBaseEntity(ctx, false, &participant.BaseEntity)
	participant.User.BaseEntity = builder.BuildBaseEntity(ctx, false, &participant.User.BaseEntity)
	participant.User.Avatar = request.Avatar
	participant.Bio = request.Bio
	participant.Gender = &request.Gender
	participant.Address = &request.Address
	participant.ProvinceID = &request.ProvinceID
	participant.CityID = &request.CityID
	participant.DistrictID = &request.DistrictID
	participant.VillageID = &request.VillageID

	birthDate, err := helper.ParseDateStringToTime(request.Birthdate)
	if err != nil {
		return participant, err
	}
	participant.Birthdate = &birthDate

	if request.Action == "update" {
		participant.User.Name = request.Name
	}
	if request.Action == "register" {
		participant.User.PhoneNumber = &request.PhoneNumber
	}

	participant.LevelOfStudy = &request.LevelOfStudy
	participant.School = &request.School
	participant.GraduationYear = request.GraduationYear
	participant.Major = request.Major

	participant.User.OccupationID = &request.OccupationID
	participant.User.Institution = request.CompanyName
	participant.NumOfHackathon = request.NumOfHackathon
	participant.LinkPortfolio = request.Portfolio
	participant.LinkRepository = request.Repository
	participant.LinkLinkedin = request.Linkedin
	participant.Resume = request.Resume
	participant.SpecialityID = &request.SpecialityID

	var skills []model.ParticipantSkill
	if len(request.Skills) > 0 {
		for _, v := range request.Skills {
			skills = append(skills, model.ParticipantSkill{
				ParticipantID: participant.ID,
				SkillID:       v,
			})
		}
	}

	if err = service.Repository.Update(model.UpdateParticipant{
		ID:            participant.ID,
		Participant:   participant,
		Skills:        skills,
		RemovedSkills: request.RemovedSkills,
	}); err != nil {
		return participant, err
	}

	if request.Action == "register" {
		if err = service.CompleteRegistration(ctx); err != nil {
			return participant, err
		}
		participant.IsRegistered = true
	}

	return participant, nil
}

func (service *ParticipantServiceImpl) UpdateProfileAndLocation(ctx context.Context, request model.UpdateParticipantProfileRequest) (model.Participant, error) {
	authenticatedUser := ctx.Value("user").(model.User)
	participant, err := service.Repository.FindByEmail(authenticatedUser.Email)
	if err != nil {
		return participant, err
	}

	participant.BaseEntity = builder.BuildBaseEntity(ctx, false, &participant.BaseEntity)
	participant.User.BaseEntity = builder.BuildBaseEntity(ctx, false, &participant.User.BaseEntity)
	participant.User.Avatar = request.Avatar
	participant.Bio = request.Bio
	participant.Gender = &request.Gender
	participant.Address = &request.Address
	participant.ProvinceID = &request.ProvinceID
	participant.CityID = &request.CityID
	participant.DistrictID = &request.DistrictID
	participant.VillageID = &request.VillageID

	birthDate, err := helper.ParseDateStringToTime(request.Birthdate)
	if err != nil {
		return participant, err
	}
	participant.Birthdate = &birthDate

	if request.Action == "update" {
		participant.User.Name = request.Name
	}
	if request.Action == "register" {
		participant.User.PhoneNumber = &request.PhoneNumber
	}

	if err = service.Repository.Update(model.UpdateParticipant{
		ID:          participant.ID,
		Participant: participant,
	}); err != nil {
		return participant, err
	}
	return participant, nil
}

func (service *ParticipantServiceImpl) UpdateEducation(ctx context.Context, request model.UpdateParticipantEducationRequest) (model.Participant, error) {
	authenticatedUser := ctx.Value("user").(model.User)
	participant, err := service.Repository.FindByEmail(authenticatedUser.Email)
	if err != nil {
		return participant, err
	}

	participant.BaseEntity = builder.BuildBaseEntity(ctx, false, &participant.BaseEntity)
	participant.User.BaseEntity = builder.BuildBaseEntity(ctx, false, &participant.User.BaseEntity)
	participant.LevelOfStudy = &request.LevelOfStudy
	participant.School = &request.School
	participant.GraduationYear = request.GraduationYear
	participant.Major = request.Major
	if err = service.Repository.Update(model.UpdateParticipant{
		ID:          participant.ID,
		Participant: participant,
	}); err != nil {
		return participant, err
	}
	return participant, nil
}

func (service *ParticipantServiceImpl) UpdatePreference(ctx context.Context, request model.UpdateParticipantPreferenceRequest) (model.Participant, error) {
	authenticatedUser := ctx.Value("user").(model.User)
	participant, err := service.Repository.FindByEmail(authenticatedUser.Email)
	if err != nil {
		return participant, err
	}

	participant.BaseEntity = builder.BuildBaseEntity(ctx, false, &participant.BaseEntity)
	participant.User.BaseEntity = builder.BuildBaseEntity(ctx, false, &participant.User.BaseEntity)
	participant.User.OccupationID = &request.OccupationID
	participant.User.Institution = request.CompanyName
	participant.NumOfHackathon = request.NumOfHackathon
	participant.LinkPortfolio = request.Portfolio
	participant.LinkRepository = request.Repository
	participant.LinkLinkedin = request.Linkedin
	participant.Resume = request.Resume
	participant.SpecialityID = &request.SpecialityID

	var skills []model.ParticipantSkill
	if len(request.Skills) > 0 {
		for _, v := range request.Skills {
			skills = append(skills, model.ParticipantSkill{
				ParticipantID: participant.ID,
				SkillID:       v,
			})
		}
	}

	if err = service.Repository.Update(model.UpdateParticipant{
		ID:            participant.ID,
		Participant:   participant,
		Skills:        skills,
		RemovedSkills: request.RemovedSkills,
	}); err != nil {
		return participant, err
	}
	return participant, nil
}

func (service *ParticipantServiceImpl) UpdateAccount(ctx context.Context, request model.UpdateParticipantAccountRequest) (model.Participant, error) {
	authenticatedUser := ctx.Value("user").(model.User)
	participant, err := service.Repository.FindByEmail(authenticatedUser.Email)
	if err != nil {
		return participant, err
	}

	participant.User.PhoneNumber = &request.PhoneNumber
	participant.User.Username = &request.Username
	if err = service.UserRepo.Update(participant.UserID, *participant.User); err != nil {
		return participant, err
	}
	return participant, nil
}

func (service *ParticipantServiceImpl) CompleteRegistration(ctx context.Context) error {
	authenticatedUser := ctx.Value("user").(model.User)
	participant, err := service.Repository.FindByEmail(authenticatedUser.Email)
	if err != nil {
		return err
	}

	isComplete := isRegistrationCompleted(participant)
	if !isComplete {
		return e.ErrRegistrationNotCompleted
	}

	participant.BaseEntity = builder.BuildBaseEntity(ctx, false, &participant.BaseEntity)
	participant.User.BaseEntity = builder.BuildBaseEntity(ctx, false, &participant.User.BaseEntity)
	participant.IsRegistered = isComplete
	if err = service.Repository.Update(model.UpdateParticipant{
		ID:          participant.ID,
		Participant: participant,
	}); err != nil {
		return err
	}
	return nil
}

func isRegistrationCompleted(participant model.Participant) bool {
	if participant.User.PhoneNumber == nil || participant.Birthdate == nil || participant.Gender == nil || participant.Address == nil ||
		participant.ProvinceID == nil || (participant.ProvinceID != nil && *participant.ProvinceID == 0) ||
		participant.CityID == nil || (participant.CityID != nil && *participant.CityID == 0) ||
		participant.DistrictID == nil || (participant.DistrictID != nil && *participant.DistrictID == 0) ||
		participant.VillageID == nil || (participant.VillageID != nil && *participant.VillageID == 0) ||
		participant.LevelOfStudy == nil || participant.School == nil || participant.GraduationYear == 0 ||
		participant.User.OccupationID == nil || participant.SpecialityID == nil || len(participant.Skills) == 0 {
		return false
	}

	return true
}

func (service *ParticipantServiceImpl) GetList(
	filter model.FilterUser,
	pg *utils.PaginateQueryOffset,
) (response model.ListParticipantResponse, err error) {
	//Get Role Participant
	role, err := service.RoleRepository.FindByName(constants.UserParticipant)
	if err != nil {
		return
	}

	filter.RoleID = int(role.ID)
	participants, totalData, totalPage, err := service.UserRepo.Find(filter, pg)
	if err != nil {
		return
	}

	var responseData []model.ParticipantLite
	for _, v := range participants {
		var participant model.ParticipantLite
		participant.ID = v.Participant.ID
		participant.UserID = v.ID
		participant.Name = v.Name
		participant.Email = v.Email
		participant.PhoneNumber = v.PhoneNumber
		participant.IsActive = v.IsActive
		responseData = append(responseData, participant)
	}
	response.Participants = responseData
	response.TotalPage = totalPage
	response.TotalItem = totalData
	return
}

func (service *ParticipantServiceImpl) GetDetail(id uint) (participant model.Participant, err error) {
	if participant, err = service.Repository.FindDetail(id); err != nil {
		return
	}
	return
}

func (service *ParticipantServiceImpl) GetListParticipantSearch(
	filter model.FilterParticipantSearch,
	pg *utils.PaginateQueryOffset,
) (response model.ListParticipantSearchResponse, err error) {
	filter.InTeam = false
	response.Participants, response.TotalItem, response.TotalPage, err = service.Repository.Find(filter, pg)
	if err != nil {
		return
	}
	return
}
