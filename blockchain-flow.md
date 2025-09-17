# OpenEdu Blockchain Architecture Documentation

## 🎯 Overview

OpenEdu blockchain system consists of two main services (gateway):
- **openedu-core**: Main business logic, user management
- **openedu-blockchain**: Blockchain operations, wallet management, transaction processing

### Key Design Principles
- **Microservice Architecture**: Separated concerns between business logic and blockchain
- **Async Communication**: RabbitMQ message queues for service-to-service communication
- **Multi-Network Support**: NEAR
- **Security First**: Private key encryption, secure transaction handling

## 🏗️ Service Architecture

```
┌─────────────────┐                 ┌──────────────────────┐
│   openedu-core  │                 │ openedu-blockchain   │
│                 │                 │                      │
│                 │                 │ • Wallet Management  │
│ • Business Logic│                 │ • Transaction Proc.  │
│ • User Mgmt     │                 │ • NFT Minting        │
│ • Course Mgmt   │                 │ • Launchpad          │
└─────────────────┘                 └──────────────────────┘
         │                                       │
         │              RabbitMQ                 │
         └──────────── Queues Only ──────────────┘
         │                                       │
         ▼                                       ▼
┌─────────────────┐                 ┌──────────────────────┐
│   PostgreSQL    │                 │   PostgreSQL         │
│   (Core DB)     │                 │   (Blockchain DB)    │
└─────────────────┘                 └──────────────────────┘
```

### 🔒 **Security-First Communication**
- **NO REST APIs** between services
- **RabbitMQ message queues ONLY** for service-to-service communication
- **Async processing** for all blockchain operations
- **Complete isolation** between business logic and blockchain logic

### Database Separation
- **Core DB**: Users, courses, organizations, business data
- **Blockchain DB**: Wallets, transactions, sponsor wallets, blockchain-specific data


## 🔄 Communication Flow

### 1. Message Queue Pattern

**Core → Blockchain:**
```go
// openedu-core/pkg/openedu_chain/
producer.Publish(queueName, message)
```

**Blockchain → Core:**
```go
// openedu-blockchain/pkg/openedu_core/
producer.Publish(queueName, syncMessage)
```

### 2. Queue Types

#### 📤 **Core → Blockchain Queues**
| Queue Name                    | Purpose                | Handler                   |
| ----------------------------- | ---------------------- | ------------------------- |
| `wallet_create_queue`         | Create new wallets     | `CreateWallet`            |
| `sponsor_wallet_create_queue` | Create sponsor wallets | `CreateSponsorWallet`     |
| `transaction_create_queue`    | Process transactions   | `CreateTransaction`       |
| `launchpad_rpc_query`         | Launchpad operations   | `HandleLaunchpadRPCQuery` |
| `wallet_rpc_query`            | Wallet queries         | `HandleWalletRPCQuery`    |

#### 📥 **Blockchain → Core Queues**
| Queue Name                          | Purpose                 | Handler                       |
| ----------------------------------- | ----------------------- | ----------------------------- |
| `wallet_sync_queue`                 | Sync wallet data        | `SyncWallet`                  |
| `transaction_sync_queue`            | Sync transaction status | `SyncTransaction`             |
| `wallet_retrieve_get_details_queue` | Wallet details          | `HandleRetrieveWalletDetails` |

### 3. RPC vs Fire-and-Forget

**RPC Pattern** (Request-Response):
```go
// For queries that need immediate response
repliedMsg, err := producer.PublishRPC(queueName, message)
```

**Fire-and-Forget Pattern**:
```go
// For async operations
err := producer.Publish(queueName, message)
```

## 🚀 Blockchain Features

### 1. 🎯 NFT Minting System

**Supported Networks:**
- **NEAR**: Native NEAR NFT contracts

**Gas Fee Payers:**
- `Platform`: OpenEdu pays gas fees
- `Learner`: Student pays gas fees
- `Creator`: Course creator pays via sponsor wallet
- `Paymaster`: Coinbase Paymaster (BASE only)

**Transaction Types:**
- `mint_nft`: Standard NFT minting
- `mint_nft_with_permit`: EIP-712 permit-based minting (BASE)

**Flow:**
```
1. User completes course
2. Core validates completion
3. Core → Queue: MintNFT request
4. Blockchain processes minting
5. Blockchain → Queue: Transaction sync
6. Core updates certificate status
```

Flow Diagram: https://www.canva.com/design/DAGlcCklLrY/6WJDomdRtdigTHaef62e6A/view?utm_content=DAGlcCklLrY&utm_campaign=designshare&utm_medium=link2&utm_source=uniquelinks&utlId=h5954fa7160

### 2. 💰 Wallet Management

**Wallet Types:**
- **Crypto Wallets**: NEAR, USDT, USDC, OpenEdu token
- **Network**: NEAR

**Supported Operations:**
- Wallet creation and sync
- Balance tracking across networks
- Private key encryption with KMS
- Account info retrieval

**Multi-Network Support:**
```go
// Strategy pattern for different networks
type WalletStrategy interface {
    CreateWallet(req *CreateWalletRequest) (*CreateWalletResponse, error)
    GetBalance(address string) (decimal.Decimal, error)
}
```

### 3. 💸 Transfer System

**Transfer Types:**
- **Single Transfer**: One-to-one transfers
- **Batch Transfer**: One-to-many transfers
- **Cross-Network**: Transfer between different blockchains

**Supported Tokens:**
- **Native Tokens**: NEAR (NEAR)
- **Fungible Tokens**: USDT, USDC, custom tokens

**Flow:**
```
1. Core validates transfer request
2. Core → Queue: Transfer request
3. Blockchain selects appropriate strategy
4. Strategy executes transfer on blockchain
5. Blockchain → Queue: Transaction sync
6. Core updates balances
```

### 4. 🏦 Sponsor Wallet System

**Purpose**: Allow course creators to sponsor gas fees for their students

**Features:**
- **Token deposits** for gas fee sponsoring
- **Balance tracking** and management
- **Automatic fallback** to paymaster if insufficient balance

**Operations:**
- `deposit_sponsor_gas`: Add funds to sponsor wallet
- `withdraw_sponsor_gas`: Remove funds from sponsor wallet
- `init_sponsor_wallet`: Initialize new sponsor wallet

**Flow:**
```
1. Creator creates sponsor wallet
2. Creator deposits token
3. Student mints NFT
4. Gas fee deducted from sponsor wallet
```

### 5. 💳 Payment System

**Purpose**: Handle crypto payments for courses and services

**Features:**
- **Multi-token support**: USDT, USDC, custom tokens
- **Network**: NEAR
- **Automatic conversion**: Token to fiat equivalent
- **Payment validation**: Balance and network checks

**Transaction Type:**
- `payment`: Process payment transactions

### 6. 🎁 Earning Claims System

**Purpose**: Allow users to claim earned rewards and tokens

**Features:**
- **Reward distribution**: Automatic token distribution
- **Claim validation**: Verify eligibility before claiming
- **Multi-token support**: Various reward tokens
- **Batch claims**: Multiple rewards in single transaction

**Transaction Type:**
- `claim_earning`: Process earning claims

### 7. 🚀 Launchpad System

**Purpose**: Crowdfunding platform for educational projects

**Features:**
- **Pool creation** with funding goals
- **Milestone-based** fund release
- **Voting mechanisms** for milestone approval
- **Refund system** if goals not met
- **Multi-network support**

**Transaction Types:**
- `init_launchpad_pool`: Create new funding pool
- `approve_launchpad_pool`: Approve pool for funding
- `pledge_launchpad`: Pledge funds to pool
- `withdraw_launchpad_fund_to_creator`: Release funds to creator
- `claim_launchpad_refund`: Claim refund if pool fails

Flow Diagram: https://www.canva.com/design/DAGYH252Uqg/oNyh-VFamfIUs9vwaxL3sw/view?utm_content=DAGYH252Uqg&utm_campaign=designshare&utm_medium=link2&utm_source=uniquelinks&utlId=h94d9465a34

### Data Flow Between Services

```
1. Core creates wallet record (blockchain_wallet_id = null)
2. Core → Queue: Create wallet request
3. Blockchain creates wallet with private key
4. Blockchain → Queue: Wallet sync with blockchain_wallet_id
5. Core updates wallet record with blockchain_wallet_id
```

## 🛠️ Development Guide

### Adding New Blockchain Network

1. **Create Strategy Implementation:**
```go
// openedu-blockchain/services/nft/strategies/new_network.go
type NewNetworkMintNftService struct{}

func (s *NewNetworkMintNftService) MintNFT(account Account, req *MintNftRequest) (*MintNftResponse, error) {
    // Implementation
}
```

2. **Register Strategy:**
```go
// openedu-blockchain/services/services.go
nftSvc.RegisterStrategy(models.NetworkNEW, &strategies.NewNetworkMintNftService{})
```

3. **Add Network Constants:**
```go
// models/constant.go
NetworkNEW BlockchainNetwork = "new_network"
```

### Adding New Transaction Type

1. **Define DTO:**
```go
// dto/transaction.go
type NewTransactionRequest struct {
    // Fields
}
```

2. **Implement Service Method:**
```go
// services/transaction.go
func (s *TransactionService) NewTransaction(req *dto.NewTransactionRequest) (*models.Transaction, *e.AppError) {
    // Implementation
}
```

## 📝 Summary

This OpenEdu blockchain architecture provides:

### ✅ **Comprehensive Features**
- **7 Major Blockchain Features**: NFT minting, transfers, payments, sponsor wallets, earning claims, launchpad, wallet management
- **Network Support**: NEAR
- **Multiple Transaction Types**: 12+ different transaction types supported
- **Flexible Gas Fee Management**: Platform, learner, creator, and paymaster options