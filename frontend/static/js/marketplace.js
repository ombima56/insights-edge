document.addEventListener('DOMContentLoaded', function() {
  // Check if user is logged in
  const token = localStorage.getItem('token');
  if (!token) {
    // window.location.href = '/dashboard';
    return;
  }

  // Initialize marketplace
  initializeMarketplace();
  
  // Event listeners
  setupEventListeners();
});

// Initialize marketplace
async function initializeMarketplace() {
  // Load insights with default filters
  await loadInsights();
}

// Setup event listeners
function setupEventListeners() {
  // Search button
  document.getElementById('search-insights').addEventListener('click', function() {
    loadInsights();
  });
  
  // Industry filter change
  document.getElementById('industry-filter').addEventListener('change', function() {
    loadInsights();
  });
  
  // Type filter change
  document.getElementById('type-filter').addEventListener('change', function() {
    loadInsights();
  });
  
  // Sort filter change
  document.getElementById('sort-filter').addEventListener('change', function() {
    loadInsights();
  });
  
  // Modal close button
  const modalClose = document.querySelector('.modal-close');
  modalClose.addEventListener('click', function() {
    document.getElementById('insight-detail-modal').style.display = 'none';
  });
  
  // Close modal when clicking outside
  window.addEventListener('click', function(event) {
    const modal = document.getElementById('insight-detail-modal');
    if (event.target === modal) {
      modal.style.display = 'none';
    }
  });
}

// Load insights with filters
async function loadInsights(page = 1) {
  try {
    const insightsGrid = document.getElementById('insights-grid');
    insightsGrid.innerHTML = '<div class="loading-spinner">Loading insights...</div>';
    
    // Get filter values
    const industry = document.getElementById('industry-filter').value;
    const type = document.getElementById('type-filter').value;
    const sort = document.getElementById('sort-filter').value;
    
    // Build query string
    let queryString = `?page=${page}`;
    if (industry) queryString += `&industry=${encodeURIComponent(industry)}`;
    if (type) queryString += `&type=${encodeURIComponent(type)}`;
    if (sort) queryString += `&sort=${encodeURIComponent(sort)}`;
    
    const response = await fetch(`/api/insights${queryString}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    
    if (!response.ok) {
      throw new Error('Failed to load insights');
    }
    
    const data = await response.json();
    const insights = data.insights || [];
    const totalPages = data.total_pages || 1;
    
    if (insights.length === 0) {
      insightsGrid.innerHTML = `
        <div class="insights-empty">
          <p>No insights found matching your criteria.</p>
          <p>Try adjusting your filters or check back later.</p>
        </div>
      `;
      document.getElementById('pagination').innerHTML = '';
      return;
    }
    
    let insightsHTML = '';
    
    insights.forEach(insight => {
      // Truncate content for preview
      const previewContent = insight.content.length > 150 
        ? insight.content.substring(0, 150) + '...' 
        : insight.content;
      
      insightsHTML += `
        <div class="insight-card" data-id="${insight.id}">
          <div class="insight-header">
            <div class="insight-meta">
              <span class="insight-industry">${insight.industry}</span>
              <div class="insight-rating">
                <i class="fas fa-star"></i>
                <span>${insight.average_rating || 'N/A'}</span>
              </div>
            </div>
            <h3 class="insight-title">${insight.title || 'Untitled Insight'}</h3>
            <div class="insight-type">${insight.type}</div>
          </div>
          <div class="insight-body">
            <div class="insight-preview">${previewContent}</div>
          </div>
          <div class="insight-footer">
            <div class="insight-business">${insight.business_name}</div>
            <div class="insight-price">${insight.price || 10} BPT</div>
          </div>
        </div>
      `;
    });
    
    insightsGrid.innerHTML = insightsHTML;
    
    // Add click event to insight cards
    const insightCards = document.querySelectorAll('.insight-card');
    insightCards.forEach(card => {
      card.addEventListener('click', function() {
        const insightId = this.getAttribute('data-id');
        openInsightDetail(insightId);
      });
    });
    
    // Generate pagination
    generatePagination(page, totalPages);
  } catch (error) {
    console.error('Error loading insights:', error);
    document.getElementById('insights-grid').innerHTML = `
      <div class="insights-empty">
        <p>Failed to load insights. Please try again later.</p>
      </div>
    `;
    document.getElementById('pagination').innerHTML = '';
  }
}

// Generate pagination controls
function generatePagination(currentPage, totalPages) {
  currentPage = parseInt(currentPage);
  const pagination = document.getElementById('pagination');
  
  if (totalPages <= 1) {
    pagination.innerHTML = '';
    return;
  }
  
  let paginationHTML = '';
  
  // Previous button
  paginationHTML += `
    <div class="pagination-item ${currentPage === 1 ? 'disabled' : ''}" 
         data-page="${currentPage - 1}" ${currentPage === 1 ? 'disabled' : ''}>
      <i class="fas fa-chevron-left"></i>
    </div>
  `;
  
  // Page numbers
  const maxVisiblePages = 5;
  let startPage = Math.max(1, currentPage - Math.floor(maxVisiblePages / 2));
  let endPage = Math.min(totalPages, startPage + maxVisiblePages - 1);
  
  if (endPage - startPage + 1 < maxVisiblePages) {
    startPage = Math.max(1, endPage - maxVisiblePages + 1);
  }
  
  // First page
  if (startPage > 1) {
    paginationHTML += `
      <div class="pagination-item" data-page="1">1</div>
    `;
    
    if (startPage > 2) {
      paginationHTML += `<div class="pagination-ellipsis">...</div>`;
    }
  }
  
  // Page numbers
  for (let i = startPage; i <= endPage; i++) {
    paginationHTML += `
      <div class="pagination-item ${i === currentPage ? 'active' : ''}" data-page="${i}">
        ${i}
      </div>
    `;
  }
  
  // Last page
  if (endPage < totalPages) {
    if (endPage < totalPages - 1) {
      paginationHTML += `<div class="pagination-ellipsis">...</div>`;
    }
    
    paginationHTML += `
      <div class="pagination-item" data-page="${totalPages}">${totalPages}</div>
    `;
  }
  
  // Next button
  paginationHTML += `
    <div class="pagination-item ${currentPage === totalPages ? 'disabled' : ''}" 
         data-page="${currentPage + 1}" ${currentPage === totalPages ? 'disabled' : ''}>
      <i class="fas fa-chevron-right"></i>
    </div>
  `;
  
  pagination.innerHTML = paginationHTML;
  
  // Add click events to pagination items
  const paginationItems = document.querySelectorAll('.pagination-item:not(.disabled)');
  paginationItems.forEach(item => {
    item.addEventListener('click', function() {
      const page = this.getAttribute('data-page');
      loadInsights(page);
      
      // Scroll to top of insights
      document.querySelector('.marketplace-content').scrollIntoView({ behavior: 'smooth' });
    });
  });
}

// Open insight detail modal
async function openInsightDetail(insightId) {
  try {
    const modal = document.getElementById('insight-detail-modal');
    const modalContent = document.getElementById('insight-detail-content');
    
    // Show modal with loading spinner
    modal.style.display = 'flex';
    modalContent.innerHTML = '<div class="loading-spinner">Loading insight details...</div>';
    
    // Fetch insight details
    const response = await fetch(`/api/insights/${insightId}`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    
    if (!response.ok) {
      throw new Error('Failed to load insight details');
    }
    
    const insight = await response.json();
    
    // Check if user has purchased this insight
    const purchaseResponse = await fetch(`/api/insights/${insightId}/purchased`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    
    const isPurchased = purchaseResponse.ok;
    
    // Format content based on purchase status
    let contentClass = isPurchased ? '' : 'blurred';
    let overlayHTML = '';
    
    if (!isPurchased) {
      overlayHTML = `
        <div class="insight-detail-overlay">
          <h3>Purchase this insight to view full content</h3>
          <p>Unlock valuable market information for your business</p>
        </div>
      `;
    }
    
    // Update modal content
    modalContent.innerHTML = `
      <div class="insight-detail">
        <div class="insight-detail-header">
          <div class="insight-detail-meta">
            <span class="insight-detail-tag">
              <i class="fas fa-industry"></i> ${insight.industry}
            </span>
            <span class="insight-detail-tag">
              <i class="fas fa-tag"></i> ${insight.type}
            </span>
            <span class="insight-detail-tag">
              <i class="fas fa-calendar"></i> ${new Date(insight.created_at).toLocaleDateString()}
            </span>
          </div>
          <h2 class="insight-detail-title">${insight.title || 'Untitled Insight'}</h2>
          <div class="insight-detail-business">
            <div class="business-avatar">
              <i class="fas fa-building"></i>
            </div>
            <div>
              <div class="business-name">${insight.business_name}</div>
              <div class="business-industry">${insight.business_industry}</div>
            </div>
          </div>
        </div>
        <div class="insight-detail-content ${contentClass}" style="position: relative;">
          ${insight.content}
          ${overlayHTML}
        </div>
      </div>
    `;
    
    // Update purchase button
    const purchaseBtn = document.getElementById('purchase-insight-btn');
    const feedbackSection = document.getElementById('feedback-section');
    
    if (isPurchased) {
      purchaseBtn.style.display = 'none';
      feedbackSection.style.display = 'block';
      
      // Load user's existing feedback if any
      await loadUserFeedback(insightId);
    } else {
      purchaseBtn.style.display = 'block';
      purchaseBtn.textContent = `Purchase Insight (${insight.price || 10} BPT)`;
      purchaseBtn.setAttribute('data-id', insightId);
      purchaseBtn.setAttribute('data-price', insight.price || 10);
      feedbackSection.style.display = 'none';
      
      // Add click event to purchase button
      purchaseBtn.addEventListener('click', function() {
        purchaseInsight(insightId, insight.price || 10);
      });
    }
  } catch (error) {
    console.error('Error loading insight details:', error);
    document.getElementById('insight-detail-content').innerHTML = `
      <div class="error-message">
        <p>Failed to load insight details. Please try again later.</p>
      </div>
    `;
  }
}

// Purchase insight
async function purchaseInsight(insightId, price) {
  try {
    const response = await fetch(`/api/insights/${insightId}/purchase`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({
        price: price
      })
    });
    
    if (!response.ok) {
      throw new Error('Failed to purchase insight');
    }
    
    showNotification('Insight purchased successfully', 'success');
    
    // Reload insight details
    openInsightDetail(insightId);
  } catch (error) {
    console.error('Error purchasing insight:', error);
    showNotification('Failed to purchase insight', 'error');
  }
}

// Load user's feedback for an insight
async function loadUserFeedback(insightId) {
  try {
    const response = await fetch(`/api/insights/${insightId}/feedback`, {
      method: 'GET',
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    
    if (response.status === 404) {
      // No feedback yet
      resetFeedbackForm();
      return;
    }
    
    if (!response.ok) {
      throw new Error('Failed to load feedback');
    }
    
    const feedback = await response.json();
    
    // Update rating stars
    const stars = document.querySelectorAll('.rating-stars i');
    stars.forEach((star, index) => {
      if (index < feedback.rating) {
        star.classList.remove('far');
        star.classList.add('fas');
        star.classList.add('active');
      } else {
        star.classList.add('far');
        star.classList.remove('fas');
        star.classList.remove('active');
      }
    });
    
    // Update comments
    document.getElementById('feedback-comments').value = feedback.comments || '';
    
    // Update submit button text
    document.getElementById('submit-feedback-btn').textContent = 'Update Feedback';
  } catch (error) {
    console.error('Error loading feedback:', error);
    resetFeedbackForm();
  }
  
  // Add event listeners for feedback form
  setupFeedbackForm(insightId);
}

// Reset feedback form
function resetFeedbackForm() {
  // Reset stars
  const stars = document.querySelectorAll('.rating-stars i');
  stars.forEach(star => {
    star.classList.add('far');
    star.classList.remove('fas');
    star.classList.remove('active');
  });
  
  // Reset comments
  document.getElementById('feedback-comments').value = '';
  
  // Reset button text
  document.getElementById('submit-feedback-btn').textContent = 'Submit Feedback';
}

// Setup feedback form event listeners
function setupFeedbackForm(insightId) {
  // Star rating
  const stars = document.querySelectorAll('.rating-stars i');
  stars.forEach(star => {
    star.addEventListener('click', function() {
      const rating = parseInt(this.getAttribute('data-rating'));
      
      // Update stars
      stars.forEach((s, index) => {
        if (index < rating) {
          s.classList.remove('far');
          s.classList.add('fas');
          s.classList.add('active');
        } else {
          s.classList.add('far');
          s.classList.remove('fas');
          s.classList.remove('active');
        }
      });
    });
  });
  
  // Submit button
  const submitBtn = document.getElementById('submit-feedback-btn');
  submitBtn.addEventListener('click', function() {
    // Get rating
    const activeStars = document.querySelectorAll('.rating-stars i.active');
    const rating = activeStars.length;
    
    if (rating === 0) {
      showNotification('Please select a rating', 'error');
      return;
    }
    
    // Get comments
    const comments = document.getElementById('feedback-comments').value;
    
    // Submit feedback
    submitFeedback(insightId, rating, comments);
  });
}

// Submit feedback
async function submitFeedback(insightId, rating, comments) {
  try {
    const response = await fetch(`/api/insights/${insightId}/feedback`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      },
      body: JSON.stringify({
        rating: rating,
        comments: comments
      })
    });
    
    if (!response.ok) {
      throw new Error('Failed to submit feedback');
    }
    
    showNotification('Feedback submitted successfully', 'success');
  } catch (error) {
    console.error('Error submitting feedback:', error);
    showNotification('Failed to submit feedback', 'error');
  }
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
