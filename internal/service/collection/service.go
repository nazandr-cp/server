package collection

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"go.uber.org/zap"

	config "go-server/configs"
	"go-server/contracts"
	"go-server/internal/gql"
	"go-server/internal/platform/graphql"
)

// Service manages collection operations and data.
type Service struct {
	subgraphURL      string
	registryContract *contracts.ICollectionRegistry
	ethClient        *ethclient.Client
	logger           *zap.Logger
	config           config.Config
}

func NewService(cfg config.Config, logger *zap.Logger) (*Service, error) {
	ethClient, err := ethclient.Dial(cfg.RPCHTTPURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum client: %w", err)
	}

	registryContract := contracts.NewICollectionRegistry()

	return &Service{
		subgraphURL:      cfg.SubgraphURL,
		registryContract: registryContract,
		ethClient:        ethClient,
		logger:           logger,
		config:           cfg,
	}, nil
}

func (s *Service) GetCollectionConfig(ctx context.Context, collectionAddr common.Address) (*contracts.ICollectionRegistryCollection, error) {
	s.logger.Info("Getting collection config", zap.String("collection", collectionAddr.Hex()))

	// Mock response for now - would need actual registry address in config
	collection := &contracts.ICollectionRegistryCollection{
		CollectionAddress:    collectionAddr,
		CollectionType:       0, // ERC721
		YieldSharePercentage: 100,
		Vaults:               []common.Address{},
	}

	return collection, nil
}

func (s *Service) GetCollectionParticipants(ctx context.Context, collectionAddr common.Address) ([]*gql.CollectionParticipation, error) {
	query := gql.GetCollectionParticipationsQuery()

	variables := map[string]interface{}{
		"first": 1000,
		"skip":  0,
		"where": map[string]interface{}{
			"collection": collectionAddr.Hex(),
		},
	}

	client := graphql.NewClient(s.subgraphURL)
	var result gql.CollectionParticipationResponse
	err := client.QueryWithVariables(ctx, query, variables, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to query collection participations: %w", err)
	}

	return result.CollectionParticipations, nil
}

func (s *Service) CalculateCollectionMetrics(ctx context.Context, collectionAddr common.Address) (*CollectionMetrics, error) {
	participations, err := s.GetCollectionParticipants(ctx, collectionAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to get collection participants: %w", err)
	}

	metrics := &CollectionMetrics{
		CollectionAddress: collectionAddr,
		TotalParticipants: big.NewInt(int64(len(participations))),
		TotalYieldAccrued: big.NewInt(0),
		TotalSubsidies:    big.NewInt(0),
		AverageAPY:        big.NewInt(0),
	}

	var totalAPY *big.Int = big.NewInt(0)
	participantCount := int64(0)

	for _, participation := range participations {
		if participation.YieldAccrued != nil {
			metrics.TotalYieldAccrued.Add(metrics.TotalYieldAccrued, participation.YieldAccrued)
		}
		if participation.TotalSubsidies != nil {
			metrics.TotalSubsidies.Add(metrics.TotalSubsidies, participation.TotalSubsidies)
		}
		if participation.AverageAPY != nil && participation.AverageAPY.Cmp(big.NewInt(0)) > 0 {
			totalAPY.Add(totalAPY, participation.AverageAPY)
			participantCount++
		}
	}

	if participantCount > 0 {
		metrics.AverageAPY.Div(totalAPY, big.NewInt(participantCount))
	}

	s.logger.Info("Calculated collection metrics",
		zap.String("collection", collectionAddr.Hex()),
		zap.String("total_participants", metrics.TotalParticipants.String()),
		zap.String("total_yield", metrics.TotalYieldAccrued.String()),
		zap.String("average_apy", metrics.AverageAPY.String()))

	return metrics, nil
}

func (s *Service) GetAllActiveCollections(ctx context.Context) ([]common.Address, error) {
	s.logger.Info("Getting all active collections")

	// Mock response for now - would need actual registry address in config
	collections := []common.Address{}

	return collections, nil
}

func (s *Service) Close() {
	if s.ethClient != nil {
		s.ethClient.Close()
	}
	if s.logger != nil {
		s.logger.Sync()
	}
}

type CollectionMetrics struct {
	CollectionAddress common.Address
	TotalParticipants *big.Int
	TotalYieldAccrued *big.Int
	TotalSubsidies    *big.Int
	AverageAPY        *big.Int
}
