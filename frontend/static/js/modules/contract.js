// Contract configuration
const CONTRACT_ADDRESS = '0xf1620709eb1818ea618465643685c34d40e75c5a';

// Contract ABI
export const CONTRACT_ABI = [
    {
        "anonymous": false,
        "inputs": [
            {
                "indexed": true,
                "internalType": "uint256",
                "name": "id",
                "type": "uint256"
            },
            {
                "indexed": true,
                "internalType": "address",
                "name": "provider",
                "type": "address"
            },
            {
                "indexed": false,
                "internalType": "string",
                "name": "industry",
                "type": "string"
            },
            {
                "indexed": false,
                "internalType": "string",
                "name": "title",
                "type": "string"
            },
            {
                "indexed": false,
                "internalType": "uint256",
                "name": "price",
                "type": "uint256"
            }
        ],
        "name": "InsightCreated",
        "type": "event"
    },
    {
        "anonymous": false,
        "inputs": [
            {
                "indexed": true,
                "internalType": "uint256",
                "name": "id",
                "type": "uint256"
            },
            {
                "indexed": true,
                "internalType": "address",
                "name": "buyer",
                "type": "address"
            },
            {
                "indexed": true,
                "internalType": "address",
                "name": "seller",
                "type": "address"
            },
            {
                "indexed": false,
                "internalType": "uint256",
                "name": "price",
                "type": "uint256"
            }
        ],
        "name": "InsightPurchased",
        "type": "event"
    },
    {
        "inputs": [
            {
                "internalType": "string",
                "name": "_industry",
                "type": "string"
            },
            {
                "internalType": "string",
                "name": "_title",
                "type": "string"
            },
            {
                "internalType": "string",
                "name": "_description",
                "type": "string"
            },
            {
                "internalType": "uint256",
                "name": "_price",
                "type": "uint256"
            }
        ],
        "name": "createInsight",
        "outputs": [],
        "stateMutability": "nonpayable",
        "type": "function"
    },
    {
        "inputs": [
            {
                "internalType": "uint256",
                "name": "_insightId",
                "type": "uint256"
            }
        ],
        "name": "getInsight",
        "outputs": [
            {
                "internalType": "address",
                "name": "provider",
                "type": "address"
            },
            {
                "internalType": "string",
                "name": "industry",
                "type": "string"
            },
            {
                "internalType": "string",
                "name": "title",
                "type": "string"
            },
            {
                "internalType": "string",
                "name": "description",
                "type": "string"
            },
            {
                "internalType": "uint256",
                "name": "price",
                "type": "uint256"
            },
            {
                "internalType": "uint256",
                "name": "timestamp",
                "type": "uint256"
            }
        ],
        "stateMutability": "view",
        "type": "function"
    },
    {
        "inputs": [],
        "name": "insightCount",
        "outputs": [
            {
                "internalType": "uint256",
                "name": "",
                "type": "uint256"
            }
        ],
        "stateMutability": "view",
        "type": "function"
    }
];

export const loadInsights = async () => {
    try {
        const insightCount = await state.contract.methods.insightCount().call();
        const insights = [];
        
        for (let i = 1; i <= insightCount; i++) {
            const insight = await state.contract.methods.getInsight(i).call();
            insights.push({
                id: i,
                provider: insight[0],
                industry: insight[1],
                title: insight[2],
                description: insight[3],
                price: insight[4],
                timestamp: insight[5]
            });
        }
        
        state.insights = insights;
        return insights;
    } catch (error) {
        console.error('Error loading insights:', error);
        throw error;
    }
};

export const loadMyInsights = async () => {
    try {
        const account = state.account;
        const insights = await loadInsights();
        
        state.myCreatedInsights = insights.filter(
            insight => insight.provider.toLowerCase() === account.toLowerCase()
        );
        
        // Load purchased insights
        const purchasedInsights = [];
        for (let i = 1; i <= state.insights.length; i++) {
            try {
                const purchased = await state.contract.methods
                    .isPurchased(i, account)
                    .call();
                if (purchased) {
                    purchasedInsights.push(state.insights[i-1]);
                }
            } catch (err) {
                console.error(`Error checking insight ${i}:`, err);
            }
        }
        
        state.myPurchasedInsights = purchasedInsights;
        return { created: state.myCreatedInsights, purchased: purchasedInsights };
    } catch (error) {
        console.error('Error loading my insights:', error);
        throw error;
    }
};

export const createInsight = async (industry, title, description, price) => {
    try {
        const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });
        const account = accounts[0];
        
        const gasPrice = await window.ethereum.request({
            method: 'eth_gasPrice'
        });
        
        const tx = await state.contract.methods.createInsight(
            industry,
            title,
            description,
            window.web3.utils.toWei(price.toString(), 'ether')
        ).send({
            from: account,
            gasPrice: gasPrice
        });
        
        return tx;
    } catch (error) {
        console.error('Error creating insight:', error);
        throw error;
    }
};

export const purchaseInsight = async (insightId) => {
    try {
        const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });
        const account = accounts[0];
        
        const insight = state.insights.find(i => i.id === insightId);
        if (!insight) {
            throw new Error('Insight not found');
        }
        
        const gasPrice = await window.ethereum.request({
            method: 'eth_gasPrice'
        });
        
        const tx = await state.contract.methods.purchaseInsight(insightId).send({
            from: account,
            value: window.web3.utils.toWei(insight.price.toString(), 'ether'),
            gasPrice: gasPrice
        });
        
        return tx;
    } catch (error) {
        console.error('Error purchasing insight:', error);
        throw error;
    }
};
