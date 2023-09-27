package mapper

import (
	"be-sagara-hackathon/src/modules/auth/model"
	um "be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/utils/common"
	"be-sagara-hackathon/src/utils/constants"
)

func RegisterByGoogleRequestToUser(request model.RegisterByGoogleRequest, userEmail string) um.User {
	return um.User{
		Name:        request.FullName,
		Email:       userEmail,
		PhoneNumber: &request.PhoneNumber,
		UserRoleID:  4,
		AuthType:    constants.AuthTypeGoogle,
		IsActive:    true,
		BaseEntity: common.BaseEntity{
			CreatedBy: "self",
			UpdatedBy: "self",
		},
	}
}
