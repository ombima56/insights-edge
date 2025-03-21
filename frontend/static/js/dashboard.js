document.addEventListener('DOMContentLoaded', function() {
  // Check if user is logged in
  const token = localStorage.getItem('token');
  if (!token) {
    window.location.href = '/login';
    return;
  }

  // Initialize dashboard
  initializeDashboard();
  
  // Event listeners
  setupEventListeners();
});

// Initialize dashboard components
async function initializeDashboard() {
  // Load user profile
  await loadUserProfile();
  
  // Load dashboard data based on user type
  const userType = document.getElementById('user-type').textContent;
  
  if (userType === 'Business') {
    // Show business-specific sections
    document.getElementById('business-menu-item').style.display = 'flex';
    document.getElementById('insights-menu-item').style.display = 'flex';
    
    // Load business profile
    await loadBusinessProfile();
    
    // Load business insights
    await loadBusinessInsights();
    
    // Load business subscription
    await loadBusinessSubscription();
  } else {
    // Hide business-specific sections
    document.getElementById('business-menu-item').style.display = 'none';
    document.getElementById('insights-menu-item').style.display = 'none';
    
    // Load purchased insights
    await loadPurchasedInsights();
  }
  
  // Load transactions
  await loadTransactions();
  
  // Load dashboard stats
  await loadDashboardStats();
}

// Setup event listeners
function setupEventListeners() {
  // Sidebar menu navigation
  const menuItems = document.querySelectorAll('.sidebar-menu li');
  menuItems.forEach(item => {
    item.addEventListener('click', function() {
      // Remove active class from all menu items
      menuItems.forEach(i => i.classList.remove('active'));
      
      // Add active class to clicked menu item
      this.classList.add('active');
      
      // Hide all sections
      const sections = document.querySelectorAll('.dashboard-section');
      sections.forEach(section => section.classList.remove('active'));
      
      // Show selected section
      const sectionId = this.getAttribute('data-section');
      document.getElementById(sectionId).classList.add('active');
    });
  });
  
  // Logout button
  document.getElementById('logout-btn').addEventListener('click', function() {
    localStorage.removeItem('token');
    window.location.href = '/login';
  });
  
  // Register business form
  const registerBusinessForm = document.getElementById('register-business-form');
  if (registerBusinessForm) {
    registerBusinessForm.addEventListener('submit', async function(e) {
      e.preventDefault();
      await registerBusiness();
    });
  }
  
  // Create insight button and modal
  const createInsightBtn = document.getElementById('create-insight-btn');
  const createInsightModal = document.getElementById('create-insight-modal');
  const modalClose = document.querySelector('.modal-close');
  
  if (createInsightBtn) {
    createInsightBtn.addEventListener('click', function() {
      createInsightModal.style.display = 'flex';
    });
  }
  
  if (modalClose) {
    modalClose.addEventListener('click', function() {
      createInsightModal.style.display = 'none';
    });
  }
  
  // Create insight form
  const createInsightForm = document.getElementById('create-insight-form');
  if (createInsightForm) {
    createInsightForm.addEventListener('submit', async function(e) {
      e.preventDefault();
      await createInsight();
    });
  }
  
  // Purchase subscription buttons
  const purchasePlanBtns = document.querySelectorAll('.purchase-plan');
  purchasePlanBtns.forEach(btn => {
    btn.addEventListener('click', async function() {
      const plan = this.getAttribute('data-plan');
      await purchaseSubscription(plan);
    });
  });
  
  // Settings forms
  const profileSettingsForm = document.getElementById('profile-settings-form');
  if (profileSettingsForm) {
    profileSettingsForm.addEventListener('submit', async function(e) {
      e.preventDefault();
      await updateProfileSettings();
    });
  }
  
  const passwordSettingsForm = document.getElementById('password-settings-form');
  if (passwordSettingsForm) {
    passwordSettingsForm.addEventListener('submit', async function(e) {
      e.preventDefault();
      await updatePasswordSettings();
    });
  }
}

// Load user profile
async function loadUserProfile() {
  try {
    const response = await fetch('/api/user/profile', {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    
    if (!response.ok) {
      throw new Error('Failed to load user profile');
    }
    
    const userData = await response.json();
    
    // Update user profile in sidebar
    document.getElementById('username').textContent = userData.username;
    document.getElementById('user-type').textContent = userData.is_business ? 'Business' : 'Individual';
    
    // Update settings form
    document.getElementById('settings-username').value = userData.username;
    document.getElementById('settings-email').value = userData.email;
    document.getElementById('settings-wallet').value = userData.wallet_address;
    
    return userData;
  } catch (error) {
    console.error('Error loading user profile:', error);
    showNotification('Failed to load user profile', 'error');
  }
}

// Load business profile
async function loadBusinessProfile() {
  try {
    const response = await fetch('/api/business/profile', {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    
    if (response.status === 404) {
      // Business not registered yet
      document.getElementById('business-profile-content').style.display = 'none';
      document.getElementById('business-register-form').style.display = 'block';
      return null;
    }
    
    if (!response.ok) {
      throw new Error('Failed to load business profile');
    }
    
    const businessData = await response.json();
    
    // Update business profile content
    document.getElementById('business-profile-content').innerHTML = `
      <div class="business-detail">
        <span class="business-detail-label">Business Name</span>
        <div class="business-detail-value">${businessData.name}</div>
      </div>
      <div class="business-detail">
        <span class="business-detail-label">Industry</span>
        <div class="business-detail-value">${businessData.industry}</div>
      </div>
      <div class="business-detail">
        <span class="business-detail-label">Location</span>
        <div class="business-detail-value">${businessData.location || 'Not specified'}</div>
      </div>
      <div class="business-detail">
        <span class="business-detail-label">Subscription Status</span>
        <div class="business-detail-value">
          <span class="subscription-badge ${businessData.subscription_active ? 'active' : 'expired'}">
            ${businessData.subscription_active ? 'Active' : 'Expired'}
          </span>
        </div>
      </div>
    `;
    
    document.getElementById('business-profile-content').style.display = 'block';
    document.getElementById('business-register-form').style.display = 'none';
    
    return businessData;
  } catch (error) {
    console.error('Error loading business profile:', error);
    showNotification('Failed to load business profile', 'error');
  }
}

// Register business
async function registerBusiness() {
  try {
    const businessName = document.getElementById('business-name').value;
    const industry = document.getElementById('business-industry').value;
    const location = document.getElementById('business-location').value;
    
    const response = await fetch('/api/business/register', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({
        name: businessName,
        industry: industry,
        location: location
      })
    });
    
    if (!response.ok) {
      throw new Error('Failed to register business');
    }
    
    showNotification('Business registered successfully', 'success');
    
    // Reload business profile
    await loadBusinessProfile();
  } catch (error) {
    console.error('Error registering business:', error);
    showNotification('Failed to register business', 'error');
  }
}

// Load business insights
async function loadBusinessInsights() {
  try {
    const insightsList = document.getElementById('insights-list');
    insightsList.innerHTML = '<div class="loading-spinner">Loading insights...</div>';
    
    const response = await fetch('/api/insights?business=true', {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    
    if (!response.ok) {
      throw new Error('Failed to load insights');
    }
    
    const insights = await response.json();
    
    if (insights.length === 0) {
      insightsList.innerHTML = `
        <div class="insights-empty">
          <p>You haven't created any insights yet.</p>
          <p>Click the "Create New Insight" button to get started.</p>
        </div>
      `;
      return;
    }
    
    let insightsHTML = '';
    
    insights.forEach(insight => {
      insightsHTML += `
        <div class="insight-card" data-id="${insight.id}">
          <div class="insight-type">${insight.type}</div>
          <div class="insight-title">${insight.title}</div>
          <div class="insight-preview">${insight.content.substring(0, 150)}...</div>
          <div class="insight-meta">
            <div class="insight-rating">
              <i class="fas fa-star"></i>
              <span>${insight.average_rating || 'No ratings'}</span>
            </div>
            <div class="insight-date">${new Date(insight.created_at).toLocaleDateString()}</div>
          </div>
        </div>
      `;
    });
    
    insightsList.innerHTML = insightsHTML;
    
    // Add click event to insight cards
    const insightCards = document.querySelectorAll('.insight-card');
    insightCards.forEach(card => {
      card.addEventListener('click', function() {
        const insightId = this.getAttribute('data-id');
        viewInsightDetails(insightId);
      });
    });
  } catch (error) {
    console.error('Error loading insights:', error);
    document.getElementById('insights-list').innerHTML = `
      <div class="insights-empty">
        <p>Failed to load insights. Please try again later.</p>
      </div>
    `;
  }
}

// Create new insight
async function createInsight() {
  try {
    const industry = document.getElementById('insight-industry').value;
    const type = document.getElementById('insight-type').value;
    const content = document.getElementById('insight-data').value;
    
    const response = await fetch('/api/insights/list', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({
        industry: industry,
        type: type,
        content: content
      })
    });
    
    if (!response.ok) {
      throw new Error('Failed to create insight');
    }
    
    showNotification('Insight created successfully', 'success');
    
    // Close modal and reset form
    document.getElementById('create-insight-modal').style.display = 'none';
    document.getElementById('create-insight-form').reset();
    
    // Reload insights
    await loadBusinessInsights();
    
    // Reload dashboard stats
    await loadDashboardStats();
  } catch (error) {
    console.error('Error creating insight:', error);
    showNotification('Failed to create insight', 'error');
  }
}

// Load purchased insights
async function loadPurchasedInsights() {
  try {
    const insightsList = document.getElementById('insights-list');
    insightsList.innerHTML = '<div class="loading-spinner">Loading insights...</div>';
    
    const response = await fetch('/api/insights/purchased', {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    
    if (!response.ok) {
      throw new Error('Failed to load purchased insights');
    }
    
    const insights = await response.json();
    
    if (insights.length === 0) {
      insightsList.innerHTML = `
        <div class="insights-empty">
          <p>You haven't purchased any insights yet.</p>
          <p>Visit the marketplace to discover and purchase insights.</p>
        </div>
      `;
      return;
    }
    
    let insightsHTML = '';
    
    insights.forEach(insight => {
      insightsHTML += `
        <div class="insight-card" data-id="${insight.id}">
          <div class="insight-type">${insight.type}</div>
          <div class="insight-title">${insight.title}</div>
          <div class="insight-preview">${insight.content.substring(0, 150)}...</div>
          <div class="insight-meta">
            <div class="insight-rating">
              <i class="fas fa-star"></i>
              <span>${insight.average_rating || 'No ratings'}</span>
            </div>
            <div class="insight-date">${new Date(insight.created_at).toLocaleDateString()}</div>
          </div>
        </div>
      `;
    });
    
    insightsList.innerHTML = insightsHTML;
    
    // Add click event to insight cards
    const insightCards = document.querySelectorAll('.insight-card');
    insightCards.forEach(card => {
      card.addEventListener('click', function() {
        const insightId = this.getAttribute('data-id');
        viewInsightDetails(insightId);
      });
    });
  } catch (error) {
    console.error('Error loading purchased insights:', error);
    document.getElementById('insights-list').innerHTML = `
      <div class="insights-empty">
        <p>Failed to load insights. Please try again later.</p>
      </div>
    `;
  }
}

// Load business subscription
async function loadBusinessSubscription() {
  try {
    const subscriptionContainer = document.getElementById('current-subscription');
    subscriptionContainer.innerHTML = '<div class="loading-spinner">Loading subscription...</div>';
    
    const response = await fetch('/api/business/subscription', {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    
    if (response.status === 404) {
      // No subscription
      subscriptionContainer.innerHTML = `
        <div class="current-plan">
          <div class="current-plan-header">
            <div class="current-plan-name">No Active Subscription</div>
          </div>
          <p>You don't have an active subscription. Choose a plan below to subscribe.</p>
        </div>
      `;
      return;
    }
    
    if (!response.ok) {
      throw new Error('Failed to load subscription');
    }
    
    const subscription = await response.json();
    
    const expiryDate = new Date(subscription.expiry_date);
    const isActive = expiryDate > new Date();
    
    subscriptionContainer.innerHTML = `
      <div class="current-plan">
        <div class="current-plan-header">
          <div class="current-plan-name">${subscription.plan_name} Plan</div>
          <div class="current-plan-status ${isActive ? 'active' : 'expired'}">
            ${isActive ? 'Active' : 'Expired'}
          </div>
        </div>
        <div class="current-plan-info">
          <div class="plan-info-item">
            <div class="plan-info-label">Subscription ID</div>
            <div class="plan-info-value">${subscription.id}</div>
          </div>
          <div class="plan-info-item">
            <div class="plan-info-label">Start Date</div>
            <div class="plan-info-value">${new Date(subscription.start_date).toLocaleDateString()}</div>
          </div>
          <div class="plan-info-item">
            <div class="plan-info-label">Expiry Date</div>
            <div class="plan-info-value">${new Date(subscription.expiry_date).toLocaleDateString()}</div>
          </div>
          <div class="plan-info-item">
            <div class="plan-info-label">Amount Paid</div>
            <div class="plan-info-value">${subscription.amount} BPT</div>
          </div>
        </div>
      </div>
    `;
  } catch (error) {
    console.error('Error loading subscription:', error);
    document.getElementById('current-subscription').innerHTML = `
      <div class="current-plan">
        <div class="current-plan-header">
          <div class="current-plan-name">Error</div>
        </div>
        <p>Failed to load subscription information. Please try again later.</p>
      </div>
    `;
  }
}

// Purchase subscription
async function purchaseSubscription(plan) {
  try {
    let amount = 0;
    switch (plan) {
      case 'basic':
        amount = 50;
        break;
      case 'premium':
        amount = 100;
        break;
      case 'enterprise':
        amount = 200;
        break;
      default:
        throw new Error('Invalid plan');
    }
    
    const response = await fetch('/api/business/subscription/purchase', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({
        plan: plan,
        amount: amount
      })
    });
    
    if (!response.ok) {
      throw new Error('Failed to purchase subscription');
    }
    
    showNotification(`Successfully subscribed to ${plan} plan`, 'success');
    
    // Reload subscription
    await loadBusinessSubscription();
    
    // Reload transactions
    await loadTransactions();
    
    // Reload dashboard stats
    await loadDashboardStats();
  } catch (error) {
    console.error('Error purchasing subscription:', error);
    showNotification('Failed to purchase subscription', 'error');
  }
}

// Load transactions
async function loadTransactions() {
  try {
    const transactionsList = document.getElementById('transactions-list');
    transactionsList.innerHTML = '<div class="loading-spinner">Loading transactions...</div>';
    
    const response = await fetch('/api/transactions', {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    
    if (!response.ok) {
      throw new Error('Failed to load transactions');
    }
    
    const transactions = await response.json();
    
    if (transactions.length === 0) {
      transactionsList.innerHTML = `
        <div class="transactions-empty">
          <p>No transactions found.</p>
        </div>
      `;
      return;
    }
    
    let transactionsHTML = `
      <table>
        <thead>
          <tr>
            <th>Transaction ID</th>
            <th>Type</th>
            <th>Amount</th>
            <th>Date</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
    `;
    
    transactions.forEach(transaction => {
      transactionsHTML += `
        <tr>
          <td>${transaction.id}</td>
          <td>${transaction.type}</td>
          <td>${transaction.amount} BPT</td>
          <td>${new Date(transaction.timestamp).toLocaleString()}</td>
          <td>
            <span class="transaction-status ${transaction.status.toLowerCase()}">
              ${transaction.status}
            </span>
          </td>
        </tr>
      `;
    });
    
    transactionsHTML += `
        </tbody>
      </table>
    `;
    
    transactionsList.innerHTML = transactionsHTML;
  } catch (error) {
    console.error('Error loading transactions:', error);
    document.getElementById('transactions-list').innerHTML = `
      <div class="transactions-empty">
        <p>Failed to load transactions. Please try again later.</p>
      </div>
    `;
  }
}

// Load dashboard stats
async function loadDashboardStats() {
  try {
    const response = await fetch('/api/dashboard/stats', {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    
    if (!response.ok) {
      throw new Error('Failed to load dashboard stats');
    }
    
    const stats = await response.json();
    
    // Update stats
    document.getElementById('insights-count').textContent = stats.insights_count || 0;
    document.getElementById('purchases-count').textContent = stats.purchases_count || 0;
    document.getElementById('avg-rating').textContent = stats.average_rating || '0.0';
    document.getElementById('wallet-balance').textContent = `${stats.wallet_balance || 0} BPT`;
    
    // Update recent activity
    const activityList = document.getElementById('activity-list');
    
    if (!stats.recent_activity || stats.recent_activity.length === 0) {
      activityList.innerHTML = '<div class="activity-empty">No recent activity</div>';
      return;
    }
    
    let activityHTML = '';
    
    stats.recent_activity.forEach(activity => {
      let icon = 'fa-info-circle';
      
      switch (activity.type) {
        case 'insight_listed':
          icon = 'fa-lightbulb';
          break;
        case 'insight_purchased':
          icon = 'fa-shopping-cart';
          break;
        case 'subscription_purchased':
          icon = 'fa-credit-card';
          break;
        case 'feedback_submitted':
          icon = 'fa-comment';
          break;
      }
      
      activityHTML += `
        <div class="activity-item">
          <div class="activity-icon">
            <i class="fas ${icon}"></i>
          </div>
          <div class="activity-info">
            <div class="activity-title">${activity.description}</div>
            <div class="activity-time">${new Date(activity.timestamp).toLocaleString()}</div>
          </div>
        </div>
      `;
    });
    
    activityList.innerHTML = activityHTML;
  } catch (error) {
    console.error('Error loading dashboard stats:', error);
  }
}

// Update profile settings
async function updateProfileSettings() {
  try {
    const username = document.getElementById('settings-username').value;
    
    const response = await fetch('/api/user/profile', {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({
        username: username
      })
    });
    
    if (!response.ok) {
      throw new Error('Failed to update profile');
    }
    
    showNotification('Profile updated successfully', 'success');
    
    // Reload user profile
    await loadUserProfile();
  } catch (error) {
    console.error('Error updating profile:', error);
    showNotification('Failed to update profile', 'error');
  }
}

// Update password settings
async function updatePasswordSettings() {
  try {
    const currentPassword = document.getElementById('current-password').value;
    const newPassword = document.getElementById('new-password').value;
    const confirmPassword = document.getElementById('confirm-password').value;
    
    if (newPassword !== confirmPassword) {
      showNotification('New passwords do not match', 'error');
      return;
    }
    
    const response = await fetch('/api/user/password', {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({
        current_password: currentPassword,
        new_password: newPassword
      })
    });
    
    if (!response.ok) {
      throw new Error('Failed to update password');
    }
    
    showNotification('Password updated successfully', 'success');
    
    // Reset form
    document.getElementById('password-settings-form').reset();
  } catch (error) {
    console.error('Error updating password:', error);
    showNotification('Failed to update password', 'error');
  }
}

// View insight details
function viewInsightDetails(insightId) {
  // This would open a modal with insight details
  // For simplicity, we'll just redirect to the insight page
  window.location.href = `/insight/${insightId}`;
}

// Show notification
function showNotification(message, type = 'info') {
  // Create notification element
  const notification = document.createElement('div');
  notification.className = `notification ${type}`;
  notification.innerHTML = `
    <div class="notification-content">
      <i class="fas ${type === 'success' ? 'fa-check-circle' : type === 'error' ? 'fa-exclamation-circle' : 'fa-info-circle'}"></i>
      <span>${message}</span>
    </div>
    <button class="notification-close">&times;</button>
  `;
  
  // Add to document
  document.body.appendChild(notification);
  
  // Show notification
  setTimeout(() => {
    notification.classList.add('show');
  }, 10);
  
  // Auto hide after 5 seconds
  setTimeout(() => {
    notification.classList.remove('show');
    setTimeout(() => {
      notification.remove();
    }, 300);
  }, 5000);
  
  // Close button
  const closeBtn = notification.querySelector('.notification-close');
  closeBtn.addEventListener('click', () => {
    notification.classList.remove('show');
    setTimeout(() => {
      notification.remove();
    }, 300);
  });
}
