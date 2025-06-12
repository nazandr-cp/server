# EpochScheduler Go Server

This Go server implements the EpochScheduler service for monitoring and managing epoch transitions in the lend.fam ecosystem.

## Features

- Monitors blockchain for epoch boundaries
- Interacts with EpochManager smart contract
- Automatically starts new epochs when current epochs end
- Handles first epoch initialization
- Configurable polling intervals
- Graceful shutdown handling

## Configuration

The scheduler is configured using environment variables:

- `RPC_URL`: Ethereum node RPC URL (default: `http://localhost:8545`)
- `EPOCH_MANAGER_ADDRESS`: Address of the deployed EpochManager contract (required)
- `PRIVATE_KEY`: Private key of the account that owns the EpochManager contract (required)
- `POLLING_INTERVAL`: Polling interval in seconds (default: `30`)

## Usage

### Running the scheduler

```bash
# Set environment variables
export RPC_URL="https://eth-mainnet.alchemyapi.io/v2/your-api-key"
export EPOCH_MANAGER_ADDRESS="0x1234567890123456789012345678901234567890"
export PRIVATE_KEY="your-private-key-without-0x-prefix"
export POLLING_INTERVAL="30"

# Run the scheduler
go run main.go
```

### Building the binary

```bash
go build -o epochscheduler main.go
./epochscheduler
```

## Architecture

### Components

1. **EpochManager Interface** (`contracts/epochmanager/epochmanager.go`)
   - Defines the interface for interacting with the EpochManager smart contract
   - Contains placeholder implementation (to be replaced with abigen-generated bindings)

2. **Scheduler** (`epochscheduler/scheduler.go`)
   - Main scheduler logic
   - Monitors epoch transitions
   - Handles transaction submission for starting new epochs

3. **Main** (`main.go`)
   - Application entry point
   - Configuration management
   - Graceful shutdown handling

### Key Methods

- `NewScheduler()`: Creates and initializes a new scheduler instance
- `Run()`: Main scheduler loop with configurable polling
- `checkAndProcessEpochTransition()`: Checks if epoch transition is needed
- `startFirstEpoch()`: Handles the first epoch initialization
- `attemptToStartNewEpoch()`: Starts a new epoch after the previous one ends

## Dependencies

- `github.com/ethereum/go-ethereum`: Ethereum client library
- Standard Go libraries for time, crypto, and system operations

## TODO

- Replace placeholder EpochManager interface with actual abigen-generated bindings
- Implement automated processing pipeline triggers
- Add comprehensive error handling and retry logic
- Add metrics and monitoring
- Add configuration file support
- Add tests