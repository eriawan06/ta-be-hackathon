package builder

import (
	um "be-sagara-hackathon/src/modules/user/model"
	"be-sagara-hackathon/src/utils/common"
	"context"
	"time"
)

func BuildBaseEntity(ctx context.Context, isCreate bool, be *common.BaseEntity) common.BaseEntity {
	var entity common.BaseEntity

	if isCreate {
		entity = common.BaseEntity{
			CreatedAt: time.Now(),
			CreatedBy: ctx.Value("user").(um.User).Email,
			UpdatedAt: time.Now(),
			UpdatedBy: ctx.Value("user").(um.User).Email,
		}
	} else {
		entity = *be
		entity.UpdatedAt = time.Now()
		entity.UpdatedBy = ctx.Value("user").(um.User).Email
	}

	return entity
}
