package gql

import (
	"math/big"
)

type WeightFunctionType string

const (
	WeightFunctionTypeLinear      WeightFunctionType = "LINEAR"
	WeightFunctionTypeExponential WeightFunctionType = "EXPONENTIAL"
)

type CollectionType string

const (
	CollectionTypeERC721  CollectionType = "ERC721"
	CollectionTypeERC1155 CollectionType = "ERC1155"
)

type EpochStatus string

const (
	EpochStatusActive     EpochStatus = "ACTIVE"
	EpochStatusProcessing EpochStatus = "PROCESSING"
	EpochStatusCompleted  EpochStatus = "COMPLETED"
	EpochStatusFailed     EpochStatus = "FAILED"
)

type RoleType string

const (
	RoleTypeAdmin             RoleType = "ADMIN"
	RoleTypeEpochManager      RoleType = "EPOCH_MANAGER"
	RoleTypeCollectionManager RoleType = "COLLECTION_MANAGER"
	RoleTypeSubsidyManager    RoleType = "SUBSIDY_MANAGER"
	RoleTypeLendingManager    RoleType = "LENDING_MANAGER"
)

type CollectionsVault struct {
	ID                       string                     `json:"id"`
	CTokenMarket             *CTokenMarket              `json:"cTokenMarket"`
	TotalShares              *big.Int                   `json:"totalShares"`
	TotalDeposits            *big.Int                   `json:"totalDeposits"`
	TotalCTokens             *big.Int                   `json:"totalCTokens"`
	GlobalDepositIndex       *big.Int                   `json:"globalDepositIndex"`
	TotalPrincipalDeposited  *big.Int                   `json:"totalPrincipalDeposited"`
	CollectionRegistry       *CollectionRegistry        `json:"collectionRegistry"`
	EpochManager             *EpochManager              `json:"epochManager"`
	LendingManager           *LendingManager            `json:"lendingManager"`
	DebtSubsidizer           *DebtSubsidizer            `json:"debtSubsidizer"`
	CollectionParticipations []*CollectionParticipation `json:"collectionParticipations"`
	MerkleDistributions      []*MerkleDistribution      `json:"merkleDistributions"`
	EpochAllocations         []*EpochVaultAllocation    `json:"epochAllocations"`
	CreatedAtBlock           *big.Int                   `json:"createdAtBlock"`
	CreatedAtTimestamp       *big.Int                   `json:"createdAtTimestamp"`
	UpdatedAtBlock           *big.Int                   `json:"updatedAtBlock"`
	UpdatedAtTimestamp       *big.Int                   `json:"updatedAtTimestamp"`
}

type CollectionRegistry struct {
	ID                     string                   `json:"id"`
	TotalCollections       *big.Int                 `json:"totalCollections"`
	TotalActiveCollections *big.Int                 `json:"totalActiveCollections"`
	Owner                  string                   `json:"owner"`
	Roles                  []*AccountRoleAssignment `json:"roles"`
	Collections            []*Collection            `json:"collections"`
	CreatedAtBlock         *big.Int                 `json:"createdAtBlock"`
	CreatedAtTimestamp     *big.Int                 `json:"createdAtTimestamp"`
	UpdatedAtBlock         *big.Int                 `json:"updatedAtBlock"`
	UpdatedAtTimestamp     *big.Int                 `json:"updatedAtTimestamp"`
}

type EpochManager struct {
	ID                    string                   `json:"id"`
	CurrentEpochID        *big.Int                 `json:"currentEpochId"`
	CurrentEpoch          *Epoch                   `json:"currentEpoch"`
	EpochDuration         *big.Int                 `json:"epochDuration"`
	ProcessingBuffer      *big.Int                 `json:"processingBuffer"`
	MinimumYieldThreshold *big.Int                 `json:"minimumYieldThreshold"`
	TotalEpochs           *big.Int                 `json:"totalEpochs"`
	TotalYieldDistributed *big.Int                 `json:"totalYieldDistributed"`
	Owner                 string                   `json:"owner"`
	Roles                 []*AccountRoleAssignment `json:"roles"`
	Epochs                []*Epoch                 `json:"epochs"`
	Vault                 *CollectionsVault        `json:"vault"`
	CreatedAtBlock        *big.Int                 `json:"createdAtBlock"`
	CreatedAtTimestamp    *big.Int                 `json:"createdAtTimestamp"`
	UpdatedAtBlock        *big.Int                 `json:"updatedAtBlock"`
	UpdatedAtTimestamp    *big.Int                 `json:"updatedAtTimestamp"`
}

type LendingManager struct {
	ID                      string                   `json:"id"`
	SupportedMarkets        *big.Int                 `json:"supportedMarkets"`
	TotalMarketParticipants *big.Int                 `json:"totalMarketParticipants"`
	GlobalCollateralFactor  *big.Int                 `json:"globalCollateralFactor"`
	LiquidationIncentive    *big.Int                 `json:"liquidationIncentive"`
	TotalSupplyVolume       *big.Int                 `json:"totalSupplyVolume"`
	TotalBorrowVolume       *big.Int                 `json:"totalBorrowVolume"`
	TotalLiquidationVolume  *big.Int                 `json:"totalLiquidationVolume"`
	Owner                   string                   `json:"owner"`
	Roles                   []*AccountRoleAssignment `json:"roles"`
	Vault                   *CollectionsVault        `json:"vault"`
	CreatedAtBlock          *big.Int                 `json:"createdAtBlock"`
	CreatedAtTimestamp      *big.Int                 `json:"createdAtTimestamp"`
	UpdatedAtBlock          *big.Int                 `json:"updatedAtBlock"`
	UpdatedAtTimestamp      *big.Int                 `json:"updatedAtTimestamp"`
}

type DebtSubsidizer struct {
	ID                        string                   `json:"id"`
	TotalSubsidyPool          *big.Int                 `json:"totalSubsidyPool"`
	TotalSubsidiesDistributed *big.Int                 `json:"totalSubsidiesDistributed"`
	TotalSubsidiesRemaining   *big.Int                 `json:"totalSubsidiesRemaining"`
	TotalEligibleUsers        *big.Int                 `json:"totalEligibleUsers"`
	SubsidyRate               *big.Int                 `json:"subsidyRate"`
	MaxSubsidyPerUser         *big.Int                 `json:"maxSubsidyPerUser"`
	SubsidyDuration           *big.Int                 `json:"subsidyDuration"`
	Owner                     string                   `json:"owner"`
	Roles                     []*AccountRoleAssignment `json:"roles"`
	SubsidyDistributions      []*SubsidyDistribution   `json:"subsidyDistributions"`
	Vaults                    []*CollectionsVault      `json:"vaults"`
	VaultAdditions            []*VaultAddition         `json:"vaultAdditions"`
	WhitelistedCollections    []*CollectionWhitelist   `json:"whitelistedCollections"`
	CreatedAtBlock            *big.Int                 `json:"createdAtBlock"`
	CreatedAtTimestamp        *big.Int                 `json:"createdAtTimestamp"`
	UpdatedAtBlock            *big.Int                 `json:"updatedAtBlock"`
	UpdatedAtTimestamp        *big.Int                 `json:"updatedAtTimestamp"`
}

type Epoch struct {
	ID                            string                  `json:"id"`
	EpochNumber                   *big.Int                `json:"epochNumber"`
	EpochManager                  *EpochManager           `json:"epochManager"`
	Status                        EpochStatus             `json:"status"`
	StartTimestamp                *big.Int                `json:"startTimestamp"`
	EndTimestamp                  *big.Int                `json:"endTimestamp"`
	ProcessingStartedTimestamp    *big.Int                `json:"processingStartedTimestamp"`
	ProcessingCompletedTimestamp  *big.Int                `json:"processingCompletedTimestamp"`
	TotalYieldAvailable           *big.Int                `json:"totalYieldAvailable"`
	TotalYieldAllocated           *big.Int                `json:"totalYieldAllocated"`
	TotalYieldDistributed         *big.Int                `json:"totalYieldDistributed"`
	RemainingYield                *big.Int                `json:"remainingYield"`
	TotalSubsidiesDistributed     *big.Int                `json:"totalSubsidiesDistributed"`
	TotalEligibleUsers            *big.Int                `json:"totalEligibleUsers"`
	TotalParticipatingCollections *big.Int                `json:"totalParticipatingCollections"`
	ParticipantCount              *big.Int                `json:"participantCount"`
	ProcessingTimeMs              *big.Int                `json:"processingTimeMs"`
	EstimatedProcessingTime       *big.Int                `json:"estimatedProcessingTime"`
	ProcessingGasUsed             *big.Int                `json:"processingGasUsed"`
	ProcessingTransactionCount    *big.Int                `json:"processingTransactionCount"`
	VaultAllocations              []*EpochVaultAllocation `json:"vaultAllocations"`
	SubsidyDistributions          []*SubsidyDistribution  `json:"subsidyDistributions"`
	MerkleDistributions           []*MerkleDistribution   `json:"merkleDistributions"`
	UserEligibilities             []*UserEpochEligibility `json:"userEligibilities"`
	CreatedAtBlock                *big.Int                `json:"createdAtBlock"`
	CreatedAtTimestamp            *big.Int                `json:"createdAtTimestamp"`
	UpdatedAtBlock                *big.Int                `json:"updatedAtBlock"`
	UpdatedAtTimestamp            *big.Int                `json:"updatedAtTimestamp"`
}

type CollectionParticipation struct {
	ID                     string                           `json:"id"`
	Collection             *Collection                      `json:"collection"`
	Vault                  *CollectionsVault                `json:"vault"`
	PrincipalShares        *big.Int                         `json:"principalShares"`
	PrincipalDeposited     *big.Int                         `json:"principalDeposited"`
	TotalCTokens           *big.Int                         `json:"totalCTokens"`
	GlobalDepositIndex     *big.Int                         `json:"globalDepositIndex"`
	LastGlobalDepositIndex *big.Int                         `json:"lastGlobalDepositIndex"`
	YieldAccrued           *big.Int                         `json:"yieldAccrued"`
	YieldClaimed           *big.Int                         `json:"yieldClaimed"`
	TotalYieldGenerated    *big.Int                         `json:"totalYieldGenerated"`
	IsBorrowBased          bool                             `json:"isBorrowBased"`
	RewardSharePercentage  *big.Int                         `json:"rewardSharePercentage"`
	WeightFunctionType     WeightFunctionType               `json:"weightFunctionType"`
	WeightFunctionP1       *big.Int                         `json:"weightFunctionP1"`
	WeightFunctionP2       *big.Int                         `json:"weightFunctionP2"`
	SecondsAccumulated     *big.Int                         `json:"secondsAccumulated"`
	SecondsClaimed         *big.Int                         `json:"secondsClaimed"`
	TotalSubsidies         *big.Int                         `json:"totalSubsidies"`
	TotalSubsidiesClaimed  *big.Int                         `json:"totalSubsidiesClaimed"`
	AverageAPY             *big.Int                         `json:"averageAPY"`
	TotalParticipants      *big.Int                         `json:"totalParticipants"`
	AccountSubsidies       []*AccountSubsidiesPerCollection `json:"accountSubsidies"`
	CreatedAtBlock         *big.Int                         `json:"createdAtBlock"`
	CreatedAtTimestamp     *big.Int                         `json:"createdAtTimestamp"`
	UpdatedAtBlock         *big.Int                         `json:"updatedAtBlock"`
	UpdatedAtTimestamp     *big.Int                         `json:"updatedAtTimestamp"`
}

type UserEpochEligibility struct {
	ID                    string      `json:"id"`
	User                  *Account    `json:"user"`
	Epoch                 *Epoch      `json:"epoch"`
	Collection            *Collection `json:"collection"`
	NftBalance            *big.Int    `json:"nftBalance"`
	BorrowBalance         *big.Int    `json:"borrowBalance"`
	HoldingDuration       *big.Int    `json:"holdingDuration"`
	IsEligible            bool        `json:"isEligible"`
	SubsidyReceived       *big.Int    `json:"subsidyReceived"`
	YieldShare            *big.Int    `json:"yieldShare"`
	BonusMultiplier       *big.Int    `json:"bonusMultiplier"`
	CalculatedAtBlock     *big.Int    `json:"calculatedAtBlock"`
	CalculatedAtTimestamp *big.Int    `json:"calculatedAtTimestamp"`
}

type SubsidyDistribution struct {
	ID                   string            `json:"id"`
	Epoch                *Epoch            `json:"epoch"`
	DebtSubsidizer       *DebtSubsidizer   `json:"debtSubsidizer"`
	User                 *Account          `json:"user"`
	Collection           *Collection       `json:"collection"`
	Vault                *CollectionsVault `json:"vault"`
	SubsidyAmount        *big.Int          `json:"subsidyAmount"`
	BorrowAmountBefore   *big.Int          `json:"borrowAmountBefore"`
	BorrowAmountAfter    *big.Int          `json:"borrowAmountAfter"`
	NftBalance           *big.Int          `json:"nftBalance"`
	WeightedContribution *big.Int          `json:"weightedContribution"`
	GasUsed              *big.Int          `json:"gasUsed"`
	BlockNumber          *big.Int          `json:"blockNumber"`
	Timestamp            *big.Int          `json:"timestamp"`
	TransactionHash      string            `json:"transactionHash"`
}

type MerkleDistribution struct {
	ID              string            `json:"id"`
	Epoch           *Epoch            `json:"epoch"`
	Vault           *CollectionsVault `json:"vault"`
	MerkleRoot      string            `json:"merkleRoot"`
	TotalAmount     *big.Int          `json:"totalAmount"`
	TotalClaims     *big.Int          `json:"totalClaims"`
	BlockNumber     *big.Int          `json:"blockNumber"`
	Timestamp       *big.Int          `json:"timestamp"`
	TransactionHash string            `json:"transactionHash"`
}

type Collection struct {
	ID                     string                     `json:"id"`
	ContractAddress        string                     `json:"contractAddress"`
	Name                   string                     `json:"name"`
	Symbol                 string                     `json:"symbol"`
	TotalSupply            *big.Int                   `json:"totalSupply"`
	CollectionType         CollectionType             `json:"collectionType"`
	Registry               *CollectionRegistry        `json:"registry"`
	IsActive               bool                       `json:"isActive"`
	YieldSharePercentage   *big.Int                   `json:"yieldSharePercentage"`
	WeightFunctionType     WeightFunctionType         `json:"weightFunctionType"`
	WeightFunctionP1       *big.Int                   `json:"weightFunctionP1"`
	WeightFunctionP2       *big.Int                   `json:"weightFunctionP2"`
	MinBorrowAmount        *big.Int                   `json:"minBorrowAmount"`
	MaxBorrowAmount        *big.Int                   `json:"maxBorrowAmount"`
	TotalNFTsDeposited     *big.Int                   `json:"totalNFTsDeposited"`
	TotalBorrowVolume      *big.Int                   `json:"totalBorrowVolume"`
	TotalYieldGenerated    *big.Int                   `json:"totalYieldGenerated"`
	TotalSubsidiesReceived *big.Int                   `json:"totalSubsidiesReceived"`
	Participations         []*CollectionParticipation `json:"participations"`
	YieldAccruals          []*CollectionYieldAccrual  `json:"yieldAccruals"`
	RegisteredAtBlock      *big.Int                   `json:"registeredAtBlock"`
	RegisteredAtTimestamp  *big.Int                   `json:"registeredAtTimestamp"`
	UpdatedAtBlock         *big.Int                   `json:"updatedAtBlock"`
	UpdatedAtTimestamp     *big.Int                   `json:"updatedAtTimestamp"`
}

type CTokenMarket struct {
	ID                        string              `json:"id"`
	Symbol                    string              `json:"symbol"`
	Name                      string              `json:"name"`
	Decimals                  int                 `json:"decimals"`
	TotalSupply               *big.Int            `json:"totalSupply"`
	TotalBorrows              *big.Int            `json:"totalBorrows"`
	TotalReserves             *big.Int            `json:"totalReserves"`
	ExchangeRate              *big.Int            `json:"exchangeRate"`
	InterestAccumulated       *big.Int            `json:"interestAccumulated"`
	CashPrior                 *big.Int            `json:"cashPrior"`
	BorrowIndex               *big.Int            `json:"borrowIndex"`
	CollateralFactor          *big.Int            `json:"collateralFactor"`
	LiquidationIncentive      *big.Int            `json:"liquidationIncentive"`
	ReserveFactor             *big.Int            `json:"reserveFactor"`
	BaseRatePerBlock          *big.Int            `json:"baseRatePerBlock"`
	MultiplierPerBlock        *big.Int            `json:"multiplierPerBlock"`
	JumpMultiplierPerBlock    *big.Int            `json:"jumpMultiplierPerBlock"`
	Kink                      *big.Int            `json:"kink"`
	Vaults                    []*CollectionsVault `json:"vaults"`
	AccountMarkets            []*AccountMarket    `json:"accountMarkets"`
	LastExchangeRateTimestamp *big.Int            `json:"lastExchangeRateTimestamp"`
	UpdatedAtBlock            *big.Int            `json:"updatedAtBlock"`
	UpdatedAtTimestamp        *big.Int            `json:"updatedAtTimestamp"`
}

type Account struct {
	ID                           string                           `json:"id"`
	TotalSecondsClaimed          *big.Int                         `json:"totalSecondsClaimed"`
	TotalSubsidiesReceived       *big.Int                         `json:"totalSubsidiesReceived"`
	TotalYieldEarned             *big.Int                         `json:"totalYieldEarned"`
	TotalBorrowVolume            *big.Int                         `json:"totalBorrowVolume"`
	TotalNFTsOwned               *big.Int                         `json:"totalNFTsOwned"`
	TotalCollectionsParticipated *big.Int                         `json:"totalCollectionsParticipated"`
	Markets                      []*AccountMarket                 `json:"markets"`
	AccountSubsidies             []*AccountSubsidiesPerCollection `json:"accountSubsidies"`
	RoleAssignments              []*AccountRoleAssignment         `json:"roleAssignments"`
	UserEligibilities            []*UserEpochEligibility          `json:"userEligibilities"`
	FirstInteractionBlock        *big.Int                         `json:"firstInteractionBlock"`
	FirstInteractionTimestamp    *big.Int                         `json:"firstInteractionTimestamp"`
	UpdatedAtBlock               *big.Int                         `json:"updatedAtBlock"`
	UpdatedAtTimestamp           *big.Int                         `json:"updatedAtTimestamp"`
}

type AccountMarket struct {
	ID                     string                           `json:"id"`
	Account                *Account                         `json:"account"`
	CTokenMarket           *CTokenMarket                    `json:"cTokenMarket"`
	SupplyBalance          *big.Int                         `json:"supplyBalance"`
	BorrowBalance          *big.Int                         `json:"borrowBalance"`
	CollateralBalance      *big.Int                         `json:"collateralBalance"`
	SupplyIndex            *big.Int                         `json:"supplyIndex"`
	BorrowIndex            *big.Int                         `json:"borrowIndex"`
	AccountSubsidies       []*AccountSubsidiesPerCollection `json:"accountSubsidies"`
	EnteredMarketBlock     *big.Int                         `json:"enteredMarketBlock"`
	EnteredMarketTimestamp *big.Int                         `json:"enteredMarketTimestamp"`
	UpdatedAtBlock         *big.Int                         `json:"updatedAtBlock"`
	UpdatedAtTimestamp     *big.Int                         `json:"updatedAtTimestamp"`
}

type AccountSubsidiesPerCollection struct {
	ID                      string                   `json:"id"`
	Account                 *Account                 `json:"account"`
	Vault                   *CollectionsVault        `json:"vault"`
	Collection              *Collection              `json:"collection"`
	AccountMarket           *AccountMarket           `json:"accountMarket"`
	CollectionParticipation *CollectionParticipation `json:"collectionParticipation"`
	BalanceNFT              *big.Int                 `json:"balanceNFT"`
	WeightedBalance         *big.Int                 `json:"weightedBalance"`
	SecondsAccumulated      *big.Int                 `json:"secondsAccumulated"`
	SecondsClaimed          *big.Int                 `json:"secondsClaimed"`
	SubsidiesAccrued        *big.Int                 `json:"subsidiesAccrued"`
	SubsidiesClaimed        *big.Int                 `json:"subsidiesClaimed"`
	AverageHoldingPeriod    *big.Int                 `json:"averageHoldingPeriod"`
	TotalRewardsEarned      *big.Int                 `json:"totalRewardsEarned"`
	UpdatedAtBlock          *big.Int                 `json:"updatedAtBlock"`
	UpdatedAtTimestamp      *big.Int                 `json:"updatedAtTimestamp"`
}

type Role struct {
	ID                 string                   `json:"id"`
	RoleType           RoleType                 `json:"roleType"`
	Name               string                   `json:"name"`
	Description        string                   `json:"description"`
	Assignments        []*AccountRoleAssignment `json:"assignments"`
	CreatedAtBlock     *big.Int                 `json:"createdAtBlock"`
	CreatedAtTimestamp *big.Int                 `json:"createdAtTimestamp"`
}

type AccountRoleAssignment struct {
	ID                 string              `json:"id"`
	Account            *Account            `json:"account"`
	Role               *Role               `json:"role"`
	Registry           *CollectionRegistry `json:"registry"`
	EpochManager       *EpochManager       `json:"epochManager"`
	LendingManager     *LendingManager     `json:"lendingManager"`
	DebtSubsidizer     *DebtSubsidizer     `json:"debtSubsidizer"`
	IsActive           bool                `json:"isActive"`
	GrantedBy          string              `json:"grantedBy"`
	GrantedAtBlock     *big.Int            `json:"grantedAtBlock"`
	GrantedAtTimestamp *big.Int            `json:"grantedAtTimestamp"`
	RevokedAtBlock     *big.Int            `json:"revokedAtBlock"`
	RevokedAtTimestamp *big.Int            `json:"revokedAtTimestamp"`
}

type VaultAddition struct {
	ID                    string          `json:"id"`
	DebtSubsidizer        *DebtSubsidizer `json:"debtSubsidizer"`
	VaultAddress          string          `json:"vaultAddress"`
	CTokenAddress         string          `json:"cTokenAddress"`
	LendingManagerAddress string          `json:"lendingManagerAddress"`
	AddedAtBlock          *big.Int        `json:"addedAtBlock"`
	AddedAtTimestamp      *big.Int        `json:"addedAtTimestamp"`
	TransactionHash       string          `json:"transactionHash"`
}

type CollectionWhitelist struct {
	ID                     string          `json:"id"`
	DebtSubsidizer         *DebtSubsidizer `json:"debtSubsidizer"`
	VaultAddress           string          `json:"vaultAddress"`
	CollectionAddress      string          `json:"collectionAddress"`
	IsActive               bool            `json:"isActive"`
	WhitelistedAtBlock     *big.Int        `json:"whitelistedAtBlock"`
	WhitelistedAtTimestamp *big.Int        `json:"whitelistedAtTimestamp"`
	RemovedAtBlock         *big.Int        `json:"removedAtBlock"`
	RemovedAtTimestamp     *big.Int        `json:"removedAtTimestamp"`
	TransactionHash        string          `json:"transactionHash"`
}

type EpochVaultAllocation struct {
	ID                    string            `json:"id"`
	Epoch                 *Epoch            `json:"epoch"`
	Vault                 *CollectionsVault `json:"vault"`
	YieldAllocated        *big.Int          `json:"yieldAllocated"`
	SubsidiesDistributed  *big.Int          `json:"subsidiesDistributed"`
	RemainingYield        *big.Int          `json:"remainingYield"`
	ParticipantCount      *big.Int          `json:"participantCount"`
	AverageSubsidyPerUser *big.Int          `json:"averageSubsidyPerUser"`
	UtilizationRate       *big.Int          `json:"utilizationRate"`
	CreatedAtBlock        *big.Int          `json:"createdAtBlock"`
	CreatedAtTimestamp    *big.Int          `json:"createdAtTimestamp"`
	UpdatedAtBlock        *big.Int          `json:"updatedAtBlock"`
	UpdatedAtTimestamp    *big.Int          `json:"updatedAtTimestamp"`
}

type CollectionYieldAccrual struct {
	ID                 string      `json:"id"`
	Collection         *Collection `json:"collection"`
	YieldAmount        *big.Int    `json:"yieldAmount"`
	CumulativeYield    *big.Int    `json:"cumulativeYield"`
	GlobalDepositIndex *big.Int    `json:"globalDepositIndex"`
	SourceType         string      `json:"sourceType"`
	SourceTransaction  string      `json:"sourceTransaction"`
	BlockNumber        *big.Int    `json:"blockNumber"`
	Timestamp          *big.Int    `json:"timestamp"`
	TransactionHash    string      `json:"transactionHash"`
}

type SystemState struct {
	ID                        string   `json:"id"`
	ActiveEpochID             string   `json:"activeEpochId"`
	TotalVaults               *big.Int `json:"totalVaults"`
	TotalCollections          *big.Int `json:"totalCollections"`
	TotalUsers                *big.Int `json:"totalUsers"`
	TotalValueLocked          *big.Int `json:"totalValueLocked"`
	TotalYieldDistributed     *big.Int `json:"totalYieldDistributed"`
	TotalSubsidiesDistributed *big.Int `json:"totalSubsidiesDistributed"`
	SystemUtilizationRate     *big.Int `json:"systemUtilizationRate"`
	AverageAPY                *big.Int `json:"averageAPY"`
	LastUpdatedBlock          *big.Int `json:"lastUpdatedBlock"`
	LastUpdatedTimestamp      *big.Int `json:"lastUpdatedTimestamp"`
}

type DailyMetrics struct {
	ID                        string   `json:"id"`
	Date                      string   `json:"date"`
	DailyVolumeUSD            *big.Int `json:"dailyVolumeUSD"`
	DailyTransactionCount     *big.Int `json:"dailyTransactionCount"`
	DailyActiveUsers          *big.Int `json:"dailyActiveUsers"`
	DailyYieldDistributed     *big.Int `json:"dailyYieldDistributed"`
	DailySubsidiesDistributed *big.Int `json:"dailySubsidiesDistributed"`
	AverageAPY                *big.Int `json:"averageAPY"`
	UtilizationRate           *big.Int `json:"utilizationRate"`
	Timestamp                 *big.Int `json:"timestamp"`
}

type QueryResponse struct {
	Data   interface{}    `json:"data"`
	Errors []GraphQLError `json:"errors,omitempty"`
}

type CollectionsVaultResponse struct {
	CollectionsVaults []*CollectionsVault `json:"collectionsVaults"`
}

type EpochManagerResponse struct {
	EpochManagers []*EpochManager `json:"epochManagers"`
}

type DebtSubsidizerResponse struct {
	DebtSubsidizers []*DebtSubsidizer `json:"debtSubsidizers"`
}

type EpochResponse struct {
	Epochs []*Epoch `json:"epochs"`
}

type CollectionParticipationResponse struct {
	CollectionParticipations []*CollectionParticipation `json:"collectionParticipations"`
}

type UserEpochEligibilityResponse struct {
	UserEpochEligibilities []*UserEpochEligibility `json:"userEpochEligibilities"`
}

type SubsidyDistributionResponse struct {
	SubsidyDistributions []*SubsidyDistribution `json:"subsidyDistributions"`
}

type MerkleDistributionResponse struct {
	MerkleDistributions []*MerkleDistribution `json:"merkleDistributions"`
}
