# insights-edge
A blockchain-based platform delivering actionable insights to small businesses by analyzing trends and competitor data. Users earn crypto tokens for data contributions, fostering smarter decisions and growth.

##  1. Features  
- **Actionable Insights**: Aggregates and analyzes public data to provide personalized recommendations for small businesses.  
- **Token Rewards**: Users earn tokens for contributing or engaging with the platform.  
- **Secure and Transparent**: Built on blockchain for enhanced security and data integrity.  
- **Intuitive Interface**: A user-friendly interface with dashboard and marketplace features.
- **Business Dashboard**: Centralized dashboard for businesses to manage their profiles, insights, and subscriptions.
- **Insights Marketplace**: Browse, purchase, and provide feedback on business insights from various industries.

## 2. Tools and Technologies  

### **1. Blockchain**
- **[Celo](https://celo.org/)**: For building the blockchain infrastructure and issuing platform tokens.  
- **[Solidity](https://soliditylang.org/)**: For developing smart contracts to handle tokenization and rewards.  

### **2. Backend**
- **[Go](https://go.dev/)**: For high-performance API development, ensuring efficient data aggregation and analysis.  

### **3. Frontend**
- **HTML/CSS/JavaScript**: For building an interactive and responsive user interface.
- **Modern UI Components**: Responsive design with grid layouts and interactive elements.
- **Client-side Authentication**: JWT-based authentication for secure user sessions.

# Setup Instructions  

### **1. Prerequisites**  
Ensure you have the following installed:  
- Go (1.19 or later)  
- Solidity compiler (via Remix or Hardhat)  
- Celo Wallet and CLI  
- Node.js and npm/yarn  

### **2. Clone the Repository**  
```bash
git clone https://github.com/ombima56/insights-edge.git
cd insights-edge
```

## 3. Backend Setup (Go)

### 1. Navigate to the backend folder:
```sh
cd backend
```

### 2. Install dependencies:
```sh
go mod tidy
```

### 3. Run the server:
```sh
go run main.go
```

## 4. Smart Contracts

### 1. Navigate to the contracts directory:
```sh
cd contracts
```

### 2. Compile and deploy contracts using Hardhat:
```sh
npx hardhat compile
npx hardhat run scripts/deploy.js --network celo
```

## 5. Frontend Features

### Dashboard
The dashboard provides businesses with a comprehensive overview of their account:

- **Profile Management**: Update business information and settings
- **Insights Management**: Create, edit, and monitor published insights
- **Subscription Management**: View and manage subscription plans
- **Transaction History**: Track all platform transactions
- **Analytics**: View performance metrics for published insights

### Marketplace
The marketplace allows users to discover and purchase business insights:

- **Search & Filter**: Find insights by industry, type, and other criteria
- **Insight Preview**: View summaries before purchasing
- **Purchase System**: Securely buy insights using platform tokens
- **Feedback System**: Rate and comment on purchased insights
- **Recommendations**: Discover relevant insights based on interests

## Usage
### 1. Access the Platform:
Open your browser and navigate to http://localhost:3000.

### 2. Connect Wallet:
Use the Celo Wallet to connect and interact with the platform.

### 3. Explore Insights:
-    Browse data-driven insights in the marketplace.
-    Earn tokens by contributing anonymized data or engaging with content.

### 4. Business Dashboard:
-    Register your business and create a profile.
-    Publish insights and track their performance.
-    Manage subscriptions and view transaction history.

## Roadmap
### Phase 1: MVP Development

-    Build core backend services for data aggregation.
-    Develop and deploy tokenization smart contracts.
-    Create a basic frontend interface for users.

### Phase 2: Tokenized Rewards System

-    Implement smart contracts for reward distribution.
-    Integrate token balances and redemption features.

### Phase 3: Advanced Analytics

-    Add machine learning for predictive insights.
-    Expand data sources and visualization tools.

### Phase 4: Scaling and Refinement

-    Deploy on mainnet.
-    Optimize for scalability and performance.

## License

This project is licensed under [MIT](https://github.com/Adamur-Tribe/insights-edge/blob/main/LICENSE).

## Authors:
[John E. Odhiambo](https://github.com/johneliud)

[Vinolia Esao](https://github.com/Vinolia-E)

[Hillary Ombima](https://github.com/ombima56)
