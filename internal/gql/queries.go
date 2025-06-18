package gql

func GetCollectionsVaultQuery() string {
	return `{
		collectionsVaults {
			id
			totalShares
			totalDeposits
			totalCTokens
			globalDepositIndex
			totalPrincipalDeposited
			cTokenMarket {
				id
				symbol
				name
				decimals
				totalSupply
				totalBorrows
				exchangeRate
				collateralFactor
			}
			collectionRegistry {
				id
				totalCollections
				totalActiveCollections
				owner
			}
			epochManager {
				id
				currentEpochId
				epochDuration
				totalEpochs
			}
			debtSubsidizer {
				id
				totalSubsidyPool
				totalSubsidiesDistributed
				totalEligibleUsers
			}
			createdAtBlock
			createdAtTimestamp
			updatedAtBlock
			updatedAtTimestamp
		}
	}`
}

func GetCollectionsVaultByIDQuery() string {
	return `query GetCollectionsVault($id: ID!) {
		collectionsVault(id: $id) {
			id
			totalShares
			totalDeposits
			totalCTokens
			globalDepositIndex
			totalPrincipalDeposited
			cTokenMarket {
				id
				symbol
				name
				decimals
				totalSupply
				totalBorrows
				exchangeRate
				collateralFactor
				liquidationIncentive
				reserveFactor
				baseRatePerBlock
				multiplierPerBlock
				jumpMultiplierPerBlock
				kink
				lastExchangeRateTimestamp
			}
			collectionRegistry {
				id
				totalCollections
				totalActiveCollections
				owner
				collections {
					id
					contractAddress
					name
					symbol
					isActive
					yieldSharePercentage
					weightFunctionType
					totalNFTsDeposited
					totalBorrowVolume
				}
			}
			epochManager {
				id
				currentEpochId
				currentEpoch {
					id
					epochNumber
					status
					startTimestamp
					endTimestamp
					totalYieldAvailable
					totalYieldDistributed
					participantCount
				}
				epochDuration
				processingBuffer
				minimumYieldThreshold
				totalEpochs
				totalYieldDistributed
			}
			debtSubsidizer {
				id
				totalSubsidyPool
				totalSubsidiesDistributed
				totalSubsidiesRemaining
				totalEligibleUsers
				subsidyRate
				maxSubsidyPerUser
				subsidyDuration
			}
			lendingManager {
				id
				supportedMarkets
				totalMarketParticipants
				globalCollateralFactor
				liquidationIncentive
				totalSupplyVolume
				totalBorrowVolume
				totalLiquidationVolume
			}
			collectionParticipations {
				id
				collection {
					id
					contractAddress
					name
					symbol
					collectionType
				}
				principalShares
				principalDeposited
				yieldAccrued
				totalSubsidies
				rewardSharePercentage
				weightFunctionType
				averageAPY
				totalParticipants
			}
			createdAtBlock
			createdAtTimestamp
			updatedAtBlock
			updatedAtTimestamp
		}
	}`
}

func GetEpochManagerQuery() string {
	return `{
		epochManagers {
			id
			currentEpochId
			currentEpoch {
				id
				epochNumber
				status
				startTimestamp
				endTimestamp
				totalYieldAvailable
				totalYieldDistributed
				participantCount
			}
			epochDuration
			processingBuffer
			minimumYieldThreshold
			totalEpochs
			totalYieldDistributed
			owner
			createdAtBlock
			createdAtTimestamp
			updatedAtBlock
			updatedAtTimestamp
		}
	}`
}

func GetEpochManagerByIDQuery() string {
	return `query GetEpochManager($id: ID!) {
		epochManager(id: $id) {
			id
			currentEpochId
			currentEpoch {
				id
				epochNumber
				status
				startTimestamp
				endTimestamp
				processingStartedTimestamp
				processingCompletedTimestamp
				totalYieldAvailable
				totalYieldAllocated
				totalYieldDistributed
				remainingYield
				totalSubsidiesDistributed
				totalEligibleUsers
				totalParticipatingCollections
				participantCount
				processingTimeMs
				estimatedProcessingTime
				processingGasUsed
				processingTransactionCount
			}
			epochDuration
			processingBuffer
			minimumYieldThreshold
			totalEpochs
			totalYieldDistributed
			owner
			epochs(orderBy: epochNumber, orderDirection: desc, first: 10) {
				id
				epochNumber
				status
				startTimestamp
				endTimestamp
				totalYieldDistributed
				participantCount
			}
			vault {
				id
				totalShares
				totalDeposits
			}
			createdAtBlock
			createdAtTimestamp
			updatedAtBlock
			updatedAtTimestamp
		}
	}`
}

func GetDebtSubsidizerQuery() string {
	return `{
		debtSubsidizers {
			id
			totalSubsidyPool
			totalSubsidiesDistributed
			totalSubsidiesRemaining
			totalEligibleUsers
			subsidyRate
			maxSubsidyPerUser
			subsidyDuration
			owner
			createdAtBlock
			createdAtTimestamp
			updatedAtBlock
			updatedAtTimestamp
		}
	}`
}

func GetDebtSubsidizerByIDQuery() string {
	return `query GetDebtSubsidizer($id: ID!) {
		debtSubsidizer(id: $id) {
			id
			totalSubsidyPool
			totalSubsidiesDistributed
			totalSubsidiesRemaining
			totalEligibleUsers
			subsidyRate
			maxSubsidyPerUser
			subsidyDuration
			owner
			vaults {
				id
				totalShares
				totalDeposits
				globalDepositIndex
			}
			subsidyDistributions(orderBy: timestamp, orderDirection: desc, first: 100) {
				id
				epoch {
					id
					epochNumber
				}
				user {
					id
				}
				collection {
					id
					contractAddress
					name
				}
				subsidyAmount
				borrowAmountBefore
				borrowAmountAfter
				nftBalance
				weightedContribution
				timestamp
				transactionHash
			}
			whitelistedCollections {
				id
				vaultAddress
				collectionAddress
				isActive
				whitelistedAtTimestamp
			}
			createdAtBlock
			createdAtTimestamp
			updatedAtBlock
			updatedAtTimestamp
		}
	}`
}

func GetEpochsQuery() string {
	return `query GetEpochs($first: Int, $skip: Int, $orderBy: String, $orderDirection: String) {
		epochs(first: $first, skip: $skip, orderBy: $orderBy, orderDirection: $orderDirection) {
			id
			epochNumber
			status
			startTimestamp
			endTimestamp
			totalYieldAvailable
			totalYieldDistributed
			totalSubsidiesDistributed
			totalEligibleUsers
			totalParticipatingCollections
			participantCount
			processingTimeMs
			epochManager {
				id
				currentEpochId
			}
			createdAtBlock
			createdAtTimestamp
			updatedAtBlock
			updatedAtTimestamp
		}
	}`
}

func GetEpochByIDQuery() string {
	return `query GetEpoch($id: ID!) {
		epoch(id: $id) {
			id
			epochNumber
			status
			startTimestamp
			endTimestamp
			processingStartedTimestamp
			processingCompletedTimestamp
			totalYieldAvailable
			totalYieldAllocated
			totalYieldDistributed
			remainingYield
			totalSubsidiesDistributed
			totalEligibleUsers
			totalParticipatingCollections
			participantCount
			processingTimeMs
			estimatedProcessingTime
			processingGasUsed
			processingTransactionCount
			epochManager {
				id
				epochDuration
				processingBuffer
			}
			vaultAllocations {
				id
				vault {
					id
				}
				yieldAllocated
				subsidiesDistributed
				remainingYield
				participantCount
				averageSubsidyPerUser
				utilizationRate
			}
			subsidyDistributions {
				id
				user {
					id
				}
				collection {
					id
					contractAddress
					name
				}
				subsidyAmount
				nftBalance
				weightedContribution
			}
			merkleDistributions {
				id
				vault {
					id
				}
				merkleRoot
				totalAmount
				totalClaims
			}
			userEligibilities {
				id
				user {
					id
				}
				collection {
					id
					contractAddress
					name
				}
				nftBalance
				borrowBalance
				holdingDuration
				isEligible
				subsidyReceived
				yieldShare
			}
			createdAtBlock
			createdAtTimestamp
			updatedAtBlock
			updatedAtTimestamp
		}
	}`
}

func GetCurrentEpochQuery() string {
	return `{
		epochs(orderBy: epochNumber, orderDirection: desc, first: 1, where: {status: ACTIVE}) {
			id
			epochNumber
			status
			startTimestamp
			endTimestamp
			totalYieldAvailable
			totalYieldDistributed
			participantCount
			epochManager {
				id
				epochDuration
			}
		}
	}`
}

func GetCollectionParticipationsQuery() string {
	return `query GetCollectionParticipations($first: Int, $skip: Int, $where: CollectionParticipation_filter) {
		collectionParticipations(first: $first, skip: $skip, where: $where) {
			id
			collection {
				id
				contractAddress
				name
				symbol
				collectionType
				isActive
				yieldSharePercentage
				weightFunctionType
			}
			vault {
				id
				totalShares
				totalDeposits
			}
			principalShares
			principalDeposited
			totalCTokens
			globalDepositIndex
			yieldAccrued
			yieldClaimed
			totalYieldGenerated
			isBorrowBased
			rewardSharePercentage
			weightFunctionType
			secondsAccumulated
			secondsClaimed
			totalSubsidies
			totalSubsidiesClaimed
			averageAPY
			totalParticipants
			createdAtBlock
			createdAtTimestamp
			updatedAtBlock
			updatedAtTimestamp
		}
	}`
}

func GetCollectionParticipationByIDQuery() string {
	return `query GetCollectionParticipation($id: ID!) {
		collectionParticipation(id: $id) {
			id
			collection {
				id
				contractAddress
				name
				symbol
				totalSupply
				collectionType
				isActive
				yieldSharePercentage
				weightFunctionType
				weightFunctionP1
				weightFunctionP2
				minBorrowAmount
				maxBorrowAmount
				totalNFTsDeposited
				totalBorrowVolume
				totalYieldGenerated
				totalSubsidiesReceived
			}
			vault {
				id
				totalShares
				totalDeposits
				totalCTokens
				globalDepositIndex
			}
			principalShares
			principalDeposited
			totalCTokens
			globalDepositIndex
			lastGlobalDepositIndex
			yieldAccrued
			yieldClaimed
			totalYieldGenerated
			isBorrowBased
			rewardSharePercentage
			weightFunctionType
			weightFunctionP1
			weightFunctionP2
			secondsAccumulated
			secondsClaimed
			totalSubsidies
			totalSubsidiesClaimed
			averageAPY
			totalParticipants
			accountSubsidies {
				id
				account {
					id
				}
				balanceNFT
				weightedBalance
				secondsAccumulated
				secondsClaimed
				subsidiesAccrued
				subsidiesClaimed
			}
			createdAtBlock
			createdAtTimestamp
			updatedAtBlock
			updatedAtTimestamp
		}
	}`
}

func GetUserEpochEligibilitiesQuery() string {
	return `query GetUserEpochEligibilities($first: Int, $skip: Int, $where: UserEpochEligibility_filter) {
		userEpochEligibilities(first: $first, skip: $skip, where: $where) {
			id
			user {
				id
				totalSubsidiesReceived
				totalYieldEarned
				totalNFTsOwned
			}
			epoch {
				id
				epochNumber
				status
				startTimestamp
				endTimestamp
			}
			collection {
				id
				contractAddress
				name
				symbol
				collectionType
			}
			nftBalance
			borrowBalance
			holdingDuration
			isEligible
			subsidyReceived
			yieldShare
			bonusMultiplier
			calculatedAtBlock
			calculatedAtTimestamp
		}
	}`
}

func GetUserEpochEligibilityByUserAndEpochQuery() string {
	return `query GetUserEpochEligibility($userId: ID!, $epochId: ID!) {
		userEpochEligibilities(where: {user: $userId, epoch: $epochId}) {
			id
			user {
				id
				totalSubsidiesReceived
				totalYieldEarned
				totalBorrowVolume
				totalNFTsOwned
				totalCollectionsParticipated
			}
			epoch {
				id
				epochNumber
				status
				startTimestamp
				endTimestamp
				totalYieldDistributed
				totalSubsidiesDistributed
			}
			collection {
				id
				contractAddress
				name
				symbol
				collectionType
				isActive
				yieldSharePercentage
				weightFunctionType
			}
			nftBalance
			borrowBalance
			holdingDuration
			isEligible
			subsidyReceived
			yieldShare
			bonusMultiplier
			calculatedAtBlock
			calculatedAtTimestamp
		}
	}`
}

func GetSubsidyDistributionsQuery() string {
	return `query GetSubsidyDistributions($first: Int, $skip: Int, $where: SubsidyDistribution_filter) {
		subsidyDistributions(first: $first, skip: $skip, where: $where, orderBy: timestamp, orderDirection: desc) {
			id
			epoch {
				id
				epochNumber
				status
			}
			debtSubsidizer {
				id
				totalSubsidyPool
			}
			user {
				id
				totalSubsidiesReceived
			}
			collection {
				id
				contractAddress
				name
				symbol
			}
			vault {
				id
			}
			subsidyAmount
			borrowAmountBefore
			borrowAmountAfter
			nftBalance
			weightedContribution
			gasUsed
			blockNumber
			timestamp
			transactionHash
		}
	}`
}

func GetSubsidyDistributionsByEpochQuery() string {
	return `query GetSubsidyDistributionsByEpoch($epochId: ID!, $first: Int, $skip: Int) {
		subsidyDistributions(
			where: {epoch: $epochId}, 
			first: $first, 
			skip: $skip, 
			orderBy: subsidyAmount, 
			orderDirection: desc
		) {
			id
			epoch {
				id
				epochNumber
			}
			user {
				id
				totalSubsidiesReceived
				totalYieldEarned
			}
			collection {
				id
				contractAddress
				name
				symbol
				collectionType
			}
			vault {
				id
			}
			subsidyAmount
			borrowAmountBefore
			borrowAmountAfter
			nftBalance
			weightedContribution
			gasUsed
			blockNumber
			timestamp
			transactionHash
		}
	}`
}

func GetMerkleDistributionsQuery() string {
	return `query GetMerkleDistributions($first: Int, $skip: Int, $where: MerkleDistribution_filter) {
		merkleDistributions(first: $first, skip: $skip, where: $where, orderBy: timestamp, orderDirection: desc) {
			id
			epoch {
				id
				epochNumber
				status
			}
			vault {
				id
				totalShares
				totalDeposits
			}
			merkleRoot
			totalAmount
			totalClaims
			blockNumber
			timestamp
			transactionHash
		}
	}`
}

func GetMerkleDistributionByEpochAndVaultQuery() string {
	return `query GetMerkleDistribution($epochId: ID!, $vaultId: ID!) {
		merkleDistributions(where: {epoch: $epochId, vault: $vaultId}) {
			id
			epoch {
				id
				epochNumber
				status
				startTimestamp
				endTimestamp
			}
			vault {
				id
				totalShares
				totalDeposits
				globalDepositIndex
			}
			merkleRoot
			totalAmount
			totalClaims
			blockNumber
			timestamp
			transactionHash
		}
	}`
}

func GetAccountQuery() string {
	return `query GetAccount($id: ID!) {
		account(id: $id) {
			id
			totalSecondsClaimed
			totalSubsidiesReceived
			totalYieldEarned
			totalBorrowVolume
			totalNFTsOwned
			totalCollectionsParticipated
			markets {
				id
				cTokenMarket {
					id
					symbol
					name
				}
				supplyBalance
				borrowBalance
				collateralBalance
			}
			accountSubsidies {
				id
				collection {
					id
					contractAddress
					name
				}
				balanceNFT
				secondsAccumulated
				subsidiesAccrued
				subsidiesClaimed
			}
			userEligibilities {
				id
				epoch {
					id
					epochNumber
				}
				collection {
					id
					name
				}
				isEligible
				subsidyReceived
			}
			firstInteractionBlock
			firstInteractionTimestamp
			updatedAtBlock
			updatedAtTimestamp
		}
	}`
}

func GetSystemStateQuery() string {
	return `{
		systemStates {
			id
			activeEpochId
			totalVaults
			totalCollections
			totalUsers
			totalValueLocked
			totalYieldDistributed
			totalSubsidiesDistributed
			systemUtilizationRate
			averageAPY
			lastUpdatedBlock
			lastUpdatedTimestamp
		}
	}`
}

func GetDailyMetricsQuery() string {
	return `query GetDailyMetrics($first: Int, $skip: Int, $orderBy: String, $orderDirection: String) {
		dailyMetrics(first: $first, skip: $skip, orderBy: $orderBy, orderDirection: $orderDirection) {
			id
			date
			dailyVolumeUSD
			dailyTransactionCount
			dailyActiveUsers
			dailyYieldDistributed
			dailySubsidiesDistributed
			averageAPY
			utilizationRate
			timestamp
		}
	}`
}

func GetCollectionDepositsQuery() string {
	return `query GetCollectionDeposits($first: Int, $skip: Int, $where: Deposit_filter) {
		deposits(first: $first, skip: $skip, where: $where) {
			id
			user {
				id
			}
			collection {
				id
				contractAddress
			}
			vault {
				id
			}
			amount
			shares
			timestamp
			transactionHash
		}
	}`
}

func GetBorrowsQuery() string {
	return `query GetBorrows($first: Int, $skip: Int, $where: Borrow_filter) {
		borrows(first: $first, skip: $skip, where: $where) {
			id
			user {
				id
			}
			collection {
				id
				contractAddress
			}
			vault {
				id
			}
			amount
			timestamp
			transactionHash
		}
	}`
}

func GetSubsidyClaimsQuery() string {
	return `query GetSubsidyClaims($first: Int, $skip: Int, $where: SubsidyClaim_filter) {
		subsidyClaims(first: $first, skip: $skip, where: $where) {
			id
			user {
				id
			}
			epoch {
				id
				epochNumber
			}
			collection {
				id
				contractAddress
			}
			amount
			timestamp
			transactionHash
		}
	}`
}

func GetEpochQuery() string {
	return `query GetEpoch($id: ID!) {
		epoch(id: $id) {
			id
			epochNumber
			status
			startTimestamp
			endTimestamp
			totalYieldAvailable
			totalYieldDistributed
			participantCount
		}
	}`
}

func GetUserBalanceQuery() string {
	return `query GetUserBalance($id: ID!) {
		account(id: $id) {
			id
			totalYieldEarned
			totalSubsidiesReceived
			totalBorrowVolume
			totalNFTsOwned
			vaultBalance
		}
	}`
}

type QueryBuilder struct {
	query     string
	variables map[string]interface{}
}

func NewQueryBuilder(query string) *QueryBuilder {
	return &QueryBuilder{
		query:     query,
		variables: make(map[string]interface{}),
	}
}

func (qb *QueryBuilder) AddVariable(key string, value interface{}) *QueryBuilder {
	qb.variables[key] = value
	return qb
}

func (qb *QueryBuilder) Build() (string, map[string]interface{}) {
	return qb.query, qb.variables
}

func BuildPaginationVars(first, skip int) map[string]interface{} {
	return map[string]interface{}{
		"first": first,
		"skip":  skip,
	}
}

func BuildOrderVars(orderBy, orderDirection string) map[string]interface{} {
	return map[string]interface{}{
		"orderBy":        orderBy,
		"orderDirection": orderDirection,
	}
}

func BuildFilterVars(filters map[string]interface{}) map[string]interface{} {
	return filters
}

func BuildListQueryVars(first, skip int, orderBy, orderDirection string, filters map[string]interface{}) map[string]interface{} {
	vars := make(map[string]interface{})

	if first > 0 {
		vars["first"] = first
	}
	if skip > 0 {
		vars["skip"] = skip
	}
	if orderBy != "" {
		vars["orderBy"] = orderBy
	}
	if orderDirection != "" {
		vars["orderDirection"] = orderDirection
	}
	if len(filters) > 0 {
		vars["where"] = filters
	}

	return vars
}
