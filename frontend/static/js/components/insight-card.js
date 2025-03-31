import { shortenAddress } from '../modules/ui.js';

export const createInsightCard = (insight) => {
    const card = document.createElement('div');
    card.className = 'insight-card';
    
    card.innerHTML = `
        <div class="insight-header">
            <h3>${insight.title}</h3>
            <span class="industry-tag">${insight.industry}</span>
        </div>
        <div class="insight-content">
            <p>${insight.description}</p>
            <div class="insight-meta">
                <div class="price">Price: ${insight.price} ETH</div>
                <div class="provider">Provider: ${shortenAddress(insight.provider)}</div>
                <div class="timestamp">${new Date(insight.timestamp * 1000).toLocaleDateString()}</div>
            </div>
        </div>
        <div class="insight-actions">
            <button onclick="handlePurchase(${insight.id})" class="purchase-btn">Purchase</button>
            <button onclick="openInsightModal(${insight.id})" class="details-btn">Details</button>
        </div>
    `;
    
    return card;
};
