import { CONTRACT_ABI, loadInsights, loadMyInsights, createInsight, purchaseInsight } from './modules/contract.js';
import { updateUI, showToast, changeTab, updateIndustryFilterOptions, changePage } from './modules/ui.js';
import { createInsightCard } from './components/insight-card.js';

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
    purchasedIds: new Set(),
    lastLoadTime: 0
};

// Initialize Application
const initApp = async () => {
    try {
        if (window.ethereum) {
            window.web3 = new Web3(window.ethereum);
            await window.ethereum.enable();
            
            // Initialize contract
            state.contract = new window.web3.eth.Contract(
                CONTRACT_ABI,
                CONTRACT_ADDRESS
            );
            
            // Check if already connected
            const accounts = await window.web3.eth.getAccounts();
            if (accounts.length > 0) {
                state.account = accounts[0];
                state.isConnected = true;
                updateUI();
            }
            
            // Load initial data
            await loadInsights();
            await loadMyInsights();
            
            // Update UI
            updateIndustryFilterOptions();
            renderInsights();
            renderMyInsights();
            
            // Subscribe to events
            subscribeToEvents();
        } else {
            showToast('Please install MetaMask to use this application', 'error');
        }
    } catch (error) {
        console.error('Error initializing app:', error);
        showToast('Error initializing application', 'error');
    }
};

// Setup Event Listeners
const setupEventListeners = () => {
    connectWalletBtn.addEventListener('click', connectWallet);
    industryFilter.addEventListener('change', updateFilters);
    searchInput.addEventListener('input', updateFilters);
    prevPageBtn.addEventListener('click', () => changePage('prev'));
    nextPageBtn.addEventListener('click', () => changePage('next'));
    createInsightBtn.addEventListener('click', () => createInsightModal.style.display = 'block');
    closeModalBtn.addEventListener('click', () => insightModal.style.display = 'none');
    createInsightForm.addEventListener('submit', handleCreateInsight);
};

// Connect Wallet
const connectWallet = async () => {
    try {
        if (state.isConnected) {
            await disconnectWallet();
            return;
        }
        
        const accounts = await window.ethereum.request({ method: 'eth_requestAccounts' });
        state.account = accounts[0];
        state.isConnected = true;
        
        // Load user data
        await loadMyInsights();
        
        updateUI();
        showToast('Wallet connected successfully');
    } catch (error) {
        console.error('Error connecting wallet:', error);
        showToast('Error connecting wallet', 'error');
    }
};

// Disconnect Wallet
const disconnectWallet = () => {
    state.account = null;
    state.isConnected = false;
    state.myCreatedInsights = [];
    state.myPurchasedInsights = [];
    updateUI();
    showToast('Wallet disconnected');
};

// Subscribe to Contract Events
const subscribeToEvents = () => {
    state.contract.events.InsightCreated({
        fromBlock: 'latest'
    }).on('data', async (event) => {
        await loadInsights();
        updateIndustryFilterOptions();
        renderInsights();
        showToast('New insight created');
    });

    state.contract.events.InsightPurchased({
        fromBlock: 'latest'
    }).on('data', async (event) => {
        await loadInsights();
        await loadMyInsights();
        renderInsights();
        renderMyInsights();
        showToast('Insight purchased');
    });
};

// Render Functions
const renderInsights = () => {
    const filteredInsights = state.insights
        .filter(insight => {
            if (state.currentFilter && insight.industry !== state.currentFilter) return false;
            if (state.searchQuery && !insight.title.toLowerCase().includes(state.searchQuery.toLowerCase())) return false;
            return true;
        });
    
    const startIndex = (state.currentPage - 1) * state.itemsPerPage;
    const endIndex = startIndex + state.itemsPerPage;
    const pageInsights = filteredInsights.slice(startIndex, endIndex);
    
    insightsContainer.innerHTML = '';
    pageInsights.forEach(insight => {
        const card = createInsightCard(insight);
        insightsContainer.appendChild(card);
    });
    
    // Update pagination
    const totalPages = Math.ceil(filteredInsights.length / state.itemsPerPage);
    pageIndicator.textContent = `${state.currentPage} of ${totalPages}`;
    prevPageBtn.disabled = state.currentPage === 1;
    nextPageBtn.disabled = state.currentPage === totalPages;
};

const renderMyInsights = () => {
    myInsightsContainer.innerHTML = '';
    
    // Show created insights
    const createdSection = document.createElement('div');
    createdSection.className = 'my-insights-section';
    createdSection.innerHTML = '<h3>Created Insights</h3>';
    state.myCreatedInsights.forEach(insight => {
        const card = createInsightCard(insight);
        createdSection.appendChild(card);
    });
    myInsightsContainer.appendChild(createdSection);
    
    // Show purchased insights
    const purchasedSection = document.createElement('div');
    purchasedSection.className = 'my-insights-section';
    purchasedSection.innerHTML = '<h3>Purchased Insights</h3>';
    state.myPurchasedInsights.forEach(insight => {
        const card = createInsightCard(insight);
        purchasedSection.appendChild(card);
    });
    myInsightsContainer.appendChild(purchasedSection);
};

const openInsightModal = (insight) => {
    state.selectedInsight = insight;
    insightModalContent.innerHTML = `
        <h2>${insight.title}</h2>
        <p><strong>Industry:</strong> ${insight.industry}</p>
        <p><strong>Description:</strong> ${insight.description}</p>
        <p><strong>Price:</strong> ${insight.price} ETH</p>
        <p><strong>Provider:</strong> ${shortenAddress(insight.provider)}</p>
        <p><strong>Created:</strong> ${new Date(insight.timestamp * 1000).toLocaleString()}</p>
        <button onclick="handlePurchase(${insight.id})" class="purchase-btn">Purchase</button>
    `;
    insightModal.style.display = 'block';
};

const handleCreateInsight = async (e) => {
    e.preventDefault();
    try {
        const industry = document.getElementById('industry').value;
        const title = document.getElementById('title').value;
        const description = document.getElementById('description').value;
        const price = document.getElementById('price').value;
        
        await createInsight(industry, title, description, price);
        createInsightModal.style.display = 'none';
        createInsightForm.reset();
        
        await loadInsights();
        updateIndustryFilterOptions();
        renderInsights();
        
        showToast('Insight created successfully');
    } catch (error) {
        console.error('Error creating insight:', error);
        showToast('Error creating insight', 'error');
    }
};

const handlePurchaseInsight = async (insightId) => {
    try {
        await purchaseInsight(insightId);
        await loadInsights();
        await loadMyInsights();
        renderInsights();
        renderMyInsights();
        
        showToast('Insight purchased successfully');
    } catch (error) {
        console.error('Error purchasing insight:', error);
        showToast('Error purchasing insight', 'error');
    }
};

const updateFilters = () => {
    state.currentFilter = industryFilter.value;
    state.searchQuery = searchInput.value;
    state.currentPage = 1;
    renderInsights();
};

// Initialize the application
window.addEventListener('DOMContentLoaded', () => {
    initApp();
    setupEventListeners();
});
