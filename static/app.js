// Application State
const state = {
    web3: null,
    contract: null,
    account: null,
    isConnected: false,
    insightCount: 0,
    insights: [],
    myCreatedInsights: [],
    myPurchasedInsights: [],
    currentPage: 1,
    itemsPerPage: 9,
    selectedInsight: null,
    industryFilters: new Set(),
    currentFilter: '',
    searchQuery: '',
    purchasedIds: new Set()
};

// Contract Address 
const CONTRACT_ADDRESS = '0xf1620709eb1818ea618465643685c34d40e75c5a'; 

// Contract ABI from provided JSON
const CONTRACT_ABI = [
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
    },
    {
        "inputs": [
            {
                "internalType": "uint256",
                "name": "",
                "type": "uint256"
            }
        ],
        "name": "insights",
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
        "inputs": [
            {
                "internalType": "uint256",
                "name": "_insightId",
                "type": "uint256"
            }
        ],
        "name": "purchaseInsight",
        "outputs": [],
        "stateMutability": "payable",
        "type": "function"
    }
];

// DOM Elements
const connectWalletBtn = document.getElementById('connect-wallet');
const connectionStatus = document.getElementById('connection-status');
const accountAddress = document.getElementById('account-address');
const insightsContainer = document.getElementById('insights-container');
const createInsightForm = document.getElementById('create-insight-form');
const industryFilter = document.getElementById('industry-filter');
const searchInput = document.getElementById('search-insights');
const prevPageBtn = document.getElementById('prev-page');
const nextPageBtn = document.getElementById('next-page');
const pageInfo = document.getElementById('page-info');
const createdInsightsContainer = document.getElementById('created-insights-container');
const purchasedInsightsContainer = document.getElementById('purchased-insights-container');
const modal = document.getElementById('insight-details-modal');
const closeModal = document.querySelector('.close');
const modalTitle = document.getElementById('modal-insight-title');
const modalIndustry = document.getElementById('modal-insight-industry');
const modalProvider = document.getElementById('modal-insight-provider');
const modalPrice = document.getElementById('modal-insight-price');
const modalTimestamp = document.getElementById('modal-insight-timestamp');
const modalDescription = document.getElementById('modal-insight-description');
const purchaseBtn = document.getElementById('purchase-insight-btn');
const toast = document.getElementById('toast');

// Initialize Application
async function initApp() {
    // Setup event listeners
    setupEventListeners();
    
    // Check if Web3 is already injected
    if (window.ethereum) {
        try {
            state.web3 = new Web3(window.ethereum);
            
            // Check if already connected
            const accounts = await state.web3.eth.getAccounts();
            if (accounts && accounts.length > 0) {
                await connectWallet();
            }
        } catch (error) {
            console.error('Error initializing Web3:', error);
            showToast('Failed to initialize Web3', 'error');
        }
    } else {
        showToast('Please install MetaMask to use this application', 'warning');
    }
}

// Setup Event Listeners
function setupEventListeners() {
    // Wallet connection
    connectWalletBtn.addEventListener('click', connectWallet);
    
    
    // Tab navigation
    const tabBtns = document.querySelectorAll('.tab-btn');
    tabBtns.forEach(btn => {
        btn.addEventListener('click', () => {
            const tabName = btn.getAttribute('data-tab');
            changeTab(tabName, btn, '.tab-btn', '.tab-content');
        });
    });
    
    // Sub-tab navigation
    const subTabBtns = document.querySelectorAll('.sub-tab-btn');
    subTabBtns.forEach(btn => {
        btn.addEventListener('click', () => {
            const subTabName = btn.getAttribute('data-subtab');
            changeTab(subTabName, btn, '.sub-tab-btn', '.sub-tab-content');
        });
    });
    
    // Form submission
    createInsightForm.addEventListener('submit', handleCreateInsight);
    
    // Filters
    industryFilter.addEventListener('change', updateFilters);
    searchInput.addEventListener('input', updateFilters);
    
    // Pagination
    prevPageBtn.addEventListener('click', () => changePage(-1));
    nextPageBtn.addEventListener('click', () => changePage(1));
    
    // Modal
    closeModal.addEventListener('click', () => modal.style.display = 'none');
    window.addEventListener('click', (e) => {
        if (e.target === modal) modal.style.display = 'none';
    });
    
    // Purchase button
    purchaseBtn.addEventListener('click', handlePurchaseInsight);
    
    // Metamask account change
    if (window.ethereum) {
        window.ethereum.on('accountsChanged', async (accounts) => {
            if (accounts.length === 0) {
                // User logged out
                disconnectWallet();
            } else {
                // Account changed
                state.account = accounts[0];
                updateUI();
                await loadInsights();
            }
        });
        
        // Network change
        window.ethereum.on('chainChanged', () => {
            window.location.reload();
        });
    }
}

async function connectWallet() {
    try {
        if (!state.web3) {
            if (window.ethereum) {
                state.web3 = new Web3(window.ethereum);
            } else {
                showToast('Please install MetaMask to use this application', 'warning');
                return;
            }
        }

        // Request account access
        const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });

        if (!accounts || accounts.length === 0) {
            showToast('No accounts found. Please unlock MetaMask.', 'error');
            return;
        }

        state.account = accounts[0];
        state.isConnected = true;

        // Initialize contract
        state.contract = new state.web3.eth.Contract(CONTRACT_ABI, CONTRACT_ADDRESS);

        // Update UI
        updateUI();
        await loadInsights();
        subscribeToEvents();
        
        showToast('Wallet connected successfully!', 'success');
    } catch (error) {
        console.error('Error connecting wallet:', error);
        showToast('Failed to connect wallet', 'error');
    }
}


// Disconnect Wallet
function disconnectWallet() {
    state.account = null;
    state.isConnected = false;
    state.contract = null;
    updateUI();
    showToast('Wallet disconnected', 'info');
}

// Update UI based on connection state
function updateUI() {
    if (state.isConnected) {
        connectionStatus.textContent = 'Connected';
        connectionStatus.classList.add('connected');
        connectWalletBtn.textContent = 'Disconnect';
        connectWalletBtn.removeEventListener('click', connectWallet);
        connectWalletBtn.addEventListener('click', disconnectWallet);
        accountAddress.textContent = shortenAddress(state.account);
    } else {
        connectionStatus.textContent = 'Not Connected';
        connectionStatus.classList.remove('connected');
        connectWalletBtn.textContent = 'Connect Wallet';
        connectWalletBtn.removeEventListener('click', disconnectWallet);
        connectWalletBtn.addEventListener('click', connectWallet);
        accountAddress.textContent = '';
        
        // Clear containers
        insightsContainer.innerHTML = '';
        createdInsightsContainer.innerHTML = '';
        purchasedInsightsContainer.innerHTML = '';
    }
}

// Subscribe to Contract Events
function subscribeToEvents() {
    if (!state.contract) return;
    
    // Listen for InsightCreated events
    state.contract.events.InsightCreated()
        .on('data', async (event) => {
            const { id, provider, industry, title, price } = event.returnValues;
            showToast(`New insight created: ${title}`, 'info');
            await loadInsights();
        })
        .on('error', console.error);
    
    // Listen for InsightPurchased events
    state.contract.events.InsightPurchased()
        .on('data', async (event) => {
            const { id, buyer, seller, price } = event.returnValues;
            if (buyer.toLowerCase() === state.account.toLowerCase()) {
                showToast('Purchase successful!', 'success');
                state.purchasedIds.add(parseInt(id));
                await loadMyInsights();
            }
        })
        .on('error', console.error);
}

// Load all insights from the contract
async function loadInsights() {
    if (!state.contract) return;
    
    try {
        // Get total insight count
        const count = await state.contract.methods.insightCount().call();
        state.insightCount = parseInt(count);
        
        // Clear existing insights
        state.insights = [];
        state.industryFilters.clear();
        
        // Load all insights
        for (let i = 1; i <= state.insightCount; i++) {
            const insight = await state.contract.methods.getInsight(i).call();
            
            // Add id to the insight object
            insight.id = i;
            
            // Add to insights array
            state.insights.push(insight);
            
            // Add industry to filters
            if (insight.industry && insight.industry.trim() !== '') {
                state.industryFilters.add(insight.industry);
            }
        }
        
        // Update industry filter dropdown
        updateIndustryFilterOptions();
        
        // Display insights
        renderInsights();
        
        // Load insights related to current user
        await loadMyInsights();
        
    } catch (error) {
        console.error('Error loading insights:', error);
        showToast('Failed to load insights', 'error');
    }
}

// Load insights created by or purchased by the current user
async function loadMyInsights() {
    if (!state.contract || !state.account) return;
    
    try {
        // Clear existing user insights
        state.myCreatedInsights = [];
        state.myPurchasedInsights = [];
        
        // Filter created insights
        state.myCreatedInsights = state.insights.filter(
            insight => insight.provider.toLowerCase() === state.account.toLowerCase()
        );
        
        // Load purchased insights
        // Note: This is simplified - in a real app, you might need a separate mapping in the contract
        // to track purchases by user
        const events = await state.contract.getPastEvents('InsightPurchased', {
            filter: { buyer: state.account },
            fromBlock: 0,
            toBlock: 'latest'
        });
        
        // Extract purchased insight IDs
        const purchasedIds = events.map(event => parseInt(event.returnValues.id));
        state.purchasedIds = new Set(purchasedIds);
        
        // Get full insight data for purchased insights
        for (const id of purchasedIds) {
            const insight = state.insights.find(i => parseInt(i.id) === id);
            if (insight) state.myPurchasedInsights.push(insight);
        }
        
        // Render user insights
        renderMyInsights();
        
    } catch (error) {
        console.error('Error loading user insights:', error);
        showToast('Failed to load your insights', 'error');
    }
}

// Update industry filter dropdown options
function updateIndustryFilterOptions() {
    // Clear existing options (except "All Industries")
    while (industryFilter.options.length > 1) {
        industryFilter.remove(1);
    }
    
    // Add industry options
    state.industryFilters.forEach(industry => {
        const option = document.createElement('option');
        option.value = industry;
        option.textContent = industry;
        industryFilter.appendChild(option);
    });
}

// Render all insights with pagination and filtering
function renderInsights() {
    insightsContainer.innerHTML = '';
    
    // Apply filters
    const filteredInsights = state.insights.filter(insight => {
        // Check industry filter
        const industryMatch = state.currentFilter === '' || 
                             insight.industry === state.currentFilter;
        
        // Check search query
        const searchMatch = state.searchQuery === '' || 
                           insight.title.toLowerCase().includes(state.searchQuery.toLowerCase());
        
        return industryMatch && searchMatch;
    });
    
    // Calculate pagination
    const startIndex = (state.currentPage - 1) * state.itemsPerPage;
    const endIndex = Math.min(startIndex + state.itemsPerPage, filteredInsights.length);
    const paginatedInsights = filteredInsights.slice(startIndex, endIndex);
    
    // Update pagination UI
    pageInfo.textContent = `Page ${state.currentPage} of ${Math.max(1, Math.ceil(filteredInsights.length / state.itemsPerPage))}`;
    prevPageBtn.disabled = state.currentPage === 1;
    nextPageBtn.disabled = endIndex >= filteredInsights.length;
    
    // No insights to display
    if (paginatedInsights.length === 0) {
        insightsContainer.innerHTML = '<p class="no-results">No insights found matching your filters</p>';
        return;
    }
    
    // Render each insight card
    paginatedInsights.forEach(insight => {
        const insightCard = createInsightCard(insight);
        insightsContainer.appendChild(insightCard);
    });
}

// Render user-specific insights
function renderMyInsights() {
    // Render created insights
    createdInsightsContainer.innerHTML = '';
    if (state.myCreatedInsights.length === 0) {
        createdInsightsContainer.innerHTML = '<p class="no-results">You haven\'t created any insights yet</p>';
    } else {
        state.myCreatedInsights.forEach(insight => {
            const insightCard = createInsightCard(insight);
            createdInsightsContainer.appendChild(insightCard);
        });
    }
    
    // Render purchased insights
    purchasedInsightsContainer.innerHTML = '';
    if (state.myPurchasedInsights.length === 0) {
        purchasedInsightsContainer.innerHTML = '<p class="no-results">You haven\'t purchased any insights yet</p>';
    } else {
        state.myPurchasedInsights.forEach(insight => {
            const insightCard = createInsightCard(insight);
            purchasedInsightsContainer.appendChild(insightCard);
        });
    }
}

// Create insight card element
function createInsightCard(insight) {
    const card = document.createElement('div');
    card.className = 'insight-card';
    
    // Format price from wei to ETH
    const priceInEth = state.web3.utils.fromWei(insight.price.toString(), 'ether');
    
    // Format timestamp to date
    const date = new Date(insight.timestamp * 1000);
    const formattedDate = date.toLocaleDateString();
    
    card.innerHTML = `
        <div class="insight-card-header">
            <span class="insight-card-industry">${insight.industry}</span>
        </div>
        <div class="insight-card-body">
            <h3 class="insight-card-title">${insight.title}</h3>
            <p class="insight-card-description">Provider: ${shortenAddress(insight.provider)} • Created: ${formattedDate}</p>
        </div>
        <div class="insight-card-footer">
            <div class="insight-card-price">${priceInEth} ETH</div>
            <button class="view-insight-btn btn btn-primary">View Details</button>
        </div>
    `;
    
    // Add click event to show details modal
    card.addEventListener('click', () => {
        state.selectedInsight = insight;
        openInsightModal(insight);
    });
    
    return card;
}

// Open insight details modal
function openInsightModal(insight) {
    modalTitle.textContent = insight.title;
    modalIndustry.textContent = insight.industry;
    modalProvider.textContent = shortenAddress(insight.provider);
    modalPrice.textContent = state.web3.utils.fromWei(insight.price.toString(), 'ether');
    
    // Format timestamp
    const date = new Date(insight.timestamp * 1000);
    modalTimestamp.textContent = date.toLocaleString();
    
    // Check if the user has already purchased the insight
    const hasAccess = state.purchasedIds.has(parseInt(insight.id)) || 
                       insight.provider.toLowerCase() === state.account.toLowerCase();
    
    // Show or hide description based on purchase status
    if (hasAccess) {
        modalDescription.textContent = insight.description;
        purchaseBtn.style.display = 'none';
    } else {
        modalDescription.textContent = 'Purchase this insight to view the complete description.';
        purchaseBtn.style.display = 'block';
    }
    
    // Show the modal
    modal.style.display = 'block';
}

// Handle creating a new insight
async function handleCreateInsight(e) {
    e.preventDefault();
    
    if (!state.isConnected) {
        showToast('Please connect your wallet first', 'warning');
        return;
    }
    
    // Get form values
    const industry = document.getElementById('industry').value;
    const title = document.getElementById('title').value;
    const description = document.getElementById('description').value;
    const price = document.getElementById('price').value;
    
    // Convert price to wei
    const priceInWei = state.web3.utils.toWei(price, 'ether');
    
    try {
        // Send transaction to create insight
        await state.contract.methods.createInsight(industry, title, description, priceInWei)
            .send({ from: state.account });
        
        // Reset form
        createInsightForm.reset();
        
        // Switch to browse tab
        const browseTabBtn = document.querySelector('[data-tab="browse"]');
        changeTab('browse', browseTabBtn, '.tab-btn', '.tab-content');
        
        showToast('Insight created successfully!', 'success');
        
    } catch (error) {
        console.error('Error creating insight:', error);
        showToast('Failed to create insight', 'error');
    }
}

// Handle purchasing an insight
async function handlePurchaseInsight() {
    if (!state.isConnected || !state.selectedInsight) {
        return;
    }
    
    try {
        // Get insight price
        const price = state.selectedInsight.price;
        const insightId = state.selectedInsight.id;
        
        // Send transaction to purchase insight
        await state.contract.methods.purchaseInsight(insightId)
            .send({ from: state.account, value: price });
        
        // Close modal and refresh
        modal.style.display = 'none';
        
    } catch (error) {
        console.error('Error purchasing insight:', error);
        showToast('Failed to purchase insight', 'error');
    }
}

// Update filters and render insights
function updateFilters() {
    state.currentFilter = industryFilter.value;
    state.searchQuery = searchInput.value;
    state.currentPage = 1; // Reset to first page
    renderInsights();
}

// Change page
function changePage(direction) {
    state.currentPage += direction;
    if (state.currentPage < 1) state.currentPage = 1;
    renderInsights();
}

// Change active tab
function changeTab(tabName, activeBtn, btnSelector, contentSelector) {
    // Update active button
    document.querySelectorAll(btnSelector).forEach(btn => {
        btn.classList.remove('active');
    });
    activeBtn.classList.add('active');
    
    // Show active content
    document.querySelectorAll(contentSelector).forEach(content => {
        content.classList.remove('active');
    });
    document.getElementById(tabName).classList.add('active');
}

function showToast(message, type = 'info') {
    if (!toast) {
        console.error('Toast element not found');
        return;
    }
    toast.textContent = message;
    toast.className = `toast ${type}`;
    toast.style.display = 'block';
    setTimeout(() => {
        toast.style.display = 'none';
    }, 3000);
}


// Helper to shorten Ethereum addresses
function shortenAddress(address) {
    return `${address.substring(0, 6)}...${address.substring(address.length - 4)}`;
}

// Initialize the application when the DOM is loaded
document.addEventListener('DOMContentLoaded', initApp);