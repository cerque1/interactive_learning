package persistent

import (
	"database/sql"
	"interactive_learning/internal/repo"
	"interactive_learning/internal/repo/persistent"
)

type UnitOfWorkImpl struct {
	db *sql.DB
	tx *sql.Tx

	userRepoWrite                   repo.UsersRepoWrite
	cardRepoWrite                   repo.CardRepoWrite
	moduleRepoWrite                 repo.ModuleRepoWrite
	categoryRepoWrite               repo.CategoryRepoWrite
	categoryModulesRepoWrite        repo.CategoryModulesRepoWrite
	resultRepoWrite                 repo.ResultsRepoWrite
	cardsResultsRepoWrite           repo.CardsResultsRepoWrite
	modulesResultsRepoWrite         repo.ModulesResultsRepoWrite
	categoryModulesResultsRepoWrite repo.CategoryModulesResultsRepoWrite

	userRepoRead                   repo.UsersRepoRead
	cardRepoRead                   repo.CardRepoRead
	moduleRepoRead                 repo.ModuleRepoRead
	categoryRepoRead               repo.CategoryRepoRead
	categoryModulesRepoRead        repo.CategoryModulesRepoRead
	resultRepoRead                 repo.ResultsRepoRead
	cardsResultsRepoRead           repo.CardsResultsRepoRead
	modulesResultsRepoRead         repo.ModulesResultsRepoRead
	categoryModulesResultsRepoRead repo.CategoryModulesResultsRepoRead
}

func NewUnitOfWork(db *sql.DB) *UnitOfWorkImpl {
	return &UnitOfWorkImpl{
		db: db,
	}
}

func (uow *UnitOfWorkImpl) Begin() error {
	tx, err := uow.db.Begin()
	if err != nil {
		return err
	}

	uow.tx = tx
	userRepo := persistent.NewUsersRepo(tx)
	cardRepo := persistent.NewCardsRepo(tx)
	moduleRepo := persistent.NewModulesRepo(tx)
	categoryRepo := persistent.NewCategoryRepo(tx)
	categoryModulesRepo := persistent.NewCategoryModulesRepo(tx)
	resultRepo := persistent.NewResultsRepo(tx)
	cardsResultsRepo := persistent.NewCardsResultsRepo(tx)
	modulesResultsRepo := persistent.NewModulesResultsRepo(tx)
	categoryModulesResultsRepo := persistent.NewCategoryModulesResultsRepo(tx)

	uow.userRepoRead = userRepo
	uow.userRepoWrite = userRepo
	uow.cardRepoRead = cardRepo
	uow.cardRepoWrite = cardRepo
	uow.moduleRepoRead = moduleRepo
	uow.moduleRepoWrite = moduleRepo
	uow.categoryRepoRead = categoryRepo
	uow.categoryRepoWrite = categoryRepo
	uow.categoryModulesRepoRead = categoryModulesRepo
	uow.categoryModulesRepoWrite = categoryModulesRepo
	uow.resultRepoRead = resultRepo
	uow.resultRepoWrite = resultRepo
	uow.cardsResultsRepoRead = cardsResultsRepo
	uow.cardsResultsRepoWrite = cardsResultsRepo
	uow.modulesResultsRepoRead = modulesResultsRepo
	uow.modulesResultsRepoWrite = modulesResultsRepo
	uow.categoryModulesResultsRepoRead = categoryModulesResultsRepo
	uow.categoryModulesResultsRepoWrite = categoryModulesResultsRepo

	return nil
}

func (uow *UnitOfWorkImpl) Commit() error {
	if uow.tx == nil {
		return nil
	}
	return uow.tx.Commit()
}

func (uow *UnitOfWorkImpl) Rollback() error {
	if uow.tx == nil {
		return nil
	}
	return uow.tx.Rollback()
}

func (uow *UnitOfWorkImpl) Close() error {
	if uow.tx != nil {
		uow.tx.Rollback()
	}
	return nil
}

func (uow *UnitOfWorkImpl) GetUsersRepoWriter() repo.UsersRepoWrite {
	return uow.userRepoWrite
}

func (uow *UnitOfWorkImpl) GetCardRepoWriter() repo.CardRepoWrite {
	return uow.cardRepoWrite
}

func (uow *UnitOfWorkImpl) GetModuleRepoWriter() repo.ModuleRepoWrite {
	return uow.moduleRepoWrite
}

func (uow *UnitOfWorkImpl) GetCategoryRepoWriter() repo.CategoryRepoWrite {
	return uow.categoryRepoWrite
}

func (uow *UnitOfWorkImpl) GetCategoryModulesRepoWriter() repo.CategoryModulesRepoWrite {
	return uow.categoryModulesRepoWrite
}

func (uow *UnitOfWorkImpl) GetResultsRepoWriter() repo.ResultsRepoWrite {
	return uow.resultRepoWrite
}

func (uow *UnitOfWorkImpl) GetCardsResultsRepoWriter() repo.CardsResultsRepoWrite {
	return uow.cardsResultsRepoWrite
}

func (uow *UnitOfWorkImpl) GetModulesResultsRepoWriter() repo.ModulesResultsRepoWrite {
	return uow.modulesResultsRepoWrite
}

func (uow *UnitOfWorkImpl) GetCategoryModulesResultsRepoWriter() repo.CategoryModulesResultsRepoWrite {
	return uow.categoryModulesResultsRepoWrite
}

func (uow *UnitOfWorkImpl) GetUsersRepoReader() repo.UsersRepoRead {
	return uow.userRepoRead
}

func (uow *UnitOfWorkImpl) GetCardRepoReader() repo.CardRepoRead {
	return uow.cardRepoRead
}

func (uow *UnitOfWorkImpl) GetModuleRepoReader() repo.ModuleRepoRead {
	return uow.moduleRepoRead
}

func (uow *UnitOfWorkImpl) GetCategoryRepoReader() repo.CategoryRepoRead {
	return uow.categoryRepoRead
}

func (uow *UnitOfWorkImpl) GetCategoryModulesRepoReader() repo.CategoryModulesRepoRead {
	return uow.categoryModulesRepoRead
}

func (uow *UnitOfWorkImpl) GetResultsRepoReader() repo.ResultsRepoRead {
	return uow.resultRepoRead
}

func (uow *UnitOfWorkImpl) GetCardsResultsRepoReader() repo.CardsResultsRepoRead {
	return uow.cardsResultsRepoRead
}

func (uow *UnitOfWorkImpl) GetModulesResultsRepoReader() repo.ModulesResultsRepoRead {
	return uow.modulesResultsRepoRead
}

func (uow *UnitOfWorkImpl) GetCategoryModulesResultsRepoReader() repo.CategoryModulesResultsRepoRead {
	return uow.categoryModulesResultsRepoRead
}
