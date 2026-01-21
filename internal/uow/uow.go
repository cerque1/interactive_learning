package uow

import "interactive_learning/internal/repo"

type UnitOfWork interface {
	Begin() error
	Commit() error
	Rollback() error
	Close() error

	GetUsersRepoWriter() repo.UsersRepoWrite
	GetCardRepoWriter() repo.CardRepoWrite
	GetModuleRepoWriter() repo.ModuleRepoWrite
	GetCategoryRepoWriter() repo.CategoryRepoWrite
	GetCategoryModulesRepoWriter() repo.CategoryModulesRepoWrite
	GetResultsRepoWriter() repo.ResultsRepoWrite
	GetCardsResultsRepoWriter() repo.CardsResultsRepoWrite
	GetModulesResultsRepoWriter() repo.ModulesResultsRepoWrite
	GetCategoryModulesResultsRepoWriter() repo.CategoryModulesResultsRepoWrite
	GetSelectedRepoWriter() repo.SelectedRepoWrite

	GetUsersRepoReader() repo.UsersRepoRead
	GetCardRepoReader() repo.CardRepoRead
	GetModuleRepoReader() repo.ModuleRepoRead
	GetCategoryRepoReader() repo.CategoryRepoRead
	GetCategoryModulesRepoReader() repo.CategoryModulesRepoRead
	GetResultsRepoReader() repo.ResultsRepoRead
	GetCardsResultsRepoReader() repo.CardsResultsRepoRead
	GetModulesResultsRepoReader() repo.ModulesResultsRepoRead
	GetCategoryModulesResultsRepoReader() repo.CategoryModulesResultsRepoRead
	GetSelectedRepoReader() repo.SelectedRepoRead
}
