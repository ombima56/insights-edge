const connectWalletBtn = document.getElementById('connect-wallet');
const connectionStatus = document.getElementById('connection-status');
const insightsContainer = document.getElementById('insights-container');
const myInsightsContainer = document.getElementById('my-insights-container');
const insightModal = document.getElementById('insight-modal');
const insightModalContent = document.getElementById('insight-modal-content');
const closeModalBtn = document.getElementById('close-modal');
const createInsightBtn = document.getElementById('create-insight-btn');
const createInsightModal = document.getElementById('create-insight-modal');
const createInsightForm = document.getElementById('create-insight-form');
const industryFilter = document.getElementById('industry-filter');
const searchInput = document.getElementById('search-input');
const prevPageBtn = document.getElementById('prev-page');
const nextPageBtn = document.getElementById('next-page');
const pageIndicator = document.getElementById('page-indicator');

export const updateUI = () => {
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
};

export const showToast = (message, type = 'info') => {
    const toast = document.createElement('div');
    toast.className = `toast ${type}`;
    toast.textContent = message;
    
    document.body.appendChild(toast);
    
    setTimeout(() => {
        toast.remove();
    }, 3000);
};

export const changeTab = (tabName, activeBtn, btnSelector, contentSelector) => {
    // Hide all tabs
    document.querySelectorAll('.tab-content').forEach(tab => tab.style.display = 'none');
    document.querySelectorAll('.tab-btn').forEach(btn => btn.classList.remove('active'));
    
    // Show selected tab
    document.querySelector(contentSelector).style.display = 'block';
    document.querySelector(btnSelector).classList.add('active');
    
    // Update URL hash
    window.location.hash = tabName;
};

export const shortenAddress = (address) => {
    return `${address.substring(0, 6)}...${address.substring(address.length - 4)}`;
};

export const updateIndustryFilterOptions = () => {
    const uniqueIndustries = [...new Set(state.insights.map(i => i.industry))];
    industryFilter.innerHTML = '<option value="">All Industries</option>';
    uniqueIndustries.forEach(industry => {
        const option = document.createElement('option');
        option.value = industry;
        option.textContent = industry;
        industryFilter.appendChild(option);
    });
};

export const changePage = (direction) => {
    if (direction === 'prev' && state.currentPage > 1) {
        state.currentPage--;
    } else if (direction === 'next' && state.currentPage < Math.ceil(state.insights.length / state.itemsPerPage)) {
        state.currentPage++;
    }
    
    renderInsights();
};
