# Go Server Project Structure

This Go server follows the Standard Go Project Layout with industry best practices for a monolithic architecture.

## Directory Structure

```
go-server/
├── cmd/                          # Main applications
│   └── server/                   # Server application entry point
│       └── main.go              # Application main function
├── configs/                      # Configuration packages
│   └── config.go                # Configuration loading and management
├── internal/                     # Private application code
│   ├── api/                     # API layer
│   │   ├── handlers/            # HTTP handlers
│   │   └── middleware/          # HTTP middleware
│   ├── platform/                # Platform-specific code
│   │   ├── ethereum/            # Ethereum blockchain integration
│   │   ├── graphql/             # GraphQL client utilities
│   │   └── websocket/           # WebSocket functionality
│   └── service/                 # Business logic services
│       ├── datacollector/       # Data collection service
│       ├── epoch/               # Epoch management service
│       └── subsidy/             # Subsidy distribution service
├── pkg/                         # Public packages (reusable)
│   ├── merkletree/              # Merkle tree utilities
│   └── scheduler/               # Generic scheduler utilities
├── contracts/                   # Smart contract bindings
├── deployments/                 # Deployment configurations
└── docs/                        # Documentation
```

## Key Principles

### Separation of Concerns
- **cmd/**: Application entry points
- **internal/**: Private business logic (cannot be imported by external packages)
- **pkg/**: Reusable public packages
- **configs/**: Configuration management

### API Layer
- **handlers/**: HTTP request handlers
- **middleware/**: Cross-cutting concerns (logging, CORS, auth)

### Platform Layer
- **ethereum/**: Blockchain interaction utilities
- **graphql/**: GraphQL client functionality
- **websocket/**: Real-time communication

### Service Layer
- **datacollector/**: Aggregates blockchain data
- **epoch/**: Manages epoch lifecycle
- **subsidy/**: Handles reward distribution

### Package Organization
- Each package has a single responsibility
- Clear dependency hierarchy (no circular dependencies)
- Platform-specific code is isolated

## Import Strategy
- Use package aliases to maintain backward compatibility
- Prefer explicit imports over dot imports
- Group imports: standard library, external packages, internal packages

## Best Practices
- All packages follow Go naming conventions
- Error handling follows Go idioms
- Configuration is centralized
- Logging is structured (using zap)
- Graceful shutdown is implemented