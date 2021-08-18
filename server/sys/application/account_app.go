package application

import (
	"mayfly-go/base/biz"
	"mayfly-go/base/model"
	"mayfly-go/base/utils"
	"mayfly-go/server/sys/domain/entity"
	"mayfly-go/server/sys/domain/repository"
	"mayfly-go/server/sys/infrastructure/persistence"

	"gorm.io/gorm"
)

type Account interface {
	GetAccount(condition *entity.Account, cols ...string) error

	GetPageList(condition *entity.Account, pageParam *model.PageParam, toEntity interface{}, orderBy ...string) *model.PageResult

	Create(account *entity.Account)

	Update(account *entity.Account)

	Delete(id uint64)
}

type accountAppImpl struct {
	accountRepo repository.Account
}

var AccountApp Account = &accountAppImpl{
	accountRepo: persistence.AccountDao,
}

// 根据条件获取账号信息
func (a *accountAppImpl) GetAccount(condition *entity.Account, cols ...string) error {
	return a.accountRepo.GetAccount(condition, cols...)
}

func (a *accountAppImpl) GetPageList(condition *entity.Account, pageParam *model.PageParam, toEntity interface{}, orderBy ...string) *model.PageResult {
	return a.accountRepo.GetPageList(condition, pageParam, toEntity)
}

func (a *accountAppImpl) Create(account *entity.Account) {
	biz.IsTrue(a.GetAccount(&entity.Account{Username: account.Username}) != nil, "该账号用户名已存在")
	// 默认密码为账号用户名
	account.Password = utils.Md5(account.Username)
	account.Status = entity.AccountEnableStatus
	a.accountRepo.Insert(account)
}

func (a *accountAppImpl) Update(account *entity.Account) {
	// 禁止更新用户名，防止误传被更新
	account.Username = ""
	a.accountRepo.Update(account)
}

func (a *accountAppImpl) Delete(id uint64) {
	err := model.Tx(
		func(db *gorm.DB) error {
			// 删除account表信息
			return db.Delete(new(entity.Account), "id = ?", id).Error
		},
		func(db *gorm.DB) error {
			// 删除账号关联的角色信息
			accountRole := &entity.AccountRole{AccountId: id}
			return db.Where(accountRole).Delete(accountRole).Error
		},
	)
	biz.ErrIsNilAppendErr(err, "删除失败：%s")
}