package interactivelearning

import (
	errors_mapper "interactive_learning/internal/mappers/errors"
	"interactive_learning/internal/repo"
	"interactive_learning/internal/uow"
	"sync"
)

type UseCase struct {
	unitOfWorkFactory func() uow.UnitOfWork

	tokenStorage                   repo.TokenStorage
	usersRepoRead                  repo.UsersRepoRead
	cardsRepoRead                  repo.CardRepoRead
	moduleRepoRead                 repo.ModuleRepoRead
	categoryRepoRead               repo.CategoryRepoRead
	categoryModulesRepoRead        repo.CategoryModulesRepoRead
	resultsRepoRead                repo.ResultsRepoRead
	cardsResultsRepoRead           repo.CardsResultsRepoRead
	modulesResultsRepoRead         repo.ModulesResultsRepoRead
	categoryModulesResultsRepoRead repo.CategoryModulesResultsRepoRead
	selectedRepoRead               repo.SelectedRepoRead

	usersMutex                  sync.Mutex
	cardMutex                   sync.Mutex
	moduleMutex                 sync.Mutex
	categoryMutex               sync.Mutex
	categoryModulesMutex        sync.Mutex
	resultsMutex                sync.Mutex
	cardsResultsMutex           sync.Mutex
	modulesResultsMutex         sync.Mutex
	categoryModulesResultsMutex sync.Mutex
	selectedMutex               sync.Mutex

	errorsMapper *errors_mapper.DomainsErrorsMapper
}

func New(unitOfWorkFactory func() uow.UnitOfWork,
	tokenStorage repo.TokenStorage,
	usersRepoRead repo.UsersRepoRead,
	cardsRepoRead repo.CardRepoRead,
	moduleRepoRead repo.ModuleRepoRead,
	categoryRepoRead repo.CategoryRepoRead,
	categoryModulesRepoRead repo.CategoryModulesRepoRead,
	resultsRepoRead repo.ResultsRepoRead,
	cardsResultsRepoRead repo.CardsResultsRepoRead,
	modulesResultsRepoRead repo.ModulesResultsRepoRead,
	categoryModulesResultsRepoRead repo.CategoryModulesResultsRepoRead,
	selectedRepoRead repo.SelectedRepoRead,
	errorsMapper *errors_mapper.DomainsErrorsMapper) *UseCase {

	return &UseCase{unitOfWorkFactory: unitOfWorkFactory,
		tokenStorage:                   tokenStorage,
		usersRepoRead:                  usersRepoRead,
		cardsRepoRead:                  cardsRepoRead,
		moduleRepoRead:                 moduleRepoRead,
		categoryRepoRead:               categoryRepoRead,
		categoryModulesRepoRead:        categoryModulesRepoRead,
		resultsRepoRead:                resultsRepoRead,
		cardsResultsRepoRead:           cardsResultsRepoRead,
		modulesResultsRepoRead:         modulesResultsRepoRead,
		categoryModulesResultsRepoRead: categoryModulesResultsRepoRead,
		selectedRepoRead:               selectedRepoRead,
		errorsMapper:                   errorsMapper,
	}
}
